package api

import (
	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/ginutil"
)

// LicenseTypeList
//
//	@Summary		许可证类型下拉框
//	@Description	许可证id + name
//	@Tags			许可证监控
//	@Accept			json
//	@Produce		json
//	@Response		200	{object}	dto.LicenseTypeListResponse
//	@Router			/licenseManagers/typeList [get]
func (r *apiRoute) LicenseTypeList(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	resp, err := r.licManagerService.LicenseTypeList(ctx)
	if err != nil {
		logger.Errorf("get license manager list err: %v", err)
		ginutil.Error(ctx, errcode.ErrFailedConfigModuleList, errcode.MsgFailedConfigModuleList)
		return
	}

	ginutil.Success(ctx, resp)
}
