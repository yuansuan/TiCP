package api

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/yuansuan/ticp/PSP/psp/internal/common"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/oplog"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/approve"
	"github.com/yuansuan/ticp/PSP/psp/internal/job/consts"
	"github.com/yuansuan/ticp/PSP/psp/internal/job/dto"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/floatutil"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/ginutil"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/timeutil"
)

// GetJobStatisticsTotalCPUTime
//
//	@Summary		获取作业统计总核时
//	@Description	获取作业统计总核时接口
//	@Tags			作业
//	@Accept			json
//	@Produce		json
//	@Param			param	query		dto.GetJobStatisticsTotalCPUTimeRequest	true	"请求参数"
//	@Response		200		{object}	dto.GetJobStatisticsTotalCPUTimeResponse
//	@Router			/job/statistics/totalCPUTime [get]
func (r *apiRoute) GetJobStatisticsTotalCPUTime(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	req := &dto.GetJobStatisticsTotalCPUTimeRequest{}
	if err := ctx.BindQuery(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	cpuTimeTotal, err := r.jobService.GetJobCPUTimeTotal(ctx, req.QueryType, req.ComputeType, req.Names, req.ProjectIDs, req.StartTime, req.EndTime)
	if err != nil {
		logger.Errorf("get job statistics total cpu time err: %v, req: [%+v]", err, req)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrJobFailCPUTimeTotal)
		return
	}

	ginutil.Success(ctx, &dto.GetJobStatisticsTotalCPUTimeResponse{
		CPUTime: floatutil.NumberToFloatStr(cpuTimeTotal, common.DecimalPlaces),
	})
}

// GetJobStatisticsOverview
//
//	@Summary		获取作业统计总览
//	@Description	获取作业统计总览接口
//	@Tags			作业
//	@Accept			json
//	@Produce		json
//	@Param			param	query		dto.GetJobStatisticsOverviewRequest	true	"请求参数"
//	@Response		200		{object}	dto.GetJobStatisticsOverviewResponse
//	@Router			/job/statistics/overview [get]
func (r *apiRoute) GetJobStatisticsOverview(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	req := &dto.GetJobStatisticsOverviewRequest{}
	if err := ctx.BindQuery(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	if req.PageSize == 0 {
		req.PageSize = common.DefaultMaxPageSize
	}

	overviews, total, err := r.jobService.GetJobStatisticsOverview(ctx, req.QueryType, req.ComputeType, req.Names, req.ProjectIDs, req.StartTime, req.EndTime, req.PageIndex, req.PageSize)
	if err != nil {
		logger.Errorf("get job statistics overview err: %v, req: [%+v]", err, req)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrJobFailStatisticsOverview)
		return
	}

	ginutil.Success(ctx, &dto.GetJobStatisticsOverviewResponse{Overviews: overviews, Total: total})
}

// GetJobStatisticsDetail
//
//	@Summary		获取作业统计明细
//	@Description	获取作业统计明细接口
//	@Tags			作业
//	@Accept			json
//	@Produce		json
//	@Param			param	query		dto.GetJobStatisticsDetailRequest	true	"请求参数"
//	@Response		200		{object}	dto.GetJobStatisticsDetailResponse
//	@Router			/job/statistics/detail [get]
func (r *apiRoute) GetJobStatisticsDetail(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	req := &dto.GetJobStatisticsDetailRequest{}
	if err := ctx.BindQuery(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	if req.PageSize == 0 {
		req.PageSize = common.DefaultMaxPageSize
	}

	jobDetails, total, err := r.jobService.GetJobStatisticsDetail(ctx, req.QueryType, req.ComputeType, req.Names, req.ProjectIDs, req.StartTime, req.EndTime, req.PageIndex, req.PageSize)
	if err != nil {
		logger.Errorf("get job statistics detail err: %v, req: [%+v]", err, req)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrJobFailStatisticsDetail)
		return
	}

	ginutil.Success(ctx, &dto.GetJobStatisticsDetailResponse{JobDetails: jobDetails, Total: total})
}

// GetJobStatisticsExport
//
//	@Summary		导出统计数据
//	@Description	获取作业统计明细接口
//	@Tags			作业
//	@Accept			json
//	@Produce		json
//	@Param			param	query	dto.GetJobStatisticsExportRequest	true	"请求参数"
//	@Response		200
//	@Router			/job/statistics/export [get]
func (r *apiRoute) GetJobStatisticsExport(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	req := &dto.GetJobStatisticsExportRequest{}
	if err := ctx.BindQuery(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	err := r.jobService.GetJobStatisticsExport(ctx, req.QueryType, req.ComputeType, req.ShowType, req.Names, req.ProjectIDs, req.StartTime, req.EndTime)
	if err != nil {
		logger.Errorf("get job statistics detail err: %v, req: [%+v]", err, req)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrJobFailStatisticsExport)
		return
	}

	var obj string = consts.JobQueryUser
	var showTypeView string = consts.JobDetail

	if req.ShowType == consts.JobStatisticsShowTypeOverview {
		showTypeView = consts.JobOverview
	} else if req.ShowType == consts.JobStatisticsShowTypeDetail {
		showTypeView = consts.JobDetail
	}

	if req.QueryType == consts.JobStatisticsQueryTypeUser {
		obj = consts.JobQueryUser
	} else if req.QueryType == consts.JobStatisticsQueryTypeApp {
		obj = consts.JobQueryApp
	}

	oplog.GetInstance().SaveAuditLogInfo(ctx, approve.OperateTypeEnum_JOB_MANAGER, fmt.Sprintf("用户%v导出【%v】-【%v】界面且提交时间范围为：%s-%s的作业统计数据", ginutil.GetUserName(ctx), obj, showTypeView, timeutil.FormatTime(time.Unix(req.StartTime, 0), common.DatetimeFormat), timeutil.FormatTime(time.Unix(req.EndTime, 0), common.DatetimeFormat)))

}

// GetTop5ProjectInfo
//
//	@Summary		获取作业核时统计排名前5的项目信息
//	@Description	获取作业核时统计排名前5的项目信息接口
//	@Tags			作业
//	@Accept			json
//	@Produce		json
//	@Param			param	query		dto.GetTop5ProjectInfoRequest	true	"请求参数"
//	@Response		200		{object}	dto.GetTop5ProjectInfoResponse
//	@Router			/job/statistics/top5ProjectInfo [get]
func (r *apiRoute) GetTop5ProjectInfo(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	req := &dto.GetTop5ProjectInfoRequest{}
	if err := ctx.BindQuery(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	if req.Start == 0 || req.End == 0 || (req.Start > req.End) {
		logger.Errorf("reqeust params invalid, start: %d, end: %d", req.Start, req.End)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	res, err := r.jobService.GetTop5ProjectInfo(ctx, req.Start, req.End)
	if err != nil {
		logger.Errorf("get top5 project info err: %v, req: [%+v]", err, req)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrJobFailGetTop5ProjectInfo)
		return
	}

	ginutil.Success(ctx, res)
}
