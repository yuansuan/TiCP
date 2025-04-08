package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/oplog"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/approve"
	"github.com/yuansuan/ticp/PSP/psp/internal/job/dto"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/ginutil"
)

// JobTerminate
//
//	@Summary		作业终止
//	@Description	作业终止接口
//	@Tags			作业
//	@Accept			json
//	@Produce		json
//	@Param			param	query		dto.JobTerminateRequest	true	"请求参数"
//	@Response		200		{object}	nil
//	@Router			/job/terminate [post]
func (r *apiRoute) JobTerminate(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	var req = dto.JobTerminateRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	if err := r.jobService.JobTerminate(ctx, req.OutJobID, req.ComputeType); err != nil {
		logger.Errorf("terminate job [%v] err: %v", req.OutJobID, err)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrJobFailTerminate)
		return
	}

	job, _ := r.jobService.GetJobDetailByOutID(ctx, req.OutJobID, req.ComputeType)
	var jobName string
	if job != nil {
		jobName = job.Name
	}
	oplog.GetInstance().SaveAuditLogInfo(ctx, approve.OperateTypeEnum_JOB_MANAGER, fmt.Sprintf("用户%v终止求解作业[%v]", ginutil.GetUserName(ctx), jobName))

	ginutil.Success(ctx, nil)
}
