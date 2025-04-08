package api

import (
	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/internal/job/dto"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/ginutil"
)

// JobSnapshotList
//
//	@Summary		作业云图集
//	@Description	作业云图集接口
//	@Tags			作业
//	@Accept			json
//	@Produce		json
//	@Param			param	query		dto.JobSnapshotListRequest	true	"请求参数"
//	@Response		200		{object}	dto.JobSnapshotListResponse
//	@Router			/job/snapshots [get]
func (r *apiRoute) JobSnapshotList(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	req := &dto.JobSnapshotListRequest{}
	if err := ctx.BindQuery(&req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	response, err := r.jobService.GetJobSnapshotList(ctx, req.JobID)
	if err != nil {
		logger.Errorf("get job [%v] snapshots err: %v", req.JobID, err)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrJobSnapshotFailList)
		return
	}

	ginutil.Success(ctx, response)
}

// JobSnapshot
//
//	@Summary		作业云图
//	@Description	作业云图接口
//	@Tags			作业
//	@Accept			json
//	@Produce		json
//	@Param			param	query		dto.JobSnapshotRequest	true	"请求参数"
//	@Response		200		{object}	dto.JobSnapshotResponse
//	@Router			/job/snapshot [get]
func (r *apiRoute) JobSnapshot(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	req := &dto.JobSnapshotRequest{}
	if err := ctx.BindQuery(&req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	response, err := r.jobService.GetJobSnapshot(ctx, req.JobID, req.Path)
	if err != nil {
		logger.Errorf("get job [%v] snapshot [%v] err: %v", req.JobID, req.Path, err)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrJobSnapshotFailGet)
		return
	}

	ginutil.Success(ctx, response)
}
