package job

import (
	"context"
	"net/http"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/common/openapi-go/utils/payby"
	iam_api "github.com/yuansuan/ticp/common/project-root-iam/iam-api"
	iam_client "github.com/yuansuan/ticp/common/project-root-iam/iam-client"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/config"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	api "github.com/yuansuan/ticp/common/project-root-api/common"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/jobbatchget"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/jobcpuusage"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/jobcreate"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/jobdelete"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/jobget"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/joblist"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/jobmonitorchart"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/jobpreschedule"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/jobresidual"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/jobresume"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/jobsnapshotlist"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/jobterminate"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/jobtransmitresume"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/jobtransmitsuspend"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/system/jobneedsyncfile"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/system/jobsyncfilestate"
	schema "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common/validation"
	app "github.com/yuansuan/ticp/iPaaS/project-root/internal/job/api/v1/application"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao/models"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/service/v1/application"
	jobservice "github.com/yuansuan/ticp/iPaaS/project-root/internal/job/service/v1/job"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/util"
)

// Handler ...
type Handler struct {
	appSrv      application.Service
	jobSrv      jobservice.Service
	userChecker util.UserChecker
}

// NewJobHandler ...
func NewJobHandler(control *app.Controller, jobSrv jobservice.Service, checker util.UserChecker) *Handler {
	appService := control.GetService()
	return &Handler{
		appSrv:      appService,
		jobSrv:      jobSrv,
		userChecker: checker,
	}
}

func selfAllow(opUser, owner string) bool {
	// 验证是否是自己的作业
	return opUser == owner
}

// InvalidJobID 无效的作业ID
func (h *Handler) InvalidJobID(c *gin.Context) {
	common.ErrorResp(c, http.StatusBadRequest, api.InvalidJobID, "invalid job id")
}

// Get 获取作业
func (h *Handler) Get(c *gin.Context) {
	req := &jobget.Request{}
	if err := c.ShouldBindUri(req); err != nil {
		common.InvalidParams(c, "invalid params uri, "+err.Error())
		return
	}
	userID, err := ValidateUserInfo(c)
	if JudgeGetError(c, err) {
		return
	}
	if JudgeGetError(c, ValidateJobID(c, req.JobID)) {
		return
	}

	ctx := logging.AppendWith(c, "func", "job.GetQuota", "userID", userID, "jobID", req.JobID)
	job, err := h.jobSrv.Get(ctx, req, userID, func(opUser string, owner string) bool {
		if selfAllow(opUser, owner) {
			return true
		}
		isYsProductUser, err := h.userChecker.IsYsProductUser(snowflake.MustParseString(opUser))
		if err != nil {
			logging.GetLogger(ctx).Infof("check user is ys product user error: %v", err)
			return false
		}
		return isYsProductUser
	}, false)
	if JudgeGetError(c, err) {
		return
	}

	schemaJob := util.ModelToOpenAPIJob(job)
	schemaJob, err = util.TrimAdminParams(schemaJob)
	if err != nil {
		logging.GetLogger(ctx).Warnf("TrimAdminParams err: %v", err)
		common.InternalServerError(c, "TrimAdminParams err!")
		return
	}
	common.SuccessResp(c, schemaJob)
}

// BatchGet 批量获取作业
func (h *Handler) BatchGet(c *gin.Context) {
	req := &jobbatchget.Request{}
	if err := c.ShouldBindJSON(req); err != nil {
		common.InvalidParams(c, "invalid params, "+err.Error())
		return
	}
	userID, err := ValidateUserInfo(c)
	if JudgeBatchGetError(c, err) {
		return
	}
	err = ValidateBatchGet(c, req)
	if JudgeBatchGetError(c, err) {
		return
	}

	queryIDs := []string{}
	for _, id := range req.JobIDs {
		if err := ValidateID(c, id); err != nil {
			continue
		}
		queryIDs = append(queryIDs, id)
	}

	ctx := logging.AppendWith(c, "func", "job.BatchGet", "userID", userID)
	jobs, err := h.jobSrv.BatchGet(ctx, queryIDs, userID)
	if JudgeBatchGetError(c, err) {
		return
	}

	schemaJobs := []*schema.JobInfo{}
	for _, job := range jobs {
		schemaJob := util.ModelToOpenAPIJob(job)
		schemaJob, err := util.TrimAdminParams(schemaJob)
		if err != nil {
			logging.GetLogger(ctx).Infof("TrimAdminParams err: %v", err)
			continue
		}
		schemaJobs = append(schemaJobs, schemaJob)
	}

	common.SuccessResp(c, schemaJobs)
}

