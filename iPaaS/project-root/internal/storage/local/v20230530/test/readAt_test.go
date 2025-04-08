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
	readAtAPI "github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/readAt"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao/model"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/filemode"
	v20230530 "github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/local/v20230530"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/common/mock"
	"xorm.io/xorm"
)

type ReadAt struct {
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

func (read *ReadAt) SetupSuite() {
	read.T().Log("setup suite")
	logger, err := logging.NewLogger(logging.WithReleaseLevel(logging.DevelopmentLevel))
	if err != nil {
		panic(fmt.Sprintf("init logger failed: %v", err))
	}

	ctrl := gomock.NewController(read.T())
	mockStorageQuotaDao := dao.NewMockStorageQuotaDao(ctrl)
	read.mockStorageQuotaDao = mockStorageQuotaDao
	read.mockStorageQuotaDao.EXPECT().Get(gomock.Any(), gomock.Any()).Return(true, &model.StorageQuota{StorageLimit: 1000}, nil).AnyTimes()
	mockStorageOperationLogDao := dao.NewMockStorageOperationLogDao(ctrl)
	read.mockStorageOperationLogDao = mockStorageOperationLogDao
	read.mockStorageOperationLogDao.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mockEngine, _ := mock.Engine()
	mockEngine.SetLogger(middleware.NewXormLogger(read.logger, true))
	read.engine = mockEngine

	read.logger = logger
	currentDir, storage := CreateLocalStorage(read.userID, v20230530.SoftLink, read.mockStorageQuotaDao, read.mockStorageOperationLogDao, read.engine)
	read.storage = storage

	read.userID = "4TiSBX39DtN"
	read.testDir = "/test-read"

	tmpDir := currentDir + "/" + read.userID + read.testDir
	err = os.MkdirAll(tmpDir, filemode.Directory)
	if err != nil {
		panic(fmt.Sprintf("create temp dir failed: %v", err))
	}
	read.tmpDir = tmpDir
	filePath := filepath.Join(read.tmpDir, "file-1")
	data := make([]byte, 5*1024)
	rand.Read(data)
	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		panic(fmt.Sprintf("Failed to create file: %v", err))
	}
	read.T().Log("File created:", filePath)

}

func (read *ReadAt) TearDownSuite() {
	read.T().Log("teardown suite")
	_ = os.RemoveAll(filepath.Dir(read.tmpDir))
}

func (read *ReadAt) SetupTest() {
	read.T().Log("setup test")
	recorder := httptest.NewRecorder()
	ctx := GetTestGinContext(recorder)
	ctx.Set(logging.LoggerName, read.logger)
	read.ctx = ctx
	read.recorder = recorder
}

func (read *ReadAt) TearDownTest() {
	read.T().Log("teardown test")
}

func TestRead(t *testing.T) {
	suite.Run(t, new(ReadAt))
}

func (read *ReadAt) TestReadAtSuccessFile() {
	req := readAtAPI.Request{
		Path:       "/" + read.userID + read.testDir + "/file-1",
		Length:     100,
		Compressor: common.GZIP,
	}
	BindRequest(read.ctx, req)
	read.storage.ReadAt(read.ctx)

	resp := &readAtAPI.Response{}
	err := defaultReadResolver(read.recorder.Result(), resp)
	if !assert.Nil(read.T(), err) {
		read.T().Logf("%v: %v", time.Now().Format("2006-01-02 15:04:05"), err)
	}
}

// ------------------------ error case ------------------------

func (read *ReadAt) TestReadAtInvalidPath() {
	req := readAtAPI.Request{
		Path:   read.userID + read.testDir + "/file-1",
		Length: 100,
	}
	BindRequest(read.ctx, req)
	read.storage.ReadAt(read.ctx)

	resp := &readAtAPI.Response{}
	_ = defaultReadResolver(read.recorder.Result(), resp)
	if assert.NotNil(read.T(), resp) {
		assert.Equal(read.T(), read.recorder.Code, http.StatusBadRequest)
		assert.Contains(read.T(), string(resp.Data), common.InvalidPath)
	}
}

func (read *ReadAt) TestReadAtPathNotFound() {
	req := readAtAPI.Request{
		Path:   "/" + read.userID + read.testDir + "/not-exist",
		Length: 100,
	}
	BindRequest(read.ctx, req)
	read.storage.ReadAt(read.ctx)

	resp := &readAtAPI.Response{}
	_ = defaultReadResolver(read.recorder.Result(), resp)
	if assert.NotNil(read.T(), resp) {
		assert.Equal(read.T(), read.recorder.Code, http.StatusNotFound)
		assert.Contains(read.T(), string(resp.Data), common.PathNotFound)
	}
}

func (read *ReadAt) TestReadAtPathIsDir() {
	req := readAtAPI.Request{
		Path:   "/" + read.userID + read.testDir,
		Length: 100,
	}
	BindRequest(read.ctx, req)
	read.storage.ReadAt(read.ctx)

	resp := &readAtAPI.Response{}
	_ = defaultReadResolver(read.recorder.Result(), resp)
	if assert.NotNil(read.T(), resp) {
		assert.Equal(read.T(), read.recorder.Code, http.StatusBadRequest)
		assert.Contains(read.T(), string(resp.Data), common.InvalidPath)
	}
}

func (read *ReadAt) TestReadAtInvalidLength() {
	req := readAtAPI.Request{
		Path:   "/" + read.userID + read.testDir + "/file-1",
		Length: -1,
	}
	BindRequest(read.ctx, req)
	read.storage.ReadAt(read.ctx)

	resp := &readAtAPI.Response{}
	_ = defaultReadResolver(read.recorder.Result(), resp)
	if assert.NotNil(read.T(), resp) {
		assert.Equal(read.T(), read.recorder.Code, http.StatusBadRequest)
		assert.Contains(read.T(), string(resp.Data), common.InvalidLength)
	}
}

func (read *ReadAt) TestReadAtInvalidOffset() {
	req := readAtAPI.Request{
		Path:   "/" + read.userID + read.testDir + "/file-1",
		Length: 100,
		Offset: -1,
	}
	BindRequest(read.ctx, req)
	read.storage.ReadAt(read.ctx)

	resp := &readAtAPI.Response{}
	_ = defaultReadResolver(read.recorder.Result(), resp)
	if assert.NotNil(read.T(), resp) {
		assert.Equal(read.T(), read.recorder.Code, http.StatusBadRequest)
		assert.Contains(read.T(), string(resp.Data), common.InvalidOffset)
	}
}

func (read *ReadAt) TestReadAtInvalidCompressor() {
	req := readAtAPI.Request{
		Path:       "/" + read.userID + read.testDir + "/file-1",
		Length:     100,
		Compressor: "invalid",
	}
	BindRequest(read.ctx, req)
	read.storage.ReadAt(read.ctx)

	resp := &readAtAPI.Response{}
	_ = defaultReadResolver(read.recorder.Result(), resp)
	if assert.NotNil(read.T(), resp) {
		assert.Equal(read.T(), read.recorder.Code, http.StatusBadRequest)
		assert.Contains(read.T(), string(resp.Data), common.InvalidCompressor)
	}
}

func defaultReadResolver(resp *http.Response, ret *readAtAPI.Response) error {

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
