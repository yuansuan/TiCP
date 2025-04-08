package job

import (
	"context"
	"fmt"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/jobresidual"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao/models"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/common/with"
	"xorm.io/xorm"
)

type ResidualSuite struct {
	jobServiceSuite

	// data
	allow allowFunc
}

func TestResidual(t *testing.T) {
	suite.Run(t, new(ResidualSuite))
}

type JobResidualTestCase struct {
	name           string
	id             string
	expectedError  error
	mockExpectFunc func()
	setReq         func()
}

func (tc *JobResidualTestCase) Run(s *ResidualSuite) {
	s.Run(tc.name, func() {
		if tc.mockExpectFunc != nil {
			tc.mockExpectFunc()
		}

		if tc.setReq != nil {
			tc.setReq()
		}

		// do
		_, err := s.jobSrv.GetResidual(s.ctx, &jobresidual.Request{
			JobID: tc.id,
		}, snowflake.ID(10086), s.allow)

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

func (s *ResidualSuite) TestGetResidual() {
	mockTransaction := func() {
		s.mockJobDao.EXPECT().Transaction(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, action func(context.Context) error) error {
				_, err := s.engine.Transaction(func(db *xorm.Session) (interface{}, error) {
					return nil, action(with.KeepSession(ctx, db))
				})
				return err
			})
	}

	tastCases := []JobResidualTestCase{
		{
			name: "normal",
			id:   "4ER",
			mockExpectFunc: func() {
				mockTransaction()

				s.sqlmock.ExpectBegin()
				s.mockJobDao.EXPECT().Get(gomock.Any(), snowflake.ID(12345), false, gomock.Any()).Return(&models.Job{
					ID:     snowflake.ID(12345),
					UserID: snowflake.ID(10086),
				}, nil).Times(1)

				s.mockResidualDao.EXPECT().GetJobResidual(gomock.Any(), snowflake.ID(12345)).Return(&models.Residual{
					JobID:    snowflake.ID(12345),
					Finished: true,
				}, nil).Times(1)

				s.sqlmock.ExpectCommit()
			},
			setReq: func() {
				s.allow = func(userID, jobUserID string) bool {
					return true
				}
			},
		},
		{
			name: "get job error",
			id:   "ERR",
			mockExpectFunc: func() {
				mockTransaction()

				s.sqlmock.ExpectBegin()
				s.mockJobDao.EXPECT().Get(gomock.Any(), snowflake.ID(130723), false, gomock.Any()).Return(nil, fmt.Errorf("get job error")).Times(1)
				s.sqlmock.ExpectRollback()
			},
			expectedError: fmt.Errorf("get job error"),
			setReq: func() {
				s.allow = func(userID, jobUserID string) bool {
					return true
				}
			},
		},
		{
			name: "no permission",
			id:   "oTher",
			mockExpectFunc: func() {
				mockTransaction()

				s.sqlmock.ExpectBegin()
				s.mockJobDao.EXPECT().Get(gomock.Any(), snowflake.ID(258968227), false, gomock.Any()).Return(&models.Job{
					ID:     snowflake.ID(258968227),
					UserID: snowflake.ID(258968227),
				}, nil).Times(1)
				s.sqlmock.ExpectRollback()
			},
			expectedError: common.ErrJobAccessDenied,
			setReq: func() {
				s.allow = func(userID, jobUserID string) bool {
					return false
				}
			},
		},
		{
			name: "get job monitor chart error",
			id:   "ERR",
			mockExpectFunc: func() {
				mockTransaction()

				s.sqlmock.ExpectBegin()
				s.mockJobDao.EXPECT().Get(gomock.Any(), snowflake.ID(130723), false, gomock.Any()).Return(&models.Job{
					ID:     snowflake.ID(130723),
					UserID: snowflake.ID(10086),
				}, nil).Times(1)
				s.mockResidualDao.EXPECT().GetJobResidual(gomock.Any(), snowflake.ID(130723)).Return(nil, fmt.Errorf("get job monitor chart error")).Times(1)
				s.sqlmock.ExpectRollback()
			},
			expectedError: fmt.Errorf("get job monitor chart error"),
			setReq: func() {
				s.allow = func(userID, jobUserID string) bool {
					return true
				}
			},
		},
	}

	for _, tc := range tastCases {
		tc.Run(s)
	}
}