// List 获取所有作业
func (h *Handler) List(c *gin.Context) {
	req := &joblist.Request{}
	if err := c.ShouldBindQuery(req); err != nil {
		common.InvalidParams(c, "invalid params query, "+err.Error())
		return
	}
	userID, err := ValidateUserInfo(c)
	if JudgeListError(c, err) {
		return
	}
	err = ValidateList(c, req)
	if JudgeListError(c, err) {
		return
	}
	ctx := logging.AppendWith(c, "func", "job.List", "userID", userID)
	total, jobs, err := h.jobSrv.List(ctx, req, userID, snowflake.Zero(), false, false)
	if JudgeListError(c, err) {
		return
	}

	jobsResp := make([]*schema.JobInfo, 0, len(jobs))
	for _, job := range jobs {
		openAPIJob := util.ModelToOpenAPIJob(job)
		openAPIJob, err = util.TrimAdminParams(openAPIJob)
		if err != nil {
			logging.GetLogger(ctx).Warnf("TrimAdminParams err: %v", err)
			continue
		}
		jobsResp = append(jobsResp, openAPIJob)
	}
	data := &joblist.Data{
		Total: total,
		Jobs:  jobsResp,
	}
	common.SuccessResp(c, &data)
}

// Create 创建作业
func (h *Handler) Create(c *gin.Context) {
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
		err := ValidateNoPreScheduleParams(c, req)
		if JudgeSubmitError(c, err) {
			return
		}
	} else { //预调度情况
		err := util.FillScheduleInfo(c, req, scheduleInfo)
		if JudgeSubmitError(c, err) {
			return
		}
	}

	billEnabled := config.GetConfig().BillEnabled
	// 验证代支付
	payByUserId := snowflake.ID(0)
	if billEnabled && req.PayBy != "" {
		reqPayBy, err := payby.ParseToken(req.PayBy)
		if err != nil {
			common.ErrorResp(c, http.StatusBadRequest, api.InvalidArgumentPayBy, "invalid payBy")
			return
		}
		iamClient := iam_client.NewClient(config.GetConfig().OpenAPIEndpoint,
			config.GetConfig().AK, config.GetConfig().AS)
		resp, err := iamClient.GetSecret(&iam_api.GetSecretRequest{
			AccessKeyId: reqPayBy.GetAccessKeyID(),
		})
		if err != nil {
			if strings.Contains(err.Error(), "secret not found") {
				logger.Infof("accessKeyId not found. accessKeyId: %s", reqPayBy.GetAccessKeyID())
				common.InvalidParams(c, "accessKeyId not found")
				return
			}
			logger.Warnf("get secret error, err: %v", err)
			common.InternalServerError(c, "internal server error")
			return
		}
		payByUserId, err = snowflake.ParseString(resp.YSId)
		if JudgeSubmitError(c, err) {
			return
		}
	}

	// 验证参数
	appInfo, inputZone, outputZone, chargeParams, err := h.ValidateCreate(c, req, userID,
		func(ctx context.Context, appID snowflake.ID) bool { // AppQuotaChecker
			al, err3 := h.appSrv.AppsAllow().GetAllow(ctx, appID)
			if err3 == nil {
				logger.Infof("create job, get allow success, allow: %v", al)
				return true
			}
			aq, err2 := h.appSrv.AppsQuota().GetQuota(ctx, appID, userID)
			if err2 == nil {
				logger.Infof("create job, get quota success, quota: %v", aq)
				return true
			}
			if billEnabled && req.PayBy != "" && errors.Is(err2, common.ErrAppQuotaNotFound) {
				_, err2 = h.appSrv.AppsQuota().GetQuota(ctx, appID, payByUserId)
				if err2 == nil {
					return true
				}
				if errors.Is(err2, common.ErrAppQuotaNotFound) {
					logger.Infof("app quota not found, appID: %d, err: %v", appID, err2)
					return false
				}
			}
			logger.Warnf("create job, get quota error, err: %v", err2)
			return false
		},
		func(ctx context.Context, userID snowflake.ID) (bool, error) { // SharedChecker
			return h.userChecker.IsYsProductUser(userID)
		},
		func(ctx context.Context, userID snowflake.ID) (bool, error) { // AllocTypeChecker
			return h.userChecker.IsYsProductUser(userID)
		})

	if JudgeSubmitError(c, err) {
		return
	}

	ctx := logging.AppendWith(c, "func", "job.Create", "userID", userID)
	jobID, err := h.jobSrv.Create(ctx, req, userID, chargeParams, appInfo, scheduleInfo,
		func(ctx context.Context, jobID snowflake.ID) (*models.Job, error) {
			return util.ConvertJobModel(ctx, logger, req, userID, jobID, appInfo,
				inputZone, outputZone, nil, scheduleInfo)
		})
	if JudgeSubmitError(c, err) {
		return
	}
	common.SuccessResp(c, jobcreate.Data{
		JobID: jobID,
	})
}

