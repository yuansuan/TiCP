package dao

//go:generate mockgen -self_package=github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao -destination mock_job_dao.go -package dao github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao JobDao,ResidualDao
import (
	"context"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/admin/joblistfiltered"
	"time"

	"xorm.io/xorm"

	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/consts"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao/models"
	app "github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao/models"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
)

// ApplicationDao Application interface
type ApplicationDao interface {
	GetApp(ctx context.Context, id snowflake.ID) (*app.Application, error)
	ListApps(ctx context.Context, userID snowflake.ID, publishStatus consts.PublishStatus) ([]*app.Application, int64, error)
	AddApp(ctx context.Context, appInfo *app.Application) error
	UpdateApp(ctx context.Context, appInfo *app.Application) error
}

// ApplicationQuotaDao ApplicationQuota dao interface
type ApplicationQuotaDao interface {
	GetByUser(ctx context.Context, session *xorm.Session, applicationID, userID snowflake.ID, forUpdate bool) (*models.ApplicationQuota, error)
}

// ApplicationAllowDao ApplicationAllow dao interface
type ApplicationAllowDao interface {
	GetByAppId(ctx context.Context, session *xorm.Session, applicationID snowflake.ID) (*models.ApplicationAllow, error)
}

// JobDao job dao interface
type JobDao interface {
	Engine() *xorm.Engine
	Transaction(ctx context.Context, action func(context.Context) error) error

	Get(ctx context.Context, jobID snowflake.ID, forUpdate bool, withDelete bool) (*models.Job, error)
	BatchGet(ctx context.Context, jobIDs []snowflake.ID, userID snowflake.ID, forUpdate bool, withDelete bool) ([]*models.Job, error)
	ListJobs(ctx context.Context, offset, limit int, userID, appID snowflake.ID, zone, jobState string, withDelete, isSystemFailed bool) (int64, []*models.Job, error)
	ListJobsFiltered(ctx context.Context, offset, limit int, in *joblistfiltered.Request, userID, appID snowflake.ID) (int64, []*models.Job, error)
	GetJobMonitorChart(ctx context.Context, jobId snowflake.ID, forUpdate bool) (*models.MonitorChart, error)
	GetUnfinishedbMonitorChart(ctx context.Context) ([]*models.MonitorChart, error)
	UpdateSubmitJob(ctx context.Context, job *models.Job) (int64, error)
	UpdateSchedulingReason(ctx context.Context, job *models.Job) error
	ListJobsBySubStates(ctx context.Context, subState ...int) (int64, []*models.Job, error)
	ListSchedulerTransferJobs(ctx context.Context) (int64, []*models.Job, error)
	ListInputHpcFinalSyncingJobs(ctx context.Context) (int64, []*models.Job, error)
	ListNeedFileSyncJobs(ctx context.Context, zone string, offset, limit int64) ([]*app.Job, int64, error)
	ListShouldPostPaidJobs(ctx context.Context) ([]*models.Job, error)
	GetBill(ctx context.Context, jobId snowflake.ID) (*models.Bill, bool, error)
	InsertBill(ctx context.Context, bill *models.Bill) error
	UpdateBilledDurationAndBillTimeByJobId(ctx context.Context, jobId snowflake.ID, billedDuration int64, billTime time.Time) error
	MarkJobPaidFinished(ctx context.Context, jobId snowflake.ID) error
	GetPreSchedule(ctx context.Context, preScheduleID snowflake.ID) (*models.PreSchedule, bool, error)
}

type ResidualDao interface {
	GetJobResidual(ctx context.Context, jobId snowflake.ID) (*models.Residual, error)
	GetUnfinishedResidual(ctx context.Context) ([]*models.Residual, error)
	InsertResidual(ctx context.Context, residual *models.Residual) error
	UpdateResidual(ctx context.Context, residual *models.Residual) error
}
