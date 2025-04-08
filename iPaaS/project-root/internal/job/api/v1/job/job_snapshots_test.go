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
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/jobsnapshotlist"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/common/mock"
)

type JobSnapshotsSuite struct {
	JobSuite

	// data
	mockSnapshots map[string][]string
}

func TestJobSnapshots(t *testing.T) {
	suite.Run(t, new(JobSnapshotsSuite))
}

func (s *JobSnapshotsSuite) SetupTest() {
	s.JobSuite.SetupTest()
	s.mockSnapshots = map[string][]string{
		"4ER": {"1", "2", "3"},
	}
}

// TestCase structure
type jobSnapshotsTestCase struct {
	JobTestCase
	pathParams    []gin.Param
	setHeaderFunc func(ctx *gin.Context)
}

func (s *JobSnapshotsSuite) reqFunc(pathParams []gin.Param, setHeaderFunc func(ctx *gin.Context)) func() {
	return func() {
		// Set custom headers
		if setHeaderFunc != nil {
			setHeaderFunc(s.ctx)
		}

		mock.HTTPRequest(s.ctx, http.MethodGet, nil, pathParams, nil)
	}
}

func (s *JobSnapshotsSuite) doFunc() {
	s.handler.Snapshots(s.ctx)
}

func (s *JobSnapshotsSuite) respFunc(tc *jobSnapshotsTestCase) func() {
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

func (s *JobSnapshotsSuite) TestSnapshots() {
	testCases := []jobSnapshotsTestCase{
		{
			JobTestCase: JobTestCase{
				Name:              "test get snapshots success",
				ExpectedHTTPCode:  http.StatusOK,
				ExpectedErrorCode: "",
				MockExpectFunc: func() {
					s.mockJobSrv.EXPECT().ListJobSnapshot(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(s.mockSnapshots, nil)
				},
			},
			pathParams: []gin.Param{{Key: "JobID", Value: "4ER"}},
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "h9z")
			},
		},
	}

	for _, tc := range testCases {
		testJobSnapshots := tc.MakeTestFunc(s.reqFunc(tc.pathParams, tc.setHeaderFunc), s.doFunc, s.respFunc(&tc))
		testJobSnapshots(s.T(), &s.JobSuite)
	}
}

func (s *JobSnapshotsSuite) TestSnapshotsValidate() {
	testCases := []jobSnapshotsTestCase{
		{
			JobTestCase: JobTestCase{
				Name:              "test get snapshots bind error",
				ExpectedHTTPCode:  http.StatusBadRequest,
				ExpectedErrorCode: api.InvalidJobID,
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test get snapshots invalid user id",
				ExpectedHTTPCode:  http.StatusBadRequest,
				ExpectedErrorCode: api.InvalidUserID,
			},
			pathParams: []gin.Param{{Key: "JobID", Value: "4ER"}},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test get snapshots invalid job id",
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
		testJobSnapshots := tc.MakeTestFunc(s.reqFunc(tc.pathParams, tc.setHeaderFunc), s.doFunc, s.respFunc(&tc))
		testJobSnapshots(s.T(), &s.JobSuite)
	}
}

func (s *JobSnapshotsSuite) TestSnapshotsJudge() {
	testCases := []jobSnapshotsTestCase{
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
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "h9z")
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test get snapshots access denied",
				ExpectedHTTPCode:  http.StatusForbidden,
				ExpectedErrorCode: api.JobAccessDenied,
				MockExpectFunc: func() {
					s.mockJobSrv.EXPECT().ListJobSnapshot(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, common.ErrJobAccessDenied)
				},
			},
			pathParams: []gin.Param{{Key: "JobID", Value: "4ER"}},
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "h9z")
			},
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
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "h9z")
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test get snapshots app not enable snapshot",
				ExpectedHTTPCode:  http.StatusNotFound,
				ExpectedErrorCode: api.JobSnapshotNotFound,
				MockExpectFunc: func() {
					s.mockJobSrv.EXPECT().ListJobSnapshot(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, common.ErrJobSnapshotNotFound)
				},
			},
			pathParams: []gin.Param{{Key: "JobID", Value: "4ER"}},
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "h9z")
			},
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
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "h9z")
			},
		},
	}

	for _, tc := range testCases {
		testJobSnapshots := tc.MakeTestFunc(s.reqFunc(tc.pathParams, tc.setHeaderFunc), s.doFunc, s.respFunc(&tc))
		testJobSnapshots(s.T(), &s.JobSuite)
	}
}
