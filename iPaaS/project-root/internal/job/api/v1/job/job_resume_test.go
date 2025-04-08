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
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/jobresume"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/common/mock"
)

// JobResumeSuite
type JobResumeSuite struct {
	JobSuite

	// data
}

func TestJobResume(t *testing.T) {
	suite.Run(t, new(JobResumeSuite))
}

type jobResumeTestCase struct {
	JobTestCase
	pathParams    []gin.Param
	setHeaderFunc func(ctx *gin.Context)
}

func (s *JobResumeSuite) reqFunc(pathParams []gin.Param, setHeaderFunc func(ctx *gin.Context)) func() {
	return func() {
		// Set custom headers
		if setHeaderFunc != nil {
			setHeaderFunc(s.ctx)
		}

		mock.HTTPRequest(s.ctx, http.MethodPatch, nil, pathParams, nil)
	}
}

func (s *JobResumeSuite) doFunc() {
	s.handler.Resume(s.ctx)
}

func (s *JobResumeSuite) respFunc(tc *jobResumeTestCase) func() {
	return func() {
		got := s.w.Body.Bytes()
		var resp jobresume.Response

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

func (s *JobResumeSuite) TestResume() {
	testCases := []jobResumeTestCase{
		{
			JobTestCase: JobTestCase{
				Name:              "test resume job success",
				ExpectedHTTPCode:  http.StatusOK,
				ExpectedErrorCode: "",
				MockExpectFunc: func() {
					s.mockJobSrv.EXPECT().Resume(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				},
			},
			pathParams: []gin.Param{{Key: "JobID", Value: "4ER"}},
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "h9z")
			},
		},
	}

	for _, tc := range testCases {
		testJobResume := tc.MakeTestFunc(s.reqFunc(tc.pathParams, tc.setHeaderFunc), s.doFunc, s.respFunc(&tc))
		testJobResume(s.T(), &s.JobSuite)
	}
}

func (s *JobResumeSuite) TestResumeValidate() {
	testCases := []jobResumeTestCase{
		{
			JobTestCase: JobTestCase{
				Name:              "test resume job bind error",
				ExpectedHTTPCode:  http.StatusBadRequest,
				ExpectedErrorCode: api.InvalidArgumentErrorCode,
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test resume job invalid user id",
				ExpectedHTTPCode:  http.StatusBadRequest,
				ExpectedErrorCode: api.InvalidUserID,
			},
			pathParams: []gin.Param{{Key: "JobID", Value: "4ER"}},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test resume job invalid job id",
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
		testJobResume := tc.MakeTestFunc(s.reqFunc(tc.pathParams, tc.setHeaderFunc), s.doFunc, s.respFunc(&tc))
		testJobResume(s.T(), &s.JobSuite)
	}
}

func (s *JobResumeSuite) TestResumeJudge() {
	testCases := []jobResumeTestCase{
		{
			JobTestCase: JobTestCase{
				Name:              "test resume job id not found",
				ExpectedHTTPCode:  http.StatusNotFound,
				ExpectedErrorCode: api.JobIDNotFound,
				MockExpectFunc: func() {
					s.mockJobSrv.EXPECT().Resume(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(common.ErrJobIDNotFound)
				},
			},
			pathParams: []gin.Param{{Key: "JobID", Value: "NtFd"}},
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "h9z")
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test resume job access denied",
				ExpectedHTTPCode:  http.StatusForbidden,
				ExpectedErrorCode: api.JobAccessDenied,
				MockExpectFunc: func() {
					s.mockJobSrv.EXPECT().Resume(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(common.ErrJobAccessDenied)
				},
			},
			pathParams: []gin.Param{{Key: "JobID", Value: "4ER"}},
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "vmR")
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test resume job state not allow resume",
				ExpectedHTTPCode:  http.StatusForbidden,
				ExpectedErrorCode: api.JobStateNotAllowResume,
				MockExpectFunc: func() {
					s.mockJobSrv.EXPECT().Resume(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(common.ErrJobStateNotAllowResume)
				},
			},
			pathParams: []gin.Param{{Key: "JobID", Value: "4ER"}},
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "vmR")
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test resume job internal error",
				ExpectedHTTPCode:  http.StatusInternalServerError,
				ExpectedErrorCode: api.InternalServerErrorCode,
				MockExpectFunc: func() {
					s.mockJobSrv.EXPECT().Resume(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(common.ErrInternalServer)
				},
			},
			pathParams: []gin.Param{{Key: "JobID", Value: "4ER"}},
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "h9z")
			},
		},
	}

	for _, tc := range testCases {
		testJobResume := tc.MakeTestFunc(s.reqFunc(tc.pathParams, tc.setHeaderFunc), s.doFunc, s.respFunc(&tc))
		testJobResume(s.T(), &s.JobSuite)
	}
}
