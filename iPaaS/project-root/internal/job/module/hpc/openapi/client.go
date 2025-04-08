package openapi

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/yuansuan/ticp/common/go-kit/logging"
	openapi "github.com/yuansuan/ticp/common/openapi-go"
	"github.com/yuansuan/ticp/common/openapi-go/credential"
	"github.com/yuansuan/ticp/common/project-root-api/hpc/v1/cpuusage"
	"github.com/yuansuan/ticp/common/project-root-api/hpc/v1/job"
	"github.com/yuansuan/ticp/common/project-root-api/hpc/v1/resource"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/config"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/module"
)

func DefaultCred(cfg config.CustomT) *credential.Credential {
	return credential.NewCredential(cfg.GetAK(), cfg.GetAS())
}

func DefaultRetryCondition(resp *http.Response, err error) bool {
	if err != nil {
		logging.Default().Info(err)
		return true
	}

	if resp != nil && resp.StatusCode >= http.StatusInternalServerError {
		logging.Default().Info("StatusCode: %d", resp.StatusCode)
		return true
	}

	return false
}

type openapiClientPool struct {
	pool *sync.Pool
}

func newOpenapiClientPool() *openapiClientPool {
	cp := &openapiClientPool{
		pool: &sync.Pool{},
	}

	cp.pool.New = func() interface{} {
		opts := []openapi.Option{
			openapi.WithBaseURL(openapi.DefaultBaseURL),
			openapi.WithTimeout(module.DefaultTimeout),
			openapi.WithRetryTimes(module.DefaultRetryTimes),
			openapi.WithRetryInterval(module.DefaultRetryInterval),
			openapi.WithRetryCondition(DefaultRetryCondition),
		}
		openapiClient, err := openapi.NewClient(credential.NewCredential("", ""), opts...)
		if err != nil {
			logging.Default().Errorf("new client failed, %v", err)
			return nil
		}

		return openapiClient
	}

	return cp
}

func (cp *openapiClientPool) Get(endpoint string, timeout time.Duration, cred *credential.Credential) (*openapi.Client, error) {
	cli, ok := cp.pool.Get().(*openapi.Client)
	if !ok {
		return nil, fmt.Errorf("get from pool cannot convert to *openapi.Client")
	}

	// 动态设置baseUrl/timeout/credential
	cli.SetBaseUrl(endpoint)
	cli.SetTimeout(timeout)
	if err := cli.SetCredential(cred); err != nil {
		return nil, fmt.Errorf("openapi client set credentail failed, %w", err)
	}

	return cli, nil
}

func (cp *openapiClientPool) Put(openapiClient *openapi.Client) {
	if openapiClient == nil {
		return
	}

	// 丢回pool时重新至空
	openapiClient.SetBaseUrl("")
	openapiClient.SetTimeout(module.DefaultTimeout)
	_ = openapiClient.SetCredential(credential.NewCredential("", ""))
	cp.pool.Put(openapiClient)
}

type client struct {
	OpenapiClientPool *openapiClientPool
}

var _client *client
var once sync.Once

func Client() *client {
	once.Do(func() {
		_client = newClient()
	})

	return _client
}

func newClient() *client {
	return &client{
		OpenapiClientPool: newOpenapiClientPool(),
	}
}

func (c *client) GetResource(endpoint string) (*resource.SystemGetResponse, error) {
	cli, err := c.OpenapiClientPool.Get(endpoint, module.DefaultTimeout, DefaultCred(config.GetConfig()))
	if err != nil {
		return nil, fmt.Errorf("get openapi client failed, %w", err)
	}
	defer c.OpenapiClientPool.Put(cli)

	return cli.HPC.Resource.System.Get()
}

func (c *client) PostJob(endpoint string, timeout time.Duration, req job.SystemPostRequest) (*job.SystemPostResponse, error) {
	cli, err := c.OpenapiClientPool.Get(endpoint, timeout, DefaultCred(config.GetConfig()))
	if err != nil {
		return nil, fmt.Errorf("get openapi client failed, %w", err)
	}
	defer c.OpenapiClientPool.Put(cli)

	api := cli.HPC.Job.System.Post
	return api(
		api.Application(req.Application),
		api.Command(req.Command),
		api.Environment(req.Environment),
		api.Inputs(req.Inputs),
		api.Output(req.Output),
		api.Override(req.Override),
		api.Resource(req.Resource),
		api.Queue(req.Queue),
		api.IdempotentId(req.IdempotentID),
		api.CustomStateRule(req.CustomStateRule),
		api.JobSchedulerSubmitFlags(req.JobSchedulerSubmitFlags),
	)
}

func (c *client) CancelJob(endpoint string, timeout time.Duration, jobID string) (*job.SystemCancelResponse, error) {
	cli, err := c.OpenapiClientPool.Get(endpoint, timeout, DefaultCred(config.GetConfig()))
	if err != nil {
		return nil, fmt.Errorf("get openapi client failed, %w", err)
	}
	defer c.OpenapiClientPool.Put(cli)

	api := cli.HPC.Job.System.Cancel
	return api(
		api.JobID(jobID),
	)
}

func (c *client) GetJob(endpoint string, timeout time.Duration, jobID string) (*job.SystemGetResponse, error) {
	cli, err := c.OpenapiClientPool.Get(endpoint, timeout, DefaultCred(config.GetConfig()))
	if err != nil {
		return nil, fmt.Errorf("get openapi client failed, %w", err)
	}
	defer c.OpenapiClientPool.Put(cli)

	api := cli.HPC.Job.System.Get
	return api(
		api.JobID(jobID),
	)
}

func (c *client) GetJobCpuUsage(endpoint string, timeout time.Duration, jobID string) (*cpuusage.SystemGetResponse, error) {
	cli, err := c.OpenapiClientPool.Get(endpoint, timeout, DefaultCred(config.GetConfig()))
	if err != nil {
		return nil, fmt.Errorf("get openapi client failed, %w", err)
	}
	defer c.OpenapiClientPool.Put(cli)

	api := cli.HPC.Job.System.GetCpuUsage
	return api(
		api.JobID(jobID),
	)
}
