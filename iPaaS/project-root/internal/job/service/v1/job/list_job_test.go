package job

import (
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/joblist"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao/models"
)

type ListSuite struct {
	jobServiceSuite

	// data
}

func TestList(t *testing.T) {
	suite.Run(t, new(ListSuite))
}

type JobListTestCase struct {
	name           string
	mockExpectFunc func()
}

func (tc *JobListTestCase) Run(s *ListSuite) {
	s.Run(tc.name, func() {
		if tc.mockExpectFunc != nil {
			tc.mockExpectFunc()
		}

		// do
		ps, po := int64(0), int64(0)
		_, _, err := s.jobSrv.List(s.ctx, &joblist.Request{
			PageSize:   &ps,
			PageOffset: &po,
		}, snowflake.ID(10086), snowflake.Zero(), false, false)

		s.NoError(err)
	})
}

func (s *ListSuite) TestList() {
	tastCases := []JobListTestCase{
		{
			name: "normal",
			mockExpectFunc: func() {
				s.mockJobDao.EXPECT().ListJobs(gomock.Any(), 0, 0, snowflake.ID(10086), snowflake.Zero(), "", "", gomock.Any(), gomock.Any()).Return(int64(0), []*models.Job{
					{
						ID: 12345,
					},
				}, nil)
			},
		},
	}

	for _, tc := range tastCases {
		tc.Run(s)
	}
}
