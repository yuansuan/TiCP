package job

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/jobtransmitresume"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/consts"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/common/with"
	"xorm.io/xorm"
)

// TransmitResume 传输恢复
func (srv *jobService) TransmitResume(ctx context.Context, req *jobtransmitresume.Request, userID snowflake.ID, allow allowFunc) error {
	logger := logging.GetLogger(ctx)

	logger.Info("job transmit resume start")
	defer logger.Info("job transmit resume end")

	jobID := snowflake.MustParseString(req.JobID)
	return srv.jobdao.Transaction(ctx, func(c context.Context) error {

		// JobID存在性验证
		forUpdate := true
		job, err := srv.jobdao.Get(c, jobID, forUpdate, false)
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
		if !fileState.CanTransmitResume() {
			logger.Warnf("job can not transmit resume, file sync state: %v", job.FileSyncState)
			return errors.WithMessage(common.ErrJobStateNotAllowTransmitResume, fmt.Sprintf("job can not transmit resume, file sync state: %v", job.FileSyncState))
		}

		if fileState == consts.FileSyncStatePaused {
			// 修改传输状态为resuming
			job.FileSyncState = consts.FileSyncStateResuming.String()
			updateResult := int64(0)
			err = with.DefaultSession(c, func(db *xorm.Session) error {
				updateResult, err = db.ID(jobID).ForUpdate().Cols("file_sync_state").Update(job)
				return err
			})
			if err != nil {
				logger.Errorf("session.Update job error: %v", err)
				return err // internal error
			}
			logger.With("updateResult", updateResult).Infof("update job info complete")
		}
		return nil
	})
}