// Delete 删除作业
func (h *Handler) Delete(c *gin.Context) {
	req := &jobdelete.Request{}
	if err := c.ShouldBindUri(req); err != nil {
		common.InvalidParams(c, "invalid params uri, "+err.Error())
		return
	}
	userID, err := ValidateUserInfo(c)
	if JudgeDeleteError(c, err) {
		return
	}
	if JudgeDeleteError(c, ValidateJobID(c, req.JobID)) {
		return
	}

	ctx := logging.AppendWith(c, "func", "job.Delete", "userID", userID, "jobID", req.JobID)
	err = h.jobSrv.Delete(ctx, req, userID, selfAllow)
	if JudgeDeleteError(c, err) {
		return
	}
	common.SuccessResp(c, nil)
}

// Terminate 终止作业
func (h *Handler) Terminate(c *gin.Context) {
	req := &jobterminate.Request{}
	if err := c.ShouldBindUri(req); err != nil {
		common.InvalidParams(c, "invalid params uri, "+err.Error())
		return
	}
	userID, err := ValidateUserInfo(c)
	if JudgeTerminateError(c, err) {
		return
	}
	if JudgeTerminateError(c, ValidateJobID(c, req.JobID)) {
		return
	}
	ctx := logging.AppendWith(c, "func", "job.Terminate", "userID", userID, "jobID", req.JobID)
	err = h.jobSrv.Terminate(ctx, req, userID, selfAllow)
	if JudgeTerminateError(c, err) {
		return
	}
	common.SuccessResp(c, nil)
}

// Resume 恢复作业
func (h *Handler) Resume(c *gin.Context) {
	req := &jobresume.Request{}
	if err := c.ShouldBindUri(req); err != nil {
		common.InvalidParams(c, "invalid params uri, "+err.Error())
		return
	}
	userID, err := ValidateUserInfo(c)
	if JudgeResumeError(c, err) {
		return
	}
	if JudgeResumeError(c, ValidateJobID(c, req.JobID)) {
		return
	}
	ctx := logging.AppendWith(c, "func", "job.Resume", "userID", userID, "jobID", req.JobID)
	err = h.jobSrv.Resume(ctx, req, userID, selfAllow)
	if JudgeResumeError(c, err) {
		return
	}
	common.SuccessResp(c, nil)
}

// TransmitSuspend 传输暂停
func (h *Handler) TransmitSuspend(c *gin.Context) {
	req := &jobtransmitsuspend.Request{}
	if err := c.ShouldBindUri(req); err != nil {
		common.InvalidParams(c, "invalid params uri, "+err.Error())
		return
	}
	userID, err := ValidateUserInfo(c)
	if JudgeTransmitSuspendError(c, err) {
		return
	}
	if JudgeTransmitSuspendError(c, ValidateJobID(c, req.JobID)) {
		return
	}
	ctx := logging.AppendWith(c, "func", "job.TransmitSuspend", "userID", userID, "jobID", req.JobID)
	err = h.jobSrv.TransmitSuspend(ctx, req, userID, selfAllow)
	if JudgeTransmitSuspendError(c, err) {
		return
	}
	common.SuccessResp(c, nil)
}

