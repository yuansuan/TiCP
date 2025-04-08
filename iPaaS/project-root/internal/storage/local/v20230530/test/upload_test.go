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
	uploadComplete "github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/upload/complete"
	uploadInit "github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/upload/init"
	uploadSlice "github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/upload/slice"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao/model"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/filemode"
	v20230530 "github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/local/v20230530"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/common/mock"
	"xorm.io/xorm"
)

type Upload struct {
	suite.Suite
	logger                     *logging.Logger
	storage                    *v20230530.Storage
	tmpDir                     string
	userID                     string
	testDir                    string
	curDir                     string
	ctx                        *gin.Context
	recorder                   *httptest.ResponseRecorder
	engine                     *xorm.Engine
	mockUploadInfoDao          *dao.MockUploadInfoDao
	mockStorageQuotaDao        *dao.MockStorageQuotaDao
	mockStorageOperationLogDao *dao.MockStorageOperationLogDao
}

func (upload *Upload) SetupSuite() {
	upload.T().Log("setup suite")

	logger, err := logging.NewLogger(logging.WithReleaseLevel(logging.DevelopmentLevel))
	if err != nil {
		panic(fmt.Sprintf("init logger failed: %v", err))
	}

	ctrl := gomock.NewController(upload.T())
	mockUploadInfoDao := dao.NewMockUploadInfoDao(ctrl)
	upload.mockUploadInfoDao = mockUploadInfoDao
	upload.mockUploadInfoDao.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	mockStorageQuotaDao := dao.NewMockStorageQuotaDao(ctrl)
	upload.mockStorageQuotaDao = mockStorageQuotaDao
	upload.mockStorageQuotaDao.EXPECT().Get(gomock.Any(), gomock.Any()).Return(true, &model.StorageQuota{StorageLimit: 1000}, nil).AnyTimes()
	mockStorageOperationLogDao := dao.NewMockStorageOperationLogDao(ctrl)
	upload.mockStorageOperationLogDao = mockStorageOperationLogDao
	upload.mockStorageOperationLogDao.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mockEngine, _ := mock.Engine()
	mockEngine.SetLogger(middleware.NewXormLogger(upload.logger, true))
	upload.engine = mockEngine

	upload.logger = logger
	currentDir, storage := CreateLocalStorageWithUploadDao(upload.userID, v20230530.SoftLink, upload.mockUploadInfoDao, upload.mockStorageQuotaDao, upload.mockStorageOperationLogDao, upload.engine)
	upload.storage = storage

	upload.userID = "4TiSBX39DtN"
	upload.testDir = "/test-upload"
	upload.curDir = currentDir

	tmpDir := currentDir + "/" + upload.userID + upload.testDir
	err = os.MkdirAll(tmpDir, filemode.Directory)
	if err != nil {
		panic(fmt.Sprintf("create temp dir failed: %v", err))
	}
	upload.tmpDir = tmpDir

	filePath := filepath.Join(upload.tmpDir, "file-1")
	err = os.WriteFile(filePath, nil, 0644)
	if err != nil {
		panic(fmt.Sprintf("Failed to create file: %v", err))
	}
	upload.T().Log("File created:", filePath)

}

func (upload *Upload) SetupTest() {
	upload.T().Log("setup test")

	recorder := httptest.NewRecorder()
	ctx := GetTestGinContext(recorder)
	ctx.Set(logging.LoggerName, upload.logger)
	upload.ctx = ctx
	upload.recorder = recorder
}

func (upload *Upload) TearDownTest() {
	upload.T().Log("teardown test")
}

func (upload *Upload) TearDownSuite() {
	upload.T().Log("teardown suite")
	_ = os.RemoveAll(filepath.Dir(upload.tmpDir))
	_ = os.RemoveAll(filepath.Join(upload.curDir, ".uploading"))
}

func TestUpload(t *testing.T) {
	suite.Run(t, new(Upload))
}

