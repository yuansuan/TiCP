package job

import (
	"encoding/json"
	"github.com/davecgh/go-spew/spew"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/common/project-root-api/common"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/admin/joblistfiltered"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/common/mock"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao/models"
)

type AdminJobListFilteredSuite struct {
	JobSuite

	mockJobAdminList []*models.Job
}

func TestAdminJobListFiltered(t *testing.T) {
	suite.Run(t, new(AdminJobListFilteredSuite))
}

func (s *AdminJobListFilteredSuite) SetupSuite() {
	s.JobSuite.SetupSuite()
	s.mockConfig()
}

func (s *AdminJobListFilteredSuite) SetupTest() {
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

type adminJobListFilteredTestCase struct {
	JobTestCase
	queryParams url.Values
	respFunc    func()
}

func (s *AdminJobListFilteredSuite) reqFunc(queryParams url.Values) func() {
	return func() {
		mock.HTTPRequest(s.ctx, http.MethodGet, nil, nil, queryParams)
	}
}

func (s *AdminJobListFilteredSuite) doFunc() {
	s.handler.AdminJobListFiltered(s.ctx)
}

func (s *AdminJobListFilteredSuite) respFunc(tc *adminJobListFilteredTestCase) func() {
	return func() {
		got := s.w.Body.Bytes()
		var resp joblistfiltered.Response

		err := json.Unmarshal(got, &resp)
		if !s.NoError(err) {
			return
		}
		s.T().Log(spew.Sdump(resp))

		tc.ErrorCode = resp.ErrorCode
		tc.ErrorMsg = resp.ErrorMsg

		if resp.Data != nil {
			s.Equal(int64(len(s.mockJobAdminList)), resp.Data.Total, "Expected total to be %v, got %v", len(s.mockJobAdminList), resp.Data.Total)
		}
	}
}

func (s *AdminJobListFilteredSuite) TestAdminJobListFiltered() {
	testCases := []adminJobListFilteredTestCase{
		{
			JobTestCase: JobTestCase{
				Name:              "test admin job list filtered success",
				ExpectedHTTPCode:  http.StatusOK,
				ExpectedErrorCode: "",
				MockExpectFunc: func() {
					s.mockJobSrv.EXPECT().ListFiltered(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(int64(2), s.mockJobAdminList, nil)
				},
			},
			queryParams: url.Values{
				"PageOffset": []string{"1"},
				"PageSize":   []string{"10"},
				"UserID":     []string{"h9z"},
				"AppID":      []string{"12345"},
				"JobID":      []string{"123456"},
			},
			respFunc: func() {
				got := s.w.Body.Bytes()
				var resp joblistfiltered.Data

				err := json.Unmarshal(got, &resp)
				if !s.NoError(err) {
					return
				}
				s.T().Log(spew.Sdump(resp))

				// 验证 total 字段
				s.Equal(int64(len(s.mockJobAdminList)), resp.Total, "Expected total to be %v, got %v", len(s.mockJobAdminList), resp.Total)
			},
		},

		{
			JobTestCase: JobTestCase{
				Name:              "test admin job list filtered invalid params",
				ExpectedHTTPCode:  http.StatusBadRequest,
				ExpectedErrorCode: common.InvalidArgumentErrorCode,
			},
			queryParams: url.Values{
				"PageOffset": []string{"AAA"},
				"PageSize":   []string{"1"},
			},
			respFunc: func() {
				got := s.w.Body.Bytes()
				var resp joblistfiltered.Data

				err := json.Unmarshal(got, &resp)
				if !s.NoError(err) {
					return
				}
				s.T().Log(spew.Sdump(resp))

				// 验证 total 字段
				s.Equal(int64(len(s.mockJobAdminList)), resp.Total, "Expected total to be %v, got %v", len(s.mockJobAdminList), resp.Total)
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test admin job list filtered internal error",
				ExpectedHTTPCode:  http.StatusInternalServerError,
				ExpectedErrorCode: common.InternalServerErrorCode,
				MockExpectFunc: func() {
					s.mockJobSrv.EXPECT().ListFiltered(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(int64(0), nil, errors.New("some error"))
				},
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test admin list job invalid user id",
				ExpectedHTTPCode:  http.StatusBadRequest,
				ExpectedErrorCode: common.InvalidUserID,
			},
			queryParams: url.Values{
				"UserID": []string{"O"},
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test admin list job invalid zone",
				ExpectedHTTPCode:  http.StatusBadRequest,
				ExpectedErrorCode: common.InvalidArgumentZone,
			},
			queryParams: url.Values{
				"Zone": []string{"az-zhigu"},
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test admin list job invalid page offset",
				ExpectedHTTPCode:  http.StatusBadRequest,
				ExpectedErrorCode: common.InvalidPageOffset,
			},
			queryParams: url.Values{
				"PageOffset": []string{"-1"},
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test admin list job invalid page size",
				ExpectedHTTPCode:  http.StatusBadRequest,
				ExpectedErrorCode: common.InvalidPageSize,
			},
			queryParams: url.Values{
				"PageSize": []string{"0"},
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test admin job list filtered invalid JobID",
				ExpectedHTTPCode:  http.StatusBadRequest,
				ExpectedErrorCode: common.InvalidJobID,
			},
			queryParams: url.Values{
				"JobID": []string{"Oo0ccxs"},
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test admin job list filtered invalid AccountID",
				ExpectedHTTPCode:  http.StatusBadRequest,
				ExpectedErrorCode: common.InternalErrorInvalidAccountId,
			},
			queryParams: url.Values{
				"AccountID": []string{"Oo0ccxs"},
			},
		},
		{
			JobTestCase: JobTestCase{
				Name:              "test admin job list filtered invalid AppID",
				ExpectedHTTPCode:  http.StatusBadRequest,
				ExpectedErrorCode: common.InvalidAppID,
			},
			queryParams: url.Values{
				"AppID": []string{"Oo0ccxs"},
			},
		},
	}

	for _, tc := range testCases {
		testJobAdminListFiltered := tc.MakeTestFunc(s.reqFunc(tc.queryParams), s.doFunc, s.respFunc(&tc))
		testJobAdminListFiltered(s.T(), &s.JobSuite)
	}
}
