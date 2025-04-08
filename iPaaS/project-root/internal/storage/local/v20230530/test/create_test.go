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
	createAPI "github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/create"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao/model"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/filemode"
	v20230530 "github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/local/v20230530"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/common/mock"
	"xorm.io/xorm"
)

type Create struct {
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

func (create *Create) SetupSuite() {
	create.T().Log("setup suite")
	logger, err := logging.NewLogger(logging.WithReleaseLevel(logging.DevelopmentLevel))
	if err != nil {
		panic(fmt.Sprintf("init logger failed: %v", err))
	}

	ctrl := gomock.NewController(create.T())
	mockStorageQuotaDao := dao.NewMockStorageQuotaDao(ctrl)
	create.mockStorageQuotaDao = mockStorageQuotaDao
	create.mockStorageQuotaDao.EXPECT().Get(gomock.Any(), gomock.Any()).Return(true, &model.StorageQuota{StorageLimit: 1000}, nil).AnyTimes()
	mockStorageOperationLogDao := dao.NewMockStorageOperationLogDao(ctrl)
	create.mockStorageOperationLogDao = mockStorageOperationLogDao
	create.mockStorageOperationLogDao.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mockEngine, _ := mock.Engine()
	mockEngine.SetLogger(middleware.NewXormLogger(create.logger, true))
	create.engine = mockEngine

	create.logger = logger
	currentDir, storage := CreateLocalStorage(create.userID, v20230530.SoftLink, create.mockStorageQuotaDao, create.mockStorageOperationLogDao, create.engine)
	create.storage = storage

	create.userID = "4TiSBX39DtN"
	create.testDir = "/test-create"

	tmpDir := currentDir + "/" + create.userID + create.testDir
	err = os.MkdirAll(tmpDir, filemode.Directory)
	if err != nil {
		panic(fmt.Sprintf("create temp dir failed: %v", err))
	}
	create.tmpDir = tmpDir
	filePath := filepath.Join(create.tmpDir, "file-1")
	err = os.WriteFile(filePath, nil, 0644)
	if err != nil {
		panic(fmt.Sprintf("Failed to create file: %v", err))
	}
	create.T().Log("File created:", filePath)

}

func (create *Create) TearDownSuite() {
	create.T().Log("teardown suite")
	_ = os.RemoveAll(filepath.Dir(create.tmpDir))
}

func (create *Create) SetupTest() {
	create.T().Log("setup test")
	recorder := httptest.NewRecorder()
	ctx := GetTestGinContext(recorder)
	ctx.Set(logging.LoggerName, create.logger)
	create.ctx = ctx
	create.recorder = recorder
}

func (create *Create) TearDownTest() {
	create.T().Log("teardown test")
}

func TestCreate(t *testing.T) {
	suite.Run(t, new(Create))
}

func (create *Create) TestCreateSuccess() {
	req := createAPI.Request{
		Path: "/" + create.userID + create.testDir + "/file",
		Size: 1000,
	}
	err := BindJsonRequest(create.ctx, req)
	if err != nil {
		panic(fmt.Sprintf("bind request failed: %v", err))
	}
	create.storage.Create(create.ctx)

	resp := &createAPI.Response{}
	err = ParseResp(create.recorder.Result(), resp)
	if !assert.Nil(create.T(), err) {
		create.T().Logf("%v: %v", time.Now().Format("2006-01-02 15:04:05"), err)
	} else {
		spew.Dump(resp.Response)
	}
}

// ------------------------ error case ------------------------

func (create *Create) TestCreateSuccessInvalidPath() {
	req := createAPI.Request{
		Path: create.userID + create.testDir + "/file-1",
	}
	err := BindJsonRequest(create.ctx, req)
	if err != nil {
		panic(fmt.Sprintf("bind request failed: %v", err))
	}
	create.storage.Create(create.ctx)

	resp := &createAPI.Response{}
	_ = ParseResp(create.recorder.Result(), resp)
	if assert.NotNil(create.T(), resp) {
		assert.Equal(create.T(), create.recorder.Code, http.StatusBadRequest)
		assert.Equal(create.T(), resp.ErrorCode, common.InvalidPath)
	}
}

func (create *Create) TestCreateSuccessInvalidSize() {
	req := createAPI.Request{
		Path: "/" + create.userID + create.testDir + "/file-1",
		Size: -1,
	}
	err := BindJsonRequest(create.ctx, req)
	if err != nil {
		panic(fmt.Sprintf("bind request failed: %v", err))
	}
	create.storage.Create(create.ctx)

	resp := &createAPI.Response{}
	_ = ParseResp(create.recorder.Result(), resp)
	if assert.NotNil(create.T(), resp) {
		assert.Equal(create.T(), create.recorder.Code, http.StatusBadRequest)
		assert.Equal(create.T(), resp.ErrorCode, common.InvalidSize)
	}
}

func (create *Create) TestCreateSuccessFileExist() {
	req := createAPI.Request{
		Path: "/" + create.userID + create.testDir + "/file-1",
	}
	err := BindJsonRequest(create.ctx, req)
	if err != nil {
		panic(fmt.Sprintf("bind request failed: %v", err))
	}
	create.storage.Create(create.ctx)

	resp := &createAPI.Response{}
	_ = ParseResp(create.recorder.Result(), resp)
	if assert.NotNil(create.T(), resp) {
		assert.Equal(create.T(), create.recorder.Code, http.StatusBadRequest)
		assert.Equal(create.T(), resp.ErrorCode, common.PathExists)
	}
}
