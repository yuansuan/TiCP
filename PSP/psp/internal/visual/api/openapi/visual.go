package openapi

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"google.golang.org/grpc/status"

	"github.com/yuansuan/ticp/PSP/psp/internal/common"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/oplog"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/approve"
	"github.com/yuansuan/ticp/PSP/psp/internal/visual/dto"
	"github.com/yuansuan/ticp/PSP/psp/internal/visual/dto/openapi"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
	"github.com/yuansuan/ticp/PSP/psp/pkg/tracelog"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/ginutil"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/structutil"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/strutil"
)

// ListSoftware
//
//	@Summary		软件列表
//	@Description	软件列表接口
//	@Tags			可视化-软件
//	@Accept			json
//	@Produce		json
//	@Param			param	query		openapi.ListSoftwareRequest	true	"请求参数"
//	@Response		200		{object}	openapi.ListSoftwareResponse
//	@Router			/openapi/vis/software [get]
func (s *RouteOpenapiService) ListSoftware(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	req := &openapi.ListSoftwareRequest{}
	if err := ctx.BindQuery(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	if req.PageSize == 0 {
		req.PageSize = common.DefaultMaxPageSize
	}

	userID := ginutil.GetUserID(ctx)
	username := ginutil.GetUserName(ctx)

	loginUserID := ginutil.GetUserID(ctx)

	softwares, total, err := s.visualService.ListSoftware(ctx, userID, "", "", "published", username,
		true, false, false, req.PageIndex, req.PageSize, snowflake.ID(loginUserID))
	if err != nil {
		logger.Errorf("list software err: %v, req: [%+v]", err, req)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrVisualListSoftwareFailed)
		return
	}

	rsp := &openapi.ListSoftwareResponse{}
	if err = structutil.CopyStruct(rsp, &dto.ListSoftwareResponse{Softwares: softwares, Total: total}); err != nil {
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrVisualListSoftwareFailed)
		return
	}

	ginutil.Success(ctx, rsp)
}

// ListHardware
//
//	@Summary		硬件列表
//	@Description	硬件列表接口
//	@Tags			可视化-硬件
//	@Accept			json
//	@Produce		json
//	@Param			param	query		openapi.ListHardwareRequest	true	"请求参数"
//	@Response		200		{object}	openapi.ListHardwareResponse
//	@Router			/openapi/vis/hardware [get]
func (s *RouteOpenapiService) ListHardware(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	req := &openapi.ListHardwareRequest{}
	if err := ctx.BindQuery(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	if req.PageSize == 0 {
		req.PageSize = common.DefaultMaxPageSize
	}

	loginUserID := ginutil.GetUserID(ctx)

	hardwares, total, err := s.visualService.ListHardware(ctx, "", false, false,
		0, 0, 0, req.PageIndex, req.PageSize, snowflake.ID(loginUserID))
	if err != nil {
		logger.Errorf("list hardware err: %v, req: [%+v]", err, req)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrVisualListHardwareFailed)
		return
	}

	rsp := &openapi.ListHardwareResponse{}
	if err = structutil.CopyStruct(rsp, &dto.ListHardwareResponse{Hardwares: hardwares, Total: total}); err != nil {
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrVisualListSoftwareFailed)
		return
	}

	ginutil.Success(ctx, rsp)
}

// CloseSession
//
//	@Summary		关闭会话
//	@Description	关闭会话接口
//	@Tags			可视化-会话
//	@Accept			json
//	@Produce		json
//	@Param			param	body		openapi.CloseSessionRequest	true	"关闭会话参数"
//	@Response		200		{object}	openapi.CloseSessionResponse
//	@Router			/openapi/vis/session/close [post]
func (s *RouteOpenapiService) CloseSession(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	req := &openapi.CloseSessionRequest{}
	if err := ctx.BindJSON(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	if err := s.validate.Struct(req); err != nil {
		ginutil.Error(ctx, errcode.ErrInvalidParam, err.Error())
		return
	}

	tracelog.Info(ctx, fmt.Sprintf("user: [%v] close session req: [%+v]", ginutil.GetUserID(ctx), req))

	outSessionID, err := s.visualService.CloseSession(ctx, req.SessionID, req.ExitReason, false)
	if err != nil {
		logger.Errorf("close session err: %v, req: [%+v]", err, req)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrVisualCloseSessionFailed)
		return
	}

	oplog.GetInstance().SaveAuditLogInfo(ctx, approve.OperateTypeEnum_VIS_MANAGER, fmt.Sprintf("【OPENAPi】用户%v关闭会话[%v]", ginutil.GetUserName(ctx), outSessionID))

	ginutil.Success(ctx, &openapi.CloseSessionResponse{Success: true})
}

