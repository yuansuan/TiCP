package job

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/jobdelete"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/consts"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao/models"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/common/with"
	"xorm.io/xorm"
)

// Delete 删除作业
func (srv *jobService) Delete(ctx context.Context, req *jobdelete.Request, userID snowflake.ID, allow allowFunc) error {
	logger := logging.GetLogger(ctx)
	logger.Info("job delete start")
	defer logger.Info("job delete end")

	jobID := snowflake.MustParseString(req.JobID)
	return srv.jobdao.Transaction(ctx, func(c context.Context) error {
		// JobID存在性验证
		job, err := srv.jobdao.Get(c, jobID, true, true) // 删除作业接口保证幂等，这里存在性校验查询包含已删除的作业
		if err != nil {
			if !errors.Is(err, common.ErrJobIDNotFound) {
				logger.Warnf("get Job error! err: %v", err)
			}

			return err // internal error OR job not exist
		}

		logger.Infof("job info: %+v", job)

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

		if !state.IsFinal() {
			logger.Warnf("job is not in final state: %v", job.SubState)
			return errors.WithMessage(common.ErrJobStateNotAllowDelete, fmt.Sprintf("job state is %v, not in final state [%v]", state.StateString(), consts.FinalString()))
		}

		// 判断文件传输状态是否是可删除的状态
		if state.IsCompleted() && job.FileSyncState != "" && !consts.FileSyncState(job.FileSyncState).CanDelete() {
			logger.Warnf("job file sync state can not delete: %v", job.FileSyncState)
			return errors.WithMessage(common.ErrJobStateNotAllowDelete, fmt.Sprintf("job file sync state is %v, can not delete", job.FileSyncState))
		}

		updateResult := int64(0)
		err = with.DefaultSession(c, func(db *xorm.Session) error {
			// 修改作业is_deleted字段
			updateResult, err = db.ID(jobID).Cols("is_deleted").Update(&models.Job{
				ID:        jobID,
				IsDeleted: 1,
			})
			return err
		})
		if err != nil {
			logger.Warnf("session.Update job error: %v", err)
			return err // internal error
		}

		logger.With("updateResult", updateResult).Infof("update job info complete")
		return nil
	})
}
