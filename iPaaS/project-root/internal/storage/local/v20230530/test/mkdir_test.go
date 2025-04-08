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
	mkdirAPI "github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/mkdir"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao/model"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/filemode"
	v20230530 "github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/local/v20230530"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/common/mock"
	"xorm.io/xorm"
)

type Mkdir struct {
	suite.Suite
	logger                     *logging.Logger
	storage                    *v20230530.Storage
	tmpDir                     string
	userID                     string
	testDir                    string
	ctx                        *gin.Context
	recorder                   *httptest.ResponseRecorder
	engine                     *xorm.Engine
	mockStorageQuotaDao        *dao.MockStorageQuotaDao
	mockStorageOperationLogDao *dao.MockStorageOperationLogDao
}

func (mk *Mkdir) SetupSuite() {
	mk.T().Log("setup suite")
	logger, err := logging.NewLogger(logging.WithReleaseLevel(logging.DevelopmentLevel))
	if err != nil {
		panic(fmt.Sprintf("init logger failed: %v", err))
	}

	ctrl := gomock.NewController(mk.T())
	mockStorageQuotaDao := dao.NewMockStorageQuotaDao(ctrl)
	mk.mockStorageQuotaDao = mockStorageQuotaDao
	mk.mockStorageQuotaDao.EXPECT().Get(gomock.Any(), gomock.Any()).Return(true, &model.StorageQuota{StorageLimit: 1000}, nil).AnyTimes()
	mockStorageOperationLogDao := dao.NewMockStorageOperationLogDao(ctrl)
	mk.mockStorageOperationLogDao = mockStorageOperationLogDao
	mk.mockStorageOperationLogDao.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mockEngine, _ := mock.Engine()
	mockEngine.SetLogger(middleware.NewXormLogger(mk.logger, true))
	mk.engine = mockEngine

	mk.logger = logger
	currentDir, storage := CreateLocalStorage(mk.userID, v20230530.SoftLink, mk.mockStorageQuotaDao, mk.mockStorageOperationLogDao, mk.engine)
	mk.storage = storage

	mk.userID = "4TiSBX39DtN"
	mk.testDir = "/test-mkdir"

	tmpDir := currentDir + "/" + mk.userID + mk.testDir
	err = os.MkdirAll(tmpDir+"/test-dir", filemode.Directory)
	if err != nil {
		panic(fmt.Sprintf("create temp dir failed: %v", err))
	}
	mk.tmpDir = tmpDir

}

func (mk *Mkdir) TearDownSuite() {
	mk.T().Log("teardown suite")
	_ = os.RemoveAll(filepath.Dir(mk.tmpDir))
}

func (mk *Mkdir) SetupTest() {
	mk.T().Log("setup test")
	recorder := httptest.NewRecorder()
	ctx := GetTestGinContext(recorder)
	ctx.Set(logging.LoggerName, mk.logger)
	mk.ctx = ctx
	mk.recorder = recorder
}

func (mk *Mkdir) TearDownTest() {
	mk.T().Log("teardown test")
}

func TestMkdir(t *testing.T) {
	suite.Run(t, new(Mkdir))
}

func (mk *Mkdir) TestMkdirSuccess() {
	req := mkdirAPI.Request{
		Path: "/" + mk.userID + mk.testDir + "/test",
	}
	err := BindJsonRequest(mk.ctx, req)
	if err != nil {
		panic(fmt.Sprintf("bind request failed: %v", err))
	}
	mk.storage.Mkdir(mk.ctx)

	resp := &mkdirAPI.Response{}
	err = ParseResp(mk.recorder.Result(), resp)
	if !assert.Nil(mk.T(), err) {
		mk.T().Logf("%v: %v", time.Now().Format("2006-01-02 15:04:05"), err)
	} else {
		spew.Dump(resp.Response)
	}
}

func (mk *Mkdir) TestMkdirSuccessIgnoreExist() {
	req := mkdirAPI.Request{
		Path:        "/" + mk.userID + mk.testDir + "/test-dir",
		IgnoreExist: true,
	}
	err := BindJsonRequest(mk.ctx, req)
	if err != nil {
		panic(fmt.Sprintf("bind request failed: %v", err))
	}
	mk.storage.Mkdir(mk.ctx)

	resp := &mkdirAPI.Response{}
	err = ParseResp(mk.recorder.Result(), resp)
	if !assert.Nil(mk.T(), err) {
		mk.T().Logf("%v: %v", time.Now().Format("2006-01-02 15:04:05"), err)
	} else {
		spew.Dump(resp.Response)
	}
}

// ------------------------ error case ------------------------

func (mk *Mkdir) TestMkdirInvalidPath() {
	req := mkdirAPI.Request{
		Path: mk.userID + mk.testDir + "/test-dir",
	}
	err := BindJsonRequest(mk.ctx, req)
	if err != nil {
		panic(fmt.Sprintf("bind request failed: %v", err))
	}
	mk.storage.Mkdir(mk.ctx)

	resp := &mkdirAPI.Response{}
	_ = ParseResp(mk.recorder.Result(), resp)
	if assert.NotNil(mk.T(), resp) {
		assert.Equal(mk.T(), mk.recorder.Code, http.StatusBadRequest)
		assert.Equal(mk.T(), resp.ErrorCode, common.InvalidPath)
	}
}

func (mk *Mkdir) TestMkdirPathExist() {
	req := mkdirAPI.Request{
		Path: "/" + mk.userID + mk.testDir + "/test-dir",
	}
	err := BindJsonRequest(mk.ctx, req)
	if err != nil {
		panic(fmt.Sprintf("bind request failed: %v", err))
	}
	mk.storage.Mkdir(mk.ctx)

	resp := &mkdirAPI.Response{}
	_ = ParseResp(mk.recorder.Result(), resp)
	if assert.NotNil(mk.T(), resp) {
		assert.Equal(mk.T(), mk.recorder.Code, http.StatusBadRequest)
		assert.Equal(mk.T(), resp.ErrorCode, common.PathExists)
	}
}
