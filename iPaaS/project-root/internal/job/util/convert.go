package util

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/config"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/consts"
	"go.uber.org/zap"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	hpc "github.com/yuansuan/ticp/common/project-root-api/hpc/v1/job"
	adminjobcreate "github.com/yuansuan/ticp/common/project-root-api/job/v1/admin/jobcreate"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/jobcreate"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/jobpreschedule"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/system/jobneedsyncfile"
	schema "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao/models"
)

// InvalidTime ...
var InvalidTime, _ = time.Parse(time.RFC3339, "1970-01-01T00:00:00+00:00")

// ModelToOpenAPIJob 转换为openapi作业 普通作业
// get/batchget/list job用
func ModelToOpenAPIJob(job *models.Job) *schema.JobInfo {
	if job == nil {
		return nil
	}
	var stateDesc string
	state, b := consts.GetStateBySubState(job.SubState)
	if b {
		stateDesc = state.StateString()
	}
	endTime := job.EndTime
	executionDuration := job.ExecutionDuration
	return &schema.JobInfo{
		ID:            job.ID.String(),
		Name:          job.Name,
		JobState:      stateDesc,
		StateReason:   job.StateReason,
		FileSyncState: job.FileSyncState,
		AllocResource: &schema.AllocResource{
			Cores:  int(job.ResourceAssignCpus),
			Memory: int(job.ResourceAssignMemory),
		},
		AllocType:        job.AllocType,
		ExecHostNum:      job.ExecHostNum,
		Zone:             job.Zone,
		Workdir:          job.WorkDir,
		Parameters:       job.Params,
		NoRound:          job.NoRound,
		PreScheduleID:    job.PreScheduleID,
		OutputDir:        job.OutputDir,
		NoNeededPaths:    job.NoNeededPaths,
		NeededPaths:      job.NeededPaths,
		PendingTime:      ModelTimeToString(job.PendingTime),
		RunningTime:      ModelTimeToString(job.RunningTime),
		TerminatingTime:  ModelTimeToString(job.TerminatingTime),
		SuspendingTime:   ModelTimeToString(job.SuspendingTime),
		SuspendedTime:    ModelTimeToString(job.SuspendedTime),
		EndTime:          ModelTimeToString(endTime),
		CreateTime:       ModelTimeToString(job.CreateTime),
		UpdateTime:       ModelTimeToString(job.UpdateTime),
		FileReadyTime:    ModelTimeToString(job.UploadTime),
		TransmittingTime: ModelTimeToString(job.TransmittingTime),
		TransmittedTime:  ModelTimeToString(job.DownloadTime),
		DownloadProgress: &schema.DownloadProgress{
			Progress: &schema.Progress{
				TotalSize: int(job.DownloadFileSizeTotal),
				Progress:  calPercent(job.DownloadFileSizeCurrent, job.DownloadFileSizeTotal),
			},
		},
		UploadProgress: &schema.UploadProgress{
			Progress: &schema.Progress{
				TotalSize: int(job.UploadFileSizeTotal),
				Progress:  calPercent(job.UploadFileSizeCurrent, job.UploadFileSizeTotal),
			},
		},
		ExecutionDuration: int(executionDuration),
		ExitCode:          job.ExitCode,
		IsSystemFailed:    job.IsSystemFailed == 1,
		StdoutPath:        AddSuffixSlash(job.WorkDir) + consts.DefaultStdout,
		StderrPath:        AddSuffixSlash(job.WorkDir) + consts.DefaultStderr,
	}
}

