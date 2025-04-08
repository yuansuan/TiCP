package test

import (
	"bytes"
	"fmt"
	"github.com/yuansuan/ticp/common/project-root-api/common"
	createAPI "github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/create"
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
	truncateAPI "github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/truncate"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao/model"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/filemode"
	v20230530 "github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/local/v20230530"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/common/mock"
	"xorm.io/xorm"
)

const ShortContent = "Hello world!"

type Truncate struct {
	suite.Suite
	logger                     *logging.Logger
	storage                    *v20230530.Storage
	tmpDir                     string
	userID                     string
	testDir                    string
	file1Path                  string
	ctx                        *gin.Context
	recorder                   *httptest.ResponseRecorder
	engine                     *xorm.Engine
	mockStorageQuotaDao        *dao.MockStorageQuotaDao
	mockStorageOperationLogDao *dao.MockStorageOperationLogDao
}

func (truncate *Truncate) SetupSuite() {
	truncate.T().Log("setup suite")
	logger, err := logging.NewLogger(logging.WithReleaseLevel(logging.DevelopmentLevel))
	if err != nil {
		panic(fmt.Sprintf("init logger failed: %v", err))
	}

	ctrl := gomock.NewController(truncate.T())
	mockStorageQuotaDao := dao.NewMockStorageQuotaDao(ctrl)
	truncate.mockStorageQuotaDao = mockStorageQuotaDao
	truncate.mockStorageQuotaDao.EXPECT().Get(gomock.Any(), gomock.Any()).Return(true, &model.StorageQuota{StorageLimit: 1000}, nil).AnyTimes()
	mockStorageOperationLogDao := dao.NewMockStorageOperationLogDao(ctrl)
	truncate.mockStorageOperationLogDao = mockStorageOperationLogDao
	truncate.mockStorageOperationLogDao.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mockEngine, _ := mock.Engine()
	mockEngine.SetLogger(middleware.NewXormLogger(truncate.logger, true))
	truncate.engine = mockEngine

	truncate.logger = logger
	currentDir, storage := CreateLocalStorage(truncate.userID, v20230530.SoftLink, truncate.mockStorageQuotaDao, truncate.mockStorageOperationLogDao, truncate.engine)
	truncate.storage = storage

	truncate.userID = "4TiSBX39DtN"
	truncate.testDir = "/test-truncate"

	tmpDir := currentDir + "/" + truncate.userID + truncate.testDir
	err = os.MkdirAll(tmpDir, filemode.Directory)
	if err != nil {
		panic(fmt.Sprintf("truncate temp dir failed: %v", err))
	}
	truncate.tmpDir = tmpDir
	truncate.file1Path = filepath.Join(truncate.tmpDir, "file-1")
	err = os.WriteFile(truncate.file1Path, []byte(ShortContent), 0644)
	if err != nil {
		panic(fmt.Sprintf("Failed to truncate file: %v", err))
	}
	truncate.T().Log("File created:", truncate.file1Path)

}

func (truncate *Truncate) TearDownSuite() {
	truncate.T().Log("teardown suite")
	_ = os.RemoveAll(filepath.Dir(truncate.tmpDir))
}

func (truncate *Truncate) SetupTest() {
	truncate.T().Log("setup test")
	recorder := httptest.NewRecorder()
	ctx := GetTestGinContext(recorder)
	ctx.Set(logging.LoggerName, truncate.logger)
	truncate.ctx = ctx
	truncate.recorder = recorder
}

func (truncate *Truncate) TearDownTest() {
	truncate.T().Log("teardown test")
}

func TestTruncate(t *testing.T) {
	suite.Run(t, new(Truncate))
}

func (truncate *Truncate) TestTruncateNotCreateSuccess() {
	req := truncateAPI.Request{
		Path: "/" + truncate.userID + truncate.testDir + "/file",
		Size: 1000,
	}
	err := BindJsonRequest(truncate.ctx, req)
	if err != nil {
		panic(fmt.Sprintf("bind request failed: %v", err))
	}
	truncate.storage.Truncate(truncate.ctx)

	resp := &createAPI.Response{}
	err = ParseResp(truncate.recorder.Result(), resp)
	if !assert.Nil(truncate.T(), err) {
		truncate.T().Logf("%v: %v", time.Now().Format("2006-01-02 15:04:05"), err)
	} else {
		spew.Dump(resp.Response)
	}
}

func (truncate *Truncate) TestTruncateCreateSuccess() {
	req := truncateAPI.Request{
		Path:              "/" + truncate.userID + truncate.testDir + "/file",
		Size:              1000,
		CreateIfNotExists: true,
	}
	err := BindJsonRequest(truncate.ctx, req)
	if err != nil {
		panic(fmt.Sprintf("bind request failed: %v", err))
	}
	truncate.storage.Truncate(truncate.ctx)

	resp := &createAPI.Response{}
	err = ParseResp(truncate.recorder.Result(), resp)
	if !assert.Nil(truncate.T(), err) {
		truncate.T().Logf("%v: %v", time.Now().Format("2006-01-02 15:04:05"), err)
	} else {
		spew.Dump(resp.Response)
	}
}

// ------------------------ error case ------------------------

func (truncate *Truncate) TestTruncateSuccessInvalidPath() {
	req := truncateAPI.Request{
		Path: truncate.userID + truncate.testDir + "/file-1",
	}
	err := BindJsonRequest(truncate.ctx, req)
	if err != nil {
		panic(fmt.Sprintf("bind request failed: %v", err))
	}
	truncate.storage.Truncate(truncate.ctx)

	resp := &truncateAPI.Response{}
	_ = ParseResp(truncate.recorder.Result(), resp)
	if assert.NotNil(truncate.T(), resp) {
		assert.Equal(truncate.T(), truncate.recorder.Code, http.StatusBadRequest)
		assert.Equal(truncate.T(), resp.ErrorCode, common.InvalidPath)
	}
}

func (truncate *Truncate) TestTruncateSuccessInvalidSize() {
	req := truncateAPI.Request{
		Path: "/" + truncate.userID + truncate.testDir + "/file-1",
		Size: -1,
	}
	err := BindJsonRequest(truncate.ctx, req)
	if err != nil {
		panic(fmt.Sprintf("bind request failed: %v", err))
	}
	truncate.storage.Truncate(truncate.ctx)

	resp := &truncateAPI.Response{}
	_ = ParseResp(truncate.recorder.Result(), resp)
	if assert.NotNil(truncate.T(), resp) {
		assert.Equal(truncate.T(), truncate.recorder.Code, http.StatusBadRequest)
		assert.Equal(truncate.T(), resp.ErrorCode, common.InvalidSize)
	}
}

func (truncate *Truncate) TestTruncateSuccessAppend() {
	req := truncateAPI.Request{
		Path: "/" + truncate.userID + truncate.testDir + "/file-1",
		Size: 100,
	}
	err := BindJsonRequest(truncate.ctx, req)
	if err != nil {
		panic(fmt.Sprintf("bind request failed: %v", err))
	}
	truncate.storage.Truncate(truncate.ctx)

	resp := &truncateAPI.Response{}
	_ = ParseResp(truncate.recorder.Result(), resp)
	if assert.NotNil(truncate.T(), resp) {
		assert.Equal(truncate.T(), http.StatusOK, truncate.recorder.Code)
	}

	buf, err := os.ReadFile(truncate.file1Path)
	if err != nil {
		panic(fmt.Sprintf("read file failed: %v", err))
	}
	assert.Equal(truncate.T(), 0, bytes.Compare(buf[:len(ShortContent)], []byte(ShortContent)))
}
