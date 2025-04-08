package api

import (
	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/internal/monitor/dto"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/ginutil"
)

// JobInfo
//
//	@Summary		作业状态数量统计, 包括：过去24小时作业状态数量及实时数量、过去24小时应用维度作业数量、过去24小时用户维度作业数量
//	@Description	作业数量统计
//	@Tags			集群监控
//	@Accept			json
//	@Produce		json
//	@Param			param	body		dto.Request	true	"请求参数"
//	@Response		200		{object}	dto.JobResponse
//	@Router			/dashboard/jobInfo [get]
func (r *apiRoute) JobInfo(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	var req = dto.Request{}
	if err := ctx.BindQuery(&req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	jobResRange, jobResLatest, err := r.dashboardService.GetJobInfo(ctx, &req)
	if err != nil {
		logger.Errorf("node detail err: %v", err)
		ginutil.Error(ctx, errcode.ErrJobStatusInfoFail, errcode.MsgJobStatusInfoFail)
		return
	}

	ginutil.Success(ctx, dto.JobResponse{
		JobResRange:  jobResRange,
		JobResLatest: jobResLatest,
	})
}
