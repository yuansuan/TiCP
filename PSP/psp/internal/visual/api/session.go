package api

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
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
	"github.com/yuansuan/ticp/PSP/psp/pkg/tracelog"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/ginutil"
)

// ListSession
//
//	@Summary		会话列表
//	@Description	会话列表接口
//	@Tags			可视化-会话
//	@Accept			json
//	@Produce		json
//	@Param			param	query		dto.ListSessionsRequest	true	"请求参数"
//	@Response		200		{object}	dto.ListSessionsResponse
//	@Router			/vis/session [get]
func (s *RouteService) ListSession(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	req := &dto.ListSessionsRequest{}
	if err := ctx.BindQuery(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	if req.PageSize == 0 {
		req.PageSize = common.DefaultMaxPageSize
	}

	loginUserID := ginutil.GetUserID(ctx)

	sessions, total, err := s.visualService.ListSession(ctx, req.HardwareIDs, req.SoftwareIDs, req.ProjectIDs,
		req.Statuses, req.UserName, req.IsAdmin, req.PageIndex, req.PageSize, snowflake.ID(loginUserID))
	if err != nil {
		logger.Errorf("list session err: %v, req: [%+v]", err, req)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrVisualListSessionFailed)
		return
	}
	ginutil.Success(ctx, &dto.ListSessionsResponse{Sessions: sessions, Total: total})
}

