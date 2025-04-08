package impl

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/yuansuan/ticp/common/go-kit/logging"

	mainconfig "github.com/yuansuan/ticp/PSP/psp/cmd/config"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi/config"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/project"
	"github.com/yuansuan/ticp/PSP/psp/internal/job/consts"
	"github.com/yuansuan/ticp/PSP/psp/internal/job/dao"
	"github.com/yuansuan/ticp/PSP/psp/internal/job/dao/model"
	"github.com/yuansuan/ticp/PSP/psp/internal/job/dto"
	"github.com/yuansuan/ticp/PSP/psp/internal/job/service"
	"github.com/yuansuan/ticp/PSP/psp/internal/job/service/client"
	"github.com/yuansuan/ticp/PSP/psp/internal/job/util"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/strutil"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/timeutil"
	"github.com/yuansuan/ticp/PSP/psp/pkg/xtype"
)

var (
	JobNotFoundError = errors.New("the job is not found")
)

type jobServiceImpl struct {
	rpc            *client.GRPC
	sid            *snowflake.Node
	localCfg       *config.Local
	localAPI       *openapi.OpenAPI
	jobDao         dao.JobDao
	jobAttrDao     dao.JobAttrDao
	jobTimelineDao dao.JobTimelineDao
}

func NewJobService() (service.JobService, error) {
	jobDao, err := dao.NewJobDao()
	if err != nil {
		return nil, err
	}

	jobAttrDao, err := dao.NewJobAttrDao()
	if err != nil {
		return nil, err
	}

	jobTimelineDao, err := dao.NewJobTimelineDao()
	if err != nil {
		return nil, err
	}

	localAPI, err := openapi.NewLocalAPI()
	if err != nil {
		return nil, err
	}

	sid, err := snowflake.GetInstance()
	if err != nil {
		return nil, err
	}

	rpc := client.GetInstance()

	localCfg := config.GetConfig().Local
	if localCfg == nil {
		return nil, errors.New("openapi configuration is invalid")
	}

	jobService := &jobServiceImpl{
		rpc:            rpc,
		sid:            sid,
		localCfg:       localCfg,
		localAPI:       localAPI,
		jobDao:         jobDao,
		jobAttrDao:     jobAttrDao,
		jobTimelineDao: jobTimelineDao,
	}

	return jobService, nil
}

// GetJobDetail 获取作业详情
func (s *jobServiceImpl) GetJobDetail(ctx context.Context, jobID string) (*dto.JobDetailInfo, error) {
	logger := logging.GetLogger(ctx)

	sid, err := snowflake.ParseString(jobID)
	if err != nil {
		return nil, err
	}

	has, job, err := s.jobDao.GetJobDetail(ctx, sid)
	if err != nil {
		logger.Errorf("get the job [%v] detail err: %v", sid.Int64(), err)
		return nil, err
	}

	if !has {
		logger.Infof("the job [%v] is not found", sid.Int64())
		return nil, JobNotFoundError
	}

	timelines, err := s.GetJobTimeline(ctx, job.Id.String(), job.UploadTaskId, job.State, job.DataState)
	if err != nil {
		logger.Errorf("get job [%v] timeline err: %v", job.Id, err)
	}

	detail := util.ConvertJob2Detail(job)

	if timelines != nil {
		detail.Timelines = timelines
	}

	envs, err := s.ParseEnvs(ctx, job.Id, consts.JobAttrKeySubmitEnvs)
	if err != nil {
		logger.Errorf("get job [%v] timeline err: %v", job.Id, err)
	}

	forbidEnv := envs[consts.ForbidDownload]
	if !strutil.IsEmpty(forbidEnv) {
		detail.FileFilterRegs = strings.Split(forbidEnv, consts.Semicolon)
	}

	return detail, nil
}

