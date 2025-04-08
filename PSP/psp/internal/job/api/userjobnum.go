package api

import (
	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/internal/job/dto"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/ginutil"
)

// UserJobNum 用户作业数
func (r *apiRoute) UserJobNum(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	var req = dto.JobNumRequest{}
	if err := ctx.BindQuery(&req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}
	if req.Start <= 0 || req.End <= 0 {
		logger.Errorf("reqeust params start: %v ,end: %v", req.Start, req.End)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	resp, err := r.jobService.UserJobNum(ctx, req.Start, req.End)
	if err != nil {
		logger.Errorf("get app job num  [start:%v,end:%v] err: %v", req.Start, req.End, err)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrJobFailAppJobNum)
		return
	}

	ginutil.Success(ctx, resp)
}