// ModelToAdminOpenAPIJob 转换为openapi作业 admin作业
func ModelToAdminOpenAPIJob(job *models.Job) *schema.AdminJobInfo {
	openApiJob := ModelToOpenAPIJob(job)
	if openApiJob == nil {
		return nil
	}
	var jobInfo schema.AdminJobInfo
	jobInfo.JobInfo = *openApiJob
	// admin 部分
	jobInfo.Queue = job.Queue
	jobInfo.Priority = job.Priority
	jobInfo.OriginJobID = job.OriginJobID
	jobInfo.ExecHosts = job.ExecHosts
	jobInfo.SubmitTime = ModelTimeToString(job.SubmitTime)
	jobInfo.UserID = job.UserID.String()
	jobInfo.HPCJobID = job.HPCJobID
	jobInfo.IsDeleted = job.IsDeleted == 1
	return &jobInfo
}

// TrimAdminParams 去除adminParams信息
func TrimAdminParams(jobInfo *schema.JobInfo) (*schema.JobInfo, error) {
	admminParams := models.AdminParams{}
	err := json.Unmarshal([]byte(jobInfo.Parameters), &admminParams)
	if err != nil {
		return jobInfo, err
	}
	params := admminParams.Params
	paramsStr, err := json.Marshal(params)
	if err != nil {
		return jobInfo, err
	}
	jobInfo.Parameters = string(paramsStr)
	return jobInfo, nil
}

// ModelToOpenAPIJobResidual 获取作业残差图
func ModelToOpenAPIJobResidual(residual *models.Residual) (*schema.Residual, error) {
	if residual.Content == "" {
		return &schema.Residual{}, nil
	}
	return ResidualUnmarshal(residual.Content)
}

// ModelToOpenAPIJobMonitorCharts 获取作业监控图表
func ModelToOpenAPIJobMonitorCharts(monitorChart *models.MonitorChart) ([]*schema.MonitorChart, error) {
	if monitorChart.Content == "" {
		return nil, nil
	}
	return MonitorChartUnmarshal(monitorChart.Content)
}

// ModelTimeToString 时间转字符串
func ModelTimeToString(timeInput time.Time) string {
	var cstSh, _ = time.LoadLocation("Asia/Shanghai")
	formatTime := timeInput.In(cstSh).Format(time.RFC3339)
	return formatTime
}

func calPercent(numerator int64, denominator int64) int {
	if denominator == 0 {
		return 0
	}
	percentage := (numerator * 100) / denominator
	return int(percentage)
}

// AssembleHPCJobRequest model转hpc req
func AssembleHPCJobRequest(ctx context.Context, logger *zap.SugaredLogger, job *models.Job,
	appImage string, envVars, schedulerSubmitFlags map[string]string,
	noTransfer, localImage bool, coresPerNode int) hpc.SystemPostRequest {
	inputSrc := job.InputDir
	input := schema.JobInHPCInputStorage{
		Type: consts.FileType(job.InputType).ToHPCFileType(),
		Src:  inputSrc,
	}
	inputs := []schema.JobInHPCInputStorage{input}
	output := &schema.JobInHPCOutputStorage{
		Type:          consts.FileType(job.OutputType).ToHPCFileType(),
		Dst:           job.OutputDir,
		NoNeededPaths: job.NoNeededPaths,
		NeededPaths:   job.NeededPaths,
	}
	if job.InputType == "" {
		inputs = nil
	}
	if job.OutputType == "" {
		output = nil
	}
	if noTransfer {
		inputs = nil
		output = nil
	}
	resource := schema.JobInHPCResource{
		Cores:     int(job.ResourceUsageCpus), // 用户申请核数
		AllocType: job.AllocType,
	}
	if coresPerNode > 0 {
		resource.CoresPerNode = &coresPerNode // 用户指定每个节点的核数
	}
	return hpc.SystemPostRequest{
		Application: AddAppImagePrefix(appImage, localImage),
		Environment: envVars,
		Command:     job.Command,
		Override: schema.JobInHPCOverride{
			Enable:  true, // 修改workdir
			WorkDir: job.WorkDir,
		},
		Resource:     resource,
		Inputs:       inputs,
		Output:       output,
		Queue:        job.Queue,
		IdempotentID: job.ID.String(), // 以paas平台job id作为hpc 的 幂等id
		CustomStateRule: &schema.JobInHPCCustomStateRule{
			KeyStatement: job.CustomStateRuleKeyStatement,
			ResultState:  job.CustomStateRuleResultState,
		},
		JobSchedulerSubmitFlags: schedulerSubmitFlags,
	}
}

