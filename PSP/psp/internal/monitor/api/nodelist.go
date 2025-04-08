package api

import (
	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/internal/monitor/dto"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/ginutil"
)

// NodeList
//
//	@Summary		节点列表
//	@Description	获取集群中的节点列表，支持按照节点名称模糊查询及分页查询
//	@Tags			节点管理
//	@Accept			json
//	@Produce		json
//	@Param			param	body		dto.NodeListRequest	true	"请求参数"
//	@Response		200		{object}	dto.NodeListResponse
//	@Router			/node/list [get]
func (r *apiRoute) NodeList(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	var req = dto.NodeListRequest{}
	if err := ctx.BindQuery(&req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	if req.PageIndex < 1 {
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	if req.PageSize < 1 || req.PageSize > 1000 {
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	list, total, err := r.nodeService.GetNodeList(ctx, req.NodeName, req.PageIndex, req.PageSize)
	if err != nil {
		logger.Errorf("get node list err: %v", err)
		ginutil.Error(ctx, errcode.ErrNodeListFail, errcode.MsgNodeListFail)
		return
	}

	ginutil.Success(ctx, dto.NodeListResponse{
		Total:        total,
		NodeInfoList: list,
	})
}
