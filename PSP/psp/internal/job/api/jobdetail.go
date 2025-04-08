package api

import (
	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/internal/job/dto"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/ginutil"
)

// JobDetail
//
//	@Summary		作业详情
//	@Description	作业详情接口
//	@Tags			作业
//	@Accept			json
//	@Produce		json
//	@Param			param	query		dto.JobDetailRequest	true	"请求参数"
//	@Response		200		{object}	dto.JobDetailInfo
//	@Router			/job/detail [get]
func (r *apiRoute) JobDetail(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	var req = dto.JobDetailRequest{}
	if err := ctx.BindQuery(&req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	detail, err := r.jobService.GetJobDetail(ctx, req.JobID)
	if err != nil {
		logger.Errorf("get job [%v] detail err: %v", req.JobID, err)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrJobFailGet)
		return
	}

	ginutil.Success(ctx, detail)
}

// JobSetDetail
//
//	@Summary		作业集详情
//	@Description	作业集详情接口
//	@Tags			作业
//	@Accept			json
//	@Produce		json
//	@Param			param	query		dto.JobSetDetailRequest	true	"请求参数"
//	@Success		200		{object}	dto.JobSetDetailResponse
//	@Router			/job/jobSetDetail [get]
func (r *apiRoute) JobSetDetail(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	var req = dto.JobSetDetailRequest{}
	if err := ctx.BindQuery(&req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	loginUserID := ginutil.GetUserID(ctx)

	jobSetInfo, jobList, err := r.jobService.GetJobSetDetail(ctx, req.JobSetID, snowflake.ID(loginUserID))
	if err != nil {
		logger.Errorf("get job set: [%v] detail err: %v", req.JobSetID, err)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrJobFailJobSetGet)
		return
	}

	ginutil.Success(ctx, &dto.JobSetDetailResponse{
		JobSetInfo: jobSetInfo,
		JobList:    jobList,
	})
}
