package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/middleware"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/common/project-root-api/common"
	quotaAPI "github.com/yuansuan/ticp/common/project-root-api/storage/quota/api"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao/model"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/pathchecker"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/common/mock"
	"xorm.io/xorm"
)

type GetQuotaAdmin struct {
	suite.Suite
	logger              *logging.Logger
	ctx                 *gin.Context
	recorder            *httptest.ResponseRecorder
	engine              *xorm.Engine
	mockStorageQuotaDao *dao.MockStorageQuotaDao
	ctrl                *gomock.Controller
}

func (admin *GetQuotaAdmin) SetupSuite() {
	admin.T().Log("setup suite")
	logger, err := logging.NewLogger(logging.WithReleaseLevel(logging.DevelopmentLevel))
	if err != nil {
		panic(fmt.Sprintf("init logger failed: %v", err))
	}

	ctrl := gomock.NewController(admin.T())
	admin.ctrl = ctrl
	mockEngine, _ := mock.Engine()
	mockEngine.SetLogger(middleware.NewXormLogger(admin.logger, true))
	admin.engine = mockEngine
	admin.logger = logger
}

func (admin *GetQuotaAdmin) TearDownSuite() {
	admin.T().Log("teardown suite")
}

func (admin *GetQuotaAdmin) SetupTest() {
	admin.T().Log("setup test")
	recorder := httptest.NewRecorder()
	ctx := GetTestGinContext(recorder)
	ctx.Set(logging.LoggerName, admin.logger)
	admin.ctx = ctx
	admin.recorder = recorder
}

func (admin *GetQuotaAdmin) TearDownTest() {
	admin.T().Log("teardown test")
}
func (admin *GetQuotaAdmin) mockQuota(mockDaoFunc func(*dao.MockStorageQuotaDao)) *Quota {
	pathchecker := pathchecker.PathAccessCheckerImpl{
		AuthEnabled:       true,
		IamClient:         nil,
		UserIDKey:         userIDKey,
		UserAppKeyInQuery: userAppKeyInQuery,
	}
	mockStorageQuotaDao := dao.NewMockStorageQuotaDao(admin.ctrl)
	mockDaoFunc(mockStorageQuotaDao)
	return NewQuota(mockStorageQuotaDao, admin.engine, pathchecker)
}

func TestGetQuotaAdmin(t *testing.T) {
	suite.Run(t, new(GetQuotaAdmin))
}

func (admin *GetQuotaAdmin) TestGetQuotaAdminSuccess() {
	admin.ctx.AddParam(UserIDKey, userID)

	quota := admin.mockQuota(func(mockStorageQuotaDao *dao.MockStorageQuotaDao) {
		mockStorageQuotaDao.EXPECT().Get(gomock.Any(), gomock.Any()).Return(true, &model.StorageQuota{StorageUsage: 500, StorageLimit: 1000}, nil).AnyTimes()
	})

	quota.GetStorageQuotaAdmin(admin.ctx)
	resp := &quotaAPI.GetStorageQuotaResponse{}
	err := ParseResp(admin.recorder.Result(), resp)
	if !assert.Nil(admin.T(), err) {
		admin.T().Logf("%v: %v", time.Now().Format("2006-01-02 15:04:05"), err)
	} else {
		spew.Dump(resp)
	}
}

// ------------------------ error case ------------------------

func (admin *GetQuotaAdmin) TestGetQuotaAdminNoUserID() {
	admin.ctx.AddParam(UserIDKey, "")

	quota := admin.mockQuota(func(mockStorageQuotaDao *dao.MockStorageQuotaDao) {
		mockStorageQuotaDao.EXPECT().Get(gomock.Any(), gomock.Any()).Return(true, &model.StorageQuota{StorageUsage: 500, StorageLimit: 1000}, nil).AnyTimes()
	})

	quota.GetStorageQuotaAdmin(admin.ctx)
	resp := &quotaAPI.GetStorageQuotaResponse{}
	_ = ParseResp(admin.recorder.Result(), resp)
	if assert.NotNil(admin.T(), resp) {
		assert.Equal(admin.T(), admin.recorder.Code, http.StatusBadRequest)
		assert.Equal(admin.T(), resp.ErrorCode, common.InvalidUserID)
	}
}

func (admin *GetQuotaAdmin) TestGetQuotaAdminInvalidUserID() {
	admin.ctx.AddParam(UserIDKey, InvalidUserID)
	quota := admin.mockQuota(func(mockStorageQuotaDao *dao.MockStorageQuotaDao) {
		mockStorageQuotaDao.EXPECT().Get(gomock.Any(), gomock.Any()).Return(true, &model.StorageQuota{StorageUsage: 500, StorageLimit: 1000}, nil).AnyTimes()
	})

	quota.GetStorageQuotaAdmin(admin.ctx)
	resp := &quotaAPI.GetStorageQuotaResponse{}
	_ = ParseResp(admin.recorder.Result(), resp)
	if assert.NotNil(admin.T(), resp) {
		assert.Equal(admin.T(), admin.recorder.Code, http.StatusBadRequest)
		assert.Equal(admin.T(), resp.ErrorCode, common.InvalidUserID)
	}
}

func (admin *GetQuotaAdmin) TestGetQuotaAdminDBError() {
	admin.ctx.AddParam(UserIDKey, userID)

	quota := admin.mockQuota(func(mockStorageQuotaDao *dao.MockStorageQuotaDao) {
		mockStorageQuotaDao.EXPECT().Get(gomock.Any(), gomock.Any()).Return(true, nil, errors.New("db error")).AnyTimes()
	})

	quota.GetStorageQuotaAdmin(admin.ctx)
	resp := &quotaAPI.GetStorageQuotaResponse{}
	_ = ParseResp(admin.recorder.Result(), resp)
	if assert.NotNil(admin.T(), resp) {
		assert.Equal(admin.T(), admin.recorder.Code, http.StatusInternalServerError)
		assert.Equal(admin.T(), resp.ErrorCode, common.InternalServerErrorCode)
	}
}

func (admin *GetQuotaAdmin) TestGetQuotaAdminNotFound() {
	admin.ctx.AddParam(UserIDKey, userID)

	quota := admin.mockQuota(func(mockStorageQuotaDao *dao.MockStorageQuotaDao) {
		mockStorageQuotaDao.EXPECT().Get(gomock.Any(), gomock.Any()).Return(false, nil, nil).AnyTimes()
	})

	quota.GetStorageQuotaAdmin(admin.ctx)
	resp := &quotaAPI.GetStorageQuotaResponse{}
	_ = ParseResp(admin.recorder.Result(), resp)
	if assert.NotNil(admin.T(), resp) {
		assert.Equal(admin.T(), admin.recorder.Code, http.StatusNotFound)
		assert.Equal(admin.T(), resp.ErrorCode, common.StorageQuotaNotFound)
	}
}