func (upload *Upload) TestUploadSuccess() {
	path := "/" + upload.userID + upload.testDir + "/a.txt"
	initReq := uploadInit.Request{
		Path: path,
		Size: 1024 * 5,
	}
	err := BindJsonRequest(upload.ctx, initReq)
	if err != nil {
		panic(fmt.Sprintf("bind request failed: %v", err))
	}

	upload.storage.UploadInit(upload.ctx)

	resp := &uploadInit.Response{}
	err = ParseResp(upload.recorder.Result(), resp)
	if !assert.Nil(upload.T(), err) {
		upload.T().Logf("%v: %v", time.Now().Format("2006-01-02 15:04:05"), err)
	} else {
		spew.Dump(resp.Data)
		UploadID := resp.Data.UploadID

		sliceReq := uploadSlice.Request{
			UploadID: UploadID,
			Length:   1024 * 5,
		}

		recorder := httptest.NewRecorder()
		ctx := GetTestGinContext(recorder)
		ctx.Set(logging.LoggerName, upload.logger)
		BindRequest(ctx, sliceReq)
		err, reader := GenerateRandomData(5)
		if err != nil {
			panic(err)
		}
		ctx.Request.Body = io.NopCloser(reader)

		upload.storage.UploadSlice(ctx)
		resp := &uploadSlice.Response{}
		err = ParseResp(recorder.Result(), resp)
		if !assert.Nil(upload.T(), err) {
			upload.T().Logf("%v: %v", time.Now().Format("2006-01-02 15:04:05"), err)
		} else {
			spew.Dump(resp)

			completeReq := uploadComplete.Request{
				UploadID: UploadID,
				Path:     path,
			}

			recorder := httptest.NewRecorder()
			ctx := GetTestGinContext(recorder)
			ctx.Set(logging.LoggerName, upload.logger)

			err := BindJsonRequest(ctx, completeReq)
			if err != nil {
				panic(fmt.Sprintf("bind request failed: %v", err))
			}

			upload.storage.UploadComplete(ctx)
			resp := &uploadComplete.Response{}
			err = ParseResp(recorder.Result(), resp)
			if !assert.Nil(upload.T(), err) {
				upload.T().Logf("%v: %v", time.Now().Format("2006-01-02 15:04:05"), err)
			} else {
				spew.Dump(resp)
			}
		}
	}
}

// ------------------------ error case ------------------------
func (upload *Upload) TestUploadInitSizeInvalid() {

	path := "/" + upload.userID + upload.testDir + "/a.txt"
	initReq := uploadInit.Request{
		Path: path,
		Size: -1,
	}
	err := BindJsonRequest(upload.ctx, initReq)
	if err != nil {
		panic(fmt.Sprintf("bind request failed: %v", err))
	}

	upload.storage.UploadInit(upload.ctx)

	resp := &uploadInit.Response{}
	err = ParseResp(upload.recorder.Result(), resp)
	_ = ParseResp(upload.recorder.Result(), resp)
	if assert.NotNil(upload.T(), resp) {
		assert.Equal(upload.T(), upload.recorder.Code, http.StatusBadRequest)
		assert.Equal(upload.T(), resp.ErrorCode, common.InvalidSize)
	}
}

func (upload *Upload) TestUploadInitSizeTooLarge() {

	path := "/" + upload.userID + upload.testDir + "/a.txt"
	initReq := uploadInit.Request{
		Path: path,
		Size: 1024 * 1024 * 1024 * 1024 * 2,
	}
	err := BindJsonRequest(upload.ctx, initReq)
	if err != nil {
		panic(fmt.Sprintf("bind request failed: %v", err))
	}

	upload.storage.UploadInit(upload.ctx)

	resp := &uploadInit.Response{}
	_ = ParseResp(upload.recorder.Result(), resp)
	if assert.NotNil(upload.T(), resp) {
		assert.Equal(upload.T(), upload.recorder.Code, http.StatusBadRequest)
		assert.Equal(upload.T(), resp.ErrorCode, common.SizeTooLarge)
	}
}

func (upload *Upload) TestUploadInitInvalidPath() {

	path := upload.userID + upload.testDir + "/a.txt"
	initReq := uploadInit.Request{
		Path: path,
		Size: 100 * 1024 * 1024,
	}
	err := BindJsonRequest(upload.ctx, initReq)
	if err != nil {
		panic(fmt.Sprintf("bind request failed: %v", err))
	}

	upload.storage.UploadInit(upload.ctx)

	resp := &uploadInit.Response{}
	_ = ParseResp(upload.recorder.Result(), resp)
	if assert.NotNil(upload.T(), resp) {
		assert.Equal(upload.T(), upload.recorder.Code, http.StatusBadRequest)
		assert.Equal(upload.T(), resp.ErrorCode, common.InvalidPath)
	}
}

