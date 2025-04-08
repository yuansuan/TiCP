package statemachine

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	registry "github.com/yuansuan/ticp/iPaaS/standard-compute/tools/singularity/registry"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/tools/singularity/registry/image"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/yuansuan/ticp/common/go-kit/logging"
	"xorm.io/xorm"

	"github.com/yuansuan/ticp/common/project-root-api/hpc/jobstate"

	"github.com/yuansuan/ticp/iPaaS/standard-compute/config"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/backend"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/backend/job"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/dao"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/dao/models"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/log"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/oshelp"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/statemachine/jobsubstate"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/storage"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/storage/client"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/util"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/pkg/with"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/pkg/xhttp"
)

const (
	// _DefaultPushJobRetryTimes 推送任务到远算的默认重试次数
	_DefaultPushJobRetryTimes = 3

	// _DefaultWatcherInterval 默认的作业状态间隔时长
	_DefaultWatcherInterval = 3 * time.Second

	// _DefaultKillJobRetryTimes 强制取消作业的默认重试次数
	_DefaultKillJobRetryTimes = 10

	// _DefaultMarkJobCompletedRetryTimes 标记任务已完成的默认重试次数
	_DefaultMarkJobCompletedRetryTimes = 20
)

// StateMachine 任务状态机
type StateMachine struct {
	cfg           *config.Config
	storage       *storage.Manager
	backend       backend.Provider
	singularity   registry.Client
	dao           *dao.Dao
	db            *xorm.Engine
	storageClient *client.Client

	watcherInterval time.Duration
}

// Start 启动状态机
func (m *StateMachine) Start(ctx context.Context, j *job.Job) {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
		}
	}()

	ctx = log.WithTraceLoggerAndJobId(ctx, logging.Default(), j.Id)
	j.TraceLogger = log.GetJobTraceLogger(ctx)
	for {
		if j.State == jobstate.Completed || j.State == jobstate.Canceled {
			break
		}

		action := m.getAction(j.State)
		if err := action(ctx, j); err != nil {
			// 由于用户取消作业或任务超时导致的任务停止
			if isJobCanceled(err) {
				// 消息作业的处理函数应该是不可中断的
				m.OnCanceled(m.BackgroundContext(), j)
				break
			}

			// 其他非取消导致的错误
			m.OnFailed(ctx, j, err)
			break
		}
	}

	m.Completed(ctx, j)
}

