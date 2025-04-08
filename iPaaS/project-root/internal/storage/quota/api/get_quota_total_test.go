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
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/pathchecker"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/common/mock"
	"xorm.io/xorm"
)

type GetQuotaTotal struct {
	suite.Suite
	logger              *logging.Logger
	ctx                 *gin.Context
	recorder            *httptest.ResponseRecorder
	engine              *xorm.Engine
	mockStorageQuotaDao *dao.MockStorageQuotaDao
	ctrl                *gomock.Controller
}

func (total *GetQuotaTotal) SetupSuite() {
	total.T().Log("setup suite")
	logger, err := logging.NewLogger(logging.WithReleaseLevel(logging.DevelopmentLevel))
	if err != nil {
		panic(fmt.Sprintf("init logger failed: %v", err))
	}

	ctrl := gomock.NewController(total.T())
	total.ctrl = ctrl
	mockEngine, _ := mock.Engine()
	mockEngine.SetLogger(middleware.NewXormLogger(total.logger, true))
	total.engine = mockEngine
	total.logger = logger
}

func (total *GetQuotaTotal) TearDownSuite() {
	total.T().Log("teardown suite")
}

func (total *GetQuotaTotal) SetupTest() {
	total.T().Log("setup test")
	recorder := httptest.NewRecorder()
	ctx := GetTestGinContext(recorder)
	ctx.Set(logging.LoggerName, total.logger)
	total.ctx = ctx
	total.recorder = recorder
}

func (total *GetQuotaTotal) TearDownTest() {
	total.T().Log("teardown test")
}
func (total *GetQuotaTotal) mockQuota(mockDaoFunc func(*dao.MockStorageQuotaDao)) *Quota {
	pathchecker := pathchecker.PathAccessCheckerImpl{
		AuthEnabled:       true,
		IamClient:         nil,
		UserIDKey:         userIDKey,
		UserAppKeyInQuery: userAppKeyInQuery,
	}
	mockStorageQuotaDao := dao.NewMockStorageQuotaDao(total.ctrl)
	mockDaoFunc(mockStorageQuotaDao)
	return NewQuota(mockStorageQuotaDao, total.engine, pathchecker)
}

func TestGetQuotaTotal(t *testing.T) {
	suite.Run(t, new(GetQuotaTotal))
}

func (total *GetQuotaTotal) TestGetQuotaTotalSuccess() {
	total.ctx.AddParam(UserIDKey, userID)

	quota := total.mockQuota(func(mockStorageQuotaDao *dao.MockStorageQuotaDao) {
		mockStorageQuotaDao.EXPECT().Total(gomock.Any()).Return(float64(1000), nil).AnyTimes()
	})

	quota.GetStorageQuotaTotal(total.ctx)
	resp := &quotaAPI.GetStorageQuotaResponse{}
	err := ParseResp(total.recorder.Result(), resp)
	if !assert.Nil(total.T(), err) {
		total.T().Logf("%v: %v", time.Now().Format("2006-01-02 15:04:05"), err)
	} else {
		spew.Dump(resp)
	}
}

// ------------------------ error case ------------------------

func (total *GetQuotaTotal) TestGetQuotaTotalDBError() {
	total.ctx.AddParam(UserIDKey, userID)

	quota := total.mockQuota(func(mockStorageQuotaDao *dao.MockStorageQuotaDao) {
		mockStorageQuotaDao.EXPECT().Total(gomock.Any()).Return(float64(0), errors.New("error db")).AnyTimes()
	})

	quota.GetStorageQuotaTotal(total.ctx)
	resp := &quotaAPI.GetStorageQuotaResponse{}
	_ = ParseResp(total.recorder.Result(), resp)
	if assert.NotNil(total.T(), resp) {
		assert.Equal(total.T(), total.recorder.Code, http.StatusInternalServerError)
		assert.Equal(total.T(), resp.ErrorCode, common.InternalServerErrorCode)
	}
}
