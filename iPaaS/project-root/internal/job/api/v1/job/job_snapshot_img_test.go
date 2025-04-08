package job

import (
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	api "github.com/yuansuan/ticp/common/project-root-api/common"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/jobsnapshotget"

	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/common/mock"
)

type JobSnapshotImgSuite struct {
	JobSuite

	// data
	mockSnapshotImg string
}

func TestJobSnapshotImg(t *testing.T) {
	suite.Run(t, new(JobSnapshotImgSuite))
}

func (s *JobSnapshotImgSuite) SetupTest() {
	s.JobSuite.SetupTest()
	s.mockSnapshotImg = "mockSnapshotImg"
}

// TestCase structure
type jobSnapshotImgTestCase struct {
	JobTestCase
	pathParams    []gin.Param
	queryParams   url.Values
	setHeaderFunc func(ctx *gin.Context)
}

func (s *JobSnapshotImgSuite) reqFunc(pathParams []gin.Param, queryParams url.Values, setHeaderFunc func(ctx *gin.Context)) func() {
	return func() {
		// Set custom headers
		if setHeaderFunc != nil {
			setHeaderFunc(s.ctx)
		}

		mock.HTTPRequest(s.ctx, http.MethodGet, nil, pathParams, queryParams)
	}
}

func (s *JobSnapshotImgSuite) doFunc() {
	s.handler.SnapshotImg(s.ctx)
}

func (s *JobSnapshotImgSuite) respFunc(tc *jobSnapshotImgTestCase) func() {
	return func() {
		got := s.w.Body.Bytes()
		var resp jobsnapshotget.Response

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

func (s *JobSnapshotImgSuite) TestSnapshotImg() {
	testCases := []jobSnapshotImgTestCase{
		{
			JobTestCase: JobTestCase{
				Name:              "test get snapshot img success",
				ExpectedHTTPCode:  http.StatusOK,
				ExpectedErrorCode: "",
				MockExpectFunc: func() {
					s.mockJobSrv.EXPECT().GetJobSnapshot(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(s.mockSnapshotImg, nil)
				},
			},
			pathParams: []gin.Param{{Key: "JobID", Value: "4ER"}},
			queryParams: url.Values{
				"Path": []string{"1"},
			},
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "h9z")
			},
		},
	}

	for _, tc := range testCases {
		testJobSnapshotImg := tc.MakeTestFunc(s.reqFunc(tc.pathParams, tc.queryParams, tc.setHeaderFunc), s.doFunc, s.respFunc(&tc))
		testJobSnapshotImg(s.T(), &s.JobSuite)
	}
}

func (s *JobSnapshotImgSuite) TestSnapshotImgValidate() {
	testCases := []jobSnapshotImgTestCase{
		{
			JobTestCase: JobTestCase{
				Name:              "test get snapshot img empty job id",
				ExpectedHTTPCode:  http.StatusBadRequest,
				ExpectedErrorCode: api.InvalidJobID,
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test get snapshot img empty path",
				ExpectedHTTPCode:  http.StatusBadRequest,
				ExpectedErrorCode: api.InvalidPath,
			},
			pathParams: []gin.Param{{Key: "JobID", Value: "4ER"}},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test get snapshot img invalid path",
				ExpectedHTTPCode:  http.StatusBadRequest,
				ExpectedErrorCode: api.InvalidPath,
			},
			pathParams: []gin.Param{{Key: "JobID", Value: "4ER"}},
			queryParams: url.Values{
				"Path": []string{"../"},
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test get snapshot img invalid user id",
				ExpectedHTTPCode:  http.StatusBadRequest,
				ExpectedErrorCode: api.InvalidUserID,
			},
			pathParams: []gin.Param{{Key: "JobID", Value: "4ER"}},
			queryParams: url.Values{
				"Path": []string{"1"},
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test get snapshot img invalid job id",
				ExpectedHTTPCode:  http.StatusBadRequest,
				ExpectedErrorCode: api.InvalidJobID,
			},
			pathParams: []gin.Param{{Key: "JobID", Value: "O"}}, // O is invalid
			queryParams: url.Values{
				"Path": []string{"1"},
			},
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "h9z")
			},
		},
	}

	for _, tc := range testCases {
		testJobSnapshotImg := tc.MakeTestFunc(s.reqFunc(tc.pathParams, tc.queryParams, tc.setHeaderFunc), s.doFunc, s.respFunc(&tc))
		testJobSnapshotImg(s.T(), &s.JobSuite)
	}
}