func (upload *Upload) TestUploadInitPathExist() {

	path := "/" + upload.userID + upload.testDir + "/file-1"
	initReq := uploadInit.Request{
		Path: path,
		Size: 100 * 1024 * 1024,
	}
	err := BindJsonRequest(upload.ctx, initReq)
	if err != nil {
		panic(fmt.Sprintf("bind request failed: %v", err))
	}

	upload.storage.UploadInit(upload.ctx)

	resp := &uploadInit.Response{}
	_ = ParseResp(upload.recorder.Result(), resp)
	if assert.NotNil(upload.T(), resp) {
		assert.Equal(upload.T(), upload.recorder.Code, http.StatusBadRequest)
		assert.Equal(upload.T(), resp.ErrorCode, common.PathExists)
	}
}

func (upload *Upload) TestUploadSliceUploadIDNotFound() {

	path := "/" + upload.userID + upload.testDir + "/a.txt"
	initReq := uploadInit.Request{
		Path: path,
		Size: 1024 * 5,
	}
	err := BindJsonRequest(upload.ctx, initReq)
	if err != nil {
		panic(fmt.Sprintf("bind request failed: %v", err))
	}

	upload.storage.UploadInit(upload.ctx)

	resp := &uploadInit.Response{}
	err = ParseResp(upload.recorder.Result(), resp)
	if !assert.Nil(upload.T(), err) {
		upload.T().Logf("%v: %v", time.Now().Format("2006-01-02 15:04:05"), err)
	} else {
		spew.Dump(resp.Data)

		sliceReq := uploadSlice.Request{
			Length: 1024 * 5,
		}

		recorder := httptest.NewRecorder()
		ctx := GetTestGinContext(recorder)
		ctx.Set(logging.LoggerName, upload.logger)
		BindRequest(ctx, sliceReq)
		err, reader := GenerateRandomData(5)
		if err != nil {
			panic(err)
		}
		ctx.Request.Body = io.NopCloser(reader)

		upload.storage.UploadSlice(ctx)
		resp := &uploadSlice.Response{}
		_ = ParseResp(recorder.Result(), resp)
		if assert.NotNil(upload.T(), resp) {
			assert.Equal(upload.T(), recorder.Code, http.StatusNotFound)
			assert.Equal(upload.T(), resp.ErrorCode, common.UploadIDNotFound)
		}
	}
}

func (upload *Upload) TestUploadSliceLengthTooLarge() {

	path := "/" + upload.userID + upload.testDir + "/a.txt"
	initReq := uploadInit.Request{
		Path: path,
		Size: 1024 * 5,
	}
	err := BindJsonRequest(upload.ctx, initReq)
	if err != nil {
		panic(fmt.Sprintf("bind request failed: %v", err))
	}

	upload.storage.UploadInit(upload.ctx)

	resp := &uploadInit.Response{}
	err = ParseResp(upload.recorder.Result(), resp)
	if !assert.Nil(upload.T(), err) {
		upload.T().Logf("%v: %v", time.Now().Format("2006-01-02 15:04:05"), err)
	} else {
		spew.Dump(resp.Data)
		uploadID := resp.Data.UploadID

		sliceReq := uploadSlice.Request{
			UploadID: uploadID,
			Length:   1024 * 1024 * 100,
		}

		recorder := httptest.NewRecorder()
		ctx := GetTestGinContext(recorder)
		ctx.Set(logging.LoggerName, upload.logger)
		BindRequest(ctx, sliceReq)
		err, reader := GenerateRandomData(5)
		if err != nil {
			panic(err)
		}
		ctx.Request.Body = io.NopCloser(reader)

		upload.storage.UploadSlice(ctx)
		resp := &uploadSlice.Response{}
		_ = ParseResp(recorder.Result(), resp)
		if assert.NotNil(upload.T(), resp) {
			assert.Equal(upload.T(), recorder.Code, http.StatusBadRequest)
			assert.Equal(upload.T(), resp.ErrorCode, common.LengthTooLarge)
		}
	}
}

