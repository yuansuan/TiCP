package service

import (
	"context"

	"github.com/gin-gonic/gin"

	"github.com/yuansuan/ticp/PSP/psp/internal/monitor/dto"
)

type NodeService interface {
	GetNodeInfo(ctx context.Context, nodeName string) (*dto.NodeDetail, error)
	GetNodeList(ctx context.Context, nodeName string, index, size int64) ([]*dto.NodeInfo, int64, error)
	NodeOperate(ctx context.Context, nodeNames []string, operation string) error
	NodeCoreNum(ctx context.Context) (*dto.CoreStatistics, error)
}

type DashboardService interface {
	GetClusterInfo(ctx context.Context) (*dto.ClusterInfo, []*dto.NodeDetail, *dto.Disk, error)
	GetResourceInfo(ctx context.Context, req *dto.Request) ([]*dto.ValueStruct, []*dto.ValueStruct, []*dto.ValueStruct, error)
	GetJobInfo(ctx context.Context, req *dto.Request) (jobResRange []*dto.JobStatusValue, jobResLatest []*dto.JobStatusValue, err error)
}

type ReportService interface {
	GetHostResourceMetricUtAvg(ctx *gin.Context, req *dto.UniteReportReq) (*dto.ResourceUtAvgReportResp, error)
	GetDiskUtAvg(ctx *gin.Context, req *dto.UniteReportReq) (*dto.DiskUtAvgReportResp, error)
	GetCPUTimeSum(ctx *gin.Context, req *dto.UniteReportReq) (*dto.CPUTimeSumMetricsResp, error)
	GetJobCount(ctx *gin.Context, req *dto.UniteReportReq) (*dto.JobCountMetricResp, error)
	GetJobDeliverCount(ctx *gin.Context, req *dto.UniteReportReq) (*dto.JobDeliverCountResp, error)
	GetJobWaitStatistic(ctx *gin.Context, req *dto.UniteReportReq) (*dto.JobWaitStatisticResp, error)
	GetLicenseAppUsedUtAvg(ctx *gin.Context, req *dto.UniteReportReq) (*dto.LicenseAppUsedUtAvgReportResp, error)
	GetLicenseAppModuleUsedUtAvg(ctx *gin.Context, req *dto.LicenseAppModuleUsedUtAvgReq) (*dto.LicenseAppModuleUsedUtAvgReportResp, error)
	NodeDownStatisticReport(ctx *gin.Context, req *dto.NodeDownStatisticReportReq) (*dto.NodeDownStatisticReportResp, error)
	ExportResourceUtAvg(ctx *gin.Context, req *dto.UniteReportReq) error
	ExportDiskUtAvg(ctx *gin.Context, req *dto.UniteReportReq) error
	ExportCPUTimeSum(ctx *gin.Context, req *dto.UniteReportReq) error
	ExportJobCount(ctx *gin.Context, req *dto.UniteReportReq) error
	ExportJobDeliverCount(ctx *gin.Context, req *dto.UniteReportReq) error
	ExportJobWaitStatistic(ctx *gin.Context, req *dto.UniteReportReq) error
	ExportLicenseAppUsedUtAvg(ctx *gin.Context, req *dto.UniteReportReq) error
	ExportNodeDownStatistics(ctx *gin.Context, req *dto.ExportNodeDownStatisticsReq) error
	ExportReport(ctx *gin.Context, req *dto.UniteReportReq) error
}
