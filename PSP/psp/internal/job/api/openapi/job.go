package openapi

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/yuansuan/ticp/PSP/psp/internal/common"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/oplog"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/approve"
	"github.com/yuansuan/ticp/PSP/psp/internal/job/dto"
	"github.com/yuansuan/ticp/PSP/psp/internal/job/dto/openapi"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/ginutil"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/structutil"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/strutil"
)

// CreateTempDir
//
//	@Summary		创建作业临时目录
//	@Description	创建作业临时目录接口
//	@Tags			作业
//	@Accept			json
//	@Produce		json
//	@Param			param	query		openapi.CreateTempDirRequest	true	"请求参数"
//	@Response		200		{object}	openapi.CreateTempDirResponse
//	@Router			/openapi/job/createTempDir [post]
func (r *openapiApiRoute) CreateTempDir(ctx *gin.Context) {

	logger := logging.GetLogger(ctx)

	var req = openapi.CreateTempDirRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	if err := r.validate.Struct(req); err != nil {
		ginutil.Error(ctx, errcode.ErrInvalidParam, err.Error())
		return
	}

	userName := ginutil.GetUserName(ctx)
	tempDir, err := r.jobService.CreateJobTempDir(ctx, userName, req.ComputeType)
	if err != nil {
		logger.Errorf("[%v] create job temp directory err: %v", userName, err)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrJobFailCreateTempDir)
		return
	}

	rsp := &openapi.CreateTempDirResponse{}
	if err = structutil.CopyStruct(rsp, &dto.CreateTempDirResponse{
		Path: tempDir,
	}); err != nil {
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrJobFailCreateTempDir)
		return
	}

	ginutil.Success(ctx, rsp)
}

// Submit
//
//	@Summary		作业提交
//	@Description	作业提交接口
//	@Tags			作业
//	@Accept			json
//	@Produce		json
//	@Param			param	query		openapi.JobSubmitRequest	true	"请求参数"
//	@Response		200		{object}	openapi.JobSubmitResponse
//	@Router			/openapi/job/submit [post]
func (r *openapiApiRoute) Submit(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	var req = openapi.JobSubmitRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	if err := r.validate.Struct(req); err != nil {
		ginutil.Error(ctx, errcode.ErrInvalidParam, err.Error())
		return
	}

	if strutil.IsEmpty(req.ProjectID) {
		req.ProjectID = common.PersonalProjectID.String()
	}

	userID := ginutil.GetUserID(ctx)
	userName := ginutil.GetUserName(ctx)

	fields := make([]*dto.Field, 0)
	for _, field := range req.Fields {
		fields = append(fields, &dto.Field{
			ID:     field.ID,
			Type:   field.Type,
			Value:  field.Value,
			Values: field.Values,
		})
	}

	submitParam := &dto.SubmitParam{
		AppID:     req.AppID,
		ProjectID: req.ProjectID,
		UserID:    snowflake.ID(userID),
		UserName:  userName,
		MainFiles: req.MainFiles,
		WorkDir: &dto.WorkDir{
			Path:   req.WorkDir.Path,
			IsTemp: true,
		},
		Fields:    fields,
		IsOpenApi: true,
	}
	jobIDs, err := r.jobService.JobSubmit(ctx, submitParam)
	if err != nil {
		logger.Errorf("submit job err: %v", err)
		if strings.Contains(err.Error(), "InvalidAccountStatus.NotEnoughBalance") {
			errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrJobFailNotEnoughBalance)
		} else {
			errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrJobFailSubmit)
		}

		return
	}

	ginutil.Success(ctx, &openapi.JobSubmitResponse{JobIDs: jobIDs})
}

// JobDetail
//
//	@Summary		作业详情
//	@Description	作业详情接口
//	@Tags			作业
//	@Accept			json
//	@Produce		json
//	@Param			param	query		openapi.JobDetailRequest	true	"请求参数"
//	@Response		200		{object}	openapi.JobDetailInfo
//	@Router			/openapi/job/detail [get]
func (r *openapiApiRoute) JobDetail(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	var req = openapi.JobDetailRequest{}
	if err := ctx.BindQuery(&req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	if err := r.validate.Struct(req); err != nil {
		ginutil.Error(ctx, errcode.ErrInvalidParam, err.Error())
		return
	}

	detail, err := r.jobService.GetJobDetail(ctx, req.JobID)
	if err != nil {
		logger.Errorf("get job [%v] detail err: %v", req.JobID, err)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrJobFailGet)
		return
	}

	rsp := &openapi.JobDetailInfo{}
	if err = structutil.CopyStruct(rsp, detail); err != nil {
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrJobFailCreateTempDir)
		return
	}

	ginutil.Success(ctx, rsp)
}

// JobTerminate
//
//	@Summary		作业终止
//	@Description	作业终止接口
//	@Tags			作业
//	@Accept			json
//	@Produce		json
//	@Param			param	query		openapi.JobTerminateRequest	true	"请求参数"
//	@Response		200		{object}	nil
//	@Router			/openapi/job/terminate [post]
func (r *openapiApiRoute) JobTerminate(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	var req = openapi.JobTerminateRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	if err := r.validate.Struct(req); err != nil {
		ginutil.Error(ctx, errcode.ErrInvalidParam, err.Error())
		return
	}

	outJobID, err := r.jobService.GetOutIDByJobID(ctx, snowflake.MustParseString(req.JobID))
	if err != nil {
		logger.Errorf("get job [%v] err: %v", req.JobID, err)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrJobFailGet)
		return
	}

	if err := r.jobService.JobTerminate(ctx, outJobID, req.ComputeType); err != nil {
		logger.Errorf("terminate job [%v] err: %v", req.JobID, err)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrJobFailTerminate)
		return
	}

	job, _ := r.jobService.GetJobDetailByOutID(ctx, outJobID, req.ComputeType)
	var jobName string
	if job != nil {
		jobName = job.Name
	}
	oplog.GetInstance().SaveAuditLogInfo(ctx, approve.OperateTypeEnum_JOB_MANAGER, fmt.Sprintf("【OPENAPI】用户%v终止求解作业[%v]", ginutil.GetUserName(ctx), jobName))

	ginutil.Success(ctx, nil)
}
