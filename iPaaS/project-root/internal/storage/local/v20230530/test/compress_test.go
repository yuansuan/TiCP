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
	"github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/compress/start"
	"github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/compress/status"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao/model"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/filemode"
	v20230530 "github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/local/v20230530"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/local/v20230530/linker"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/common/mock"
	"xorm.io/xorm"
)

type Compress struct {
	suite.Suite
	logger                     *logging.Logger
	storage                    *v20230530.Storage
	tmpDir                     string
	userID                     string
	testDir                    string
	curDir                     string
	fileNumber                 int
	ctx                        *gin.Context
	recorder                   *httptest.ResponseRecorder
	engine                     *xorm.Engine
	mockStorageQuotaDao        *dao.MockStorageQuotaDao
	mockStorageOperationLogDao *dao.MockStorageOperationLogDao
	testLinkDir                string
}

func (com *Compress) SetupSuite() {
	com.T().Log("setup suite")
	logger, err := logging.NewLogger(logging.WithReleaseLevel(logging.DevelopmentLevel))
	if err != nil {
		panic(fmt.Sprintf("init logger failed: %v", err))
	}

	ctrl := gomock.NewController(com.T())
	mockStorageQuotaDao := dao.NewMockStorageQuotaDao(ctrl)
	com.mockStorageQuotaDao = mockStorageQuotaDao
	com.mockStorageQuotaDao.EXPECT().Get(gomock.Any(), gomock.Any()).Return(true, &model.StorageQuota{StorageLimit: 1000}, nil).AnyTimes()
	mockStorageOperationLogDao := dao.NewMockStorageOperationLogDao(ctrl)
	com.mockStorageOperationLogDao = mockStorageOperationLogDao
	com.mockStorageOperationLogDao.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mockEngine, _ := mock.Engine()
	mockEngine.SetLogger(middleware.NewXormLogger(com.logger, true))
	com.engine = mockEngine

	com.logger = logger
	com.userID = "4TiSBX39DtN"
	currentDir, storage := CreateLocalStorage(com.userID, v20230530.SoftLink, com.mockStorageQuotaDao, com.mockStorageOperationLogDao, com.engine)
	com.storage = storage
	com.curDir = currentDir

	com.testDir = "/test-compress"
	com.fileNumber = 5
	com.testLinkDir = "/test-link"

	tmpDir := currentDir + "/" + com.userID + com.testDir
	linkDir := currentDir + "/" + com.userID + com.testLinkDir
	err = os.MkdirAll(tmpDir, filemode.Directory)
	if err != nil {
		panic(fmt.Sprintf("create temp dir failed: %v", err))
	}
	com.tmpDir = tmpDir
	for i := 0; i < com.fileNumber; i++ {
		filePath := filepath.Join(com.tmpDir, "file-"+fmt.Sprintf("%d", i+1))
		err := os.WriteFile(filePath, nil, 0644)
		if err != nil {
			panic(fmt.Sprintf("Failed to create file: %v", err))
		}
		com.T().Log("File created:", filePath)
	}
	filePath := filepath.Join(com.tmpDir, "file.zip")
	err = os.WriteFile(filePath, nil, 0644)
	if err != nil {
		panic(fmt.Sprintf("Failed to create file: %v", err))
	}
	com.T().Log("File created:", filePath)

	// softlink dir
	sl := &linker.SoftLink{}
	err = sl.Link(tmpDir, linkDir)
}

func (com *Compress) TearDownSuite() {
	com.T().Log("teardown suite")
	_ = os.RemoveAll(filepath.Dir(com.tmpDir))
	_ = os.RemoveAll(filepath.Join(com.curDir, ".compress"))
}

func (com *Compress) SetupTest() {
	com.T().Log("setup test")
	recorder := httptest.NewRecorder()
	ctx := GetTestGinContext(recorder)
	ctx.Set(logging.LoggerName, com.logger)
	com.ctx = ctx
	com.recorder = recorder
}

func (com *Compress) TearDownTest() {
	com.T().Log("teardown test")
}

func TestCompress(t *testing.T) {
	suite.Run(t, new(Compress))
}