// StartSession
//
//	@Summary		创建会话
//	@Description	创建会话接口
//	@Tags			可视化-会话
//	@Accept			json
//	@Produce		json
//	@Param			param	body		dto.StartSessionRequest	true	"请求参数"
//	@Response		200		{object}	dto.StartSessionResponse
//	@Router			/vis/session [post]
func (s *RouteService) StartSession(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	req := &dto.StartSessionRequest{}
	if err := ctx.BindJSON(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	if req.HardwareID == "" || req.SoftwareID == "" {
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
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

	oplog.GetInstance().SaveAuditLogInfo(ctx, approve.OperateTypeEnum_VIS_MANAGER, fmt.Sprintf("用户%v创建会话[%v]", ginutil.GetUserName(ctx), outSessionID))

	ginutil.Success(ctx, &dto.StartSessionResponse{SessionID: sessionID})
}

// GetMountInfo
//
//	@Summary		获取挂载信息
//	@Description	获取挂载信息接口
//	@Tags			可视化-会话
//	@Accept			json
//	@Produce		json
//	@Param			param	query		dto.GetMountInfoRequest	true	"请求参数"
//	@Response		200		{object}	dto.GetMountInfoResponse
//	@Router			/vis/session/getMountInfo [get]
func (s *RouteService) GetMountInfo(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	req := &dto.GetMountInfoRequest{}
	if err := ctx.BindQuery(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	loginUserID := ginutil.GetUserID(ctx)
	loginUserName := ginutil.GetUserName(ctx)

	res, err := s.visualService.GetMountInfo(ctx, req.ProjectID, loginUserName, snowflake.ID(loginUserID))
	if err != nil {
		logger.Errorf("get mount info err: %v, req: [%+v]", err, req)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrVisualGetMountInfoFailed)
		return
	}

	ginutil.Success(ctx, res)
}

// PowerOffSession
//
//	@Summary		关机指定会话
//	@Description	关机指定会话接口
//	@Tags			可视化-会话
//	@Accept			json
//	@Produce		json
//	@Param			param	body		dto.PowerOffSessionRequest	true	"关闭会话参数"
//	@Response		200		{object}	dto.PowerOffSessionResponse
//	@Router			/vis/session/powerOff [post]
func (s *RouteService) PowerOffSession(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	req := &dto.PowerOffSessionRequest{}
	if err := ctx.BindJSON(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	if req.SessionID == "" {
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	tracelog.Info(ctx, fmt.Sprintf("user: [%v] power off session req: [%+v]", ginutil.GetUserID(ctx), req))

	outSessionID, err := s.visualService.PowerOffSession(ctx, req.SessionID)
	if err != nil {
		logger.Errorf("power off session err: %v, req: [%+v]", err, req)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrVisualPowerOffSessionFailed)
		return
	}

	oplog.GetInstance().SaveAuditLogInfo(ctx, approve.OperateTypeEnum_VIS_MANAGER, fmt.Sprintf("用户%v关机指定会话[%v]", ginutil.GetUserName(ctx), outSessionID))

	ginutil.Success(ctx, &dto.PowerOffSessionResponse{Success: true})
}

// PowerOnSession
//
//	@Summary		开机指定会话
//	@Description	开机指定会话接口
//	@Tags			可视化-会话
//	@Accept			json
//	@Produce		json
//	@Param			param	body		dto.PowerOnSessionRequest	true	"关闭会话参数"
//	@Response		200		{object}	dto.PowerOnSessionResponse
//	@Router			/vis/session/powerOn [post]
func (s *RouteService) PowerOnSession(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	req := &dto.PowerOnSessionRequest{}
	if err := ctx.BindJSON(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	if req.SessionID == "" {
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	tracelog.Info(ctx, fmt.Sprintf("user: [%v] power on session req: [%+v]", ginutil.GetUserID(ctx), req))

	outSessionID, err := s.visualService.PowerOnSession(ctx, req.SessionID)
	if err != nil {
		logger.Errorf("power on session err: %v, req: [%+v]", err, req)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrVisualPowerOnSessionFailed)
		return
	}

	oplog.GetInstance().SaveAuditLogInfo(ctx, approve.OperateTypeEnum_VIS_MANAGER, fmt.Sprintf("用户%v开机指定会话[%v]", ginutil.GetUserName(ctx), outSessionID))

	ginutil.Success(ctx, &dto.PowerOnSessionResponse{Success: true})
}

// CloseSession
//
//	@Summary		关闭会话
//	@Description	关闭会话接口
//	@Tags			可视化-会话
//	@Accept			json
//	@Produce		json
//	@Param			param	body		dto.CloseSessionRequest	true	"关闭会话参数"
//	@Response		200		{object}	dto.CloseSessionResponse
//	@Router			/vis/session/close [post]
func (s *RouteService) CloseSession(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	req := &dto.CloseSessionRequest{}
	if err := ctx.BindJSON(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	if req.SessionID == "" {
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}
	if req.Admin && req.ExitReason == "" {
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	tracelog.Info(ctx, fmt.Sprintf("user: [%v] close session req: [%+v]", ginutil.GetUserID(ctx), req))

	outSessionID, err := s.visualService.CloseSession(ctx, req.SessionID, req.ExitReason, req.Admin)
	if err != nil {
		logger.Errorf("close session err: %v, req: [%+v]", err, req)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrVisualCloseSessionFailed)
		return
	}

	oplog.GetInstance().SaveAuditLogInfo(ctx, approve.OperateTypeEnum_VIS_MANAGER, fmt.Sprintf("用户%v关闭会话[%v]", ginutil.GetUserName(ctx), outSessionID))

	ginutil.Success(ctx, &dto.CloseSessionResponse{Success: true})
}

// RebootSession
//
//	@Summary		重启会话
//	@Description	重启会话接口
//	@Tags			可视化-会话
//	@Accept			json
//	@Produce		json
//	@Param			param	body		dto.RebootSessionRequest	true	"重启会话参数"
//	@Response		200		{object}	dto.RebootSessionResponse
//	@Router			/vis/session/reboot [post]
func (s *RouteService) RebootSession(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	req := &dto.RebootSessionRequest{}
	if err := ctx.BindJSON(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
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

	oplog.GetInstance().SaveAuditLogInfo(ctx, approve.OperateTypeEnum_VIS_MANAGER, fmt.Sprintf("用户%v重启会话[%v]", ginutil.GetUserName(ctx), ready))

	ginutil.Success(ctx, &dto.RebootSessionResponse{Status: ready})
}

// ReadySession
//
//	@Summary		会话启动状态
//	@Description	会话启动状态接口
//	@Tags			可视化-会话
//	@Accept			json
//	@Produce		json
//	@Param			param	query		dto.ReadySessionRequest	true	"请求参数"
//	@Response		200		{object}	dto.ReadySessionResponse
//	@Router			/vis/session/ready [get]
func (s *RouteService) ReadySession(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	req := &dto.ReadySessionRequest{}
	if err := ctx.BindQuery(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	if req.SessionID == "" {
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	ready, err := s.visualService.ReadySession(ctx, req.SessionID)
	if err != nil {
		logger.Errorf("ready session err: %v, req: [%+v]", err, req)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrVisualReadySessionFailed)
		return
	}
	ginutil.Success(ctx, &dto.ReadySessionResponse{Ready: ready})
}

// GetRemoteAppURL
//
//	@Summary		远程应用URL
//	@Description	远程应用URL接口
//	@Tags			可视化-会话
//	@Accept			json
//	@Produce		json
//	@Param			param	query		dto.GetRemoteAppURLRequest	true	"请求参数"
//	@Response		200		{object}	dto.GetRemoteAppURLResponse
//	@Router			/vis/session/remoteAppUrl [get]
func (s *RouteService) GetRemoteAppURL(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	req := &dto.GetRemoteAppURLRequest{}
	if err := ctx.BindQuery(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	if req.SessionID == "" || req.RemoteAppName == "" {
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	url, err := s.visualService.GetRemoteAppURL(ctx, req.SessionID, req.RemoteAppName)
	if err != nil {
		logger.Errorf("get remote app url err: %v, req: [%+v]", err, req)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrVisualGetRemoteAppURLFailed)
		return
	}
	ginutil.Success(ctx, &dto.GetRemoteAppURLResponse{URL: url})
}

// ListUsedProjectNames
//
//	@Summary		获取已使用项目名称列表
//	@Description	获取已使用项目名称列表接口
//	@Tags			可视化-会话
//	@Accept			json
//	@Produce		json
//	@Param			param	query		dto.ListUsedProjectNamesRequest	true	"请求参数"
//	@Response		200		{object}	dto.ListUsedProjectNamesResponse
//	@Router			/vis/session/projectNames [get]
func (s *RouteService) ListUsedProjectNames(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	req := &dto.ListUsedProjectNamesRequest{}
	if err := ctx.BindQuery(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	username := ""
	if req.HasUsed {
		username = ginutil.GetUserName(ctx)
	}

	names, err := s.visualService.ListUsedProjectNames(ctx, username)
	if err != nil {
		logger.Errorf("list used project names err: %v, req: [%+v]", err, req)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrVisualListUsedProjectNamesFailed)
		return
	}

	ginutil.Success(ctx, &dto.ListUsedProjectNamesResponse{Names: names})
}

// ExportSessionInfo
//
//	@Summary		导出会话信息
//	@Description	导出会话信息接口
//	@Tags			可视化-会话
//	@Accept			json
//	@Produce		json
//	@Param			param	query	dto.ExportSessionInfoRequest	true	"请求参数"
//	@Response		200
//	@Router			/vis/session/export [get]
func (s *RouteService) ExportSessionInfo(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	req := &dto.ExportSessionInfoRequest{}
	if err := ctx.BindQuery(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	if req.Start == 0 || req.End == 0 {
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	err := s.visualService.ExportSessionInfo(ctx, req.Start, req.End)
	if err != nil {
		logger.Errorf("export session info err: %v, req: [%+v]", err, req)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrVisualExportSessionInfoFailed)
		return
	}
}
