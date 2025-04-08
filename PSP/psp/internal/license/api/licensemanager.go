package api

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"google.golang.org/grpc/status"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/oplog"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/approve"
	"github.com/yuansuan/ticp/PSP/psp/internal/license/consts"
	"github.com/yuansuan/ticp/PSP/psp/internal/license/dto"
	"github.com/yuansuan/ticp/PSP/psp/internal/license/util"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/ginutil"
)

// ListLicenseManager
//
//	@Summary		license manager列表
//	@Description	license manager列表接口
//	@Tags			许可证监控
//	@Accept			json
//	@Produce		json
//	@Param			param	query		dto.LicenseManagerListRequest	true	"请求参数"
//	@Response		200		{object}	dto.LicenseManagerListResponse
//	@Router			/licenseManagers [get]
func (r *apiRoute) ListLicenseManager(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	var req = dto.LicenseManagerListRequest{}
	if err := ctx.BindQuery(&req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	resp, err := r.licManagerService.LicenseManagerList(ctx, req.LicenseType)
	if err != nil {
		logger.Errorf("get license manager list err: %v", err)
		ginutil.Error(ctx, errcode.ErrFailedLicenseManagerList, errcode.MsgFailedLicenseManagerList)
		return
	}

	ginutil.Success(ctx, resp)
}

// LicenseManageInfo
//
//	@Summary		详情, 包括：许可证类型信息、许可证详情、模块信息
//	@Description	查看许可证信息
//	@Tags			许可证监控
//	@Accept			json
//	@Produce		json
//	@Response		200	{object}	dto.LicenseManagerData
//	@Router			/licenseManagers/:id [get]
func (r *apiRoute) LicenseManageInfo(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	id, ok := getResourceId(ctx)
	if !ok {
		return
	}

	resp, err := r.licManagerService.LicenseManagerInfo(ctx, id)
	if err != nil {
		logger.Errorf("get license manager list err: %v", err)
		ginutil.Error(ctx, errcode.ErrFailedLicenseManagerList, errcode.MsgFailedLicenseManagerList)
		return
	}

	ginutil.Success(ctx, resp)
}

func getResourceId(ctx *gin.Context) (string, bool) {
	id := ctx.Param("id")
	if id == "" {
		logging.GetLogger(ctx).Warnf("invalid params invalid id")
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return "", false
	}
	return id, true
}

// AddLicenseManager
//
//	@Summary		新增license manager
//	@Description	新增license manager
//	@Tags			许可证监控
//	@Accept			json
//	@Produce		json
//	@Param			param	query		dto.AddLicenseManagerRequest	true	"请求参数"
//	@Response		200		{object}	dto.AddLicenseManagerResponse
//	@Router			/licenseManagers [post]
func (r *apiRoute) AddLicenseManager(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	req := dto.AddLicenseManagerRequest{}
	if err := ctx.BindJSON(&req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	resp, err := r.licManagerService.AddLicenseManager(ctx, &req)
	if err != nil {
		logger.Errorf("add license manager error! err:%v", err)
		if status.Code(err) == errcode.ErrFailedAppTypeRepeat {
			ginutil.Error(ctx, errcode.ErrFailedAppTypeRepeat, errcode.MsgFailedAppTypeRepeat)
		} else {
			ginutil.Error(ctx, errcode.ErrFailedLicenseManagerAdd, errcode.MsgFailedLicenseManagerAdd)
		}

		return
	}

	oplog.GetInstance().SaveAuditLogInfo(ctx, approve.OperateTypeEnum_LICENSE_MANAGER, fmt.Sprintf("用户%v添加许可证[%v]", ginutil.GetUserName(ctx), resp.Id))

	ginutil.Success(ctx, resp)
}

// EditLicenseManager
//
//	@Summary		编辑license manager
//	@Description	编辑license manager
//	@Tags			许可证监控
//	@Accept			json
//	@Produce		json
//	@Param			param	query	dto.EditLicenseManagerRequest	true	"请求参数"
//	@Router			/licenseManagers/:id [put]
func (r *apiRoute) EditLicenseManager(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	req := dto.EditLicenseManagerRequest{}
	if err := ctx.BindJSON(&req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	id, ok := util.GetResourceId(ctx)
	if !ok {
		return
	} else {
		req.Id = id
	}

	err := r.licManagerService.EditLicenseManager(ctx, &req)
	if err != nil {
		logger.Errorf("add license manager error! err:%v", err)
		if status.Code(err) == errcode.ErrFailedAppTypeRepeat {
			ginutil.Error(ctx, errcode.ErrFailedAppTypeRepeat, errcode.MsgFailedAppTypeRepeat)
		} else {
			ginutil.Error(ctx, errcode.ErrFailedLicenseManagerEdit, errcode.MsgFailedLicenseManagerEdit)
		}
		return
	}

	oplog.GetInstance().SaveAuditLogInfo(ctx, approve.OperateTypeEnum_LICENSE_MANAGER, fmt.Sprintf("用户%v编辑许可证[%v]", ginutil.GetUserName(ctx), req.Id))

	ginutil.Success(ctx, nil)
}

// DeleteLicenseManager
//
//	@Summary		删除license manager
//	@Description	删除license manager
//	@Tags			许可证监控
//	@Accept			json
//	@Produce		json
//	@Router			/licenseManagers/:id [delete]
func (r *apiRoute) DeleteLicenseManager(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	id, ok := util.GetResourceId(ctx)
	if !ok {
		return
	}

	err := r.licManagerService.DeleteLicenseManager(ctx, id)
	if err != nil {
		logger.Errorf("delete license manager error! err:%v", err)
		if status.Code(err) == errcode.ErrFailedLicenseManagerDeleteBind {
			ginutil.Error(ctx, errcode.ErrFailedLicenseManagerDeleteBind, errcode.MsgFailedLicenseManagerDeleteBind)
		} else if strings.Contains(fmt.Sprintf("%v", err), consts.DeleteLicenseManagerPaasErrCode) {
			ginutil.Error(ctx, errcode.ErrFailedLicenseManagerDeleteExist, errcode.MsgFailedLicenseManagerDeleteExist)
		} else {
			ginutil.Error(ctx, errcode.ErrFailedLicenseManagerDelete, errcode.MsgFailedLicenseManagerDelete)
		}

		return
	}

	oplog.GetInstance().SaveAuditLogInfo(ctx, approve.OperateTypeEnum_LICENSE_MANAGER, fmt.Sprintf("用户%v删除许可证[%v]", ginutil.GetUserName(ctx), id))

	ginutil.Success(ctx, nil)
}