// Preparing 作业准备阶段
// 1. 下载镜像
// 2. 下载输入文件至工作目录
// 3. 提交计算至计算集群
func (m *StateMachine) Preparing(ctx context.Context, j *job.Job) error {
	j.TraceLogger.Infof("job preparing: %v", j.Id)

	err := m.determineWorkspace(j)
	if err != nil {
		err = fmt.Errorf("determine workspace failed, %w", err)
		j.TraceLogger.Error(err)
		return errWrap(err, "determine workspace failed")
	}

	// singularity模式需要下载镜像，非singularity模式不需要下载镜像
	switch j.AppMode {
	case models.ImageAppMode:
		// download singularity image
		locator, err := image.FromString(j.SingularityImage)
		if err != nil {
			return errWrap(err, "fsm:singularity:locator")
		}

		// 开始下载 Singularity 镜像
		j.SubState = jobsubstate.PreparingPulling
		j.StateReason = fmt.Sprintf("pulling image: %s", locator)
		if err = m.UpdateAndPushJob(ctx, j); err != nil {
			return errWrap(err, "fsm:preparing")
		}

		locally, err := m.singularity.Pull(ctx, locator, registry.WithPullBlocking())
		if err != nil {
			return errWrap(err, "fsm:singularity:pull")
		}

		j.AppPath, err = locally.RealPath()
		if err != nil {
			return errWrap(err, "fsm:singularity:realpath")
		}

	case models.LocalAppMode:
		// check if local app exist
		if j.AppPath == "" {
			return errWrap(errors.New("app path is empty"), "fsm:check local app")
		}

		if !util.IsFileExist(j.AppPath) {
			return errWrap(fmt.Errorf("app not exist which path is [%s]", j.AppPath), "fsm:check local app")
		}
	default:
		return errWrap(fmt.Errorf("unsupported app mode %s", j.AppMode), "fsm:preparing application")
	}

	// 开始下载作业输入文件
	j.SubState = jobsubstate.PreparingDownloading
	j.StateReason = "downloading input files"
	if err = m.UpdateAndPushJob(ctx, j); err != nil {
		return errWrap(err, "fsm:preparing")
	}

	if err = m.storage.Download(ctx, j.Id, j.Inputs, j.Workspace); err != nil {
		err = fmt.Errorf("download all inputs failed, %w", err)
		log.Error(err)
		return errWrap(err, "fsm:input:download")
	}

	// 创建作业的准备文件目录
	opts := make([]oshelp.Option, 0)
	username := m.cfg.BackendProvider.SchedulerCommon.SubmitSysUser
	if username != "" {
		opts = append(opts, oshelp.WithChown(username))
	}

	preparedPath := filepath.Join(m.cfg.PreparedFilePath, fmt.Sprintf("%d", j.Id))

	if err := oshelp.Mkdir(preparedPath, opts...); err != nil {
		return errWrap(err, "fsm:mkdir prepared path")
	}

	j.PreparedFilePath = preparedPath

	// 开始提交作业到调度器上
	j.SubState = jobsubstate.PreparingSubmitting
	j.StateReason = "submitting the job"
	if err = m.UpdateAndPushJob(ctx, j); err != nil {
		return errWrap(err, "fsm:preparing")
	}

	// Replace the reserved flag in the command to check if the job failed during the prepared script
	j.Command = util.ReplaceCommand(j.Command, util.PreparedFlag, util.PreparedCmd(j, m.cfg.PreparedFilePath))

	// 提交作业并获取到作业的ID
	if j.OriginJobId, err = m.backend.Submit(ctx, j); err != nil {
		return errWrap(err, "fsm:backend:submit")
	}

	// 进入 Pending 状态
	j.SubState = jobsubstate.PendingWaitingSchedule
	j.StateReason = "waiting schedule"
	j.State = nextState(j.State, j.ControlBitTerminate, j.IsTimeout)
	if err = m.UpdateAndPushJob(ctx, j); err != nil {
		return errWrap(err, "fsm:preparing:push")
	}

	return nil
}

func (m *StateMachine) determineWorkspace(j *job.Job) error {
	if j.Workspace != "" {
		return nil
	}

	if j.IsOverride {
		// endpoint which parsed from workdir is useless for now, should replace it by config['hpc_storage_address']
		_, relativeWorkspace, err := util.ParseRawStorageUrl(j.WorkDir)
		if err != nil {
			err = fmt.Errorf("parse raw storage url failed, %w", err)
			log.Error(err)
			return err
		}

		workspaceAbsPath, err := m.storageClient.RealPath(m.cfg.HpcStorageAddress, relativeWorkspace)
		if err != nil {
			err = fmt.Errorf("get abs path failed, %w", err)
			log.Error(err)
			return fmt.Errorf("storage client get abs path failed, %w", err)
		}
		if workspaceAbsPath == "" {
			err = fmt.Errorf("base abs path failed, %w", err)
			log.Error(err)
			return err
		}

		j.Workspace = workspaceAbsPath
	} else {
		j.Workspace = m.backend.NewWorkspace()
	}

	opts := make([]oshelp.Option, 0)
	username := m.cfg.BackendProvider.SchedulerCommon.SubmitSysUser
	if username != "" {
		opts = append(opts, oshelp.WithChown(username))
	}

	if err := oshelp.Mkdir(j.Workspace, opts...); err != nil {
		return fmt.Errorf("mkdir workspace [%s] failed, %w", j.Workspace, err)
	}

	return nil
}

// Pending 排队中
// 检查作业状态 同步作业信息 若作业进入计算进入下一阶段
func (m *StateMachine) Pending(ctx context.Context, j *job.Job) error {
	j.TraceLogger.Infof("job pending: %v", j.Id)

	// 等待调度器中的作业进入 Running 状态或 Completed 状态
	err := m.watchState(ctx, j, func(nj *job.Job) (watcherState, error) {
		if nj.BackendJobState == job.StatePending {
			return holdWatcher, nil // 继续等待
		}
		// 从 Pending 转换到其他状态了就可以停止了
		return stopWatcher, nil
	})
	if err != nil {
		return errWrap(err, "fsm:pending:watch")
	}

	// 进入 Running 状态
	j.SubState = jobsubstate.RunningWaitingResult
	j.StateReason = "waiting for the result of the job"

	j.State = nextState(j.State, j.ControlBitTerminate, j.IsTimeout)
	if err = m.UpdateAndPushJob(ctx, j); err != nil {
		return errWrap(err, "fsm:pending:push")
	}

	return nil
}

