package job

import (
	"context"
	"sync"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/admin/joblistfiltered"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/jobcpuusage"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/jobcreate"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/jobdelete"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/jobget"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/joblist"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/jobmonitorchart"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/jobpreschedule"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/jobresidual"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/jobresume"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/jobsnapshotget"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/jobsnapshotlist"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/jobterminate"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/jobtransmitresume"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/jobtransmitsuspend"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/system/jobneedsyncfile"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/system/jobsyncfilestate"
	schema "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao/models"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/service/scheduler/monitorchart"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/service/scheduler/residual"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/service/v1/application"
)

type allowFunc func(opUser, owner string) bool
type createConvertFunc func(ctx context.Context, jobID snowflake.ID) (*models.Job, error)

//go:generate mockgen -self_package=github.com/yuansuan/ticp/iPaaS/project-root/internal/job/service/v1/job -destination mock_job_srv.go -package job github.com/yuansuan/ticp/iPaaS/project-root/internal/job/service/v1/job Service

// Service 服务
type Service interface {
	Get(ctx context.Context, req *jobget.Request, userID snowflake.ID, allow allowFunc, withDelete bool) (*models.Job, error)
	BatchGet(ctx context.Context, ids []string, userID snowflake.ID) ([]*models.Job, error)
	List(ctx context.Context, req *joblist.Request, userID, appID snowflake.ID, withDelete, isSystemFailed bool) (int64, []*models.Job, error)
	ListFiltered(ctx context.Context, req *joblistfiltered.Request, userID, appID snowflake.ID) (int64, []*models.Job, error)
	Create(ctx context.Context, req *jobcreate.Request, userID snowflake.ID, chargeParams schema.ChargeParams, appInfo *models.Application, scheduleInfo *models.PreSchedule, convert createConvertFunc) (string, error)
	Delete(ctx context.Context, req *jobdelete.Request, userID snowflake.ID, allow allowFunc) error
	Terminate(ctx context.Context, req *jobterminate.Request, userID snowflake.ID, allow allowFunc) error
	Resume(ctx context.Context, req *jobresume.Request, userID snowflake.ID, allow allowFunc) error
	TransmitResume(ctx context.Context, req *jobtransmitresume.Request, userID snowflake.ID, allow allowFunc) error
	TransmitSuspend(ctx context.Context, req *jobtransmitsuspend.Request, userID snowflake.ID, allow allowFunc) error
	GetResidual(ctx context.Context, req *jobresidual.Request, userID snowflake.ID, allow allowFunc) (*schema.Residual, error)
	GetMonitorChart(ctx context.Context, req *jobmonitorchart.Request, userID snowflake.ID, allow allowFunc) (*models.MonitorChart, error)
	GetJobSnapshot(ctx context.Context, appSrv application.AppSrv, req *jobsnapshotget.Request, userID snowflake.ID, allow allowFunc) (string, error)
	ListJobSnapshot(ctx context.Context, appSrv application.AppSrv, req *jobsnapshotlist.Request, userID snowflake.ID, allow allowFunc) (map[string][]string, error)
	PreSchedule(ctx context.Context, req *jobpreschedule.Request, zones schema.Zones, userID snowflake.ID, appInfo *models.Application) (jobpreschedule.Data, error)
	GetPreSchedule(ctx context.Context, preScheduleID string) (*models.PreSchedule, bool, error)
	ListNeedSyncFileJobs(ctx context.Context, req *jobneedsyncfile.Request) (*jobneedsyncfile.Data, error)
	UpdateSyncFileState(ctx context.Context, req *jobsyncfilestate.Request, jobIDStr string) error
	GetCpuUsage(ctx context.Context, req *jobcpuusage.Request, userID snowflake.ID, allow allowFunc) (*schema.JobCpuUsage, error)
}

// jobService
type jobService struct {
	snowflake.IDGen
	jobdao          dao.JobDao
	residualdao     dao.ResidualDao
	residualHandler residual.ResidualHandler
	jobPlugins      []JobPlugin
}

var _ Service = (*jobService)(nil)

var jobSrv Service

var once sync.Once

func (js *jobService) RegisterJobPlugin(plugin ...JobPlugin) {
	js.jobPlugins = append(js.jobPlugins, plugin...)
}

// newJobService return a new JobSrv.
func newJobService(idgen snowflake.IDGen, jobdao dao.JobDao, residualdao dao.ResidualDao, residualHandler residual.ResidualHandler) *jobService {
	js := &jobService{
		IDGen:           idgen,
		jobdao:          jobdao,
		residualdao:     residualdao,
		residualHandler: residualHandler,
	}

	js.RegisterJobPlugin(
		residual.NewResidualPlugin(idgen, residualdao), // 残差图
		monitorchart.NewMonitorChartPlugin(),           // 监控图表
	)
	return js
}

// NewJobService return the JobSrv instance.
func NewJobService(idgen snowflake.IDGen, jobdao dao.JobDao, residualdao dao.ResidualDao, residualHandler residual.ResidualHandler) Service {
	once.Do(func() {
		jobSrv = newJobService(idgen, jobdao, residualdao, residualHandler)
	})
	return jobSrv
}

type JobPlugin interface {
	Name() string
	Insert(ctx context.Context, app *models.Application, job *models.Job)
}