func (upload *Upload) TestUploadSliceInvalidLength() {

	path := "/" + upload.userID + upload.testDir + "/a.txt"
	initReq := uploadInit.Request{
		Path: path,
		Size: 1024 * 5,
	}
	err := BindJsonRequest(upload.ctx, initReq)
	if err != nil {
		panic(fmt.Sprintf("bind request failed: %v", err))
	}

	upload.storage.UploadInit(upload.ctx)

	resp := &uploadInit.Response{}
	err = ParseResp(upload.recorder.Result(), resp)
	if !assert.Nil(upload.T(), err) {
		upload.T().Logf("%v: %v", time.Now().Format("2006-01-02 15:04:05"), err)
	} else {
		spew.Dump(resp.Data)
		uploadID := resp.Data.UploadID

		sliceReq := uploadSlice.Request{
			UploadID: uploadID,
			Length:   -1,
		}

		recorder := httptest.NewRecorder()
		ctx := GetTestGinContext(recorder)
		ctx.Set(logging.LoggerName, upload.logger)
		BindRequest(ctx, sliceReq)
		err, reader := GenerateRandomData(5)
		if err != nil {
			panic(err)
		}
		ctx.Request.Body = io.NopCloser(reader)

		upload.storage.UploadSlice(ctx)
		resp := &uploadSlice.Response{}
		_ = ParseResp(recorder.Result(), resp)
		if assert.NotNil(upload.T(), resp) {
			assert.Equal(upload.T(), recorder.Code, http.StatusBadRequest)
			assert.Equal(upload.T(), resp.ErrorCode, common.InvalidLength)
		}
	}
}

func (upload *Upload) TestUploadSliceInvalidOffset() {

	path := "/" + upload.userID + upload.testDir + "/a.txt"
	initReq := uploadInit.Request{
		Path: path,
		Size: 1024 * 5,
	}
	err := BindJsonRequest(upload.ctx, initReq)
	if err != nil {
		panic(fmt.Sprintf("bind request failed: %v", err))
	}

	upload.storage.UploadInit(upload.ctx)

	resp := &uploadInit.Response{}
	err = ParseResp(upload.recorder.Result(), resp)
	if !assert.Nil(upload.T(), err) {
		upload.T().Logf("%v: %v", time.Now().Format("2006-01-02 15:04:05"), err)
	} else {
		spew.Dump(resp.Data)
		uploadID := resp.Data.UploadID

		sliceReq := uploadSlice.Request{
			UploadID: uploadID,
			Offset:   -1,
			Length:   100,
		}

		recorder := httptest.NewRecorder()
		ctx := GetTestGinContext(recorder)
		ctx.Set(logging.LoggerName, upload.logger)
		BindRequest(ctx, sliceReq)
		err, reader := GenerateRandomData(5)
		if err != nil {
			panic(err)
		}
		ctx.Request.Body = io.NopCloser(reader)

		upload.storage.UploadSlice(ctx)
		resp := &uploadSlice.Response{}
		_ = ParseResp(recorder.Result(), resp)
		if assert.NotNil(upload.T(), resp) {
			assert.Equal(upload.T(), recorder.Code, http.StatusBadRequest)
			assert.Equal(upload.T(), resp.ErrorCode, common.InvalidOffset)
		}
	}
}

func (upload *Upload) TestUploadSliceOffsetTooLarge() {

	path := "/" + upload.userID + upload.testDir + "/a.txt"
	initReq := uploadInit.Request{
		Path: path,
		Size: 1024 * 5,
	}
	err := BindJsonRequest(upload.ctx, initReq)
	if err != nil {
		panic(fmt.Sprintf("bind request failed: %v", err))
	}

	upload.storage.UploadInit(upload.ctx)

	resp := &uploadInit.Response{}
	err = ParseResp(upload.recorder.Result(), resp)
	if !assert.Nil(upload.T(), err) {
		upload.T().Logf("%v: %v", time.Now().Format("2006-01-02 15:04:05"), err)
	} else {
		spew.Dump(resp.Data)
		uploadID := resp.Data.UploadID

		sliceReq := uploadSlice.Request{
			UploadID: uploadID,
			Offset:   1024 * 10,
			Length:   100,
		}

		recorder := httptest.NewRecorder()
		ctx := GetTestGinContext(recorder)
		ctx.Set(logging.LoggerName, upload.logger)
		BindRequest(ctx, sliceReq)
		err, reader := GenerateRandomData(5)
		if err != nil {
			panic(err)
		}
		ctx.Request.Body = io.NopCloser(reader)

		upload.storage.UploadSlice(ctx)
		resp := &uploadSlice.Response{}
		_ = ParseResp(recorder.Result(), resp)
		if assert.NotNil(upload.T(), resp) {
			assert.Equal(upload.T(), recorder.Code, http.StatusBadRequest)
			assert.Equal(upload.T(), resp.ErrorCode, common.InvalidOffset)
		}
	}
}

