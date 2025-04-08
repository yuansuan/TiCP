package dao

import (
	"context"
	"fmt"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/admin/joblistfiltered"
	v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
	"time"

	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/pkg/common/with"

	"xorm.io/xorm"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/consts"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao/models"
)

// JobDaoImpl dao
type JobDaoImpl struct {
	db *xorm.Engine
}

// NewJobDaoImpl 实现
func NewJobDaoImpl(db *xorm.Engine) *JobDaoImpl {
	return &JobDaoImpl{
		db: db,
	}
}

// Engine engine
func (j *JobDaoImpl) Engine() *xorm.Engine {
	return j.db
}

// Transaction 事务，通过with.KeepSession保证action中所有的with.DefaultSession都使用同一个事务session
func (j *JobDaoImpl) Transaction(ctx context.Context, action func(context.Context) error) error {
	_, err := j.db.Transaction(func(db *xorm.Session) (interface{}, error) {
		return nil, action(with.KeepSession(ctx, db))
	})
	return err
}

// Get 查询作业
func (j *JobDaoImpl) Get(ctx context.Context, jobID snowflake.ID,
	forUpdate bool, withDelete bool) (*models.Job, error) {
	job := models.Job{}
	err := with.DefaultSession(ctx, func(db *xorm.Session) error {
		if forUpdate {
			db = db.ForUpdate()
		}

		if !withDelete {
			db = db.Where("is_deleted = ?", 0)
		}

		b, err := db.ID(jobID).Get(&job)
		if err != nil {
			return err
		}
		if !b {
			return common.ErrJobIDNotFound
		}

		return nil
	})
	return &job, err
}

// BatchGet 批量查询作业
func (j *JobDaoImpl) BatchGet(ctx context.Context, jobIDs []snowflake.ID,
	userID snowflake.ID, forUpdate bool, withDelete bool) ([]*models.Job, error) {
	jobs := []*models.Job{}

	err := with.DefaultSession(ctx, func(db *xorm.Session) (err error) {
		if forUpdate {
			db = db.ForUpdate()
		}

		if !withDelete {
			db = db.Where("is_deleted = ?", 0)
		}

		err = db.Where("user_id = ?", userID).In("id", jobIDs).Find(&jobs)

		return err
	})
	if err != nil {
		return nil, err
	}
	return jobs, nil
}

// ListJobs 跟据用户ID、区域、作业状态查询作业 (不包含已删除的作业)
func (j *JobDaoImpl) ListJobs(ctx context.Context, offset, limit int,
	userID, appID snowflake.ID, zone, jobState string, withDelete,
	isSystemFailed bool) (int64, []*models.Job, error) {
	jobs := []*models.Job{}
	count := int64(0)

	err := with.DefaultSession(ctx, func(db *xorm.Session) (err error) {
		db.Table(consts.TableJob)

		if zone != "" {
			db = db.Where("zone = ?", zone)
		}

		if jobState != "" {
			state, _ := consts.GetStateValue(jobState)
			db.Where("state = ?", state)
		}

		if userID != 0 {
			db.Where("user_id = ?", userID)
		}

		if appID != 0 {
			db.Where("app_id = ?", appID)
		}

		if !withDelete {
			db.Where("is_deleted = ?", 0)
		}

		if isSystemFailed {
			db.Where("is_system_failed = ?", 1)
		}

		count, err = db.
			Limit(limit, offset).
			OrderBy("create_time desc").
			FindAndCount(&jobs)

		return err
	})
	if err != nil {
		return 0, nil, err
	}
	return count, jobs, nil
}

