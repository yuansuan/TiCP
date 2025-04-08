package job

import (
	"fmt"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao/models"
)

type BatchGetSuite struct {
	jobServiceSuite
}

func TestBatchGet(t *testing.T) {
	suite.Run(t, new(BatchGetSuite))
}

type JobBatchGetTestCase struct {
	name           string
	ids            []string
	mockExpectFunc func()
	expectedError  error
}

func (tc *JobBatchGetTestCase) Run(s *BatchGetSuite) {
	s.Run(tc.name, func() {
		if tc.mockExpectFunc != nil {
			tc.mockExpectFunc()
		}

		// do
		_, err := s.jobSrv.BatchGet(s.ctx, tc.ids, 0)

		// assert
		if tc.expectedError != nil {
			if s.Error(err) {
				s.ErrorContains(err, tc.expectedError.Error())
			}
			return
		}

		s.NoError(err)
	})
}

func (s *BatchGetSuite) TestBatchGet() {
	tastCases := []JobBatchGetTestCase{
		{
			name: "normal",
			ids: []string{
				"4ER",
				"h9z",
			},
			mockExpectFunc: func() {
				s.mockJobDao.EXPECT().BatchGet(gomock.Any(), []snowflake.ID{12345, 54321}, snowflake.Zero(), false, gomock.Any()).Return([]*models.Job{
					{
						ID: snowflake.ID(12345),
					},
					{
						ID: snowflake.ID(54321),
					},
				}, nil).Times(1)
			},
			expectedError: nil,
		},
		{
			name: "with error",
			ids:  []string{},
			mockExpectFunc: func() {
				s.mockJobDao.EXPECT().BatchGet(gomock.Any(), []snowflake.ID{}, snowflake.Zero(), false, gomock.Any()).Return(nil, fmt.Errorf("error")).Times(1)
			},
			expectedError: fmt.Errorf("error"),
		},
	}

	for _, tc := range tastCases {
		tc.Run(s)
	}
}
