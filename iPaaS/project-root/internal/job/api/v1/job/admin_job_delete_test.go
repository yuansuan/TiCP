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
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/admin/jobdelete"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/common/mock"
)

// JobAdminDeleteSuite
type JobAdminDeleteSuite struct {
	JobSuite

	// data
}

func TestJobAdminDelete(t *testing.T) {
	suite.Run(t, new(JobAdminDeleteSuite))
}

type jobAdminDeleteTestCase struct {
	JobTestCase
	pathParams []gin.Param
}

func (s *JobAdminDeleteSuite) reqFunc(pathParams []gin.Param) func() {
	return func() {
		mock.HTTPRequest(s.ctx, http.MethodDelete, nil, pathParams, nil)
	}
}

func (s *JobAdminDeleteSuite) doFunc() {
	s.handler.AdminDelete(s.ctx)
}

func (s *JobAdminDeleteSuite) respFunc(tc *jobAdminDeleteTestCase) func() {
	return func() {
		got := s.w.Body.Bytes()
		var resp jobdelete.Response

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

func (s *JobAdminDeleteSuite) TestAdminDelete() {
	testCase := []jobAdminDeleteTestCase{
		{
			JobTestCase: JobTestCase{
				Name:              "test admin delete job success",
				ExpectedHTTPCode:  http.StatusOK,
				ExpectedErrorCode: "",
				MockExpectFunc: func() {
					s.mockJobSrv.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				},
			},
			pathParams: []gin.Param{{Key: "JobID", Value: "4ER"}},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test admin delete job bind error",
				ExpectedHTTPCode:  http.StatusBadRequest,
				ExpectedErrorCode: api.InvalidArgumentErrorCode,
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test admin delete job invalid job id",
				ExpectedHTTPCode:  http.StatusBadRequest,
				ExpectedErrorCode: api.InvalidJobID,
			},
			pathParams: []gin.Param{{Key: "JobID", Value: "O"}}, // O is invalid
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test admin delete job id not found",
				ExpectedHTTPCode:  http.StatusNotFound,
				ExpectedErrorCode: api.JobIDNotFound,
				MockExpectFunc: func() {
					s.mockJobSrv.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(common.ErrJobIDNotFound)
				},
			},
			pathParams: []gin.Param{{Key: "JobID", Value: "NtFd"}},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test admin delete job state not allow delete",
				ExpectedHTTPCode:  http.StatusForbidden,
				ExpectedErrorCode: api.JobStateNotAllowDelete,
				MockExpectFunc: func() {
					s.mockJobSrv.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(common.ErrJobStateNotAllowDelete)
				},
			},
			pathParams: []gin.Param{{Key: "JobID", Value: "4ER"}},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test admin delete job internal error",
				ExpectedHTTPCode:  http.StatusInternalServerError,
				ExpectedErrorCode: api.InternalServerErrorCode,
				MockExpectFunc: func() {
					s.mockJobSrv.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(common.ErrInternalServer)
				},
			},
			pathParams: []gin.Param{{Key: "JobID", Value: "4ER"}},
		},
	}

	for _, tc := range testCase {
		testJobAdminDelete := tc.MakeTestFunc(s.reqFunc(tc.pathParams), s.doFunc, s.respFunc(&tc))
		testJobAdminDelete(s.T(), &s.JobSuite)
	}
}
