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
	"xorm.io/xorm"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
)

type ApplicationAllowSuite struct {
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
	mockAppAllowStore *store.MockApplicationAllowStore

	// field
	mockIDGen *snowflake.MockIDGen

	// service
	aqSrv *applicationAllowService

	// data
}

func TestApplicationAllowSuite(t *testing.T) {
	suite.Run(t, new(ApplicationAllowSuite))
}

func (s *ApplicationAllowSuite) SetupSuite() {
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
	applicationAllowStore := store.NewMockApplicationAllowStore(s.ctrl)
	s.mockAppAllowStore = applicationAllowStore
	s.mockStorage.EXPECT().Applications().Return(applicationStore).AnyTimes()
	s.mockStorage.EXPECT().ApplicationAllow().Return(applicationAllowStore).AnyTimes()
	s.mockStorage.EXPECT().Engine().Return(mockEngine).AnyTimes()

	// mock idgen
	mockIDGen := snowflake.NewMockIDGen(s.ctrl)
	s.mockIDGen = mockIDGen

	// service
	srv := &service{store: mockStorage}
	s.aqSrv = newAppAllowService(srv, mockIDGen)
}

func (s *ApplicationAllowSuite) TearDownSuite() {
	s.T().Log("teardown suite")
	s.engine.Close()
}

func (s *ApplicationAllowSuite) SetupTest() {
	s.T().Log("setup test")
}

func (s *ApplicationAllowSuite) TearDownTest() {
	s.T().Log("teardown test")
}

func (s *ApplicationAllowSuite) SetupSubTest() {
	s.T().Log("setup sub test")
}

func (s *ApplicationAllowSuite) TearDownSubTest() {
	s.T().Log("teardown sub test")
}

func (s *ApplicationAllowSuite) TestGetAllow() {
	s.Run("normal", func() {
		s.mockAppAllowStore.EXPECT().GetByAppId(gomock.Any(), nil,
			snowflake.ID(12345)).Return(&models.ApplicationAllow{
			ID:            10000, // 3Yq
			ApplicationID: snowflake.ID(12345),
			CreateTime:    time.Now(),
			UpdateTime:    time.Now(),
		}, nil)

		aq, err := s.aqSrv.GetAllow(s.ctx, snowflake.ID(12345))
		if !s.NoError(err) {
			return
		}

		s.Equal("3Yq", aq.ID)
	})

	s.Run("get allow error", func() {
		s.mockAppAllowStore.EXPECT().GetByAppId(gomock.Any(), nil,
			snowflake.ID(12345)).Return(nil, fmt.Errorf("get allow error"))
		aq, err := s.aqSrv.GetAllow(s.ctx, snowflake.ID(12345))
		if s.Error(err) {
			s.Equal("get allow error", err.Error())
			s.Nil(aq)
		}
	})

	s.Run("app allow not found", func() {
		s.mockAppAllowStore.EXPECT().GetByAppId(gomock.Any(), nil,
			snowflake.ID(12345)).Return(nil, xorm.ErrNotExist)
		aq, err := s.aqSrv.GetAllow(s.ctx, snowflake.ID(12345))
		if s.Error(err) {
			s.Equal(common.ErrAppAllowNotFound, err)
			s.Nil(aq)
		}
	})
}

func (s *ApplicationAllowSuite) TestAddAllow() {
	s.Run("normal", func() {
		s.mockIDGen.EXPECT().GenID(gomock.Any()).Return(snowflake.ID(10000), nil).Times(1)

		s.sqlmock.ExpectBegin()
		s.sqlmock.ExpectQuery("SELECT (.+) FROM `application` WHERE `id`=(.+) LIMIT 1 FOR UPDATE").
			WithArgs(12345).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "type", "version", "image", "publish_status", "bin_path"}).
				AddRow(12345, "test2", "test", "test", "test2", "unpublished", "test"))
		s.sqlmock.ExpectExec("INSERT INTO `application_allow`").
			WithArgs(10000, 12345, sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(0, 1))

		s.sqlmock.ExpectCommit()

		aq, err := s.aqSrv.AddAllow(s.ctx, snowflake.ID(12345))
		if !s.NoError(err) {
			return
		}

		s.Equal("3Yq", aq.ID)
	})
	s.Run("insert allow error", func() {
		s.mockIDGen.EXPECT().GenID(gomock.Any()).Return(snowflake.ID(10000), nil).Times(1)

		s.sqlmock.ExpectBegin()
		s.sqlmock.ExpectQuery("SELECT (.+) FROM `application` WHERE `id`=(.+) LIMIT 1 FOR UPDATE").
			WithArgs(12345).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "type", "version", "image", "publish_status", "bin_path"}).
				AddRow(12345, "test2", "test", "test", "test2", "unpublished", "test"))
		s.sqlmock.ExpectExec("INSERT INTO `application_allow`").
			WillReturnError(fmt.Errorf("insert allow error"))

		s.sqlmock.ExpectRollback()

		_, err := s.aqSrv.AddAllow(s.ctx, snowflake.ID(12345))
		if s.Error(err) {
			s.Equal("insert allow error", err.Error())
		}
	})

	s.Run("allow already exists", func() {
		s.mockIDGen.EXPECT().GenID(gomock.Any()).Return(snowflake.ID(10000), nil).Times(1)

		s.sqlmock.ExpectBegin()
		s.sqlmock.ExpectQuery("SELECT (.+) FROM `application` WHERE `id`=(.+) LIMIT 1 FOR UPDATE").
			WithArgs(12345).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "type", "version", "image", "publish_status", "bin_path"}).
				AddRow(12345, "test2", "test", "test", "test2", "unpublished", "test"))
		e := &mysql.MySQLError{
			Number: 1062,
		}
		s.sqlmock.ExpectExec("INSERT INTO `application_allow`").
			WillReturnError(e)

		s.sqlmock.ExpectRollback()

		_, err := s.aqSrv.AddAllow(s.ctx, snowflake.ID(12345))
		if s.Error(err) {
			s.Equal(common.ErrAppAllowAlreadyExist, err)
		}
	})
}

func (s *ApplicationAllowSuite) TestDeleteAllow() {
	s.Run("normal", func() {
		s.sqlmock.ExpectExec("DELETE FROM `application_allow`").
			WithArgs(12345).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := s.aqSrv.DeleteAllow(s.ctx, snowflake.ID(12345))
		if !s.NoError(err) {
			return
		}
	})

	s.Run("delete allow error", func() {
		s.sqlmock.ExpectExec("DELETE FROM `application_allow`").
			WillReturnError(fmt.Errorf("delete allow error"))

		err := s.aqSrv.DeleteAllow(s.ctx, snowflake.ID(12345))
		if s.Error(err) {
			s.Equal("delete allow error", err.Error())
		}
	})

	s.Run("allow not found", func() {
		s.sqlmock.ExpectExec("DELETE FROM `application_allow`").
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := s.aqSrv.DeleteAllow(s.ctx, snowflake.ID(12345))
		if s.Error(err) {
			s.Equal(common.ErrAppAllowNotFound, err)
		}
	})
}
