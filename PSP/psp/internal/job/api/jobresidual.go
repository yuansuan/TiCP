package api

import (
	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/internal/job/dto"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/ginutil"
)

// JobResidual
//
//	@Summary		作业残差图
//	@Description	作业残差图接口
//	@Tags			作业
//	@Accept			json
//	@Produce		json
//	@Param			param	query		dto.JobResidualRequest	true	"请求参数"
//	@Response		200		{object}	dto.JobResidualResponse
//	@Router			/job/residual [get]
func (r *apiRoute) JobResidual(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	req := &dto.JobResidualRequest{}
	if err := ctx.BindQuery(&req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	response, err := r.jobService.GetJobResidual(ctx, req.JobID)
	if err != nil {
		logger.Errorf("get job [%v] residual err: %v", req.JobID, err)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrJobResidualFailGet)
		return
	}

	ginutil.Success(ctx, response)
}
