package api

import (
	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/internal/job/dto"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/ginutil"
)

// Resubmit
//
//	@Summary		作业重新提交
//	@Description	作业重新提交接口(准备算力文件以及返回上次作业提交的参数)
//	@Tags			作业
//	@Accept			json
//	@Produce		json
//	@Param			param	body		dto.ResubmitRequest	true	"请求参数"
//	@Response		200		{object}	dto.ResubmitResponse
//	@Router			/job/resubmit [post]
func (r *apiRoute) Resubmit(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	req := &dto.ResubmitRequest{}
	if err := ctx.BindJSON(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	if req.JobId == "" {
		logger.Errorf("job id empty")
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	userID := ginutil.GetUserID(ctx)
	userName := ginutil.GetUserName(ctx)

	res, err := r.jobService.JobResubmit(ctx, req, snowflake.ID(userID), userName)
	if err != nil {
		logger.Errorf("job resubmit err: %v", err)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrJobResubmitFailed)
		return
	}

	ginutil.Success(ctx, res)
}
