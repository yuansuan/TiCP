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
type ApplicationAllowSuite struct {
	// suite
	suite.Suite

	// base
	ctrl *gomock.Controller

	ctx    *gin.Context
	logger *logging.Logger

	engine  *xorm.Engine
	sqlmock sqlmock.Sqlmock

	// storage
	storage *ApplicationAllow
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
	s.storage = &ApplicationAllow{
		db: s.engine,
	}
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

func (s *ApplicationAllowSuite) TestGetByApplication() {
	s.Run("get success", func() {
		s.sqlmock.ExpectQuery("SELECT (.+) FROM `application_allow` WHERE `application_id`=(.+) LIMIT 1").
			WithArgs(12345).
			WillReturnRows(sqlmock.NewRows([]string{"id", "application_id", "created_at", "updated_at"}).
				AddRow(10000, 12345, nil, nil))

		item, err := s.storage.GetByAppId(s.ctx, nil, snowflake.ID(12345))
		if !s.NoError(err) {
			return
		}

		s.Equal(snowflake.ID(10000), item.ID)
		s.Equal(snowflake.ID(12345), item.ApplicationID)
	})

	s.Run("get error", func() {
		s.sqlmock.ExpectQuery("SELECT (.+) FROM `application_allow` WHERE `application_id`=(.+)  LIMIT 1").
			WithArgs(99999).
			WillReturnError(fmt.Errorf("db error"))

		item, err := s.storage.GetByAppId(s.ctx, nil, snowflake.ID(99999))
		if s.Error(err) {
			s.EqualError(err, "db error")
			s.Nil(item)
		}
	})

	s.Run("get not found", func() {
		s.sqlmock.ExpectQuery("SELECT (.+) FROM `application_allow` WHERE `application_id`=(.+)  LIMIT 1").
			WithArgs(99999).
			WillReturnRows(sqlmock.NewRows([]string{"id", "application_id", "created_at", "updated_at"}))

		item, err := s.storage.GetByAppId(s.ctx, nil, snowflake.ID(99999))
		if s.Error(err) {
			s.Equal(err, xorm.ErrNotExist)
			s.Nil(item)
		}
	})
}