func (j *JobDaoImpl) ListJobsFiltered(ctx context.Context, offset, limit int,
	in *joblistfiltered.Request, userID, appID snowflake.ID) (int64, []*models.Job, error) {
	jobs := []*models.Job{}
	count := int64(0)

	err := with.DefaultSession(ctx, func(db *xorm.Session) (err error) {
		db.Table(consts.TableJob)
		if in.Zone != "" {
			db.Where("zone = ?", in.Zone)
		}
		if in.JobState != "" {
			state, _ := consts.GetStateValue(in.JobState)
			db.Where("state = ?", state)
		}
		if in.FileSyncState != "" {
			db.Where("file_sync_state =?", in.FileSyncState)
		}
		if in.Name != "" {
			db.Where("name = ?", in.Name)
		}
		if in.JobID != "" {
			JobID, err := snowflake.ParseString(in.JobID)
			if err != nil {
				return fmt.Errorf("invalid JobID: %w", err)
			}
			db.Where("id = ?", JobID.Int64())
		}
		if in.AccountID != "" {
			AccountID, err := snowflake.ParseString(in.AccountID)
			if err != nil {
				return fmt.Errorf("invalid AccountID: %w", err)
			}
			db.Where("account_id = ?", AccountID.Int64())
		}
		if userID != 0 {
			db.Where("user_id = ?", userID)
		}
		if appID != 0 {
			db.Where("app_id = ?", appID)
		}
		if !in.WithDelete {
			db.Where("is_deleted = ?", 0)
		}
		if in.IsSystemFailed {
			db.Where("is_system_failed = ?", 1)
		}
		if !in.StartTime.IsZero() {
			db.Where("create_time >= ?", in.StartTime)
		}
		if !in.EndTime.IsZero() {
			db.Where("create_time < ?", in.EndTime)
		}
		count, err = db.Limit(limit, offset).
			OrderBy("create_time desc").
			FindAndCount(&jobs)
		return err
	})
	if err != nil {
		logger := logging.GetLogger(ctx).With("function", "ListJobsFiltered",
			"offset", offset,
			"limit", limit,
			"userID", userID,
			"appID", appID,
			"request", in)
		logger.Warnf("ListJobsFiltered failed: %f", err)
		return 0, nil, err
	}
	return count, jobs, nil
}

// GetJobResidual 获取作业残差图
func (j *JobDaoImpl) GetJobResidual(ctx context.Context, jobID snowflake.ID) (*models.Residual, error) {
	residual := new(models.Residual)
	err := with.DefaultSession(ctx, func(db *xorm.Session) error {
		exist, err := db.Where("job_id = ?", jobID).Get(residual)
		if err != nil {
			return err
		}
		if !exist {
			return common.ErrJobResidualNotFound
		}

		return nil
	})
	return residual, err
}

// GetUnfinishedResidual 获取未完成的残差图
func (j *JobDaoImpl) GetUnfinishedResidual(ctx context.Context) ([]*models.Residual, error) {
	residuals := make([]*models.Residual, 0)
	err := with.DefaultSession(ctx, func(db *xorm.Session) error {
		return db.Where("finished = ?", false).Find(&residuals)
	})
	return residuals, err
}

// InsertResidual 插入残差图
func (j *JobDaoImpl) InsertResidual(ctx context.Context, residual *models.Residual) error {
	return with.DefaultSession(ctx, func(db *xorm.Session) error {
		_, err := db.Insert(residual)
		return err
	})
}

// UpdateResidual 更新残差图
func (j *JobDaoImpl) UpdateResidual(ctx context.Context, residual *models.Residual) error {
	return with.DefaultSession(ctx, func(db *xorm.Session) error {
		_, err := db.ID(residual.ID).UseBool().Update(residual)
		return err
	})
}

// GetJobMonitorChart 获取作业监控图表
func (j *JobDaoImpl) GetJobMonitorChart(ctx context.Context, jobId snowflake.ID,
	forUpdate bool) (*models.MonitorChart, error) {
	monitorChart := new(models.MonitorChart)
	err := with.DefaultSession(ctx, func(db *xorm.Session) error {
		if forUpdate {
			db = db.ForUpdate()
		}

		exist, err := db.Where("job_id = ?", jobId).Get(monitorChart)
		if err != nil {
			return err
		}
		if !exist {
			return common.ErrJobMonitorChartNotFound
		}

		return nil
	})
	return monitorChart, err
}

