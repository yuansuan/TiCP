package api

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/common/go-kit/logging/trace"
	"github.com/yuansuan/ticp/common/project-root-api/hpc/jobstate"
	"github.com/yuansuan/ticp/common/project-root-api/hpc/v1/job"
	v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/config"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/dao"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/dao/models"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/log"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/response"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/util"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/pkg/errorcode"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/pkg/with"
)

const (
	IdempotentIDKey = "IdempotentID"
)

func PostJobs(c *gin.Context) {
	var err error
	logger := trace.GetLogger(c)

	s, err := getState(c)
	if err = response.InternalErrorIfError(c, err, errorcode.InternalServerError); err != nil {
		err = fmt.Errorf("get state from gin ctx failed, %w", err)
		logger.Error(err)
		return
	}

	ctx := c.Request.Context()
	idempotentId := c.Query(IdempotentIDKey)
	if idempotentId != "" {
		exist, j, err := dao.Default.GetJobByIdempotentId(ctx, idempotentId)
		if err = response.InternalErrorIfError(c, err, errorcode.DatabaseInternalServerError); err != nil {
			err = fmt.Errorf("get job by idempotent id failed, %w", err)
			logger.Error(err)
			return
		}
		if exist {
			response.OK(c, j.ToHTTPModel())
			return
		}
	}

	req := new(job.SystemPostRequest)
	err = bindPostJobsReq(req, c)
	if err = response.BadRequestIfError(c, err, errorcode.InvalidArgument); err != nil {
		err = fmt.Errorf("bind post jobs request failed, %w", err)
		logger.Error(err)
		return
	}
	logger.Infof("post jobs request: %#v", req)

	if req.CustomStateRule != nil && req.CustomStateRule.KeyStatement != "" {
		if req.CustomStateRule.ResultState != jobstate.Completed.String() &&
			req.CustomStateRule.ResultState != jobstate.Failed.String() {
			err = fmt.Errorf("CustomStateRule.ResultState should only be in [ completed | failed ]")
			_ = response.BadRequestIfError(c, err, errorcode.InvalidCustomStateRule)
			logger.Error(err)
			return
		}
	}

	if req.Resource.CoresPerNode != nil && *req.Resource.CoresPerNode <= 0 {
		err = fmt.Errorf("Resourcece.CorePerNode cannot less or equal than 0")
		_ = response.BadRequestIfError(c, err, errorcode.InvalidArgumentCoresPerNode)
		logger.Error(err)
		return
	}

	j, err := newJobModel(req, s.Conf, idempotentId) //获取 j，即*models.Job
	if err = response.InternalErrorIfError(c, err, errorcode.InternalServerError); err != nil {
		err = fmt.Errorf("new job model failed, %w", err)
		logger.Error(err)
		return
	}

	jobID, err := dao.Default.InsertJobWithGenerateID(ctx, j) // j 落库；包括 RequestCores:int64(req.Resource.Cores)
	// 之后向slurm提交时再用jobID查 信息
	if err = response.InternalErrorIfError(c, err, errorcode.DatabaseInternalServerError); err != nil {
		err = fmt.Errorf("insert job to db failed, %w", err)
		logger.Error(err)
		return
	}
	s.Factory.ProduceJob(c, jobID)

	response.OK(c, j.ToHTTPModel())
}

func bindPostJobsReq(req *job.SystemPostRequest, c *gin.Context) error {
	if err := c.ShouldBind(req); err != nil {
		return fmt.Errorf("bind request body failed, %w", err)
	}

	return nil
}

