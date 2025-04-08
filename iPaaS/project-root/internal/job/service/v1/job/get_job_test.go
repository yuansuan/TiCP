package job

import (
	"fmt"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/jobget"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao/models"
)

type GetSuite struct {
	jobServiceSuite

	// data
	allow allowFunc
}

func TestGet(t *testing.T) {
	suite.Run(t, new(GetSuite))
}

type JobGetTestCase struct {
	name           string
	id             string
	expectedError  error
	mockExpectFunc func()
	setReq         func()
}

func (tc *JobGetTestCase) Run(s *GetSuite) {
	s.Run(tc.name, func() {
		if tc.mockExpectFunc != nil {
			tc.mockExpectFunc()
		}

		if tc.setReq != nil {
			tc.setReq()
		}

		// do
		_, err := s.jobSrv.Get(s.ctx, &jobget.Request{
			JobID: tc.id,
		}, snowflake.ID(10086), s.allow, false)

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

func (s *GetSuite) TestGet() {
	tastCases := []JobGetTestCase{
		{
			name: "normal",
			id:   "4ER",
			mockExpectFunc: func() {
				s.mockJobDao.EXPECT().Get(s.ctx, snowflake.ID(12345), false, gomock.Any()).Return(&models.Job{
					ID:     snowflake.ID(12345),
					UserID: snowflake.ID(10086),
				}, nil)
			},
			expectedError: nil,
			setReq: func() {
				s.allow = func(userID, jobUserID string) bool {
					return true
				}
			},
		},
		{
			name: "with error",
			id:   "4ER",
			mockExpectFunc: func() {
				s.mockJobDao.EXPECT().Get(s.ctx, snowflake.ID(12345), false, gomock.Any()).Return(nil, fmt.Errorf("error"))
			},
			expectedError: fmt.Errorf("error"),
			setReq: func() {
				s.allow = func(userID, jobUserID string) bool {
					return true
				}
			},
		},
		{
			name: "no permission",
			id:   "4ER",
			mockExpectFunc: func() {
				s.mockJobDao.EXPECT().Get(s.ctx, snowflake.ID(12345), false, gomock.Any()).Return(&models.Job{
					ID:     snowflake.ID(12345),
					UserID: snowflake.ID(10086),
				}, nil)
			},
			expectedError: common.ErrJobAccessDenied,
			setReq: func() {
				s.allow = func(userID, jobUserID string) bool {
					return false
				}
			},
		},
	}

	for _, tc := range tastCases {
		tc.Run(s)
	}
}