// GetJobSetDetail 获取作业集详情
func (s *jobServiceImpl) GetJobSetDetail(ctx context.Context, jobSetID string, loginUserID snowflake.ID) (*dto.JobSetInfo, []*dto.JobListInfo, error) {
	logger := logging.GetLogger(ctx)

	if jobSetID == "" {
		return nil, nil, fmt.Errorf("the job set: [%v] not exist", jobSetID)
	}

	isAdmin, _, err := s.checkAndGetProjects(ctx, loginUserID)
	if err != nil {
		logger.Errorf("check and get projects err: %v", err)
		return nil, nil, err
	}

	jobs, total, err := s.jobDao.GetJobList(ctx, &dto.JobFilter{JobSetID: jobSetID}, &xtype.Page{Index: 1, Size: 1000}, nil, isAdmin, loginUserID)
	if err != nil {
		logger.Errorf("get the job list err: %v", err)
		return nil, nil, err
	}

	jobSetInfo := &dto.JobSetInfo{JobCount: total}

	startTime, endTime, submitTime := time.Time{}, time.Time{}, time.Time{}
	successCount, failureCount, execDuration, endTimeFlag := 0, 0, 0, true
	jobList := make([]*dto.JobListInfo, 0, len(jobs))
	for i, job := range jobs {
		if i == 0 {
			jobSetInfo.ProjectId = job.ProjectId.String()
			jobSetInfo.ProjectName = job.ProjectName
			jobSetInfo.JobSetId = job.JobSetId.String()
			jobSetInfo.JobSetName = job.JobSetName
			jobSetInfo.JobType = job.Type
			jobSetInfo.AppId = job.AppId.String()
			jobSetInfo.AppName = job.AppName
			jobSetInfo.UserId = job.UserId.String()
			jobSetInfo.UserName = job.UserName
		}
		if startTime.IsZero() && !job.StartTime.Equal(timeutil.DefaultDateTime) {
			startTime = job.StartTime
		}
		if endTime.IsZero() && !job.EndTime.Equal(timeutil.DefaultDateTime) {
			endTime = job.EndTime
		}
		if submitTime.IsZero() && !job.SubmitTime.Equal(timeutil.DefaultDateTime) {
			submitTime = job.SubmitTime
		}

		execDuration += job.ExecDuration

		if job.State == consts.APIJobStateCompleted {
			successCount++
		} else if job.State == consts.APIJobStateFailed {
			failureCount++
		}

		if endTimeFlag && job.EndTime.Equal(timeutil.DefaultDateTime) {
			endTimeFlag = false
		}
		if !job.EndTime.Equal(timeutil.DefaultDateTime) && job.StartTime.Before(startTime) {
			startTime = job.StartTime
		}
		if endTimeFlag && !job.EndTime.Equal(timeutil.DefaultDateTime) && job.EndTime.After(endTime) {
			endTime = job.EndTime
		}
		if !job.EndTime.Equal(timeutil.DefaultDateTime) && job.SubmitTime.Before(submitTime) {
			submitTime = job.SubmitTime
		}

		jobInfo := util.ConvertJob2ListInfo(job)
		if jobInfo != nil {
			jobList = append(jobList, jobInfo)
		}
	}

	jobSetInfo.ExecDuration = strconv.Itoa(execDuration)
	jobSetInfo.SuccessCount = int64(successCount)
	jobSetInfo.FailureCount = int64(failureCount)
	jobSetInfo.StartTime = timeutil.DefaultFormatTime(startTime)
	if endTimeFlag {
		jobSetInfo.EndTime = timeutil.DefaultFormatTime(endTime)
	} else {
		jobSetInfo.EndTime = ""
	}

	return jobSetInfo, jobList, err
}

// GetJobDetailByOutID 获取作业详情
func (s *jobServiceImpl) GetJobDetailByOutID(ctx context.Context, outJobID, jobType string) (*model.Job, error) {
	logger := logging.GetLogger(ctx)

	has, job, err := s.jobDao.GetJobByOutID(ctx, outJobID, jobType)
	if err != nil {
		logger.Errorf("get job [%v] err: %v", outJobID, err)
		return nil, err
	}

	if !has {
		logger.Infof("the job [%v] is not found", outJobID)
		return nil, JobNotFoundError
	}

	return job, nil
}

// GetJobList 获取作业列表
func (s *jobServiceImpl) GetJobList(ctx context.Context, filter *dto.JobFilter, page *xtype.Page, orderSort *xtype.OrderSort, loginUserID snowflake.ID) ([]*model.Job, int64, error) {
	logger := logging.GetLogger(ctx)

	isAdmin, projectIds, err := s.checkAndGetProjects(ctx, loginUserID)
	if err != nil {
		logger.Errorf("check and get projects err: %v", err)
		return nil, 0, err
	}

	if !isAdmin {
		if len(filter.ProjectIDs) == 0 {
			filter.ProjectIDs = projectIds
		}
		filter.UserNames = make([]string, 0)
	}

	jobs, total, err := s.jobDao.GetJobList(ctx, filter, page, orderSort, isAdmin, loginUserID)
	if err != nil {
		logger.Errorf("get the job list err: %v", err)
		return nil, 0, err
	}

	return jobs, total, nil
}

