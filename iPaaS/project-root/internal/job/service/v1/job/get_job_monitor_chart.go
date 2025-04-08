package job

import (
	"context"

	"github.com/pkg/errors"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/jobmonitorchart"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao/models"
)

// GetMonitorChart 获取作业监控图表
func (srv *jobService) GetMonitorChart(ctx context.Context, req *jobmonitorchart.Request, userID snowflake.ID, allow allowFunc) (*models.MonitorChart, error) {
	logger := logging.GetLogger(ctx)
	logger.Info("job get monitor chart start")
	defer logger.Info("job get monitor chart end")

	jobID := snowflake.MustParseString(req.JobID)
	monitorChart := &models.MonitorChart{}
	err := srv.jobdao.Transaction(ctx, func(c context.Context) error {
		// JobID存在性验证
		job, err := srv.jobdao.Get(c, jobID, false, false)
		if err != nil {
			if !errors.Is(err, common.ErrJobIDNotFound) {
				logger.Warnf("get Job error! err: %v", err)
			}
			return err // internal error OR job not exist
		}

		// 用户权限验证
		if !allow(userID.String(), job.UserID.String()) {
			logger.Warnf("no permission to operate other's job")
			return errors.WithMessage(common.ErrJobAccessDenied, "no permission to operate other's job")
		}

		// 获取作业监控图表
		monitorChart, err = srv.jobdao.GetJobMonitorChart(c, jobID, false)
		if err != nil {
			logger.Warnf("session.GetJobMonitorChart error: %v", err)
			return err
		}
		return nil
	})

	return monitorChart, err
}
