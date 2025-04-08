package application

import (
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/middleware"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao/models"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao/store"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/common/mock"
	hydra_lcp "github.com/yuansuan/ticp/iPaaS/sso/protos"
	"xorm.io/xorm"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
)

type ApplicationQuotaSuite struct {
	// suite
	suite.Suite

	// base
	ctrl   *gomock.Controller
	ctx    *gin.Context
	logger *logging.Logger

	engine  *xorm.Engine
	sqlmock sqlmock.Sqlmock

	// storage
	mockStorage       *store.MockFactoryNew
	mockAppStore      *store.MockApplicationStore
	mockAppQuotaStore *store.MockApplicationQuotaStore

	// field
	mockIDGen     *snowflake.MockIDGen
	mockUserGeter *MockUserGeter

	// service
	aqSrv *applicationQuotaService

	// data
}

func TestApplicationQuotaSuite(t *testing.T) {
	suite.Run(t, new(ApplicationQuotaSuite))
}

func (s *ApplicationQuotaSuite) SetupSuite() {
	s.T().Log("setup suite")

	// mock base controller
	s.ctrl = gomock.NewController(s.T())

	// mock gin context
	w := httptest.NewRecorder()
	s.ctx = mock.GinContext(w)

	// logger
	logger, err := logging.NewLogger(logging.WithReleaseLevel(logging.DevelopmentLevel))
	if err != nil {
		panic(fmt.Sprintf("init logger failed: %v", err))
	}
	s.ctx.Set(logging.LoggerName, logger)
	s.logger = logger

	// mock mysql engine
	mockEngine, mocksql := mock.Engine()
	mockEngine.SetLogger(middleware.NewXormLogger(s.logger, true))
	s.engine = mockEngine
	s.sqlmock = mocksql

	// mock store
	mockStorage := store.NewMockFactoryNew(s.ctrl)
	s.mockStorage = mockStorage
	applicationStore := store.NewMockApplicationStore(s.ctrl)
	s.mockAppStore = applicationStore
	applicationQuotaStore := store.NewMockApplicationQuotaStore(s.ctrl)
	s.mockAppQuotaStore = applicationQuotaStore
	s.mockStorage.EXPECT().Applications().Return(applicationStore).AnyTimes()
	s.mockStorage.EXPECT().ApplicationQuota().Return(applicationQuotaStore).AnyTimes()
	s.mockStorage.EXPECT().Engine().Return(mockEngine).AnyTimes()

	// mock idgen
	mockIDGen := snowflake.NewMockIDGen(s.ctrl)
	s.mockIDGen = mockIDGen

	// mock user geter
	mockUserGeter := NewMockUserGeter(s.ctrl)
	s.mockUserGeter = mockUserGeter

	// service
	srv := &service{store: mockStorage}
	mockHelper := struct {
		*snowflake.MockIDGen
		*MockUserGeter
	}{
		mockIDGen,
		mockUserGeter,
	}
	s.aqSrv = newAppQuotaService(srv, mockHelper)
}

func (s *ApplicationQuotaSuite) TearDownSuite() {
	s.T().Log("teardown suite")
	s.engine.Close()
}

func (s *ApplicationQuotaSuite) SetupTest() {
	s.T().Log("setup test")
}

func (s *ApplicationQuotaSuite) TearDownTest() {
	s.T().Log("teardown test")
}

func (s *ApplicationQuotaSuite) SetupSubTest() {
	s.T().Log("setup sub test")
}

func (s *ApplicationQuotaSuite) TearDownSubTest() {
	s.T().Log("teardown sub test")
}

func (s *ApplicationQuotaSuite) TestGetQuota() {
	s.Run("normal", func() {
		s.mockAppQuotaStore.EXPECT().GetByUser(gomock.Any(), nil, snowflake.ID(12345), snowflake.ID(54321), false).Return(&models.ApplicationQuota{
			ID:            10000, // 3Yq
			ApplicationID: snowflake.ID(12345),
			YsID:          snowflake.ID(54321),
			CreateTime:    time.Now(),
			UpdateTime:    time.Now(),
		}, nil)

		aq, err := s.aqSrv.GetQuota(s.ctx, snowflake.ID(12345), snowflake.ID(54321))
		if !s.NoError(err) {
			return
		}

		s.Equal("3Yq", aq.ID)
	})

	s.Run("get quota error", func() {
		s.mockAppQuotaStore.EXPECT().GetByUser(gomock.Any(), nil, snowflake.ID(12345), snowflake.ID(54321), false).Return(nil, fmt.Errorf("get quota error"))

		aq, err := s.aqSrv.GetQuota(s.ctx, snowflake.ID(12345), snowflake.ID(54321))

		if s.Error(err) {
			s.Equal("get quota error", err.Error())
			s.Nil(aq)
		}
	})

	s.Run("app quota not found", func() {
		s.mockAppQuotaStore.EXPECT().GetByUser(gomock.Any(), nil, snowflake.ID(12345), snowflake.ID(54321), false).Return(nil, xorm.ErrNotExist)

		aq, err := s.aqSrv.GetQuota(s.ctx, snowflake.ID(12345), snowflake.ID(54321))

		if s.Error(err) {
			s.Equal(common.ErrAppQuotaNotFound, err)
			s.Nil(aq)
		}
	})
}

