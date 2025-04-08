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
	quotaAPI "github.com/yuansuan/ticp/common/project-root-api/storage/quota/admin"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/dao/model"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/pathchecker"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/common/mock"
	"xorm.io/xorm"
)

type ListQuota struct {
	suite.Suite
	logger              *logging.Logger
	ctx                 *gin.Context
	recorder            *httptest.ResponseRecorder
	engine              *xorm.Engine
	mockStorageQuotaDao *dao.MockStorageQuotaDao
	ctrl                *gomock.Controller
}

func (list *ListQuota) SetupSuite() {
	list.T().Log("setup suite")
	logger, err := logging.NewLogger(logging.WithReleaseLevel(logging.DevelopmentLevel))
	if err != nil {
		panic(fmt.Sprintf("init logger failed: %v", err))
	}

	ctrl := gomock.NewController(list.T())
	list.ctrl = ctrl
	mockEngine, _ := mock.Engine()
	mockEngine.SetLogger(middleware.NewXormLogger(list.logger, true))
	list.engine = mockEngine
	list.logger = logger
}

func (list *ListQuota) TearDownSuite() {
	list.T().Log("teardown suite")
}

func (list *ListQuota) SetupTest() {
	list.T().Log("setup test")
	recorder := httptest.NewRecorder()
	ctx := GetTestGinContext(recorder)
	ctx.Set(logging.LoggerName, list.logger)
	list.ctx = ctx
	list.recorder = recorder
}

func (list *ListQuota) TearDownTest() {
	list.T().Log("teardown test")
}
func (list *ListQuota) mockQuota(mockDaoFunc func(*dao.MockStorageQuotaDao)) *Quota {
	pathchecker := pathchecker.PathAccessCheckerImpl{
		AuthEnabled:       true,
		IamClient:         nil,
		UserIDKey:         userIDKey,
		UserAppKeyInQuery: userAppKeyInQuery,
	}
	mockStorageQuotaDao := dao.NewMockStorageQuotaDao(list.ctrl)
	mockDaoFunc(mockStorageQuotaDao)
	return NewQuota(mockStorageQuotaDao, list.engine, pathchecker)
}

func TestListQuota(t *testing.T) {
	suite.Run(t, new(ListQuota))
}

func (list *ListQuota) TestListQuotaSuccess() {
	res := []*model.StorageQuota{
		{
			StorageUsage: 500,
			StorageLimit: 1000,
		},
		{
			StorageUsage: 200,
			StorageLimit: 500,
		},
		{
			StorageUsage: 50,
			StorageLimit: 100,
		},
	}

	quota := list.mockQuota(func(mockStorageQuotaDao *dao.MockStorageQuotaDao) {
		mockStorageQuotaDao.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any()).Return(res, nil, int64(5), int64(10)).AnyTimes()
	})

	BindRequest(list.ctx, quotaAPI.ListStorageQuotaRequest{PageOffset: 0, PageSize: 100})
	quota.ListStorageQuota(list.ctx)

	resp := &quotaAPI.ListStorageQuotaResponse{}
	err := ParseResp(list.recorder.Result(), resp)
	if !assert.Nil(list.T(), err) {
		list.T().Logf("%v: %v", time.Now().Format("2006-01-02 15:04:05"), err)
	} else {
		spew.Dump(resp)
	}
}

// ------------------------ error case ------------------------

func (list *ListQuota) TestListQuotaInvalidPageSize() {

	quota := list.mockQuota(func(mockStorageQuotaDao *dao.MockStorageQuotaDao) {
		mockStorageQuotaDao.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any()).Return([]*model.StorageQuota{}, nil, int64(5), int64(10)).AnyTimes()
	})

	BindRequest(list.ctx, quotaAPI.ListStorageQuotaRequest{PageOffset: 0, PageSize: -1})
	quota.ListStorageQuota(list.ctx)

	resp := &quotaAPI.ListStorageQuotaResponse{}
	_ = ParseResp(list.recorder.Result(), resp)
	if assert.NotNil(list.T(), resp) {
		assert.Equal(list.T(), list.recorder.Code, http.StatusBadRequest)
		assert.Equal(list.T(), resp.ErrorCode, common.InvalidPageSize)
	}
}

func (list *ListQuota) TestListQuotaInvalidPageOffset() {

	quota := list.mockQuota(func(mockStorageQuotaDao *dao.MockStorageQuotaDao) {
		mockStorageQuotaDao.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any()).Return([]*model.StorageQuota{}, nil, int64(5), int64(10)).AnyTimes()
	})

	BindRequest(list.ctx, quotaAPI.ListStorageQuotaRequest{PageOffset: -1, PageSize: 0})
	quota.ListStorageQuota(list.ctx)

	resp := &quotaAPI.ListStorageQuotaResponse{}
	_ = ParseResp(list.recorder.Result(), resp)
	if assert.NotNil(list.T(), resp) {
		assert.Equal(list.T(), list.recorder.Code, http.StatusBadRequest)
		assert.Equal(list.T(), resp.ErrorCode, common.InvalidPageOffset)
	}
}

func (list *ListQuota) TestListQuotaDBError() {

	quota := list.mockQuota(func(mockStorageQuotaDao *dao.MockStorageQuotaDao) {
		mockStorageQuotaDao.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any()).Return([]*model.StorageQuota{}, errors.New("db error"), int64(5), int64(10)).AnyTimes()
	})

	BindRequest(list.ctx, quotaAPI.ListStorageQuotaRequest{PageOffset: 0, PageSize: 100})
	quota.ListStorageQuota(list.ctx)

	resp := &quotaAPI.ListStorageQuotaResponse{}
	_ = ParseResp(list.recorder.Result(), resp)
	if assert.NotNil(list.T(), resp) {
		assert.Equal(list.T(), list.recorder.Code, http.StatusInternalServerError)
		assert.Equal(list.T(), resp.ErrorCode, common.InternalServerErrorCode)
	}
}
