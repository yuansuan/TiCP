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
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/storage/pathchecker"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/common/mock"
	"xorm.io/xorm"
)

type PutQuota struct {
	suite.Suite
	logger              *logging.Logger
	ctx                 *gin.Context
	recorder            *httptest.ResponseRecorder
	engine              *xorm.Engine
	mockStorageQuotaDao *dao.MockStorageQuotaDao
	ctrl                *gomock.Controller
}

func (put *PutQuota) SetupSuite() {
	put.T().Log("setup suite")
	logger, err := logging.NewLogger(logging.WithReleaseLevel(logging.DevelopmentLevel))
	if err != nil {
		panic(fmt.Sprintf("init logger failed: %v", err))
	}

	ctrl := gomock.NewController(put.T())
	put.ctrl = ctrl
	mockEngine, _ := mock.Engine()
	mockEngine.SetLogger(middleware.NewXormLogger(put.logger, true))
	put.engine = mockEngine
	put.logger = logger
}

func (put *PutQuota) TearDownSuite() {
	put.T().Log("teardown suite")
}

func (put *PutQuota) SetupTest() {
	put.T().Log("setup test")
	recorder := httptest.NewRecorder()
	ctx := GetTestGinContext(recorder)
	ctx.Set(logging.LoggerName, put.logger)
	put.ctx = ctx
	put.recorder = recorder
}

func (put *PutQuota) TearDownTest() {
	put.T().Log("teardown test")
}

func (put *PutQuota) mockQuota(mockDaoFunc func(*dao.MockStorageQuotaDao)) *Quota {
	pathchecker := pathchecker.PathAccessCheckerImpl{
		AuthEnabled:       true,
		IamClient:         nil,
		UserIDKey:         userIDKey,
		UserAppKeyInQuery: userAppKeyInQuery,
	}
	mockStorageQuotaDao := dao.NewMockStorageQuotaDao(put.ctrl)
	mockDaoFunc(mockStorageQuotaDao)
	return NewQuota(mockStorageQuotaDao, put.engine, pathchecker)
}

func TestPutQuota(t *testing.T) {
	suite.Run(t, new(PutQuota))
}

func (put *PutQuota) TestPutQuotaInsertSuccess() {
	put.ctx.AddParam(UserIDKey, userID)

	put.ctx.Request.Header.Set(userIDKey, userID)
	put.ctx.Request.Header.Set(userAppKeyInQuery, ak)

	err := BindJsonRequest(put.ctx, quotaAPI.PutStorageQuotaRequest{StorageLimit: 500, UserID: userID})
	if err != nil {
		panic(fmt.Sprintf("bind request failed: %v", err))
	}

	quota := put.mockQuota(func(mockStorageQuotaDao *dao.MockStorageQuotaDao) {
		mockStorageQuotaDao.EXPECT().Get(gomock.Any(), gomock.Any()).Return(false, nil, nil).AnyTimes()
		mockStorageQuotaDao.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	})

	quota.PutStorageQuota(put.ctx)

	resp := &quotaAPI.PutStorageQuotaResponse{}
	err = ParseResp(put.recorder.Result(), resp)
	if !assert.Nil(put.T(), err) {
		put.T().Logf("%v: %v", time.Now().Format("2006-01-02 15:04:05"), err)
	} else {
		spew.Dump(resp)
	}
}

