package job

import (
	"context"

	"github.com/pkg/errors"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/jobresidual"
	schema "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/config"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao/models"
	resipkg "github.com/yuansuan/ticp/iPaaS/project-root/internal/job/service/scheduler/residual"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/util"
)

// GetResidual 获取作业残差图
func (srv *jobService) GetResidual(ctx context.Context, req *jobresidual.Request, userID snowflake.ID, allow allowFunc) (*schema.Residual, error) {
	logger := logging.GetLogger(ctx)

	logger.Info("job get residual start")
	defer logger.Info("job get residual end")

	jobID := snowflake.MustParseString(req.JobID)
	job := &models.Job{}
	residual := &models.Residual{}
	err := srv.jobdao.Transaction(ctx, func(c context.Context) error {
		// JobID存在性验证
		var err error
		job, err = srv.jobdao.Get(c, jobID, false, false)
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

		// 获取作业残差图
		residual, err = srv.residualdao.GetJobResidual(c, jobID)
		if err != nil {
			logger.Warnf("session.GetJobResidual error: %v", err)
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	if !residual.Finished {
		// 实时从超算读
		resp, err := srv.residualHandler.HandlerHpcResidual(ctx, job, residual, config.GetConfig().Zones)
		if err != nil {
			if errors.Is(err, resipkg.ErrResidualEmpty) || errors.Is(err, common.ErrPathNotFound) {
				return nil, nil // 作业运行初的一段时间stdout.log不会立即生成，生成也会有一小段时间内容为空，暂且认为是正常现象
			}
			return nil, errors.Wrap(common.ErrHpcResidual, err.Error())
		}

		return resp, nil
	} else {
		resp, err := util.ModelToOpenAPIJobResidual(residual)
		if err != nil {
			logger.Warnf("util.ModelToOpenAPIJobResidual error: %v", err)
			return nil, err
		}
		return resp, nil
	}
}
