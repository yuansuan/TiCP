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
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/admin/jobsnapshotget"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/common/mock"
)

type JobAdminSnapshotImgSuite struct {
	JobSuite

	// data
	mockSnapshotImg string
}

func TestJobAdminSnapshotImg(t *testing.T) {
	suite.Run(t, new(JobAdminSnapshotImgSuite))
}

func (s *JobAdminSnapshotImgSuite) SetupTest() {
	s.JobSuite.SetupTest()
	s.mockSnapshotImg = "mockSnapshotImg"
}

// TestCase structure
type jobAdminSnapshotImgTestCase struct {
	JobTestCase
	pathParams  []gin.Param
	queryParams url.Values
}

func (s *JobAdminSnapshotImgSuite) reqFunc(pathParams []gin.Param, queryParams url.Values) func() {
	return func() {
		mock.HTTPRequest(s.ctx, http.MethodGet, nil, pathParams, queryParams)
	}
}

func (s *JobAdminSnapshotImgSuite) doFunc() {
	s.handler.AdminSnapshotImg(s.ctx)
}

func (s *JobAdminSnapshotImgSuite) respFunc(tc *jobAdminSnapshotImgTestCase) func() {
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

func (s *JobAdminSnapshotImgSuite) TestAdminSnapshotImg() {
	testCases := []jobAdminSnapshotImgTestCase{
		{
			JobTestCase: JobTestCase{
				Name:              "test admin get snapshot img success",
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
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test admin get snapshot img empty job id",
				ExpectedHTTPCode:  http.StatusBadRequest,
				ExpectedErrorCode: api.InvalidJobID,
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test admin get snapshot img invalid job id",
				ExpectedHTTPCode:  http.StatusBadRequest,
				ExpectedErrorCode: api.InvalidJobID,
			},
			pathParams: []gin.Param{{Key: "JobID", Value: "O"}}, // O is invalid
			queryParams: url.Values{
				"Path": []string{"1"},
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test admin get snapshot img job id not found",
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
		},
	}

	for _, tc := range testCases {
		testJobAdminSnapshotImg := tc.MakeTestFunc(s.reqFunc(tc.pathParams, tc.queryParams), s.doFunc, s.respFunc(&tc))
		testJobAdminSnapshotImg(s.T(), &s.JobSuite)
	}
}
