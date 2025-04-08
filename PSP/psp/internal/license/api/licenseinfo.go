package api

import (
	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"google.golang.org/grpc/status"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/internal/license/dto"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/ginutil"
)

// AddLicenseInfo
//
//	@Summary		新增许可证信息
//	@Description	包括许可证名称、地址、端口、MAC地址、有效时间、工具路径等
//	@Tags			许可证监控
//	@Accept			json
//	@Produce		json
//	@Param			param	query		dto.LicenseInfoAddRequest	true	"请求参数"
//	@Response		200		{object}	dto.LicenseInfoAddResponse
//	@Router			/licenseInfos [post]
func (r *apiRoute) AddLicenseInfo(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	req := dto.LicenseInfoAddRequest{}
	if err := ctx.BindJSON(&req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	resp, err := r.licInfoService.AddLicenseInfo(ctx, &req)
	if err != nil {
		logger.Errorf("add license info error! err:%v", err)
		if status.Code(err) == errcode.ErrFailedLicenseNameRepeat {
			ginutil.Error(ctx, errcode.ErrFailedLicenseNameRepeat, errcode.MsgFailedLicenseNameRepeat)
		} else {
			ginutil.Error(ctx, errcode.ErrFailedLicenseInfoAdd, errcode.MsgFailedLicenseInfoAdd)
		}
		return
	}

	ginutil.Success(ctx, resp)
}

// EditLicenseInfo
//
//	@Summary		修改许可证信息
//	@Description	包括许可证名称、地址、端口、MAC地址、有效时间、工具路径等
//	@Tags			许可证监控
//	@Accept			json
//	@Produce		json
//	@Param			param	body	dto.LicenseInfoEditRequest	true	"请求参数"
//	@Router			/licenseInfos/:id [put]
func (r *apiRoute) EditLicenseInfo(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	req := dto.LicenseInfoEditRequest{}
	if err := ctx.BindJSON(&req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	err := r.licInfoService.EditLicenseInfo(ctx, &req)
	if err != nil {
		logger.Errorf("edit license info error! err:%v", err)
		if status.Code(err) == errcode.ErrFailedLicenseNameRepeat {
			ginutil.Error(ctx, errcode.ErrFailedLicenseNameRepeat, errcode.MsgFailedLicenseNameRepeat)
		} else {
			ginutil.Error(ctx, errcode.ErrFailedLicenseInfoEdit, errcode.MsgFailedLicenseInfoEdit)
		}
		return
	}

	ginutil.Success(ctx, nil)
}

// DeleteLicenseInfo
//
//	@Summary		修改许可证信息
//	@Description	包括许可证名称、地址、端口、MAC地址、有效时间、工具路径等
//	@Tags			许可证监控
//	@Accept			json
//	@Produce		json
//	@Router			/licenseInfos/:id [delete]
func (r *apiRoute) DeleteLicenseInfo(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	id, ok := getResourceId(ctx)
	if !ok {
		return
	}

	err := r.licInfoService.DeleteLicenseInfo(ctx, id)
	if err != nil {
		logger.Errorf("delete license info error! err:%v", err)
		ginutil.Error(ctx, errcode.ErrFailedLicenseInfoDel, errcode.MsgFailedLicenseInfoDel)
		return
	}

	ginutil.Success(ctx, nil)
}
