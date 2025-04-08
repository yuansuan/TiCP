package scheduler

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/common/project-root-api/proto/license"
	v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
	"go.uber.org/zap"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/config"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/consts"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao/models"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/module"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/module/hpc/openapi"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/module/rpc"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/service/v1/application"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/util"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/common/with"
)

// LicenseInfo 配置的license消息
type LicenseInfo struct {
	JobID       int64
	ServerURL   string
	PlaceHolder string
}

// Scheduler 调度器
type Scheduler struct {
	logger *zap.SugaredLogger
	AppSrv application.Service
	JobDao dao.JobDao
}

// NewScheduler 新建调度器
func NewScheduler(logger *zap.SugaredLogger, appSrv application.Service, jobDao dao.JobDao) *Scheduler {
	return &Scheduler{logger: logger, AppSrv: appSrv, JobDao: jobDao}
}

// Run 运行
func (s *Scheduler) Run(ctx context.Context) {
	s.run(ctx)
}

func (s *Scheduler) run(ctx context.Context) {
	logger := s.logger.With("func", "job.scheduler.run")
	logger.Infof("scheduler.run() start...")
	defer logger.Infof("scheduler.run() end...")

	session := boot.MW.DefaultSession(ctx)
	defer session.Close()
	ctx = with.KeepSession(ctx, session)
	// 获取zones
	zones := config.GetConfig().Zones //TODO: 可添加到Scheduler中
	// 作业调度
	s.schedulerJobs(ctx, logger, zones)
}

// schedulerJobs 调度作业
// 处理状态为Scheduling,Suspending,Terminating的作业
func (s *Scheduler) schedulerJobs(ctx context.Context, logger *zap.SugaredLogger, zones v20230530.Zones) {
	logger = logger.With("func", "job.scheduler.schedulerJobs")
	// 获取paas数据库[SubStateScheduling,SubStateSuspending,SubStateTerminating]的job
	var querySubStates = []int{consts.SubStateScheduling.SubState,
		consts.SubStateSuspending.SubState, consts.SubStateTerminating.SubState,
	}
	count, jobs, err := s.JobDao.ListJobsBySubStates(ctx, querySubStates...)
	if err != nil {
		logger.Warnf("query job error! err: %v", err)
		return
	}
	if count == 0 || len(jobs) == 0 {
		logger.Infof("no job need to scheduler")
		return
	}
	for _, job := range jobs {
		state := consts.NewState(job.State, job.SubState)
		switch {
		case state.IsScheduling():
			s.jobsubmit(ctx, zones, job)
		case state.IsTerminating():
			s.jobterminate(ctx, zones, job)
		default:
			continue
		}
	}
}

