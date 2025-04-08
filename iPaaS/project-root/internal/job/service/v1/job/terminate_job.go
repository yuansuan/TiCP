package job

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/jobterminate"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/consts"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/common/with"
	"xorm.io/xorm"
)

// Terminate 终止作业
func (srv *jobService) Terminate(ctx context.Context, req *jobterminate.Request, userID snowflake.ID, allow allowFunc) error {
	logger := logging.GetLogger(ctx)
	logger.Info("job terminate start")
	defer logger.Info("job terminate end")

	jobID := snowflake.MustParseString(req.JobID)
	return srv.jobdao.Transaction(ctx, func(c context.Context) error {
		// JobID存在性验证
		job, err := srv.jobdao.Get(c, jobID, true, false)
		if err != nil {
			if !errors.Is(err, common.ErrJobIDNotFound) {
				logger.Warnf("get Job error! err: %v", err)
			}

			return err // internal error
		}

		// 用户权限验证
		if !allow(userID.String(), job.UserID.String()) {
			logger.Warnf("no permission to operate other's job")
			return errors.WithMessage(common.ErrJobAccessDenied, "no permission to operate other's job")
		}

		// 验证作业状态
		state, ok := consts.GetStateBySubState(job.SubState)
		if !ok {
			logger.Warnf("invalid job sub state: %v", job.SubState)
			return fmt.Errorf("invalid job sub state: %v", job.SubState)
		}

		if !state.CanTerminate() {
			logger.Warnf("job can not terminate, state: %v", job.State)
			return errors.WithMessage(common.ErrJobStateNotAllowTerminate, fmt.Sprintf("job can not terminate, state: %v", job.State))
		}

		update := false
		end := false

		switch job.SubState {
		case consts.SubStateInitiated.SubState, consts.SubStateInitiallySuspended.SubState, consts.SubStateScheduling.SubState: // 提交暂停中，中央调度中
			job.State = consts.SubStateTerminated.State
			job.SubState = consts.SubStateTerminated.SubState
			end = true
			job.StateReason = consts.ParseAndUpdateStateReasonString(job.StateReason, consts.SubStateTerminated, consts.StateReasonUserCancel).String()
			update = true
		case consts.SubStateFileUploading.SubState, consts.SubStateHpcWaiting.SubState, consts.SubStateRunning.SubState, consts.SubStateSuspended.SubState:
			// 文件上传中、HPC等待中、运行中、暂停中
			job.State = consts.SubStateTerminating.State
			job.SubState = consts.SubStateTerminating.SubState
			job.StateReason = consts.ParseAndUpdateStateReasonString(job.StateReason, consts.SubStateTerminating, consts.StateReasonUserCancel).String()
			update = true
		default:
			// do nothing
		}

		now := time.Now()
		job.TerminatingTime = now
		job.UserCancel = consts.UserCancel

		if update {
			updateResult := int64(0)
			err = with.DefaultSession(c, func(db *xorm.Session) error {
				db = db.ID(jobID)
				if end {
					job.EndTime = now
					db = db.Cols("end_time")
				}

				// 作业被终止的时候，还未提交到hpc，则传输状态直接变为完成
				if job.HPCJobID == "" {
					job.FileSyncState = consts.FileSyncStateCompleted.String()
					job.DownloadFinished = 1
					job.DownloadTime = now
					job.TransmittingTime = now
				}

				// 修改作业状态为Terminating
				updateResult, err = db.Cols("state", "sub_state", "terminating_time", "state_reason", "user_cancel", "file_sync_state", "download_finished", "download_time", "transmitting_time").Update(job)
				return err
			})
			if err != nil {
				logger.Warnf("session.Update job error: %v", err)
				return err // internal error
			}
			logger.With("updateResult", updateResult).Infof("update job info complete")
		}

		return nil
	})
}