// Running 运行中
// 检查作业状态 同步作业信息 若作业完成进入下一阶段
func (m *StateMachine) Running(ctx context.Context, j *job.Job) error {
	j.TraceLogger.Infof("job running: %v", j.Id)

	// 等待作业运行完成或等到任务超时返回错误
	var ddl *time.Time
	if j.RunningTime != nil {
		ddl = util.PTime(j.RunningTime.Add(time.Duration(j.Timeout) * time.Second))
	}
	err := m.watchState(ctx, j, func(j *job.Job) (watcherState, error) {
		if j.BackendJobState == job.StateCompleted {
			return stopWatcher, nil
		}

		// 如果有配置超时时间并且已经超时的话直接取消任务
		if j.RunningTime != nil {
			if ddl == nil {
				ddl = util.PTime(j.RunningTime.Add(time.Duration(j.Timeout) * time.Second))
			}

			if ddl.After(*j.RunningTime) && time.Now().After(*ddl) {
				j.IsTimeout = true
				j.StateReason = "job canceled by timeout"
				return stopWatcher, ErrJobCanceled
			}
		}

		return holdWatcher, nil
	})
	if err != nil {
		return errWrap(err, "fsm:running:watch")
	}

	// 进入 Completing 状态(不处理SubState)
	j.State = nextState(j.State, j.ControlBitTerminate, j.IsTimeout)
	if err = m.UpdateAndPushJob(ctx, j); err != nil {
		return errWrap(err, "fsm:running:push")
	}

	return nil
}

// Completing 完成中
// 同步工作目录至输出位置
func (m *StateMachine) Completing(ctx context.Context, j *job.Job) error {
	j.TraceLogger.Infof("job completing: %v", j.Id)

	// 回传计算结果
	j.SubState = jobsubstate.CompletingUploading
	j.StateReason = "uploading the result of the job"

	var err error

	if err = m.ensureJobState(j); err != nil {
		return errWrap(err, "fsm:ensureJobState")
	}

	if err = m.UpdateAndPushJob(ctx, j); err != nil {
		return errWrap(err, "fsm:completing:push")
	}

	return nil
}

// Completed 任务完成处理
func (m *StateMachine) Completed(ctx context.Context, j *job.Job) {
	j.TraceLogger.Infof("job completed: %d", j.Id)
	completedTime := time.Now()

	err := m.withRetry(_DefaultMarkJobCompletedRetryTimes,
		func(curr int, lastErr error, log util.WideLogger) (util.ReplayState, error) {
			if lastErr != nil {
				log("failed to complete the job", "times", curr, "sc_job", j.Id, "lastErr", lastErr)
			}

			// 标记任务的完成时间
			j.CompletedTime = &completedTime
			if err := m.UpdateAndPushJob(ctx, j); err != nil {
				return util.AutoStopRetry(errWrap(err, "fsm:completing:push"))
			}

			return util.StopReplay, nil
		},
	)
	if err != nil {
		j.TraceLogger.Errorw("failed to complete the job", "sc_job", j.Id, "error", err)
	}

	// 在远程服务上标记任务已完成
	m.OnCompleted(ctx, j)
}

// UpdateAndPushJob 将 job 更新到数据库中并推送到远算云上
// 如果在更新数据库或者推送状态到云上的过程中发生了错误, 比如网络错误, 数据库错误等问题
// 我们不应该结束任务, 而应该通过日志的方式报告错误并让运维人员去处理
// 如果发生的错误的要求取消任务, 此时应该返回这个错误并让最外层函数去去执行相对应的操作
func (m *StateMachine) UpdateAndPushJob(ctx context.Context, j *job.Job) error {
	ctx = m.BackgroundContext() // 强制忽略用户的取消信号
	if err := m.dao.UpdateJob(ctx, j.Job); err != nil {
		log.Errorw("update job failed", "sc_job", j.JobID(), "error", err)
	}

	return m.PushJob(ctx, j)
}