// ConvertJobModel 参数转换成model
func ConvertJobModel(ctx context.Context, logger *zap.SugaredLogger, req *jobcreate.Request,
	userID snowflake.ID, jobID snowflake.ID, appInfo *models.Application,
	inputZone, outputZone consts.Zone, schedulerSubmitFlags map[string]string,
	scheduleInfo *models.PreSchedule) (*models.Job, error) {
	now := time.Now()
	adminParams := models.AdminParams{
		Params:                  req.Params,
		JobSchedulerSubmitFlags: schedulerSubmitFlags,
	}
	params, err := json.Marshal(adminParams)
	if err != nil {
		logger.Errorf("json.Marshal error: %v", err)
		return nil, err // internal error
	}
	// 输入处理
	input := &jobcreate.Input{Type: "", Source: "", Destination: ""}
	workDir := ""
	if req.Params.Input != nil { // nil when job preSchedule
		input = req.Params.Input
		inputYsID := ParseYsID(input.Source) // 不带'/'前缀
		if inputYsID == "" {
			return nil, errors.WithMessagef(common.ErrInvalidArgumentInput, "input ys_id is empty")
		}
		inputPath := ParsePath(input.Source) // 带'/'前缀
		if inputPath == "" {
			return nil, errors.WithMessagef(common.ErrInvalidArgumentInput, "input path is empty")
		}
		workDir = input.Destination
		// workDir如果为空，默认取input.src 里的path
		if workDir == "" {
			workDir = fmt.Sprintf("%s%s", inputYsID, inputPath)
		}
		if req.Params.TmpWorkdir {
			// workDir路径前增加一个唯一的临时前缀目录
			workYsID := ParseYsIDWithOutDomain(workDir) // destYsID不带'/'前缀
			if workYsID == "" {
				return nil, errors.WithMessagef(common.ErrInvalidArgumentInput, "input dest ys_id is empty")
			}
			workPath := ParsePathWithOutDomain(workDir) // destPath带'/'前缀
			if workPath == "" {
				return nil, errors.WithMessagef(common.ErrInvalidArgumentInput, "input dest ys_id is empty")
			}
			tmpworkdir := consts.TmpWorkdirPrefix + jobID.String()
			// workDir的第一级是user_id, 添加tmpworkdir需要加在第二级
			// - ys_id/tmpworkdir/job_id/path
			workDir = fmt.Sprintf("%s/%s%s", workYsID, tmpworkdir, workPath)
		}

		// 如果路径末尾没有/，则加上/
		workDir = AddSuffixSlash(workDir)
	}
	// 输出处理
	output := req.Params.Output
	if output == nil {
		output = &jobcreate.Output{
			Type:    "",
			Address: "",
		}
	} else {
		// 如果路径末尾没有/，则加上/
		if !strings.HasSuffix(output.Address, "/") {
			output.Address = fmt.Sprintf("%s/", output.Address)
		}
	}
	zone := req.Zone
	zones := config.GetConfig().Zones
	if zone != "" {
		zoneInfo, ok := zones[zone]
		if ok && zoneInfo.HPCEndpoint != "" {
			workDir = fmt.Sprintf("%s/%s", zoneInfo.HPCEndpoint, workDir)
		}
	}
	command := req.Params.Application.Command
	cores := int64(*req.Params.Resource.Cores)
	memory := int64(*req.Params.Resource.Memory)
	allocType := req.AllocType
	timeout := req.Timeout
	if timeout == 0 {
		timeout = -1
	}
	// 初始化model.Job数据
	state := consts.SubStateScheduling
	stateReason := consts.UserSubmitReason
	pendingTime := now
	// SubmitWithSuspend 为true时，作业状态为InitiallySuspended
	if req.Params.SubmitWithSuspend {
		state = consts.SubStateInitiallySuspended
		stateReason = consts.InitiallySuspendedReason
		pendingTime = InvalidTime
	}
	sr := consts.NewStateReason(state, stateReason)
	customStateRule := &jobcreate.CustomStateRule{
		KeyStatement: "",
		ResultState:  "",
	}
	if req.Params.CustomStateRule != nil {
		customStateRule = req.Params.CustomStateRule
	}
	preScheduleID := ""
	if scheduleInfo != nil { // 预调度这里应该覆盖的参数
		zone = scheduleInfo.Zone
		workDir = scheduleInfo.WorkDir
		preScheduleID = scheduleInfo.ID.String()
	}
	job := &models.Job{
		// 作业基本信息
		ID:        jobID,
		Name:      req.Name,
		Comment:   req.Comment,
		UserID:    userID,
		JobSource: "", // 预留
		// 状态信息
		State:         state.State,
		SubState:      state.SubState,
		StateReason:   sr.String(),
		ExitCode:      "",
		FileSyncState: consts.FileSyncStateNone.String(),
		// 参数信息
		Params:                      string(params),
		UserZone:                    zone,
		Timeout:                     timeout,
		FileClassifier:              "", // 预留
		ResourceUsageCpus:           cores,
		ResourceUsageMemory:         memory,
		AllocType:                   allocType,
		CustomStateRuleKeyStatement: customStateRule.KeyStatement,
		CustomStateRuleResultState:  customStateRule.ResultState,
		NoRound:                     req.NoRound,
		PreScheduleID:               preScheduleID,
		// 作业运行信息
		Zone:                 zone,
		ResourceAssignCpus:   0,
		ResourceAssignMemory: 0,
		Command:              command,
		WorkDir:              workDir,
		OriginJobID:          "",
		Queue:                "",
		Priority:             0,
		ExecHosts:            "",
		ExecHostNum:          0,
		ExecutionDuration:    0,
		// 文件信息
		InputType:               input.Type,
		InputDir:                input.Source,
		Destination:             input.Destination,
		OutputType:              output.Type,
		OutputDir:               output.Address,
		NoNeededPaths:           output.NoNeededPaths,
		NeededPaths:             output.NeededPaths,
		FileInputStorageZone:    inputZone.String(),
		FileOutputStorageZone:   outputZone.String(),
		DownloadFileSizeTotal:   0,
		DownloadFileSizeCurrent: 0,
		UploadFileSizeTotal:     0,
		UploadFileSizeCurrent:   0,
		// 应用信息
		AppID:   appInfo.ID,
		AppName: appInfo.Name,
		// 标志信息
		UserCancel:       0,
		IsFileReady:      0,
		DownloadFinished: 0,
		IsSystemFailed:   0,
		IsDeleted:        0,
		// 时间信息
		CreateTime:       now,
		UpdateTime:       now,
		UploadTime:       InvalidTime,
		DownloadTime:     InvalidTime,
		PendingTime:      pendingTime,
		RunningTime:      InvalidTime,
		TerminatingTime:  InvalidTime,
		TransmittingTime: InvalidTime,
		SuspendingTime:   InvalidTime,
		SuspendedTime:    InvalidTime,
		SubmitTime:       InvalidTime,
		EndTime:          InvalidTime,
	}
	return job, nil
}