func (upload *Upload) TestUploadSliceOffsetPlusLengthTooLarge() {

	path := "/" + upload.userID + upload.testDir + "/a.txt"
	initReq := uploadInit.Request{
		Path: path,
		Size: 1024 * 5,
	}
	err := BindJsonRequest(upload.ctx, initReq)
	if err != nil {
		panic(fmt.Sprintf("bind request failed: %v", err))
	}

	upload.storage.UploadInit(upload.ctx)

	resp := &uploadInit.Response{}
	err = ParseResp(upload.recorder.Result(), resp)
	if !assert.Nil(upload.T(), err) {
		upload.T().Logf("%v: %v", time.Now().Format("2006-01-02 15:04:05"), err)
	} else {
		spew.Dump(resp.Data)
		uploadID := resp.Data.UploadID

		sliceReq := uploadSlice.Request{
			UploadID: uploadID,
			Offset:   1024 * 4,
			Length:   1024 * 3,
		}

		recorder := httptest.NewRecorder()
		ctx := GetTestGinContext(recorder)
		ctx.Set(logging.LoggerName, upload.logger)
		BindRequest(ctx, sliceReq)
		err, reader := GenerateRandomData(5)
		if err != nil {
			panic(err)
		}
		ctx.Request.Body = io.NopCloser(reader)

		upload.storage.UploadSlice(ctx)
		resp := &uploadSlice.Response{}
		_ = ParseResp(recorder.Result(), resp)
		if assert.NotNil(upload.T(), resp) {
			assert.Equal(upload.T(), recorder.Code, http.StatusBadRequest)
			assert.Equal(upload.T(), resp.ErrorCode, common.InvalidLength)
		}
	}
}

func (upload *Upload) TestUploadSliceInvalidSlice() {

	path := "/" + upload.userID + upload.testDir + "/a.txt"
	initReq := uploadInit.Request{
		Path: path,
		Size: 1024 * 5,
	}
	err := BindJsonRequest(upload.ctx, initReq)
	if err != nil {
		panic(fmt.Sprintf("bind request failed: %v", err))
	}

	upload.storage.UploadInit(upload.ctx)

	resp := &uploadInit.Response{}
	err = ParseResp(upload.recorder.Result(), resp)
	if !assert.Nil(upload.T(), err) {
		upload.T().Logf("%v: %v", time.Now().Format("2006-01-02 15:04:05"), err)
	} else {
		spew.Dump(resp.Data)
		uploadID := resp.Data.UploadID

		sliceReq := uploadSlice.Request{
			UploadID: uploadID,
			Offset:   100,
			Length:   100,
		}

		recorder := httptest.NewRecorder()
		ctx := GetTestGinContext(recorder)
		ctx.Set(logging.LoggerName, upload.logger)
		BindRequest(ctx, sliceReq)

		upload.storage.UploadSlice(ctx)
		resp := &uploadSlice.Response{}
		_ = ParseResp(recorder.Result(), resp)
		if assert.NotNil(upload.T(), resp) {
			assert.Equal(upload.T(), recorder.Code, http.StatusBadRequest)
			assert.Equal(upload.T(), resp.ErrorCode, common.InvalidArgumentErrorCode)
		}
	}
}