func (s *JobSnapshotImgSuite) TestSnapshotImgJudge() {
	testCases := []jobSnapshotImgTestCase{
		{
			JobTestCase: JobTestCase{
				Name:              "test get snapshot img job id not found",
				ExpectedHTTPCode:  http.StatusNotFound,
				ExpectedErrorCode: api.JobIDNotFound,
				MockExpectFunc: func() {
					s.mockJobSrv.EXPECT().GetJobSnapshot(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("", common.ErrJobIDNotFound)
				},
			},
			pathParams: []gin.Param{{Key: "JobID", Value: "NtFd"}},
			queryParams: url.Values{
				"Path": []string{"1"},
			},
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "h9z")
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test get snapshot img access denied",
				ExpectedHTTPCode:  http.StatusForbidden,
				ExpectedErrorCode: api.JobAccessDenied,
				MockExpectFunc: func() {
					s.mockJobSrv.EXPECT().GetJobSnapshot(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("", common.ErrJobAccessDenied)
				},
			},
			pathParams: []gin.Param{{Key: "JobID", Value: "4ER"}},
			queryParams: url.Values{
				"Path": []string{"1"},
			},
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "h9z")
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test get snapshot img app not found",
				ExpectedHTTPCode:  http.StatusNotFound,
				ExpectedErrorCode: api.AppIDNotFoundErrorCode,
				MockExpectFunc: func() {
					s.mockJobSrv.EXPECT().GetJobSnapshot(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("", common.ErrAppIDNotFound)
				},
			},
			pathParams: []gin.Param{{Key: "JobID", Value: "4ER"}},
			queryParams: url.Values{
				"Path": []string{"1"},
			},
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "h9z")
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test get snapshot img app not enable snapshot",
				ExpectedHTTPCode:  http.StatusNotFound,
				ExpectedErrorCode: api.JobSnapshotNotFound,
				MockExpectFunc: func() {
					s.mockJobSrv.EXPECT().GetJobSnapshot(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("", common.ErrJobSnapshotNotFound)
				},
			},
			pathParams: []gin.Param{{Key: "JobID", Value: "4ER"}},
			queryParams: url.Values{
				"Path": []string{"1"},
			},
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "h9z")
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test get snapshot img path not found",
				ExpectedHTTPCode:  http.StatusNotFound,
				ExpectedErrorCode: api.PathNotFound,
				MockExpectFunc: func() {
					s.mockJobSrv.EXPECT().GetJobSnapshot(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("", common.ErrPathNotFound)
				},
			},
			pathParams: []gin.Param{{Key: "JobID", Value: "4ER"}},
			queryParams: url.Values{
				"Path": []string{"1"},
			},
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "h9z")
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test get snapshot img internal error",
				ExpectedHTTPCode:  http.StatusInternalServerError,
				ExpectedErrorCode: api.InternalServerErrorCode,
				MockExpectFunc: func() {
					s.mockJobSrv.EXPECT().GetJobSnapshot(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("", common.ErrInternalServer)
				},
			},
			pathParams: []gin.Param{{Key: "JobID", Value: "4ER"}},
			queryParams: url.Values{
				"Path": []string{"1"},
			},
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "h9z")
			},
		},
	}

	for _, tc := range testCases {
		testJobSnapshotImg := tc.MakeTestFunc(s.reqFunc(tc.pathParams, tc.queryParams, tc.setHeaderFunc), s.doFunc, s.respFunc(&tc))
		testJobSnapshotImg(s.T(), &s.JobSuite)
	}
}
