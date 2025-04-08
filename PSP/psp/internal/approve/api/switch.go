package api

import (
	"github.com/gin-gonic/gin"

	"github.com/yuansuan/ticp/PSP/psp/internal/approve/config"
	"github.com/yuansuan/ticp/PSP/psp/internal/approve/dto"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/ginutil"
)

// ThreePersonManagement
//
//	@Summary		查询三员管理开关状态
//	@Description	查询三员管理开关状态
//	@Tags			审批
//	@Produce		json
//	@Param			param	query		dto.ThreeStateRequest	true	"入参"
//	@Response		200		{object}	dto.ThreeStateResponse
//	@Router			/approve/threePersonManagement [get]
func (s *RouteService) ThreePersonManagement(ctx *gin.Context) {

	if config.GetConfig().ThreePersonManagement {
		ginutil.Success(ctx, dto.ThreeStateResponse{
			State: true,
		})
		return
	}

	ginutil.Success(ctx, dto.ThreeStateResponse{
		State: false,
	})
	return
}