func (com *Compress) TestCompressSuccess() {
	req := start.Request{
		Paths:      []string{"/" + com.userID + com.testDir + "/file-1"},
		TargetPath: "/" + com.userID + com.testDir + "/file-1.zip",
	}
	err := BindJsonRequest(com.ctx, req)
	if err != nil {
		panic(fmt.Sprintf("bind request failed: %v", err))
	}
	com.storage.CompressStart(com.ctx)

	resp := &start.Response{}
	err = ParseResp(com.recorder.Result(), resp)
	if !assert.Nil(com.T(), err) {
		com.T().Logf("%v: %v", time.Now().Format("2006-01-02 15:04:05"), err)
	} else {
		spew.Dump(resp)
		time.Sleep(3 * time.Second)
		req := status.Request{
			CompressID: resp.Data.CompressID,
		}
		recorder := httptest.NewRecorder()
		ctx := GetTestGinContext(recorder)
		ctx.Set(logging.LoggerName, com.logger)
		BindRequest(ctx, req)

		com.storage.CompressStatus(ctx)
		resp := &status.Response{}
		err = ParseResp(recorder.Result(), resp)
		if !assert.Nil(com.T(), err) {
			com.T().Logf("%v: %v", time.Now().Format("2006-01-02 15:04:05"), err)
		} else {
			spew.Dump(resp)
		}
	}
}

func (con *Compress) TestCompressLinkDirSuccess() {
	req := start.Request{
		Paths:      []string{"/" + con.userID + con.testLinkDir + "/file-1"},
		TargetPath: "/" + con.userID + con.testLinkDir + "/file-1.zip",
	}
	err := BindJsonRequest(con.ctx, req)
	if err != nil {
		panic(fmt.Sprintf("bind request failed: %v", err))
	}
	con.storage.CompressStart(con.ctx)

	resp := &start.Response{}
	err = ParseResp(con.recorder.Result(), resp)
	if !assert.Nil(con.T(), err) {
		con.T().Logf("%v: %v", time.Now().Format("2006-01-02 15:04:05"), err)
	} else {
		spew.Dump(resp)
		time.Sleep(3 * time.Second)
		req := status.Request{
			CompressID: resp.Data.CompressID,
		}
		recorder := httptest.NewRecorder()
		ctx := GetTestGinContext(recorder)
		ctx.Set(logging.LoggerName, con.logger)
		BindRequest(ctx, req)

		con.storage.CompressStatus(ctx)
		resp := &status.Response{}
		err = ParseResp(recorder.Result(), resp)
		if !assert.Nil(con.T(), err) {
			con.T().Logf("%v: %v", time.Now().Format("2006-01-02 15:04:05"), err)
		} else {
			spew.Dump(resp)
		}
	}
}

// ------------------------ error case ------------------------

func (com *Compress) TestCompressEmptyPaths() {
	req := start.Request{
		TargetPath: "/" + com.userID + com.testDir + "/file-1.zip",
	}
	err := BindJsonRequest(com.ctx, req)
	if err != nil {
		panic(fmt.Sprintf("bind request failed: %v", err))
	}
	com.storage.CompressStart(com.ctx)

	resp := &start.Response{}
	_ = ParseResp(com.recorder.Result(), resp)
	if assert.NotNil(com.T(), resp) {
		assert.Equal(com.T(), com.recorder.Code, http.StatusBadRequest)
		assert.Equal(com.T(), resp.ErrorCode, common.InvalidPath)
	}
}

func (com *Compress) TestCompressInvalidPaths() {
	req := start.Request{
		Paths:      []string{com.userID + com.testDir + "/file-1"},
		TargetPath: "/" + com.userID + com.testDir + "/file-1.zip",
	}
	err := BindJsonRequest(com.ctx, req)
	if err != nil {
		panic(fmt.Sprintf("bind request failed: %v", err))
	}
	com.storage.CompressStart(com.ctx)

	resp := &start.Response{}
	_ = ParseResp(com.recorder.Result(), resp)
	if assert.NotNil(com.T(), resp) {
		assert.Equal(com.T(), com.recorder.Code, http.StatusBadRequest)
		assert.Equal(com.T(), resp.ErrorCode, common.InvalidPath)
	}
}

