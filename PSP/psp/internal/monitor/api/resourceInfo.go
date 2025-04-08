package api

import (
	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/internal/monitor/dto"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/ginutil"
)

// ResourceInfo
//
//	@Summary		节点资源信息, 包括：过去24小时CPU利用率， 过去24小时内存利用率，过去24小时磁盘IO速率
//	@Description	节点基础资源信息
//	@Tags			集群监控
//	@Accept			json
//	@Produce		json
//	@Param			param	body		dto.Request	true	"请求参数"
//	@Response		200		{object}	dto.ResourceResponse
//	@Router			/dashboard/resourceInfo [get]
func (r *apiRoute) ResourceInfo(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	var req = dto.Request{}
	if err := ctx.BindQuery(&req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	cpuUtAvg, ioUtAvg, memUtAvg, err := r.dashboardService.GetResourceInfo(ctx, &req)
	if err != nil {
		logger.Errorf("node detail err: %v", err)
		ginutil.Error(ctx, errcode.ErrResourceInfoFail, errcode.MsgResourceInfoFail)
		return
	}

	ginutil.Success(ctx, dto.ResourceResponse{
		MetricCpuUtAvg: cpuUtAvg,
		MetricIoUtAvg:  ioUtAvg,
		MetricMemUtAvg: memUtAvg,
	})
}
