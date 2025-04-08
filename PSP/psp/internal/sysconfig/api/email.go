package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/internal/sysconfig/dto"
	"github.com/yuansuan/ticp/PSP/psp/pkg/tracelog"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/ginutil"
)

// SetEmailConfig
//
//	@Summary		设置邮件通知配置
//	@Description	设置邮件通知接口
//	@Tags			系统配置
//	@Accept			json
//	@Produce		json
//	@Param			param	body		dto.SetEmailConfigReq	true	"请求参数"
//	@Response		200		{object}	dto.SetEmailConfigRes
//	@Router			/sysconfig/setEmailConfig [post]
func (s *RouteService) SetEmailConfig(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	req := &dto.SetEmailConfigReq{}
	if err := ctx.ShouldBindJSON(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	tracelog.Info(ctx, fmt.Sprintf("user: [%v] set email config req: [%+v]", ginutil.GetUserID(ctx), req))

	err := s.SysConfigService.SetEmailConfig(ctx, req)
	if err != nil {
		logger.Errorf("set email config err: %v, req: [%+v]", err, req)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrSysConfigSetEmailFailed)
		return
	}

	ginutil.Success(ctx, &dto.SetEmailConfigRes{})
}

// GetEmailConfig
//
//	@Summary		获取邮件通知配置
//	@Description	获取邮件通知接口
//	@Tags			系统配置
//	@Accept			json
//	@Produce		json
//	@Response		200	{object}	dto.GetEmailConfigRes
//	@Router			/sysconfig/getEmailConfig [get]
func (s *RouteService) GetEmailConfig(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	emailConfig, err := s.SysConfigService.GetEmailConfig(ctx)
	if err != nil {
		logger.Errorf("get email config err: %v, result: [%+v]", err, emailConfig)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrSysConfigGetEmailFailed)
		return
	}

	ginutil.Success(ctx, &emailConfig)
}

// SetGlobalEmail
//
//	@Summary		设置电子邮件信息
//	@Description	设置电子邮件信息接口
//	@Tags			通知-电子邮件
//	@Accept			json
//	@Produce		json
//	@Param			param	body		dto.SetGlobalEmailRequest	true	"请求参数"
//	@Response		200		{object}	dto.SetGlobalEmailResponse
//	@Router			/sysconfig/globalEmail [post]
func (s *RouteService) SetGlobalEmail(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	req := &dto.SetGlobalEmailRequest{}

	if err := ctx.ShouldBindJSON(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	tracelog.Info(ctx, fmt.Sprintf("user: [%v] set global config req: [%+v]", ginutil.GetUserID(ctx), req))

	err := s.SysConfigService.SetGlobalEmail(ctx, req.EmailConfig)
	if err != nil {
		logger.Errorf("set email info err: %v, req: [%+v]", err, req)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrNoticeSetEmailFailed)
		return
	}

	ginutil.Success(ctx, &dto.SetGlobalEmailResponse{})
}

// GetGlobalEmail
//
//	@Summary		获取邮件通知配置
//	@Description	获取邮件通知接口
//	@Tags			系统配置
//	@Accept			json
//	@Produce		json
//	@Response		200	{object}	dto.GetGlobalEmailRes
//	@Router			/sysconfig/globalEmail [get]
func (s *RouteService) GetGlobalEmail(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	emailConfig, err := s.SysConfigService.GetGlobalEmail(ctx)
	if err != nil {
		logger.Errorf("get email config err: %v, result: [%+v]", err, emailConfig)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrNoticeGetEmailFailed)
		return
	}

	ginutil.Success(ctx, &dto.GetGlobalEmailRes{
		EmailConfig: emailConfig,
	})
}

// TestSendEmail
//
//	@Summary		测试发送电子邮件
//	@Description	测试发送电子邮件接口
//	@Tags			通知-电子邮件
//	@Accept			json
//	@Produce		json
//	@Response		200	{object}	dto.SendEmailTestResponse
//	@Router			/sysconfig/email/testSend [post]
func (r *RouteService) TestSendEmail(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	tracelog.Info(ctx, fmt.Sprintf("user: [%v] test send email", ginutil.GetUserID(ctx)))

	err := r.SysConfigService.SendEmail(ctx)
	if err != nil {
		logger.Errorf("test send email err: %v", err)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrNoticeTestSendEmailFailed)
		return
	}

	ginutil.Success(ctx, &dto.SendEmailTestResponse{})
}