func FillScheduleInfo(ctx context.Context, req *jobcreate.Request, scheduleInfo *models.PreSchedule) error {
	logger := logging.GetLogger(ctx).With("func", "fillScheduleInfo", "preScheduleID", req.PreScheduleID)
	req.Params.Application = jobcreate.Application{
		Command: scheduleInfo.Command,
		AppID:   scheduleInfo.AppID.String(),
	}
	mincpus, _, mems := scheduleInfo.GetResource()
	req.Params.Resource = &jobcreate.Resource{
		Cores:  mincpus, // 将预调度期望的最小核数作为作业req的核数，真正调度时会根据资源池的资源情况进行调整
		Memory: mems,
	}
	envVars, err := scheduleInfo.GetEnvVars()
	if err != nil {
		logger.Warnf("get env vars error: %v", err)
		return fmt.Errorf("get env vars error: %w", err)
	}
	req.Params.EnvVars = envVars
	req.Params.Input = nil // 置空，不允许传入
	req.Zone = scheduleInfo.Zone
	req.NoRound = scheduleInfo.Shared // noRound means shared
	return nil
}

func ConvertPreScheduleModel(ctx context.Context, logger *zap.SugaredLogger,
	req *jobpreschedule.Request, userID, preScheduleID snowflake.ID,
	appInfo *models.Application, zone consts.Zone) (*models.PreSchedule, error) {
	workDir := fmt.Sprintf("%s/%s/%s", userID, consts.PreScheduleDir, preScheduleID)
	workDir = AddSuffixSlash(workDir)
	params, err := json.Marshal(req.Params)
	if err != nil {
		logger.Errorf("params json.Marshal error: %v", err)
		return nil, err // internal error
	}
	zones := config.GetConfig().Zones
	if zone != "" {
		zoneInfo, ok := zones[zone.String()]
		if ok && zoneInfo.HPCEndpoint != "" {
			workDir = fmt.Sprintf("%s/%s", zoneInfo.HPCEndpoint, workDir)
		}
	}
	command := req.Params.Application.Command
	appID := appInfo.ID
	appName := appInfo.Name
	mincores := int64(*req.Params.Resource.MinCores)
	maxcores := int64(*req.Params.Resource.MaxCores)
	memory := int64(*req.Params.Resource.Memory)
	// req.Params.EnvVars -> string
	envs, err := json.Marshal(req.Params.EnvVars)
	if err != nil {
		logger.Errorf("envs json.Marshal error: %v", err)
		return nil, err // internal error
	}
	preSchedule := &models.PreSchedule{
		ID:              preScheduleID,
		Params:          string(params),
		ExpectedMinCpus: mincores,
		ExpectedMaxCpus: maxcores,
		ExpectedMemory:  memory,
		Shared:          req.Shared,
		Fixed:           req.Fixed,
		Zone:            zone.String(),
		Command:         command,
		WorkDir:         workDir,
		AppID:           appID,
		AppName:         appName,
		Envs:            string(envs),
	}
	return preSchedule, nil
}

