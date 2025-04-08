package job

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/mitchellh/mapstructure"
	"github.com/ory/viper"
	"github.com/stretchr/testify/suite"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	app "github.com/yuansuan/ticp/iPaaS/project-root/internal/job/api/v1/application"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/config"
	appSrv "github.com/yuansuan/ticp/iPaaS/project-root/internal/job/service/v1/application"
	jobSrv "github.com/yuansuan/ticp/iPaaS/project-root/internal/job/service/v1/job"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/util"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/common/mock"
)

// suite
type JobSuite struct {
	// suite
	suite.Suite

	// base
	ctrl *gomock.Controller

	w      *httptest.ResponseRecorder
	ctx    *gin.Context
	logger *logging.Logger

	// service
	mockAppSrv      *appSrv.MockService
	mockAppASrv     *appSrv.MockAppSrv
	mockAppQuotaSrv *appSrv.MockAppQuotaSrv
	mockAppAllowSrv *appSrv.MockAppAllowSrv

	mockJobSrv *jobSrv.MockService

	// field
	mockUserChecker *util.MockUserChecker

	// handler
	handler *Handler
}

func (s *JobSuite) SetupSuite() {
	s.T().Log("setup suite")

	// mock base controller
	s.ctrl = gomock.NewController(s.T())

	// mock gin context
	s.w = httptest.NewRecorder()
	s.ctx = mock.GinContext(s.w)

	// logger
	logger, err := logging.NewLogger(logging.WithDefaultLogConfigOption())
	if !s.NoError(err) {
		return
	}
	s.logger = logger
	s.ctx.Set(logging.LoggerName, s.logger)

	// mock service
	mockAppSrv := appSrv.NewMockService(s.ctrl)
	mockAppASrv := appSrv.NewMockAppSrv(s.ctrl)
	mockAppQuotaSrv := appSrv.NewMockAppQuotaSrv(s.ctrl)
	mockAppAllowSrv := appSrv.NewMockAppAllowSrv(s.ctrl)

	mockAppSrv.EXPECT().Apps().Return(mockAppASrv).AnyTimes()
	mockAppSrv.EXPECT().AppsQuota().Return(mockAppQuotaSrv).AnyTimes()
	mockAppSrv.EXPECT().AppsAllow().Return(mockAppAllowSrv).AnyTimes()

	s.mockAppSrv = mockAppSrv
	s.mockAppASrv = mockAppASrv
	s.mockAppQuotaSrv = mockAppQuotaSrv
	s.mockAppAllowSrv = mockAppAllowSrv
	s.mockJobSrv = jobSrv.NewMockService(s.ctrl)

	// mock field
	mockUserChecker := util.NewMockUserChecker(s.ctrl)
	s.mockUserChecker = mockUserChecker

	// handler
	mockController := app.MockApplicationController(s.mockAppSrv)
	s.handler = NewJobHandler(mockController, s.mockJobSrv, s.mockUserChecker)
}

func (s *JobSuite) TearDownSuite() {
	s.T().Log("teardown suite")
}

func (s *JobSuite) SetupTest() {
	s.T().Log("setup test")

	s.w = httptest.NewRecorder()
	s.ctx = mock.GinContext(s.w)
	s.ctx.Set(logging.LoggerName, s.logger)
	s.ctx.Request.Header.Set("x-ys-request-id", gofakeit.UUID())
}

func (s *JobSuite) TearDownTest() {
	s.T().Log("teardown suite")
}

func (s *JobSuite) SetupSubTest() {
	s.T().Log("setup sub test")
	s.w = httptest.NewRecorder()
	s.ctx = mock.GinContext(s.w)
	s.ctx.Set(logging.LoggerName, s.logger)
	s.ctx.Request.Header.Set("x-ys-request-id", gofakeit.UUID())
}

func (s *JobSuite) TearDownSubTest() {
	s.T().Log("teardown sub test")
}

func (s *JobSuite) mockConfig() {
	const cfgStr = `
zones:
  az-jinan:
    hpc_endpoint: https://jn_hpc_endpoint:8080
    storage_endpoint: https://jn_storage_endpoint:8899
    cloud_app_enable: true
  az-wuxi:
    storage_endpoint: https://wx_storage_endpoint:8899
    cloud_app_enable: false
bill_enabled: true
`
	viper.SetConfigType("yaml")
	err := viper.ReadConfig(strings.NewReader(cfgStr))
	if !s.NoError(err) {
		s.FailNow("read config failed")
	}
	md := mapstructure.Metadata{}
	customT := config.CustomT{}
	err = viper.Unmarshal(&customT, func(config *mapstructure.DecoderConfig) {
		config.TagName = "yaml"
		config.Metadata = &md
	})
	if !s.NoError(err) {
		s.FailNow("unmarshal config failed")
	}
	config.Custom = customT
}

type JobTestCase struct {
	Name              string
	ExpectedHTTPCode  int
	ExpectedErrorCode string
	MockExpectFunc    func()
	ErrorCode         string
	ErrorMsg          string
}

func (tc *JobTestCase) MakeTestFunc(req, do, resp func()) func(t *testing.T, s *JobSuite) {
	return func(t *testing.T, s *JobSuite) {
		s.Run(tc.Name, func() {
			// Mock EXPECT
			if tc.MockExpectFunc != nil {
				tc.MockExpectFunc()
			}

			if req != nil {
				req()
			}

			if do != nil {
				do()
			}

			if resp != nil {
				resp()
			}

			// assert
			if tc.ExpectedHTTPCode != 0 {
				s.Equal(tc.ExpectedHTTPCode, s.w.Code)
			}

			if tc.ExpectedErrorCode != "" {
				if !s.Equal(tc.ExpectedErrorCode, tc.ErrorCode) {
					t.Logf("ErrorMsg: %s", tc.ErrorMsg)
				}
			}
		})
	}
}