// jobsubmit 作业提交
func (s *Scheduler) jobsubmit(ctx context.Context, zones v20230530.Zones, job *models.Job) {
	logger := s.logger.With("func", "job.scheduler.jobsubmit", "job_id", job.ID,
		"job_name", job.Name, "user_id", job.UserID, "app_id", job.AppID,
		"app_name", job.AppName, "user_zone", job.UserZone)
	logger.Infof("jobsubmit start...")
	defer logger.Infof("jobsubmit end...")

	// 	参数转换
	params := &models.AdminParams{}
	err := json.Unmarshal([]byte(job.Params), params)
	if err != nil {
		// update job stateReason
		logger.Warnf("unmarshal job params error! err: %v", err)
		s.setJobScheduling(ctx, job, "unmarshal job params error")
		return
	}

	// 获取应用信息
	appInfo, err := s.AppSrv.Apps().GetApp(ctx, job.AppID)
	if err != nil {
		// update job stateReason
		logger.Warnf("get app error! err: %v", err)
		s.setJobScheduling(ctx, job, "get app error")
		return
	}

	// job.NoRound means shared
	// shared或者app.NeedLimitCore为true时，不取整
	shared := job.NoRound || appInfo.NeedLimitCore
	hpcClient := openapi.Client()

	appSpecifyQueueMap := util.ToStringMap(appInfo.SpecifyQueue)

	// 区域选择器
	zs := NewZoneSelector(zones, hpcClient, appSpecifyQueueMap, job.UserZone, job.Queue)

	// 资源聚合
	resources, err := zs.ResourceAggregation(ctx)
	if err != nil {
		// update job stateReason
		logger.Infof("resource aggregation error! err: %v", err)
		s.setJobScheduling(ctx, job, errors.Wrap(err, "resource aggregation error").Error())
		return
	}

	inputHPCZone := ""
	if job.InputType == string(consts.HpcStorage) {
		inputHPCZone = job.FileInputStorageZone
	}

	cpuRange := util.CoresRange{
		MinExpectedCores: job.ResourceUsageCpus,
		MaxExpectedCores: job.ResourceUsageCpus,
	}

	// 如果有预调度，应当去取预调度的Min,Max
	if job.PreScheduleID != "" {
		scheduleInfo, exist, err := s.JobDao.GetPreSchedule(ctx, snowflake.MustParseString(job.PreScheduleID))
		if err != nil {
			logger.Infof("get preSchedule info error: %v", err)
			s.setJobScheduling(ctx, job, "get preSchedule info error")
			return
		}
		if !exist {
			logger.Infof("preSchedule info not exist")
			s.setJobScheduling(ctx, job, "preSchedule info not exist")
			return
		}

		cpuRange.MinExpectedCores = scheduleInfo.ExpectedMinCpus
		cpuRange.MaxExpectedCores = scheduleInfo.ExpectedMaxCpus
	}

	licenseFilterParams := &LicenseFilterParams{
		IdentifierID:         job.ID,
		CPURange:             cpuRange,
		Shared:               job.NoRound, // job.NoRound means shared
		Average:              job.AllocType == "average",
		JobResourceUsageCpus: job.ResourceUsageCpus, //average模式下用这个字段，用户请求的核数
	}
	zs.RegisterFilter(
		NewBaseFilter(cpuRange, job.ResourceUsageMemory, job.UserZone, job.Queue, inputHPCZone, shared),
		NewReserveResourceFilter(cpuRange, job.Queue, appSpecifyQueueMap),
		NewLicenseFilter(zones, licenseFilterParams, appInfo, rpc.GetInstance().License.LicenseServer),
	)

	// 优选器
	selectorWeights := config.GetConfig().SelectorWeights
	rc := NewResourceCount()
	rc.SetWeight(selectorWeights[rc.Name()])
	qps := NewQueuePrioritySelector() // 队列优先级选择器
	qps.SetWeight(selectorWeights[qps.Name()])
	sf := NewStorageFirst(job.FileInputStorageZone)
	sf.SetWeight(selectorWeights[sf.Name()])
	zs.RegisterOptimalSelector(rc, qps, sf)

	selectParams := &SelectParams{
		Resources: resources,
	}
	zone, resource, err := zs.Select(ctx, selectParams, WithJobID(job.ID), WithUserID(job.UserID))
	if err != nil {
		// update job stateReason
		logger.Infof("select zone failed! err: %v", err)
		s.setJobScheduling(ctx, job, errors.Wrap(err, "select zone failed").Error())
		return
	}
	logger.Infof("zone: %s, selectZone: %s, resource: %+v", job.UserZone, zone, resource)

	job.Queue = resource.Queue

	// input判断
	// 	1. job.InputType为hpc_storage时，不进行实际的数据传输
	// 作业预调度时，input为空，也不进行传输
	noTransfer := true
	if job.PreScheduleID == "" {
		if params.Input == nil {
			// 不应该出现，检查代码错误
			logger.Errorf("job input is nil when job.PreScheduleID is empty, "+
				"should not be nil, check code error!, job: %+v", job)
			return
		}
		noTransfer = s.checkNoTransfer(job.InputType, params.Input.Destination)
	}

	//  取hpc域名
	relzone, ok := zones[zone]
	if !ok {
		logger.Warnf("get zone config failed! zone: %s", zone)
		s.setJobScheduling(ctx, job, fmt.Sprintf("get zone config failed! zone: [%s]", zone))
		return
	}
	zoneDomain := relzone.HPCEndpoint
	if zoneDomain == "" {
		logger.Warnf("get zone domain failed! zone: %s", zone)
		s.setJobScheduling(ctx, job, fmt.Sprintf("get zone domain failed! zone: [%s]", zone))
		return
	}
	if job.UserZone == "" {
		job.WorkDir = fmt.Sprintf("%s/%s", zoneDomain, job.WorkDir)
	}

	// hpcresource 是从 zone resource里面获取数据的
	hpcResource := &v20230530.Resource{
		Cpu:          resource.CPU, //空闲CPU
		Memory:       resource.Mem,
		CoresPerNode: resource.CoresPerNode, //物理节点的核数
		TotalNodeNum: resource.TotalNodeNum, //当前队列的节点数
	}

	if hpcResource.CoresPerNode == 0 {
		logger.Errorf("hpc resource coresPerNode is 0!")
		s.setJobScheduling(ctx, job, "hpc resource coresPerNode is 0!")
		return
	}

	// coresPerNode 按理来说是机器单节点核数，但是在 `非average`且`noround= true` 的条件下,
	// 	它会变成用户请求的核数，即 <= 单节点核数
	var totalCores int64
	var coresPerNode int64
	if job.AllocType == "average" {
		// average类型直接使用机器的核数配置; 实际上 这个值无论多少在hpc_sc里面都不影响提交作业
		coresPerNode = hpcResource.CoresPerNode
		if shared {
			totalCores = job.ResourceUsageCpus
		} else { // 非shared用户至少要用一个单机节点核数
			totalCores = max(coresPerNode, job.ResourceUsageCpus)
		}
	} else {
		var resourceUsageCpus int64
		coresPerNode, resourceUsageCpus, err = util.CalculateResourceUsage(cpuRange,
			hpcResource.CoresPerNode, hpcResource.Cpu, shared)
		if err != nil {
			logger.Warnf("calculate resource usage error! err: %v", err)
			s.setJobScheduling(ctx, job, "calculate resource usage error")
			return
		}
		totalCores = resourceUsageCpus
	}

	// job.ResourceUsageCpus = resourceUsageCpus 原来的逻辑 ，怎么手动改用户提交的参数？TODO 可能的bug
	// 应该改 实际分配的核数吧 job.ResourceAssignCpus？But the later check uses this value
	job.ResourceUsageCpus = totalCores

	// job.ResourceUsageMemory
	//  资源不足应当排队
	resourceOK, err := s.checkResource(ctx, job, hpcResource)
	if err != nil {
		// update job stateReason
		logger.Warnf("check resource error! err: %v", err)
		s.setJobScheduling(ctx, job, "check resource error")
		return
	}
	if !resourceOK {
		// update job stateReason
		logger.Info("resource not enough!")
		s.setJobScheduling(ctx, job, "resource not enough")
		return
	}

	appImage := appInfo.Image
	// 本地镜像
	localImage := false
	ysAppBin := ""

	if appInfo.BinPath != "" {
		m := util.ToStringMap(appInfo.BinPath)
		p, ok := m[zone]
		if !ok {
			logging.Default().Debugf("get bin path failed! zone: %s", zone)
		} else {
			localImage = true
			appImage = p
			ysAppBin = p
		}
	}

	envVars := params.EnvVars
	if envVars == nil {
		envVars = make(map[string]string)
	}
	if job.Command == "" {
		// 非命令行应用取应用的默认命令
		job.Command = appInfo.Command

		// 忽略appinfo里不存在的参数
		extentionParams := make(map[string]v20230530.ExtentionParam)
		err = json.Unmarshal([]byte(appInfo.ExtentionParams), &extentionParams)
		if err != nil {
			// update job stateReason
			logger.Warnf("unmarshal extentionParams error! err: %v", err)
			s.setJobScheduling(ctx, job, "unmarshal extentionParams error")
			return
		}

		validEnv := make(map[string]string, len(envVars))
		for k, v := range envVars {
			// 检查key是否存在
			if _, ok := extentionParams[k]; !ok {
				continue // 忽略不存在的参数
			}
			validEnv[k] = v
		}
		envVars = validEnv
	} else {
		// 添加上 AppPreparedFlag
		// 这样用户命令的执行均在AppPreparedFlag之后，失败即都为用户失败
		job.Command = fmt.Sprintf("%s\n%s", consts.AppPreparedFlag, job.Command)
	}

	// 添加YS_APP_BIN环境变量
	envVars["YS_APP_BIN"] = ysAppBin // 镜像应用 或 分区应用未配置时 为空
	// 一切准备就绪后，进行license check
	if config.GetConfig().ChangeLicense && appInfo.LicManagerId > 0 {
		// appInfo.LicManagerId > 0 代表商业软件
		// 	license不足应当排队等待
		licenseOK, licEnvs, err := s.checkLicense(ctx, job, zoneDomain, appInfo.LicManagerId)
		if err != nil {
			// update job stateReason
			logger.Warnf("check license error! err: %v", err)
			s.setJobScheduling(ctx, job, "license check error")
			return
		}
		if !licenseOK {
			// update job stateReason
			logger.Info("license not enough!")
			s.setJobScheduling(ctx, job, "license not enough")
			return
		}
		// 添加license 环境变量
		for _, env := range licEnvs {
			tmp := strings.SplitN(env, "=", 2)
			if len(tmp) != 2 {
				logger.Errorf("bad license envs: %s", env)
				s.setJobScheduling(ctx, job, "license not enough")
				return
			}
			envVars[tmp[0]] = tmp[1]
		}
	}

	schedulerSubmitFlags := params.JobSchedulerSubmitFlags
	if _, ok := envVars["YS_MAIN_FILE"]; ok {
		envVars["YS_MAIN_FILE"] = fmt.Sprintf(`"%s"`, envVars["YS_MAIN_FILE"])
	}
	hpcJobReq := util.AssembleHPCJobRequest(ctx, logger, job, appImage, envVars,
		schedulerSubmitFlags, noTransfer, localImage, int(coresPerNode))
	logger.Infof("hpc submit request: %+v", hpcJobReq)

	hpcJobResp, err := hpcClient.PostJob(zoneDomain, module.DefaultTimeout, hpcJobReq)
	if err != nil {
		logger.Errorf("hpc job create api error: %v", err)
		s.setJobScheduling(ctx, job, "hpc job create api error")
		// 如果提交失败，释放license
		if config.GetConfig().ChangeLicense && appInfo.LicManagerId > 0 {
			releaseLicense(ctx, job.ID)
		}
		return
	}
	hpcJob := hpcJobResp.Data
	logger.Infof("hpc submit response: %+v", hpcJob)

	newState := consts.SubStateFileUploading
	job.State = newState.State
	job.SubState = newState.SubState
	nsr := consts.ParseAndUpdateStateReasonString(job.StateReason,
		consts.SubStateFileUploading, "hpc job create success")
	job.StateReason = nsr.String()

	// 修改作业运行信息
	now := time.Now()
	hpcJobID := hpcJob.ID
	job.HPCJobID = hpcJobID
	job.SubmitTime = now
	job.UpdateTime = now
	job.Zone = zone

	//  Update DB
	//  应该修改的字段: hpc_job_id, submit_time, update_time, zone, resource_usage_cpus,
	// 	resource_usage_memory, state_reason, command, state, sub_state, work_dir, alloc_type
	updateResult, err := s.JobDao.UpdateSubmitJob(ctx, job)
	if err != nil {
		logger.Errorf("update job info after submit error! err: %v", err)
		s.setJobScheduling(ctx, job, "update job info after submit error")
		return
	}
	logger.With("updateResult", updateResult).Infof("update job info complete")
}

