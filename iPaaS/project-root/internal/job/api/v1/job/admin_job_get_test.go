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
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/admin/jobget"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao/models"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/common/mock"
)

type JobAdminGetSuite struct {
	JobSuite

	// data
	mockJob *models.Job
}

func TestJobAdminGet(t *testing.T) {
	suite.Run(t, new(JobAdminGetSuite))
}

func (s *JobAdminGetSuite) SetupTest() {
	s.JobSuite.SetupTest()
	s.mockJob = &models.Job{
		ID:     snowflake.ID(12345),
		UserID: snowflake.ID(54321),
	}
}

type jobAdminGetTestCase struct {
	JobTestCase
	pathParams []gin.Param
}

func (s *JobAdminGetSuite) reqFunc(pathParams []gin.Param) func() {
	return func() {
		mock.HTTPRequest(s.ctx, http.MethodGet, nil, pathParams, nil)
	}
}

func (s *JobAdminGetSuite) doFunc() {
	s.handler.AdminGet(s.ctx)
}

func (s *JobAdminGetSuite) respFunc(tc *jobAdminGetTestCase) func() {
	return func() {
		got := s.w.Body.Bytes()
		var resp jobget.Response

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

func (s *JobAdminGetSuite) TestAdminGet() {
	testCase := []jobAdminGetTestCase{
		{
			JobTestCase: JobTestCase{
				Name:              "test admin get job success",
				ExpectedHTTPCode:  http.StatusOK,
				ExpectedErrorCode: "",
				MockExpectFunc: func() {
					s.mockJobSrv.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(s.mockJob, nil)
				},
			},
			pathParams: []gin.Param{{Key: "JobID", Value: "4ER"}},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test admin get job bind error",
				ExpectedHTTPCode:  http.StatusBadRequest,
				ExpectedErrorCode: api.InvalidArgumentErrorCode,
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test admin get job invalid job id",
				ExpectedHTTPCode:  http.StatusBadRequest,
				ExpectedErrorCode: api.InvalidJobID,
			},
			pathParams: []gin.Param{{Key: "JobID", Value: "O"}}, // O is invalid
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test admin get job id not found",
				ExpectedHTTPCode:  http.StatusNotFound,
				ExpectedErrorCode: api.JobIDNotFound,
				MockExpectFunc: func() {
					s.mockJobSrv.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, common.ErrJobIDNotFound)
				},
			},
			pathParams: []gin.Param{{Key: "JobID", Value: "NtFd"}},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test admin get job internal error",
				ExpectedHTTPCode:  http.StatusInternalServerError,
				ExpectedErrorCode: api.InternalServerErrorCode,
				MockExpectFunc: func() {
					s.mockJobSrv.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, common.ErrInternalServer)
				},
			},
			pathParams: []gin.Param{{Key: "JobID", Value: "4ER"}},
		},
	}

	for _, tc := range testCase {
		testJobAdminGet := tc.MakeTestFunc(s.reqFunc(tc.pathParams), s.doFunc, s.respFunc(&tc))
		testJobAdminGet(s.T(), &s.JobSuite)
	}
}
