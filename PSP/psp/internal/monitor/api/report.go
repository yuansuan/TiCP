package api

import (
	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/internal/monitor/dto"
	"github.com/yuansuan/ticp/PSP/psp/internal/monitor/util"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/ginutil"
)

// ResourceUtAvgReport
//
//	@Summary		资源类型报表, 包括 cpu平均利用率(type=CPU_UT_AVG)， 内存利用率(type=MEM_UT_AVG)，磁盘吞吐率(type=TOTAL_IO_UT_AVG)
//	@Description	资源类型报表
//	@Tags			报表
//	@Accept			json
//	@Produce		json
//	@Param			param	body		dto.UniteReportReq	true	"请求参数"
//	@Response		200		{object}	dto.ResourceUtAvgReportResp
//	@Router			/report/resourceUtAvg [get]
func (r *apiRoute) ResourceUtAvgReport(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	req := dto.UniteReportReq{}

	if err := ctx.BindQuery(&req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	if util.ValidReportType(req.ReportType) {
		logger.Errorf("report type invalid")
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	resp, err := r.reportService.GetHostResourceMetricUtAvg(ctx, &req)
	if err != nil {
		logger.Errorf("get report data err: %v", err)
		ginutil.Error(ctx, errcode.ErrCommonReportGetFailed, errcode.MsgCommonReportFailGet)
		return
	}

	ginutil.Success(ctx, resp)
}

// DiskUTAvgReport
//
//	@Summary		磁盘使用率报表, 包括 磁盘使用率(type=DISK_UT_AVG)
//	@Description	磁盘使用率报表
//	@Tags			报表
//	@Accept			json
//	@Produce		json
//	@Param			param	body		dto.UniteReportReq	true	"请求参数"
//	@Response		200		{object}	dto.DiskUtAvgReportResp
//	@Router			/report/diskUtAvg [get]
func (r *apiRoute) DiskUTAvgReport(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	req := dto.UniteReportReq{}

	if err := ctx.BindQuery(&req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	resp, err := r.reportService.GetDiskUtAvg(ctx, &req)
	if err != nil {
		logger.Errorf("get report data err: %v", err)
		ginutil.Error(ctx, errcode.ErrCommonReportGetFailed, errcode.MsgCommonReportFailGet)
		return
	}

	ginutil.Success(ctx, resp)
}

// CPUTimeReport
//
//	@Summary		核时使用情况报表, 包括 核时使用情况(type=CPU_TIME_SUM)
//	@Description	核时使用情况报表
//	@Tags			报表
//	@Accept			json
//	@Produce		json
//	@Param			param	body		dto.UniteReportReq	true	"请求参数"
//	@Response		200		{object}	dto.CPUTimeSumMetricsResp
//	@Router			/report/cpuTimeSum [get]
func (r *apiRoute) CPUTimeReport(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	req := dto.UniteReportReq{}

	if err := ctx.BindQuery(&req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	resp, err := r.reportService.GetCPUTimeSum(ctx, &req)
	if err != nil {
		logger.Errorf("get report data err: %v", err)
		ginutil.Error(ctx, errcode.ErrCommonReportGetFailed, errcode.MsgCommonReportFailGet)
		return
	}

	ginutil.Success(ctx, resp)
}

// JobCountReport
//
//	@Summary		作业投递情况(type=JOB_COUNT)
//	@Description	作业投递情况报表
//	@Tags			报表
//	@Accept			json
//	@Produce		json
//	@Param			param	body		dto.UniteReportReq	true	"请求参数"
//	@Response		200		{object}	dto.JobCountMetricResp
//	@Router			/report/jobCount [get]
func (r *apiRoute) JobCountReport(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	req := dto.UniteReportReq{}

	if err := ctx.BindQuery(&req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	resp, err := r.reportService.GetJobCount(ctx, &req)
	if err != nil {
		logger.Errorf("get report data err: %v", err)
		ginutil.Error(ctx, errcode.ErrCommonReportGetFailed, errcode.MsgCommonReportFailGet)
		return
	}

	ginutil.Success(ctx, resp)
}

// JobDeliverCountReport
//
//	@Summary		作业用户数和作业数情况(type=JOB_DELIVER_COUNT)
//	@Description	作业用户数和作业数情况
//	@Tags			报表
//	@Accept			json
//	@Produce		json
//	@Param			param	body		dto.UniteReportReq	true	"请求参数"
//	@Response		200		{object}	dto.JobDeliverCountResp
//	@Router			/report/jobDeliverCount [get]
func (r *apiRoute) JobDeliverCountReport(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	req := dto.UniteReportReq{}

	if err := ctx.BindQuery(&req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	resp, err := r.reportService.GetJobDeliverCount(ctx, &req)
	if err != nil {
		logger.Errorf("get report data err: %v", err)
		ginutil.Error(ctx, errcode.ErrCommonReportGetFailed, errcode.MsgCommonReportFailGet)
		return
	}

	ginutil.Success(ctx, resp)
}

// JobWaitStatisticReport
//
//	@Summary		作业等待情况(type=JOB_WAIT_STATISTIC)
//	@Description	作业等待情况
//	@Tags			报表
//	@Accept			json
//	@Produce		json
//	@Param			param	body		dto.UniteReportReq	true	"请求参数"
//	@Response		200		{object}	dto.JobWaitStatisticResp
//	@Router			/report/jobWaitStatistic [get]
func (r *apiRoute) JobWaitStatisticReport(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	req := dto.UniteReportReq{}

	if err := ctx.BindQuery(&req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	resp, err := r.reportService.GetJobWaitStatistic(ctx, &req)
	if err != nil {
		logger.Errorf("get report data err: %v", err)
		ginutil.Error(ctx, errcode.ErrCommonReportGetFailed, errcode.MsgCommonReportFailGet)
		return
	}

	ginutil.Success(ctx, resp)
}

// LicenseAppUsedAvgReport
//
//	@Summary		license app使用情况(type=LICENSE_APP_USED_UT_AVG)
//	@Description	license app使用情况
//	@Tags			报表
//	@Accept			json
//	@Produce		json
//	@Param			param	body		dto.UniteReportReq	true	"请求参数"
//	@Response		200		{object}	dto.LicenseAppUsedUtAvgReportResp
//	@Router			/report/licenseAppUsedUtAvg [get]
func (r *apiRoute) LicenseAppUsedAvgReport(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	req := dto.UniteReportReq{}

	if err := ctx.BindQuery(&req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	if req.Start != 0 && req.End != 0 && req.Start > req.End {
		logger.Errorf("report time rang invalid, end time before start time")
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	resp, err := r.reportService.GetLicenseAppUsedUtAvg(ctx, &req)
	if err != nil {
		logger.Errorf("get report data err: %v", err)
		ginutil.Error(ctx, errcode.ErrCommonReportGetFailed, errcode.MsgCommonReportFailGet)
		return
	}

	ginutil.Success(ctx, resp)
}

// LicenseAppModuleUsedUtAvgReport
//
//	@Summary		license app使用情况(type=LICENSE_APP_MODULE_USED_UT_AVG)
//	@Description	license app使用情况
//	@Tags			报表
//	@Accept			json
//	@Produce		json
//	@Param			param	body		dto.LicenseAppModuleUsedUtAvgReq	true	"请求参数"
//	@Response		200		{object}	dto.LicenseAppModuleUsedUtAvgReportResp
//	@Router			/report/licenseAppModuleUsedUtAvg [get]
func (r *apiRoute) LicenseAppModuleUsedUtAvgReport(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	req := dto.LicenseAppModuleUsedUtAvgReq{}

	if err := ctx.BindQuery(&req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	if req.Start != 0 && req.End != 0 && req.Start > req.End {
		logger.Errorf("report time rang invalid, end time before start time")
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	resp, err := r.reportService.GetLicenseAppModuleUsedUtAvg(ctx, &req)
	if err != nil {
		logger.Errorf("get report data err: %v", err)
		ginutil.Error(ctx, errcode.ErrCommonReportGetFailed, errcode.MsgCommonReportFailGet)
		return
	}

	ginutil.Success(ctx, resp)
}

// NodeDownStatisticReport
//
//	@Summary		节点下线统计(type=NODE_DOWN_STATISTIC)
//	@Description	节点下线统计接口
//	@Tags			报表
//	@Accept			json
//	@Produce		json
//	@Param			param	body		dto.NodeDownStatisticReportReq	true	"请求参数"
//	@Response		200		{object}	dto.NodeDownStatisticReportResp
//	@Router			/report/nodeDownStatistics [get]
func (r *apiRoute) NodeDownStatisticReport(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	req := dto.NodeDownStatisticReportReq{}

	if err := ctx.BindQuery(&req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	if req.Start == 0 || req.End == 0 || req.Start > req.End {
		logger.Errorf("report time rang invalid, end time before start time")
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	resp, err := r.reportService.NodeDownStatisticReport(ctx, &req)
	if err != nil {
		logger.Errorf("get report data err: %v", err)
		ginutil.Error(ctx, errcode.ErrCommonReportGetFailed, errcode.MsgCommonReportFailGet)
		return
	}

	ginutil.Success(ctx, resp)
}

// ExportNodeDownStatistics
//
//	@Summary		节点下线统计导出
//	@Description	节点下线统计导出接口
//	@Tags			报表
//	@Accept			json
//	@Produce		json
//	@Param			param	query	dto.ExportNodeDownStatisticsReq	true	"请求参数"
//	@Router			/report/export/nodeDownStatistics [get]
func (r *apiRoute) ExportNodeDownStatistics(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	req := dto.ExportNodeDownStatisticsReq{}

	if err := ctx.BindQuery(&req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	if req.Start == 0 || req.End == 0 || req.Start > req.End {
		logger.Errorf("report time rang invalid, end time before start time")
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	err := r.reportService.ExportNodeDownStatistics(ctx, &req)
	if err != nil {
		logger.Errorf("export report data err: %v", err)
		ginutil.Error(ctx, errcode.ErrCommonReportExportFailed, errcode.MsgCommonReportFailExport)
		return
	}
}

// ExportResourceUtAvg
//
//	@Summary		资源类型报表导出, 包括 cpu平均利用率(type=CPU_UT_AVG)， 内存利用率(type=MEM_UT_AVG)，磁盘吞吐率(type=TOTAL_IO_UT_AVG)
//	@Description	资源类型报表导出
//	@Tags			报表
//	@Accept			json
//	@Produce		json
//	@Param			param	query	dto.UniteReportReq	true	"请求参数"
//	@Router			/report/export/resourceUtAvg [get]
func (r *apiRoute) ExportResourceUtAvg(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	req := dto.UniteReportReq{}

	if err := ctx.BindQuery(&req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	if req.Start == 0 || req.End == 0 || req.Start > req.End {
		logger.Errorf("report time rang invalid, end time before start time")
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	err := r.reportService.ExportResourceUtAvg(ctx, &req)
	if err != nil {
		logger.Errorf("export report data err: %v", err)
		ginutil.Error(ctx, errcode.ErrCommonReportExportFailed, errcode.MsgCommonReportFailExport)
		return
	}
}

// ExportDiskUtAvg
//
//	@Summary		磁盘吞吐率报表导出
//	@Description	磁盘吞吐率报表导出接口
//	@Tags			报表
//	@Accept			json
//	@Produce		json
//	@Param			param	query	dto.UniteReportReq	true	"请求参数"
//	@Router			/report/export/diskUtAvg [get]
func (r *apiRoute) ExportDiskUtAvg(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	req := dto.UniteReportReq{}

	if err := ctx.BindQuery(&req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	if req.Start == 0 || req.End == 0 || req.Start > req.End {
		logger.Errorf("report time rang invalid, end time before start time")
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	err := r.reportService.ExportDiskUtAvg(ctx, &req)
	if err != nil {
		logger.Errorf("export report data err: %v", err)
		ginutil.Error(ctx, errcode.ErrCommonReportExportFailed, errcode.MsgCommonReportFailExport)
		return
	}
}

// ExportCPUTimeSum
//
//	@Summary		CPU平均利用率报表导出
//	@Description	CPU平均利用率报表导出接口
//	@Tags			报表
//	@Accept			json
//	@Produce		json
//	@Param			param	query	dto.UniteReportReq	true	"请求参数"
//	@Router			/report/export/cpuTimeSum [get]
func (r *apiRoute) ExportCPUTimeSum(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	req := dto.UniteReportReq{}

	if err := ctx.BindQuery(&req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	if req.Start == 0 || req.End == 0 || req.Start > req.End {
		logger.Errorf("report time rang invalid, end time before start time")
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	err := r.reportService.ExportCPUTimeSum(ctx, &req)
	if err != nil {
		logger.Errorf("export report data err: %v", err)
		ginutil.Error(ctx, errcode.ErrCommonReportExportFailed, errcode.MsgCommonReportFailExport)
		return
	}
}

// ExportJobCount
//
//	@Summary		作业投递情况报表导出
//	@Description	作业投递情况报表导出接口
//	@Tags			报表
//	@Accept			json
//	@Produce		json
//	@Param			param	query	dto.UniteReportReq	true	"请求参数"
//	@Router			/report/export/jobCount [get]
func (r *apiRoute) ExportJobCount(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	req := dto.UniteReportReq{}

	if err := ctx.BindQuery(&req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	if req.Start == 0 || req.End == 0 || req.Start > req.End {
		logger.Errorf("report time rang invalid, end time before start time")
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	err := r.reportService.ExportJobCount(ctx, &req)
	if err != nil {
		logger.Errorf("export report data err: %v", err)
		ginutil.Error(ctx, errcode.ErrCommonReportExportFailed, errcode.MsgCommonReportFailExport)
		return
	}
}

// ExportJobDeliverCount
//
//	@Summary		用户数和作业投递数量报表导出
//	@Description	用户数和作业投递数量报表导出接口
//	@Tags			报表
//	@Accept			json
//	@Produce		json
//	@Param			param	query	dto.UniteReportReq	true	"请求参数"
//	@Router			/report/export/jobDeliverCount [get]
func (r *apiRoute) ExportJobDeliverCount(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	req := dto.UniteReportReq{}

	if err := ctx.BindQuery(&req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	if req.Start == 0 || req.End == 0 || req.Start > req.End {
		logger.Errorf("report time rang invalid, end time before start time")
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	err := r.reportService.ExportJobDeliverCount(ctx, &req)
	if err != nil {
		logger.Errorf("export report data err: %v", err)
		ginutil.Error(ctx, errcode.ErrCommonReportExportFailed, errcode.MsgCommonReportFailExport)
		return
	}
}

// ExportJobWaitStatistic
//
//	@Summary		作业等待时间统计报表导出
//	@Description	作业等待时间统计报表导出接口
//	@Tags			报表
//	@Accept			json
//	@Produce		json
//	@Param			param	query	dto.UniteReportReq	true	"请求参数"
//	@Router			/report/export/jobWaitStatistic [get]
func (r *apiRoute) ExportJobWaitStatistic(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	req := dto.UniteReportReq{}

	if err := ctx.BindQuery(&req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	if req.Start == 0 || req.End == 0 || req.Start > req.End {
		logger.Errorf("report time rang invalid, end time before start time")
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	err := r.reportService.ExportJobWaitStatistic(ctx, &req)
	if err != nil {
		logger.Errorf("export report data err: %v", err)
		ginutil.Error(ctx, errcode.ErrCommonReportExportFailed, errcode.MsgCommonReportFailExport)
		return
	}
}

// ExportLicenseAppUsedUtAvg
//
//	@Summary		license app使用情况报表导出
//	@Description	license app使用情况报表导出接口
//	@Tags			报表
//	@Accept			json
//	@Produce		json
//	@Param			param	query	dto.UniteReportReq	true	"请求参数"
//	@Router			/report/export/licenseAppUsedUtAvg [get]
func (r *apiRoute) ExportLicenseAppUsedUtAvg(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	req := dto.UniteReportReq{}

	if err := ctx.BindQuery(&req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	if req.Start == 0 || req.End == 0 || req.Start > req.End {
		logger.Errorf("report time rang invalid, end time before start time")
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	err := r.reportService.ExportLicenseAppUsedUtAvg(ctx, &req)
	if err != nil {
		logger.Errorf("export report data err: %v", err)
		ginutil.Error(ctx, errcode.ErrCommonReportExportFailed, errcode.MsgCommonReportFailExport)
		return
	}
}

// ReportExport
//
//	@Summary		报表导出,具体type类型同查询
//	@Description	报表导出
//	@Tags			报表
//	@Accept			json
//	@Produce		json
//	@Param			param	body	dto.UniteReportReq	true	"请求参数"
//	@Router			/report/export [get]
func (r *apiRoute) ReportExport(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	req := dto.UniteReportReq{}

	if err := ctx.BindQuery(&req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	if util.ValidReportType(req.ReportType) {
		logger.Errorf("report type invalid")
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	if req.Start != 0 && req.End != 0 && req.Start > req.End {
		logger.Errorf("report time rang invalid, end time before start time")
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	err := r.reportService.ExportReport(ctx, &req)
	if err != nil {
		logger.Errorf("get report data err: %v", err)
		ginutil.Error(ctx, errcode.ErrCommonReportExportFailed, errcode.MsgCommonReportFailExport)
		return
	}
}
