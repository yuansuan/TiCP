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
	"github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/ls"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao/model"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/filemode"
	v20230530 "github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/local/v20230530"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/local/v20230530/linker"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/common/mock"
	"xorm.io/xorm"
)

type Ls struct {
	suite.Suite
	logger              *logging.Logger
	storage             *v20230530.Storage
	tmpDir              string
	userID              string
	testDir             string
	fileNumber          int
	ctx                 *gin.Context
	recorder            *httptest.ResponseRecorder
	engine              *xorm.Engine
	mockStorageQuotaDao *dao.MockStorageQuotaDao
	testLinkDir         string
}

func (l *Ls) SetupSuite() {
	l.T().Log("setup suite")
	logger, err := logging.NewLogger(logging.WithReleaseLevel(logging.DevelopmentLevel))
	if err != nil {
		panic(fmt.Sprintf("init logger failed: %v", err))
	}

	ctrl := gomock.NewController(l.T())
	mockStorageQuotaDao := dao.NewMockStorageQuotaDao(ctrl)
	l.mockStorageQuotaDao = mockStorageQuotaDao
	l.mockStorageQuotaDao.EXPECT().Get(gomock.Any(), gomock.Any()).Return(true, &model.StorageQuota{StorageLimit: 1000}, nil).AnyTimes()
	mockEngine, _ := mock.Engine()
	mockEngine.SetLogger(middleware.NewXormLogger(l.logger, true))
	l.engine = mockEngine

	l.logger = logger
	currentDir, storage := CreateLocalStorage(l.userID, v20230530.SoftLink, l.mockStorageQuotaDao, nil, l.engine)
	l.storage = storage

	l.userID = "4TiSBX39DtN"
	l.testDir = "/test-ls"
	l.fileNumber = 5
	l.testLinkDir = "/test-link"

	tmpDir := currentDir + "/" + l.userID + l.testDir
	linkDir := currentDir + "/" + l.userID + l.testLinkDir
	err = os.MkdirAll(tmpDir, filemode.Directory)
	if err != nil {
		panic(fmt.Sprintf("create temp dir failed: %v", err))
	}
	l.tmpDir = tmpDir
	for i := 0; i < l.fileNumber; i++ {
		filePath := filepath.Join(l.tmpDir, "file-"+fmt.Sprintf("%d", i+1))
		err := os.WriteFile(filePath, nil, 0644)
		if err != nil {
			panic(fmt.Sprintf("Failed to create file: %v", err))
		}
		l.T().Log("File created:", filePath)
	}

	sl := &linker.SoftLink{}
	err = sl.Link(l.tmpDir, linkDir)
}

func (l *Ls) TearDownSuite() {
	l.T().Log("teardown suite")
	_ = os.RemoveAll(filepath.Dir(l.tmpDir))
}

func (l *Ls) SetupTest() {
	l.T().Log("setup test")
	recorder := httptest.NewRecorder()
	ctx := GetTestGinContext(recorder)
	ctx.Set(logging.LoggerName, l.logger)
	l.ctx = ctx
	l.recorder = recorder
}

func (l *Ls) TearDownTest() {
	l.T().Log("teardown test")
}

func TestLs(t *testing.T) {
	suite.Run(t, new(Ls))
}

func (l *Ls) TestLsSuccess() {
	req := ls.Request{
		Path:     "/" + l.userID + l.testDir,
		PageSize: 100,
	}
	BindRequest(l.ctx, req)
	l.storage.Ls(l.ctx)

	resp := &ls.Response{}
	err := ParseResp(l.recorder.Result(), resp)
	if !assert.Nil(l.T(), err) {
		l.T().Logf("%v: %v", time.Now().Format("2006-01-02 15:04:05"), err)
	} else {
		spew.Dump(resp.Data)
	}
}