func newJobModel(req *job.SystemPostRequest, cfg *config.Config, idempotentId string) (*models.Job, error) {
	var err error
	for i := range req.Inputs {
		if req.Inputs[i].Type == v20230530.HPCStorageType {
			req.Inputs[i].Src, err = util.ReplaceEndpoint(req.Inputs[i].Src, cfg.HpcStorageAddress)
			if err != nil {
				return nil, err
			}
		}
	}

	if req.Output != nil && req.Output.Type == v20230530.HPCStorageType {
		req.Output.Dst, err = util.ReplaceEndpoint(req.Output.Dst, cfg.HpcStorageAddress)
		if err != nil {
			return nil, err
		}
	}

	input, err := jsoniter.MarshalToString(req.Inputs)
	if err != nil {
		return nil, err
	}

	output, err := jsoniter.MarshalToString(req.Output)
	if err != nil {
		return nil, err
	}

	customStateRule, err := jsoniter.MarshalToString(req.CustomStateRule)
	if err != nil {
		return nil, err
	}

	envs, err := parseEnvs(req.Environment)
	if err != nil {
		return nil, err
	}

	schedulerSubmitFlags, err := jsoniter.MarshalToString(req.JobSchedulerSubmitFlags)
	if err != nil {
		return nil, err
	}

	j := &models.Job{
		IdempotentId:         idempotentId,
		State:                jobstate.Preparing,
		Inputs:               input,
		Output:               output,
		CustomStateRule:      customStateRule,
		EnvVars:              envs,
		Command:              req.Command,
		RequestCores:         int64(req.Resource.Cores),
		AllocType:            req.Resource.AllocType,
		IsOverride:           req.Override.Enable,
		WorkDir:              req.Override.WorkDir,
		SchedulerSubmitFlags: schedulerSubmitFlags,
	}

	if req.Queue == "" {
		j.Queue = cfg.BackendProvider.SchedulerCommon.DefaultQueue
	} else {
		j.Queue = req.Queue
	}

	appMode, appPath, err := ensureApplication(req.Application)
	if err != nil {
		return nil, fmt.Errorf("ensure application failed, %w", err)
	}

	j.AppMode = appMode
	if appMode == models.ImageAppMode {
		j.SingularityImage = appPath
	} else if appMode == models.LocalAppMode {
		j.AppPath = appPath
	}

	if req.Resource.CoresPerNode != nil {
		j.CoresPerNode = int64(*req.Resource.CoresPerNode) //先从job请求参数里读取单节点核数
	} else { //再从配置文件里读取单节点核数
		j.CoresPerNode = int64(cfg.BackendProvider.SchedulerCommon.CoresPerNode[j.Queue])
	}

	return j, nil
}

func parseEnvs(env map[string]string) (string, error) {
	envs := make([]string, 0, len(env))
	for k, v := range env {
		envs = append(envs, fmt.Sprintf("%s=%s", k, v))
	}

	return jsoniter.MarshalToString(envs)
}

// image:application:0.0.1
// local:/app_dir/application
func ensureApplication(application string) (appMode models.AppMode, appPath string, err error) {
	fields := strings.Split(application, ":")
	if len(fields) < 2 {
		err = fmt.Errorf("invalid application format, %s", application)
	}

	appMode, err = models.StrToAppMode(fields[0])
	if err != nil {
		log.Errorf("appMode [%s] unsupported, %v", fields[0], err)
		return
	}

	appPath = strings.Join(fields[1:], ":")
	return
}

