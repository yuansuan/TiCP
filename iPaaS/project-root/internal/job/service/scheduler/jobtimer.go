package scheduler

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/davecgh/go-spew/spew"
	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/common/project-root-api/hpc/jobstate"
	"github.com/yuansuan/ticp/common/project-root-api/proto/license"
	schema "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/config"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/consts"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao/models"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/module"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/module/hpc/openapi"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/module/rpc"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/module/wx"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/service/alarm"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/service/v1/application"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/util"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/common/with"
	"go.uber.org/zap"
	"xorm.io/xorm"
)

const (
	warnPendingCount = 20
)

// JobTimer 作业更新定时器
type JobTimer struct {
	JobDao         dao.JobDao
	AppSrv         application.Service
	Sender         wx.Sender
	LongRunningJob sync.Map
}

// NewJobTimer 创建作业更新定时器
func NewJobTimer(jobDao dao.JobDao, appSrv application.Service, sender wx.Sender) *JobTimer {
	return &JobTimer{
		JobDao: jobDao,
		AppSrv: appSrv,
		Sender: sender,
	}
}

// Run 定时器执行函数
func (s *JobTimer) Run(ctx context.Context) {
	s.syncJobInfo(ctx)
}

// syncJobInfo 行为：定时查询hpc作业信息，更新到数据库
func (s *JobTimer) syncJobInfo(ctx context.Context) {
	logger := logging.GetLogger(ctx).With("func", "syncJobInfo")
	logger.Info("job timer sync hpc job is running.......")

	session := boot.MW.DefaultSession(ctx)
	defer session.Close()
	ctx = with.KeepSession(ctx, session)
	jobs := make([]*models.Job, 0)
	_, jobsForUpdate, err := s.getJobsForUpdate(ctx)
	if err != nil {
		logger.Warnf("query job error! err: %v", err)
		return
	}

	// 顺便统计Pending的作业数量
	pendingCount, runningCount, terminatingCount := 0, 0, 0
	for _, job := range jobsForUpdate {
		jobs = append(jobs, job)

		if job.State == consts.Pending {
			pendingCount++
		}
		if job.State == consts.Running {
			runningCount++
		}
		if job.State == consts.Terminating {
			terminatingCount++
		}
	}

	if pendingCount > warnPendingCount {
		logger.Warnf("pending job count: %d, running %d, terminating %d",
			pendingCount, runningCount, terminatingCount)
	}
	if len(jobs) == 0 {
		logger.Info("no jobs sync!")
		return
	}

	zones := config.GetConfig().Zones
	updated := s.processJobs(ctx, jobs, zones)
	logger.Infof("syncJobInfo job timer sync hpc job end, update count: %d", updated)
}

func (s *JobTimer) getJobsForUpdate(ctx context.Context) (int64, []*models.Job, error) {
	var querySubStates = []int{
		consts.SubStateFileUploading.SubState,
		consts.SubStateHpcWaiting.SubState,
		consts.SubStateRunning.SubState,
		consts.SubStateTerminating.SubState,
		consts.SubStateSuspending.SubState,
		consts.SubStateUnknown.SubState,
	}
	return s.JobDao.ListJobsBySubStates(ctx, querySubStates...)
}

func (s *JobTimer) processJobs(ctx context.Context, jobs []*models.Job, zones schema.Zones) int {
	updated := 0
	for _, job := range jobs {
		if err := s.processJob(ctx, job, zones); err == nil {
			updated++
		}
	}
	return updated
}

