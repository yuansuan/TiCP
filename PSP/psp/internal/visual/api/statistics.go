package api

import (
	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/yuansuan/ticp/PSP/psp/internal/common"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/internal/visual/dto"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/ginutil"
)

// DurationStatistic
//
//	@Summary		会话时长统计
//	@Description	会话时长统计接口
//	@Tags			可视化-统计
//	@Accept			json
//	@Produce		json
//	@Param			param	query		dto.DurationStatisticRequest	true	"请求参数"
//	@Response		200		{object}	dto.DurationStatisticResponse
//	@Router			/vis/statistic/duration [get]
func (s *RouteService) DurationStatistic(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	req := &dto.DurationStatisticRequest{}
	if err := ctx.BindQuery(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	if req.StartTime == "" || req.EndTime == "" {
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	statistics, err := s.visualService.DurationStatistic(ctx, req.AppIDs, req.StartTime, req.EndTime)
	if err != nil {
		logger.Errorf("duration statistic err: %v, req: [%+v]", err, req)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrVisualDurationStatisticFailed)
		return
	}
	ginutil.Success(ctx, &dto.DurationStatisticResponse{Statistics: statistics})
}

// ListHistoryDuration
//
//	@Summary		会话时长列表
//	@Description	会话时长列表接口
//	@Tags			可视化-统计
//	@Accept			json
//	@Produce		json
//	@Param			param	query		dto.ListHistoryDurationRequest	true	"请求参数"
//	@Response		200		{object}	dto.ListHistoryDurationResponse
//	@Router			/vis/statistic/duration/list [get]
func (s *RouteService) ListHistoryDuration(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	req := &dto.ListHistoryDurationRequest{}
	if err := ctx.BindQuery(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	if req.StartTime == "" || req.EndTime == "" {
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	if req.PageSize == 0 {
		req.PageSize = common.DefaultMaxPageSize
	}

	statistics, total, err := s.visualService.ListHistoryDuration(ctx, req.AppIDs, req.StartTime, req.EndTime,
		req.PageIndex, req.PageSize)
	if err != nil {
		logger.Errorf("list history duration err: %v, req: [%+v]", err, req)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrVisualListHistoryDurationFailed)
		return
	}
	ginutil.Success(ctx, &dto.ListHistoryDurationResponse{Statistics: statistics, Total: total})
}

// SessionUsageDurationStatistic
//
//	@Summary		会话时长统计
//	@Description	会话时长统计接口
//	@Tags			可视化-统计
//	@Accept			json
//	@Produce		json
//	@Param			param	query		dto.SessionUsageDurationStatisticRequest	true	"请求参数"
//	@Response		200		{object}	dto.SessionUsageDurationStatisticResponse
//	@Router			/vis/statistic/report/duration [get]
func (s *RouteService) SessionUsageDurationStatistic(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	req := &dto.SessionUsageDurationStatisticRequest{}
	if err := ctx.BindQuery(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	if req.Start == 0 || req.End == 0 {
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	statistics, err := s.visualService.SessionUsageDurationStatistic(ctx, req.Start, req.End)
	if err != nil {
		logger.Errorf("session usage duration statistic err: %v, req: [%+v]", err, req)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrVisualSessionUsageDurationStatisticFailed)
		return
	}

	ginutil.Success(ctx, statistics)
}

// ExportUsageDurationStatistic
//
//	@Summary		导出会话时长统计
//	@Description	导出会话时长统计接口
//	@Tags			可视化-统计
//	@Accept			json
//	@Produce		json
//	@Param			param	query	dto.ExportUsageDurationStatisticRequest	true	"请求参数"
//	@Response		200
//	@Router			/vis/statistic/report/duration/export [get]
func (s *RouteService) ExportUsageDurationStatistic(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	req := &dto.ExportUsageDurationStatisticRequest{}
	if err := ctx.BindQuery(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	if req.Start == 0 || req.End == 0 {
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	err := s.visualService.ExportUsageDurationStatistic(ctx, req.Start, req.End)
	if err != nil {
		logger.Errorf("export usage duration statistic err: %v, req: [%+v]", err, req)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrVisualExportUsageDurationStatisticFailed)
		return
	}
}

// SessionCreateNumberStatistic
//
//	@Summary		会话创建数量统计
//	@Description	会话创建数量统计接口
//	@Tags			可视化-统计
//	@Accept			json
//	@Produce		json
//	@Param			param	query		dto.SessionCreateNumberStatisticRequest	true	"请求参数"
//	@Response		200		{object}	dto.SessionCreateNumberStatisticResponse
//	@Router			/vis/statistic/report/createNumber [get]
func (s *RouteService) SessionCreateNumberStatistic(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	req := &dto.SessionCreateNumberStatisticRequest{}
	if err := ctx.BindQuery(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	if req.Start == 0 || req.End == 0 {
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	statistics, err := s.visualService.SessionCreateNumberStatistic(ctx, req.Start, req.End)
	if err != nil {
		logger.Errorf("session create number statistic err: %v, req: [%+v]", err, req)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrVisualSessionCreateNumberStatisticFailed)
		return
	}

	ginutil.Success(ctx, statistics)
}

// SessionNumberStatusStatistic
//
//	@Summary		会话数量状态统计
//	@Description	会话数量状态统计接口
//	@Tags			可视化-统计
//	@Accept			json
//	@Produce		json
//	@Param			param	query		dto.SessionNumberStatusStatisticRequest	true	"请求参数"
//	@Response		200		{object}	dto.SessionNumberStatusStatisticResponse
//	@Router			/vis/statistic/report/numberStatus [get]
func (s *RouteService) SessionNumberStatusStatistic(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	req := &dto.SessionNumberStatusStatisticRequest{}
	if err := ctx.BindQuery(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	if req.Start == 0 || req.End == 0 {
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	statistics, err := s.visualService.SessionNumberStatusStatistic(ctx, req.Start, req.End)
	if err != nil {
		logger.Errorf("session number status statistic err: %v, req: [%+v]", err, req)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrVisualSessionNumberStatusStatisticFailed)
		return
	}

	ginutil.Success(ctx, statistics)
}
