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
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/admin/jobsnapshotlist"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/common/mock"
)

type JobAdminSnapshotsSuite struct {
	JobSuite

	// data
	mockSnapshots map[string][]string
}

func TestJobAdminSnapshots(t *testing.T) {
	suite.Run(t, new(JobAdminSnapshotsSuite))
}

func (s *JobAdminSnapshotsSuite) SetupTest() {
	s.JobSuite.SetupTest()
	s.mockSnapshots = map[string][]string{
		"4ER": {"1", "2", "3"},
	}
}

// TestCase structure
type jobAdminSnapshotsTestCase struct {
	JobTestCase
	pathParams []gin.Param
}

func (s *JobAdminSnapshotsSuite) reqFunc(pathParams []gin.Param) func() {
	return func() {
		mock.HTTPRequest(s.ctx, http.MethodGet, nil, pathParams, nil)
	}
}

func (s *JobAdminSnapshotsSuite) doFunc() {
	s.handler.AdminSnapshots(s.ctx)
}

func (s *JobAdminSnapshotsSuite) respFunc(tc *jobAdminSnapshotsTestCase) func() {
	return func() {
		got := s.w.Body.Bytes()
		var resp jobsnapshotlist.Response

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

func (s *JobAdminSnapshotsSuite) TestAdminSnapshots() {
	testCases := []jobAdminSnapshotsTestCase{
		{
			JobTestCase: JobTestCase{
				Name:              "test admin get snapshots success",
				ExpectedHTTPCode:  http.StatusOK,
				ExpectedErrorCode: "",
				MockExpectFunc: func() {
					s.mockJobSrv.EXPECT().ListJobSnapshot(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(s.mockSnapshots, nil)
				},
			},
			pathParams: []gin.Param{{Key: "JobID", Value: "4ER"}},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test get snapshots bind error",
				ExpectedHTTPCode:  http.StatusBadRequest,
				ExpectedErrorCode: api.InvalidJobID,
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test get snapshots invalid job id",
				ExpectedHTTPCode:  http.StatusBadRequest,
				ExpectedErrorCode: api.InvalidJobID,
			},
			pathParams: []gin.Param{{Key: "JobID", Value: "O"}}, // O is invalid
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test get snapshots job id not found",
				ExpectedHTTPCode:  http.StatusNotFound,
				ExpectedErrorCode: api.JobIDNotFound,
				MockExpectFunc: func() {
					s.mockJobSrv.EXPECT().ListJobSnapshot(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, common.ErrJobIDNotFound)
				},
			},
			pathParams: []gin.Param{{Key: "JobID", Value: "NtFd"}},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test get snapshots app not found",
				ExpectedHTTPCode:  http.StatusNotFound,
				ExpectedErrorCode: api.AppIDNotFoundErrorCode,
				MockExpectFunc: func() {
					s.mockJobSrv.EXPECT().ListJobSnapshot(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, common.ErrAppIDNotFound)
				},
			},
			pathParams: []gin.Param{{Key: "JobID", Value: "4ER"}},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test get snapshots internal error",
				ExpectedHTTPCode:  http.StatusInternalServerError,
				ExpectedErrorCode: api.InternalServerErrorCode,
				MockExpectFunc: func() {
					s.mockJobSrv.EXPECT().ListJobSnapshot(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, common.ErrInternalServer)
				},
			},
			pathParams: []gin.Param{{Key: "JobID", Value: "4ER"}},
		},
	}

	for _, tc := range testCases {
		testJobAdminSnapshots := tc.MakeTestFunc(s.reqFunc(tc.pathParams), s.doFunc, s.respFunc(&tc))
		testJobAdminSnapshots(s.T(), &s.JobSuite)
	}
}
