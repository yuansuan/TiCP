package test

import (
	"fmt"
	"math/rand"
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
	"github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/copyRange"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao/model"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/filemode"
	v20230530 "github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/local/v20230530"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/common/mock"
	"xorm.io/xorm"
)

type CopyRange struct {
	suite.Suite
	logger                     *logging.Logger
	storage                    *v20230530.Storage
	tmpDir                     string
	userID                     string
	testDir                    string
	srcDir                     string
	destDir                    string
	ctx                        *gin.Context
	recorder                   *httptest.ResponseRecorder
	engine                     *xorm.Engine
	mockStorageQuotaDao        *dao.MockStorageQuotaDao
	mockStorageOperationLogDao *dao.MockStorageOperationLogDao
}

func (cr *CopyRange) SetupSuite() {
	cr.T().Log("setup suite")
	logger, err := logging.NewLogger(logging.WithReleaseLevel(logging.DevelopmentLevel))
	if err != nil {
		panic(fmt.Sprintf("init logger failed: %v", err))
	}

	ctrl := gomock.NewController(cr.T())
	mockStorageQuotaDao := dao.NewMockStorageQuotaDao(ctrl)
	cr.mockStorageQuotaDao = mockStorageQuotaDao
	cr.mockStorageQuotaDao.EXPECT().Get(gomock.Any(), gomock.Any()).Return(true, &model.StorageQuota{StorageLimit: 1000}, nil).AnyTimes()
	mockStorageOperationLogDao := dao.NewMockStorageOperationLogDao(ctrl)
	cr.mockStorageOperationLogDao = mockStorageOperationLogDao
	cr.mockStorageOperationLogDao.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mockEngine, _ := mock.Engine()
	mockEngine.SetLogger(middleware.NewXormLogger(cr.logger, true))
	cr.engine = mockEngine

	cr.logger = logger
	currentDir, storage := CreateLocalStorage(cr.userID, v20230530.SoftLink, cr.mockStorageQuotaDao, cr.mockStorageOperationLogDao, cr.engine)
	cr.storage = storage

	cr.userID = "4TiSBX39DtN"
	cr.testDir = "/test-copy-range"

	tmpDir := currentDir + "/" + cr.userID + cr.testDir
	cr.tmpDir = tmpDir
	cr.srcDir = "/" + cr.userID + cr.testDir + "/src"
	cr.destDir = "/" + cr.userID + cr.testDir + "/dest"

	err = os.MkdirAll(filepath.Join(cr.tmpDir, "src"), filemode.Directory)
	if err != nil {
		panic(fmt.Sprintf("create src dir failed: %v", err))
	}

	filePath := filepath.Join(cr.tmpDir, "src", "file-1")
	data := make([]byte, 5*1024*1024)
	rand.Read(data)
	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		panic(fmt.Sprintf("Failed to create file: %v", err))
	}
	cr.T().Log("File created:", filePath)

	err = os.MkdirAll(filepath.Join(cr.tmpDir, "dest"), filemode.Directory)
	if err != nil {
		panic(fmt.Sprintf("create src dir failed: %v", err))
	}

	filePath = filepath.Join(cr.tmpDir, "dest", "file-1")
	err = os.WriteFile(filePath, nil, 0644)
	if err != nil {
		panic(fmt.Sprintf("Failed to create file: %v", err))
	}
	cr.T().Log("File created:", filePath)

}

func (cr *CopyRange) TearDownSuite() {
	cr.T().Log("teardown suite")
	_ = os.RemoveAll(filepath.Dir(cr.tmpDir))
}

func (cr *CopyRange) SetupTest() {
	cr.T().Log("setup test")
	recorder := httptest.NewRecorder()
	ctx := GetTestGinContext(recorder)
	ctx.Set(logging.LoggerName, cr.logger)
	cr.ctx = ctx
	cr.recorder = recorder
}

func (cr *CopyRange) TearDownTest() {
	cr.T().Log("teardown test")
}