func (l *Ls) TestLsLinkDirSuccess() {
	req := ls.Request{
		Path:     "/" + l.userID + l.testLinkDir,
		PageSize: 100,
	}
	BindRequest(l.ctx, req)
	l.storage.Ls(l.ctx)

	resp := &ls.Response{}
	err := ParseResp(l.recorder.Result(), resp)
	if !assert.Nil(l.T(), err) {
		l.T().Logf("%v: %v", time.Now().Format("2006-01-02 15:04:05"), err)
	} else {
		spew.Dump(resp.Data)
	}
}

// ------------------------ error case ------------------------

func (l *Ls) TestLsInvalidPageOffset() {
	req := ls.Request{
		Path:       "/" + l.userID + l.testDir,
		PageSize:   0,
		PageOffset: -1,
	}
	BindRequest(l.ctx, req)
	l.storage.Ls(l.ctx)

	resp := &ls.Response{}
	_ = ParseResp(l.recorder.Result(), resp)
	if assert.NotNil(l.T(), resp) {
		assert.Equal(l.T(), l.recorder.Code, http.StatusBadRequest)
		assert.Equal(l.T(), resp.ErrorCode, common.InvalidPageOffset)
	}
}

func (l *Ls) TestLsInvalidPageSize() {
	req := ls.Request{
		Path:     "/" + l.userID + l.testDir,
		PageSize: 0,
	}
	BindRequest(l.ctx, req)
	l.storage.Ls(l.ctx)

	resp := &ls.Response{}
	_ = ParseResp(l.recorder.Result(), resp)
	if assert.NotNil(l.T(), resp) {
		assert.Equal(l.T(), l.recorder.Code, http.StatusBadRequest)
		assert.Equal(l.T(), resp.ErrorCode, common.InvalidPageSize)
	}
}

func (l *Ls) TestLsInvalidPath() {
	req := ls.Request{
		Path:     l.userID + l.testDir,
		PageSize: 100,
	}
	BindRequest(l.ctx, req)
	l.storage.Ls(l.ctx)

	resp := &ls.Response{}
	_ = ParseResp(l.recorder.Result(), resp)
	if assert.NotNil(l.T(), resp) {
		assert.Equal(l.T(), l.recorder.Code, http.StatusBadRequest)
		assert.Equal(l.T(), resp.ErrorCode, common.InvalidPath)
	}
}

func (l *Ls) TestLsInvalidRegex() {
	req := ls.Request{
		Path:         "/" + l.userID + l.testDir,
		PageSize:     100,
		FilterRegexp: "[",
	}
	BindRequest(l.ctx, req)
	l.storage.Ls(l.ctx)

	resp := &ls.Response{}
	_ = ParseResp(l.recorder.Result(), resp)
	if assert.NotNil(l.T(), resp) {
		assert.Equal(l.T(), l.recorder.Code, http.StatusBadRequest)
		assert.Equal(l.T(), resp.ErrorCode, common.InvalidRegexp)
	}
}

func (l *Ls) TestLsPathNotFound() {
	req := ls.Request{
		Path:     "/" + l.userID + "/not-exist",
		PageSize: 100,
	}
	BindRequest(l.ctx, req)
	l.storage.Ls(l.ctx)

	resp := &ls.Response{}
	_ = ParseResp(l.recorder.Result(), resp)
	if assert.NotNil(l.T(), resp) {
		assert.Equal(l.T(), l.recorder.Code, http.StatusNotFound)
		assert.Equal(l.T(), resp.ErrorCode, common.PathNotFound)
	}
}

func (l *Ls) TestLsPathIsFile() {
	req := ls.Request{
		Path:     l.userID + l.testDir + "/file-1",
		PageSize: 100,
	}
	BindRequest(l.ctx, req)
	l.storage.Ls(l.ctx)

	resp := &ls.Response{}
	_ = ParseResp(l.recorder.Result(), resp)
	if assert.NotNil(l.T(), resp) {
		assert.Equal(l.T(), l.recorder.Code, http.StatusBadRequest)
		assert.Equal(l.T(), resp.ErrorCode, common.InvalidPath)
	}
}
