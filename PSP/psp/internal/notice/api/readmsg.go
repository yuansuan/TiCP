package api

import (
	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/internal/notice/dto"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/ginutil"
)

// ReadMessage
//
//	@Summary		消息设置已读
//	@Description	消息设置已读接口
//	@Tags			通知-消息管理
//	@Accept			json
//	@Produce		json
//	@Param			param	body	dto.ReadMessageRequest	true	"请求参数"
//	@Response		200		"正常返回表示成功，否则返回错误码和错误信息"
//	@Router			/notice/read [put]
func (r *apiRoute) ReadMessage(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	var req = &dto.ReadMessageRequest{}
	if err := ctx.ShouldBindJSON(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	userID := ginutil.GetUserID(ctx)
	if err := r.messageService.ReadMessage(ctx, userID, req.MessageIDs); err != nil {
		logger.Errorf("send websocket message err: %v", err)
		ginutil.Error(ctx, errcode.ErrNoticeFailRead, errcode.MsgNoticeFailRead)
		return
	}

	ginutil.Success(ctx, nil)
}

// ReadAllMessage
//
//	@Summary		消息全部设置已读
//	@Description	消息全部设置已读接口
//	@Tags			通知-消息管理
//	@Accept			json
//	@Produce		json
//	@Response		200	"正常返回表示成功，否则返回错误码和错误信息"
//	@Router			/notice/readAll [put]
func (r *apiRoute) ReadAllMessage(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	userID := ginutil.GetUserID(ctx)
	if err := r.messageService.ReadAllMessage(ctx, userID); err != nil {
		logger.Errorf("send websocket message err: %v", err)
		ginutil.Error(ctx, errcode.ErrNoticeFailRead, errcode.MsgNoticeFailRead)
		return
	}

	ginutil.Success(ctx, nil)
}
