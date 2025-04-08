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
	uploadFileAPI "github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/upload/file"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao/model"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/filemode"
	v20230530 "github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/local/v20230530"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/common/mock"
	"xorm.io/xorm"
)

type UploadFile struct {
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

func (uploadFile *UploadFile) SetupSuite() {
	uploadFile.T().Log("setup suite")
	logger, err := logging.NewLogger(logging.WithReleaseLevel(logging.DevelopmentLevel))
	if err != nil {
		panic(fmt.Sprintf("init logger failed: %v", err))
	}

	ctrl := gomock.NewController(uploadFile.T())
	mockStorageQuotaDao := dao.NewMockStorageQuotaDao(ctrl)
	uploadFile.mockStorageQuotaDao = mockStorageQuotaDao
	uploadFile.mockStorageQuotaDao.EXPECT().Get(gomock.Any(), gomock.Any()).Return(true, &model.StorageQuota{StorageLimit: 1000}, nil).AnyTimes()
	mockStorageOperationLogDao := dao.NewMockStorageOperationLogDao(ctrl)
	uploadFile.mockStorageOperationLogDao = mockStorageOperationLogDao
	uploadFile.mockStorageOperationLogDao.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mockEngine, _ := mock.Engine()
	mockEngine.SetLogger(middleware.NewXormLogger(uploadFile.logger, true))
	uploadFile.engine = mockEngine

	uploadFile.logger = logger
	currentDir, storage := CreateLocalStorage(uploadFile.userID, v20230530.SoftLink, uploadFile.mockStorageQuotaDao, uploadFile.mockStorageOperationLogDao, uploadFile.engine)
	uploadFile.storage = storage

	uploadFile.userID = "4TiSBX39DtN"
	uploadFile.testDir = "/test-upload-file"

	tmpDir := currentDir + "/" + uploadFile.userID + uploadFile.testDir
	err = os.MkdirAll(tmpDir, filemode.Directory)
	if err != nil {
		panic(fmt.Sprintf("create temp dir failed: %v", err))
	}
	uploadFile.tmpDir = tmpDir
	filePath := filepath.Join(uploadFile.tmpDir, "file-1")
	err = os.WriteFile(filePath, nil, 0644)
	if err != nil {
		panic(fmt.Sprintf("Failed to create file: %v", err))
	}
	uploadFile.T().Log("File created:", filePath)

}

func (uploadFile *UploadFile) TearDownSuite() {
	uploadFile.T().Log("teardown suite")
	_ = os.RemoveAll(filepath.Dir(uploadFile.tmpDir))
}

func (uploadFile *UploadFile) SetupTest() {
	uploadFile.T().Log("setup test")
	recorder := httptest.NewRecorder()
	ctx := GetTestGinContext(recorder)
	ctx.Set(logging.LoggerName, uploadFile.logger)
	uploadFile.ctx = ctx
	uploadFile.recorder = recorder
}

func (uploadFile *UploadFile) TearDownTest() {
	uploadFile.T().Log("teardown test")
}

func TestUploadFile(t *testing.T) {
	suite.Run(t, new(UploadFile))
}

func (uploadFile *UploadFile) TestUploadFileSuccess() {

	req := uploadFileAPI.Request{
		Path: "/" + uploadFile.userID + uploadFile.testDir + "/file-2",
	}
	err, reader := GenerateRandomData(5)
	if err != nil {
		panic(err)
	}
	BindRequest(uploadFile.ctx, req)
	uploadFile.ctx.Request.Body = io.NopCloser(reader)
	uploadFile.storage.UploadFile(uploadFile.ctx)

	resp := &uploadFileAPI.Response{}
	err = ParseResp(uploadFile.recorder.Result(), resp)
	if !assert.Nil(uploadFile.T(), err) {
		uploadFile.T().Logf("%v: %v", time.Now().Format("2006-01-02 15:04:05"), err)
	} else {
		spew.Dump(resp)
	}
}

func (uploadFile *UploadFile) TestUploadFileSuccessOverwrite() {

	req := uploadFileAPI.Request{
		Path:      "/" + uploadFile.userID + uploadFile.testDir + "/file-1",
		Overwrite: true,
	}
	err, reader := GenerateRandomData(5)
	if err != nil {
		panic(err)
	}
	BindRequest(uploadFile.ctx, req)
	uploadFile.ctx.Request.Body = io.NopCloser(reader)
	uploadFile.storage.UploadFile(uploadFile.ctx)

	resp := &uploadFileAPI.Response{}
	err = ParseResp(uploadFile.recorder.Result(), resp)
	if !assert.Nil(uploadFile.T(), err) {
		uploadFile.T().Logf("%v: %v", time.Now().Format("2006-01-02 15:04:05"), err)
	} else {
		spew.Dump(resp)
	}
}

// ------------------------ error case ------------------------

func (uploadFile *UploadFile) TestUploadFileInvalidPath() {

	req := uploadFileAPI.Request{
		Path:      uploadFile.userID + uploadFile.testDir + "/file-1",
		Overwrite: true,
	}
	err, reader := GenerateRandomData(5)
	if err != nil {
		panic(err)
	}
	BindRequest(uploadFile.ctx, req)
	uploadFile.ctx.Request.Body = io.NopCloser(reader)
	uploadFile.storage.UploadFile(uploadFile.ctx)

	resp := &uploadFileAPI.Response{}
	_ = ParseResp(uploadFile.recorder.Result(), resp)
	if assert.NotNil(uploadFile.T(), resp) {
		assert.Equal(uploadFile.T(), uploadFile.recorder.Code, http.StatusBadRequest)
		assert.Equal(uploadFile.T(), resp.ErrorCode, common.InvalidPath)
	}
}

func (uploadFile *UploadFile) TestUploadFilePathExists() {

	req := uploadFileAPI.Request{
		Path: "/" + uploadFile.userID + uploadFile.testDir + "/file-1",
	}
	err, reader := GenerateRandomData(5)
	if err != nil {
		panic(err)
	}
	BindRequest(uploadFile.ctx, req)
	uploadFile.ctx.Request.Body = io.NopCloser(reader)
	uploadFile.storage.UploadFile(uploadFile.ctx)

	resp := &uploadFileAPI.Response{}
	_ = ParseResp(uploadFile.recorder.Result(), resp)
	if assert.NotNil(uploadFile.T(), resp) {
		assert.Equal(uploadFile.T(), uploadFile.recorder.Code, http.StatusBadRequest)
		assert.Equal(uploadFile.T(), resp.ErrorCode, common.PathExists)
	}
}

func (uploadFile *UploadFile) TestUploadFileSizeTooLarge() {
	filesize := 101 * 1024 //101m
	req := uploadFileAPI.Request{
		Path: "/" + uploadFile.userID + uploadFile.testDir + "/file-3",
	}
	err, reader := GenerateRandomData(int64(filesize))
	if err != nil {
		panic(err)
	}
	BindRequest(uploadFile.ctx, req)
	uploadFile.ctx.Request.Body = io.NopCloser(reader)
	uploadFile.storage.UploadFile(uploadFile.ctx)

	resp := &uploadFileAPI.Response{}
	_ = ParseResp(uploadFile.recorder.Result(), resp)
	if assert.NotNil(uploadFile.T(), resp) {
		assert.Equal(uploadFile.T(), uploadFile.recorder.Code, http.StatusBadRequest)
		assert.Equal(uploadFile.T(), resp.ErrorCode, common.SizeTooLarge)
	}
}
