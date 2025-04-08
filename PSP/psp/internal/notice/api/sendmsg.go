package api

import (
	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/internal/notice/dto"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/ginutil"
)

// SendWebsocketMessage
//
//	@Summary		发送websocket消息
//	@Description	发送websocket消息接口
//	@Tags			通知-消息管理
//	@Accept			json
//	@Produce		json
//	@Param			param	body	dto.WebsocketMessage	true	"请求参数"
//	@Response		200		"正常返回表示成功，否则返回错误码和错误信息"
//	@Router			/notice/producer [post]
func (r *apiRoute) SendWebsocketMessage(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	var req = &dto.WebsocketMessage{}
	if err := ctx.ShouldBindJSON(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	if err := r.messageService.SendWebsocketMessage(ctx, req); err != nil {
		logger.Errorf("send websocket message err: %v", err)
		ginutil.Error(ctx, errcode.ErrNoticeFailSend, errcode.MsgNoticeFailSend)
		return
	}

	ginutil.Success(ctx, nil)
}
