package job

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/admin/jobcpuusage"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/admin/jobcreate"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/admin/jobdelete"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/admin/jobget"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/admin/joblist"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/admin/joblistfiltered"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/admin/jobmonitorchart"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/admin/jobresidual"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/admin/jobsnapshotget"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/admin/jobsnapshotlist"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/admin/jobterminate"
	schema "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common/validation"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao/models"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/util"
)

func adminAllow(opUser, owner string) bool {
	// admin不验证
	return true
}

// AdminGet 管理员获取指定作业
func (h *Handler) AdminGet(c *gin.Context) {
	req := &jobget.Request{}
	if err := c.ShouldBindUri(req); err != nil {
		common.InvalidParams(c, "invalid params uri, "+err.Error())
		return
	}

	if JudgeGetError(c, ValidateJobID(c, req.JobID)) {
		return
	}

	ctx := logging.AppendWith(c, "func", "job.AdminGet", "jobID", req.JobID)
	jobInfo, err := h.jobSrv.Get(ctx, &req.Request, snowflake.Zero(), adminAllow, true)
	if JudgeGetError(c, err) {
		return
	}

	common.SuccessResp(c, util.ModelToAdminOpenAPIJob(jobInfo))
}

// AdminList 管理员获取所有作业
func (h *Handler) AdminList(c *gin.Context) {
	req := &joblist.Request{}
	if err := c.ShouldBindQuery(req); err != nil {
		common.InvalidParams(c, "invalid params query, "+err.Error())
		return
	}

	err := ValidateList(c, &req.Request)
	if JudgeListError(c, err) {
		return
	}

	userID, appID := snowflake.Zero(), snowflake.Zero()
	if len(req.UserID) != 0 {
		err = ValidateUserID(c, req.UserID)
		if JudgeListError(c, err) {
			return
		}
		userID = snowflake.MustParseString(req.UserID)
	}

	if len(req.AppID) != 0 {
		err = ValidateAppID(c, req.AppID)
		if JudgeListError(c, err) {
			return
		}
		appID = snowflake.MustParseString(req.AppID)
	}

	ctx := logging.AppendWith(c, "func", "job.AdminList")
	total, jobs, err := h.jobSrv.List(ctx, &req.Request, userID, appID, req.WithDelete, req.IsSystemFailed)
	if JudgeListError(c, err) {
		return
	}

	jobsResp := make([]*schema.AdminJobInfo, 0, len(jobs))
	for _, job := range jobs {
		openAPIJob := util.ModelToAdminOpenAPIJob(job)
		jobsResp = append(jobsResp, openAPIJob)
	}

	resp := &joblist.Data{
		Jobs:  jobsResp,
		Total: total,
	}

	common.SuccessResp(c, &resp)
}

// AdminJobListFiltered 管理员获取带有过滤条件的作业列表
func (h *Handler) AdminJobListFiltered(c *gin.Context) {
	logger := logrus.WithFields(logrus.Fields{
		"func": "AdminJobListFiltered",
	})
	logger.Info("Entering AdminJobListFiltered with query parameters: ", c.Request.URL.Query())
	req := &joblistfiltered.Request{}
	if err := c.ShouldBindQuery(req); err != nil {
		common.InvalidParams(c, "invalid params query, "+err.Error())
		return
	}
	logger.Infof("Request object: %+v", req)
	err := ValidateList(c, &req.Request.Request)
	if JudgeListError(c, err) {
		return
	}
	userID, appID := snowflake.Zero(), snowflake.Zero()
	if len(req.UserID) > 0 {
		err = ValidateUserID(c, req.UserID)
		if JudgeListError(c, err) {
			return
		}
		userID = snowflake.MustParseString(req.UserID)
	}
	if len(req.AppID) > 0 {
		err = ValidateAppID(c, req.AppID)
		if JudgeListError(c, err) {
			return
		}
		appID = snowflake.MustParseString(req.AppID)
	}
	if len(req.JobID) > 0 {
		err = ValidateJobID(c, req.JobID)
		if JudgeListError(c, err) {
			return
		}
	}
	if len(req.AccountID) > 0 {
		err = ValidateAccountID(c, req.AccountID)
		if JudgeListError(c, err) {
			return
		}
	}
	ctx := logging.AppendWith(c, "func", "job.ListFiltered", "userID", userID)
	logger.Infof("Calling ListFiltered with parameters: req=%+v, userID=%v, appID=%v", req, userID, appID)
	total, jobs, err := h.jobSrv.ListFiltered(ctx, req, userID, appID)
	if JudgeListError(c, err) {
		return
	}
	logger.Infof("ListFiltered successful: total=%v, jobs=%v", total, jobs)
	jobsResp := make([]*schema.AdminJobInfo, 0, len(jobs))
	for _, job := range jobs {
		openAPIJob := util.ModelToAdminOpenAPIJob(job)
		jobsResp = append(jobsResp, openAPIJob)
	}
	resp := &joblistfiltered.Data{
		Jobs:  jobsResp,
		Total: total,
	}
	common.SuccessResp(c, &resp)
}