func CancelJob(c *gin.Context) {
	logger := trace.GetLogger(c)

	req := new(job.SystemCancelRequest)
	err := c.ShouldBindUri(req)
	if err = response.BadRequestIfError(c, err, errorcode.InvalidArgument); err != nil {
		err = fmt.Errorf("cancel job bad request, %w", err)
		logger.Error(err)
		return
	}

	jobId, err := snowflake.ParseString(req.JobID)
	if err = response.BadRequestIfError(c, err, errorcode.InvalidJobID); err != nil {
		err = fmt.Errorf("parse JobID %s to snowflake id failed, %w", req.JobID, err)
		logger.Error(err)
		return
	}

	s, err := getState(c)
	if err = response.InternalErrorIfError(c, err, errorcode.InternalServerError); err != nil {
		err = fmt.Errorf("get state from gin ctx failed, %w", err)
		logger.Error(err)
		return
	}

	// FIXME for update transaction
	ctx := c.Request.Context()
	exist, j, err := dao.Default.GetJob(ctx, jobId.Int64())
	if err = response.InternalErrorIfError(c, err, errorcode.DatabaseInternalServerError); err != nil {
		err = fmt.Errorf("get job failed where job_id = %d, %w", jobId, err)
		logger.Error(err)
		return
	}
	if !exist {
		err = fmt.Errorf("job not found")
		logger.Error(err)
		_ = response.NotfoundIfError(c, err, errorcode.JobNotFound)
		return
	}

	jobModel := new(v20230530.JobInHPC)
	if j.State == jobstate.Cancelling || j.State == jobstate.Canceled {
		jobModel = j.ToHTTPModel()
		response.OK(c, jobModel)
		return
	}

	// TODO Completing not support to cancel job for now
	if j.State == jobstate.Completed || j.State == jobstate.Completing {
		err = fmt.Errorf("job state %s cannot be cancelled", j.State)
		logger.Error(err)
		_ = response.ForbiddenIfError(c, err, errorcode.CancelJobForbidden)
		return
	}

	err = dao.Default.TerminateJob(ctx, jobId.Int64())
	if err = response.InternalErrorIfError(c, err, errorcode.InternalServerError); err != nil {
		err = fmt.Errorf("terminate job failed, %w", err)
		logger.Error(err)
		return
	}
	s.Factory.Cancel(jobId.Int64())

	jobModel = j.ToHTTPModel()
	jobModel.Status = jobstate.Cancelling
	response.OK(c, jobModel)
}

func GetJob(c *gin.Context) {
	var err error
	logger := trace.GetLogger(c)

	req := new(job.SystemGetRequest)
	err = bindGetJobReq(req, c)
	if err = response.BadRequestIfError(c, err, errorcode.InvalidArgument); err != nil {
		err = fmt.Errorf("bind get job request failed, %w", err)
		logger.Error(err)
		return
	}

	jobID, err := snowflake.ParseString(req.JobID)
	if err = response.BadRequestIfError(c, err, errorcode.InvalidJobID); err != nil {
		err = fmt.Errorf("parse job id %s to snowflake id failed, %w", req.JobID, err)
		logger.Error(err)
		return
	}

	ctx := c.Request.Context()
	exist, j, err := dao.Default.GetJob(ctx, jobID.Int64())
	if err = response.InternalErrorIfError(c, err, errorcode.DatabaseInternalServerError); err != nil {
		err = fmt.Errorf("get job from db failed, %w", err)
		logger.Error(err)
		return
	}
	if !exist {
		err = fmt.Errorf("job not found")
		_ = response.NotfoundIfError(c, err, errorcode.JobNotFound)
		logger.Error(err)
		return
	}

	response.OK(c, j.ToHTTPModel())
}

func bindGetJobReq(req *job.SystemGetRequest, c *gin.Context) error {
	return c.ShouldBindUri(req)
}

const (
	defaultPageOffset = 0
	defaultPageSize   = 10
)

func GetJobs(c *gin.Context) {
	logger := trace.GetLogger(c)

	ids, err := parseJobIds(c.Query("ids"))
	if err = response.BadRequestIfError(c, err, errorcode.InvalidJobID); err != nil {
		err = fmt.Errorf("parse job ids failed, %w", err)
		logger.Error(err)
		return
	}

	pageOffset, err := parseStrToIntWithDefault(c.Query(util.PageOffsetKey), defaultPageOffset)
	if err = response.BadRequestIfError(c, err, errorcode.InvalidPageOffset); err != nil {
		err = fmt.Errorf("invalid PageOffset %s, err: %w", c.Query(util.PageOffsetKey), err)
		logger.Error(err)
		return
	}

	pageSize, err := parseStrToIntWithDefault(c.Query(util.PageSizeKey), defaultPageSize)
	if err = response.BadRequestIfError(c, err, errorcode.InvalidPageSize); err != nil {
		err = fmt.Errorf("invalid PageSize %s, err: %w", c.Query(util.PageSizeKey), err)
		logger.Error(err)
		return
	}

	jobs, err := dao.Default.GetJobs(c.Request.Context(), dao.GetJobsArg{
		IDs:        ids,
		State:      c.Query("Status"),
		PageOffset: pageOffset,
		PageSize:   pageSize,
	})
	if err = response.InternalErrorIfError(c, err, errorcode.DatabaseInternalServerError); err != nil {
		err = fmt.Errorf("get jobs from db failed, %w", err)
		logger.Error(err)
		return
	}

	// TODO pager
	resp := &job.SystemListResponseData{
		Jobs: make([]*v20230530.JobInHPC, 0),
	}
	for _, j := range jobs {
		resp.Jobs = append(resp.Jobs, j.ToHTTPModel())
	}

	response.OK(c, resp)
}

