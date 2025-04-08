package job

import (
	"context"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/middleware"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao/models"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/common/mock"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/common/with"
	"xorm.io/xorm"
)

type jobServiceSuite struct {
	// suite
	suite.Suite

	// base
	ctrl    *gomock.Controller
	ctx     context.Context
	engine  *xorm.Engine
	sqlmock sqlmock.Sqlmock

	// field
	mockJobDao      *dao.MockJobDao
	mockResidualDao *dao.MockResidualDao
	mockIDGen       *snowflake.MockIDGen

	// service
	jobSrv *jobService
}

func (s *jobServiceSuite) SetupSuite() {
	s.T().Log("setup suite")

	// mock base controller
	s.ctrl = gomock.NewController(s.T())

	// mock gin context
	w := httptest.NewRecorder()
	ctx := mock.GinContext(w)

	// logger
	logger, err := logging.NewLogger(logging.WithReleaseLevel(logging.DevelopmentLevel))
	if err != nil {
		panic(fmt.Sprintf("init logger failed: %v", err))
	}
	ctx.Set(logging.LoggerName, logger)

	// mock mysql engine
	mockEngine, mocksql := mock.Engine()
	mockEngine.SetLogger(middleware.NewXormLogger(logger, true))
	s.engine = mockEngine
	s.sqlmock = mocksql

	db := s.engine.Context(ctx)
	s.ctx = with.KeepSession(ctx, db)

	// mock store
	mockJobDao := dao.NewMockJobDao(s.ctrl)
	s.mockJobDao = mockJobDao
	mockResidualDao := dao.NewMockResidualDao(s.ctrl)
	s.mockResidualDao = mockResidualDao

	// mock idgen
	mockIDGen := snowflake.NewMockIDGen(s.ctrl)
	s.mockIDGen = mockIDGen

	// service
	srv := &jobService{
		IDGen:       mockIDGen,
		jobdao:      mockJobDao,
		residualdao: mockResidualDao,
		jobPlugins:  []JobPlugin{&mockJobPlugin{}}, // mock job plugins
	}
	s.jobSrv = srv
}

type mockJobPlugin struct{}

func (m *mockJobPlugin) Insert(ctx context.Context, app *models.Application, job *models.Job) {
	logger := logging.GetLogger(ctx)
	logger.Infof("job %s, %s inserted", m.Name(), job.ID.String())
}

func (m *mockJobPlugin) Name() string {
	return "mockJobPlugin"
}

func TestNewJobService(t *testing.T) {
	NewJobService(nil, nil, nil, nil)
}
