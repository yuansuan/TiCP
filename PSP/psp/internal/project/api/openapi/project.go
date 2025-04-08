package openapi

import (
	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/internal/project/dto"
	"github.com/yuansuan/ticp/PSP/psp/internal/project/dto/openapi"
	"github.com/yuansuan/ticp/PSP/psp/internal/project/util"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/ginutil"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/structutil"
)

// ProjectList
//
//	@Summary		项目列表
//	@Description	项目列表接口
//	@Tags			项目
//	@Accept			json
//	@Produce		json
//	@Param			param	query		dto.ProjectListRequest	true	"请求参数"
//	@Response		200		{object}	dto.ProjectListResponse
//	@Router			/openapi/project/list [post]
func (r *RouteOpenapiService) ProjectList(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	req := openapi.ProjectListRequest{}
	if err := ctx.BindJSON(&req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	loginUserID := ginutil.GetUserID(ctx)

	resp, err := r.projectService.ProjectList(ctx, &dto.ProjectListRequest{
		Page:        req.Page,
		ProjectName: req.ProjectName,
		State:       req.State,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		IsSysMenu:   false,
	}, snowflake.ID(loginUserID))
	if err != nil {
		logger.Errorf("List project err: %v", err)
		util.CheckErrorIfNoPermission(ctx, err, errcode.ErrProjectList)
		return
	}

	rsp := &openapi.ProjectListResponse{}
	if err = structutil.CopyStruct(rsp, resp); err != nil {
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrVisualListSoftwareFailed)
		return
	}

	ginutil.Success(ctx, rsp)
}
