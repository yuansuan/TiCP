package test

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/middleware"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/common/project-root-api/common"
	api "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
	"github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/batchDownload"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/filemode"
	v20230530 "github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/local/v20230530"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/local/v20230530/linker"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/common/mock"
	"xorm.io/xorm"
)

type BatchDownload struct {
	suite.Suite
	logger                     *logging.Logger
	storage                    *v20230530.Storage
	tmpDir                     string
	userID                     string
	testDir                    string
	fileNumber                 int
	DirNumber                  int
	ctx                        *gin.Context
	recorder                   *httptest.ResponseRecorder
	engine                     *xorm.Engine
	mockStorageOperationLogDao *dao.MockStorageOperationLogDao
	testLinkDir                string
}

func (b *BatchDownload) SetupSuite() {
	b.T().Log("setup suite")
	logger, err := logging.NewLogger(logging.WithReleaseLevel(logging.DevelopmentLevel))
	if err != nil {
		panic(fmt.Sprintf("init logger failed: %v", err))
	}
	b.logger = logger
	ctrl := gomock.NewController(b.T())
	mockStorageOperationLogDao := dao.NewMockStorageOperationLogDao(ctrl)
	b.mockStorageOperationLogDao = mockStorageOperationLogDao
	b.mockStorageOperationLogDao.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mockEngine, _ := mock.Engine()
	mockEngine.SetLogger(middleware.NewXormLogger(b.logger, true))
	b.engine = mockEngine
	currentDir, storage := CreateLocalStorage("", v20230530.SoftLink, nil, b.mockStorageOperationLogDao, b.engine)
	b.storage = storage

	b.userID = "4TiSBX39DtN"
	b.testDir = "/test-batch-download"
	b.fileNumber = 5
	b.testLinkDir = "/test-link"

	tmpDir := currentDir + "/" + b.userID + b.testDir
	linkDir := currentDir + "/" + b.userID + b.testLinkDir
	err = os.MkdirAll(tmpDir, filemode.Directory)
	if err != nil {
		panic(fmt.Sprintf("create temp dir failed: %v", err))
	}

	b.tmpDir = tmpDir
	for i := 0; i < b.fileNumber; i++ {
		filePath := filepath.Join(b.tmpDir, "file-"+fmt.Sprintf("%d", i+1))
		data := make([]byte, 10*1024*1024)
		rand.Read(data)
		err := os.WriteFile(filePath, data, 0644)
		if err != nil {
			panic(fmt.Sprintf("Failed to create file: %v", err))
		}
		b.T().Log("File created:", filePath)
	}
	for i := 0; i < b.DirNumber; i++ {
		dirPath := filepath.Join(b.tmpDir, "dir-"+fmt.Sprintf("%d", i+1))
		err := os.Mkdir(dirPath, 0644)
		if err != nil {
			panic(fmt.Sprintf("Failed to create dir: %v", err))
		}
		b.T().Log("Dir created:", dirPath)
	}

	sl := &linker.SoftLink{}
	err = sl.Link(b.tmpDir, linkDir)
}

func (b *BatchDownload) TearDownSuite() {
	b.T().Log("teardown suite")
	_ = os.RemoveAll(filepath.Dir(b.tmpDir))
}

func (b *BatchDownload) SetupTest() {
	b.T().Log("setup test")
	recorder := httptest.NewRecorder()
	ctx := GetTestGinContext(recorder)
	ctx.Set(logging.LoggerName, b.logger)
	b.ctx = ctx
	b.recorder = recorder
}

func (b *BatchDownload) TearDownTest() {
	b.T().Log("teardown test")
}

func TestBatchDownload(t *testing.T) {
	suite.Run(t, new(BatchDownload))
}

