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
	remove "github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/rm"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao/model"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/filemode"
	v20230530 "github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/local/v20230530"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/common/mock"
	"xorm.io/xorm"
)

type Rm struct {
	suite.Suite
	logger                     *logging.Logger
	storage                    *v20230530.Storage
	tmpDir                     string
	userID                     string
	testDir                    string
	fileNumber                 int
	ctx                        *gin.Context
	recorder                   *httptest.ResponseRecorder
	engine                     *xorm.Engine
	mockStorageQuotaDao        *dao.MockStorageQuotaDao
	mockStorageOperationLogDao *dao.MockStorageOperationLogDao
}

func (rm *Rm) SetupSuite() {
	rm.T().Log("setup suite")
	logger, err := logging.NewLogger(logging.WithReleaseLevel(logging.DevelopmentLevel))
	if err != nil {
		panic(fmt.Sprintf("init logger failed: %v", err))
	}

	ctrl := gomock.NewController(rm.T())
	mockStorageQuotaDao := dao.NewMockStorageQuotaDao(ctrl)
	rm.mockStorageQuotaDao = mockStorageQuotaDao
	rm.mockStorageQuotaDao.EXPECT().Get(gomock.Any(), gomock.Any()).Return(true, &model.StorageQuota{StorageLimit: 1000}, nil).AnyTimes()
	mockStorageOperationLogDao := dao.NewMockStorageOperationLogDao(ctrl)
	rm.mockStorageOperationLogDao = mockStorageOperationLogDao
	rm.mockStorageOperationLogDao.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mockEngine, _ := mock.Engine()
	mockEngine.SetLogger(middleware.NewXormLogger(rm.logger, true))
	rm.engine = mockEngine

	rm.logger = logger
	currentDir, storage := CreateLocalStorage(rm.userID, v20230530.SoftLink, rm.mockStorageQuotaDao, rm.mockStorageOperationLogDao, rm.engine)
	rm.storage = storage

	rm.userID = "4TiSBX39DtN"
	rm.testDir = "/test-rm"

	tmpDir := currentDir + "/" + rm.userID + rm.testDir
	err = os.MkdirAll(tmpDir, filemode.Directory)
	if err != nil {
		panic(fmt.Sprintf("create temp dir failed: %v", err))
	}
	rm.tmpDir = tmpDir
	filePath := filepath.Join(rm.tmpDir, "file-1")
	err = os.WriteFile(filePath, nil, 0644)
	if err != nil {
		panic(fmt.Sprintf("Failed to create file: %v", err))
	}
	rm.T().Log("File created:", filePath)

}

func (rm *Rm) TearDownSuite() {
	rm.T().Log("teardown suite")
	_ = os.RemoveAll(filepath.Dir(rm.tmpDir))
}

func (rm *Rm) SetupTest() {
	rm.T().Log("setup test")
	recorder := httptest.NewRecorder()
	ctx := GetTestGinContext(recorder)
	ctx.Set(logging.LoggerName, rm.logger)
	rm.ctx = ctx
	rm.recorder = recorder
}

func (rm *Rm) TearDownTest() {
	rm.T().Log("teardown test")
}

func TestRm(t *testing.T) {
	suite.Run(t, new(Rm))
}

func (rm *Rm) TestRmSuccess() {
	req := remove.Request{
		Path: "/" + rm.userID + rm.testDir + "/file-1",
	}
	err := BindJsonRequest(rm.ctx, req)
	if err != nil {
		panic(fmt.Sprintf("bind request failed: %v", err))
	}
	rm.storage.Rm(rm.ctx)

	resp := &remove.Response{}
	err = ParseResp(rm.recorder.Result(), resp)
	if !assert.Nil(rm.T(), err) {
		rm.T().Logf("%v: %v", time.Now().Format("2006-01-02 15:04:05"), err)
	} else {
		spew.Dump(resp.Response)
	}
}

func (rm *Rm) TestRmSuccessNotExist() {
	req := remove.Request{
		Path:           "/" + rm.userID + rm.testDir + "/not-exist",
		IgnoreNotExist: true,
	}
	err := BindJsonRequest(rm.ctx, req)
	if err != nil {
		panic(fmt.Sprintf("bind request failed: %v", err))
	}
	rm.storage.Rm(rm.ctx)

	resp := &remove.Response{}
	err = ParseResp(rm.recorder.Result(), resp)
	if !assert.Nil(rm.T(), err) {
		rm.T().Logf("%v: %v", time.Now().Format("2006-01-02 15:04:05"), err)
	} else {
		spew.Dump(resp.Response)
	}
}

// ------------------------ error case ------------------------

func (rm *Rm) TestRmSuccessInvalidPath() {
	req := remove.Request{
		Path:           rm.userID + rm.testDir + "/file-1",
		IgnoreNotExist: true,
	}
	err := BindJsonRequest(rm.ctx, req)
	if err != nil {
		panic(fmt.Sprintf("bind request failed: %v", err))
	}
	rm.storage.Rm(rm.ctx)

	resp := &remove.Response{}
	_ = ParseResp(rm.recorder.Result(), resp)
	if assert.NotNil(rm.T(), resp) {
		assert.Equal(rm.T(), rm.recorder.Code, http.StatusBadRequest)
		assert.Equal(rm.T(), resp.ErrorCode, common.InvalidPath)
	}
}

func (rm *Rm) TestRmSuccessPathNotFound() {
	req := remove.Request{
		Path: "/" + rm.userID + rm.testDir + "/not-exist",
	}
	err := BindJsonRequest(rm.ctx, req)
	if err != nil {
		panic(fmt.Sprintf("bind request failed: %v", err))
	}
	rm.storage.Rm(rm.ctx)

	resp := &remove.Response{}
	_ = ParseResp(rm.recorder.Result(), resp)
	if assert.NotNil(rm.T(), resp) {
		assert.Equal(rm.T(), rm.recorder.Code, http.StatusNotFound)
		assert.Equal(rm.T(), resp.ErrorCode, common.PathNotFound)
	}
}