func (upload *Upload) TestUploadCompleteUploadIDNotFound() {

	path := "/" + upload.userID + upload.testDir + "/a.txt"
	initReq := uploadInit.Request{
		Path: path,
		Size: 1024 * 5,
	}
	err := BindJsonRequest(upload.ctx, initReq)
	if err != nil {
		panic(fmt.Sprintf("bind request failed: %v", err))
	}

	upload.storage.UploadInit(upload.ctx)

	resp := &uploadInit.Response{}
	err = ParseResp(upload.recorder.Result(), resp)
	if !assert.Nil(upload.T(), err) {
		upload.T().Logf("%v: %v", time.Now().Format("2006-01-02 15:04:05"), err)
	} else {
		spew.Dump(resp.Data)
		uploadID := resp.Data.UploadID

		sliceReq := uploadSlice.Request{
			UploadID: uploadID,
			Length:   1024 * 5,
		}

		recorder := httptest.NewRecorder()
		ctx := GetTestGinContext(recorder)
		ctx.Set(logging.LoggerName, upload.logger)
		BindRequest(ctx, sliceReq)
		err, reader := GenerateRandomData(5)
		if err != nil {
			panic(err)
		}
		ctx.Request.Body = io.NopCloser(reader)

		upload.storage.UploadSlice(ctx)
		resp := &uploadSlice.Response{}
		err = ParseResp(recorder.Result(), resp)
		if !assert.Nil(upload.T(), err) {
			upload.T().Logf("%v: %v", time.Now().Format("2006-01-02 15:04:05"), err)
		} else {
			spew.Dump(resp)

			completeReq := uploadComplete.Request{
				Path: path,
			}

			recorder := httptest.NewRecorder()
			ctx := GetTestGinContext(recorder)
			ctx.Set(logging.LoggerName, upload.logger)

			err := BindJsonRequest(ctx, completeReq)
			if err != nil {
				panic(fmt.Sprintf("bind request failed: %v", err))
			}

			upload.storage.UploadComplete(ctx)
			resp := &uploadComplete.Response{}
			_ = ParseResp(recorder.Result(), resp)
			if assert.NotNil(upload.T(), resp) {
				assert.Equal(upload.T(), recorder.Code, http.StatusNotFound)
				assert.Equal(upload.T(), resp.ErrorCode, common.UploadIDNotFound)
			}
		}
	}
}

func (upload *Upload) TestUploadCompleteInvalidPath() {

	path := "/" + upload.userID + upload.testDir + "/a.txt"
	initReq := uploadInit.Request{
		Path: path,
		Size: 1024 * 5,
	}
	err := BindJsonRequest(upload.ctx, initReq)
	if err != nil {
		panic(fmt.Sprintf("bind request failed: %v", err))
	}

	upload.storage.UploadInit(upload.ctx)

	resp := &uploadInit.Response{}
	err = ParseResp(upload.recorder.Result(), resp)
	if !assert.Nil(upload.T(), err) {
		upload.T().Logf("%v: %v", time.Now().Format("2006-01-02 15:04:05"), err)
	} else {
		spew.Dump(resp.Data)
		uploadID := resp.Data.UploadID

		sliceReq := uploadSlice.Request{
			UploadID: uploadID,
			Length:   1024 * 5,
		}

		recorder := httptest.NewRecorder()
		ctx := GetTestGinContext(recorder)
		ctx.Set(logging.LoggerName, upload.logger)
		BindRequest(ctx, sliceReq)
		err, reader := GenerateRandomData(5)
		if err != nil {
			panic(err)
		}
		ctx.Request.Body = io.NopCloser(reader)

		upload.storage.UploadSlice(ctx)
		resp := &uploadSlice.Response{}
		err = ParseResp(recorder.Result(), resp)
		if !assert.Nil(upload.T(), err) {
			upload.T().Logf("%v: %v", time.Now().Format("2006-01-02 15:04:05"), err)
		} else {
			spew.Dump(resp)

			completeReq := uploadComplete.Request{
				UploadID: uploadID,
				Path:     upload.userID + upload.testDir + "/a.txt",
			}

			recorder := httptest.NewRecorder()
			ctx := GetTestGinContext(recorder)
			ctx.Set(logging.LoggerName, upload.logger)

			err := BindJsonRequest(ctx, completeReq)
			if err != nil {
				panic(fmt.Sprintf("bind request failed: %v", err))
			}

			upload.storage.UploadComplete(ctx)
			resp := &uploadComplete.Response{}
			_ = ParseResp(recorder.Result(), resp)
			if assert.NotNil(upload.T(), resp) {
				assert.Equal(upload.T(), recorder.Code, http.StatusBadRequest)
				assert.Equal(upload.T(), resp.ErrorCode, common.InvalidPath)
			}
		}
	}
}