// PushJob 推送作业状态等信息到云服务上
// 如果在推送作业信息时发生了网络、接口等问题，应该通过日志的方式报告错误并让运维人员去处理
// 如果发生的错误的要求取消任务, 此时应该返回这个错误并让最外层函数去去执行相对应的操作
func (m *StateMachine) PushJob(ctx context.Context, j *job.Job) error {
	ctx = m.BackgroundContext() // 强制忽略用户的取消信号

	err := m.withRetry(_DefaultPushJobRetryTimes,
		func(curr int, lastErr error, log util.WideLogger) (util.ReplayState, error) {
			if lastErr != nil {
				log("failed to push the job", "times", curr, "sc_job", j.Id, "lastErr", lastErr)
			}

			return util.AutoStopRetry(m.WebHook(ctx, j))
		},
	)
	if err != nil {
		log.Errorw("push job failed", "sc_job", j.JobID(), "error", err)
	}

	return nil
}

// OnFailed 任务发生错误时的处理方法
func (m *StateMachine) OnFailed(ctx context.Context, j *job.Job, jobErr error) {
	j.TraceLogger.Warnf("job failed: %d, err: %v", j.Id, jobErr)

	if isJobUserFailed(jobErr) {
		j.IsUserFailed = true
	}

	// 更新本地数据库中的作业记录
	if err := m.dao.FailedJob(ctx, j.Job.Id, jobErr.Error()); err != nil {
		j.TraceLogger.Errorw("mark failed the job in database", "error", err, "sc_job", j.Id)
	}

	// 更新远算云作业失败信息
	j.State = jobstate.Failed
	j.StateReason = jobErr.Error()
	if err := m.UpdateAndPushJob(ctx, j); err != nil {
		j.TraceLogger.Errorw("push failed job", "error", err, "sc_job", j.Id)
	}
}

// OnCanceled 由于用户取消或者超时导致的作业停止
func (m *StateMachine) OnCanceled(ctx context.Context, j *job.Job) {
	j.TraceLogger.Infof("job canceled: %d", j.Id)

	// 只有在将作业提交到调度器之后才需要执行杀死操作
	if j.State == jobstate.Running || j.State == jobstate.Pending {
		err := m.withFastRetry(_DefaultKillJobRetryTimes,
			func(curr int, lastErr error, log util.WideLogger) (util.ReplayState, error) {
				if lastErr != nil {
					log("failed to kill the job", "sc_job", j.Id, "times", curr, "lastError", lastErr)
				}

				return util.AutoStopRetry(m.KillJob(ctx, j))
			},
		)
		if err != nil {
			j.TraceLogger.Errorw("failed to kill the job", "sc_job", j.Id, "error", err)
		}
	}

	exist, controlBitTerminate, err := dao.Default.GetJobControlBitTerminateById(ctx, j.Id)
	if err != nil {
		j.TraceLogger.Warnf("get job control_bit_terminate from db where job id = %d failed, %v", j.Id, err)
		controlBitTerminate = j.ControlBitTerminate
	}
	if !exist {
		j.TraceLogger.Warnf("job not exist in db where id = %d", j.Id)
		controlBitTerminate = j.ControlBitTerminate
	}

	j.State = jobstate.Canceled
	if controlBitTerminate {
		j.SubState = jobsubstate.CanceledByUser
		j.StateReason = "canceled by user"
	} else {
		j.SubState = jobsubstate.CanceledByTimeout
		j.StateReason = "canceled by timeout"
		j.TraceLogger.Warnf("job canceled: %d by timeout", j.Id)
	}
	if err := m.UpdateAndPushJob(ctx, j); err != nil {
		j.TraceLogger.Errorw("failed to push the job", "sc_job", j.Id, "error", err)
	}
}

// OnCompleted 在远程云服务上标记任务已经完成
func (m *StateMachine) OnCompleted(_ context.Context, j *job.Job) {
	_ = m.withRetry(_DefaultMarkJobCompletedRetryTimes,
		func(curr int, lastErr error, log util.WideLogger) (util.ReplayState, error) {
			if lastErr != nil {
				log("failed to mark the job completed", "sc_job", j.Id, "times", curr, "lastError", lastErr)
			}

			return util.AutoStopRetry(m.WebHook(m.BackgroundContext(), j))
		},
	)
}

// StateAction 每个状态相应的处理方法
type StateAction func(ctx context.Context, j *job.Job) error

// getAction 获取当前方法状态对应的动作
func (m *StateMachine) getAction(s jobstate.State) StateAction {
	switch s {
	case jobstate.Preparing:
		return m.Preparing
	case jobstate.Pending:
		return m.Pending
	case jobstate.Running:
		return m.Running
	case jobstate.Completing:
		return m.Completing
	}

	return nil
}

