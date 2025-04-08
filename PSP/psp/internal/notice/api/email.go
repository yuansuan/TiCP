package api

import (
	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/internal/notice/consts"
	"github.com/yuansuan/ticp/PSP/psp/internal/notice/dto"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/ginutil"
)

// GetEmail
//
//	@Summary		获取电子邮件信息
//	@Description	获取电子邮件信息接口
//	@Tags			通知-电子邮件
//	@Accept			json
//	@Produce		json
//	@Param			param	query		dto.GetEmailRequest	true	"请求参数"
//	@Response		200		{object}	dto.GetEmailResponse
//	@Router			/notice/email [get]
func (r *apiRoute) GetEmail(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	req := &dto.GetEmailRequest{}
	if err := ctx.BindQuery(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	email, err := r.emailService.GetEmail(ctx)
	if err != nil {
		logger.Errorf("get email info err: %v, req: [%+v]", err, req)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrNoticeGetEmailFailed)
		return
	}
	ginutil.Success(ctx, &dto.GetEmailResponse{Email: email})
}

// SetEmail
//
//	@Summary		设置电子邮件信息
//	@Description	设置电子邮件信息接口
//	@Tags			通知-电子邮件
//	@Accept			json
//	@Produce		json
//	@Param			param	body		dto.SetEmailRequest	true	"请求参数"
//	@Response		200		{object}	dto.SetEmailResponse
//	@Router			/notice/email [post]
func (r *apiRoute) SetEmail(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	req := &dto.SetEmailRequest{}
	if err := ctx.ShouldBindJSON(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	if req.Email == nil || req.Email.Setting == nil {
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	err := r.emailService.SetEmail(ctx, req.Email)
	if err != nil {
		logger.Errorf("set email info err: %v, req: [%+v]", err, req)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrNoticeSetEmailFailed)
		return
	}
	ginutil.Success(ctx, &dto.SetEmailResponse{})
}

// SendEmail
//
//	@Summary		发送电子邮件
//	@Description	发送电子邮件接口
//	@Tags			通知-电子邮件
//	@Accept			json
//	@Produce		json
//	@Param			param	body		dto.SendEmailRequest	true	"请求参数"
//	@Response		200		{object}	dto.SendEmailResponse
//	@Router			/notice/email/send [post]
func (r *apiRoute) SendEmail(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	req := &dto.SendEmailRequest{}
	if err := ctx.ShouldBindJSON(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	if req.Receiver == "" || req.EmailTemplate == "" {
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	err := r.emailService.SendEmail(ctx, req.Receiver, req.EmailTemplate, req.JsonData)
	if err != nil {
		logger.Errorf("send email err: %v, req: [%+v]", err, req)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrNoticeSendEmailFailed)
		return
	}
	ginutil.Success(ctx, &dto.SendEmailResponse{})
}

// TestSendEmail
//
//	@Summary		测试发送电子邮件
//	@Description	测试发送电子邮件接口
//	@Tags			通知-电子邮件
//	@Accept			json
//	@Produce		json
//	@Param			param	body		dto.TestSendEmailRequest	true	"请求参数"
//	@Response		200		{object}	dto.TestSendEmailResponse
//	@Router			/notice/email/testSend [post]
func (r *apiRoute) TestSendEmail(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	req := &dto.TestSendEmailRequest{}
	if err := ctx.ShouldBindJSON(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	if req.Receiver == "" {
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	err := r.emailService.SendEmail(ctx, req.Receiver, consts.DefaultTemplate, "")
	if err != nil {
		logger.Errorf("test send email err: %v, req: [%+v]", err, req)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrNoticeTestSendEmailFailed)
		return
	}
	ginutil.Success(ctx, &dto.TestSendEmailResponse{})
}
