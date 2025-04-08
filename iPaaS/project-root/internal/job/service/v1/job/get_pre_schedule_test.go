package job

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao/models"
)

type GetPreScheduleSuite struct {
	jobServiceSuite

	// data
}

func TestGetPreSchedule(t *testing.T) {
	suite.Run(t, new(GetPreScheduleSuite))
}

type GetPreScheduleTestCase struct {
	name           string
	id             string
	expectedError  error
	mockExpectFunc func()
	setReq         func()
}

func (tc *GetPreScheduleTestCase) Run(s *GetPreScheduleSuite) {
	s.Run(tc.name, func() {
		if tc.mockExpectFunc != nil {
			tc.mockExpectFunc()
		}

		if tc.setReq != nil {
			tc.setReq()
		}

		// do
		_, _, err := s.jobSrv.GetPreSchedule(s.ctx, tc.id)

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

func (s *GetPreScheduleSuite) TestGetPreSchedule() {
	tastCases := []GetPreScheduleTestCase{
		{
			name: "normal",
			id:   "4ER",
			mockExpectFunc: func() {
				s.mockJobDao.EXPECT().GetPreSchedule(s.ctx, snowflake.ID(12345)).Return(&models.PreSchedule{
					ID: snowflake.ID(12345),
				}, true, nil)
			},
			expectedError: nil,
		},
	}

	for _, tc := range tastCases {
		tc.Run(s)
	}
}