func (s *JobTimer) processJob(ctx context.Context, job *models.Job, zones schema.Zones) error {
	logger := logging.GetLogger(ctx).With("func", "processJob", "jobId", job.ID)
	hpcJobID := snowflake.MustParseString(job.HPCJobID).String()
	if hpcJobID == "" || hpcJobID == "0" {
		logger.Warnf("get hpc job error! originalJobId error! originalJobId: %v", hpcJobID)
		return fmt.Errorf("get hpc job error! originalJobId error! originalJobId: %v", hpcJobID)
	}
	logger = logger.With("hpcJobId", hpcJobID)

	zone, ok := zones[job.Zone]
	if !ok {
		logger.Warnf("get hpc job error! zone error! zone: %v", job.Zone)
		return fmt.Errorf("get hpc job error! zone error! zone: %v", job.Zone)
	}
	zoneDomain := zone.HPCEndpoint
	if zoneDomain == "" {
		logger.Warnf("get hpc job error! zoneDomain error! zoneDomain: %v", zoneDomain)
		return fmt.Errorf("get hpc job error! zoneDomain error! zoneDomain: %v", zoneDomain)
	}
	jobResp, err := openapi.Client().GetJob(zoneDomain, module.DefaultTimeout, hpcJobID)
	if err != nil {
		logger.Warnf("get hpc job error! original job id: %v,  err: %v", hpcJobID, err)
		return err
	}
	logger.Infof("hpc sync job id: %v result: %v", hpcJobID, jobResp)
	if jobResp == nil || jobResp.Data == nil {
		logger.Warnf("get hpc job error! original job id: %v", hpcJobID)
		return fmt.Errorf("get hpc job error! original job id: %v", hpcJobID)
	}

	hpcJob := jobResp.Data
	logger.Infof("hpc get job info : %s", spew.Sdump(hpcJob))

	updater := NewJobUpdater(logger)
	jobModel := updater.Init(job, hpcJob)
	if jobModel == nil {
		logger.Warnf("get hpc job error! original job id: %v", hpcJobID)
		return fmt.Errorf("get hpc job error! original job id: %v", hpcJobID)
	}
	for _, fn := range updater.allFunc() {
		jobModel = fn(jobModel, hpcJob)
	}

	jobID := snowflake.MustParseString(job.ID.String())
	err = s.JobDao.Transaction(ctx, func(c context.Context) error {
		var dbJob models.Job
		err := with.DefaultSession(c, func(db *xorm.Session) error {
			_, err := db.ID(jobID).ForUpdate().Get(&dbJob)
			return err
		})
		if err != nil {
			logger.Warnf("Get job info error! job id: %v, err: %v", jobID, err)
			return err
		}

		var dbApp models.Application
		err = with.DefaultSession(c, func(db *xorm.Session) error {
			_, err := db.ID(dbJob.AppID).Get(&dbApp)
			return err
		})
		if err != nil {
			logger.Warnf("Get app info error! app id: %v, err: %v", dbJob.AppID, err)
			return err
		}

		// 防止更新期间，作业被终止
		if dbJob.State != updater.OldJob.State {
			logger.Infof("Overdue operation status, oldState: %v, dbState: %v,  Job id: %v",
				updater.OldJob.State, dbJob.State, jobID)
			return nil // 无需打印告警日志
		}

		if updater.OldState.ToFinal(updater.NewState) {
			// 释放license
			if config.GetConfig().ChangeLicense && dbApp.LicManagerId > 0 {
				err = releaseLicense(ctx, jobID)
				if err != nil {
					// TODO: important, 释放license失败，应该报警
					logger.Warnf("release license error!, job id: %v,  err: %v", jobID, err)
					return err
				}
			}
		}

		jobModel.UpdateTime = time.Now()
		err = with.DefaultSession(c, func(db *xorm.Session) error {
			// 文件下载大小以sync-runner 上传的数据库数据为准，不再以hpc的下载数据为准,因此使用原值
			_, err = db.ID(jobID).Omit("download_file_size_total", "download_file_size_current",
				"file_sync_state", "download_finished", "download_time").Update(jobModel)
			return err
		})
		if err != nil {
			logger.Warnf("update job info error! job id: %v, err: %v", jobID, err)
			return err
		}
		logger.Info("update job info success!")
		return nil
	})
	if err != nil {
		logger.Warnf("update job info error!, job id: %v,  err: %v", jobID, err)
		return err
	}

	go s.checkLongRunningJob(ctx, jobModel) // 检查长时间运行的作业可异步
	return nil
}

func (s *JobTimer) checkLongRunningJob(ctx context.Context, job *models.Job) {
	threshold := config.GetConfig().LongRunningJobThreshold
	if threshold > 0 {
		// 如果是长时间运行的作业，需要记录下来并告警
		if job.State == consts.Running && job.IsLongRunning(threshold) {
			if _, ok := s.LongRunningJob.Load(job.ID); !ok {
				s.LongRunningJob.Store(job.ID, job)
				// 发送告警信息到企微webhook

				alarm.SendLongRunningJobAlarm(ctx, s.Sender, job, threshold)
			}
		} else {
			// 如果作业不是Running状态，就从长时间运行的作业列表中移除
			s.LongRunningJob.Delete(job.ID)
		}
	}
}

// releaseLicense 释放license
func releaseLicense(ctx context.Context, jobID snowflake.ID) error {
	logger := logging.GetLogger(ctx).With("func", "releaseLicense", "job_id", jobID)
	req := &license.ReleaseRequest{
		JobId: jobID.Int64(),
	}
	_, err := rpc.GetInstance().License.LicenseServer.ReleaseLicense(ctx, req)
	if err != nil {
		logger.Warnf("job %s release License networks: %v", jobID.String(), err)
		return err
	}
	return nil
}

// JobUpdater 更新
type JobUpdater struct {
	now      time.Time
	OldState consts.State
	NewState consts.State
	OldJob   *models.Job

	logger *zap.SugaredLogger
}

type jobUpdateFunc func(job *models.Job, hpcJob *schema.JobInHPC) *models.Job

