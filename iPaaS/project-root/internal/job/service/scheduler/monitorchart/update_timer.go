package monitorchart

import (
	"context"
	"fmt"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	schema "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/config"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/consts"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao/models"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/module/storage"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/service/v1/application"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/util"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/common/with"
	"xorm.io/xorm"
)

// UpdateTimer 监控图表更新定时器
type UpdateTimer struct {
	JobDao dao.JobDao
	AppSrv application.Service
}

// NewUpdateTimer 新建监控图表更新定时器
func NewUpdateTimer(jobDao dao.JobDao, appSrv application.Service) *UpdateTimer {
	return &UpdateTimer{
		JobDao: jobDao,
		AppSrv: appSrv,
	}
}

// Run 定时器执行函数
func (s *UpdateTimer) Run(ctx context.Context) {
	s.run(ctx)
}

func (s *UpdateTimer) run(ctx context.Context) {
	logger := logging.GetLogger(ctx).With("func", "monitorchart.UpdateTimer.run")
	logger.Info("start update monitor chart")

	// 获取所有未完成的监控图表
	unfinishedCharts, err := s.getUnfinishedMonitorChart(ctx, logger)
	if err != nil {
		logger.Warnf("get unfinished monitor chart error! err: %v", err)
		return
	}

	if len(unfinishedCharts) == 0 {
		logger.Info("no unfinished monitor chart")
		return
	}

	zones := config.GetConfig().Zones

	updated := s.updateMonitorCharts(ctx, zones, unfinishedCharts)
	logger.Infof("update monitor chart finished, updated: %d", updated)

}

func (s *UpdateTimer) getUnfinishedMonitorChart(ctx context.Context, logger *logging.Logger) ([]*models.MonitorChart, error) {
	// 获取所有未完成的监控图表
	monitorCharts, err := s.JobDao.GetUnfinishedbMonitorChart(ctx)
	if err != nil {
		return nil, err
	}

	return monitorCharts, nil
}

func (s *UpdateTimer) updateMonitorCharts(ctx context.Context, zones schema.Zones, unfinishedCharts []*models.MonitorChart) int {
	updated := 0
	for _, unfinishedChart := range unfinishedCharts {
		err := s.updateMonitorChart(ctx, unfinishedChart, zones)
		if err != nil {
			logging.GetLogger(ctx).Warnf("update monitor chart error! err: %v", err)
			continue
		}
		updated++
	}
	return updated
}

func (s *UpdateTimer) updateMonitorChart(ctx context.Context, unfinishedChart *models.MonitorChart, zones schema.Zones) error {
	logger := logging.GetLogger(ctx)

	finished := unfinishedChart.Finished

	// 获取作业信息
	job, err := s.JobDao.Get(ctx, unfinishedChart.JobID, false, true)
	if err != nil {
		if errors.Is(err, common.ErrJobIDNotFound) && !finished {
			unfinishedChart.Finished = true
			unfinishedChart.FailedReason = "job not found"
			err = s.updateMonitorChartInDB(ctx, logger, unfinishedChart)
			if err != nil {
				return err
			}
		}
		return err
	}

	failedReason := ""
	jobState := consts.NewState(job.State, job.SubState)
	if !jobState.IsRunning() && !jobState.IsTerminating() && !jobState.IsFinal() {
		// 作业状态不需要更新
		return nil
	}

	finished = jobState.IsFinal()

	content, err := s.handlerMonitorChart(ctx, job, unfinishedChart, zones)
	if err != nil {
		failedReason = err.Error()
		logger.Info("monitorChart: handlerMonitorChart err: ", err)
		if finished && !checkFinishError(err) { // 不是需要结束的Err，不更新finished，下一次重试
			finished = false
		}
	}

	// 更新monitor chart
	if (content != "" && content != unfinishedChart.Content) || finished == true || failedReason != "" {
		unfinishedChart.Content = content
		unfinishedChart.Finished = finished
		unfinishedChart.FailedReason = failedReason

		err = s.updateMonitorChartInDB(ctx, logger, unfinishedChart)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *UpdateTimer) handlerMonitorChart(ctx context.Context, job *models.Job, unfinishedChart *models.MonitorChart, zones schema.Zones) (string, error) {
	logger := logging.GetLogger(ctx)
	monitorChartReg := unfinishedChart.MonitorChartRegexp
	if monitorChartReg == "" {
		monitorChartReg = consts.DefaultMonitorChartRegexp
	}
	monitorChartParser := unfinishedChart.MonitorChartParser
	if monitorChartParser == "" {
		logger.Warnf("monitorChart: parser is empty")
		return "", errors.WithMessage(ErrMonitorChartParser, "parser is empty")
	}

	p := NewParser(monitorChartParser)
	if p == nil {
		logger.Warnf("monitorChart: unsupported parser: %v", monitorChartParser)
		return "", errors.WithMessage(ErrMonitorChartParser, fmt.Sprintf("unsupported parser: %v", monitorChartParser))
	}

	zone, ok := zones[job.Zone]
	if !ok {
		logger.Warnf("monitorChart: zone not found: %v", job.Zone)
		return "", errors.WithMessage(ErrMonitorChartJobInfo, fmt.Sprintf("zone not found: %v", job.Zone))
	}

	zoneDomain := zone.HPCEndpoint
	if zoneDomain == "" {
		logger.Warnf("monitorChart: zone domain is empty")
		return "", errors.WithMessage(ErrMonitorChartJobInfo, "zone domain is empty")
	}

	clientParams := storage.ClientParams{
		Endpoint: zoneDomain,
		Timeout:  0, // 0 for no timeout
		AdminAPI: true,
	}

	workdir := job.WorkDir
	if workdir == "" {
		logger.Warnf("monitorChart: workdir is empty")
		return "", errors.WithMessage(ErrMonitorChartJobInfo, "workdir is empty")
	}

	workdir = strings.TrimPrefix(workdir, zoneDomain)
	workdir = util.AddPrefixSlash(workdir)
	workdir = util.AddSuffixSlash(workdir)

	logger.Info("monitorChart: readAt path: ", workdir+monitorChartReg)
	logger.Info("monitorChart: readAt client base url: ", clientParams.Endpoint)

	result, err := readAndParseMonitorChart(ctx, clientParams, workdir, monitorChartReg, p)
	if err != nil {
		logger.Warnf("monitorChart: readAndParseMonitorChart err: %v", err)
		return "", err
	}

	// store result
	monitorChartData := ConvertMonitorChart(result)
	content, err := util.MonitorChartMarshal(monitorChartData)
	if err != nil {
		logger.Warnf("monitorChart: marshal: %v", err)
		return "", errors.WithMessage(ErrMonitorChartMarshal, err.Error())
	}

	return content, nil
}

func (s *UpdateTimer) updateMonitorChartInDB(ctx context.Context, logger *logging.Logger, chart *models.MonitorChart) error {
	logger.Infof("%s", spew.Sdump(chart))

	err := with.DefaultSession(ctx, func(db *xorm.Session) error {
		_, err := db.ID(chart.ID).UseBool().Update(chart)
		return err
	})
	return err
}
