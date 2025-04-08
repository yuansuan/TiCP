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
	linkapi "github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/link"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao/model"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/filemode"
	v20230530 "github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/local/v20230530"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/common/mock"
	"xorm.io/xorm"
)

type Link struct {
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

func (link *Link) SetupSuite() {
	link.T().Log("setup suite")
	logger, err := logging.NewLogger(logging.WithReleaseLevel(logging.DevelopmentLevel))
	if err != nil {
		panic(fmt.Sprintf("init logger failed: %v", err))
	}

	ctrl := gomock.NewController(link.T())
	mockStorageQuotaDao := dao.NewMockStorageQuotaDao(ctrl)
	link.mockStorageQuotaDao = mockStorageQuotaDao
	link.mockStorageQuotaDao.EXPECT().Get(gomock.Any(), gomock.Any()).Return(true, &model.StorageQuota{StorageLimit: 1000}, nil).AnyTimes()
	mockStorageOperationLogDao := dao.NewMockStorageOperationLogDao(ctrl)
	link.mockStorageOperationLogDao = mockStorageOperationLogDao
	link.mockStorageOperationLogDao.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mockEngine, _ := mock.Engine()
	mockEngine.SetLogger(middleware.NewXormLogger(link.logger, true))
	link.engine = mockEngine

	link.logger = logger
	currentDir, storage := CreateLocalStorage(link.userID, v20230530.SoftLink, link.mockStorageQuotaDao, link.mockStorageOperationLogDao, link.engine)
	link.storage = storage

	link.userID = "4TiSBX39DtN"
	link.testDir = "/test-link"
	link.srcDir = "/" + link.userID + link.testDir + "/src"
	link.destDir = "/" + link.userID + link.testDir + "/dest"

	tmpDir := currentDir + "/" + link.userID + link.testDir
	err = os.MkdirAll(tmpDir, filemode.Directory)
	if err != nil {
		panic(fmt.Sprintf("create temp dir failed: %v", err))
	}

	link.tmpDir = tmpDir
	err = os.MkdirAll(filepath.Join(link.tmpDir, "src"), filemode.Directory)
	if err != nil {
		panic(fmt.Sprintf("create src dir failed: %v", err))
	}

	filePath := filepath.Join(link.tmpDir, "src", "file-1")
	err = os.WriteFile(filePath, nil, 0644)
	if err != nil {
		panic(fmt.Sprintf("Failed to create file: %v", err))
	}
	link.T().Log("File created:", filePath)

}

func (link *Link) TearDownSuite() {
	link.T().Log("teardown suite")
	_ = os.RemoveAll(filepath.Dir(link.tmpDir))
}

func (link *Link) SetupTest() {
	link.T().Log("setup test")
	recorder := httptest.NewRecorder()
	ctx := GetTestGinContext(recorder)
	ctx.Set(logging.LoggerName, link.logger)
	link.ctx = ctx
	link.recorder = recorder
}

func (link *Link) TearDownTest() {
	link.T().Log("teardown test")
}

func TestLink(t *testing.T) {
	suite.Run(t, new(Link))
}

func (link *Link) TestSoftLinkDirSuccess() {
	req := linkapi.Request{
		SrcPath:  link.srcDir,
		DestPath: link.destDir + "/dir-1",
	}
	_, storage := CreateLocalStorage(link.userID, v20230530.SoftLink, link.mockStorageQuotaDao, link.mockStorageOperationLogDao, link.engine)
	recorder := httptest.NewRecorder()
	ctx := GetTestGinContext(recorder)
	ctx.Set(logging.LoggerName, link.logger)
	err := BindJsonRequest(ctx, req)
	if err != nil {
		panic(fmt.Sprintf("bind request failed: %v", err))
	}
	storage.Link(ctx)

	resp := &linkapi.Response{}
	err = ParseResp(recorder.Result(), resp)
	if !assert.Nil(link.T(), err) {
		link.T().Logf("%v: %v", time.Now().Format("2006-01-02 15:04:05"), err)
	} else {
		spew.Dump(resp)
	}
}

func (link *Link) TestSoftLinkFileSuccess() {
	req := linkapi.Request{
		SrcPath:  link.srcDir + "/file-1",
		DestPath: link.destDir + "/soft-link",
	}
	err := BindJsonRequest(link.ctx, req)
	if err != nil {
		panic(fmt.Sprintf("bind request failed: %v", err))
	}
	link.storage.Link(link.ctx)

	resp := &linkapi.Response{}
	err = ParseResp(link.recorder.Result(), resp)
	if !assert.Nil(link.T(), err) {
		link.T().Logf("%v: %v", time.Now().Format("2006-01-02 15:04:05"), err)
	} else {
		spew.Dump(resp)
	}
}

func (link *Link) TestHardLinkFileSuccess() {
	req := linkapi.Request{
		SrcPath:  link.srcDir + "/file-1",
		DestPath: link.destDir + "/hard-link",
	}
	_, storage := CreateLocalStorage(link.userID, v20230530.HardLink, link.mockStorageQuotaDao, link.mockStorageOperationLogDao, link.engine)
	recorder := httptest.NewRecorder()
	ctx := GetTestGinContext(recorder)
	ctx.Set(logging.LoggerName, link.logger)
	err := BindJsonRequest(ctx, req)
	if err != nil {
		panic(fmt.Sprintf("bind request failed: %v", err))
	}
	storage.Link(ctx)

	resp := &linkapi.Response{}
	err = ParseResp(recorder.Result(), resp)
	if !assert.Nil(link.T(), err) {
		link.T().Logf("%v: %v", time.Now().Format("2006-01-02 15:04:05"), err)
	} else {
		spew.Dump(resp)
	}
}

func (link *Link) TestHardLinkDirSuccess() {
	req := linkapi.Request{
		SrcPath:  link.srcDir,
		DestPath: link.destDir + "/hard-link-2",
	}
	_, storage := CreateLocalStorage(link.userID, v20230530.HardLink, link.mockStorageQuotaDao, link.mockStorageOperationLogDao, link.engine)
	recorder := httptest.NewRecorder()
	ctx := GetTestGinContext(recorder)
	ctx.Set(logging.LoggerName, link.logger)
	err := BindJsonRequest(ctx, req)
	if err != nil {
		panic(fmt.Sprintf("bind request failed: %v", err))
	}
	storage.Link(ctx)

	resp := &linkapi.Response{}
	err = ParseResp(recorder.Result(), resp)
	if !assert.Nil(link.T(), err) {
		link.T().Logf("%v: %v", time.Now().Format("2006-01-02 15:04:05"), err)
	} else {
		spew.Dump(resp)
	}
}

// ------------------------ error case ------------------------

func (link *Link) TestLinkInvalidSrcPath() {
	req := linkapi.Request{
		SrcPath:  link.srcDir + "../file",
		DestPath: link.destDir,
	}
	err := BindJsonRequest(link.ctx, req)
	if err != nil {
		panic(fmt.Sprintf("bind request failed: %v", err))
	}
	link.storage.Link(link.ctx)

	resp := &linkapi.Response{}
	_ = ParseResp(link.recorder.Result(), resp)
	if assert.NotNil(link.T(), resp) {
		assert.Equal(link.T(), link.recorder.Code, http.StatusBadRequest)
		assert.Equal(link.T(), resp.ErrorCode, common.InvalidPath)
	}
}

func (link *Link) TestLinkInvalidDestPath() {
	req := linkapi.Request{
		SrcPath:  link.srcDir,
		DestPath: link.destDir + "../file",
	}
	err := BindJsonRequest(link.ctx, req)
	if err != nil {
		panic(fmt.Sprintf("bind request failed: %v", err))
	}
	link.storage.Link(link.ctx)

	resp := &linkapi.Response{}
	_ = ParseResp(link.recorder.Result(), resp)
	if assert.NotNil(link.T(), resp) {
		assert.Equal(link.T(), link.recorder.Code, http.StatusBadRequest)
		assert.Equal(link.T(), resp.ErrorCode, common.InvalidPath)
	}
}

func (link *Link) TestLinkSrcNotFound() {
	req := linkapi.Request{
		SrcPath:  link.srcDir + "/not-exist",
		DestPath: link.destDir,
	}
	err := BindJsonRequest(link.ctx, req)
	if err != nil {
		panic(fmt.Sprintf("bind request failed: %v", err))
	}
	link.storage.Link(link.ctx)

	resp := &linkapi.Response{}
	_ = ParseResp(link.recorder.Result(), resp)
	if assert.NotNil(link.T(), resp) {
		assert.Equal(link.T(), link.recorder.Code, http.StatusNotFound)
		assert.Equal(link.T(), resp.ErrorCode, common.SrcPathNotFound)
	}
}

func (link *Link) TestLinkDestExist() {
	req := linkapi.Request{
		SrcPath:  link.srcDir,
		DestPath: link.srcDir,
	}
	_, storage := CreateLocalStorage(link.userID, v20230530.SoftLink, link.mockStorageQuotaDao, nil, link.engine)
	recorder := httptest.NewRecorder()
	ctx := GetTestGinContext(recorder)
	ctx.Set(logging.LoggerName, link.logger)
	err := BindJsonRequest(ctx, req)
	if err != nil {
		panic(fmt.Sprintf("bind request failed: %v", err))
	}
	storage.Link(ctx)

	resp := &linkapi.Response{}
	_ = ParseResp(recorder.Result(), resp)
	if assert.NotNil(link.T(), resp) {
		assert.Equal(link.T(), recorder.Code, http.StatusBadRequest)
		assert.Equal(link.T(), resp.ErrorCode, common.DestPathExists)
	}
}