// ConvertAdminJobModel 参数转换成model
func ConvertAdminJobModel(ctx context.Context, logger *zap.SugaredLogger,
	req *adminjobcreate.Request, userID snowflake.ID, jobID snowflake.ID,
	appInfo *models.Application, inputZone, outputZone consts.Zone, queue string,
	scheduleInfo *models.PreSchedule) (*models.Job, error) {
	job, err := ConvertJobModel(ctx, logger, &req.Request, userID, jobID, appInfo,
		inputZone, outputZone, req.JobSchedulerSubmitFlags, scheduleInfo)
	if err != nil {
		return nil, err
	}
	job.Queue = queue
	return job, nil
}

// HpcModelToYsJobModel hpc job model to ys job model
func HpcModelToYsJobModel(job *schema.JobInHPC) *models.Job {
	// 纯粹的转换，不做任何逻辑处理
	if job == nil {
		return nil
	}
	input := schema.JobInHPCInputStorage{Type: "", Src: "", Dst: ""}
	if len(job.Inputs) > 0 {
		input = job.Inputs[0]
	}
	newState := consts.ConvertHpcStateToYsState(job.Status)
	return &models.Job{
		Queue:                 job.Queue,
		OriginJobID:           job.SchedulerID,
		ResourceUsageCpus:     int64(job.Resource.Cores),
		AllocType:             job.Resource.AllocType,
		InputType:             string(input.Type),
		InputDir:              input.Src,
		Destination:           input.Dst,
		OutputDir:             job.Output.Dst,
		OutputType:            string(job.Output.Type),
		State:                 newState.State,
		SubState:              newState.SubState,
		StateReason:           job.StateReason,
		RunningTime:           TimeParse(job.RunningTime),
		ResourceAssignCpus:    int64(job.AllocCores),
		ExitCode:              job.ExitCode,
		ExecutionDuration:     int64(job.ExecutionDuration),
		UploadFileSizeTotal:   int64(job.DownloadProgress.Total),   // hpcJob.DownloadProgress 才是上传相关
		UploadFileSizeCurrent: int64(job.DownloadProgress.Current), // hpcJob.DownloadProgress 才是上传相关
		Priority:              job.Priority,
		ExecHosts:             job.ExecHosts,
		ExecHostNum:           job.ExecHostsNum,
	}
}

