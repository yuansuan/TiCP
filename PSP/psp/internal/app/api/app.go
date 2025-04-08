package api

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"

	mainconfig "github.com/yuansuan/ticp/PSP/psp/cmd/config"
	"github.com/yuansuan/ticp/PSP/psp/internal/app/consts"
	"github.com/yuansuan/ticp/PSP/psp/internal/app/dto"
	"github.com/yuansuan/ticp/PSP/psp/internal/common"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/oplog"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/approve"
	"github.com/yuansuan/ticp/PSP/psp/pkg/tracelog"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/ginutil"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/strutil"
)

// ListApp
//
//	@Summary		计算应用列表
//	@Description	计算应用列表接口
//	@Tags			计算应用
//	@Accept			json
//	@Produce		json
//	@Param			param	query		dto.ListAppRequest	true	"请求参数"
//	@Response		200		{object}	dto.ListAppResponse
//	@Router			/app/list [get]
func (s *RouteService) ListApp(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	req := &dto.ListAppRequest{}
	if err := ctx.BindQuery(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	userID := ginutil.GetUserID(ctx)

	templates, err := s.appService.ListApp(ctx, userID, req.ComputeType, req.State, req.HasPermission, req.Desktop)
	if err != nil {
		logger.Errorf("list app err: %v, req: [%+v]", err, req)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrAppListAppFailed)
		return
	}
	ginutil.Success(ctx, &dto.ListAppResponse{Apps: templates})
}

// ListTemplate
//
//	@Summary		计算应用模版列表
//	@Description	计算应用模版列表接口
//	@Tags			计算应用
//	@Accept			json
//	@Produce		json
//	@Param			param	query		dto.ListTemplateRequest	true	"请求参数"
//	@Response		200		{object}	dto.ListTemplateResponse
//	@Router			/app/template/list [get]
func (s *RouteService) ListTemplate(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	req := &dto.ListTemplateRequest{}
	if err := ctx.BindQuery(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}
	templates, err := s.appService.ListApp(ctx, 0, req.ComputeType, req.State, false, false)
	if err != nil {
		logger.Errorf("list app template err: %v, req: [%+v]", err, req)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrAppListTemplateFailed)
		return
	}
	ginutil.Success(ctx, &dto.ListTemplateResponse{Apps: templates})
}

