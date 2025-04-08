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
	move "github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/mv"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao/model"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/filemode"
	v20230530 "github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/local/v20230530"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/common/mock"
	"xorm.io/xorm"
)

type Move struct {
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

func (mv *Move) SetupSuite() {
	mv.T().Log("setup suite")
	logger, err := logging.NewLogger(logging.WithReleaseLevel(logging.DevelopmentLevel))
	if err != nil {
		panic(fmt.Sprintf("init logger failed: %v", err))
	}

	ctrl := gomock.NewController(mv.T())
	mockStorageQuotaDao := dao.NewMockStorageQuotaDao(ctrl)
	mv.mockStorageQuotaDao = mockStorageQuotaDao
	mv.mockStorageQuotaDao.EXPECT().Get(gomock.Any(), gomock.Any()).Return(true, &model.StorageQuota{StorageLimit: 1000}, nil).AnyTimes()
	mockStorageOperationLogDao := dao.NewMockStorageOperationLogDao(ctrl)
	mv.mockStorageOperationLogDao = mockStorageOperationLogDao
	mv.mockStorageOperationLogDao.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mockEngine, _ := mock.Engine()
	mockEngine.SetLogger(middleware.NewXormLogger(mv.logger, true))
	mv.engine = mockEngine

	mv.logger = logger
	currentDir, storage := CreateLocalStorage(mv.userID, v20230530.SoftLink, mv.mockStorageQuotaDao, mv.mockStorageOperationLogDao, mv.engine)
	mv.storage = storage

	mv.userID = "4TiSBX39DtN"
	mv.testDir = "/test-mv"

	tmpDir := currentDir + "/" + mv.userID + mv.testDir
	mv.tmpDir = tmpDir
	mv.srcDir = "/" + mv.userID + mv.testDir + "/src"
	mv.destDir = "/" + mv.userID + mv.testDir + "/dest"

	err = os.MkdirAll(filepath.Join(mv.tmpDir, "src"), filemode.Directory)
	if err != nil {
		panic(fmt.Sprintf("create src dir failed: %v", err))
	}

	filePath := filepath.Join(mv.tmpDir, "src", "file-1")
	err = os.WriteFile(filePath, nil, 0644)
	if err != nil {
		panic(fmt.Sprintf("Failed to create file: %v", err))
	}
	mv.T().Log("File created:", filePath)

}

func (mv *Move) TearDownSuite() {
	mv.T().Log("teardown suite")
	_ = os.RemoveAll(filepath.Dir(mv.tmpDir))
}

func (mv *Move) SetupTest() {
	mv.T().Log("setup test")
	recorder := httptest.NewRecorder()
	ctx := GetTestGinContext(recorder)
	ctx.Set(logging.LoggerName, mv.logger)
	mv.ctx = ctx
	mv.recorder = recorder
}

func (mv *Move) TearDownTest() {
	mv.T().Log("teardown test")
}

func TestMove(t *testing.T) {
	suite.Run(t, new(Move))
}

func (mv *Move) TestMoveSuccessFolder() {
	req := move.Request{
		SrcPath:  mv.srcDir,
		DestPath: mv.destDir,
	}
	err := BindJsonRequest(mv.ctx, req)
	if err != nil {
		panic(fmt.Sprintf("bind request failed: %v", err))
	}
	mv.storage.Mv(mv.ctx)

	resp := &move.Response{}
	err = ParseResp(mv.recorder.Result(), resp)
	if !assert.Nil(mv.T(), err) {
		mv.T().Logf("%v: %v", time.Now().Format("2006-01-02 15:04:05"), err)
	} else {
		spew.Dump(resp)
	}
}

func (mv *Move) TestMoveSuccessFile() {
	req := move.Request{
		SrcPath:  mv.srcDir + "/file-1",
		DestPath: "/" + mv.userID + mv.testDir + "/new-dir/file-1",
	}
	err := BindJsonRequest(mv.ctx, req)
	if err != nil {
		panic(fmt.Sprintf("bind request failed: %v", err))
	}
	mv.storage.Mv(mv.ctx)

	resp := &move.Response{}
	err = ParseResp(mv.recorder.Result(), resp)
	if !assert.Nil(mv.T(), err) {
		mv.T().Logf("%v: %v", time.Now().Format("2006-01-02 15:04:05"), err)
	} else {
		spew.Dump(resp)
	}
}

// ------------------------ error case ------------------------

func (mv *Move) TestMoveInvalidSrcPath() {
	req := move.Request{
		SrcPath:  mv.srcDir + "../file-1",
		DestPath: mv.destDir,
	}
	err := BindJsonRequest(mv.ctx, req)
	if err != nil {
		panic(fmt.Sprintf("bind request failed: %v", err))
	}
	mv.storage.Mv(mv.ctx)

	resp := &move.Response{}
	_ = ParseResp(mv.recorder.Result(), resp)
	if assert.NotNil(mv.T(), resp) {
		assert.Equal(mv.T(), mv.recorder.Code, http.StatusBadRequest)
		assert.Equal(mv.T(), resp.ErrorCode, common.InvalidPath)
	}
}

func (mv *Move) TestMoveInvalidDestPath() {
	req := move.Request{
		SrcPath:  mv.srcDir,
		DestPath: mv.destDir + "../file-1",
	}
	err := BindJsonRequest(mv.ctx, req)
	if err != nil {
		panic(fmt.Sprintf("bind request failed: %v", err))
	}
	mv.storage.Mv(mv.ctx)

	resp := &move.Response{}
	_ = ParseResp(mv.recorder.Result(), resp)
	if assert.NotNil(mv.T(), resp) {
		assert.Equal(mv.T(), mv.recorder.Code, http.StatusBadRequest)
		assert.Equal(mv.T(), resp.ErrorCode, common.InvalidPath)
	}
}

func (mv *Move) TestMoveSrcPathNotFound() {
	req := move.Request{
		SrcPath:  mv.srcDir + "/not-exist",
		DestPath: mv.destDir,
	}
	err := BindJsonRequest(mv.ctx, req)
	if err != nil {
		panic(fmt.Sprintf("bind request failed: %v", err))
	}
	mv.storage.Mv(mv.ctx)

	resp := &move.Response{}
	_ = ParseResp(mv.recorder.Result(), resp)
	if assert.NotNil(mv.T(), resp) {
		assert.Equal(mv.T(), mv.recorder.Code, http.StatusNotFound)
		assert.Equal(mv.T(), resp.ErrorCode, common.SrcPathNotFound)
	}
}

func (mv *Move) TestMoveDestPathExist() {
	req := move.Request{
		SrcPath:  mv.srcDir,
		DestPath: mv.srcDir,
	}
	err := BindJsonRequest(mv.ctx, req)
	if err != nil {
		panic(fmt.Sprintf("bind request failed: %v", err))
	}
	mv.storage.Mv(mv.ctx)

	resp := &move.Response{}
	_ = ParseResp(mv.recorder.Result(), resp)
	if assert.NotNil(mv.T(), resp) {
		assert.Equal(mv.T(), mv.recorder.Code, http.StatusBadRequest)
		assert.Equal(mv.T(), resp.ErrorCode, common.DestPathExists)
	}
}
