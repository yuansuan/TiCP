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
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/admin/jobresidual"
	schema "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/common/mock"
)

type JobAdminResidualSuite struct {
	JobSuite

	// data
	mockResidual *schema.Residual
}

func TestJobAdminResidual(t *testing.T) {
	suite.Run(t, new(JobAdminResidualSuite))
}

func (s *JobAdminResidualSuite) SetupTest() {
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
type jobAdminResidualTestCase struct {
	JobTestCase
	pathParams []gin.Param
}

func (s *JobAdminResidualSuite) reqFunc(pathParams []gin.Param) func() {
	return func() {
		mock.HTTPRequest(s.ctx, http.MethodGet, nil, pathParams, nil)
	}
}

func (s *JobAdminResidualSuite) doFunc() {
	s.handler.AdminResidual(s.ctx)
}

func (s *JobAdminResidualSuite) respFunc(tc *jobAdminResidualTestCase) func() {
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

func (s *JobAdminResidualSuite) TestAdminResidual() {
	testCases := []jobAdminResidualTestCase{
		{
			JobTestCase: JobTestCase{
				Name:              "test admin get residual success",
				ExpectedHTTPCode:  http.StatusOK,
				ExpectedErrorCode: "",
				MockExpectFunc: func() {
					s.mockJobSrv.EXPECT().GetResidual(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(s.mockResidual, nil)
				},
			},
			pathParams: []gin.Param{{Key: "JobID", Value: "4ER"}},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test admin get residual bind error",
				ExpectedHTTPCode:  http.StatusBadRequest,
				ExpectedErrorCode: api.InvalidJobID,
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test admin get residual invalid job id",
				ExpectedHTTPCode:  http.StatusBadRequest,
				ExpectedErrorCode: api.InvalidJobID,
			},
			pathParams: []gin.Param{{Key: "JobID", Value: "O"}}, // O is invalid
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test admin get residual job id not found",
				ExpectedHTTPCode:  http.StatusNotFound,
				ExpectedErrorCode: api.JobIDNotFound,
				MockExpectFunc: func() {
					s.mockJobSrv.EXPECT().GetResidual(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, common.ErrJobIDNotFound)
				},
			},
			pathParams: []gin.Param{{Key: "JobID", Value: "NtFd"}},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test admin get residual not found",
				ExpectedHTTPCode:  http.StatusNotFound,
				ExpectedErrorCode: api.JobResidualNotFound,
				MockExpectFunc: func() {
					s.mockJobSrv.EXPECT().GetResidual(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, common.ErrJobResidualNotFound)
				},
			},
			pathParams: []gin.Param{{Key: "JobID", Value: "4ER"}},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test admin get residual internal error",
				ExpectedHTTPCode:  http.StatusInternalServerError,
				ExpectedErrorCode: api.InternalServerErrorCode,
				MockExpectFunc: func() {
					s.mockJobSrv.EXPECT().GetResidual(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, common.ErrInternalServer)
				},
			},
			pathParams: []gin.Param{{Key: "JobID", Value: "4ER"}},
		},
	}

	for _, tc := range testCases {
		testJobAdminResidual := tc.MakeTestFunc(s.reqFunc(tc.pathParams), s.doFunc, s.respFunc(&tc))
		testJobAdminResidual(s.T(), &s.JobSuite)
	}
}
