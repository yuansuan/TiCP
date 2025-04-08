package api

import (
	"fmt"
	stdlog "log"
	nethttp "net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promlog"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/http"
	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/yuansuan/ticp/PSP/psp/internal/monitor/collector"
	"github.com/yuansuan/ticp/PSP/psp/internal/monitor/service"
	"github.com/yuansuan/ticp/PSP/psp/internal/monitor/service/impl"
)

type apiRoute struct {
	nodeService      service.NodeService
	dashboardService service.DashboardService
	reportService    service.ReportService
}

func NewAPIRoute() (*apiRoute, error) {
	nodeService, err := impl.NewNodeService()
	if err != nil {
		return nil, err
	}
	dashBoardService, err := impl.NewDashBoardService()
	if err != nil {
		return nil, err
	}

	reportService, err := impl.NewReportService()
	if err != nil {
		return nil, err
	}

	return &apiRoute{
		nodeService:      nodeService,
		dashboardService: dashBoardService,
		reportService:    reportService,
	}, nil
}

// InitAPI 初始化API服务
func InitAPI(drv *http.Driver) {
	logger := logging.Default()

	api, err := NewAPIRoute()
	if err != nil {
		logger.Errorf("init api service err: %v", err)
		panic(err)
	}

	promLogger := promlog.New(&promlog.Config{})
	handler, err := NewHandler(promLogger)
	if err != nil {
		_ = level.Error(promLogger).Log("Error occur when start server ", err)
		return
	}

	drv.GET("/metrics", gin.WrapH(handler))

	group := drv.Group("/api/v1")
	{
		nodeGroup := group.Group("/node")
		nodeGroup.GET("/detail", api.NodeDetail)
		nodeGroup.GET("/list", api.NodeList)
		nodeGroup.POST("/operate", api.NodeOperate)
		nodeGroup.GET("/coreNum", api.NodeCoreNum)

		dashboardGroup := group.Group("/dashboard")
		dashboardGroup.GET("/clusterInfo", api.ClusterInfo)
		dashboardGroup.GET("/resourceInfo", api.ResourceInfo)
		dashboardGroup.GET("/jobInfo", api.JobInfo)
	}

	{
		reportGroup := group.Group("/report")
		reportGroup.GET("/resourceUtAvg", api.ResourceUtAvgReport)
		reportGroup.GET("/diskUtAvg", api.DiskUTAvgReport)
		reportGroup.GET("/cpuTimeSum", api.CPUTimeReport)
		reportGroup.GET("/jobCount", api.JobCountReport)
		reportGroup.GET("/jobDeliverCount", api.JobDeliverCountReport)
		reportGroup.GET("/jobWaitStatistic", api.JobWaitStatisticReport)
		reportGroup.GET("/licenseAppUsedUtAvg", api.LicenseAppUsedAvgReport)
		reportGroup.GET("/licenseAppModuleUsedUtAvg", api.LicenseAppModuleUsedUtAvgReport)
		reportGroup.GET("/nodeDownStatistics", api.NodeDownStatisticReport)

		reportGroup.GET("/export", api.ReportExport)
		reportGroup.GET("/export/resourceUtAvg", api.ExportResourceUtAvg)
		reportGroup.GET("/export/diskUtAvg", api.ExportDiskUtAvg)
		reportGroup.GET("/export/cpuTimeSum", api.ExportCPUTimeSum)
		reportGroup.GET("/export/jobCount", api.ExportJobCount)
		reportGroup.GET("/export/jobDeliverCount", api.ExportJobDeliverCount)
		reportGroup.GET("/export/jobWaitStatistic", api.ExportJobWaitStatistic)
		reportGroup.GET("/export/licenseAppUsedUtAvg", api.ExportLicenseAppUsedUtAvg)
		reportGroup.GET("/export/nodeDownStatistics", api.ExportNodeDownStatistics)
	}
}

func NewHandler(logger log.Logger) (nethttp.Handler, error) {
	nc, err := collector.NewCollector(logger)
	if err != nil {
		return nil, fmt.Errorf("couldn't create collector: %s", err)
	}

	registry := prometheus.NewRegistry()
	if err := registry.Register(nc); err != nil {
		return nil, fmt.Errorf("couldn't register node collector: %s", err)
	}

	handler := promhttp.HandlerFor(
		prometheus.Gatherers{registry},
		promhttp.HandlerOpts{
			ErrorLog:      stdlog.New(log.NewStdlibAdapter(level.Error(logger)), "", 0),
			ErrorHandling: promhttp.ContinueOnError,
		})

	return handler, nil
}