// GetUnfinishedbMonitorChart 获取未完成的监控图表
func (j *JobDaoImpl) GetUnfinishedbMonitorChart(ctx context.Context) ([]*models.MonitorChart, error) {
	monitorCharts := make([]*models.MonitorChart, 0)
	err := with.DefaultSession(ctx, func(db *xorm.Session) error {
		return db.Where("finished = ?", false).Find(&monitorCharts)
	})
	return monitorCharts, err
}

// UpdateSubmitJob 更新提交的作业
func (j *JobDaoImpl) UpdateSubmitJob(ctx context.Context, job *models.Job) (int64, error) {
	result := int64(0)
	err := with.DefaultSession(ctx, func(db *xorm.Session) error {
		// 作业提交后应该修改的字段: hpc_job_id, submit_time, zone, resource_usage_cpus,
		//	resource_usage_memory, state_reason, update_time, command, state, sub_state, work_dir
		db = db.ID(job.ID).Cols("hpc_job_id", "submit_time", "zone", "resource_usage_cpus",
			"resource_usage_memory", "state_reason", "update_time", "command", "state", "sub_state", "work_dir")
		r, err := db.Update(job)
		if err != nil {
			return err
		}
		result = r

		return nil
	})
	return result, err
}

// UpdateSchedulingReason 更新调度中的作业原因
func (j *JobDaoImpl) UpdateSchedulingReason(ctx context.Context, job *models.Job) error {
	err := with.DefaultSession(ctx, func(db *xorm.Session) error {
		_, err := db.ID(job.ID).Cols("update_time", "state_reason").Update(job)
		return err
	})
	return err
}

// ListJobsBySubStates 跟据 subState 列表查询作业 (不包含已删除的作业)
func (j *JobDaoImpl) ListJobsBySubStates(ctx context.Context, subState ...int) (int64, []*models.Job, error) {
	jobs := []*models.Job{}
	count := int64(0)
	err := with.DefaultSession(ctx, func(db *xorm.Session) error {
		c, err := db.
			In("sub_state", subState).
			Where("is_deleted = ?", 0).
			FindAndCount(&jobs)
		if err != nil {
			return err
		}

		count = c
		return nil
	})
	return count, jobs, err
}

// ListSchedulerTransferJobs 查询需要传输暂停、传输恢复的作业
func (j *JobDaoImpl) ListSchedulerTransferJobs(ctx context.Context) (int64, []*models.Job, error) {
	jobs := make([]*models.Job, 0)
	count := int64(0)
	err := with.DefaultSession(ctx, func(db *xorm.Session) error {
		c, err := db.Where("is_deleted = ?", 0).Where("(sub_state IN (?, ?, ?, ?) AND file_sync_state IN (?, ?))",
			consts.SubStateRunning.SubState, consts.SubStateCompleted.SubState,
			consts.SubStateFailed.SubState, consts.SubStateTerminated.SubState,
			consts.FileSyncStatePausing, consts.FileSyncStateResuming).
			FindAndCount(&jobs)
		if err != nil {
			return err
		}

		count = c
		return nil
	})
	return count, jobs, err
}

// ListInputHpcFinalSyncingJobs 查询终态但仍在传输的作业
func (j *JobDaoImpl) ListInputHpcFinalSyncingJobs(ctx context.Context) (int64, []*models.Job, error) {
	jobs := make([]*models.Job, 0)
	count := int64(0)
	err := with.DefaultSession(ctx, func(db *xorm.Session) error {
		c, err := db.Where("is_deleted = ?", 0).
			And("sub_state In (?, ?, ?)",
				consts.SubStateCompleted.SubState, consts.SubStateTerminated.SubState, consts.SubStateFailed.SubState).
			And("input_type = ?", consts.HpcStorage).
			And("file_sync_state = '' or file_sync_state not in(?, ?, ?)",
				consts.FileSyncStateUnknown, consts.FileSyncStateFailed, consts.FileSyncStateCompleted).
			FindAndCount(&jobs)
		if err != nil {
			return err
		}

		count = c
		return nil
	})
	return count, jobs, err
}