// KillJob 调用调度器取消任务的执行
func (m *StateMachine) KillJob(_ context.Context, j *job.Job) (err error) {
	// 取消任务的执行应该是不可中断或取消的, 所以这里直接使用全局上下文对象
	if err = m.backend.Kill(m.BackgroundContext(), j); err != nil {
		log.Warnw("kill the job failed", "sc_job", j.Id, "error", err)
	}
	return
}

func (m *StateMachine) WebHook(ctx context.Context, j *job.Job) error {
	if j.Webhook == "" {
		return nil
	}

	bs, err := json.Marshal(j.ToHTTPModel())
	if err != nil {
		return err
	}
	body := bytes.NewBuffer(bs)

	client, err := xhttp.New()
	if err != nil {
		return err
	}

	_, err = client.Post(j.Webhook, "application/json", body)
	return err
}

func (m *StateMachine) BackgroundContext() context.Context {
	return context.WithValue(context.Background(), with.OrmKey, m.db)
}

func (m *StateMachine) ensureJobState(j *job.Job) error {
	if j.ControlBitTerminate {
		return fmt.Errorf("user canceled, %w", ErrJobCanceled)
	}

	// check PreparedFile exist, if not exist or error, return system failed
	appPreparedFile := util.PreparedFile(j, m.cfg.PreparedFilePath)
	if _, err := os.Stat(appPreparedFile); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("app prepared file not exist")
		}

		return fmt.Errorf("stat app prepared file failed, %w", err)
	}

	if j.CustomStateRule == nil || j.CustomStateRule.KeyStatement == "" {
		j.State = jobstate.Completed
		j.SubState = jobsubstate.CompletedAllDone
		j.StateReason = "the job completed"
		return nil
	}

	switch j.CustomStateRule.ResultState {
	case jobstate.Completed.String():
		exist, err := checkFileExistContent(j.Stdout, j.CustomStateRule.KeyStatement)
		if err != nil {
			return fmt.Errorf("check file exist content failed, %w", err)
		}
		if exist {
			j.State = jobstate.Completed
			j.SubState = jobsubstate.CompletedAllDone
			j.StateReason = "job completed"
		} else {
			return ErrJobUserFailed
		}

		return nil
	case jobstate.Failed.String():
		exist, err := checkFileExistContent(j.Stdout, j.CustomStateRule.KeyStatement)
		if err != nil {
			return fmt.Errorf("check file exist content failed, %w", err)
		}
		if exist {
			return ErrJobUserFailed
		} else {
			j.State = jobstate.Completed
			j.SubState = jobsubstate.CompletedAllDone
			j.StateReason = "job completed"
		}

		return nil
	default:
		return fmt.Errorf("unsupported job state %s to be judged by custom", j.CustomStateRule.ResultState)
	}
}

func checkFileExistContent(file string, keyStatement string) (bool, error) {
	content, err := os.ReadFile(file)
	if err != nil {
		return false, fmt.Errorf("read file %s failed, %w", file, err)
	}

	if strings.Contains(string(content), keyStatement) {
		return true, nil
	}

	return false, nil
}

// NewStateMachine 创建状态机
func NewStateMachine(cfg *config.Config, be backend.Provider, cli registry.Client, s *storage.Manager,
	dao *dao.Dao, db *xorm.Engine) *StateMachine {

	// 每次作业轮询
	wi := _DefaultWatcherInterval
	if cfg.BackendProvider.CheckAliveInterval != 0 {
		wi = time.Duration(cfg.BackendProvider.CheckAliveInterval) * time.Second
	}

	return &StateMachine{
		cfg:           cfg,
		storage:       s,
		backend:       be,
		singularity:   cli,
		dao:           dao,
		db:            db,
		storageClient: client.New(cfg),
		//fileSyncPauseMgr: newFileSyncPauseMgr(),

		watcherInterval: wi,
	}
}

func nextState(state jobstate.State, canceled bool, isTimeout bool) jobstate.State {
	switch state {
	case jobstate.Preparing:
		return jobstate.Pending
	case jobstate.Pending:
		return jobstate.Running
	case jobstate.Running:
		return jobstate.Completing
	case jobstate.Completing:
		if canceled || isTimeout {
			return jobstate.Canceled
		}
		return jobstate.Completed
	default:
		return jobstate.Unknown
	}
}
