package openapi

import (
	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/yuansuan/ticp/PSP/psp/internal/app/dto"
	"github.com/yuansuan/ticp/PSP/psp/internal/app/dto/openapi"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/ginutil"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/structutil"
)

// ListApp
//
//	@Summary		计算应用列表-openapi
//	@Description	计算应用列表接口
//	@Tags			计算应用
//	@Accept			json
//	@Produce		json
//	@Param			param	query		openapi.ListAppRequest	true	"请求参数"
//	@Response		200		{object}	openapi.ListAppResponse
//	@Router			/openapi/app/list [get]
func (s *RouteOpenapiService) ListApp(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	req := &openapi.ListAppRequest{}
	if err := ctx.BindQuery(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	if err := s.validate.Struct(req); err != nil {
		ginutil.Error(ctx, errcode.ErrInvalidParam, err.Error())
		return
	}

	userID := ginutil.GetUserID(ctx)
	templates, err := s.appService.ListApp(ctx, userID, req.ComputeType, "published", true, false)
	if err != nil {
		logger.Errorf("list app err: %v, req: [%+v]", err, req)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrAppListAppFailed)
		return
	}

	rsp := &openapi.ListAppResponse{}
	if err = structutil.CopyStruct(rsp, &dto.ListAppResponse{Apps: templates}); err != nil {
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrAppListAppFailed)
		return
	}

	ginutil.Success(ctx, rsp)
}
