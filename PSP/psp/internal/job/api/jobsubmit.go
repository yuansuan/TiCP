package api

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/internal/job/dto"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/ginutil"
)

// Submit
//
//	@Summary		作业提交
//	@Description	作业提交接口
//	@Tags			作业
//	@Accept			json
//	@Produce		json
//	@Param			param	query		dto.JobSubmitRequest	true	"请求参数"
//	@Response		200		{object}	dto.JobSubmitResponse
//	@Router			/job/submit [post]
func (r *apiRoute) Submit(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	var req = dto.JobSubmitRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	userID := ginutil.GetUserID(ctx)
	userName := ginutil.GetUserName(ctx)
	submitParam := &dto.SubmitParam{
		AppID:     req.AppID,
		ProjectID: req.ProjectID,
		UserID:    snowflake.ID(userID),
		UserName:  userName,
		QueueName: req.QueueName,
		MainFiles: req.MainFiles,
		WorkDir:   req.WorkDir,
		Fields:    req.Fields,
	}
	_, err := r.jobService.JobSubmit(ctx, submitParam)
	if err != nil {
		logger.Errorf("submit job err: %v", err)
		if strings.Contains(err.Error(), "InvalidAccountStatus.NotEnoughBalance") {
			errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrJobFailNotEnoughBalance)
		} else {
			errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrJobFailSubmit)
		}

		return
	}

	ginutil.Success(ctx, &dto.JobSubmitResponse{})
}
