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
	cp "github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/copy"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao/model"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/filemode"
	v20230530 "github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/local/v20230530"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/common/mock"
	"xorm.io/xorm"
)

type Copy struct {
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

func (copy *Copy) SetupSuite() {
	copy.T().Log("setup suite")
	logger, err := logging.NewLogger(logging.WithReleaseLevel(logging.DevelopmentLevel))
	if err != nil {
		panic(fmt.Sprintf("init logger failed: %v", err))
	}

	ctrl := gomock.NewController(copy.T())
	mockStorageQuotaDao := dao.NewMockStorageQuotaDao(ctrl)
	copy.mockStorageQuotaDao = mockStorageQuotaDao
	copy.mockStorageQuotaDao.EXPECT().Get(gomock.Any(), gomock.Any()).Return(true, &model.StorageQuota{StorageLimit: 1000}, nil).AnyTimes()
	mockStorageOperationLogDao := dao.NewMockStorageOperationLogDao(ctrl)
	copy.mockStorageOperationLogDao = mockStorageOperationLogDao
	copy.mockStorageOperationLogDao.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mockEngine, _ := mock.Engine()
	mockEngine.SetLogger(middleware.NewXormLogger(copy.logger, true))
	copy.engine = mockEngine

	copy.logger = logger
	currentDir, storage := CreateLocalStorage(copy.userID, v20230530.SoftLink, copy.mockStorageQuotaDao, copy.mockStorageOperationLogDao, copy.engine)
	copy.storage = storage

	copy.userID = "4TiSBX39DtN"
	copy.testDir = "/test-cp"

	tmpDir := currentDir + "/" + copy.userID + copy.testDir
	copy.tmpDir = tmpDir
	copy.srcDir = "/" + copy.userID + copy.testDir + "/src"
	copy.destDir = "/" + copy.userID + copy.testDir + "/dest"

	err = os.MkdirAll(filepath.Join(copy.tmpDir, "src"), filemode.Directory)
	if err != nil {
		panic(fmt.Sprintf("create src dir failed: %v", err))
	}

	for i := 0; i < 2; i++ {
		filePath := filepath.Join(copy.tmpDir, "src", "file-"+fmt.Sprintf("%d", i+1))
		err := os.WriteFile(filePath, nil, 0644)
		if err != nil {
			panic(fmt.Sprintf("Failed to create file: %v", err))
		}
		copy.T().Log("File created:", filePath)
	}

}

func (copy *Copy) TearDownSuite() {
	copy.T().Log("teardown suite")
	_ = os.RemoveAll(filepath.Dir(copy.tmpDir))
}

func (copy *Copy) SetupTest() {
	copy.T().Log("setup test")
	recorder := httptest.NewRecorder()
	ctx := GetTestGinContext(recorder)
	ctx.Set(logging.LoggerName, copy.logger)
	copy.ctx = ctx
	copy.recorder = recorder
}

func (copy *Copy) TearDownTest() {
	copy.T().Log("teardown test")
}

func TestCopy(t *testing.T) {
	suite.Run(t, new(Copy))
}

func (copy *Copy) TestCopySuccessDir() {
	req := cp.Request{
		SrcPath:  copy.srcDir,
		DestPath: copy.destDir,
	}
	err := BindJsonRequest(copy.ctx, req)
	if err != nil {
		panic(fmt.Sprintf("bind request failed: %v", err))
	}
	copy.storage.Copy(copy.ctx)

	resp := &cp.Response{}
	err = ParseResp(copy.recorder.Result(), resp)
	if !assert.Nil(copy.T(), err) {
		copy.T().Logf("%v: %v", time.Now().Format("2006-01-02 15:04:05"), err)
	} else {
		spew.Dump(resp)
	}
}

func (copy *Copy) TestCopySuccessFile() {
	req := cp.Request{
		SrcPath:  copy.srcDir + "/file-1",
		DestPath: copy.destDir + "/new-file",
	}
	err := BindJsonRequest(copy.ctx, req)
	if err != nil {
		panic(fmt.Sprintf("bind request failed: %v", err))
	}
	copy.storage.Copy(copy.ctx)

	resp := &cp.Response{}
	err = ParseResp(copy.recorder.Result(), resp)
	if !assert.Nil(copy.T(), err) {
		copy.T().Logf("%v: %v", time.Now().Format("2006-01-02 15:04:05"), err)
	} else {
		spew.Dump(resp)
	}
}

// ------------------------ error case ------------------------

func (copy *Copy) TestCopyInvalidSrcPath() {
	req := cp.Request{
		SrcPath:  copy.srcDir + "../file",
		DestPath: copy.destDir + "/file",
	}
	err := BindJsonRequest(copy.ctx, req)
	if err != nil {
		panic(fmt.Sprintf("bind request failed: %v", err))
	}
	copy.storage.Copy(copy.ctx)

	resp := &cp.Response{}
	_ = ParseResp(copy.recorder.Result(), resp)
	if assert.NotNil(copy.T(), resp) {
		assert.Equal(copy.T(), copy.recorder.Code, http.StatusBadRequest)
		assert.Equal(copy.T(), resp.ErrorCode, common.InvalidPath)
	}
}

func (copy *Copy) TestCopyInvalidDestPath() {
	req := cp.Request{
		SrcPath:  copy.srcDir + "/file-1",
		DestPath: copy.destDir + "..file",
	}
	err := BindJsonRequest(copy.ctx, req)
	if err != nil {
		panic(fmt.Sprintf("bind request failed: %v", err))
	}
	copy.storage.Copy(copy.ctx)

	resp := &cp.Response{}
	_ = ParseResp(copy.recorder.Result(), resp)
	if assert.NotNil(copy.T(), resp) {
		assert.Equal(copy.T(), copy.recorder.Code, http.StatusBadRequest)
		assert.Equal(copy.T(), resp.ErrorCode, common.InvalidPath)
	}
}

func (copy *Copy) TestCopyInvalidSrcNotFound() {
	req := cp.Request{
		SrcPath:  copy.srcDir + "/fake-file",
		DestPath: copy.destDir + "/file-2",
	}
	err := BindJsonRequest(copy.ctx, req)
	if err != nil {
		panic(fmt.Sprintf("bind request failed: %v", err))
	}
	copy.storage.Copy(copy.ctx)

	resp := &cp.Response{}
	_ = ParseResp(copy.recorder.Result(), resp)
	if assert.NotNil(copy.T(), resp) {
		assert.Equal(copy.T(), copy.recorder.Code, http.StatusNotFound)
		assert.Equal(copy.T(), resp.ErrorCode, common.SrcPathNotFound)
	}
}

func (copy *Copy) TestCopyInvalidDestExist() {
	req := cp.Request{
		SrcPath:  copy.srcDir + "/file-1",
		DestPath: copy.srcDir + "/file-2",
	}
	err := BindJsonRequest(copy.ctx, req)
	if err != nil {
		panic(fmt.Sprintf("bind request failed: %v", err))
	}
	copy.storage.Copy(copy.ctx)

	resp := &cp.Response{}
	_ = ParseResp(copy.recorder.Result(), resp)
	if assert.NotNil(copy.T(), resp) {
		assert.Equal(copy.T(), copy.recorder.Code, http.StatusBadRequest)
		assert.Equal(copy.T(), resp.ErrorCode, common.DestPathExists)
	}
}
