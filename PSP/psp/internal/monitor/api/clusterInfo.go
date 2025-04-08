package api

import (
	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/internal/monitor/dto"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/ginutil"
)

// ClusterInfo
//
//	@Summary		集群信息, 包括：集群总核数和节点数、节点列表、共享存储
//	@Description	集群相关信息
//	@Tags			集群监控
//	@Accept			json
//	@Produce		json
//	@Response		200	{object}	dto.ClusterResponse
//	@Router			/dashboard/clusterInfo [get]
func (r *apiRoute) ClusterInfo(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	info, details, disk, err := r.dashboardService.GetClusterInfo(ctx)
	if err != nil {
		logger.Errorf("node detail err: %v", err)
		ginutil.Error(ctx, errcode.ErrClusterInfoFail, errcode.MsgClusterInfoFail)
		return
	}

	ginutil.Success(ctx, dto.ClusterResponse{
		ClusterInfo: info,
		NodeList:    details,
		Disks:       disk,
	})
}
