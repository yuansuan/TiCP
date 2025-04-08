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
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/admin/jobterminate"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/common/mock"
)

// JobJobTerminateSuite
type JobJobTerminateSuite struct {
	JobSuite

	// data
}

func TestJobJobTerminate(t *testing.T) {
	suite.Run(t, new(JobJobTerminateSuite))
}

type jobAdminTerminateTestCase struct {
	JobTestCase
	pathParams []gin.Param
}

func (s *JobJobTerminateSuite) reqFunc(pathParams []gin.Param) func() {
	return func() {
		mock.HTTPRequest(s.ctx, http.MethodPatch, nil, pathParams, nil)
	}
}

func (s *JobJobTerminateSuite) doFunc() {
	s.handler.AdminTerminate(s.ctx)
}

func (s *JobJobTerminateSuite) respFunc(tc *jobAdminTerminateTestCase) func() {
	return func() {
		got := s.w.Body.Bytes()
		var resp jobterminate.Response

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

func (s *JobJobTerminateSuite) TestTerminate() {
	testCase := []jobAdminTerminateTestCase{
		{
			JobTestCase: JobTestCase{
				Name:              "test admin terminate job success",
				ExpectedHTTPCode:  http.StatusOK,
				ExpectedErrorCode: "",
				MockExpectFunc: func() {
					s.mockJobSrv.EXPECT().Terminate(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				},
			},
			pathParams: []gin.Param{{Key: "JobID", Value: "4ER"}},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test admin terminate job bind error",
				ExpectedHTTPCode:  http.StatusBadRequest,
				ExpectedErrorCode: api.InvalidArgumentErrorCode,
				MockExpectFunc:    func() {},
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test admin terminate job invalid job id",
				ExpectedHTTPCode:  http.StatusBadRequest,
				ExpectedErrorCode: api.InvalidJobID,
				MockExpectFunc:    func() {},
			},
			pathParams: []gin.Param{{Key: "JobID", Value: "O"}}, // O is invalid
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test admin terminate job id not found",
				ExpectedHTTPCode:  http.StatusNotFound,
				ExpectedErrorCode: api.JobIDNotFound,
				MockExpectFunc: func() {
					s.mockJobSrv.EXPECT().Terminate(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(common.ErrJobIDNotFound)
				},
			},
			pathParams: []gin.Param{{Key: "JobID", Value: "NtFd"}},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test admin terminate job state not allow terminate",
				ExpectedHTTPCode:  http.StatusForbidden,
				ExpectedErrorCode: api.JobStateNotAllowTerminate,
				MockExpectFunc: func() {
					s.mockJobSrv.EXPECT().Terminate(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(common.ErrJobStateNotAllowTerminate)
				},
			},
			pathParams: []gin.Param{{Key: "JobID", Value: "4ER"}},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test admin terminate job internal error",
				ExpectedHTTPCode:  http.StatusInternalServerError,
				ExpectedErrorCode: api.InternalServerErrorCode,
				MockExpectFunc: func() {
					s.mockJobSrv.EXPECT().Terminate(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(common.ErrInternalServer)
				},
			},
			pathParams: []gin.Param{{Key: "JobID", Value: "4ER"}},
		},
	}

	for _, tc := range testCase {
		testJobAdminTerminate := tc.MakeTestFunc(s.reqFunc(tc.pathParams), s.doFunc, s.respFunc(&tc))
		testJobAdminTerminate(s.T(), &s.JobSuite)
	}
}
