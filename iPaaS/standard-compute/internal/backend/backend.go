package backend

import (
	"context"
	"errors"
	v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/config"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/backend/job"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/backend/mock"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/backend/pbspro"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/backend/slurm"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/cmdhelp"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/dao/models"
)

type Provider interface {
	Submit(ctx context.Context, j *job.Job) (string, error)
	Kill(ctx context.Context, j *job.Job) error
	CheckAlive(ctx context.Context, j *job.Job) (*job.Job, error)
	NewWorkspace() string
	GetFreeResource(ctx context.Context, queues []string) (map[string]*v20230530.Resource, error)
	GetCpuUsage(ctx context.Context, j *models.Job, nodes []string, adjustFactor float64) (*v20230530.CpuUsage, error)
}

func NewProvider(cfg *config.BackendProvider) (Provider, error) {
	switch cfg.Type {
	case "slurm":
		return slurm.NewProvider(cfg.SchedulerCommon, cfg.Slurm, cmdhelp.ExecShellCmd), nil
	case "pbs-pro":
		return pbspro.NewProvider(cfg.SchedulerCommon, cfg.PbsPro, cmdhelp.ExecShellCmd), nil
	case "mock":
		return mock.NewProvider(cfg.SchedulerCommon, cfg.Mock), nil
	default:
		return nil, errors.New("unknown backend provider")
	}
}
