package job

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/admin/joblist"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/admin/joblistfiltered"
	list "github.com/yuansuan/ticp/common/project-root-api/job/v1/joblist"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao/models"
	"testing"
)

type ListFilteredSuite struct {
	suite.Suite
	ctrl       *gomock.Controller
	mockJobDao *dao.MockJobDao
	jobSrv     *jobService
}

func TestListFilteredSuite(t *testing.T) {
	suite.Run(t, new(ListFilteredSuite))
}

func (s *ListFilteredSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.mockJobDao = dao.NewMockJobDao(s.ctrl)
	s.jobSrv = &jobService{
		jobdao: s.mockJobDao,
	}
}

func (s *ListFilteredSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *ListFilteredSuite) TestListFiltered() {
	var testCases []struct {
		name           string
		mockExpectFunc func()
		expectedCount  int64
		expectedJobs   []*models.Job
		expectedError  error
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// Setup
			tc.mockExpectFunc()

			// Create request
			req := &joblistfiltered.Request{
				Request: joblist.Request{
					Request: list.Request{
						PageOffset: new(int64),
						PageSize:   new(int64),
						JobState:   "active",
						Zone:       "zone-jinan",
					},
					AppID:  "app-123",
					UserID: "user-4TiSsZonTa3",
				},
				JobID:     "job-123",
				AccountID: "account-123",
				Name:      "job-name",
			}
			*req.PageSize = 10
			*req.PageOffset = 0

			// Print the request for debugging
			s.T().Logf("Request: %+v\n", req)

			// Execute
			count, jobs, err := s.jobSrv.ListFiltered(context.Background(), req, snowflake.ID(10086), snowflake.Zero())

			// Assertions
			if tc.expectedError != nil {
				s.Error(err)
				s.Equal(tc.expectedError.Error(), err.Error())
			} else {
				s.NoError(err)
			}
			s.Equal(tc.expectedCount, count)
			s.ElementsMatch(tc.expectedJobs, jobs)
		})
	}
}
