package post

import (
	"github.com/yuansuan/ticp/common/project-root-api/hpc/v1/job"
	v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
)

type Option func(req *job.SystemPostRequest)

func (api API) IdempotentId(idempotentId string) Option {
	return func(req *job.SystemPostRequest) {
		req.IdempotentID = idempotentId
	}
}

func (api API) Application(application string) Option {
	return func(req *job.SystemPostRequest) {
		req.Application = application
	}
}

func (api API) Environment(environment map[string]string) Option {
	return func(req *job.SystemPostRequest) {
		req.Environment = environment
	}
}

func (api API) Command(command string) Option {
	return func(req *job.SystemPostRequest) {
		req.Command = command
	}
}

func (api API) Override(override v20230530.JobInHPCOverride) Option {
	return func(req *job.SystemPostRequest) {
		req.Override = override
	}
}

func (api API) Queue(queue string) Option {
	return func(req *job.SystemPostRequest) {
		req.Queue = queue
	}
}

func (api API) Resource(resource v20230530.JobInHPCResource) Option {
	return func(req *job.SystemPostRequest) {
		req.Resource = resource
	}
}

func (api API) Inputs(inputs []v20230530.JobInHPCInputStorage) Option {
	return func(req *job.SystemPostRequest) {
		req.Inputs = inputs
	}
}

func (api API) Output(output *v20230530.JobInHPCOutputStorage) Option {
	return func(req *job.SystemPostRequest) {
		req.Output = output
	}
}

func (api API) CustomStateRule(customStateRule *v20230530.JobInHPCCustomStateRule) Option {
	return func(req *job.SystemPostRequest) {
		req.CustomStateRule = customStateRule
	}
}

func (api API) JobSchedulerSubmitFlags(jobSchedulerSubmitFlags map[string]string) Option {
	return func(req *job.SystemPostRequest) {
		req.JobSchedulerSubmitFlags = jobSchedulerSubmitFlags
	}
}
