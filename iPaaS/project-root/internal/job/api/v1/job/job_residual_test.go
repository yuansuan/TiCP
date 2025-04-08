package job

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	api "github.com/yuansuan/ticp/common/project-root-api/common"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/jobresidual"
	schema "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/common/mock"
)

type JobResidualSuite struct {
	JobSuite

	// data
	mockResidual *schema.Residual
}

func TestJobResidual(t *testing.T) {
	suite.Run(t, new(JobResidualSuite))
}

func (s *JobResidualSuite) SetupTest() {
	s.JobSuite.SetupTest()
	s.mockResidual = &schema.Residual{
		Vars: []*schema.ResidualVar{
			{
				Values: []float64{1.0, 2.0, 3.0},
				Name:   "var1",
			},
		},
		AvailableXvar: []string{"x1", "x2", "x3"},
	}
}

// TestCase structure
type jobResidualTestCase struct {
	JobTestCase
	pathParams    []gin.Param
	setHeaderFunc func(ctx *gin.Context)
}

func (s *JobResidualSuite) reqFunc(pathParams []gin.Param, setHeaderFunc func(ctx *gin.Context)) func() {
	return func() {
		// Set custom headers
		if setHeaderFunc != nil {
			setHeaderFunc(s.ctx)
		}

		mock.HTTPRequest(s.ctx, http.MethodGet, nil, pathParams, nil)
	}
}

func (s *JobResidualSuite) doFunc() {
	s.handler.Residual(s.ctx)
}

func (s *JobResidualSuite) respFunc(tc *jobResidualTestCase) func() {
	return func() {
		got := s.w.Body.Bytes()
		var resp jobresidual.Response

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

func (s *JobResidualSuite) TestResidual() {
	testCases := []jobResidualTestCase{
		{
			JobTestCase: JobTestCase{
				Name:              "test get residual success",
				ExpectedHTTPCode:  http.StatusOK,
				ExpectedErrorCode: "",
				MockExpectFunc: func() {
					s.mockJobSrv.EXPECT().GetResidual(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(s.mockResidual, nil)
				},
			},
			pathParams: []gin.Param{{Key: "JobID", Value: "4ER"}},
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "h9z")
			},
		},
	}

	for _, tc := range testCases {
		testJobResidual := tc.MakeTestFunc(s.reqFunc(tc.pathParams, tc.setHeaderFunc), s.doFunc, s.respFunc(&tc))
		testJobResidual(s.T(), &s.JobSuite)
	}
}

func (s *JobResidualSuite) TestResidualValidate() {
	testCases := []jobResidualTestCase{
		{
			JobTestCase: JobTestCase{
				Name:              "test get residual bind error",
				ExpectedHTTPCode:  http.StatusBadRequest,
				ExpectedErrorCode: api.InvalidJobID,
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test get residual invalid user id",
				ExpectedHTTPCode:  http.StatusBadRequest,
				ExpectedErrorCode: api.InvalidUserID,
			},
			pathParams: []gin.Param{{Key: "JobID", Value: "4ER"}},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test get residual invalid job id",
				ExpectedHTTPCode:  http.StatusBadRequest,
				ExpectedErrorCode: api.InvalidJobID,
			},
			pathParams: []gin.Param{{Key: "JobID", Value: "O"}}, // O is invalid
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "h9z")
			},
		},
	}

	for _, tc := range testCases {
		testJobResidual := tc.MakeTestFunc(s.reqFunc(tc.pathParams, tc.setHeaderFunc), s.doFunc, s.respFunc(&tc))
		testJobResidual(s.T(), &s.JobSuite)
	}
}

func (s *JobResidualSuite) TestResidualJudge() {
	testCases := []jobResidualTestCase{
		{
			JobTestCase: JobTestCase{
				Name:              "test get residual job id not found",
				ExpectedHTTPCode:  http.StatusNotFound,
				ExpectedErrorCode: api.JobIDNotFound,
				MockExpectFunc: func() {
					s.mockJobSrv.EXPECT().GetResidual(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, common.ErrJobIDNotFound)
				},
			},
			pathParams: []gin.Param{{Key: "JobID", Value: "NtFd"}},
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "h9z")
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test get residual access denied",
				ExpectedHTTPCode:  http.StatusForbidden,
				ExpectedErrorCode: api.JobAccessDenied,
				MockExpectFunc: func() {
					s.mockJobSrv.EXPECT().GetResidual(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, common.ErrJobAccessDenied)
				},
			},
			pathParams: []gin.Param{{Key: "JobID", Value: "4ER"}},
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "h9z")
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test get residual not found",
				ExpectedHTTPCode:  http.StatusNotFound,
				ExpectedErrorCode: api.JobResidualNotFound,
				MockExpectFunc: func() {
					s.mockJobSrv.EXPECT().GetResidual(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, common.ErrJobResidualNotFound)
				},
			},
			pathParams: []gin.Param{{Key: "JobID", Value: "4ER"}},
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "h9z")
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test get residual internal error",
				ExpectedHTTPCode:  http.StatusInternalServerError,
				ExpectedErrorCode: api.InternalServerErrorCode,
				MockExpectFunc: func() {
					s.mockJobSrv.EXPECT().GetResidual(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, common.ErrInternalServer)
				},
			},
			pathParams: []gin.Param{{Key: "JobID", Value: "4ER"}},
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "h9z")
			},
		},
	}

	for _, tc := range testCases {
		testJobResidual := tc.MakeTestFunc(s.reqFunc(tc.pathParams, tc.setHeaderFunc), s.doFunc, s.respFunc(&tc))
		testJobResidual(s.T(), &s.JobSuite)
	}
}
