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
	statAPI "github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/stat"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao/model"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/filemode"
	v20230530 "github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/local/v20230530"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/local/v20230530/linker"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/common/mock"
	"xorm.io/xorm"
)

type Stat struct {
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
	testLinkDir         string
}

func (stat *Stat) SetupSuite() {
	stat.T().Log("setup suite")
	logger, err := logging.NewLogger(logging.WithReleaseLevel(logging.DevelopmentLevel))
	if err != nil {
		panic(fmt.Sprintf("init logger failed: %v", err))
	}
	ctrl := gomock.NewController(stat.T())
	mockStorageQuotaDao := dao.NewMockStorageQuotaDao(ctrl)
	stat.mockStorageQuotaDao = mockStorageQuotaDao
	stat.mockStorageQuotaDao.EXPECT().Get(gomock.Any(), gomock.Any()).Return(true, &model.StorageQuota{StorageLimit: 1000}, nil).AnyTimes()
	mockEngine, _ := mock.Engine()
	mockEngine.SetLogger(middleware.NewXormLogger(stat.logger, true))
	stat.engine = mockEngine

	stat.logger = logger
	currentDir, storage := CreateLocalStorage(stat.userID, v20230530.SoftLink, stat.mockStorageQuotaDao, nil, stat.engine)
	stat.storage = storage

	stat.userID = "4TiSBX39DtN"
	stat.testDir = "/test-stat"
	stat.testLinkDir = "/test-link"

	tmpDir := currentDir + "/" + stat.userID + stat.testDir
	linkDir := currentDir + "/" + stat.userID + stat.testLinkDir
	err = os.MkdirAll(tmpDir, filemode.Directory)
	if err != nil {
		panic(fmt.Sprintf("create temp dir failed: %v", err))
	}
	stat.tmpDir = tmpDir
	filePath := filepath.Join(stat.tmpDir, "file-1")
	err = os.WriteFile(filePath, nil, 0644)
	if err != nil {
		panic(fmt.Sprintf("Failed to create file: %v", err))
	}
	stat.T().Log("File created:", filePath)

	// softlink dir
	sl := &linker.SoftLink{}
	err = sl.Link(tmpDir, linkDir)
}

func (stat *Stat) TearDownSuite() {
	stat.T().Log("teardown suite")
	_ = os.RemoveAll(filepath.Dir(stat.tmpDir))
}

func (stat *Stat) SetupTest() {
	stat.T().Log("setup test")
	recorder := httptest.NewRecorder()
	ctx := GetTestGinContext(recorder)
	ctx.Set(logging.LoggerName, stat.logger)
	stat.ctx = ctx
	stat.recorder = recorder
}

func (stat *Stat) TearDownTest() {
	stat.T().Log("teardown test")
}

func TestStat(t *testing.T) {
	suite.Run(t, new(Stat))
}

func (stat *Stat) TestStatSuccessDir() {
	req := statAPI.Request{
		Path: "/" + stat.userID + stat.testDir,
	}
	BindRequest(stat.ctx, req)
	stat.storage.Stat(stat.ctx)

	resp := &statAPI.Response{}
	err := ParseResp(stat.recorder.Result(), resp)
	if !assert.Nil(stat.T(), err) {
		stat.T().Logf("%v: %v", time.Now().Format("2006-01-02 15:04:05"), err)
	} else {
		spew.Dump(resp.Data)
	}
}

func (stat *Stat) TestStatSuccessLinkDir() {
	req := statAPI.Request{
		Path: "/" + stat.userID + stat.testLinkDir,
	}
	BindRequest(stat.ctx, req)
	stat.storage.Stat(stat.ctx)

	resp := &statAPI.Response{}
	err := ParseResp(stat.recorder.Result(), resp)
	if !assert.Nil(stat.T(), err) {
		stat.T().Logf("%v: %v", time.Now().Format("2006-01-02 15:04:05"), err)
	} else {
		spew.Dump(resp.Data)
	}
}

func (stat *Stat) TestStatSuccessFile() {
	req := statAPI.Request{
		Path: "/" + stat.userID + stat.testDir + "/file-1",
	}
	BindRequest(stat.ctx, req)
	stat.storage.Stat(stat.ctx)

	resp := &statAPI.Response{}
	err := ParseResp(stat.recorder.Result(), resp)
	if !assert.Nil(stat.T(), err) {
		stat.T().Logf("%v: %v", time.Now().Format("2006-01-02 15:04:05"), err)
	} else {
		spew.Dump(resp.Data)
	}
}

func (stat *Stat) TestStatSuccessFileInLinkedDir() {
	req := statAPI.Request{
		Path: "/" + stat.userID + stat.testLinkDir + "/file-1",
	}
	BindRequest(stat.ctx, req)
	stat.storage.Stat(stat.ctx)

	resp := &statAPI.Response{}
	err := ParseResp(stat.recorder.Result(), resp)
	if !assert.Nil(stat.T(), err) {
		stat.T().Logf("%v: %v", time.Now().Format("2006-01-02 15:04:05"), err)
	} else {
		spew.Dump(resp.Data)
	}
}

// ------------------------ error case ------------------------

func (stat *Stat) TestStatInvalidPath() {
	req := statAPI.Request{
		Path: stat.userID + stat.testDir + "/file-1",
	}
	BindRequest(stat.ctx, req)
	stat.storage.Stat(stat.ctx)

	resp := &statAPI.Response{}
	_ = ParseResp(stat.recorder.Result(), resp)
	if assert.NotNil(stat.T(), resp) {
		assert.Equal(stat.T(), stat.recorder.Code, http.StatusBadRequest)
		assert.Equal(stat.T(), resp.ErrorCode, common.InvalidPath)
	}
}

func (stat *Stat) TestStatPathNotFound() {
	req := statAPI.Request{
		Path: "/" + stat.userID + stat.testDir + "/not-exist",
	}
	BindRequest(stat.ctx, req)
	stat.storage.Stat(stat.ctx)

	resp := &statAPI.Response{}
	_ = ParseResp(stat.recorder.Result(), resp)
	if assert.NotNil(stat.T(), resp) {
		assert.Equal(stat.T(), stat.recorder.Code, http.StatusNotFound)
		assert.Equal(stat.T(), resp.ErrorCode, common.PathNotFound)
	}
}