// AdminCreate 管理员创建作业
func (h *Handler) AdminCreate(c *gin.Context) {
	logger := logging.GetLogger(c)
	req := &jobcreate.Request{}
	if err := c.ShouldBindJSON(req); err != nil {
		common.InvalidParams(c, "invalid params, "+err.Error())
		return
	}
	userID, err := ValidateUserInfo(c)
	if JudgeSubmitError(c, err) {
		return
	}
	// 预调度校验
	scheduleInfo, preSchedule, err := h.ValidatePreScheduleParams(c, req.PreScheduleID)
	if JudgeSubmitError(c, err) {
		return
	}
	if !preSchedule { // 非预调度时这些字段必填
		err := ValidateNoPreScheduleParams(c, &req.Request)
		if JudgeSubmitError(c, err) {
			return
		}
	} else {
		err := util.FillScheduleInfo(c, &req.Request, scheduleInfo)
		if JudgeSubmitError(c, err) {
			return
		}
	}
	// 验证参数
	appInfo, inputZone, outputZone, chargeParams, err := h.ValidateCreate(c, &req.Request, userID,
		func(ctx context.Context, appID snowflake.ID) bool {
			// admin创建作业时，暂不验证 AppQuota
			return true
		},
		func(ctx context.Context, userID snowflake.ID) (bool, error) {
			// admin创建作业时，暂不验证 Shared
			return true, nil
		},
		func(ctx context.Context, userID snowflake.ID) (bool, error) {
			// admin创建作业时，暂不验证 AllocType
			return true, nil
		})
	if JudgeSubmitError(c, err) {
		return
	}

	ctx := logging.AppendWith(c, "func", "job.AdminCreate", "userID", userID)
	jobID, err := h.jobSrv.Create(ctx, &req.Request, userID, chargeParams, appInfo, scheduleInfo, func(ctx context.Context, jobID snowflake.ID) (*models.Job, error) {
		return util.ConvertAdminJobModel(ctx, logger, req, userID, jobID, appInfo, inputZone, outputZone, req.Queue, scheduleInfo)
	})
	if JudgeSubmitError(c, err) {
		return
	}

	common.SuccessResp(c, jobcreate.Data{
		JobID: jobID,
	})
}

// AdminDelete 管理员删除作业
func (h *Handler) AdminDelete(c *gin.Context) {
	req := &jobdelete.Request{}
	if err := c.ShouldBindUri(req); err != nil {
		common.InvalidParams(c, "invalid params uri, "+err.Error())
		return
	}

	if JudgeDeleteError(c, ValidateJobID(c, req.JobID)) {
		return
	}

	ctx := logging.AppendWith(c, "func", "job.AdminDelete", "jobID", req.JobID)
	err := h.jobSrv.Delete(ctx, &req.Request, snowflake.Zero(), adminAllow)
	if JudgeDeleteError(c, err) {
		return
	}

	common.SuccessResp(c, nil)
}

// AdminTerminate 管理员终止指定作业
func (h *Handler) AdminTerminate(c *gin.Context) {
	req := &jobterminate.Request{}
	if err := c.ShouldBindUri(req); err != nil {
		common.InvalidParams(c, "invalid params uri, "+err.Error())
		return
	}

	if JudgeTerminateError(c, ValidateJobID(c, req.JobID)) {
		return
	}

	ctx := logging.AppendWith(c, "func", "job.AdminTerminate", "jobID", req.JobID)
	err := h.jobSrv.Terminate(ctx, &req.Request, snowflake.Zero(), adminAllow)
	if JudgeTerminateError(c, err) {
		return
	}

	common.SuccessResp(c, nil)
}

