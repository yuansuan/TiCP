package job

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	api "github.com/yuansuan/ticp/common/project-root-api/common"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/admin/jobmonitorchart"
	schema "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao/models"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/util"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/common/mock"
)

type JobAdminMonitorChartSuite struct {
	JobSuite

	// data
	mockMonitorChart *models.MonitorChart
}

func TestJobAdminMonitorChart(t *testing.T) {
	suite.Run(t, new(JobAdminMonitorChartSuite))
}

func (s *JobAdminMonitorChartSuite) SetupTest() {
	s.JobSuite.SetupTest()
	s.mockMonitorChart = &models.MonitorChart{
		ID:    0,
		JobID: snowflake.ID(12345),
		Content: func() string {
			r, _ := util.MonitorChartMarshal([]*schema.MonitorChart{{
				Key: "chart_key",
				Items: []*schema.MonitorChartItem{{
					Kv: []float64{1.0, 2.0, 3.0},
				}},
			}})
			return r
		}(),
		Finished:           false,
		MonitorChartRegexp: "",
		MonitorChartParser: "",
		FailedReason:       "",
	}
}

// TestCase structure
type jobAdminMonitorChartTestCase struct {
	JobTestCase
	pathParams []gin.Param
}

func (s *JobAdminMonitorChartSuite) reqFunc(pathParams []gin.Param) func() {
	return func() {
		mock.HTTPRequest(s.ctx, http.MethodGet, nil, pathParams, nil)
	}
}

func (s *JobAdminMonitorChartSuite) doFunc() {
	s.handler.AdminMonitorChart(s.ctx)
}

func (s *JobAdminMonitorChartSuite) respFunc(tc *jobAdminMonitorChartTestCase) func() {
	return func() {
		got := s.w.Body.Bytes()
		var resp jobmonitorchart.Response

		// unmarshal
		err := json.Unmarshal(got, &resp)
		if !s.NoError(err) {
			return
		}
		s.T().Log(spew.Sdump(resp))

		tc.ErrorCode = resp.ErrorCode
		tc.ErrorMsg = resp.ErrorMsg
	}
}

func (s *JobAdminMonitorChartSuite) TestGetMonitorChart() {
	testCases := []jobAdminMonitorChartTestCase{
		{
			JobTestCase: JobTestCase{
				Name:              "test admin get monitor chart job success",
				ExpectedHTTPCode:  http.StatusOK,
				ExpectedErrorCode: "",
				MockExpectFunc: func() {
					s.mockJobSrv.EXPECT().GetMonitorChart(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(s.mockMonitorChart, nil)
				},
			},
			pathParams: []gin.Param{{Key: "JobID", Value: "4ER"}},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test admin get monitor char bind error",
				ExpectedHTTPCode:  http.StatusBadRequest,
				ExpectedErrorCode: api.InvalidJobID,
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test admin get monitor char invalid job id",
				ExpectedHTTPCode:  http.StatusBadRequest,
				ExpectedErrorCode: api.InvalidJobID,
			},
			pathParams: []gin.Param{{Key: "JobID", Value: "O"}}, // O is invalid

		},
		{
			JobTestCase: JobTestCase{
				Name:              "test admin get monitor chart job id not found",
				ExpectedHTTPCode:  http.StatusNotFound,
				ExpectedErrorCode: api.JobIDNotFound,
				MockExpectFunc: func() {
					s.mockJobSrv.EXPECT().GetMonitorChart(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, common.ErrJobIDNotFound)
				},
			},
			pathParams: []gin.Param{{Key: "JobID", Value: "NtFd"}},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test admin get monitor chart not found",
				ExpectedHTTPCode:  http.StatusNotFound,
				ExpectedErrorCode: api.JobMonitorChartNotFound,
				MockExpectFunc: func() {
					s.mockJobSrv.EXPECT().GetMonitorChart(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, common.ErrJobMonitorChartNotFound)
				},
			},
			pathParams: []gin.Param{{Key: "JobID", Value: "4ER"}},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test admin get monitor chart internal error",
				ExpectedHTTPCode:  http.StatusInternalServerError,
				ExpectedErrorCode: api.InternalServerErrorCode,
				MockExpectFunc: func() {
					s.mockJobSrv.EXPECT().GetMonitorChart(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, common.ErrInternalServer)
				},
			},
			pathParams: []gin.Param{{Key: "JobID", Value: "4ER"}},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test admin get monitor chart unmarshal error",
				ExpectedHTTPCode:  http.StatusInternalServerError,
				ExpectedErrorCode: api.InternalServerErrorCode,
				MockExpectFunc: func() {
					s.mockJobSrv.EXPECT().GetMonitorChart(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&models.MonitorChart{
						Content: "errormsg",
					}, nil)
				},
			},
			pathParams: []gin.Param{{Key: "JobID", Value: "4ER"}},
		},
	}

	for _, tc := range testCases {
		testJobAdminMonitorChart := tc.MakeTestFunc(s.reqFunc(tc.pathParams), s.doFunc, s.respFunc(&tc))
		testJobAdminMonitorChart(s.T(), &s.JobSuite)
	}
}