func (b *BatchDownload) TestBatchDownloadSuccess() {
	fileName := "test.zip"
	req := batchDownload.Request{
		Paths:    []string{"/" + b.userID + b.testDir},
		FileName: fileName,
		BasePath: "/" + b.userID + "/",
	}
	err := BindJsonRequest(b.ctx, req)
	if err != nil {
		panic(fmt.Sprintf("bind request failed: %v", err))
	}
	b.storage.BatchDownload(b.ctx)
	if assert.Equal(b.T(), 200, b.recorder.Code) {
		assert.Equal(b.T(), "application/zip", b.recorder.Header().Get("Content-Type"))
		assert.Equal(b.T(), fmt.Sprintf("attachment; filename=\"%s\"", fileName), b.recorder.Header().Get("Content-Disposition"))
	}
}

func (b *BatchDownload) TestBatchDownloadLinkDirSuccess() {
	fileName := "test.zip"
	req := batchDownload.Request{
		Paths:    []string{"/" + b.userID + b.testLinkDir},
		FileName: fileName,
		BasePath: "/" + b.userID + "/",
	}
	err := BindJsonRequest(b.ctx, req)
	if err != nil {
		panic(fmt.Sprintf("bind request failed: %v", err))
	}
	b.storage.BatchDownload(b.ctx)
	if assert.Equal(b.T(), 200, b.recorder.Code) {
		assert.Equal(b.T(), "application/zip", b.recorder.Header().Get("Content-Type"))
		assert.Equal(b.T(), fmt.Sprintf("attachment; filename=\"%s\"", fileName), b.recorder.Header().Get("Content-Disposition"))
	}
}

func (b *BatchDownload) TestBatchDownloadSingleFileSuccess() {
	fileName := "file-1"
	req := batchDownload.Request{
		Paths:    []string{"/" + b.userID + b.testDir + "/" + fileName},
		FileName: fileName,
		BasePath: "/" + b.userID + "/",
	}
	err := BindJsonRequest(b.ctx, req)
	if err != nil {
		panic(fmt.Sprintf("bind request failed: %v", err))
	}
	b.storage.BatchDownload(b.ctx)
	if assert.Equal(b.T(), 200, b.recorder.Code) {
		assert.Equal(b.T(), "application/octet-stream", b.recorder.Header().Get("Content-Type"))
		assert.Equal(b.T(), fmt.Sprintf("attachment; filename=\"%s\"", fileName), b.recorder.Header().Get("Content-Disposition"))
	}
}

func (b *BatchDownload) TestBatchDownloadSingleFileInLinkDirSuccess() {
	fileName := "file-1"
	req := batchDownload.Request{
		Paths:    []string{"/" + b.userID + b.testLinkDir + "/" + fileName},
		FileName: fileName,
		BasePath: "/" + b.userID + "/",
	}
	err := BindJsonRequest(b.ctx, req)
	if err != nil {
		panic(fmt.Sprintf("bind request failed: %v", err))
	}
	b.storage.BatchDownload(b.ctx)
	if assert.Equal(b.T(), 200, b.recorder.Code) {
		assert.Equal(b.T(), "application/octet-stream", b.recorder.Header().Get("Content-Type"))
		assert.Equal(b.T(), fmt.Sprintf("attachment; filename=\"%s\"", fileName), b.recorder.Header().Get("Content-Disposition"))
	}
}

// ------------------------ error case ------------------------

func (b *BatchDownload) TestBatchDownloadEmptyPath() {
	fileName := "test.zip"
	req := batchDownload.Request{
		FileName: fileName,
		BasePath: "/" + b.userID + "/",
	}
	err := BindJsonRequest(b.ctx, req)
	if err != nil {
		panic(fmt.Sprintf("bind request failed: %v", err))
	}
	b.storage.BatchDownload(b.ctx)
	resp := &api.Response{}
	_ = ParseResp(b.recorder.Result(), resp)
	if assert.NotNil(b.T(), resp) {
		assert.Equal(b.T(), b.recorder.Code, http.StatusBadRequest)
		assert.Equal(b.T(), resp.ErrorCode, common.InvalidPath)
	}
}