func TestCopyRange(t *testing.T) {
	suite.Run(t, new(CopyRange))
}

func (cr *CopyRange) TestCopyRangeSuccessDir() {
	req := copyRange.Request{
		SrcPath:  cr.srcDir + "/file-1",
		DestPath: cr.destDir + "/file-1",
		Length:   1000,
	}
	err := BindJsonRequest(cr.ctx, req)
	if err != nil {
		panic(fmt.Sprintf("bind request failed: %v", err))
	}
	cr.storage.CopyRange(cr.ctx)

	resp := &copyRange.Response{}
	err = ParseResp(cr.recorder.Result(), resp)
	if !assert.Nil(cr.T(), err) {
		cr.T().Logf("%v: %v", time.Now().Format("2006-01-02 15:04:05"), err)
	} else {
		spew.Dump(resp)
	}
}

// ------------------------ error case ------------------------

func (cr *CopyRange) TestCopyRangeInvalidSrcPath() {
	req := copyRange.Request{
		SrcPath:  cr.srcDir + "..file-1",
		DestPath: cr.destDir + "/file-1",
		Length:   1000,
	}
	err := BindJsonRequest(cr.ctx, req)
	if err != nil {
		panic(fmt.Sprintf("bind request failed: %v", err))
	}
	cr.storage.CopyRange(cr.ctx)

	resp := &copyRange.Response{}
	_ = ParseResp(cr.recorder.Result(), resp)
	if assert.NotNil(cr.T(), resp) {
		assert.Equal(cr.T(), cr.recorder.Code, http.StatusBadRequest)
		assert.Equal(cr.T(), resp.ErrorCode, common.InvalidPath)
	}
}

func (cr *CopyRange) TestCopyRangeInvalidDestPath() {
	req := copyRange.Request{
		SrcPath:  cr.srcDir + "/file-1",
		DestPath: cr.destDir + "..file-1",
		Length:   1000,
	}
	err := BindJsonRequest(cr.ctx, req)
	if err != nil {
		panic(fmt.Sprintf("bind request failed: %v", err))
	}
	cr.storage.CopyRange(cr.ctx)

	resp := &copyRange.Response{}
	_ = ParseResp(cr.recorder.Result(), resp)
	if assert.NotNil(cr.T(), resp) {
		assert.Equal(cr.T(), cr.recorder.Code, http.StatusBadRequest)
		assert.Equal(cr.T(), resp.ErrorCode, common.InvalidPath)
	}
}

func (cr *CopyRange) TestCopyRangeSrcNotFound() {
	req := copyRange.Request{
		SrcPath:  cr.srcDir + "/not-exist",
		DestPath: cr.destDir + "file-1",
		Length:   1000,
	}
	err := BindJsonRequest(cr.ctx, req)
	if err != nil {
		panic(fmt.Sprintf("bind request failed: %v", err))
	}
	cr.storage.CopyRange(cr.ctx)

	resp := &copyRange.Response{}
	_ = ParseResp(cr.recorder.Result(), resp)
	if assert.NotNil(cr.T(), resp) {
		assert.Equal(cr.T(), cr.recorder.Code, http.StatusNotFound)
		assert.Equal(cr.T(), resp.ErrorCode, common.SrcPathNotFound)
	}
}

func (cr *CopyRange) TestCopyRangeIDestNotFound() {
	req := copyRange.Request{
		SrcPath:  cr.srcDir + "/file-1",
		DestPath: cr.destDir + "file-3",
		Length:   1000,
	}
	err := BindJsonRequest(cr.ctx, req)
	if err != nil {
		panic(fmt.Sprintf("bind request failed: %v", err))
	}
	cr.storage.CopyRange(cr.ctx)

	resp := &copyRange.Response{}
	_ = ParseResp(cr.recorder.Result(), resp)
	if assert.NotNil(cr.T(), resp) {
		assert.Equal(cr.T(), cr.recorder.Code, http.StatusNotFound)
		assert.Equal(cr.T(), resp.ErrorCode, common.DestPathNotFound)
	}
}

