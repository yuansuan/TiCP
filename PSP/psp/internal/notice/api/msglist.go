package api

import (
	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/internal/notice/dto"
	"github.com/yuansuan/ticp/PSP/psp/internal/notice/util"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/ginutil"
	"github.com/yuansuan/ticp/PSP/psp/pkg/xtype"
)

// MessageList
//
//	@Summary		获取消息列表
//	@Description	获取消息列表接口
//	@Tags			通知-消息管理
//	@Accept			json
//	@Produce		json
//	@Param			param	body		dto.MessageListRequest	true	"请求参数"
//	@Response		200		{object}	dto.MessageListResponse
//	@Router			/notice/list [post]
func (r *apiRoute) MessageList(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	var req = dto.MessageListRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	userID := ginutil.GetUserID(ctx)
	msg := &dto.MessagePage{
		UserID:    userID,
		Page:      req.Page,
		Filter:    req.Filter,
		OrderSort: req.OrderSort,
	}

	messages, total, err := r.messageService.GetMessageList(ctx, msg)
	if err != nil {
		logger.Errorf("get message list err: %v", err)
		ginutil.Error(ctx, errcode.ErrNoticeFailList, errcode.MsgNoticeFailList)
		return
	}

	messageList := make([]*dto.Message, 0, len(messages))
	for _, message := range messages {
		msg := util.ConvertMessage(message)
		if msg != nil {
			messageList = append(messageList, msg)
		}
	}

	resp := &dto.MessageListResponse{
		Page: &xtype.PageResp{
			Index: req.Page.Index,
			Size:  req.Page.Size,
			Total: total,
		},
		Messages: messageList,
	}

	ginutil.Success(ctx, resp)
}

// MessageCount
//
//	@Summary		获取消息数量统计
//	@Description	获取消息数量统计接口
//	@Tags			通知-消息管理
//	@Accept			json
//	@Produce		json
//	@Param			param	query		dto.MessageCountRequest	true	"请求参数"
//	@Response		200		{object}	dto.MessageCountResponse
//	@Router			/notice/count [get]
func (r *apiRoute) MessageCount(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	var req = &dto.MessageCountRequest{}
	if err := ctx.BindQuery(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	userID := ginutil.GetUserID(ctx)
	total, err := r.messageService.GetMessageCount(ctx, userID, req.State)
	if err != nil {
		logger.Errorf("get message list err: %v", err)
		ginutil.Error(ctx, errcode.ErrNoticeFailList, errcode.MsgNoticeFailList)
		return
	}

	resp := &dto.MessageCountResponse{
		Total: total,
	}

	ginutil.Success(ctx, resp)
}