func (upload *Upload) TestUploadCompletePathNotMatch() {

	path := "/" + upload.userID + upload.testDir + "/a.txt"
	initReq := uploadInit.Request{
		Path: path,
		Size: 1024 * 5,
	}
	err := BindJsonRequest(upload.ctx, initReq)
	if err != nil {
		panic(fmt.Sprintf("bind request failed: %v", err))
	}

	upload.storage.UploadInit(upload.ctx)

	resp := &uploadInit.Response{}
	err = ParseResp(upload.recorder.Result(), resp)
	if !assert.Nil(upload.T(), err) {
		upload.T().Logf("%v: %v", time.Now().Format("2006-01-02 15:04:05"), err)
	} else {
		spew.Dump(resp.Data)
		uploadID := resp.Data.UploadID

		sliceReq := uploadSlice.Request{
			UploadID: uploadID,
			Length:   1024 * 5,
		}

		recorder := httptest.NewRecorder()
		ctx := GetTestGinContext(recorder)
		ctx.Set(logging.LoggerName, upload.logger)
		BindRequest(ctx, sliceReq)
		err, reader := GenerateRandomData(5)
		if err != nil {
			panic(err)
		}
		ctx.Request.Body = io.NopCloser(reader)

		upload.storage.UploadSlice(ctx)
		resp := &uploadSlice.Response{}
		err = ParseResp(recorder.Result(), resp)
		if !assert.Nil(upload.T(), err) {
			upload.T().Logf("%v: %v", time.Now().Format("2006-01-02 15:04:05"), err)
		} else {
			spew.Dump(resp)

			completeReq := uploadComplete.Request{
				UploadID: uploadID,
				Path:     "/" + upload.userID + upload.testDir + "/b.txt",
			}

			recorder := httptest.NewRecorder()
			ctx := GetTestGinContext(recorder)
			ctx.Set(logging.LoggerName, upload.logger)

			err := BindJsonRequest(ctx, completeReq)
			if err != nil {
				panic(fmt.Sprintf("bind request failed: %v", err))
			}

			upload.storage.UploadComplete(ctx)
			resp := &uploadComplete.Response{}
			_ = ParseResp(recorder.Result(), resp)
			if assert.NotNil(upload.T(), resp) {
				assert.Equal(upload.T(), recorder.Code, http.StatusBadRequest)
				assert.Equal(upload.T(), resp.ErrorCode, common.PathNotMatchUploadInit)
			}
		}
	}
}

func (upload *Upload) TestUploadCompleteFileExists() {

	path := "/" + upload.userID + upload.testDir + "/a.txt"
	initReq := uploadInit.Request{
		Path: path,
		Size: 1024 * 5,
	}
	err := BindJsonRequest(upload.ctx, initReq)
	if err != nil {
		panic(fmt.Sprintf("bind request failed: %v", err))
	}

	upload.storage.UploadInit(upload.ctx)

	resp := &uploadInit.Response{}
	err = ParseResp(upload.recorder.Result(), resp)
	if !assert.Nil(upload.T(), err) {
		upload.T().Logf("%v: %v", time.Now().Format("2006-01-02 15:04:05"), err)
	} else {
		spew.Dump(resp.Data)
		uploadID := resp.Data.UploadID

		sliceReq := uploadSlice.Request{
			UploadID: uploadID,
			Length:   1024 * 5,
		}

		recorder := httptest.NewRecorder()
		ctx := GetTestGinContext(recorder)
		ctx.Set(logging.LoggerName, upload.logger)
		BindRequest(ctx, sliceReq)
		err, reader := GenerateRandomData(5)
		if err != nil {
			panic(err)
		}
		ctx.Request.Body = io.NopCloser(reader)

		upload.storage.UploadSlice(ctx)
		resp := &uploadSlice.Response{}
		err = ParseResp(recorder.Result(), resp)
		if !assert.Nil(upload.T(), err) {
			upload.T().Logf("%v: %v", time.Now().Format("2006-01-02 15:04:05"), err)
		} else {
			spew.Dump(resp)
			filePath := filepath.Join(upload.tmpDir, "a.txt")
			err = os.WriteFile(filePath, nil, 0644)
			if err != nil {
				panic(fmt.Sprintf("Failed to create file: %v", err))
			}
			completeReq := uploadComplete.Request{
				UploadID: uploadID,
				Path:     path,
			}

			recorder := httptest.NewRecorder()
			ctx := GetTestGinContext(recorder)
			ctx.Set(logging.LoggerName, upload.logger)

			err := BindJsonRequest(ctx, completeReq)
			if err != nil {
				panic(fmt.Sprintf("bind request failed: %v", err))
			}

			upload.storage.UploadComplete(ctx)
			resp := &uploadComplete.Response{}
			_ = ParseResp(recorder.Result(), resp)
			if assert.NotNil(upload.T(), resp) {
				assert.Equal(upload.T(), recorder.Code, http.StatusBadRequest)
				assert.Equal(upload.T(), resp.ErrorCode, common.PathExists)
			}
			err = os.RemoveAll(filePath)
			if err != nil {
				panic(fmt.Sprintf("Failed to remove file: %v", err))
			}
		}
	}
}