// TransmitResume 传输恢复
func (h *Handler) TransmitResume(c *gin.Context) {
	req := &jobtransmitresume.Request{}
	if err := c.ShouldBindUri(req); err != nil {
		common.InvalidParams(c, "invalid params uri, "+err.Error())
		return
	}
	userID, err := ValidateUserInfo(c)
	if JudgeTransmitResumeError(c, err) {
		return
	}
	if JudgeTransmitResumeError(c, ValidateJobID(c, req.JobID)) {
		return
	}
	ctx := logging.AppendWith(c, "func", "job.TransmitResume", "userID", userID, "jobID", req.JobID)
	err = h.jobSrv.TransmitResume(ctx, req, userID, selfAllow)
	if JudgeTransmitResumeError(c, err) {
		return
	}
	common.SuccessResp(c, nil)
}

// Residual 获取残差图
func (h *Handler) Residual(c *gin.Context) {
	logger := logging.GetLogger(c)
	req := &jobresidual.Request{}
	if err := c.ShouldBindUri(req); err != nil {
		logger.Warnf("invalid params, %s", err.Error())
		handleValidationError(c, err, req, func(fe validation.Error) bool {
			return handleJobIDRequired(c, fe)
		})
		return
	}
	userID, err := ValidateUserInfo(c)
	if JudgeGetResidualError(c, err) {
		return
	}
	if JudgeGetResidualError(c, ValidateJobID(c, req.JobID)) {
		return
	}
	ctx := logging.AppendWith(c, "func", "job.Residual", "userID", userID, "jobID", req.JobID)
	residual, err := h.jobSrv.GetResidual(ctx, req, userID, selfAllow)
	if JudgeGetResidualError(c, err) {
		return
	}
	common.SuccessResp(c, residual)
}