func (com *Compress) TestCompressInvalidTargetPath() {
	req := start.Request{
		Paths:      []string{"/" + com.userID + com.testDir + "/file-1"},
		TargetPath: com.userID + com.testDir + "/file-1.zip",
	}
	err := BindJsonRequest(com.ctx, req)
	if err != nil {
		panic(fmt.Sprintf("bind request failed: %v", err))
	}
	com.storage.CompressStart(com.ctx)

	resp := &start.Response{}
	_ = ParseResp(com.recorder.Result(), resp)
	if assert.NotNil(com.T(), resp) {
		assert.Equal(com.T(), com.recorder.Code, http.StatusBadRequest)
		assert.Equal(com.T(), resp.ErrorCode, common.InvalidTargetPath)
	}
}

func (com *Compress) TestCompressInvalidCompressFormat() {
	req := start.Request{
		Paths:      []string{"/" + com.userID + com.testDir + "/file-1"},
		TargetPath: "/" + com.userID + com.testDir + "/file-1.tar",
	}
	err := BindJsonRequest(com.ctx, req)
	if err != nil {
		panic(fmt.Sprintf("bind request failed: %v", err))
	}
	com.storage.CompressStart(com.ctx)

	resp := &start.Response{}
	_ = ParseResp(com.recorder.Result(), resp)
	if assert.NotNil(com.T(), resp) {
		assert.Equal(com.T(), com.recorder.Code, http.StatusBadRequest)
		assert.Equal(com.T(), resp.ErrorCode, common.UnsupportedCompressFileType)
	}
}

func (com *Compress) TestCompressInvalidTargetFileExist() {
	req := start.Request{
		Paths:      []string{"/" + com.userID + com.testDir + "/file-1"},
		TargetPath: "/" + com.userID + com.testDir + "/file.zip",
	}
	err := BindJsonRequest(com.ctx, req)
	if err != nil {
		panic(fmt.Sprintf("bind request failed: %v", err))
	}
	com.storage.CompressStart(com.ctx)

	resp := &start.Response{}
	_ = ParseResp(com.recorder.Result(), resp)
	if assert.NotNil(com.T(), resp) {
		assert.Equal(com.T(), com.recorder.Code, http.StatusBadRequest)
		assert.Equal(com.T(), resp.ErrorCode, common.TargetPathExists)
	}
}

func (com *Compress) TestCompressPathNotFound() {
	req := start.Request{
		Paths:      []string{"/" + com.userID + com.testDir + "/not-exist"},
		TargetPath: "/" + com.userID + com.testDir + "/file-1.zip",
	}
	err := BindJsonRequest(com.ctx, req)
	if err != nil {
		panic(fmt.Sprintf("bind request failed: %v", err))
	}
	com.storage.CompressStart(com.ctx)

	resp := &start.Response{}
	_ = ParseResp(com.recorder.Result(), resp)
	if assert.NotNil(com.T(), resp) {
		assert.Equal(com.T(), com.recorder.Code, http.StatusNotFound)
		assert.Equal(com.T(), resp.ErrorCode, common.PathNotFound)
	}
}

func (com *Compress) TestCompressCompressIDEmpty() {

	req := status.Request{
		CompressID: "",
	}
	BindRequest(com.ctx, req)

	com.storage.CompressStatus(com.ctx)
	resp := &status.Response{}
	_ = ParseResp(com.recorder.Result(), resp)
	if assert.NotNil(com.T(), resp) {
		assert.Equal(com.T(), com.recorder.Code, http.StatusBadRequest)
		assert.Equal(com.T(), resp.ErrorCode, common.InvalidCompressID)
	}
}

func (com *Compress) TestCompressTaskNotFound() {

	req := status.Request{
		CompressID: com.userID,
	}
	BindRequest(com.ctx, req)

	com.storage.CompressStatus(com.ctx)
	resp := &status.Response{}
	_ = ParseResp(com.recorder.Result(), resp)
	if assert.NotNil(com.T(), resp) {
		assert.Equal(com.T(), com.recorder.Code, http.StatusNotFound)
		assert.Equal(com.T(), resp.ErrorCode, common.CompressTaskNotFound)
	}
}
