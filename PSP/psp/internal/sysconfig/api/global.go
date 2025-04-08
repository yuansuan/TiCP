package api

import (
	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/internal/sysconfig/dto"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/ginutil"
)

// GetGlobalSysConfig
//
//	@Summary		获取全局系统配置
//	@Description	获取全局系统配置接口
//	@Tags			系统配置
//	@Accept			json
//	@Produce		json
//	@Param			param	query		dto.GetGlobalSysConfigRequest	true	"请求参数"
//	@Response		200		{object}	dto.GetGlobalSysConfigResponse
//	@Router			/sysconfig/global [get]
func (s *RouteService) GetGlobalSysConfig(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	req := &dto.GetGlobalSysConfigRequest{}
	if err := ctx.BindQuery(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	data, err := s.SysConfigService.GetGlobalSysConfig(ctx)
	if err != nil {
		logger.Errorf("get global sys config err: %v, req: [%+v]", err, req)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrSysConfigGetGlobalFailed)
		return
	}

	ginutil.Success(ctx, data)
}
