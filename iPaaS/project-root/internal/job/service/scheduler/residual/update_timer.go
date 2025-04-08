package residual

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	api "github.com/yuansuan/ticp/common/project-root-api/common"
	schema "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/config"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/consts"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao/models"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/service/v1/application"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/util"
)

// UpdateTimer 残差图更新定时器
type UpdateTimer struct {
	jobDao          dao.JobDao
	residualDao     dao.ResidualDao
	AppSrv          application.Service
	residualHandler ResidualHandler
}

// NewUpdateTimer 创建残差图更新定时器
func NewUpdateTimer(jobDao dao.JobDao, residualDao dao.ResidualDao, appSrv application.Service, residualHandler ResidualHandler) *UpdateTimer {
	return &UpdateTimer{
		jobDao:          jobDao,
		residualDao:     residualDao,
		AppSrv:          appSrv,
		residualHandler: residualHandler,
	}
}

// Run 定时器执行函数
func (s *UpdateTimer) Run(ctx context.Context) {
	s.run(ctx)
}

// 行为：获取所有未完成的残差图，更新
func (s *UpdateTimer) run(ctx context.Context) {
	logger := logging.GetLogger(ctx).With("func", "residual.UpdateTimer.run")
	logger.Info("residual update start...")

	// 获取所有未完成的残差图
	residuals, err := s.getUnfinishedResidual(ctx)
	if err != nil {
		logger.Warnf("residual update: getUnfinishedResidual err: %v", err)
		return
	}

	if len(residuals) == 0 {
		logger.Infof("no unfinished residual")
		return
	}

	zones := config.GetConfig().Zones

	// 调用残差图更新函数
	updated := s.updateResiduals(ctx, residuals, zones)
	logger.Infof("residual update end, updated: %v", updated)
}

// 获取所有未完成的残差图
func (s *UpdateTimer) getUnfinishedResidual(ctx context.Context) ([]*models.Residual, error) {
	// 获取所有未完成的残差图
	residuals, err := s.residualDao.GetUnfinishedResidual(ctx)
	if err != nil {
		return nil, err
	}

	return residuals, nil
}

func (s *UpdateTimer) updateResiduals(ctx context.Context, residuals []*models.Residual, zones schema.Zones) int {
	updated := 0
	for _, residual := range residuals {
		err := s.updateResidual(ctx, residual, zones)
		if err != nil {
			logging.GetLogger(ctx).Warnf("residual update: updateResidual err: %v", err)
			continue
		}
		updated++
	}
	return updated
}

func (s *UpdateTimer) updateResidual(ctx context.Context, residual *models.Residual, zones schema.Zones) error {
	// 获取作业信息
	job, err := s.jobDao.Get(ctx, residual.JobID, false, false)
	if errors.Is(err, common.ErrJobIDNotFound) {
		// 作业不存在，不再更新残差图了
		residual.Finished = true
		residual.FailedReason = api.JobIDNotFound
		residual.UpdateTime = time.Now()
		return s.residualDao.UpdateResidual(ctx, residual)
	} else if err != nil {
		return err
	}

	jobState := consts.NewState(job.State, job.SubState)
	if !jobState.IsRunning() && !jobState.IsTerminating() && !jobState.IsFinal() {
		// 作业状态不需要更新
		return nil
	}

	finished := jobState.IsFinal()
	failedReason := ""
	content, err := s.handlerResidual(ctx, job, residual, zones)
	if err != nil {
		failedReason = err.Error()
		logging.GetLogger(ctx).Info("residual: handlerResidual err: ", err)
		if finished && !checkFinishError(err) { // 不是需要结束的Err，不更新finished，下一次重试
			finished = false
		}
	}

	// 更新residual的字段
	if (content != "" && content != residual.Content) || finished || failedReason != "" {
		residual.Content = content
		residual.Finished = finished
		residual.FailedReason = failedReason

		err = s.residualDao.UpdateResidual(ctx, residual)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *UpdateTimer) handlerResidual(ctx context.Context, job *models.Job, residual *models.Residual, zones schema.Zones) (string, error) {
	logger := logging.GetLogger(ctx)
	residualData, err := s.residualHandler.HandlerHpcResidual(ctx, job, residual, zones)
	if err != nil {
		if errors.Is(err, ErrResidualEmpty) || errors.Is(err, common.ErrPathNotFound) {
			return "", nil // 作业运行初的一段时间stdout.log不会立即生成，生成也会有一小段时间内容为空，暂且认为是正常现象
		}
		logger.Warnf("residual: handlerResidual err: %v", err)
		return "", err
	}

	content, err := util.ResidualMarshal(residualData)
	if err != nil {
		logger.Warnf("residual: marshal: %v", err)
		return "", errors.WithMessage(ErrResidualMarshal, err.Error())
	}

	return content, nil
}