func (j *JobDaoImpl) ListNeedFileSyncJobs(ctx context.Context, zone string,
	offset, limit int64) ([]*models.Job, int64, error) {
	var jobs []*models.Job
	count := int64(0)
	err := with.DefaultSession(ctx, func(db *xorm.Session) error {
		needSyncFileJobState := []int{
			consts.Running,
			consts.Completed,
			consts.Terminated,
			consts.Failed,
		}

		needSyncFileState := []consts.FileSyncState{
			consts.FileSyncStateWaiting,
			consts.FileSyncStateSyncing,
			consts.FileSyncStatePausing,
			consts.FileSyncStateResuming,
			consts.FileSyncStateNone,
		}

		columns := "id, name, state, sub_state, file_sync_state, work_dir, output_dir, " +
			"no_needed_paths, needed_paths, file_output_storage_zone, download_file_size_total, " +
			"download_file_size_current, download_finished, transmitting_time, download_time"
		c, err := db.Cols(columns).
			Where("is_deleted = ?", 0).
			And("file_output_storage_zone = ?", zone).
			And("output_type = ?", consts.CloudStorage).
			And("download_finished = ?", 0).
			And("hpc_job_id != ''").
			In("state", needSyncFileJobState).
			In("file_sync_state", needSyncFileState).
			Limit(int(limit), int(offset)).FindAndCount(&jobs)
		if err != nil {
			return err
		}

		count = c
		return nil
	})
	return jobs, count, err
}

func (j *JobDaoImpl) ListShouldPostPaidJobs(ctx context.Context) ([]*models.Job, error) {
	jobs := make([]*models.Job, 0)
	err := with.DefaultSession(ctx, func(db *xorm.Session) error {
		err := db.Where("charge_type = ?", v20230530.PostPaid).
			Where("is_paid_finished = ?", false).
			Where("execution_duration != ?", 0).
			In("state", []int{consts.Running, consts.Suspending, consts.Suspended,
				consts.Terminating, consts.Terminated, consts.Completed, consts.Failed}).
			Find(&jobs)
		return err
	})
	return jobs, err
}

func (j *JobDaoImpl) GetBill(ctx context.Context, jobId snowflake.ID) (*models.Bill, bool, error) {
	bill := new(models.Bill)
	exist := false
	err := with.DefaultSession(ctx, func(db *xorm.Session) error {
		e, err := db.
			Where("job_id = ?", jobId).
			Get(bill)
		if err != nil {
			return err
		}

		exist = e
		return nil
	})
	return bill, exist, err
}

func (j *JobDaoImpl) InsertBill(ctx context.Context, bill *models.Bill) error {
	err := with.DefaultSession(ctx, func(db *xorm.Session) error {
		_, err := db.Insert(bill)
		return err
	})
	return err
}

func (j *JobDaoImpl) UpdateBilledDurationAndBillTimeByJobId(ctx context.Context,
	jobId snowflake.ID, billedDuration int64, billTime time.Time) error {
	err := with.DefaultSession(ctx, func(db *xorm.Session) error {
		_, err := db.
			Where("job_id = ?", jobId).
			Cols("billed_duration", "bill_time").
			Update(&models.Bill{
				BilledDuration: billedDuration,
				BillTime:       billTime,
			})
		return err
	})
	return err
}

func (j *JobDaoImpl) MarkJobPaidFinished(ctx context.Context, jobId snowflake.ID) error {
	err := with.DefaultSession(ctx, func(db *xorm.Session) error {
		_, err := db.ID(jobId).
			Cols("is_paid_finished").
			Update(&models.Job{
				IsPaidFinished: true,
			})
		return err
	})
	return err
}

func (j *JobDaoImpl) GetPreSchedule(ctx context.Context,
	preScheduleID snowflake.ID) (*models.PreSchedule, bool, error) {
	preSchedule := new(models.PreSchedule)
	exist := false
	err := with.DefaultSession(ctx, func(db *xorm.Session) error {
		e, err := db.
			Where("id = ?", preScheduleID).
			Get(preSchedule)
		if err != nil {
			return err
		}

		exist = e
		return nil
	})
	return preSchedule, exist, err
}
