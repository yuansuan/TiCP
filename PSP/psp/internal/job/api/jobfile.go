package api

import (
	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/internal/job/dto"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/ginutil"
)

// CreateTempDir
//
//	@Summary		创建作业临时目录
//	@Description	创建作业临时目录接口
//	@Tags			作业
//	@Accept			json
//	@Produce		json
//	@Param			param	query		dto.CreateTempDirRequest	true	"请求参数"
//	@Response		200		{object}	dto.CreateTempDirResponse
//	@Router			/job/createTempDir [post]
func (r *apiRoute) CreateTempDir(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	var req = dto.CreateTempDirRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	userName := ginutil.GetUserName(ctx)
	tempDir, err := r.jobService.CreateJobTempDir(ctx, userName, req.ComputeType)
	if err != nil {
		logger.Errorf("[%v] create job temp directory err: %v", userName, err)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrJobFailCreateTempDir)
		return
	}

	resp := &dto.CreateTempDirResponse{
		Path: tempDir,
	}

	ginutil.Success(ctx, resp)
}

// GetWorkSpace
//
//	@Summary		获取工作空间
//	@Description	获取工作空间接口
//	@Tags			作业
//	@Accept			json
//	@Produce		json
//	@Response		200	{object}	string
//	@Router			/job/workspace [get]
func (r *apiRoute) GetWorkSpace(ctx *gin.Context) {
	workSpace := r.jobService.GetWorkSpace(ctx)

	ginutil.Success(ctx, &dto.GetWorkSpaceResponse{Name: workSpace})
}