// MonitorChart 监控图表
func (h *Handler) MonitorChart(c *gin.Context) {
	logger := logging.GetLogger(c)
	req := &jobmonitorchart.Request{}
	if err := c.ShouldBindUri(req); err != nil {
		logger.Warnf("invalid params, %s", err.Error())
		handleValidationError(c, err, req, func(fe validation.Error) bool {
			return handleJobIDRequired(c, fe)
		})
		return
	}
	userID, err := ValidateUserInfo(c)
	if JudgeGetMonitorChartError(c, err) {
		return
	}
	if JudgeGetMonitorChartError(c, ValidateJobID(c, req.JobID)) {
		return
	}
	ctx := logging.AppendWith(c, "func", "job.MonitorChart", "userID", userID, "jobID", req.JobID)
	monitorChart, err := h.jobSrv.GetMonitorChart(ctx, req, userID, selfAllow)
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

// PreSchedule 创建作业预调度
func (h *Handler) PreSchedule(c *gin.Context) {
	logger := logging.GetLogger(c)
	req := &jobpreschedule.Request{}
	if err := c.ShouldBindJSON(req); err != nil {
		common.InvalidParams(c, "invalid params, "+err.Error())
		return
	}
	logger.Infof("preschedule job,request:%s", spew.Sdump(req))
	userID, err := ValidateUserInfo(c)
	if JudgeSubmitError(c, err) {
		return
	}

	// 验证参数
	appInfo, zones, err := h.ValidatePreSchedule(c, req, userID,
		func(ctx context.Context, appID snowflake.ID) bool { // AppQuotaChecker
			al, err3 := h.appSrv.AppsAllow().GetAllow(ctx, appID)
			if err3 == nil {
				logger.Infof("pre schedule, get allow success, allow: %v", al)
				return true
			}
			aq, err2 := h.appSrv.AppsQuota().GetQuota(ctx, appID, userID)
			if err2 == nil {
				logger.Infof("pre schedule, get quota success, quota: %v", aq)
				return true
			}

			logger.Warnf("pre schedule, get quota error, err: %v", err2)
			return false
		},
		func(ctx context.Context, userID snowflake.ID) (bool, error) { // SharedChecker
			return h.userChecker.IsYsProductUser(userID)
		},
	)
	if JudgePreScheduleError(c, err) {
		return
	}
	resp, err := h.jobSrv.PreSchedule(c, req, zones, userID, appInfo)
	if JudgePreScheduleError(c, err) {
		return
	}
	common.SuccessResp(c, resp)
}

// ListNeedSyncFileJobs 同步文件作业
func (h *Handler) ListNeedSyncFileJobs(c *gin.Context) {
	req := &jobneedsyncfile.Request{}
	if err := c.ShouldBindQuery(req); err != nil {
		common.InvalidParams(c, "invalid params query, "+err.Error())
		return
	}
	err := ValidListNeedSyncFileJobs(c, req)
	if JudgeListSyncNeedFileError(c, err) {
		return
	}
	logging.GetLogger(c).Infof("list need sync file jobs, request:%+v", req)
	resp, err := h.jobSrv.ListNeedSyncFileJobs(c, req)
	if JudgeListSyncNeedFileError(c, err) {
		return
	}
	common.SuccessResp(c, &resp)
}

// SyncFileState 更新作业同步状态信息
func (h *Handler) SyncFileState(c *gin.Context) {
	req := &jobsyncfilestate.Request{}
	jobIDStr := c.Param("JobID")
	if err := c.ShouldBindJSON(req); err != nil {
		common.InvalidParams(c, "invalid params query, "+err.Error())
		return
	}
	err := ValidSyncFileState(c, req)
	if JudgeSyncFileStateError(c, err) {
		return
	}
	if JudgeSyncFileStateError(c, ValidateJobID(c, jobIDStr)) {
		return
	}
	logging.GetLogger(c).Infof("sync file state,request:%+v", req)
	err = h.jobSrv.UpdateSyncFileState(c, req, jobIDStr)
	if JudgeSyncFileStateError(c, err) {
		return
	}
	common.SuccessResp(c, nil)
}

// Snapshots 获取云图集
func (h *Handler) Snapshots(c *gin.Context) {
	logger := logging.GetLogger(c)
	req := &jobsnapshotlist.Request{}
	if err := c.ShouldBindUri(req); err != nil {
		logger.Warnf("invalid params, %s", err.Error())
		handleValidationError(c, err, req, func(fe validation.Error) bool {
			return handleJobIDRequired(c, fe)
		})
		return
	}
	userID, err := ValidateUserInfo(c)
	if JudgeListSnapshotError(c, err) {
		return
	}
	if JudgeListSnapshotError(c, ValidateJobID(c, req.JobID)) {
		return
	}
	ctx := logging.AppendWith(c, "func", "job.Snapshots", "userID", userID, "jobID", req.JobID)
	snapshots, err := h.jobSrv.ListJobSnapshot(ctx, h.appSrv.Apps(), req, userID, selfAllow)
	if JudgeListSnapshotError(c, err) {
		return
	}
	common.SuccessResp(c, snapshots)
}

// SnapshotImg 获取云图数据
func (h *Handler) SnapshotImg(c *gin.Context) {
	logger := logging.GetLogger(c)
	JobID := c.Param("JobID")
	Path := c.Query("Path")
	req, err := ValidateSnapshotImg(c, JobID, Path)
	if JudgeGetSnapshotError(c, err) {
		return
	}
	logger.Infof("get job snapshot img,request:%+v", req)
	userID, err := ValidateUserInfo(c)
	if JudgeGetSnapshotError(c, err) {
		return
	}
	if JudgeGetSnapshotError(c, ValidateJobID(c, req.JobID)) {
		return
	}
	ctx := logging.AppendWith(c, "func", "job.SnapshotImg", "userID", userID, "jobID", req.JobID)
	img, err := h.jobSrv.GetJobSnapshot(ctx, h.appSrv.Apps(), req, userID, selfAllow)
	if JudgeGetSnapshotError(c, err) {
		return
	}
	common.SuccessResp(c, img)
}

func (h *Handler) CpuUsage(c *gin.Context) {
	req := &jobcpuusage.Request{}
	if err := c.ShouldBindUri(req); err != nil {
		common.InvalidParams(c, "invalid params uri, "+err.Error())
		return
	}
	userID, err := ValidateUserInfo(c)
	if JudgeCpuUsageError(c, err) {
		return
	}
	if JudgeCpuUsageError(c, ValidateJobID(c, req.JobID)) {
		return
	}
	ctx := logging.AppendWith(c, "func", "job.CpuUsage", "userID", userID, "jobID", req.JobID)
	cpuUsage, err := h.jobSrv.GetCpuUsage(ctx, req, userID, selfAllow)
	if JudgeCpuUsageError(c, err) {
		return
	}
	common.SuccessResp(c, cpuUsage)
}
