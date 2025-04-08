package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/yuansuan/ticp/PSP/psp/internal/common"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/oplog"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/approve"
	"github.com/yuansuan/ticp/PSP/psp/internal/visual/consts"
	"github.com/yuansuan/ticp/PSP/psp/internal/visual/dto"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
	"github.com/yuansuan/ticp/PSP/psp/pkg/tracelog"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/ginutil"
)

// ListSoftware
//
//	@Summary		软件列表
//	@Description	软件列表接口
//	@Tags			可视化-软件
//	@Accept			json
//	@Produce		json
//	@Param			param	query		dto.ListSoftwareRequest	true	"请求参数"
//	@Response		200		{object}	dto.ListSoftwareResponse
//	@Router			/vis/software [get]
func (s *RouteService) ListSoftware(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	req := &dto.ListSoftwareRequest{}
	if err := ctx.BindQuery(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	if req.HasPermission && req.HasUsed {
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	if req.PageSize == 0 {
		req.PageSize = common.DefaultMaxPageSize
	}

	userID := ginutil.GetUserID(ctx)
	username := ginutil.GetUserName(ctx)

	loginUserID := ginutil.GetUserID(ctx)

	softwares, total, err := s.visualService.ListSoftware(ctx, userID, req.Name, req.Platform, req.State, username,
		req.HasPermission, req.HasUsed, req.IsAdmin, req.PageIndex, req.PageSize, snowflake.ID(loginUserID))
	if err != nil {
		logger.Errorf("list software err: %v, req: [%+v]", err, req)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrVisualListSoftwareFailed)
		return
	}
	ginutil.Success(ctx, &dto.ListSoftwareResponse{Softwares: softwares, Total: total})
}

// AddSoftware
//
//	@Summary		新增软件
//	@Description	新增软件接口
//	@Tags			可视化-软件
//	@Accept			json
//	@Produce		json
//	@Param			param	body		dto.AddSoftwareRequest	true	"请求参数"
//	@Response		200		{object}	dto.AddSoftwareResponse
//	@Router			/vis/software [post]
func (s *RouteService) AddSoftware(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	req := &dto.AddSoftwareRequest{}
	if err := ctx.BindJSON(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	if req.Name == "" || req.Platform == "" || req.ImageID == "" {
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	tracelog.Info(ctx, fmt.Sprintf("user: [%v] add software req: [%+v]", ginutil.GetUserID(ctx), req))

	id, err := s.visualService.AddSoftware(ctx, &dto.Software{
		Name:       req.Name,
		Desc:       req.Desc,
		Platform:   req.Platform,
		ImageID:    req.ImageID,
		InitScript: req.InitScript,
		Icon:       req.Icon,
		GPUDesired: req.GPUDesired,
	})
	if err != nil {
		logger.Errorf("add software err: %v, req: [%+v]", err, req)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrVisualAddSoftwareFailed)
		return
	}

	oplog.GetInstance().SaveAuditLogInfo(ctx, approve.OperateTypeEnum_VIS_MANAGER, fmt.Sprintf("用户%v添加镜像[%v]", ginutil.GetUserName(ctx), req.Name))

	ginutil.Success(ctx, &dto.AddSoftwareResponse{ID: id})
}

// UpdateSoftware
//
//	@Summary		更新软件
//	@Description	更新软件接口
//	@Tags			可视化-软件
//	@Accept			json
//	@Produce		json
//	@Param			param	body		dto.UpdateSoftwareRequest	true	"请求参数"
//	@Response		200		{object}	dto.UpdateSoftwareResponse
//	@Router			/vis/software [put]
func (s *RouteService) UpdateSoftware(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	req := &dto.UpdateSoftwareRequest{}
	if err := ctx.BindJSON(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	if req.ID == "" || req.Name == "" || req.Platform == "" || req.ImageID == "" {
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	tracelog.Info(ctx, fmt.Sprintf("user: [%v] update software req: [%+v]", ginutil.GetUserID(ctx), req))

	err := s.visualService.UpdateSoftware(ctx, &dto.Software{
		ID:         req.ID,
		Name:       req.Name,
		Desc:       req.Desc,
		Platform:   req.Platform,
		ImageID:    req.ImageID,
		InitScript: req.InitScript,
		Icon:       req.Icon,
		GPUDesired: req.GPUDesired,
	})
	if err != nil {
		logger.Errorf("add software err: %v, req: [%+v]", err, req)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrVisualUpdateSoftwareFailed)
		return
	}

	oplog.GetInstance().SaveAuditLogInfo(ctx, approve.OperateTypeEnum_VIS_MANAGER, fmt.Sprintf("用户%v编辑镜像[%v]", ginutil.GetUserName(ctx), req.Name))

	ginutil.Success(ctx, &dto.UpdateSoftwareResponse{})
}

// DeleteSoftware
//
//	@Summary		删除软件
//	@Description	删除软件接口
//	@Tags			可视化-软件
//	@Accept			json
//	@Produce		json
//	@Param			param	query		dto.DeleteSoftwareRequest	true	"请求参数"
//	@Response		200		{object}	dto.DeleteSoftwareResponse
//	@Router			/vis/software [delete]
func (s *RouteService) DeleteSoftware(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	req := &dto.DeleteSoftwareRequest{}
	if err := ctx.BindQuery(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	if req.ID == "" {
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	tracelog.Info(ctx, fmt.Sprintf("user: [%v] delete software req: [%+v]", ginutil.GetUserID(ctx), req))

	software, _ := s.visualService.GetSoftware(ctx, snowflake.MustParseString(req.ID))
	err := s.visualService.DeleteSoftware(ctx, req.ID)
	if err != nil {
		logger.Errorf("delete software err: %v, req: [%+v]", err, req)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrVisualDeleteSoftwareFailed)
		return
	}

	var softwareName string
	if software != nil {
		softwareName = software.Name
	}
	oplog.GetInstance().SaveAuditLogInfo(ctx, approve.OperateTypeEnum_VIS_MANAGER, fmt.Sprintf("用户%v删除镜像[%v]", ginutil.GetUserName(ctx), softwareName))

	ginutil.Success(ctx, &dto.DeleteSoftwareResponse{})
}

// PublishSoftware
//
//	@Summary		发布可视化应用
//	@Description	发布可视化应用接口
//	@Tags			可视化-软件
//	@Accept			json
//	@Produce		json
//	@Param			param	body		dto.PublishSoftwareRequest	true	"请求参数"
//	@Response		200		{object}	dto.PublishSoftwareResponse
//	@Router			/vis/software/publish [put]
func (s *RouteService) PublishSoftware(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	req := &dto.PublishSoftwareRequest{}
	if err := ctx.ShouldBindJSON(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	if req.Id == "" || req.State != common.Published && req.State != common.Unpublished {
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	tracelog.Info(ctx, fmt.Sprintf("user: [%v] pulish software req: [%+v]", ginutil.GetUserID(ctx), req))

	err := s.visualService.PublishSoftware(ctx, req.Id, req.State)
	if err != nil {
		logger.Errorf("publish software err: %v, req: [%+v]", err, req)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrVisualPublishSoftwareFailed)
		return
	}

	publishState := consts.PublishMap[req.State]
	software, _ := s.visualService.GetSoftware(ctx, snowflake.MustParseString(req.Id))
	var softwareName string
	if software != nil {
		softwareName = software.Name
	}
	oplog.GetInstance().SaveAuditLogInfo(ctx, approve.OperateTypeEnum_VIS_MANAGER, fmt.Sprintf("用户%v%v镜像[%v]", ginutil.GetUserName(ctx), publishState, softwareName))

	ginutil.Success(ctx, &dto.PublishSoftwareResponse{})
}

// ListSoftwareUseStatuses
//
//	@Summary		软件使用状态列表
//	@Description	软件使用状态列表接口
//	@Tags			可视化-软件
//	@Accept			json
//	@Produce		json
//	@Param			param	query		dto.ListSoftwareUsingStatusesRequest	true	"请求参数"
//	@Response		200		{object}	dto.ListSoftwareUsingStatusesResponse
//	@Router			/vis/software/usingStatuses [get]
func (s *RouteService) ListSoftwareUseStatuses(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	userID := ginutil.GetUserID(ctx)
	username := ginutil.GetUserName(ctx)

	usingStatuses, err := s.visualService.SoftwareUseStatuses(ctx, userID, username)
	if err != nil {
		logger.Errorf("list software using statuses err: %v", err)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrVisualListSoftwareUsingStatusesFailed)
		return
	}

	ginutil.Success(ctx, &dto.ListSoftwareUsingStatusesResponse{UsingStatuses: usingStatuses})
}

// GetSoftwarePresets
//
//	@Summary		获取软件预设
//	@Description	获取软件预设接口
//	@Tags			可视化-软件
//	@Accept			json
//	@Produce		json
//	@Param			param	query		dto.GetSoftwarePresetsRequest	true	"请求参数"
//	@Response		200		{object}	dto.GetSoftwarePresetsResponse
//	@Router			/vis/software/preset [get]
func (s *RouteService) GetSoftwarePresets(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	req := &dto.GetSoftwarePresetsRequest{}
	if err := ctx.BindQuery(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}
	if req.SoftwareID == "" {
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	presets, err := s.visualService.GetSoftwarePresets(ctx, req.SoftwareID)
	if err != nil {
		logger.Errorf("get software presets err: %v, req: [%+v]", err, req)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrVisualGetSoftwarePresetsFailed)
		return
	}
	ginutil.Success(ctx, &dto.GetSoftwarePresetsResponse{Presets: presets})
}

// SetSoftwarePresets
//
//	@Summary		设置软件预设
//	@Description	设置软件预设接口
//	@Tags			可视化-软件
//	@Accept			json
//	@Produce		json
//	@Param			param	body		dto.SetSoftwarePresetsRequest	true	"请求参数"
//	@Response		200		{object}	dto.SetSoftwarePresetsResponse
//	@Router			/vis/software/preset [post]
func (s *RouteService) SetSoftwarePresets(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	req := &dto.SetSoftwarePresetsRequest{}
	if err := ctx.BindJSON(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}
	if req.SoftwareID == "" || req.Presets == nil {
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	tracelog.Info(ctx, fmt.Sprintf("user: [%v] set software presets req: [%+v]", ginutil.GetUserID(ctx), req))

	err := s.visualService.SetSoftwarePresets(ctx, req.SoftwareID, req.Presets)
	if err != nil {
		logger.Errorf("set software presets err: %v, req: [%+v]", err, req)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrVisualSetSoftwarePresetsFailed)
		return
	}

	software, _ := s.visualService.GetSoftware(ctx, snowflake.MustParseString(req.SoftwareID))
	var softwareName string
	if software != nil {
		softwareName = software.Name
	}
	oplog.GetInstance().SaveAuditLogInfo(ctx, approve.OperateTypeEnum_VIS_MANAGER, fmt.Sprintf("用户%v给[%v]镜像设置预设", ginutil.GetUserName(ctx), softwareName))

	ginutil.Success(ctx, &dto.SetSoftwarePresetsResponse{})
}

// AddRemoteApp
//
//	@Summary		新增远程应用
//	@Description	新增远程应用接口
//	@Tags			可视化-软件
//	@Accept			json
//	@Produce		json
//	@Param			param	body		dto.AddRemoteAppRequest	true	"请求参数"
//	@Response		200		{object}	dto.AddRemoteAppResponse
//	@Router			/vis/software/remote/app [post]
func (s *RouteService) AddRemoteApp(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	req := &dto.AddRemoteAppRequest{}
	if err := ctx.BindJSON(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	if req.SoftwareID == "" || req.Name == "" || req.BaseURL == "" {
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	id, err := s.visualService.AddRemoteApp(ctx, req.SoftwareID, &dto.RemoteApp{
		Name:       req.Name,
		Desc:       req.Desc,
		BaseURL:    req.BaseURL,
		Dir:        req.Dir,
		Args:       req.Args,
		Logo:       req.Logo,
		DisableGfx: req.DisableGFX,
	})
	if err != nil {
		logger.Errorf("add remote app err: %v, req: [%+v]", err, req)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrVisualAddRemoteAppFailed)
		return
	}
	ginutil.Success(ctx, &dto.AddRemoteAppResponse{ID: id})
}

// UpdateRemoteApp
//
//	@Summary		更新远程应用
//	@Description	更新远程应用接口
//	@Tags			可视化-软件
//	@Accept			json
//	@Produce		json
//	@Param			param	body		dto.UpdateRemoteAppRequest	true	"请求参数"
//	@Response		200		{object}	dto.UpdateRemoteAppResponse
//	@Router			/vis/software/remote/app [put]
func (s *RouteService) UpdateRemoteApp(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	req := &dto.UpdateRemoteAppRequest{}
	if err := ctx.BindJSON(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	if req.ID == "" || req.SoftwareID == "" || req.Name == "" || req.BaseURL == "" {
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	err := s.visualService.UpdateRemoteApp(ctx, req.SoftwareID, &dto.RemoteApp{
		ID:         req.ID,
		Name:       req.Name,
		Desc:       req.Desc,
		BaseURL:    req.BaseURL,
		Dir:        req.Dir,
		Args:       req.Args,
		Logo:       req.Logo,
		DisableGfx: req.DisableGFX,
	})
	if err != nil {
		logger.Errorf("update remote app err: %v, req: [%+v]", err, req)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrVisualUpdateRemoteAppFailed)
		return
	}
	ginutil.Success(ctx, &dto.UpdateRemoteAppResponse{})
}

// DeleteRemoteApp
//
//	@Summary		删除远程应用
//	@Description	删除远程应用接口
//	@Tags			可视化-软件
//	@Accept			json
//	@Produce		json
//	@Param			param	query		dto.DeleteRemoteAppRequest	true	"请求参数"
//	@Response		200		{object}	dto.DeleteRemoteAppResponse
//	@Router			/vis/software/remote/app [delete]
func (s *RouteService) DeleteRemoteApp(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	req := &dto.DeleteRemoteAppRequest{}
	if err := ctx.BindQuery(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	if req.ID == "" {
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	err := s.visualService.DeleteRemoteApp(ctx, req.ID)
	if err != nil {
		logger.Errorf("delete remote app err: %v, req: [%+v]", err, req)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrVisualDeleteRemoteAppFailed)
		return
	}
	ginutil.Success(ctx, &dto.DeleteRemoteAppResponse{})
}