// AdminResidual 管理员获取指定作业的残差图
func (h *Handler) AdminResidual(c *gin.Context) {
	logger := logging.GetLogger(c)
	req := &jobresidual.Request{}
	if err := c.ShouldBindUri(req); err != nil {
		logger.Warnf("invalid params, %s", err.Error())
		handleValidationError(c, err, req, func(fe validation.Error) bool {
			return handleJobIDRequired(c, fe)
		})
		return
	}

	if JudgeGetResidualError(c, ValidateJobID(c, req.JobID)) {
		return
	}

	ctx := logging.AppendWith(c, "func", "job.AdminResidual", "jobID", req.JobID)
	residual, err := h.jobSrv.GetResidual(ctx, &req.Request, snowflake.Zero(), adminAllow)
	if JudgeGetResidualError(c, err) {
		return
	}

	common.SuccessResp(c, residual)
}

// AdminMonitorChart 管理员获取监控图表
func (h *Handler) AdminMonitorChart(c *gin.Context) {
	logger := logging.GetLogger(c)
	req := &jobmonitorchart.Request{}
	if err := c.ShouldBindUri(req); err != nil {
		logger.Warnf("invalid params, %s", err.Error())
		handleValidationError(c, err, req, func(fe validation.Error) bool {
			return handleJobIDRequired(c, fe)
		})
		return
	}

	if JudgeGetMonitorChartError(c, ValidateJobID(c, req.JobID)) {
		return
	}

	ctx := logging.AppendWith(c, "func", "job.AdminMonitorChart", "jobID", req.JobID)
	monitorChart, err := h.jobSrv.GetMonitorChart(ctx, &req.Request, snowflake.Zero(), adminAllow)
	if JudgeGetMonitorChartError(c, err) {
		return
	}

	resp, err := util.ModelToOpenAPIJobMonitorCharts(monitorChart)
	if err != nil {
		logger.Warnf("util.ModelToOpenAPIJobMonitorChart error: %v", err)
		common.InternalServerError(c, "util.ModelToOpenAPIJobMonitorChart error!")
		return
	}

	common.SuccessResp(c, resp)
}

// AdminSnapshots 获取云图集
func (h *Handler) AdminSnapshots(c *gin.Context) {
	logger := logging.GetLogger(c)
	req := &jobsnapshotlist.Request{}
	if err := c.ShouldBindUri(req); err != nil {
		logger.Warnf("invalid params, %s", err.Error())
		handleValidationError(c, err, req, func(fe validation.Error) bool {
			return handleJobIDRequired(c, fe)
		})
		return
	}

	if JudgeListSnapshotError(c, ValidateJobID(c, req.JobID)) {
		return
	}

	ctx := logging.AppendWith(c, "func", "job.AdminSnapshots", "jobID", req.JobID)
	snapshots, err := h.jobSrv.ListJobSnapshot(ctx, h.appSrv.Apps(), &req.Request, snowflake.Zero(), adminAllow)
	if JudgeListSnapshotError(c, err) {
		return
	}

	common.SuccessResp(c, snapshots)
}

// AdminSnapshotImg 获取云图数据
func (h *Handler) AdminSnapshotImg(c *gin.Context) {
	JobID := c.Param("JobID")
	Path := c.Query("Path")

	r, err := ValidateSnapshotImg(c, JobID, Path)
	if JudgeGetSnapshotError(c, err) {
		return
	}

	req := &jobsnapshotget.Request{
		Request: *r,
	}

	if JudgeGetSnapshotError(c, ValidateJobID(c, req.JobID)) {
		return
	}

	ctx := logging.AppendWith(c, "func", "job.AdminSnapshotImg", "jobID", req.JobID)
	img, err := h.jobSrv.GetJobSnapshot(ctx, h.appSrv.Apps(), &req.Request, snowflake.Zero(), adminAllow)
	if JudgeGetSnapshotError(c, err) {
		return
	}

	common.SuccessResp(c, img)
}

// AdminCpuUsage 获取作业CPU使用率
func (h *Handler) AdminCpuUsage(c *gin.Context) {
	req := &jobcpuusage.Request{}
	if err := c.ShouldBindUri(req); err != nil {
		common.InvalidParams(c, "invalid params uri, "+err.Error())
		return
	}

	if JudgeCpuUsageError(c, ValidateJobID(c, req.JobID)) {
		return
	}

	ctx := logging.AppendWith(c, "func", "job.AdminCpuUsage", "jobID", req.JobID)
	cpuUsage, err := h.jobSrv.GetCpuUsage(ctx, &req.Request, snowflake.Zero(), adminAllow)
	if JudgeCpuUsageError(c, err) {
		return
	}

	common.SuccessResp(c, cpuUsage)
}