// GetJobSetNameList 获取作业集名称列表
func (s *jobServiceImpl) GetJobSetNameList(ctx context.Context, computeType string, loginUserID snowflake.ID) ([]string, error) {
	logger := logging.GetLogger(ctx)

	isAdmin, projectIds, err := s.checkAndGetProjects(ctx, loginUserID)
	if err != nil {
		logger.Errorf("check and get projects err: %v", err)
		return nil, err
	}

	appNames, err := s.jobDao.GetJobSetNameList(ctx, projectIds, computeType, isAdmin, loginUserID)
	if err != nil {
		logger.Errorf("get the job set name list err: %v", err)
		return nil, err
	}

	if appNames == nil {
		appNames = make([]string, 0)
	}

	return appNames, nil
}

// GetJobAppNameList 获取作业应用名称列表
func (s *jobServiceImpl) GetJobAppNameList(ctx context.Context, computeType string, loginUserID snowflake.ID) ([]string, error) {
	logger := logging.GetLogger(ctx)

	isAdmin, projectIds, err := s.checkAndGetProjects(ctx, loginUserID)
	if err != nil {
		logger.Errorf("check and get projects err: %v", err)
		return nil, err
	}

	appNames, err := s.jobDao.GetJobAppNameList(ctx, projectIds, computeType, isAdmin, loginUserID)
	if err != nil {
		logger.Errorf("get the job app name list err: %v", err)
		return nil, err
	}

	return appNames, nil
}

func (s *jobServiceImpl) checkAndGetProjects(ctx context.Context, loginUserID snowflake.ID) (bool, []string, error) {
	checkPermissionResponse, err := client.GetInstance().Project.CheckUserOperatorProjectsPermission(ctx,
		&project.CheckUserOperatorProjectsPermissionRequest{
			UserId: loginUserID.String(),
		},
	)
	if err != nil {
		return false, nil, err
	}

	// 管理员直接返回
	if checkPermissionResponse.Pass {
		return true, make([]string, 0), nil
	}

	projectIds, err := s.getProjectIdsByUserID(ctx, loginUserID)
	if err != nil {
		return false, nil, err
	}

	return false, projectIds, nil
}

func (s *jobServiceImpl) getProjectIdsByUserID(ctx context.Context, loginUserID snowflake.ID) ([]string, error) {
	projectsResponse, err := client.GetInstance().Project.GetMemberProjectsByUserId(ctx, &project.GetMemberProjectsByUserIdRequest{UserId: loginUserID.String(), IncludeDefault: true})
	if err != nil {
		return nil, err
	}

	projectIds := make([]string, 0, len(projectsResponse.Projects))
	for _, v := range projectsResponse.Projects {
		projectIds = append(projectIds, v.ProjectId)
	}

	return projectIds, nil
}

// GetJobComputeTypeList 获取作业计算类型列表
func (s *jobServiceImpl) GetJobComputeTypeList(ctx context.Context, loginUserID snowflake.ID) ([]*dto.ComputeTypeName, error) {
	logger := logging.GetLogger(ctx)

	isAdmin, _, err := s.checkAndGetProjects(ctx, loginUserID)
	if err != nil {
		logger.Errorf("check and get projects err: %v", err)
		return nil, err
	}

	computeTypes, err := s.jobDao.GetJobComputeTypeList(ctx, isAdmin, loginUserID)
	if err != nil {
		logger.Errorf("get the job compute type list err: %v", err)
		return nil, err
	}

	competeTypeNameMap := mainconfig.Custom.Main.ComputeTypeNames
	computeTypeNameList := make([]*dto.ComputeTypeName, 0, len(computeTypes))
	for _, v := range computeTypes {
		vName := fmt.Sprintf("[未配置]%v", v)
		if name, ok := competeTypeNameMap[v]; ok {
			vName = name
		} else {
			continue
		}
		computeTypeNameList = append(computeTypeNameList, &dto.ComputeTypeName{ComputeType: v, ShowName: vName})
	}

	return computeTypeNameList, nil
}

