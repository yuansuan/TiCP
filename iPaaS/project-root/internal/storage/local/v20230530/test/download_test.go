package test

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/middleware"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/common/project-root-api/common"
	downloadAPI "github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/download"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao/model"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/filemode"
	v20230530 "github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/local/v20230530"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/local/v20230530/linker"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/common/mock"
	"xorm.io/xorm"
)

type Download struct {
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
	testLinkDir                string
}

func (download *Download) SetupSuite() {
	download.T().Log("setup suite")
	logger, err := logging.NewLogger(logging.WithReleaseLevel(logging.DevelopmentLevel))
	if err != nil {
		panic(fmt.Sprintf("init logger failed: %v", err))
	}

	ctrl := gomock.NewController(download.T())
	mockStorageQuotaDao := dao.NewMockStorageQuotaDao(ctrl)
	download.mockStorageQuotaDao = mockStorageQuotaDao
	download.mockStorageQuotaDao.EXPECT().Get(gomock.Any(), gomock.Any()).Return(true, &model.StorageQuota{StorageLimit: 1000}, nil).AnyTimes()
	mockStorageOperationLogDao := dao.NewMockStorageOperationLogDao(ctrl)
	download.mockStorageOperationLogDao = mockStorageOperationLogDao
	download.mockStorageOperationLogDao.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mockEngine, _ := mock.Engine()
	mockEngine.SetLogger(middleware.NewXormLogger(download.logger, true))
	download.engine = mockEngine

	download.logger = logger
	currentDir, storage := CreateLocalStorage(download.userID, v20230530.SoftLink, download.mockStorageQuotaDao, download.mockStorageOperationLogDao, download.engine)
	download.storage = storage

	download.userID = "4TiSBX39DtN"
	download.testDir = "/test-download"
	download.testLinkDir = "/test-link"

	tmpDir := currentDir + "/" + download.userID + download.testDir
	linkDir := currentDir + "/" + download.userID + download.testLinkDir
	err = os.MkdirAll(tmpDir, filemode.Directory)
	if err != nil {
		panic(fmt.Sprintf("create temp dir failed: %v", err))
	}
	download.tmpDir = tmpDir
	filePath := filepath.Join(download.tmpDir, "file-1")
	data := make([]byte, 5*1024)
	rand.Read(data)
	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		panic(fmt.Sprintf("Failed to create file: %v", err))
	}
	download.T().Log("File created:", filePath)

	linkPath := filepath.Join(download.tmpDir, "link")
	err = os.Symlink("file-1", linkPath)

	// softlink dir
	sl := &linker.SoftLink{}
	err = sl.Link(tmpDir, linkDir)
}

func (download *Download) TearDownSuite() {
	download.T().Log("teardown suite")
	_ = os.RemoveAll(filepath.Dir(download.tmpDir))
}

func (download *Download) SetupTest() {
	download.T().Log("setup test")
	recorder := httptest.NewRecorder()
	ctx := GetTestGinContext(recorder)
	ctx.Set(logging.LoggerName, download.logger)
	download.ctx = ctx
	download.recorder = recorder
}

func (download *Download) TearDownTest() {
	download.T().Log("teardown test")
}

func TestDownload(t *testing.T) {
	suite.Run(t, new(Download))
}

func (download *Download) TestDownloadSuccessFile() {
	req := downloadAPI.Request{
		Path: "/" + download.userID + download.testDir + "/file-1",
	}
	BindRequest(download.ctx, req)
	download.storage.Download(download.ctx)

	resp := &downloadAPI.Response{}
	err := defaultResolver(download.recorder.Result(), resp)
	if !assert.Nil(download.T(), err) {
		download.T().Logf("%v: %v", time.Now().Format("2006-01-02 15:04:05"), err)
	}
}

func (download *Download) TestDownloadSuccessLinkDirFile() {
	req := downloadAPI.Request{
		Path: "/" + download.userID + download.testLinkDir + "/file-1",
	}
	BindRequest(download.ctx, req)
	download.storage.Download(download.ctx)

	resp := &downloadAPI.Response{}
	err := defaultResolver(download.recorder.Result(), resp)
	if !assert.Nil(download.T(), err) {
		download.T().Logf("%v: %v", time.Now().Format("2006-01-02 15:04:05"), err)
	}
}

func (download *Download) TestDownloadSuccessLink() {
	req := downloadAPI.Request{
		Path: "/" + download.userID + download.testDir + "/link",
	}
	BindRequest(download.ctx, req)
	download.storage.Download(download.ctx)

	resp := &downloadAPI.Response{}
	err := defaultResolver(download.recorder.Result(), resp)
	if !assert.Nil(download.T(), err) {
		download.T().Logf("%v: %v", time.Now().Format("2006-01-02 15:04:05"), err)
	}
}

// ------------------------ error case ------------------------

func (download *Download) TestDownloadInvalidPath() {
	req := downloadAPI.Request{
		Path: download.userID + download.testDir + "/file-1",
	}
	BindRequest(download.ctx, req)
	download.storage.Download(download.ctx)

	resp := &downloadAPI.Response{}
	_ = defaultResolver(download.recorder.Result(), resp)
	if assert.NotNil(download.T(), resp) {
		assert.Equal(download.T(), download.recorder.Code, http.StatusBadRequest)
		assert.Contains(download.T(), string(resp.Data), common.InvalidPath)
	}
}

func (download *Download) TestDownloadPathIsDir() {
	req := downloadAPI.Request{
		Path: "/" + download.userID + download.testDir,
	}
	BindRequest(download.ctx, req)
	download.storage.Download(download.ctx)

	resp := &downloadAPI.Response{}
	_ = defaultResolver(download.recorder.Result(), resp)
	if assert.NotNil(download.T(), resp) {
		assert.Equal(download.T(), download.recorder.Code, http.StatusBadRequest)
		assert.Contains(download.T(), string(resp.Data), common.InvalidPath)
	}
}

func (download *Download) TestDownloadPathNotFound() {
	req := downloadAPI.Request{
		Path: "/" + download.userID + download.testDir + "/not-exist",
	}
	BindRequest(download.ctx, req)
	download.storage.Download(download.ctx)

	resp := &downloadAPI.Response{}
	_ = defaultResolver(download.recorder.Result(), resp)
	if assert.NotNil(download.T(), resp) {
		assert.Equal(download.T(), download.recorder.Code, http.StatusNotFound)
		assert.Contains(download.T(), string(resp.Data), common.PathNotFound)
	}
}

func (download *Download) TestDownloadInvalidRange() {
	req := downloadAPI.Request{
		Path:  "/" + download.userID + download.testDir + "/file-1",
		Range: "bytes=0-0",
	}
	BindRequest(download.ctx, req)
	download.storage.Download(download.ctx)

	resp := &downloadAPI.Response{}
	_ = defaultResolver(download.recorder.Result(), resp)
	if assert.NotNil(download.T(), resp) {
		assert.Equal(download.T(), download.recorder.Code, http.StatusBadRequest)
		assert.Contains(download.T(), string(resp.Data), common.InvalidRange)
	}
}

func defaultResolver(resp *http.Response, ret *downloadAPI.Response) error {

	var err error
	ret.Data, err = io.ReadAll(resp.Body)
	if err != nil {
		return errors.Errorf("get data error: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusPartialContent {
		return errors.Errorf("http: %v, body: %v", resp.Status, string(ret.Data))
	}

	return err
}
