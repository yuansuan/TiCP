package test

import (
	"fmt"
	"io"
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
	writeAtAPI "github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/writeAt"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao/model"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/filemode"
	v20230530 "github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/local/v20230530"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/common/mock"
	"xorm.io/xorm"
)

type WriteAt struct {
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

func (write *WriteAt) SetupSuite() {
	write.T().Log("setup suite")
	logger, err := logging.NewLogger(logging.WithReleaseLevel(logging.DevelopmentLevel))
	if err != nil {
		panic(fmt.Sprintf("init logger failed: %v", err))
	}

	ctrl := gomock.NewController(write.T())
	mockStorageQuotaDao := dao.NewMockStorageQuotaDao(ctrl)
	write.mockStorageQuotaDao = mockStorageQuotaDao
	write.mockStorageQuotaDao.EXPECT().Get(gomock.Any(), gomock.Any()).Return(true, &model.StorageQuota{StorageLimit: 1000}, nil).AnyTimes()
	mockStorageOperationLogDao := dao.NewMockStorageOperationLogDao(ctrl)
	write.mockStorageOperationLogDao = mockStorageOperationLogDao
	write.mockStorageOperationLogDao.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mockEngine, _ := mock.Engine()
	mockEngine.SetLogger(middleware.NewXormLogger(write.logger, true))
	write.engine = mockEngine

	write.logger = logger
	currentDir, storage := CreateLocalStorage(write.userID, v20230530.SoftLink, write.mockStorageQuotaDao, write.mockStorageOperationLogDao, write.engine)
	write.storage = storage

	write.userID = "4TiSBX39DtN"
	write.testDir = "/test-write"

	tmpDir := currentDir + "/" + write.userID + write.testDir
	err = os.MkdirAll(tmpDir, filemode.Directory)
	if err != nil {
		panic(fmt.Sprintf("create temp dir failed: %v", err))
	}
	write.tmpDir = tmpDir
	filePath := filepath.Join(write.tmpDir, "file-1")
	err = os.WriteFile(filePath, nil, 0644)
	if err != nil {
		panic(fmt.Sprintf("Failed to create file: %v", err))
	}
	write.T().Log("File created:", filePath)

}

func (write *WriteAt) TearDownSuite() {
	write.T().Log("teardown suite")
	_ = os.RemoveAll(filepath.Dir(write.tmpDir))
}

func (write *WriteAt) SetupTest() {
	write.T().Log("setup test")
	recorder := httptest.NewRecorder()
	ctx := GetTestGinContext(recorder)
	ctx.Set(logging.LoggerName, write.logger)
	write.ctx = ctx
	write.recorder = recorder
}

func (write *WriteAt) TearDownTest() {
	write.T().Log("teardown test")
}

func TestWriteAt(t *testing.T) {
	suite.Run(t, new(WriteAt))
}

func (write *WriteAt) TestWriteAtSuccess() {

	req := writeAtAPI.Request{
		Path:       "/" + write.userID + write.testDir + "/file-1",
		Compressor: common.NONE,
		Length:     100,
	}
	err, reader := GenerateRandomData(5)
	if err != nil {
		panic(err)
	}
	BindRequest(write.ctx, req)
	write.ctx.Request.Body = io.NopCloser(reader)
	write.storage.WriteAt(write.ctx)

	resp := &writeAtAPI.Response{}
	err = ParseResp(write.recorder.Result(), resp)
	if !assert.Nil(write.T(), err) {
		write.T().Logf("%v: %v", time.Now().Format("2006-01-02 15:04:05"), err)
	} else {
		spew.Dump(resp)
	}
}

// ------------------------ error case ------------------------

func (write *WriteAt) TestWriteAtInvalidPath() {

	req := writeAtAPI.Request{
		Path:       write.userID + write.testDir + "/file-1",
		Compressor: common.NONE,
		Length:     100,
	}
	BindRequest(write.ctx, req)
	err, reader := GenerateRandomData(5)
	if err != nil {
		panic(err)
	}
	write.ctx.Request.Body = io.NopCloser(reader)
	write.storage.WriteAt(write.ctx)

	resp := &writeAtAPI.Response{}
	_ = ParseResp(write.recorder.Result(), resp)
	if assert.NotNil(write.T(), resp) {
		assert.Equal(write.T(), write.recorder.Code, http.StatusBadRequest)
		assert.Equal(write.T(), resp.ErrorCode, common.InvalidPath)
	}
}

func (write *WriteAt) TestWriteAtInvalidOffset() {
	req := writeAtAPI.Request{
		Path:       "/" + write.userID + write.testDir + "/file-1",
		Compressor: common.NONE,
		Length:     100,
		Offset:     -1,
	}
	BindRequest(write.ctx, req)
	err, reader := GenerateRandomData(5)
	if err != nil {
		panic(err)
	}
	write.ctx.Request.Body = io.NopCloser(reader)
	write.storage.WriteAt(write.ctx)

	resp := &writeAtAPI.Response{}
	_ = ParseResp(write.recorder.Result(), resp)
	if assert.NotNil(write.T(), resp) {
		assert.Equal(write.T(), write.recorder.Code, http.StatusBadRequest)
		assert.Equal(write.T(), resp.ErrorCode, common.InvalidOffset)
	}
}

func (write *WriteAt) TestWriteAtInvalidLength() {
	req := writeAtAPI.Request{
		Path:       "/" + write.userID + write.testDir + "/file-1",
		Compressor: common.NONE,
		Length:     -1,
	}
	BindRequest(write.ctx, req)
	err, reader := GenerateRandomData(5)
	if err != nil {
		panic(err)
	}
	write.ctx.Request.Body = io.NopCloser(reader)
	write.storage.WriteAt(write.ctx)

	resp := &writeAtAPI.Response{}
	_ = ParseResp(write.recorder.Result(), resp)
	if assert.NotNil(write.T(), resp) {
		assert.Equal(write.T(), write.recorder.Code, http.StatusBadRequest)
		assert.Equal(write.T(), resp.ErrorCode, common.InvalidLength)
	}
}

func (write *WriteAt) TestWriteAtPathNotFound() {
	req := writeAtAPI.Request{
		Path:       "/" + write.userID + write.testDir + "/not-exist",
		Compressor: common.NONE,
		Length:     100,
	}
	BindRequest(write.ctx, req)
	err, reader := GenerateRandomData(5)
	if err != nil {
		panic(err)
	}
	write.ctx.Request.Body = io.NopCloser(reader)
	write.storage.WriteAt(write.ctx)

	resp := &writeAtAPI.Response{}
	_ = ParseResp(write.recorder.Result(), resp)
	if assert.NotNil(write.T(), resp) {
		assert.Equal(write.T(), write.recorder.Code, http.StatusNotFound)
		assert.Equal(write.T(), resp.ErrorCode, common.PathNotFound)
	}
}

func (write *WriteAt) TestWriteAtPathIsDir() {
	req := writeAtAPI.Request{
		Path:       "/" + write.userID + write.testDir,
		Compressor: common.NONE,
		Length:     100,
	}
	BindRequest(write.ctx, req)
	err, reader := GenerateRandomData(5)
	if err != nil {
		panic(err)
	}
	write.ctx.Request.Body = io.NopCloser(reader)
	write.storage.WriteAt(write.ctx)

	resp := &writeAtAPI.Response{}
	_ = ParseResp(write.recorder.Result(), resp)
	if assert.NotNil(write.T(), resp) {
		assert.Equal(write.T(), write.recorder.Code, http.StatusBadRequest)
		assert.Equal(write.T(), resp.ErrorCode, common.InvalidPath)
	}
}

func (write *WriteAt) TestWriteAtPathInvalidCompressor() {
	req := writeAtAPI.Request{
		Path:       "/" + write.userID + write.testDir + "/file-1",
		Compressor: "fake",
		Length:     100,
	}
	BindRequest(write.ctx, req)
	err, reader := GenerateRandomData(5)
	if err != nil {
		panic(err)
	}
	write.ctx.Request.Body = io.NopCloser(reader)
	write.storage.WriteAt(write.ctx)

	resp := &writeAtAPI.Response{}
	_ = ParseResp(write.recorder.Result(), resp)
	if assert.NotNil(write.T(), resp) {
		assert.Equal(write.T(), write.recorder.Code, http.StatusBadRequest)
		assert.Equal(write.T(), resp.ErrorCode, common.InvalidCompressor)
	}
}
