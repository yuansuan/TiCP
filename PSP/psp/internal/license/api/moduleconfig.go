package api

import (
	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"google.golang.org/grpc/status"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/internal/license/dto"
	"github.com/yuansuan/ticp/PSP/psp/internal/license/util"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/ginutil"
)

// ListModuleConfig
//
//	@Summary		模块配置列表
//	@Description	模块配置列表接口
//	@Tags			许可证监控
//	@Accept			json
//	@Produce		json
//	@Response		200	{object}	dto.ModuleConfigListResponse
//	@Router			/licenseInfos/:id/moduleConfigs [get]
func (r *apiRoute) ListModuleConfig(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	id, ok := util.GetResourceId(ctx)
	if !ok {
		return
	}

	resp, err := r.moduleConfigService.ModuleConfigList(ctx, id)
	if err != nil {
		logger.Errorf("get module config list err: %v", err)
		ginutil.Error(ctx, errcode.ErrFailedConfigModuleList, errcode.MsgFailedConfigModuleList)
		return
	}

	ginutil.Success(ctx, resp)
}

// AddModuleConfig
//
//	@Summary		新增模块配置
//	@Description	新增模块配置接口
//	@Tags			许可证监控
//	@Accept			json
//	@Produce		json
//	@Param			param	query		dto.AddModuleConfigRequest	true	"请求参数"
//	@Response		200		{object}	dto.AddModuleConfigResponse
//	@Router			/moduleConfigs [post]
func (r *apiRoute) AddModuleConfig(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	var req = &dto.AddModuleConfigRequest{}
	if err := ctx.BindJSON(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	resp, err := r.moduleConfigService.AddModuleConfig(ctx, req)
	if err != nil {
		logger.Errorf("add module config err: %v", err)
		if status.Code(err) == errcode.ErrFailedModuleNameRepeat {
			ginutil.Error(ctx, errcode.ErrFailedModuleNameRepeat, errcode.MsgFailedModuleNameRepeat)
		} else {
			ginutil.Error(ctx, errcode.ErrFailedConfigModuleAdd, errcode.MsgFailedConfigModuleAdd)
		}

		return
	}

	ginutil.Success(ctx, resp)
}

// EditModuleConfig
//
//	@Summary		编辑模块配置
//	@Description	编辑模块配置接口
//	@Tags			许可证监控
//	@Accept			json
//	@Produce		json
//	@Param			param	query	dto.EditModuleConfigRequest	true	"请求参数"
//	@Router			/moduleConfigs/:id [put]
func (r *apiRoute) EditModuleConfig(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	var req = &dto.EditModuleConfigRequest{}
	if err := ctx.BindJSON(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	err := r.moduleConfigService.EditModuleConfig(ctx, req)
	if err != nil {
		logger.Errorf("edit module config err: %v", err)
		if status.Code(err) == errcode.ErrFailedModuleNameRepeat {
			ginutil.Error(ctx, errcode.ErrFailedModuleNameRepeat, errcode.MsgFailedModuleNameRepeat)
		} else {
			ginutil.Error(ctx, errcode.ErrFailedConfigModuleEdit, errcode.MsgFailedConfigModuleEdit)
		}
		return
	}

	ginutil.Success(ctx, nil)
}

// DeleteModuleConfig
//
//	@Summary		删除模块配置
//	@Description	删除模块配置接口
//	@Tags			许可证监控
//	@Accept			json
//	@Produce		json
//	@Router			/moduleConfigs/:id [delete]
func (r *apiRoute) DeleteModuleConfig(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	id, ok := util.GetResourceId(ctx)
	if !ok {
		return
	}

	err := r.moduleConfigService.DeleteModuleConfig(ctx, id)
	if err != nil {
		logger.Errorf("delete module config err: %v", err)
		ginutil.Error(ctx, errcode.ErrFailedConfigModuleDelete, errcode.MsgFailedConfigModuleDelete)
		return
	}

	ginutil.Success(ctx, nil)
}
