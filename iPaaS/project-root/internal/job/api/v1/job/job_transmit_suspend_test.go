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
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/jobtransmitsuspend"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/common/mock"
)

// JobTransmitSuspendSuite
type JobTransmitSuspendSuite struct {
	JobSuite

	// data
}

func TestJobTransmitSuspend(t *testing.T) {
	suite.Run(t, new(JobTransmitSuspendSuite))
}

// TestCase structure
type jobTransmitSuspendTestCase struct {
	JobTestCase
	pathParams    []gin.Param
	setHeaderFunc func(ctx *gin.Context)
}

func (s *JobTransmitSuspendSuite) reqFunc(pathParams []gin.Param, setHeaderFunc func(ctx *gin.Context)) func() {
	return func() {
		// Set custom headers
		if setHeaderFunc != nil {
			setHeaderFunc(s.ctx)
		}

		mock.HTTPRequest(s.ctx, http.MethodPatch, nil, pathParams, nil)
	}
}

func (s *JobTransmitSuspendSuite) doFunc() {
	s.handler.TransmitSuspend(s.ctx)
}

func (s *JobTransmitSuspendSuite) respFunc(tc *jobTransmitSuspendTestCase) func() {
	return func() {
		got := s.w.Body.Bytes()
		var resp jobtransmitsuspend.Response

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

func (s *JobTransmitSuspendSuite) TestTransmitSuspend() {
	testCases := []jobTransmitSuspendTestCase{
		{
			JobTestCase: JobTestCase{
				Name:              "test transmit suspend job success",
				ExpectedHTTPCode:  http.StatusOK,
				ExpectedErrorCode: "",
				MockExpectFunc: func() {
					s.mockJobSrv.EXPECT().TransmitSuspend(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				},
			},
			pathParams: []gin.Param{{Key: "JobID", Value: "4ER"}},
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "h9z")
			},
		},
	}

	for _, tc := range testCases {
		testJobTransmitSuspend := tc.MakeTestFunc(s.reqFunc(tc.pathParams, tc.setHeaderFunc), s.doFunc, s.respFunc(&tc))
		testJobTransmitSuspend(s.T(), &s.JobSuite)
	}
}

func (s *JobTransmitSuspendSuite) TestTransmitSuspendValidate() {
	testCases := []jobTransmitSuspendTestCase{
		{
			JobTestCase: JobTestCase{
				Name:              "test transmit suspend job bind error",
				ExpectedHTTPCode:  http.StatusBadRequest,
				ExpectedErrorCode: api.InvalidArgumentErrorCode,
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test transmit suspend job invalid user id",
				ExpectedHTTPCode:  http.StatusBadRequest,
				ExpectedErrorCode: api.InvalidUserID,
			},
			pathParams: []gin.Param{{Key: "JobID", Value: "4ER"}},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test transmit suspend job invalid job id",
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
		testJobTransmitSuspend := tc.MakeTestFunc(s.reqFunc(tc.pathParams, tc.setHeaderFunc), s.doFunc, s.respFunc(&tc))
		testJobTransmitSuspend(s.T(), &s.JobSuite)
	}
}

func (s *JobTransmitSuspendSuite) TestTransmitSuspendJudge() {
	testCases := []jobTransmitSuspendTestCase{
		{
			JobTestCase: JobTestCase{
				Name:              "test transmit suspend job id not found",
				ExpectedHTTPCode:  http.StatusNotFound,
				ExpectedErrorCode: api.JobIDNotFound,
				MockExpectFunc: func() {
					s.mockJobSrv.EXPECT().TransmitSuspend(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(common.ErrJobIDNotFound)
				},
			},
			pathParams: []gin.Param{{Key: "JobID", Value: "NtFd"}},
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "h9z")
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test transmit suspend job access denied",
				ExpectedHTTPCode:  http.StatusForbidden,
				ExpectedErrorCode: api.JobAccessDenied,
				MockExpectFunc: func() {
					s.mockJobSrv.EXPECT().TransmitSuspend(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(common.ErrJobAccessDenied)
				},
			},
			pathParams: []gin.Param{{Key: "JobID", Value: "4ER"}},
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "vmR")
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test transmit suspend job state not allow transmit suspend",
				ExpectedHTTPCode:  http.StatusForbidden,
				ExpectedErrorCode: api.JobStateNotAllowTransmitSuspend,
				MockExpectFunc: func() {
					s.mockJobSrv.EXPECT().TransmitSuspend(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(common.ErrJobStateNotAllowTransmitSuspend)
				},
			},
			pathParams: []gin.Param{{Key: "JobID", Value: "4ER"}},
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "vmR")
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test transmit suspend job internal error",
				ExpectedHTTPCode:  http.StatusInternalServerError,
				ExpectedErrorCode: api.InternalServerErrorCode,
				MockExpectFunc: func() {
					s.mockJobSrv.EXPECT().TransmitSuspend(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(common.ErrInternalServer)
				},
			},
			pathParams: []gin.Param{{Key: "JobID", Value: "4ER"}},
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "h9z")
			},
		},
	}

	for _, tc := range testCases {
		testJobTransmitSuspend := tc.MakeTestFunc(s.reqFunc(tc.pathParams, tc.setHeaderFunc), s.doFunc, s.respFunc(&tc))
		testJobTransmitSuspend(s.T(), &s.JobSuite)
	}
}