func (cr *CopyRange) TestCopyRangeSrcIsDir() {
	req := copyRange.Request{
		SrcPath:  cr.srcDir,
		DestPath: cr.destDir + "file-1",
		Length:   1000,
	}
	err := BindJsonRequest(cr.ctx, req)
	if err != nil {
		panic(fmt.Sprintf("bind request failed: %v", err))
	}
	cr.storage.CopyRange(cr.ctx)

	resp := &copyRange.Response{}
	_ = ParseResp(cr.recorder.Result(), resp)
	if assert.NotNil(cr.T(), resp) {
		assert.Equal(cr.T(), cr.recorder.Code, http.StatusBadRequest)
		assert.Equal(cr.T(), resp.ErrorCode, common.InvalidPath)
	}
}

func (cr *CopyRange) TestCopyRangeDestIsDir() {
	req := copyRange.Request{
		SrcPath:  cr.srcDir + "/file-1",
		DestPath: cr.srcDir,
		Length:   1000,
	}
	err := BindJsonRequest(cr.ctx, req)
	if err != nil {
		panic(fmt.Sprintf("bind request failed: %v", err))
	}
	cr.storage.CopyRange(cr.ctx)

	resp := &copyRange.Response{}
	_ = ParseResp(cr.recorder.Result(), resp)
	if assert.NotNil(cr.T(), resp) {
		assert.Equal(cr.T(), cr.recorder.Code, http.StatusBadRequest)
		assert.Equal(cr.T(), resp.ErrorCode, common.InvalidPath)
	}
}

func (cr *CopyRange) TestCopyRangeInvalidLength() {
	req := copyRange.Request{
		SrcPath:  cr.srcDir + "/file-1",
		DestPath: cr.destDir + "/file-1",
		Length:   -1,
	}
	err := BindJsonRequest(cr.ctx, req)
	if err != nil {
		panic(fmt.Sprintf("bind request failed: %v", err))
	}
	cr.storage.CopyRange(cr.ctx)

	resp := &copyRange.Response{}
	_ = ParseResp(cr.recorder.Result(), resp)
	if assert.NotNil(cr.T(), resp) {
		assert.Equal(cr.T(), cr.recorder.Code, http.StatusBadRequest)
		assert.Equal(cr.T(), resp.ErrorCode, common.InvalidLength)
	}
}

func (cr *CopyRange) TestCopyRangeInvalidSrcOffset() {
	req := copyRange.Request{
		SrcPath:   cr.srcDir + "/file-1",
		DestPath:  cr.destDir + "/file-1",
		Length:    0,
		SrcOffset: -1,
	}
	err := BindJsonRequest(cr.ctx, req)
	if err != nil {
		panic(fmt.Sprintf("bind request failed: %v", err))
	}
	cr.storage.CopyRange(cr.ctx)

	resp := &copyRange.Response{}
	_ = ParseResp(cr.recorder.Result(), resp)
	if assert.NotNil(cr.T(), resp) {
		assert.Equal(cr.T(), cr.recorder.Code, http.StatusBadRequest)
		assert.Equal(cr.T(), resp.ErrorCode, common.InvalidSrcOffset)
	}
}

func (cr *CopyRange) TestCopyRangeInvalidDestOffset() {
	req := copyRange.Request{
		SrcPath:    cr.srcDir + "/file-1",
		DestPath:   cr.destDir + "/file-1",
		Length:     0,
		DestOffset: -1,
	}
	err := BindJsonRequest(cr.ctx, req)
	if err != nil {
		panic(fmt.Sprintf("bind request failed: %v", err))
	}
	cr.storage.CopyRange(cr.ctx)

	resp := &copyRange.Response{}
	_ = ParseResp(cr.recorder.Result(), resp)
	if assert.NotNil(cr.T(), resp) {
		assert.Equal(cr.T(), cr.recorder.Code, http.StatusBadRequest)
		assert.Equal(cr.T(), resp.ErrorCode, common.InvalidDestOffset)
	}
}
