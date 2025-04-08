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
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	api "github.com/yuansuan/ticp/common/project-root-api/common"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/joblist"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao/models"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/common/mock"
)

// JobListSuite
type JobListSuite struct {
	JobSuite

	// data
	mockJobList []*models.Job
}

func TestJobList(t *testing.T) {
	suite.Run(t, new(JobListSuite))
}

func (s *JobListSuite) SetupSuite() {
	s.JobSuite.SetupSuite()
	s.mockConfig()
}

func (s *JobListSuite) SetupTest() {
	s.JobSuite.SetupTest()
	p := models.AdminParams{}
	ps, _ := json.Marshal(p)
	s.mockJobList = []*models.Job{
		{
			ID:     snowflake.ID(12345),
			UserID: snowflake.ID(54321),
			Params: string(ps),
		},
		{
			ID:     snowflake.ID(56789),
			UserID: snowflake.ID(98765),
			Params: string(ps),
		},
	}
}

type jobListTestCase struct {
	JobTestCase
	queryParams   url.Values
	setHeaderFunc func(ctx *gin.Context)
}

func (s *JobListSuite) reqFunc(queryParams url.Values, setHeaderFunc func(ctx *gin.Context)) func() {
	return func() {
		// Set custom headers
		if setHeaderFunc != nil {
			setHeaderFunc(s.ctx)
		}

		mock.HTTPRequest(s.ctx, http.MethodGet, nil, nil, queryParams)
	}
}

func (s *JobListSuite) doFunc() {
	s.handler.List(s.ctx)
}

func (s *JobListSuite) respFunc(tc *jobListTestCase) func() {
	return func() {
		got := s.w.Body.Bytes()
		var resp joblist.Response

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

func (s *JobListSuite) TestList() {
	testCase := []jobListTestCase{
		{
			JobTestCase: JobTestCase{
				Name:              "test list job success",
				ExpectedHTTPCode:  http.StatusOK,
				ExpectedErrorCode: "",
				MockExpectFunc: func() {
					s.mockJobSrv.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(int64(2), s.mockJobList, nil)
				},
			},
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "h9z")
			},
		},
	}

	for _, tc := range testCase {
		testJobList := tc.MakeTestFunc(s.reqFunc(tc.queryParams, tc.setHeaderFunc), s.doFunc, s.respFunc(&tc))
		testJobList(s.T(), &s.JobSuite)
	}
}

func (s *JobListSuite) TestListValidate() {
	testCase := []jobListTestCase{
		{
			JobTestCase: JobTestCase{
				Name:              "test list job invalid params",
				ExpectedHTTPCode:  http.StatusBadRequest,
				ExpectedErrorCode: api.InvalidArgumentErrorCode,
			},
			queryParams: url.Values{
				"PageOffset": []string{"AAA"},
				"PageSize":   []string{"1"},
			},
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "h9z")
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test list job invalid user id",
				ExpectedHTTPCode:  http.StatusBadRequest,
				ExpectedErrorCode: api.InvalidUserID,
			},
			queryParams: url.Values{},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test list job invalid page size",
				ExpectedHTTPCode:  http.StatusBadRequest,
				ExpectedErrorCode: api.InvalidPageSize,
			},
			queryParams: url.Values{
				"PageSize": []string{"0"},
			},
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "h9z")
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test list job invalid page offset",
				ExpectedHTTPCode:  http.StatusBadRequest,
				ExpectedErrorCode: api.InvalidPageOffset,
			},
			queryParams: url.Values{
				"PageOffset": []string{"-1"},
			},
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "h9z")
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test list job invalid zone",
				ExpectedHTTPCode:  http.StatusBadRequest,
				ExpectedErrorCode: api.InvalidArgumentZone,
			},
			queryParams: url.Values{
				"Zone": []string{"az-zhigu"},
			},
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "h9z")
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test list job invalid job state",
				ExpectedHTTPCode:  http.StatusBadRequest,
				ExpectedErrorCode: api.InvalidArgumentJobState,
			},
			queryParams: url.Values{
				"JobState": []string{"AAAA"},
			},
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "h9z")
			},
		},
	}

	for _, tc := range testCase {
		testJobList := tc.MakeTestFunc(s.reqFunc(tc.queryParams, tc.setHeaderFunc), s.doFunc, s.respFunc(&tc))
		testJobList(s.T(), &s.JobSuite)
	}
}

func (s *JobListSuite) TestListJudge() {
	testCase := []jobListTestCase{
		{
			JobTestCase: JobTestCase{
				Name:              "test list job internal error",
				ExpectedHTTPCode:  http.StatusInternalServerError,
				ExpectedErrorCode: api.InternalServerErrorCode,
				MockExpectFunc: func() {
					s.mockJobSrv.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(int64(0), nil, common.ErrInternalServer)
				},
			},
			setHeaderFunc: func(ctx *gin.Context) {
				ctx.Request.Header.Set("x-ys-user-id", "h9z")
			},
		},
	}

	for _, tc := range testCase {
		testJobList := tc.MakeTestFunc(s.reqFunc(tc.queryParams, tc.setHeaderFunc), s.doFunc, s.respFunc(&tc))
		testJobList(s.T(), &s.JobSuite)
	}
}
