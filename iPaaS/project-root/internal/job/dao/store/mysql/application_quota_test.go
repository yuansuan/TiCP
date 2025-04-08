package mysql

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/middleware"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/common/mock"
	"xorm.io/xorm"
)

// suite
type ApplicationQuotaSuite struct {
	// suite
	suite.Suite

	// base
	ctrl *gomock.Controller

	ctx    *gin.Context
	logger *logging.Logger

	engine  *xorm.Engine
	sqlmock sqlmock.Sqlmock

	// storage
	storage *ApplicationQuota
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
	logger, err := logging.NewLogger(logging.WithDefaultLogConfigOption())
	if !s.NoError(err) {
		return
	}
	s.ctx.Set(logging.LoggerName, s.logger)
	s.logger = logger

	// mock mysql engine
	mockEngine, mocksql := mock.Engine()
	mockEngine.SetLogger(middleware.NewXormLogger(s.logger, true))
	s.engine = mockEngine
	s.sqlmock = mocksql

	// storage
	s.storage = &ApplicationQuota{
		db: s.engine,
	}
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

func (s *ApplicationQuotaSuite) TestGetByUser() {
	s.Run("get success", func() {
		s.sqlmock.ExpectQuery("SELECT (.+) FROM `application_quota` WHERE `application_id`=(.+) AND `ys_id`=(.+) LIMIT 1").
			WithArgs(12345, 54321).
			WillReturnRows(sqlmock.NewRows([]string{"id", "application_id", "ys_id", "created_at", "updated_at"}).
				AddRow(10000, 12345, 54321, nil, nil))

		quota, err := s.storage.GetByUser(s.ctx, nil, snowflake.ID(12345), snowflake.ID(54321), false)
		if !s.NoError(err) {
			return
		}

		s.Equal(snowflake.ID(10000), quota.ID)
		s.Equal(snowflake.ID(12345), quota.ApplicationID)
		s.Equal(snowflake.ID(54321), quota.YsID)
	})

	s.Run("get success,forupdate", func() {
		s.sqlmock.ExpectQuery("SELECT (.+) FROM `application_quota` WHERE `application_id`=(.+) AND `ys_id`=(.+) LIMIT 1 FOR UPDATE").
			WithArgs(12345, 54321).
			WillReturnRows(sqlmock.NewRows([]string{"id", "application_id", "ys_id", "created_at", "updated_at"}).
				AddRow(10000, 12345, 54321, nil, nil))

		quota, err := s.storage.GetByUser(s.ctx, nil, snowflake.ID(12345), snowflake.ID(54321), true)
		if !s.NoError(err) {
			return
		}

		s.Equal(snowflake.ID(10000), quota.ID)
		s.Equal(snowflake.ID(12345), quota.ApplicationID)
		s.Equal(snowflake.ID(54321), quota.YsID)
	})

	s.Run("get error", func() {
		s.sqlmock.ExpectQuery("SELECT (.+) FROM `application_quota` WHERE `application_id`=(.+) AND `ys_id`=(.+) LIMIT 1").
			WithArgs(12345, 54321).
			WillReturnError(fmt.Errorf("db error"))

		quota, err := s.storage.GetByUser(s.ctx, nil, snowflake.ID(12345), snowflake.ID(54321), false)
		if s.Error(err) {
			s.EqualError(err, "db error")
			s.Nil(quota)
		}
	})

	s.Run("get not found", func() {
		s.sqlmock.ExpectQuery("SELECT (.+) FROM `application_quota` WHERE `application_id`=(.+) AND `ys_id`=(.+) LIMIT 1").
			WithArgs(12345, 54321).
			WillReturnRows(sqlmock.NewRows([]string{"id", "application_id", "ys_id", "created_at", "updated_at"}))

		quota, err := s.storage.GetByUser(s.ctx, nil, snowflake.ID(12345), snowflake.ID(54321), false)
		if s.Error(err) {
			s.Equal(err, xorm.ErrNotExist)
			s.Nil(quota)
		}
	})
}