func parseJobIds(s string) ([]int64, error) {
	if s == "" {
		return nil, nil
	}

	ids := make([]int64, 0)
	for _, str := range strings.Split(s, ",") {
		jobId, err := snowflake.ParseString(str)
		if err != nil {
			return nil, fmt.Errorf("parse job id %s to snowflake id failed, %w", str, err)
		}

		ids = append(ids, jobId.Int64())
	}

	return ids, nil
}

// if s == "", return defaultV
func parseStrToIntWithDefault(s string, defaultV int) (int, error) {
	if s == "" {
		return defaultV, nil
	}

	v, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("parse [%s] to int failed, %w", s, err)
	}

	return v, nil
}

func DeleteJob(c *gin.Context) {
	logger := trace.GetLogger(c)

	req, err := bindDeleteJobReq(c)
	if err = response.BadRequestIfError(c, err, errorcode.InvalidArgument); err != nil {
		err = fmt.Errorf("bind delete job request failed, %w", err)
		logger.Error(err)
		return
	}

	exist := true
	allowDelete := true
	js := jobstate.State("")
	allowDeleteStates := []jobstate.State{
		jobstate.Completed,
		jobstate.Canceled,
		jobstate.Failed,
	}
	err = with.DefaultTransaction(c.Request.Context(), func(ctx context.Context) error {
		ex, j, e := dao.Default.GetJob(ctx, snowflake.MustParseString(req.JobID).Int64())
		if e != nil {
			return fmt.Errorf("get job where job_id = %s failed, %w", req.JobID, err)
		}
		if !ex {
			exist = false
			return nil
		}

		if !stateInAllowDeleteStates(j.State, allowDeleteStates) {
			allowDelete = false
			js = j.State
			return nil
		}

		if e = dao.Default.DeleteJob(ctx, j.Id); e != nil {
			return fmt.Errorf("delete job where job_id = %s failed, %w", req.JobID, err)
		}

		return nil
	})
	if err = response.InternalErrorIfError(c, err, errorcode.DatabaseInternalServerError); err != nil {
		err = fmt.Errorf("database error, %w", err)
		logger.Error(err)
		return
	}
	if !exist {
		err = fmt.Errorf("job not found")
		_ = response.NotfoundIfError(c, err, errorcode.JobNotFound)
		logger.Error(err)
		return
	}
	if !allowDelete {
		err = fmt.Errorf("job state [%s] not allowed to delete", js)
		_ = response.ForbiddenIfError(c, err, errorcode.DeleteJobForbidden)
		logger.Error(err)
		return
	}

	response.OK(c, nil)
	return
}

func bindDeleteJobReq(c *gin.Context) (*job.SystemDeleteRequest, error) {
	req := new(job.SystemDeleteRequest)
	if err := c.ShouldBindUri(req); err != nil {
		return nil, fmt.Errorf("should bind uri failed, %w", err)
	}

	return req, nil
}

func stateInAllowDeleteStates(s jobstate.State, allowDeleteStates []jobstate.State) bool {
	for _, allow := range allowDeleteStates {
		if s == allow {
			return true
		}
	}

	return false
}