// GetJobUserNameList 获取作业用户名称列表
func (s *jobServiceImpl) GetJobUserNameList(ctx context.Context, computeType string, loginUserID snowflake.ID) ([]string, error) {
	logger := logging.GetLogger(ctx)

	isAdmin, projectIds, err := s.checkAndGetProjects(ctx, loginUserID)
	if err != nil {
		logger.Errorf("check and get projects err: %v", err)
		return nil, err
	}

	userNames, err := s.jobDao.GetJobUserNameList(ctx, projectIds, computeType, isAdmin, loginUserID)
	if err != nil {
		logger.Errorf("get the job username list err: %v", err)
		return nil, err
	}

	return userNames, nil
}

// GetJobQueueNameList 获取作业队列名称列表
func (s *jobServiceImpl) GetJobQueueNameList(ctx context.Context, computeType string, loginUserID snowflake.ID) ([]string, error) {
	logger := logging.GetLogger(ctx)

	isAdmin, projectIds, err := s.checkAndGetProjects(ctx, loginUserID)
	if err != nil {
		logger.Errorf("check and get projects err: %v", err)
		return nil, err
	}

	queueNames, err := s.jobDao.GetJobQueueNameList(ctx, projectIds, computeType, isAdmin, loginUserID)
	if err != nil {
		logger.Errorf("get the job queue name list err: %v", err)
		return nil, err
	}

	return queueNames, nil
}

// AppJobNum 获取应用作业数
func (s *jobServiceImpl) AppJobNum(ctx context.Context, start, end int64) ([]*dto.AppJobInfo, int, error) {
	logger := logging.GetLogger(ctx)
	//获取应用总数
	appNames, err := s.jobDao.GetJobAppNameList(ctx, nil, "", true, 0)
	if err != nil {
		logger.Errorf("get app total num err: %v", err)
		return nil, 0, err
	}
	// 获取应用作业数
	appJobInfo, err := s.jobDao.GetAppJobNum(ctx, start, end)
	if err != nil {
		logger.Errorf("get app job num err: %v", err)
		return nil, 0, err
	}
	return appJobInfo, len(appNames), nil
}

// UserJobNum 用户作业数
func (s *jobServiceImpl) UserJobNum(ctx context.Context, start, end int64) ([]*dto.UserJobInfo, error) {
	logger := logging.GetLogger(ctx)

	// 获取用户作业数
	resp, err := s.jobDao.GetUserJobNum(ctx, start, end)
	if err != nil {
		logger.Errorf("get user job num err: %v", err)
		return nil, err
	}
	return resp, nil
}

// GetJobCPUTimeMetric 作业核时运行指标统计
func (s *jobServiceImpl) GetJobCPUTimeMetric(ctx context.Context, filter *dto.JobMetricFiler) (*dto.JobCPUTimeMetric, error) {
	status := []string{consts.JobStateCompleted, consts.JobStateTerminated}
	appMetrics, err := s.jobDao.GetJobCPUTimeMetric(ctx, filter, "app_name", status)
	if err != nil {
		return nil, err
	}

	userMetrics, err := s.jobDao.GetJobCPUTimeMetric(ctx, filter, "user_name", status)
	if err != nil {
		return nil, err
	}

	return &dto.JobCPUTimeMetric{
		AppMetrics:  appMetrics,
		UserMetrics: userMetrics,
	}, nil
}

// GetJobCountMetric 应用和用户数量统计指标
func (s *jobServiceImpl) GetJobCountMetric(ctx context.Context, filter *dto.JobMetricFiler) (*dto.JobCountMetric, error) {
	status := []string{consts.JobStateCompleted, consts.JobStateTerminated}
	appMetrics, err := s.jobDao.GetJobCountMetric(ctx, filter, "app_name", status)
	if err != nil {
		return nil, err
	}

	userMetrics, err := s.jobDao.GetJobCountMetric(ctx, filter, "user_name", status)
	if err != nil {
		return nil, err
	}

	return &dto.JobCountMetric{
		AppCountMetrics:  appMetrics,
		UserCountMetrics: userMetrics,
	}, nil
}

// GetJobDeliverCount 作业提交数量指标
func (s *jobServiceImpl) GetJobDeliverCount(ctx context.Context, filter *dto.JobMetricFiler) (*dto.JobCountMetric, error) {
	jobMetrics, err := s.jobDao.GetJobDeliverCount(ctx, filter, consts.JobDeliveryCountJob)
	if err != nil {
		return nil, err
	}

	userMetrics, err := s.jobDao.GetJobDeliverCount(ctx, filter, consts.JobDeliveryCountUser)
	if err != nil {
		return nil, err
	}

	return &dto.JobCountMetric{
		AppCountMetrics:  jobMetrics,
		UserCountMetrics: userMetrics,
	}, nil
}

