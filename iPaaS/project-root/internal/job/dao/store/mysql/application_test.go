package mysql

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/middleware"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/consts"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/common/mock"
	"xorm.io/xorm"
)

// suite
type ApplicationSuite struct {
	// suite
	suite.Suite

	// base
	ctrl *gomock.Controller

	ctx    *gin.Context
	logger *logging.Logger

	engine  *xorm.Engine
	sqlmock sqlmock.Sqlmock

	// storage
	storage *Application
}

func TestApplicationSuite(t *testing.T) {
	suite.Run(t, new(ApplicationSuite))
}

func (s *ApplicationSuite) SetupSuite() {
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
	s.storage = &Application{
		db: s.engine,
	}
}

func (s *ApplicationSuite) TearDownSuite() {
	s.T().Log("teardown suite")
	s.engine.Close()
}

func (s *ApplicationSuite) SetupTest() {
	s.T().Log("setup test")
}

func (s *ApplicationSuite) TearDownTest() {
	s.T().Log("teardown test")
}

func (s *ApplicationSuite) SetupSubTest() {
	s.T().Log("setup sub test")
}

func (s *ApplicationSuite) TearDownSubTest() {
	s.T().Log("teardown sub test")
}

func (s *ApplicationSuite) TestListApps() {
	s.Run("list success,normal", func() {
		// mock sql
		s.sqlmock.ExpectQuery("SELECT (.+) FROM `application`").
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "type", "version", "image", "publish_status", "bin_path"}).
				AddRow(1, "test", "test", "test", "test", "published", "test").
				AddRow(2, "test2", "test", "test", "test2", "unpublished", "test"))

		s.sqlmock.ExpectQuery("SELECT count(.+) FROM `application`").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).
				AddRow(2))

		// list apps
		apps, count, err := s.storage.ListApps(s.ctx, 0, consts.PublishStatusAll)
		if !s.NoError(err) {
			return
		}

		s.Len(apps, 2)
		s.Equal(int64(2), count)
		s.Equal("test", apps[0].Name)
		s.Equal("test2", apps[1].Name)

		s.T().Logf("apps: %s", spew.Sdump(apps))
	})

	s.Run("list success,publish_status", func() {
		// mock sql
		s.sqlmock.ExpectQuery("SELECT (.+) FROM `application` WHERE \\(application.publish_status = (.+)\\)").
			WithArgs("published").
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "type", "version", "image", "publish_status", "bin_path"}).
				AddRow(1, "test", "test", "test", "test", "published", "test"))

		s.sqlmock.ExpectQuery("SELECT count(.+) FROM `application` WHERE \\(application.publish_status = (.+)\\)").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).
				AddRow(1))

		// list apps
		apps, count, err := s.storage.ListApps(s.ctx, 0, consts.PublishStatusPublished)
		if !s.NoError(err) {
			return
		}

		s.Len(apps, 1)
		s.Equal(int64(1), count)
		s.Equal("test", apps[0].Name)

		s.T().Logf("apps: %s", spew.Sdump(apps))
	})

	s.Run("list success,userID", func() {
		// mock sql
		s.sqlmock.ExpectQuery("SELECT (.+) FROM `application` INNER JOIN `application_quota` ON application_quota.application_id = application.id WHERE \\(application_quota.ys_id = (.+)\\)").
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "type", "version", "image", "publish_status", "bin_path"}).
				AddRow(1, "test", "test", "test", "test", "published", "test"))

		s.sqlmock.ExpectQuery("SELECT count(.+) FROM `application` INNER JOIN `application_quota` ON application_quota.application_id = application.id WHERE \\(application_quota.ys_id = (.+)\\)").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).
				AddRow(1))

		// list apps
		apps, count, err := s.storage.ListApps(s.ctx, 1, consts.PublishStatusAll)
		if !s.NoError(err) {
			return
		}

		s.Len(apps, 1)
		s.Equal(int64(1), count)
		s.Equal("test", apps[0].Name)

		s.T().Logf("apps: %s", spew.Sdump(apps))
	})

	s.Run("list success,userID,publish_status", func() {
		// mock sql
		s.sqlmock.ExpectQuery("SELECT (.+) FROM `application` INNER JOIN `application_quota` ON application_quota.application_id = application.id WHERE \\(application.publish_status = (.+)\\) AND \\(application_quota.ys_id = (.+)\\)").
			WithArgs("published", 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "type", "version", "image", "publish_status", "bin_path"}).
				AddRow(1, "test", "test", "test", "test", "published", "test"))

		s.sqlmock.ExpectQuery("SELECT count(.+) FROM `application` INNER JOIN `application_quota` ON application_quota.application_id = application.id WHERE \\(application.publish_status = (.+)\\) AND \\(application_quota.ys_id = (.+)\\)").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).
				AddRow(1))

		// list apps
		apps, count, err := s.storage.ListApps(s.ctx, 1, consts.PublishStatusPublished)
		if !s.NoError(err) {
			return
		}

		s.Len(apps, 1)
		s.Equal(int64(1), count)
		s.Equal("test", apps[0].Name)

		s.T().Logf("apps: %s", spew.Sdump(apps))
	})

	s.Run("list error", func() {
		e := fmt.Errorf("db error")
		// mock sql
		s.sqlmock.ExpectQuery("SELECT (.+) FROM `application`").
			WillReturnError(e)

		// list apps
		apps, count, err := s.storage.ListApps(s.ctx, 0, consts.PublishStatusAll)
		if !s.Error(err) {
			return
		}

		if s.Error(err) {
			s.EqualError(err, "list application error: db error")
			s.Nil(apps)
			s.Equal(int64(0), count)
		}

		s.T().Logf("apps: %s", spew.Sdump(apps))
	})
}