// NewJobUpdater 新建
func NewJobUpdater(logger *zap.SugaredLogger) *JobUpdater {
	return &JobUpdater{
		logger: logger,
	}
}

func (u *JobUpdater) allFunc() []jobUpdateFunc {
	funcs := []jobUpdateFunc{
		u.updateUpload,
		u.updateState,
		u.updateSystemFailed,
		u.updateExitCode,
		u.updateEndtime,
	}
	return funcs
}

// Init 初始化
func (u *JobUpdater) Init(job *models.Job, hpcJob *schema.JobInHPC) *models.Job {
	jobModel := util.HpcModelToYsJobModel(hpcJob)
	if jobModel == nil {
		return nil
	}

	now := time.Now()
	jobModel.ID = job.ID
	jobModel.HPCJobID = job.HPCJobID
	jobModel.Zone = job.Zone
	oldState := consts.NewState(job.State, job.SubState)
	newState := consts.NewState(jobModel.State, jobModel.SubState)
	jobModel.StateReason = consts.ParseAndUpdateStateReasonString(job.StateReason,
		newState, jobModel.StateReason).String()

	u.now = now
	u.OldJob = job
	u.OldState = oldState
	u.NewState = newState
	return jobModel
}

func (u *JobUpdater) updateSystemFailed(job *models.Job, hpcJob *schema.JobInHPC) *models.Job {
	// 根据hpc的状态以及IsUserFailed字段更新作业失败状态
	if hpcJob.Status == jobstate.Failed && !hpcJob.IsUserFailed {
		if job.IsSystemFailed == 0 {
			u.logger.Errorf("job system failed! Admin please follow in time!\n"+
				"hpc info: SchedulerID:[%s], Queue:[%s], StateReason :[%s], ExitCode:[%s], ExecHosts:[%s]",
				hpcJob.SchedulerID, hpcJob.Queue, hpcJob.StateReason, hpcJob.ExitCode, hpcJob.ExecHosts)
		}
		job.IsSystemFailed = 1
	}
	return job
}

func (u *JobUpdater) updateExitCode(job *models.Job, hpcJob *schema.JobInHPC) *models.Job {
	// 完成状态根据程序退出码修正
	if u.NewState.IsCompleted() {
		nonZero, err := util.IsNonZeroExitCode(hpcJob.ExitCode)
		if err != nil {
			u.logger.Errorf("parse exit code error: %v", err)
			return job
		}
		if nonZero {
			job.State = consts.SubStateFailed.State
			job.SubState = consts.SubStateFailed.SubState
			job.StateReason = consts.ParseAndUpdateStateReasonString(u.OldJob.StateReason,
				consts.SubStateFailed, consts.StateReasonNonZeroExitCode).String()
			u.NewState = consts.SubStateFailed
		}
	}
	return job
}

func (u *JobUpdater) updateEndtime(job *models.Job, hpcJob *schema.JobInHPC) *models.Job {
	// 其他状态转到终态，就修改endTime=now
	if u.OldState.ToFinal(u.NewState) {
		job.EndTime = u.now
	}

	return job
}

func (u *JobUpdater) updateUpload(job *models.Job, hpcJob *schema.JobInHPC) *models.Job {
	// 作业上传状态和上传时间
	if u.OldJob.IsFileReady == 0 && job.UploadFileSizeCurrent >= job.UploadFileSizeTotal {
		job.IsFileReady = 1
		job.UploadTime = u.now
	}
	return job
}

func (u *JobUpdater) updateState(job *models.Job, hpcJob *schema.JobInHPC) *models.Job {
	// 作业状态
	if !stateNeedUpdate(consts.NewState(u.OldJob.State, u.OldState.SubState), hpcJob.Status) {
		job.State = u.OldJob.State
		job.SubState = u.OldJob.SubState
		job.StateReason = u.OldJob.StateReason
		u.NewState = consts.NewState(job.State, job.SubState)
	}

	// 我们认为如果作业状态是终止中，那么无论hpc上是何种终态，都只能转到已终止
	if consts.NewState(u.OldJob.State, u.OldJob.SubState) == consts.SubStateTerminating &&
		consts.ConvertHpcStateToYsState(hpcJob.Status).IsFinal() {
		job.State = consts.SubStateTerminated.State
		job.SubState = consts.SubStateTerminated.SubState
		job.StateReason = consts.ParseAndUpdateStateReasonString(u.OldJob.StateReason,
			consts.SubStateTerminated, consts.StateReasonUserCancel).String()
		u.NewState = consts.SubStateTerminated
	}
	return job
}

func stateNeedUpdate(oldState consts.State, hpcState jobstate.State) bool {
	if oldState == consts.SubStateTerminating && hpcState != jobstate.Canceled {
		return false
	}
	return true
}