func (put *PutQuota) TestPutQuotaUpdateSuccess() {
	put.ctx.AddParam(UserIDKey, userID)

	put.ctx.Request.Header.Set(userIDKey, userID)
	put.ctx.Request.Header.Set(userAppKeyInQuery, ak)

	err := BindJsonRequest(put.ctx, quotaAPI.PutStorageQuotaRequest{StorageLimit: 500, UserID: userID})
	if err != nil {
		panic(fmt.Sprintf("bind request failed: %v", err))
	}

	quota := put.mockQuota(func(mockStorageQuotaDao *dao.MockStorageQuotaDao) {
		mockStorageQuotaDao.EXPECT().Get(gomock.Any(), gomock.Any()).Return(true, nil, nil).AnyTimes()
		mockStorageQuotaDao.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	})

	quota.PutStorageQuota(put.ctx)

	resp := &quotaAPI.PutStorageQuotaResponse{}
	err = ParseResp(put.recorder.Result(), resp)
	if !assert.Nil(put.T(), err) {
		put.T().Logf("%v: %v", time.Now().Format("2006-01-02 15:04:05"), err)
	} else {
		spew.Dump(resp)
	}
}

// ------------------------ error case ------------------------

func (put *PutQuota) TestPutQuotaNoUserID() {
	put.ctx.AddParam(UserIDKey, "")

	put.ctx.Request.Header.Set(userIDKey, userID)
	put.ctx.Request.Header.Set(userAppKeyInQuery, ak)

	err := BindJsonRequest(put.ctx, quotaAPI.PutStorageQuotaRequest{StorageLimit: 500, UserID: userID})
	if err != nil {
		panic(fmt.Sprintf("bind request failed: %v", err))
	}

	quota := put.mockQuota(func(mockStorageQuotaDao *dao.MockStorageQuotaDao) {
		mockStorageQuotaDao.EXPECT().Get(gomock.Any(), gomock.Any()).Return(true, nil, nil).AnyTimes()
	})

	quota.PutStorageQuota(put.ctx)

	resp := &quotaAPI.PutStorageQuotaResponse{}
	_ = ParseResp(put.recorder.Result(), resp)
	if assert.NotNil(put.T(), resp) {
		assert.Equal(put.T(), put.recorder.Code, http.StatusBadRequest)
		assert.Equal(put.T(), resp.ErrorCode, common.InvalidUserID)
	}
}

func (put *PutQuota) TestPutQuotaInvalidUserID() {
	put.ctx.AddParam(UserIDKey, InvalidUserID)

	put.ctx.Request.Header.Set(userIDKey, userID)
	put.ctx.Request.Header.Set(userAppKeyInQuery, ak)

	err := BindJsonRequest(put.ctx, quotaAPI.PutStorageQuotaRequest{StorageLimit: 500, UserID: InvalidUserID})
	if err != nil {
		panic(fmt.Sprintf("bind request failed: %v", err))
	}

	quota := put.mockQuota(func(mockStorageQuotaDao *dao.MockStorageQuotaDao) {
		mockStorageQuotaDao.EXPECT().Get(gomock.Any(), gomock.Any()).Return(true, nil, nil).AnyTimes()
	})

	quota.PutStorageQuota(put.ctx)

	resp := &quotaAPI.PutStorageQuotaResponse{}
	_ = ParseResp(put.recorder.Result(), resp)
	if assert.NotNil(put.T(), resp) {
		assert.Equal(put.T(), put.recorder.Code, http.StatusBadRequest)
		assert.Equal(put.T(), resp.ErrorCode, common.InvalidUserID)
	}
}

func (put *PutQuota) TestPutQuotaInvalidStorageLimit() {
	put.ctx.AddParam(UserIDKey, userID)

	put.ctx.Request.Header.Set(userIDKey, userID)
	put.ctx.Request.Header.Set(userAppKeyInQuery, ak)

	err := BindJsonRequest(put.ctx, quotaAPI.PutStorageQuotaRequest{StorageLimit: -500, UserID: userID})
	if err != nil {
		panic(fmt.Sprintf("bind request failed: %v", err))
	}

	quota := put.mockQuota(func(mockStorageQuotaDao *dao.MockStorageQuotaDao) {
	})
	quota.PutStorageQuota(put.ctx)

	resp := &quotaAPI.PutStorageQuotaResponse{}
	_ = ParseResp(put.recorder.Result(), resp)
	if assert.NotNil(put.T(), resp) {
		assert.Equal(put.T(), put.recorder.Code, http.StatusBadRequest)
		assert.Equal(put.T(), resp.ErrorCode, common.InvalidStorageLimit)
	}
}

