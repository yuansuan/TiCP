package api

import (
	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/yuansuan/ticp/PSP/psp/internal/approve/dto"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/ginutil"
)

// List
//
//	@Summary		操作日志列表
//	@Description	操作日志列表
//	@Tags			操作日志
//	@Accept			json
//	@Produce		json
//	@Param			param	body		dto.AuditLogListRequest		true	"入参"
//	@Success		200		{object}	dto.AuditLogListResponse	"正常返回表示成功，否则返回错误码和错误信息"
//	@Router			/auditlog/list [post]
func (s *RouteService) List(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	req := &dto.AuditLogListRequest{}
	if err := ctx.BindJSON(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	data, err := s.AuditLogService.List(ctx, req)
	if err != nil {
		logger.Errorf("get audit list err: %v, req: [%+v]", err, req)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrApproveAuditLogList)
		return
	}

	ginutil.Success(ctx, data)
}

// ListAll
//
//	@Summary		操作日志列表
//	@Description	操作日志列表
//	@Tags			操作日志
//	@Accept			json
//	@Produce		json
//	@Param			param	body		dto.AuditLogListAllRequest	true	"入参"
//	@Response		200		{object}	dto.AuditLogListAllResponse	"正常返回表示成功，否则返回错误码和错误信息"
//	@Router			/auditlog/listAll [post]
func (s *RouteService) ListAll(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	req := &dto.AuditLogListAllRequest{}
	if err := ctx.BindJSON(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	data, err := s.AuditLogService.ListAll(ctx, req)
	if err != nil {
		logger.Errorf("ListAll err: %v, req: [%+v]", err, req)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrApproveAuditLogList)
		return
	}

	ginutil.Success(ctx, data)
}

// Export
//
//	@Summary		操作日志列表-导出
//	@Description	操作日志列表-导出
//	@Tags			操作日志
//	@Accept			json
//	@Produce		json
//	@Param			param	body	dto.AuditLogListRequest	true	"入参"
//	@Router			/auditlog/export [post]
func (s *RouteService) Export(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	req := &dto.AuditLogListRequest{}
	if err := ctx.BindJSON(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	err := s.AuditLogService.Export(ctx, req)
	if err != nil {
		logger.Errorf("export audit list err: %v, req: [%+v]", err, req)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrApproveAuditLogList)
		return
	}

}

// ExportAll
//
//	@Summary		操作日志列表-导出所有信息
//	@Description	操作日志列表-导出所有信息
//	@Tags			操作日志
//	@Accept			json
//	@Produce		json
//	@Param			param	body	dto.AuditLogListAllRequest	true	"入参"
//	@Router			/auditlog/exportAll [post]
func (s *RouteService) ExportAll(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	req := &dto.AuditLogListAllRequest{}
	if err := ctx.BindJSON(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	err := s.AuditLogService.ExportAll(ctx, req)
	if err != nil {
		logger.Errorf("export audit list err: %v, req: [%+v]", err, req)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrApproveAuditLogList)
		return
	}

}