// GetAppInfo
//
//	@Summary		计算应用模版详情
//	@Description	计算应用模版详情接口
//	@Tags			计算应用
//	@Accept			json
//	@Produce		json
//	@Param			param	query		dto.GetAppInfoRequest	true	"请求参数"
//	@Response		200		{object}	dto.GetAppInfoResponse
//	@Router			/app/template [get]
func (s *RouteService) GetAppInfo(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	req := &dto.GetAppInfoRequest{}
	if err := ctx.BindQuery(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	if req.Name == "" || req.ComputeType == "" {
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	app, err := s.appService.GetAppInfo(ctx, &dto.GetAppInfoServiceRequest{
		Name:        req.Name,
		ComputeType: req.ComputeType,
	})
	if err != nil {
		logger.Errorf("get app info err: %v, req: [%+v]", err, req)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrAppGetAppInfoFailed)
		return
	}
	ginutil.Success(ctx, &dto.GetAppInfoResponse{App: app})
}

// AddApp
//
//	@Summary		新增计算应用模版
//	@Description	新增计算应用模版接口
//	@Tags			计算应用
//	@Accept			json
//	@Produce		json
//	@Param			param	body		dto.AddAppRequest	true	"请求参数"
//	@Response		200		{object}	dto.AddAppResponse
//	@Router			/app/template [post]
func (s *RouteService) AddApp(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	req := &dto.AddAppRequest{}
	if err := ctx.ShouldBindJSON(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	if req.NewType == "" || req.NewVersion == "" {
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	tracelog.Info(ctx, fmt.Sprintf("user: [%v] add app type: [%v], version: [%v]", ginutil.GetUserID(ctx), req.NewType, req.NewVersion))

	err := s.appService.AddApp(ctx, &dto.AddAppServiceRequest{
		NewType:           req.NewType,
		NewVersion:        req.NewVersion,
		BaseName:          req.BaseName,
		ComputeType:       req.ComputeType,
		Description:       req.Description,
		Icon:              req.Icon,
		Image:             req.Image,
		Licenses:          req.Licenses,
		ResidualLogParser: req.ResidualLogParser,
		CloudOutAppId:     req.CloudOutAppID,
		EnableResidual:    req.EnableResidual,
		EnableSnapshot:    req.EnableSnapshot,
		Queues:            req.Queues,
		BinPath:           req.BinPath,
		SchedulerParam:    req.SchedulerParam,
	})
	if err != nil {
		logger.Errorf("add app err: %v, req: [%+v]", err, req)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrAppAddAppFailed)
		return
	}

	oplog.GetInstance().SaveAuditLogInfo(ctx, approve.OperateTypeEnum_APP_MANAGER, fmt.Sprintf("用户%v新建本地应用模版[%v]", ginutil.GetUserName(ctx), fmt.Sprintf("%v%v%v", req.NewType, common.Blank, req.NewVersion)))

	ginutil.Success(ctx, &dto.AddAppResponse{})
}

// UpdateApp
//
//	@Summary		更新计算应用模版
//	@Description	更新计算应用模版接口
//	@Tags			计算应用
//	@Accept			json
//	@Produce		json
//	@Param			param	body		dto.UpdateAppRequest	true	"请求参数"
//	@Response		200		{object}	dto.UpdateAppResponse
//	@Router			/app/template [put]
func (s *RouteService) UpdateApp(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	req := &dto.UpdateAppRequest{}
	if err := ctx.ShouldBindJSON(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	if req.App == nil || strings.Contains(req.App.Version, "_") {
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	if req.App != nil {
		tracelog.Info(ctx, fmt.Sprintf("user: [%v] update app: [%v]", ginutil.GetUserID(ctx), req.App.ID))
	}

	err := s.appService.UpdateApp(ctx, req.App, req.BaseName)
	if err != nil {
		logger.Errorf("update app err: %v, req: [%+v]", err, req)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrAppUpdateAppFailed)
		return
	}

	var content string
	if strutil.IsNotEmpty(req.BaseName) {
		content = fmt.Sprintf("用户%v另存为本地应用模版[%v]", ginutil.GetUserName(ctx), fmt.Sprintf("%v%v%v", req.App.Type, common.Blank, req.App.Version))
	} else {
		content = fmt.Sprintf("用户%v更新本地应用模版[%v]", ginutil.GetUserName(ctx), fmt.Sprintf("%v%v%v", req.App.Type, common.Blank, req.App.Version))
	}

	oplog.GetInstance().SaveAuditLogInfo(ctx, approve.OperateTypeEnum_APP_MANAGER, content)

	ginutil.Success(ctx, &dto.UpdateAppResponse{})
}

// DeleteApp
//
//	@Summary		删除计算应用模版
//	@Description	删除计算应用模版接口
//	@Tags			计算应用
//	@Accept			json
//	@Produce		json
//	@Param			param	query		dto.DeleteAppRequest	true	"请求参数"
//	@Response		200		{object}	dto.DeleteAppResponse
//	@Router			/app/template [delete]
func (s *RouteService) DeleteApp(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	req := &dto.DeleteAppRequest{}
	if err := ctx.BindQuery(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	if req.Name == "" || req.ComputeType == "" {
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	tracelog.Info(ctx, fmt.Sprintf("user: [%v] delete app req: [%+v]", ginutil.GetUserID(ctx), req))

	err := s.appService.DeleteApp(ctx, req.Name, req.ComputeType)
	if err != nil {
		logger.Errorf("delete app err: %v, req: [%+v]", err, req)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrAppDeleteAppFailed)
		return
	}

	computeType := mainconfig.Custom.Main.ComputeTypeNames[req.ComputeType]
	oplog.GetInstance().SaveAuditLogInfo(ctx, approve.OperateTypeEnum_APP_MANAGER, fmt.Sprintf("用户%v删除%v应用模版[%v]", ginutil.GetUserName(ctx), computeType, req.Name))

	ginutil.Success(ctx, &dto.DeleteAppResponse{})
}

// PublishApp
//
//	@Summary		发布计算应用模版
//	@Description	发布计算应用模版接口
//	@Tags			计算应用
//	@Accept			json
//	@Produce		json
//	@Param			param	body		dto.PublishAppRequest	true	"请求参数"
//	@Response		200		{object}	dto.PublishAppResponse
//	@Router			/app/template/publish [put]
func (s *RouteService) PublishApp(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	req := &dto.PublishAppRequest{}
	if err := ctx.ShouldBindJSON(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	if len(req.Names) == 0 || req.State != common.Published && req.State != common.Unpublished && req.ComputeType != "" {
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	tracelog.Info(ctx, fmt.Sprintf("user: [%v] publish app req: [%+v]", ginutil.GetUserID(ctx), req))

	err := s.appService.PublishApp(ctx, req.Names, req.ComputeType, req.State)
	if err != nil {
		logger.Errorf("publish app err: %v, req: [%+v]", err, req)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrAppPublishAppFailed)
		return
	}

	computeType := mainconfig.Custom.Main.ComputeTypeNames[req.ComputeType]
	publishState := consts.PublishMap[req.State]

	oplog.GetInstance().SaveAuditLogInfo(ctx, approve.OperateTypeEnum_APP_MANAGER, fmt.Sprintf("用户%v%v%v应用模版%v", ginutil.GetUserName(ctx), publishState, computeType, req.Names))

	ginutil.Success(ctx, &dto.PublishAppResponse{})
}

// SyncAppContent
//
//	@Summary		同步应用模版信息到选中的模版
//	@Description	同步应用模版信息到选中的模版接口
//	@Tags			计算应用
//	@Accept			json
//	@Produce		json
//	@Param			param	body		dto.SyncAppContentRequest	true	"请求参数"
//	@Response		200		{object}	dto.SyncAppContentResponse
//	@Router			/app/template/syncAppContent [post]
func (s *RouteService) SyncAppContent(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	req := &dto.SyncAppContentRequest{}
	if err := ctx.ShouldBindJSON(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	if req.BaseAppId == "" || len(req.SyncAppIds) == 0 {
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	tracelog.Info(ctx, fmt.Sprintf("user: [%v] sync app content req: [%+v]", ginutil.GetUserID(ctx), req))

	err := s.appService.SyncAppContent(ctx, req.BaseAppId, req.SyncAppIds)
	if err != nil {
		logger.Errorf("sync app content err: %v", err)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrAppSyncAppContentFailed)
		return
	}

	ginutil.Success(ctx, &dto.SyncAppContentResponse{})
}

// ListZone
//
//	@Summary		区域列表
//	@Description	区域列表接口
//	@Tags			计算应用
//	@Accept			json
//	@Produce		json
//	@Response		200	{object}	[]string
//	@Router			/app/zone [get]
func (s *RouteService) ListZone(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	zoneList, err := s.appService.ListZone(ctx)
	if err != nil {
		logger.Errorf("list zone err: %v", err)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrAppListZoneFailed)
		return
	}
	ginutil.Success(ctx, zoneList)
}

// ListQueue
//
//	@Summary		队列列表
//	@Description	队列列表接口
//	@Tags			计算应用
//	@Accept			json
//	@Produce		json
//	@Param			param	query		dto.ListQueueRequest	true	"请求参数"
//	@Response		200		{object}	dto.ListQueueResponse
//	@Router			/app/queue [get]
func (s *RouteService) ListQueue(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	req := &dto.ListQueueRequest{}
	if err := ctx.ShouldBindQuery(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	queues, err := s.appService.ListQueue(ctx, req.AppId)
	if err != nil {
		logger.Errorf("list queue err: %v", err)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrAppListQueueFailed)
		return
	}

	ginutil.Success(ctx, &dto.ListQueueResponse{Queues: queues})
}

// ListLicense
//
//	@Summary		许可证列表
//	@Description	许可证列表接口
//	@Tags			计算应用
//	@Accept			json
//	@Produce		json
//	@Param			param	query		dto.ListLicenseRequest	true	"请求参数"
//	@Response		200		{object}	dto.ListLicenseResponse
//	@Router			/app/license [get]
func (s *RouteService) ListLicense(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	req := &dto.ListLicenseRequest{}
	if err := ctx.ShouldBindQuery(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	licenses, err := s.appService.ListLicense(ctx)
	if err != nil {
		logger.Errorf("list license err: %v", err)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrAppListLicenseFailed)
		return
	}

	ginutil.Success(ctx, &dto.ListLicenseResponse{Licenses: licenses})
}

// GetSchedulerResourceKey
//
//	@Summary		获取调度器资源 Key 信息
//	@Description	获取调度器资源 Key 信息接口
//	@Tags			计算应用
//	@Accept			json
//	@Produce		json
//	@Param			param	query		dto.GetSchedulerResourceKeyRequest	true	"请求参数"
//	@Response		200		{object}	dto.GetSchedulerResourceKeyResponse
//	@Router			/app/schedulerResourceKey [get]
func (s *RouteService) GetSchedulerResourceKey(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	req := &dto.GetSchedulerResourceKeyRequest{}
	if err := ctx.ShouldBindQuery(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	keys, err := s.appService.GetSchedulerResourceKey(ctx)
	if err != nil {
		logger.Errorf("get scheduler resource key err: %v", err)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrAppGetSchedulerResourceKeyFailed)
		return
	}

	ginutil.Success(ctx, &dto.GetSchedulerResourceKeyResponse{Keys: keys})
}

// GetSchedulerResourceValue
//
//	@Summary		获取调度器资源 Value 信息
//	@Description	获取调度器资源 Value 信息接口
//	@Tags			计算应用
//	@Accept			json
//	@Produce		json
//	@Param			param	query		dto.GetSchedulerResourceValueRequest	true	"请求参数"
//	@Response		200		{object}	dto.GetSchedulerResourceValueResponse
//	@Router			/app/schedulerResourceValue [get]
func (s *RouteService) GetSchedulerResourceValue(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	req := &dto.GetSchedulerResourceValueRequest{}
	if err := ctx.ShouldBindQuery(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	items, err := s.appService.GetSchedulerResourceValue(ctx, req.AppId, req.ResourceType, req.ResourceSubType)
	if err != nil {
		logger.Errorf("get scheduler resource value err: %v", err)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrAppGetSchedulerResourceValueFailed)
		return
	}

	ginutil.Success(ctx, &dto.GetSchedulerResourceValueResponse{Items: items})
}
