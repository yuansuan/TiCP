package job

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/system/jobsyncfilestate"
	"xorm.io/xorm"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/consts"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao/models"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/util"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/common/with"
)

type UpdateSyncFileStateSuite struct {
	jobServiceSuite

	// data
}

func TestUpdateSyncFileState(t *testing.T) {
	suite.Run(t, new(UpdateSyncFileStateSuite))
}

func (s *UpdateSyncFileStateSuite) TestUpdateSyncFileState() {
	getJob := &models.Job{FileSyncState: consts.FileSyncStateSyncing.String()}
	downloadTime, _ := util.ParseTime("2024-04-08T15:41:39+08:00", time.RFC3339)
	s.mockJobDao.EXPECT().Transaction(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, action func(context.Context) error) error {
			_, err := s.engine.Transaction(func(db *xorm.Session) (interface{}, error) {
				return nil, action(with.KeepSession(ctx, db))
			})
			return err
		})
	s.mockJobDao.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(getJob, nil).Times(1)

	_ = downloadTime

	s.sqlmock.ExpectBegin()
	s.sqlmock.ExpectExec("UPDATE `job`").
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.sqlmock.ExpectCommit()

	err := s.jobSrv.UpdateSyncFileState(s.ctx, &jobsyncfilestate.Request{
		DownloadFileSizeCurrent: int64(1024),
		DownloadFinished:        false,
		DownloadFinishedTime:    "2024-04-08T15:41:39+08:00",
	}, snowflake.ID(10086).String())
	if !s.NoError(err) {
		return
	}
}

func (s *UpdateSyncFileStateSuite) TestJobError() {
	s.mockJobDao.EXPECT().Transaction(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, action func(context.Context) error) error {
			_, err := s.engine.Transaction(func(db *xorm.Session) (interface{}, error) {
				return nil, action(with.KeepSession(ctx, db))
			})
			return err
		})
	s.mockJobDao.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error job")).Times(1)

	s.sqlmock.ExpectBegin()
	s.sqlmock.ExpectRollback()

	err := s.jobSrv.UpdateSyncFileState(s.ctx, &jobsyncfilestate.Request{}, snowflake.ID(0).String())
	s.Error(err)
	s.T().Log(err)
}

func (s *UpdateSyncFileStateSuite) TestStateError() {
	getJob := &models.Job{FileSyncState: consts.FileSyncStateCompleted.String()}
	s.mockJobDao.EXPECT().Transaction(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, action func(context.Context) error) error {
			_, err := s.engine.Transaction(func(db *xorm.Session) (interface{}, error) {
				return nil, action(with.KeepSession(ctx, db))
			})
			return err
		})
	s.mockJobDao.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(getJob, nil).Times(1)

	s.sqlmock.ExpectBegin()
	s.sqlmock.ExpectRollback()

	err := s.jobSrv.UpdateSyncFileState(s.ctx, &jobsyncfilestate.Request{}, snowflake.ID(10086).String())
	s.Error(err)
	s.T().Log(err)
}

func (s *UpdateSyncFileStateSuite) TestUpdateFailed() {
	getJob := &models.Job{}
	s.mockJobDao.EXPECT().Transaction(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, action func(context.Context) error) error {
			_, err := s.engine.Transaction(func(db *xorm.Session) (interface{}, error) {
				return nil, action(with.KeepSession(ctx, db))
			})
			return err
		})
	s.mockJobDao.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(getJob, nil).Times(1)

	s.engine.DB().Begin()
	s.sqlmock.ExpectBegin()
	s.sqlmock.ExpectExec("UPDATE `job`").
		WillReturnError(fmt.Errorf("update error"))
	s.sqlmock.ExpectRollback()

	err := s.jobSrv.UpdateSyncFileState(s.ctx, &jobsyncfilestate.Request{
		DownloadFileSizeCurrent: int64(1024),
		DownloadFinished:        false,
		DownloadFinishedTime:    "2024-04-08T15:41:39+08:00",
	}, snowflake.ID(10086).String())
	s.Error(err)
	s.T().Log(err)
}
