package job

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/jobresume"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/consts"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao/models"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/common/with"
	"xorm.io/xorm"
)

// Resume 恢复作业
func (srv *jobService) Resume(ctx context.Context, req *jobresume.Request, userID snowflake.ID, allow allowFunc) error {
	logger := logging.GetLogger(ctx)
	logger.Info("resume job start")
	defer logger.Info("resume job end")

	// 作业读取
	jobID := snowflake.MustParseString(req.JobID)
	return srv.jobdao.Transaction(ctx, func(c context.Context) error {
		job, err := srv.jobdao.Get(c, jobID, true, false)
		if err != nil {
			if !errors.Is(err, common.ErrJobIDNotFound) {
				logger.Warnf("get Job error! err: %v", err)
			}

			return err
		}

		// 用户权限验证
		if !allow(userID.String(), job.UserID.String()) {
			logger.Warnf("no permission to operate other's job")
			return errors.WithMessage(common.ErrJobAccessDenied, "no permission to operate other's job")
		}

		// 状态验证，Suspended 和 InitiallySuspended 可以恢复
		state, ok := consts.GetStateBySubState(job.SubState)
		if !ok {
			logger.Warnf("invalid job state! job id: %v, job state: %v", jobID, state)
			return fmt.Errorf("invalid job state! job id: %v, job state: %v", jobID, state)
		}

		if !state.CanResume() {
			logger.Warnf("job can not resume, state: %v", job.State)
			return common.ErrJobStateNotAllowResume
		}

		// 更新job 状态, subStateInitiallySuspended --> pending, suspended --> ? TODO:记录原状态
		if state.SubState == consts.SubStateInitiallySuspended.SubState {
			sr := consts.ParseAndUpdateStateReasonString(job.StateReason, consts.SubStateScheduling, "Job Scheduling (Resume)")
			now := time.Now()
			updateResult := int64(0)
			err = with.DefaultSession(c, func(db *xorm.Session) error {
				updateResult, err = db.ID(jobID).Cols("state", "sub_state", "state_reason", "update_time", "pending_time").
					Update(&models.Job{
						ID:          jobID,
						State:       consts.SubStateScheduling.State,
						SubState:    consts.SubStateScheduling.SubState,
						StateReason: sr.String(),
						UpdateTime:  now,
						PendingTime: now,
					})
				return err
			})
			if err != nil {
				logger.Warnf("session.Update job error: %v", err)
				return err // internal error
			}

			logger.Infof("update job info complete, job id: %v, update result: %v", jobID, updateResult)
		}
		return nil
	})
}