// jobterminate 作业终止
func (s *Scheduler) jobterminate(ctx context.Context, zones v20230530.Zones, job *models.Job) {
	logger := s.logger.With("func", "job.scheduler.jobterminate", "job_id", job.ID,
		"hpc_job_id", job.HPCJobID, "zone", job.Zone, "user_id", job.UserID)
	logger.Infof("jobterminate start...")

	// 取hpc域名
	zone, ok := zones[job.Zone]
	if !ok {
		logger.Warnf("get zone config failed! zone: %s", job.Zone)
		return
	}

	zoneDomain := zone.HPCEndpoint

	if zoneDomain == "" {
		logger.Warnf("get zone domain failed! zone: %s", job.Zone)
		return
	}

	hpcClient := openapi.Client()

	resp, err := hpcClient.CancelJob(zoneDomain, module.DefaultTimeout, job.HPCJobID)
	if err != nil {
		logger.Warnf("hpc job terminate api error: %v", err)
		return
	}

	logger.Infof("hpc job terminate response: %+v", resp)

	logger.Infof("jobterminate end...")
}

// checkNoTransfer 检查是否不进行实际的数据传输
func (s *Scheduler) checkNoTransfer(inputType string, inputDest string) bool {
	// job.InputType为hpc_storage且inpuDest==""时，不进行实际的数据传输
	return inputType == consts.HpcStorage.String() && inputDest == ""
}

