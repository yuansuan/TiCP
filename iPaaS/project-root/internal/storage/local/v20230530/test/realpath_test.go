package test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/middleware"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/common/project-root-api/common"
	realpathAPI "github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/realpath"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao/model"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/filemode"
	v20230530 "github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/local/v20230530"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/common/mock"
	"xorm.io/xorm"
)

type RealPath struct {
	suite.Suite
	logger              *logging.Logger
	storage             *v20230530.Storage
	tmpDir              string
	userID              string
	testDir             string
	ctx                 *gin.Context
	recorder            *httptest.ResponseRecorder
	engine              *xorm.Engine
	mockStorageQuotaDao *dao.MockStorageQuotaDao
}

func (realpath *RealPath) SetupSuite() {
	realpath.T().Log("setup suite")
	logger, err := logging.NewLogger(logging.WithReleaseLevel(logging.DevelopmentLevel))
	if err != nil {
		panic(fmt.Sprintf("init logger failed: %v", err))
	}

	ctrl := gomock.NewController(realpath.T())
	mockStorageQuotaDao := dao.NewMockStorageQuotaDao(ctrl)
	realpath.mockStorageQuotaDao = mockStorageQuotaDao
	realpath.mockStorageQuotaDao.EXPECT().Get(gomock.Any(), gomock.Any()).Return(true, &model.StorageQuota{StorageLimit: 1000}, nil).AnyTimes()
	mockEngine, _ := mock.Engine()
	mockEngine.SetLogger(middleware.NewXormLogger(realpath.logger, true))
	realpath.engine = mockEngine

	realpath.logger = logger
	currentDir, storage := CreateLocalStorage(realpath.userID, v20230530.SoftLink, realpath.mockStorageQuotaDao, nil, realpath.engine)
	realpath.storage = storage

	realpath.userID = "4TiSBX39DtN"
	realpath.testDir = "/test-realpath"

	tmpDir := currentDir + "/" + realpath.userID + realpath.testDir
	err = os.MkdirAll(tmpDir, filemode.Directory)
	if err != nil {
		panic(fmt.Sprintf("create temp dir failed: %v", err))
	}
	realpath.tmpDir = tmpDir
	filePath := filepath.Join(realpath.tmpDir, "file-1")
	err = os.WriteFile(filePath, nil, 0644)
	if err != nil {
		panic(fmt.Sprintf("Failed to create file: %v", err))
	}
	realpath.T().Log("File created:", filePath)
}

func (realpath *RealPath) TearDownSuite() {
	realpath.T().Log("teardown suite")
	_ = os.RemoveAll(filepath.Dir(realpath.tmpDir))
}

func (realpath *RealPath) SetupTest() {
	realpath.T().Log("setup test")
	recorder := httptest.NewRecorder()
	ctx := GetTestGinContext(recorder)
	ctx.Set(logging.LoggerName, realpath.logger)
	realpath.ctx = ctx
	realpath.recorder = recorder
}

func (realpath *RealPath) TearDownTest() {
	realpath.T().Log("teardown test")
}

func TestRealpath(t *testing.T) {
	suite.Run(t, new(RealPath))
}

func (realpath *RealPath) TestRealpathSuccess() {
	req := realpathAPI.Request{
		RelativePath: "/" + realpath.userID + realpath.testDir,
	}
	BindRequest(realpath.ctx, req)
	realpath.storage.Realpath(realpath.ctx)

	resp := &realpathAPI.Response{}
	err := ParseResp(realpath.recorder.Result(), resp)
	if !assert.Nil(realpath.T(), err) {
		realpath.T().Logf("%v: %v", time.Now().Format("2006-01-02 15:04:05"), err)
	} else {
		spew.Dump(resp.Data)
	}
}

// ------------------------ error case ------------------------

func (realpath *RealPath) TestRealpathInvalidPath() {
	req := realpathAPI.Request{
		RelativePath: realpath.userID + realpath.testDir,
	}
	BindRequest(realpath.ctx, req)
	realpath.storage.Realpath(realpath.ctx)

	resp := &realpathAPI.Response{}
	_ = ParseResp(realpath.recorder.Result(), resp)
	if assert.NotNil(realpath.T(), resp) {
		assert.Equal(realpath.T(), realpath.recorder.Code, http.StatusBadRequest)
		assert.Equal(realpath.T(), resp.ErrorCode, common.InvalidPath)
	}
}
