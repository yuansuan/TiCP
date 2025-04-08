package job

import (
	"context"
	"time"

	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/common/with"
	"xorm.io/xorm"

	"github.com/pkg/errors"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/system/jobsyncfilestate"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/consts"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/util"
)

func (srv *jobService) UpdateSyncFileState(ctx context.Context, req *jobsyncfilestate.Request, jobIDStr string) error {
	jobID := snowflake.MustParseString(jobIDStr)
	return srv.jobdao.Transaction(ctx, func(ctx context.Context) error {
		job, err := srv.jobdao.Get(ctx, jobID, true, true)
		if err != nil {
			if !errors.Is(err, common.ErrJobIDNotFound) {
				logging.GetLogger(ctx).Warnf("get Job:%s error! err: %v", jobIDStr, err)
			}

			return err
		}

		targetFileSyncState := consts.FileSyncState(req.FileSyncState)
		curJobFileSyncState := consts.FileSyncState(job.FileSyncState)
		logging.GetLogger(ctx).Infof("UpdateSyncFileState job: %v current file sync state is: %v, target file sync state: %v", jobIDStr, curJobFileSyncState, targetFileSyncState)

		if curJobFileSyncState.IsFinal() {
			return errors.WithMessagef(common.ErrJobFileSyncStateUpdateFailed, "update state error, current file sync state is final state: %v, request state:%v, job id: %v ", job.FileSyncState, targetFileSyncState, jobID)
		}

		if curJobFileSyncState == consts.FileSyncStatePaused {
			return errors.WithMessagef(common.ErrJobFileSyncStateUpdateFailed, "update state error, current file sync state is paused state: %v, request state:%v,  job id: %v ", job.FileSyncState, targetFileSyncState, jobID)
		}

		job.UpdateTime = time.Now()
		job.DownloadFileSizeCurrent = req.DownloadFileSizeCurrent
		job.DownloadFileSizeTotal = req.DownloadFileSizeTotal
		if req.DownloadFinished {
			if !targetFileSyncState.IsFinal() {
				return errors.WithMessagef(common.ErrJobFileSyncStateUpdateFailed, "update state error, download finished shouled be final state, request file sycn state: %v", targetFileSyncState)
			}

			job.DownloadFinished = 1
			if req.DownloadFinishedTime != "" {
				downloadTime, _ := util.ParseTime(req.DownloadFinishedTime, time.RFC3339)
				job.DownloadTime = downloadTime
			}

			// 如果已经下载完成，矫正当前下载数据
			if job.DownloadFileSizeTotal != job.DownloadFileSizeCurrent {
				job.DownloadFileSizeCurrent = job.DownloadFileSizeTotal
			}
		} else {
			job.DownloadFinished = 0
		}

		if job.State == consts.Completed || job.State == consts.Failed || job.State == consts.Terminated {
			job.TransmittingTime = job.EndTime
		}

		if curJobFileSyncState == consts.FileSyncStatePausing {
			job.FileSyncState = consts.FileSyncStatePaused.String()
		} else {
			job.FileSyncState = targetFileSyncState.String()
		}

		return with.DefaultSession(ctx, func(session *xorm.Session) error {
			_, err = session.ID(jobID).Cols("file_sync_state, download_file_size_current, download_file_size_total, download_finished, download_time, transmitting_time, update_time").Update(job)
			return err
		})
	})
}
