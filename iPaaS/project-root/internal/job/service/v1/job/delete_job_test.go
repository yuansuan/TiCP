package job

import (
	"context"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/jobdelete"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/consts"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao/models"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/common/with"
	"xorm.io/xorm"
)

type DeleteSuite struct {
	jobServiceSuite

	// data
	allow allowFunc
}

func TestDelete(t *testing.T) {
	suite.Run(t, new(DeleteSuite))
}

type JobDeleteTestCase struct {
	name           string
	id             string
	expectedError  error
	mockExpectFunc func()
	setReq         func()
}

func (tc JobDeleteTestCase) Run(s *DeleteSuite) {
	s.Run(tc.name, func() {
		if tc.mockExpectFunc != nil {
			tc.mockExpectFunc()
		}

		if tc.setReq != nil {
			tc.setReq()
		}

		// do
		err := s.jobSrv.Delete(s.ctx, &jobdelete.Request{
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

func (s *DeleteSuite) TestDelete() {
	mockTransaction := func() {
		s.mockJobDao.EXPECT().Transaction(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, action func(context.Context) error) error {
				_, err := s.engine.Transaction(func(db *xorm.Session) (interface{}, error) {
					return nil, action(with.KeepSession(ctx, db))
				})
				return err
			})
	}
	tastCases := []JobDeleteTestCase{
		{
			name: "normal",
			id:   "4ER",
			mockExpectFunc: func() {
				mockTransaction()

				s.sqlmock.ExpectBegin()
				s.mockJobDao.EXPECT().Get(gomock.Any(), snowflake.ID(12345), true, gomock.Any()).Return(&models.Job{
					ID:       snowflake.ID(12345),
					UserID:   snowflake.ID(10086),
					State:    consts.SubStateCompleted.State,
					SubState: consts.SubStateCompleted.SubState,
				}, nil).Times(1)

				s.sqlmock.ExpectExec("UPDATE `job`").
					WithArgs(1, snowflake.ID(12345)).
					WillReturnResult(sqlmock.NewResult(1, 1))
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
					State:    consts.SubStateCompleted.State,
					SubState: consts.SubStateCompleted.SubState,
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
			expectedError: fmt.Errorf("invalid job sub state"),
			setReq: func() {
				s.allow = func(userID, jobUserID string) bool {
					return true
				}
			},
		},
		{
			name: "job state not allow delete",
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
			expectedError: common.ErrJobStateNotAllowDelete,
			setReq: func() {
				s.allow = func(userID, jobUserID string) bool {
					return true
				}
			},
		},
		{
			name: "job file state not allow delete",
			id:   "jobState",
			mockExpectFunc: func() {
				mockTransaction()

				s.sqlmock.ExpectBegin()
				s.mockJobDao.EXPECT().Get(gomock.Any(), snowflake.ID(40588360944263), true, gomock.Any()).Return(&models.Job{
					ID:            snowflake.ID(40588360944263),
					UserID:        snowflake.ID(10086),
					State:         consts.SubStateCompleted.State,
					SubState:      consts.SubStateCompleted.SubState,
					FileSyncState: consts.FileSyncStateSyncing.String(),
				}, nil).Times(1)
				s.sqlmock.ExpectRollback()
			},
			expectedError: common.ErrJobStateNotAllowDelete,
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
					State:    consts.SubStateCompleted.State,
					SubState: consts.SubStateCompleted.SubState,
				}, nil).Times(1)

				s.sqlmock.ExpectExec("UPDATE `job`").
					WithArgs(1, snowflake.ID(12345)).
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
