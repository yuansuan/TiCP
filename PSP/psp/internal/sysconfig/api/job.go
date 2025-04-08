package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/internal/sysconfig/dto"
	"github.com/yuansuan/ticp/PSP/psp/pkg/tracelog"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/ginutil"
)

// GetJobConfig
//
//	@Summary		获取作业配置
//	@Description	获取作业配置接口
//	@Tags			系统配置
//	@Accept			json
//	@Produce		json
//	@Param			param	query		dto.GetJobConfigRequest	true	"请求参数"
//	@Response		200		{object}	dto.GetJobConfigResponse
//	@Router			/sysconfig/getJobConfig [get]
func (s *RouteService) GetJobConfig(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	req := &dto.GetJobConfigRequest{}
	if err := ctx.ShouldBindQuery(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	jobConfig, err := s.SysConfigService.GetJobConfig(ctx)
	if err != nil {
		logger.Errorf("get job burst config err: %v, req: [%+v]", err, req)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrSysConfigGetJobFailed)
		return
	}

	ginutil.Success(ctx, jobConfig)
}

// SetJobConfig
//
//	@Summary		设置作业配置
//	@Description	设置作业配置接口
//	@Tags			系统配置
//	@Accept			json
//	@Produce		json
//	@Param			param	body		dto.SetJobConfigRequest	true	"请求参数"
//	@Response		200		{object}	dto.SetJobConfigResponse
//	@Router			/sysconfig/setJobConfig [post]
func (s *RouteService) SetJobConfig(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	req := &dto.SetJobConfigRequest{}
	if err := ctx.ShouldBindJSON(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	tracelog.Info(ctx, fmt.Sprintf("user: [%v] set job config req: [%+v]", ginutil.GetUserID(ctx), req))

	err := s.SysConfigService.SetJobConfig(ctx, req)
	if err != nil {
		logger.Errorf("set job burst config err: %v, req: [%+v]", err, req)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrSysConfigSetJobFailed)
		return
	}

	ginutil.Success(ctx, &dto.SetJobConfigResponse{})
}

// GetJobBurstConfig
//
//	@Summary		获取作业爆发配置
//	@Description	获取作业爆发配置接口
//	@Tags			系统配置
//	@Accept			json
//	@Produce		json
//	@Param			param	query		dto.GetJobBurstConfigRequest	true	"请求参数"
//	@Response		200		{object}	dto.GetJobBurstConfigResponse
//	@Router			/sysconfig/getJobBurstConfig [get]
func (s *RouteService) GetJobBurstConfig(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	req := &dto.GetJobBurstConfigRequest{}
	if err := ctx.ShouldBindQuery(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	jobBurstConfig, err := s.SysConfigService.GetJobBurstConfig(ctx)
	if err != nil {
		logger.Errorf("get job burst config err: %v, req: [%+v]", err, req)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrSysConfigGetJobBurstFailed)
		return
	}

	ginutil.Success(ctx, jobBurstConfig)
}

// SetJobBurstConfig
//
//	@Summary		设置作业爆发配置
//	@Description	设置作业爆发配置接口
//	@Tags			系统配置
//	@Accept			json
//	@Produce		json
//	@Param			param	body		dto.SetJobBurstConfigRequest	true	"请求参数"
//	@Response		200		{object}	dto.SetJobBurstConfigResponse
//	@Router			/sysconfig/setJobBurstConfig [post]
func (s *RouteService) SetJobBurstConfig(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	req := &dto.SetJobBurstConfigRequest{}
	if err := ctx.ShouldBindJSON(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	tracelog.Info(ctx, fmt.Sprintf("user: [%v] set job burst config req: [%+v]", ginutil.GetUserID(ctx), req))

	err := s.SysConfigService.SetJobBurstConfig(ctx, req)
	if err != nil {
		logger.Errorf("set job burst config err: %v, req: [%+v]", err, req)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrSysConfigSetJobBurstFailed)
		return
	}

	ginutil.Success(ctx, &dto.SetJobBurstConfigResponse{})
}