func (b *BatchDownload) TestBatchDownloadEmptyFileName() {
	req := batchDownload.Request{
		Paths:    []string{"/" + b.userID + b.testDir},
		BasePath: "/" + b.userID + "/",
	}
	err := BindJsonRequest(b.ctx, req)
	if err != nil {
		panic(fmt.Sprintf("bind request failed: %v", err))
	}
	b.storage.BatchDownload(b.ctx)
	resp := &api.Response{}
	_ = ParseResp(b.recorder.Result(), resp)
	if assert.NotNil(b.T(), resp) {
		assert.Equal(b.T(), b.recorder.Code, http.StatusBadRequest)
		assert.Equal(b.T(), resp.ErrorCode, common.InvalidFileName)
	}
}

func (b *BatchDownload) TestBatchDownloadInvalidPath() {
	fileName := "test.zip"
	req := batchDownload.Request{
		FileName: fileName,
		Paths:    []string{b.userID + b.testDir},
		BasePath: "/" + b.userID + "/",
	}
	err := BindJsonRequest(b.ctx, req)
	if err != nil {
		panic(fmt.Sprintf("bind request failed: %v", err))
	}
	b.storage.BatchDownload(b.ctx)

	resp := &api.Response{}
	_ = ParseResp(b.recorder.Result(), resp)
	if assert.NotNil(b.T(), resp) {
		assert.Equal(b.T(), b.recorder.Code, http.StatusBadRequest)
		assert.Equal(b.T(), resp.ErrorCode, common.InvalidPath)
	}
}

func (b *BatchDownload) TestBatchDownloadInvalidBasePath() {
	fileName := "test.zip"
	req := batchDownload.Request{
		FileName: fileName,
		Paths:    []string{"/" + b.userID + b.testDir},
		BasePath: b.userID + "/",
	}
	err := BindJsonRequest(b.ctx, req)
	if err != nil {
		panic(fmt.Sprintf("bind request failed: %v", err))
	}
	b.storage.BatchDownload(b.ctx)

	resp := &api.Response{}
	_ = ParseResp(b.recorder.Result(), resp)
	if assert.NotNil(b.T(), resp) {
		assert.Equal(b.T(), b.recorder.Code, http.StatusBadRequest)
		assert.Equal(b.T(), resp.ErrorCode, common.InvalidPath)
	}
}

func (b *BatchDownload) TestBatchDownloadInvalidFilename() {
	fileName := "test.tar"
	req := batchDownload.Request{
		FileName: fileName,
		Paths:    []string{"/" + b.userID + b.testDir},
		BasePath: "/" + b.userID + "/",
	}
	err := BindJsonRequest(b.ctx, req)
	if err != nil {
		panic(fmt.Sprintf("bind request failed: %v", err))
	}
	b.storage.BatchDownload(b.ctx)

	resp := &api.Response{}
	_ = ParseResp(b.recorder.Result(), resp)
	if assert.NotNil(b.T(), resp) {
		assert.Equal(b.T(), b.recorder.Code, http.StatusBadRequest)
		assert.Equal(b.T(), resp.ErrorCode, common.UnsupportedCompressFileType)
	}
}

func (b *BatchDownload) TestBatchDownloadFileNotFound() {
	fileName := "test.zip"
	req := batchDownload.Request{
		FileName: fileName,
		Paths:    []string{"/" + b.userID + "/not-exist"},
		BasePath: "/" + b.userID + "/",
	}
	err := BindJsonRequest(b.ctx, req)
	if err != nil {
		panic(fmt.Sprintf("bind request failed: %v", err))
	}
	b.storage.BatchDownload(b.ctx)

	resp := &api.Response{}
	_ = ParseResp(b.recorder.Result(), resp)
	if assert.NotNil(b.T(), resp) {
		assert.Equal(b.T(), b.recorder.Code, http.StatusNotFound)
		assert.Equal(b.T(), resp.ErrorCode, common.PathNotFound)
	}
}
