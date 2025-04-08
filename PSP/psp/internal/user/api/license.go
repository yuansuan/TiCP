package api

import (
	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/internal/user/dto"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/ginutil"
)

// GetMachineID
//
//	@Summary		获取机器 ID
//	@Description	获取机器 ID 接口
//	@Tags			用户-系统许可证
//	@Accept			json
//	@Produce		json
//	@Param			param	query		dto.GetMachineIDRequest	true	"请求参数"
//	@Response		200		{object}	dto.GetMachineIDResponse
//	@Router			/auth/license/machineID [get]
func (s *RouteService) GetMachineID(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	req := &dto.GetMachineIDRequest{}
	if err := ctx.BindQuery(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	machineID, err := s.LicenseService.GetMachineID(ctx)
	if err != nil {
		logger.Errorf("get machine id err: %v, req: [%+v]", err, req)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrAuthGetMachineIDFailed)
		return
	}
	ginutil.Success(ctx, &dto.GetMachineIDResponse{ID: machineID})
}

// GetLicense
//
//	@Summary		获取系统许可证详情
//	@Description	获取系统许可证详情接口
//	@Tags			用户-系统许可证
//	@Accept			json
//	@Produce		json
//	@Param			param	query		dto.GetLicenseInfoRequest	true	"请求参数"
//	@Response		200		{object}	dto.GetLicenseInfoResponse
//	@Router			/auth/license [get]
func (s *RouteService) GetLicense(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	req := &dto.GetLicenseInfoRequest{}
	if err := ctx.BindQuery(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	license, err := s.LicenseService.GetLicense(ctx)
	if err != nil {
		logger.Errorf("get license info err: %v, req: [%+v]", err, req)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrAuthGetLicenseInfoFailed)
		return
	}
	ginutil.Success(ctx, &dto.GetLicenseInfoResponse{
		Name:          license.Name,
		Version:       license.Version,
		Expiry:        license.Expiry,
		AvailableDays: license.AvailableDays,
	})
}

// UpdateLicense
//
//	@Summary		更新系统许可证信息
//	@Description	更新系统许可证信息接口
//	@Tags			用户-系统许可证
//	@Accept			json
//	@Produce		json
//	@Param			param	body		dto.UpdateLicenseInfoRequest	true	"请求参数"
//	@Response		200		{object}	dto.UpdateLicenseInfoResponse
//	@Router			/auth/license [post]
func (s *RouteService) UpdateLicense(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	req := &dto.UpdateLicenseInfoRequest{}
	if err := ctx.ShouldBindJSON(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	if req.Name == "" || req.Version == "" || req.Expiry == "" || req.MachineID == "" || req.Key == "" {
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	err := s.LicenseService.UpdateLicense(ctx, &dto.License{
		Name:      req.Name,
		Version:   req.Version,
		Expiry:    req.Expiry,
		MachineID: req.MachineID,
		Key:       req.Key,
	})
	if err != nil {
		logger.Errorf("update license info err: %v, req: [%+v]", err, req)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrAuthUpdateLicenseInfoFailed)
		return
	}
	ginutil.Success(ctx, &dto.UpdateLicenseInfoResponse{})
}