// check license
func (s *Scheduler) checkLicense(ctx context.Context, job *models.Job,
	endpoint string, licManagerId snowflake.ID) (bool, []string, error) {
	logger := s.logger.With("func", "job.scheduler.checkLicense",
		"job_id", job.ID, "hpc_job_id", job.HPCJobID, "zone", job.Zone,
		"user_id", job.UserID, "app_id", job.AppID)

	logger.Info("checkLicense start...")
	logger.Infof("cpu amount is %d", job.ResourceUsageCpus)

	consumeInfo := &license.ConsumeInfo{
		JobId:        job.ID.Int64(),
		AppId:        job.AppID.Int64(),
		Cpus:         job.ResourceUsageCpus,
		LicManagerId: licManagerId.Int64(),
		HpcEndpoint:  endpoint,
	}

	if !licManagerId.NotZero() {
		logger.Info("licManager not set, means no need to check license")
		return true, nil, nil
	}

	// 请求license的RPC
	var consumeInfos []*license.ConsumeInfo
	consumeInfos = append(consumeInfos, consumeInfo)
	req := &license.ConsumeRequest{
		Info: consumeInfos,
	}

	logger.Infof("request license server: %+v", req)
	licenseServer, err := rpc.GetInstance().License.LicenseServer.AcquireLicenses(ctx, req)
	if err != nil {
		logger.Warnf("request license server,acquire License networks: error: %v", err)
		return false, nil, err
	}
	results := licenseServer.Result
	logger.Infof("response license server: %+v", results)

	// 单个查询，存在且只有1个
	result := results[0]
	licenseStatus := result.Status
	// license未配置
	if licenseStatus == license.LicenseStatus_UNCONFIGURED {
		logger.Warnf("license not configured, LicManagerId: [%v]", licManagerId.String())
		return true, nil, nil //! 部分应用未配置
	}
	// license不够
	if licenseStatus == license.LicenseStatus_NOTENOUTH ||
		licenseStatus == license.LicenseStatus_UNPUBLISH {
		logger.Warnf("lack of license, state: [%s], if long-term alarm, "+
			"check whether the total number of licenses is insufficient or "+
			"the partition has no license configuration", licenseStatus.String())
		return false, nil, nil
	}
	// license足够
	return true, result.LicenseEnvs, nil
}

// check resource
func (s *Scheduler) checkResource(ctx context.Context, job *models.Job,
	hpcResource *v20230530.Resource) (bool, error) {
	if hpcResource == nil {
		return false, errors.New("hpc resource is nil")
	}
	if job.ResourceUsageCpus > hpcResource.Cpu {
		return false, nil
	}
	if job.ResourceUsageMemory > hpcResource.Memory {
		return false, nil
	}
	return true, nil
}

func (s *Scheduler) setJobScheduling(ctx context.Context, job *models.Job, message string) error {
	nsr := consts.ParseAndUpdateStateReasonString(job.StateReason, consts.SubStateScheduling, message)
	job.StateReason = nsr.String()
	job.UpdateTime = time.Now()

	err := s.JobDao.UpdateSchedulingReason(ctx, job)
	if err != nil {
		return err
	}
	return nil
}