func (put *PutQuota) TestPutQuotaGetDBError() {
	put.ctx.AddParam(UserIDKey, userID)

	put.ctx.Request.Header.Set(userIDKey, userID)
	put.ctx.Request.Header.Set(userAppKeyInQuery, ak)

	err := BindJsonRequest(put.ctx, quotaAPI.PutStorageQuotaRequest{StorageLimit: 500, UserID: userID})
	if err != nil {
		panic(fmt.Sprintf("bind request failed: %v", err))
	}

	quota := put.mockQuota(func(mockStorageQuotaDao *dao.MockStorageQuotaDao) {
		mockStorageQuotaDao.EXPECT().Get(gomock.Any(), gomock.Any()).Return(false, nil, errors.New("db error")).AnyTimes()
		mockStorageQuotaDao.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	})
	quota.PutStorageQuota(put.ctx)

	resp := &quotaAPI.PutStorageQuotaResponse{}
	_ = ParseResp(put.recorder.Result(), resp)
	if assert.NotNil(put.T(), resp) {
		assert.Equal(put.T(), put.recorder.Code, http.StatusInternalServerError)
		assert.Equal(put.T(), resp.ErrorCode, common.InternalServerErrorCode)
	}
}

func (put *PutQuota) TestPutQuotaInsertDBError() {
	put.ctx.AddParam(UserIDKey, userID)

	put.ctx.Request.Header.Set(userIDKey, userID)
	put.ctx.Request.Header.Set(userAppKeyInQuery, ak)

	err := BindJsonRequest(put.ctx, quotaAPI.PutStorageQuotaRequest{StorageLimit: 500, UserID: userID})
	if err != nil {
		panic(fmt.Sprintf("bind request failed: %v", err))
	}

	quota := put.mockQuota(func(mockStorageQuotaDao *dao.MockStorageQuotaDao) {
		mockStorageQuotaDao.EXPECT().Get(gomock.Any(), gomock.Any()).Return(true, nil, nil).AnyTimes()
		mockStorageQuotaDao.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("db error")).AnyTimes()
	})
	quota.PutStorageQuota(put.ctx)

	resp := &quotaAPI.PutStorageQuotaResponse{}
	_ = ParseResp(put.recorder.Result(), resp)
	if assert.NotNil(put.T(), resp) {
		assert.Equal(put.T(), put.recorder.Code, http.StatusInternalServerError)
		assert.Equal(put.T(), resp.ErrorCode, common.InternalServerErrorCode)
	}
}

func (put *PutQuota) TestPutQuotaUpdateDBError() {
	put.ctx.AddParam(UserIDKey, userID)

	put.ctx.Request.Header.Set(userIDKey, userID)
	put.ctx.Request.Header.Set(userAppKeyInQuery, ak)

	err := BindJsonRequest(put.ctx, quotaAPI.PutStorageQuotaRequest{StorageLimit: 500, UserID: userID})
	if err != nil {
		panic(fmt.Sprintf("bind request failed: %v", err))
	}

	quota := put.mockQuota(func(mockStorageQuotaDao *dao.MockStorageQuotaDao) {
		mockStorageQuotaDao.EXPECT().Get(gomock.Any(), gomock.Any()).Return(false, nil, nil).AnyTimes()
		mockStorageQuotaDao.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(errors.New("db error")).AnyTimes()
	})
	quota.PutStorageQuota(put.ctx)

	resp := &quotaAPI.PutStorageQuotaResponse{}
	_ = ParseResp(put.recorder.Result(), resp)
	if assert.NotNil(put.T(), resp) {
		assert.Equal(put.T(), put.recorder.Code, http.StatusInternalServerError)
		assert.Equal(put.T(), resp.ErrorCode, common.InternalServerErrorCode)
	}
}
