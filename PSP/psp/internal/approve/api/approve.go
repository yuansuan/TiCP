package api

import (
	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/yuansuan/ticp/PSP/psp/internal/approve/dto"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/ginutil"
)

// ApplyApprove
//
//	@Summary		发起审批
//	@Description	发起审批
//	@Tags			审批
//	@Accept			json
//	@Produce		json
//	@Param			param	body	dto.ApplyApproveRequest	true	"入参"
//	@Router			/approve/apply [post]
func (s *RouteService) ApplyApprove(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	req := &dto.ApplyApproveRequest{}
	if err := ctx.BindJSON(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	err := s.ApproveService.ApplyApprove(ctx, req)
	if err != nil {
		logger.Errorf("apply approve err: %v, req: [%+v]", err, req)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrApprovedApply)
		return
	}

	ginutil.Success(ctx, nil)
}

// CancelApprove
//
//	@Summary		取消审批
//	@Description	取消审批
//	@Tags			审批
//	@Accept			json
//	@Produce		json
//	@Param			param	body		dto.CancelApproveRequest	true	"入参"
//	@Response		200		{object}	dto.CancelApproveResponse
//	@Router			/approve/cancel [post]
func (s *RouteService) CancelApprove(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	req := &dto.CancelApproveRequest{}
	if err := ctx.BindJSON(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}
	RecordId := snowflake.MustParseString(req.Id)

	err := s.ApproveService.CancelApprove(ctx, int64(RecordId))
	if err != nil {
		logger.Errorf("CancelApprove  err: %v, req: [%+v]", err, req)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrApproveAuditLogList)
		return
	}
	ginutil.Success(ctx, dto.CancelApproveResponse{
		State: true,
	})
	return
}

// Pass
//
//	@Summary		通过审批
//	@Description	通过审批
//	@Tags			审批
//	@Accept			json
//	@Produce		json
//	@Param			param	body	dto.HandleApproveRequest	true	"入参"
//	@Router			/approve/pass [post]
func (s *RouteService) Pass(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	req := &dto.HandleApproveRequest{}
	if err := ctx.BindJSON(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	err := s.ApproveService.Pass(ctx, req)
	if err != nil {
		logger.Errorf("pass approve err: %v, req: [%+v]", err, req)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrApprovedPass)
		return
	}

	ginutil.Success(ctx, nil)
}

// Refuse
//
//	@Summary		拒绝审批
//	@Description	拒绝审批
//	@Tags			审批
//	@Accept			json
//	@Produce		json
//	@Param			param	body	dto.HandleApproveRequest	true	"入参"
//	@Router			/approve/refuse [post]
func (s *RouteService) Refuse(ctx *gin.Context) {

	logger := logging.GetLogger(ctx)
	req := &dto.HandleApproveRequest{}
	if err := ctx.BindJSON(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	err := s.ApproveService.Refuse(ctx, req)
	if err != nil {
		logger.Errorf("refuse approve err: %v, req: [%+v]", err, req)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrApprovedRefuse)
		return
	}

	ginutil.Success(ctx, nil)
}
