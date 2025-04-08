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

// GetThreePersonManagementConfig
//
//	@Summary		获取三员管理设置
//	@Description	获取三员管理设置
//	@Tags			系统配置
//	@Accept			json
//	@Produce		json
//	@Response		200	{object}	dto.GetThreePersonConfigResponse
//	@Router			/sysconfig/getThreePersonManagementConfig [get]
func (srv *RouteService) GetThreePersonManagementConfig(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	emailConfig, err := srv.SysConfigService.GetThreePersonManagementConfig(ctx)
	if err != nil {
		logger.Errorf("get threePersonManagementConfig config err: %v, result: [%+v]", err, emailConfig)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrSysConfigGetThreePersonFailed)
		return
	}
	ginutil.Success(ctx, &emailConfig)
}

// SetThreePersonManagementConfig
//
//	@Summary		获取三员管理设置
//	@Description	获取三员管理设置
//	@Tags			系统配置
//	@Accept			json
//	@Produce		json
//	@Param			param	body		dto.SetThreePersonConfigRequest	true	"请求参数"
//	@Response		200		{object}	dto.SetThreePersonConfigResponse
//	@Router			/sysconfig/setThreePersonManagementConfig [post]
func (srv *RouteService) SetThreePersonManagementConfig(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	req := &dto.SetThreePersonConfigRequest{}
	if err := ctx.ShouldBindJSON(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	tracelog.Info(ctx, fmt.Sprintf("user: [%v] set three person manager config req: [%+v]", ginutil.GetUserID(ctx), req))

	err := srv.SysConfigService.SetThreePersonManagementConfig(ctx, req)
	if err != nil {
		logger.Errorf("set three person manager config err: %v, req: [%+v]", err, req)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrSysConfigSetThreePersonFailed)
		return
	}

	ginutil.Success(ctx, &dto.SetThreePersonConfigResponse{})
}