// TimeParse time parse
func TimeParse(currentTime *time.Time) time.Time {
	if currentTime == nil {
		return InvalidTime
	}
	return *currentTime
}

// ToStringMap binPath to map,if binPath not exist return nil
func ToStringMap(binPath string) map[string]string {
	if len(binPath) == 0 {
		return nil
	}
	binPathMap := make(map[string]string)
	err := json.Unmarshal([]byte(binPath), &binPathMap)
	if err != nil {
		logging.Default().Errorf("unmarshal binPath error: %v", err)
		return nil
	}
	return binPathMap
}

// ResidualMarshal marshal residual to string
func ResidualMarshal(res *schema.Residual) (string, error) {
	data, err := json.Marshal(res)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

// ResidualUnmarshal unmarshal string to residual
func ResidualUnmarshal(s string) (*schema.Residual, error) {
	bs, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}
	content := &schema.Residual{}
	err = json.Unmarshal(bs, content)
	if err != nil {
		return nil, err
	}
	return content, nil
}

// MonitorChartMarshal marshal monitorChart to string
func MonitorChartMarshal(monitorChart []*schema.MonitorChart) (string, error) {
	data, err := json.Marshal(monitorChart)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

// MonitorChartUnmarshal unmarshal string to monitorChart
func MonitorChartUnmarshal(s string) ([]*schema.MonitorChart, error) {
	bs, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}
	content := make([]*schema.MonitorChart, 0)
	err = json.Unmarshal(bs, &content)
	if err != nil {
		return nil, err
	}
	return content, nil
}

func JobModelToNeedSyncFileJobs(jobs []*models.Job) []*jobneedsyncfile.NeedSyncFileJobInfo {
	if len(jobs) == 0 {
		return make([]*jobneedsyncfile.NeedSyncFileJobInfo, 0)
	}
	resp := make([]*jobneedsyncfile.NeedSyncFileJobInfo, 0, len(jobs))
	for _, job := range jobs {
		var stateDesc string
		state, b := consts.GetStateBySubState(job.SubState)
		if b {
			stateDesc = state.StateString()
		}
		needSyncFileJob := &jobneedsyncfile.NeedSyncFileJobInfo{
			ID:                      job.ID.String(),
			State:                   stateDesc,
			Name:                    job.Name,
			FileSyncState:           job.FileSyncState,
			WorkDir:                 job.WorkDir,
			OutputDir:               job.OutputDir,
			NoNeededPaths:           job.NoNeededPaths,
			NeededPaths:             job.NeededPaths,
			FileOutputStorageZone:   job.FileOutputStorageZone,
			DownloadFileSizeTotal:   job.DownloadFileSizeTotal,
			DownloadFileSizeCurrent: job.DownloadFileSizeCurrent,
		}
		resp = append(resp, needSyncFileJob)
	}
	return resp
}
