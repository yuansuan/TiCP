package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/oplog"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/approve"
	"github.com/yuansuan/ticp/PSP/psp/internal/monitor/consts"
	"github.com/yuansuan/ticp/PSP/psp/internal/monitor/dto"
	"github.com/yuansuan/ticp/PSP/psp/pkg/tracelog"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/ginutil"
)

// NodeOperate
//
//	@Summary		节点操作
//	@Description	支持 接受作业/拒绝作业操作
//	@Tags			节点管理
//	@Accept			json
//	@Produce		json
//	@Param			param	body	dto.NodeOperateRequest	true	"请求参数"
//	@Router			/node/operate [post]
func (r *apiRoute) NodeOperate(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	var req = dto.NodeOperateRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}
	if len(req.NodeNames) == 0 {
		logger.Errorf("nodeNames is empty")
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}
	if req.Operation != consts.NodeClose && req.Operation != consts.NodeStart {
		logger.Errorf("operation is invalid")
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	tracelog.Info(ctx, fmt.Sprintf("update nodeStatus success, updateUserID:[%v], nodeNames:[%v], operation:[%v]", ginutil.GetUserID(ctx), req.NodeNames, req.Operation))

	err := r.nodeService.NodeOperate(ctx, req.NodeNames, req.Operation)
	if err != nil {
		logger.Errorf("node operate err: %v", err)
		ginutil.Error(ctx, errcode.ErrNodeFailOperate, errcode.MsgNodeFailOperate)
		return
	}

	operation := "接受作业"
	if req.Operation == "node_close" {
		operation = "拒绝作业"
	}
	oplog.GetInstance().SaveAuditLogInfo(ctx, approve.OperateTypeEnum_NODE_MANAGER, fmt.Sprintf("用户%v设置%v%v", ginutil.GetUserName(ctx), req.NodeNames, operation))

	ginutil.Success(ctx, "success")
}
