package mock

import (
	"context"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/dao/models"
	"path/filepath"

	"github.com/google/uuid"
	v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"

	"github.com/yuansuan/ticp/iPaaS/standard-compute/config"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/backend/job"
)

type Provider struct {
	commonCfg config.SchedulerCommon
	customCfg *config.MockBackendProvider
}

func NewProvider(commonCfg config.SchedulerCommon, customCfg *config.MockBackendProvider) *Provider {
	return &Provider{
		commonCfg: commonCfg,
		customCfg: customCfg,
	}
}

func (p *Provider) Submit(ctx context.Context, j *job.Job) (string, error) {
	return "", nil
}

func (p *Provider) Kill(ctx context.Context, j *job.Job) error {
	return nil
}

func (p *Provider) CheckAlive(ctx context.Context, j *job.Job) (*job.Job, error) {
	return j, nil
}

func (p *Provider) NewWorkspace() string {
	return filepath.Join(p.commonCfg.Workspace, uuid.NewString())
}

func (p *Provider) GetFreeResource(ctx context.Context, queues []string) (map[string]*v20230530.Resource, error) {
	return nil, nil
}
func (p *Provider) GetCpuUsage(ctx context.Context, j *models.Job, nodes []string, adjustFactor float64) (*v20230530.CpuUsage, error) {
	return &v20230530.CpuUsage{}, nil
}
