package api

import (
	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/yuansuan/ticp/PSP/psp/internal/approve/dto"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/ginutil"
)

// Application
//
//	@Summary		审批申请列表
//	@Description	审批申请列表
//	@Tags			审批
//	@Accept			json
//	@Produce		json
//	@Param			param	body		dto.GetApproveListRequest	true	"入参"
//	@Response		200		{object}	dto.ApproveLogListAllResponse
//	@Router			/approve/list/application [post]
func (s *RouteService) Application(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	req := &dto.GetApproveListRequest{}
	if err := ctx.BindJSON(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	UserId := ginutil.GetUserID(ctx)

	data, err := s.ApproveService.GetApproveList(ctx, UserId, req)
	if err != nil {
		logger.Errorf("application list err: %v, req: [%+v]", err, req)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrApproveApplicationList)
		return
	}

	ginutil.Success(ctx, data)

}

// Pending
//
//	@Summary		待审批列表
//	@Description	待审批列表
//	@Tags			审批
//	@Accept			json
//	@Produce		json
//	@Param			param	body		dto.GetApprovePendingRequest	true	"入参"
//	@Response		200		{object}	dto.ApproveLogListAllResponse
//	@Router			/approve/list/pending [post]
func (s *RouteService) Pending(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	req := &dto.GetApprovePendingRequest{}
	if err := ctx.BindJSON(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	UserId := ginutil.GetUserID(ctx)

	data, err := s.ApproveService.GetApprovePendingList(ctx, UserId, req)
	if err != nil {
		logger.Errorf("Pending list err: %v, req: [%+v]", err, req)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrApprovePendingList)
		return
	}

	ginutil.Success(ctx, data)

}

// Complete
//
//	@Summary		已审批列表
//	@Description	已审批列表
//	@Tags			审批
//	@Accept			json
//	@Produce		json
//	@Param			param	body		dto.GetApproveCompleteRequest	true	"入参"
//	@Response		200		{object}	dto.ApproveLogListAllResponse
//	@Router			/approve/list/complete [post]
func (s *RouteService) Complete(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	req := &dto.GetApproveCompleteRequest{}
	if err := ctx.BindJSON(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	UserId := ginutil.GetUserID(ctx)

	data, err := s.ApproveService.GetApprovedList(ctx, UserId, req)
	if err != nil {
		logger.Errorf("Complete list err: %v, req: [%+v]", err, req)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrApprovedList)
		return
	}

	ginutil.Success(ctx, data)

}
