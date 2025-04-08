package api

import (
	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/internal/monitor/dto"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/ginutil"
)

// NodeCoreNum
//
//	@Summary		集群中节点的总核数及空闲核数
//	@Description	获取集群的总核数及空闲核数
//	@Tags			节点管理
//	@Accept			json
//	@Produce		json
//	@Response		200	{object}	dto.CoreStatisticsResponse
//	@Router			/node/coreNum [get]
func (r *apiRoute) NodeCoreNum(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	resp, err := r.nodeService.NodeCoreNum(ctx)
	if err != nil {
		logger.Errorf("node coreNum err: %v", err)
		ginutil.Error(ctx, errcode.ErrNodeFailCoreNum, errcode.MsgNodeFailCoreNum)
		return
	}

	ginutil.Success(ctx, dto.CoreStatisticsResponse{
		CoreStatistics: resp,
	})
}