// GetJobWaitStatistic 作业等待指标
func (s *jobServiceImpl) GetJobWaitStatistic(ctx context.Context, filter *dto.JobMetricFiler) (*dto.JobWaitStatistic, error) {
	status := []string{consts.JobStateCompleted, consts.JobStateTerminated}

	timeStatisticAvg, err := s.jobDao.GetJobWaitStatistic(ctx, filter, consts.JobWaitTimeStatisticAvg, status)
	if err != nil {
		return nil, err
	}

	timeStatisticMax, err := s.jobDao.GetJobWaitStatistic(ctx, filter, consts.JobWaitTimeStatisticMax, status)
	if err != nil {
		return nil, err
	}

	timeStatisticTotal, err := s.jobDao.GetJobWaitStatistic(ctx, filter, consts.JobWaitTimeStatisticTotal, status)
	if err != nil {
		return nil, err
	}

	numStatistic, err := s.jobDao.GetJobWaitStatistic(ctx, filter, consts.JobWaitNumStatistic, status)
	if err != nil {
		return nil, err
	}

	return &dto.JobWaitStatistic{
		JobWaitTimeStatisticAvg:   timeStatisticAvg,
		JobWaitTimeStatisticTotal: timeStatisticTotal,
		JobWaitTimeStatisticMax:   timeStatisticMax,
		JobWaitNumStatistic:       numStatistic,
	}, nil
}

// GetJobStatus 作业状态统计
func (s *jobServiceImpl) GetJobStatus(ctx context.Context) (map[string]int64, error) {
	// 获取24小时内 已完成、已终止的作业数
	startTime := time.Now().Add(-24 * time.Hour)
	compAndTermJobNums, err := s.jobDao.GetJobStatusNum(ctx, startTime.Unix(), []string{consts.JobStateCompleted, consts.JobStateFailed, consts.JobStateBurstFailed})
	if err != nil {
		return nil, err
	}

	// 获取正在运行和正在等待的作业数
	runAndWaitJobNums, err := s.jobDao.GetJobStatusNum(ctx, 0, []string{consts.JobStateRunning, consts.JobStatePending})
	if err != nil {
		return nil, err
	}

	jobStatusMap := make(map[string]int64)
	for _, jobNum := range compAndTermJobNums {
		jobStatusMap[jobNum.State] = jobNum.Num
	}
	for _, jobNum := range runAndWaitJobNums {
		jobStatusMap[jobNum.State] = jobNum.Num
	}
	jobStatusMap[consts.JobStateFailed] = jobStatusMap[consts.JobStateFailed] + jobStatusMap[consts.JobStateBurstFailed]
	delete(jobStatusMap, consts.JobStateBurstFailed)

	return jobStatusMap, nil
}

func (s *jobServiceImpl) ParseEnvs(ctx context.Context, jobID snowflake.ID, envKey string) (map[string]string, error) {
	logger := logging.GetLogger(ctx)

	has, attr, err := s.jobAttrDao.GetJobAttrByKey(ctx, jobID, envKey)
	if err != nil {
		return nil, err
	}

	if !has {
		return nil, errors.Errorf("the '%v' job attr is not found, the job id is %v", envKey, jobID.Int64())
	}

	var envs map[string]string
	if err = json.Unmarshal([]byte(attr.Value), &envs); err != nil {
		logger.Errorf("unmarshal job attr err: %v, the attr value: %v", err, attr.Value)
		return nil, err
	}

	return envs, nil
}

// GetOutIDByJobID 根据作业id获取paas作业id
func (s *jobServiceImpl) GetOutIDByJobID(ctx context.Context, jobID snowflake.ID) (string, error) {
	logger := logging.GetLogger(ctx)

	has, job, err := s.jobDao.GetJobDetail(ctx, jobID)

	if err != nil {
		logger.Errorf("get job [%v] err: %v", jobID, err)
		return "", err
	}

	if !has {
		logger.Infof("the job [%v] is not found", jobID)
		return "", JobNotFoundError
	}

	return job.OutJobId, nil
}