func (s *ApplicationQuotaSuite) TestAddQuota() {
	s.Run("normal", func() {
		s.mockIDGen.EXPECT().GenID(gomock.Any()).Return(snowflake.ID(10000), nil).Times(1)
		s.mockUserGeter.EXPECT().GetUser(gomock.Any(), "h9z").Return(&hydra_lcp.UserInfo{}, nil).Times(1)

		s.sqlmock.ExpectBegin()
		s.sqlmock.ExpectQuery("SELECT (.+) FROM `application` WHERE `id`=(.+) LIMIT 1 FOR UPDATE").
			WithArgs(12345).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "type", "version", "image", "publish_status", "bin_path"}).
				AddRow(12345, "test2", "test", "test", "test2", "unpublished", "test"))

		s.sqlmock.ExpectExec("INSERT INTO `application_quota`").
			WithArgs(10000, 12345, 54321, sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(0, 1))

		s.sqlmock.ExpectCommit()

		aq, err := s.aqSrv.AddQuota(s.ctx, snowflake.ID(12345), snowflake.ID(54321))
		if !s.NoError(err) {
			return
		}

		s.Equal("3Yq", aq.ID)
	})

	s.Run("get app error", func() {
		s.sqlmock.ExpectBegin()
		s.sqlmock.ExpectQuery("SELECT (.+) FROM `application` WHERE `id`=(.+) LIMIT 1 FOR UPDATE").
			WithArgs(12345).
			WillReturnError(fmt.Errorf("get app error"))

		s.sqlmock.ExpectRollback()

		_, err := s.aqSrv.AddQuota(s.ctx, snowflake.ID(12345), snowflake.ID(54321))
		if s.Error(err) {
			s.Equal("get app error", err.Error())
		}
	})

	s.Run("app not found", func() {
		s.sqlmock.ExpectBegin()
		s.sqlmock.ExpectQuery("SELECT (.+) FROM `application` WHERE `id`=(.+) LIMIT 1 FOR UPDATE").
			WithArgs(12345).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "type", "version", "image", "publish_status", "bin_path"}))

		s.sqlmock.ExpectRollback()

		_, err := s.aqSrv.AddQuota(s.ctx, snowflake.ID(12345), snowflake.ID(54321))
		if s.Error(err) {
			s.Equal(common.ErrAppIDNotFound, err)
		}
	})

	s.Run("get user error", func() {
		s.mockUserGeter.EXPECT().GetUser(gomock.Any(), "h9z").Return(nil, nil).Times(1)

		s.sqlmock.ExpectBegin()
		s.sqlmock.ExpectQuery("SELECT (.+) FROM `application` WHERE `id`=(.+) LIMIT 1 FOR UPDATE").
			WithArgs(12345).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "type", "version", "image", "publish_status", "bin_path"}).
				AddRow(12345, "test2", "test", "test", "test2", "unpublished", "test"))

		s.sqlmock.ExpectRollback()

		_, err := s.aqSrv.AddQuota(s.ctx, snowflake.ID(12345), snowflake.ID(54321))
		if s.Error(err) {
			s.Equal(common.ErrUserNotExists, err)
		}
	})

	s.Run("get user not found", func() {
		s.mockUserGeter.EXPECT().GetUser(gomock.Any(), "h9z").Return(nil, fmt.Errorf("get user error")).Times(1)

		s.sqlmock.ExpectBegin()
		s.sqlmock.ExpectQuery("SELECT (.+) FROM `application` WHERE `id`=(.+) LIMIT 1 FOR UPDATE").
			WithArgs(12345).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "type", "version", "image", "publish_status", "bin_path"}).
				AddRow(12345, "test2", "test", "test", "test2", "unpublished", "test"))

		s.sqlmock.ExpectRollback()

		_, err := s.aqSrv.AddQuota(s.ctx, snowflake.ID(12345), snowflake.ID(54321))
		if s.Error(err) {
			s.Equal("get user error", err.Error())
		}
	})

	s.Run("gen id error", func() {
		s.mockIDGen.EXPECT().GenID(gomock.Any()).Return(snowflake.Zero(), fmt.Errorf("gen id error")).Times(1)
		s.mockUserGeter.EXPECT().GetUser(gomock.Any(), "h9z").Return(&hydra_lcp.UserInfo{}, nil).Times(1)

		s.sqlmock.ExpectBegin()
		s.sqlmock.ExpectQuery("SELECT (.+) FROM `application` WHERE `id`=(.+) LIMIT 1 FOR UPDATE").
			WithArgs(12345).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "type", "version", "image", "publish_status", "bin_path"}).
				AddRow(12345, "test2", "test", "test", "test2", "unpublished", "test"))

		s.sqlmock.ExpectRollback()

		_, err := s.aqSrv.AddQuota(s.ctx, snowflake.ID(12345), snowflake.ID(54321))
		if s.Error(err) {
			s.Equal("gen id error", err.Error())
		}
	})

	s.Run("insert quota error", func() {
		s.mockIDGen.EXPECT().GenID(gomock.Any()).Return(snowflake.ID(10000), nil).Times(1)
		s.mockUserGeter.EXPECT().GetUser(gomock.Any(), "h9z").Return(&hydra_lcp.UserInfo{}, nil).Times(1)

		s.sqlmock.ExpectBegin()
		s.sqlmock.ExpectQuery("SELECT (.+) FROM `application` WHERE `id`=(.+) LIMIT 1 FOR UPDATE").
			WithArgs(12345).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "type", "version", "image", "publish_status", "bin_path"}).
				AddRow(12345, "test2", "test", "test", "test2", "unpublished", "test"))

		s.sqlmock.ExpectExec("INSERT INTO `application_quota`").
			WillReturnError(fmt.Errorf("insert quota error"))

		s.sqlmock.ExpectRollback()

		_, err := s.aqSrv.AddQuota(s.ctx, snowflake.ID(12345), snowflake.ID(54321))
		if s.Error(err) {
			s.Equal("insert quota error", err.Error())
		}
	})

	s.Run("quota already exists", func() {
		s.mockIDGen.EXPECT().GenID(gomock.Any()).Return(snowflake.ID(10000), nil).Times(1)
		s.mockUserGeter.EXPECT().GetUser(gomock.Any(), "h9z").Return(&hydra_lcp.UserInfo{}, nil).Times(1)

		s.sqlmock.ExpectBegin()
		s.sqlmock.ExpectQuery("SELECT (.+) FROM `application` WHERE `id`=(.+) LIMIT 1 FOR UPDATE").
			WithArgs(12345).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "type", "version", "image", "publish_status", "bin_path"}).
				AddRow(12345, "test2", "test", "test", "test2", "unpublished", "test"))
		e := &mysql.MySQLError{
			Number: 1062,
		}
		s.sqlmock.ExpectExec("INSERT INTO `application_quota`").
			WillReturnError(e)

		s.sqlmock.ExpectRollback()

		_, err := s.aqSrv.AddQuota(s.ctx, snowflake.ID(12345), snowflake.ID(54321))
		if s.Error(err) {
			s.Equal(common.ErrAppQuotaAlreadyExist, err)
		}
	})
}

func (s *ApplicationQuotaSuite) TestDeleteQuota() {
	s.Run("normal", func() {
		s.sqlmock.ExpectExec("DELETE FROM `application_quota`").
			WithArgs(12345, 54321).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := s.aqSrv.DeleteQuota(s.ctx, snowflake.ID(12345), snowflake.ID(54321))
		if !s.NoError(err) {
			return
		}
	})

	s.Run("delete quota error", func() {
		s.sqlmock.ExpectExec("DELETE FROM `application_quota`").
			WillReturnError(fmt.Errorf("delete quota error"))

		err := s.aqSrv.DeleteQuota(s.ctx, snowflake.ID(12345), snowflake.ID(54321))
		if s.Error(err) {
			s.Equal("delete quota error", err.Error())
		}
	})

	s.Run("quota not found", func() {
		s.sqlmock.ExpectExec("DELETE FROM `application_quota`").
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := s.aqSrv.DeleteQuota(s.ctx, snowflake.ID(12345), snowflake.ID(54321))
		if s.Error(err) {
			s.Equal(common.ErrAppQuotaNotFound, err)
		}
	})
}
