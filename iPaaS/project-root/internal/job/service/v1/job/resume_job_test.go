package job

import (
	"context"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	jobresume "github.com/yuansuan/ticp/common/project-root-api/job/v1/jobresume"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/consts"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao/models"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/common/with"
	"xorm.io/xorm"
)

type ResumeSuite struct {
	jobServiceSuite

	// data
	allow allowFunc
}

func TestResume(t *testing.T) {
	suite.Run(t, new(ResumeSuite))
}

type JobResumeTestCase struct {
	name           string
	id             string
	expectedError  error
	mockExpectFunc func()
	setReq         func()
}

func (tc *JobResumeTestCase) Run(s *ResumeSuite) {
	s.Run(tc.name, func() {
		if tc.mockExpectFunc != nil {
			tc.mockExpectFunc()
		}

		if tc.setReq != nil {
			tc.setReq()
		}

		// do
		err := s.jobSrv.Resume(s.ctx, &jobresume.Request{
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

func (s *ResumeSuite) TestResume() {
	mockTransaction := func() {
		s.mockJobDao.EXPECT().Transaction(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, action func(context.Context) error) error {
				_, err := s.engine.Transaction(func(db *xorm.Session) (interface{}, error) {
					return nil, action(with.KeepSession(ctx, db))
				})
				return err
			})
	}
	tastCases := []JobResumeTestCase{
		{
			name: "normal",
			id:   "4ER",
			mockExpectFunc: func() {
				mockTransaction()

				s.sqlmock.ExpectBegin()
				s.mockJobDao.EXPECT().Get(gomock.Any(), snowflake.ID(12345), true, gomock.Any()).Return(&models.Job{
					ID:       snowflake.ID(12345),
					UserID:   snowflake.ID(10086),
					State:    consts.SubStateInitiallySuspended.State,
					SubState: consts.SubStateInitiallySuspended.SubState,
				}, nil)
				s.sqlmock.ExpectExec("UPDATE `job`").
					WillReturnResult(sqlmock.NewResult(1, 1))

				s.sqlmock.ExpectCommit()
			},
			expectedError: nil,
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
				s.mockJobDao.EXPECT().Get(gomock.Any(), snowflake.ID(130723), true, gomock.Any()).Return(nil, fmt.Errorf("get job error")).Times(1)
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
				s.mockJobDao.EXPECT().Get(gomock.Any(), snowflake.ID(258968227), true, gomock.Any()).Return(&models.Job{
					ID:       snowflake.ID(258968227),
					UserID:   snowflake.ID(258968227),
					State:    consts.SubStateInitiallySuspended.State,
					SubState: consts.SubStateInitiallySuspended.SubState,
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
			name: "job state error",
			id:   "jobState",
			mockExpectFunc: func() {
				mockTransaction()

				s.sqlmock.ExpectBegin()
				s.mockJobDao.EXPECT().Get(gomock.Any(), snowflake.ID(40588360944263), true, gomock.Any()).Return(&models.Job{
					ID:       snowflake.ID(40588360944263),
					UserID:   snowflake.ID(10086),
					State:    -1,
					SubState: -1,
				}, nil).Times(1)
				s.sqlmock.ExpectRollback()
			},
			expectedError: fmt.Errorf("invalid job state!"),
			setReq: func() {
				s.allow = func(userID, jobUserID string) bool {
					return true
				}
			},
		},
		{
			name: "job state not allow resume",
			id:   "jobState",
			mockExpectFunc: func() {
				mockTransaction()

				s.sqlmock.ExpectBegin()
				s.mockJobDao.EXPECT().Get(gomock.Any(), snowflake.ID(40588360944263), true, gomock.Any()).Return(&models.Job{
					ID:       snowflake.ID(40588360944263),
					UserID:   snowflake.ID(10086),
					State:    consts.SubStateRunning.State,
					SubState: consts.SubStateRunning.SubState,
				}, nil).Times(1)
				s.sqlmock.ExpectRollback()
			},
			expectedError: common.ErrJobStateNotAllowResume,
			setReq: func() {
				s.allow = func(userID, jobUserID string) bool {
					return true
				}
			},
		},
		{
			name: "update error",
			id:   "4ER",
			mockExpectFunc: func() {
				mockTransaction()

				s.sqlmock.ExpectBegin()
				s.mockJobDao.EXPECT().Get(gomock.Any(), snowflake.ID(12345), true, gomock.Any()).Return(&models.Job{
					ID:       snowflake.ID(12345),
					UserID:   snowflake.ID(10086),
					State:    consts.SubStateInitiallySuspended.State,
					SubState: consts.SubStateInitiallySuspended.SubState,
				}, nil).Times(1)

				s.sqlmock.ExpectExec("UPDATE `job`").
					WillReturnError(fmt.Errorf("update error"))
				s.sqlmock.ExpectRollback()
			},
			setReq: func() {
				s.allow = func(userID, jobUserID string) bool {
					return true
				}
			},
			expectedError: fmt.Errorf("update error"),
		},
	}

	for _, tc := range tastCases {
		tc.Run(s)
	}
}
