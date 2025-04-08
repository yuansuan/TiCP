package job

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/jobtransmitsuspend"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/consts"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/common/with"
	"xorm.io/xorm"
)

// TransmitSuspend 传输暂停
func (srv *jobService) TransmitSuspend(ctx context.Context, req *jobtransmitsuspend.Request, userID snowflake.ID, allow allowFunc) error {
	logger := logging.GetLogger(ctx)

	logger.Info("job transmit suspend start")
	defer logger.Info("job transmit suspend end")

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

		// 验证作业传输状态
		fileState := consts.FileSyncState(job.FileSyncState)
		if !fileState.CanTransmitSuspended() {
			logger.Warnf("job can not transmit suspend, file sync state: %v", job.FileSyncState)
			return errors.WithMessage(common.ErrJobStateNotAllowTransmitSuspend, fmt.Sprintf("job can not transmit suspend, file sync state: %v", job.FileSyncState))
		}

		if fileState == consts.FileSyncStateSyncing {
			// 修改传输状态为pausing
			job.FileSyncState = consts.FileSyncStatePausing.String()
			updateResult := int64(0)
			err = with.DefaultSession(c, func(db *xorm.Session) error {
				updateResult, err = db.ID(jobID).Cols("file_sync_state").Update(job)
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
