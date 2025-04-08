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

const (
	userIDKey         = "x-ys-user-id"
	userAppKeyInQuery = "AccessKeyId"
	ak                = "1UJJ7KXJ45M7CVVXCYDE"
	userID            = "4TiSsZonTa3"
	InvalidUserID     = "Invalid_user_id"
)

type GetQuotaAPI struct {
	suite.Suite
	logger              *logging.Logger
	ctx                 *gin.Context
	recorder            *httptest.ResponseRecorder
	engine              *xorm.Engine
	mockStorageQuotaDao *dao.MockStorageQuotaDao
	ctrl                *gomock.Controller
}

func (api *GetQuotaAPI) SetupSuite() {
	api.T().Log("setup suite")
	logger, err := logging.NewLogger(logging.WithReleaseLevel(logging.DevelopmentLevel))
	if err != nil {
		panic(fmt.Sprintf("init logger failed: %v", err))
	}

	ctrl := gomock.NewController(api.T())
	api.ctrl = ctrl
	mockEngine, _ := mock.Engine()
	mockEngine.SetLogger(middleware.NewXormLogger(api.logger, true))
	api.engine = mockEngine
	api.logger = logger
}

func (api *GetQuotaAPI) TearDownSuite() {
	api.T().Log("teardown suite")
}

func (api *GetQuotaAPI) SetupTest() {
	api.T().Log("setup test")
	recorder := httptest.NewRecorder()
	ctx := GetTestGinContext(recorder)
	ctx.Set(logging.LoggerName, api.logger)
	api.ctx = ctx
	api.recorder = recorder
}

func (api *GetQuotaAPI) TearDownTest() {
	api.T().Log("teardown test")
}
func (api *GetQuotaAPI) mockQuota(mockDaoFunc func(*dao.MockStorageQuotaDao)) *Quota {
	pathchecker := pathchecker.PathAccessCheckerImpl{
		AuthEnabled:       true,
		IamClient:         nil,
		UserIDKey:         userIDKey,
		UserAppKeyInQuery: userAppKeyInQuery,
	}
	mockStorageQuotaDao := dao.NewMockStorageQuotaDao(api.ctrl)
	mockDaoFunc(mockStorageQuotaDao)
	return NewQuota(mockStorageQuotaDao, api.engine, pathchecker)
}

func TestGetQuotaAPI(t *testing.T) {
	suite.Run(t, new(GetQuotaAPI))
}

func (api *GetQuotaAPI) TestGetQuotaAPISuccess() {
	api.ctx.Request.Header.Set(userIDKey, userID)
	api.ctx.Request.Header.Set(userAppKeyInQuery, ak)

	quota := api.mockQuota(func(mockStorageQuotaDao *dao.MockStorageQuotaDao) {
		mockStorageQuotaDao.EXPECT().Get(gomock.Any(), gomock.Any()).Return(true, &model.StorageQuota{StorageUsage: 500, StorageLimit: 1000}, nil).AnyTimes()
	})

	quota.GetStorageQuotaAPI(api.ctx)
	resp := &quotaAPI.GetStorageQuotaResponse{}
	err := ParseResp(api.recorder.Result(), resp)
	if !assert.Nil(api.T(), err) {
		api.T().Logf("%v: %v", time.Now().Format("2006-01-02 15:04:05"), err)
	} else {
		spew.Dump(resp)
	}
}

// ------------------------ error case ------------------------

func (api *GetQuotaAPI) TestGetQuotaAPINoUserID() {
	api.ctx.Request.Header.Set(userIDKey, "")
	api.ctx.Request.Header.Set(userAppKeyInQuery, ak)

	quota := api.mockQuota(func(mockStorageQuotaDao *dao.MockStorageQuotaDao) {
		mockStorageQuotaDao.EXPECT().Get(gomock.Any(), gomock.Any()).Return(true, &model.StorageQuota{StorageUsage: 500, StorageLimit: 1000}, nil).AnyTimes()
	})

	quota.GetStorageQuotaAPI(api.ctx)
	resp := &quotaAPI.GetStorageQuotaResponse{}
	_ = ParseResp(api.recorder.Result(), resp)
	if assert.NotNil(api.T(), resp) {
		assert.Equal(api.T(), api.recorder.Code, http.StatusBadRequest)
		assert.Equal(api.T(), resp.ErrorCode, common.InvalidUserID)
	}
}

func (api *GetQuotaAPI) TestGetQuotaAPIInvalidUserID() {
	api.ctx.Request.Header.Set(userIDKey, InvalidUserID)
	api.ctx.Request.Header.Set(userAppKeyInQuery, ak)

	quota := api.mockQuota(func(mockStorageQuotaDao *dao.MockStorageQuotaDao) {
		mockStorageQuotaDao.EXPECT().Get(gomock.Any(), gomock.Any()).Return(true, &model.StorageQuota{StorageUsage: 500, StorageLimit: 1000}, nil).AnyTimes()
	})

	quota.GetStorageQuotaAPI(api.ctx)
	resp := &quotaAPI.GetStorageQuotaResponse{}
	_ = ParseResp(api.recorder.Result(), resp)
	if assert.NotNil(api.T(), resp) {
		assert.Equal(api.T(), api.recorder.Code, http.StatusBadRequest)
		assert.Equal(api.T(), resp.ErrorCode, common.InvalidUserID)
	}
}

func (api *GetQuotaAPI) TestGetQuotaAPIUserNotFound() {
	api.ctx.Request.Header.Set(userIDKey, userID)
	api.ctx.Request.Header.Set(userAppKeyInQuery, ak)

	quota := api.mockQuota(func(mockStorageQuotaDao *dao.MockStorageQuotaDao) {
		mockStorageQuotaDao.EXPECT().Get(gomock.Any(), gomock.Any()).Return(false, nil, nil).AnyTimes()
	})

	quota.GetStorageQuotaAPI(api.ctx)
	resp := &quotaAPI.GetStorageQuotaResponse{}
	_ = ParseResp(api.recorder.Result(), resp)
	if assert.NotNil(api.T(), resp) {
		assert.Equal(api.T(), api.recorder.Code, http.StatusNotFound)
		assert.Equal(api.T(), resp.ErrorCode, common.StorageQuotaNotFound)
	}
}

func (api *GetQuotaAPI) TestGetQuotaAPIDBError() {
	api.ctx.Request.Header.Set(userIDKey, userID)
	api.ctx.Request.Header.Set(userAppKeyInQuery, ak)

	quota := api.mockQuota(func(mockStorageQuotaDao *dao.MockStorageQuotaDao) {
		mockStorageQuotaDao.EXPECT().Get(gomock.Any(), gomock.Any()).Return(false, nil, errors.New("db error")).AnyTimes()
	})

	quota.GetStorageQuotaAPI(api.ctx)
	resp := &quotaAPI.GetStorageQuotaResponse{}
	_ = ParseResp(api.recorder.Result(), resp)
	if assert.NotNil(api.T(), resp) {
		assert.Equal(api.T(), api.recorder.Code, http.StatusInternalServerError)
		assert.Equal(api.T(), resp.ErrorCode, common.InternalServerErrorCode)
	}
}
