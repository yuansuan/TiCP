package api

import (
	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/internal/monitor/dto"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/ginutil"
)

// NodeDetail
//
//	@Summary		节点详情
//	@Description	根据节点名称获取某个节点的详细
//	@Tags			节点管理
//	@Accept			json
//	@Produce		json
//	@Param			param	body		dto.NodeRequest	true	"请求参数"
//	@Response		200		{object}	dto.NodeDetailResponse
//	@Router			/node/detail [get]
func (r *apiRoute) NodeDetail(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	var req = dto.NodeRequest{}
	if err := ctx.BindQuery(&req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	if req.NodeName == "" {
		logger.Errorf("request params invalid: %v", req)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}
	nodeDetail, err := r.nodeService.GetNodeInfo(ctx, req.NodeName)
	if err != nil {
		logger.Errorf("node detail err: %v", err)
		ginutil.Error(ctx, errcode.ErrNodeFailGet, errcode.MsgNodeFailGet)
		return
	}

	ginutil.Success(ctx, dto.NodeDetailResponse{
		NodeDetail: nodeDetail,
	})
}
