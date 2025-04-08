package api

import (
	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/internal/job/dto"
	"github.com/yuansuan/ticp/PSP/psp/internal/job/util"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/ginutil"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/strutil"
	"github.com/yuansuan/ticp/PSP/psp/pkg/xtype"
)

// JobList
//
//	@Summary		作业列表
//	@Description	作业列表接口
//	@Tags			作业
//	@Accept			json
//	@Produce		json
//	@Param			param	body		dto.JobListRequest	true	"请求参数"
//	@Response		200		{object}	dto.JobListResponse
//	@Router			/job/list [post]
func (r *apiRoute) JobList(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	var req = dto.JobListRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	loginUserID := ginutil.GetUserID(ctx)

	jobs, total, err := r.jobService.GetJobList(ctx, req.Filter, req.Page, req.OrderSort, snowflake.ID(loginUserID))
	if err != nil {
		logger.Errorf("get job list err: %v", err)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrJobFailList)
		return
	}

	jobList := make([]*dto.JobListInfo, 0, len(jobs))
	for _, job := range jobs {
		jobInfo := util.ConvertJob2ListInfo(job)
		if jobInfo != nil {
			jobList = append(jobList, jobInfo)
		}
	}

	resp := &dto.JobListResponse{
		Page: &xtype.PageResp{
			Index: req.Page.Index,
			Size:  req.Page.Size,
			Total: total,
		},
		Jobs: jobList,
	}

	ginutil.Success(ctx, resp)
}

// JobComputeTypeList
//
//	@Summary		获取计算类型列表
//	@Description	获取计算类型列表接口
//	@Tags			作业
//	@Accept			json
//	@Produce		json
//	@Param			param	query		dto.JobComputeTypeListRequest	true	"请求参数"
//	@Success		200		{object}	dto.JobComputeTypeListResponse
//	@Router			/job/computeTypes [get]
func (r *apiRoute) JobComputeTypeList(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	loginUserID := ginutil.GetUserID(ctx)

	computeTypeNameList, err := r.jobService.GetJobComputeTypeList(ctx, snowflake.ID(loginUserID))
	if err != nil {
		logger.Errorf("get job compute type list err: %v", err)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrJobFailComputeTypeList)
		return
	}

	ginutil.Success(ctx, &dto.JobComputeTypeListResponse{ComputeTypes: computeTypeNameList})
}

// JobSetNameList
//
//	@Summary		获取作业集名称列表
//	@Description	获取作业集名称列表接口
//	@Tags			作业
//	@Accept			json
//	@Produce		json
//	@Param			param	query		dto.JobSetNameListRequest	true	"请求参数"
//	@Success		200		{object}	dto.JobSetNameListResponse
//	@Router			/job/jobSetNames [get]
func (r *apiRoute) JobSetNameList(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	req := &dto.JobSetNameListRequest{}
	if err := ctx.BindQuery(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	loginUserID := ginutil.GetUserID(ctx)

	jobSetNames, err := r.jobService.GetJobSetNameList(ctx, req.ComputeType, snowflake.ID(loginUserID))
	if err != nil {
		logger.Errorf("get job set name list err: %v", err)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrJobFailJobSetNameList)
		return
	}

	ginutil.Success(ctx, &dto.JobSetNameListResponse{JobSetNames: jobSetNames})
}

// JobAppNameList
//
//	@Summary		获取作业应用名称列表
//	@Description	获取作业应用名称列表接口
//	@Tags			作业
//	@Accept			json
//	@Produce		json
//	@Param			param	query		dto.JobAppNameListRequest	true	"请求参数"
//	@Success		200		{object}	[]string
//	@Router			/job/appNames [get]
func (r *apiRoute) JobAppNameList(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	req := &dto.JobAppNameListRequest{}
	if err := ctx.BindQuery(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	loginUserID := ginutil.GetUserID(ctx)

	appNames, err := r.jobService.GetJobAppNameList(ctx, req.ComputeType, snowflake.ID(loginUserID))
	if err != nil {
		logger.Errorf("get job app name list err: %v", err)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrJobFailAppNameList)
		return
	}

	ginutil.Success(ctx, appNames)
}

// JobUserNameList
//
//	@Summary		获取作业用户名称列表
//	@Description	获取作业用户名称列表接口
//	@Tags			作业
//	@Accept			json
//	@Produce		json
//	@Param			param	query		dto.JobUserNameListRequest	true	"请求参数"
//	@Success		200		{object}	[]string
//	@Router			/job/userNames [get]
func (r *apiRoute) JobUserNameList(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	req := &dto.JobUserNameListRequest{}
	if err := ctx.BindQuery(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	loginUserID := ginutil.GetUserID(ctx)

	userNames, err := r.jobService.GetJobUserNameList(ctx, req.ComputeType, snowflake.ID(loginUserID))
	if err != nil {
		logger.Errorf("get job user name list err: %v", err)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrJobFailUserNameList)
		return
	}

	ginutil.Success(ctx, userNames)
}

// JobQueueNameList
//
//	@Summary		获取作业队列名称列表
//	@Description	获取作业队列名称列表接口
//	@Tags			作业
//	@Accept			json
//	@Produce		json
//	@Param			param	query		dto.JobQueueNameListRequest	true	"请求参数"
//	@Response		200		{object}	[]string
//	@Router			/job/queueNames [get]
func (r *apiRoute) JobQueueNameList(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	req := &dto.JobQueueNameListRequest{}
	if err := ctx.BindQuery(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	loginUserID := ginutil.GetUserID(ctx)

	queueNames, err := r.jobService.GetJobQueueNameList(ctx, req.ComputeType, snowflake.ID(loginUserID))
	if err != nil {
		logger.Errorf("get job queue name list err: %v", err)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrJobFailQueueNameList)
		return
	}

	names := make([]string, 0, len(queueNames))
	for _, name := range queueNames {
		if !strutil.IsEmpty(name) {
			names = append(names, name)
		}
	}

	ginutil.Success(ctx, names)
}
