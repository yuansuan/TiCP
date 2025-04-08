package job

import (
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	api "github.com/yuansuan/ticp/common/project-root-api/common"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/admin/joblist"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao/models"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/common/mock"
)

// JobAdminListSuite
type JobAdminListSuite struct {
	JobSuite

	// data
	mockJobAdminList []*models.Job
}

func TestJobAdminList(t *testing.T) {
	suite.Run(t, new(JobAdminListSuite))
}

func (s *JobAdminListSuite) SetupSuite() {
	s.JobSuite.SetupSuite()
	s.mockConfig()
}

func (s *JobAdminListSuite) SetupTest() {
	s.JobSuite.SetupTest()
	s.mockJobAdminList = []*models.Job{
		{
			ID:     snowflake.ID(12345),
			UserID: snowflake.ID(54321),
		},
		{
			ID:     snowflake.ID(56789),
			UserID: snowflake.ID(98765),
		},
	}
}

type jobAdminListTestCase struct {
	JobTestCase
	queryParams url.Values
}

func (s *JobAdminListSuite) reqFunc(queryParams url.Values) func() {
	return func() {
		mock.HTTPRequest(s.ctx, http.MethodGet, nil, nil, queryParams)
	}
}

func (s *JobAdminListSuite) doFunc() {
	s.handler.AdminList(s.ctx)
}

func (s *JobAdminListSuite) respFunc(tc *jobAdminListTestCase) func() {
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

func (s *JobAdminListSuite) TestAdminList() {
	testCase := []jobAdminListTestCase{
		{
			JobTestCase: JobTestCase{
				Name:              "test admin list job success",
				ExpectedHTTPCode:  http.StatusOK,
				ExpectedErrorCode: "",
				MockExpectFunc: func() {
					s.mockJobSrv.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(int64(2), s.mockJobAdminList, nil)
				},
			},
			queryParams: url.Values{
				"UserID": []string{"h9z"},
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test admin list job invalid params",
				ExpectedHTTPCode:  http.StatusBadRequest,
				ExpectedErrorCode: api.InvalidArgumentErrorCode,
			},
			queryParams: url.Values{
				"PageOffset": []string{"AAA"},
				"PageSize":   []string{"1"},
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test admin list job invalid page size",
				ExpectedHTTPCode:  http.StatusBadRequest,
				ExpectedErrorCode: api.InvalidPageSize,
			},
			queryParams: url.Values{
				"PageSize": []string{"0"},
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test admin list job invalid page offset",
				ExpectedHTTPCode:  http.StatusBadRequest,
				ExpectedErrorCode: api.InvalidPageOffset,
			},
			queryParams: url.Values{
				"PageOffset": []string{"-1"},
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test admin list job invalid zone",
				ExpectedHTTPCode:  http.StatusBadRequest,
				ExpectedErrorCode: api.InvalidArgumentZone,
			},
			queryParams: url.Values{
				"Zone": []string{"az-zhigu"},
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test admin list job invalid job state",
				ExpectedHTTPCode:  http.StatusBadRequest,
				ExpectedErrorCode: api.InvalidArgumentJobState,
			},
			queryParams: url.Values{
				"JobState": []string{"AAAA"},
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test admin list job invalid user id",
				ExpectedHTTPCode:  http.StatusBadRequest,
				ExpectedErrorCode: api.InvalidUserID,
			},
			queryParams: url.Values{
				"UserID": []string{"O"},
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test admin list job internal error",
				ExpectedHTTPCode:  http.StatusInternalServerError,
				ExpectedErrorCode: api.InternalServerErrorCode,
				MockExpectFunc: func() {
					s.mockJobSrv.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(int64(0), nil, common.ErrInternalServer)
				},
			},
		},
	}

	for _, tc := range testCase {
		testJobAdminList := tc.MakeTestFunc(s.reqFunc(tc.queryParams), s.doFunc, s.respFunc(&tc))
		testJobAdminList(s.T(), &s.JobSuite)
	}
}