// RebootSession
//
//	@Summary		重启会话
//	@Description	重启会话接口
//	@Tags			可视化-会话
//	@Accept			json
//	@Produce		json
//	@Param			param	body		openapi.RebootSessionRequest	true	"重启会话参数"
//	@Response		200		{object}	openapi.RebootSessionResponse
//	@Router			/openapi/vis/session/reboot [post]
func (s *RouteOpenapiService) RebootSession(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	req := &openapi.RebootSessionRequest{}
	if err := ctx.BindJSON(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	if err := s.validate.Struct(req); err != nil {
		ginutil.Error(ctx, errcode.ErrInvalidParam, err.Error())
		return
	}

	if req.SessionID == "" {
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	tracelog.Info(ctx, fmt.Sprintf("user: [%v] reboot session req: [%+v]", ginutil.GetUserID(ctx), req))

	ready, err := s.visualService.RebootSession(ctx, req.SessionID, "", false)
	if err != nil {
		logger.Errorf("reboot session err: %v, req: [%+v]", err, req)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrVisualRebootSessionFailed)
		return
	}

	oplog.GetInstance().SaveAuditLogInfo(ctx, approve.OperateTypeEnum_VIS_MANAGER, fmt.Sprintf("【OPENAPi】用户%v重启会话[%v]", ginutil.GetUserName(ctx), ready))

	ginutil.Success(ctx, &openapi.RebootSessionResponse{Status: ready})
}

// StartSession
//
//	@Summary		创建会话
//	@Description	创建会话接口
//	@Tags			可视化-会话
//	@Accept			json
//	@Produce		json
//	@Param			param	body		openapi.StartSessionRequest	true	"请求参数"
//	@Response		200		{object}	openapi.StartSessionResponse
//	@Router			/openapi/vis/session [post]
func (s *RouteOpenapiService) StartSession(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	req := &openapi.StartSessionRequest{}
	if err := ctx.BindJSON(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	if err := s.validate.Struct(req); err != nil {
		ginutil.Error(ctx, errcode.ErrInvalidParam, err.Error())
		return
	}

	if strutil.IsEmpty(req.ProjectID) {
		req.ProjectID = common.PersonalProjectID.String()
	}

	loginUserID := ginutil.GetUserID(ctx)
	loginUserName := ginutil.GetUserName(ctx)

	tracelog.Info(ctx, fmt.Sprintf("user: [%v] start session req: [%+v]", ginutil.GetUserID(ctx), req))

	sessionID, projectName, outSessionID, err := s.visualService.StartSession(ctx, req.ProjectID, req.HardwareID, req.SoftwareID, loginUserName, req.Mounts, snowflake.ID(loginUserID))
	if err != nil {
		logger.Errorf("start session err: %v, req: [%+v]", err, req)
		switch status.Code(err) {
		case errcode.ErrVisualSessionRepeatStart:
			ginutil.Error(ctx, errcode.ErrVisualSessionRepeatStart, fmt.Sprintf(errcode.VisualCodeMsg[errcode.ErrVisualSessionRepeatStart], projectName, sessionID))
			return
		}
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrVisualStartSessionFailed)
		return
	}

	oplog.GetInstance().SaveAuditLogInfo(ctx, approve.OperateTypeEnum_VIS_MANAGER, fmt.Sprintf("【OPENAPi】用户%v创建会话[%v]", ginutil.GetUserName(ctx), outSessionID))

	ginutil.Success(ctx, &openapi.StartSessionResponse{SessionID: sessionID})
}

// ListSession
//
//	@Summary		会话列表
//	@Description	会话列表接口
//	@Tags			可视化-会话
//	@Accept			json
//	@Produce		json
//	@Param			param	query		dto.ListSessionsRequest	true	"请求参数"
//	@Response		200		{object}	dto.ListSessionsResponse
//	@Router			/openapi/vis/session/list [get]
func (s *RouteOpenapiService) ListSession(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	req := &openapi.ListSessionsRequest{}
	if err := ctx.BindJSON(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	if req.PageSize == 0 {
		req.PageSize = common.DefaultMaxPageSize
	}

	loginUserID := ginutil.GetUserID(ctx)

	sessions, total, err := s.visualService.ListSession(ctx, req.HardwareIDs, req.SoftwareIDs, req.ProjectIDs,
		req.Statuses, ginutil.GetUserName(ctx), false, req.PageIndex, req.PageSize, snowflake.ID(loginUserID))
	if err != nil {
		logger.Errorf("list session err: %v, req: [%+v]", err, req)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrVisualListSessionFailed)
		return
	}

	rsp := &openapi.ListSessionsResponse{}
	if err = structutil.CopyStruct(rsp, &dto.ListSessionsResponse{Sessions: sessions, Total: total}); err != nil {
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrVisualListSessionFailed)
		return
	}

	ginutil.Success(ctx, rsp)
}

func (s *RouteOpenapiService) SessionInfo(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	req := &openapi.SessionInfoRequest{}
	if err := ctx.BindQuery(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	if err := s.validate.Struct(req); err != nil {
		ginutil.Error(ctx, errcode.ErrInvalidParam, err.Error())
		return
	}

	sessionInfo, err := s.visualService.SessionInfo(ctx, req.SessionID)

	if err != nil {
		logger.Errorf("list session err: %v, req: [%+v]", err, req)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrVisualListSessionFailed)
		return
	}
	rsp := &openapi.Session{}
	if err = structutil.CopyStruct(rsp, sessionInfo); err != nil {
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrVisualListSessionFailed)
		return
	}

	ginutil.Success(ctx, rsp)
}
