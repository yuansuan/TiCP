package delete

import (
	"github.com/yuansuan/ticp/common/project-root-api/hpc/v1/job"
)

type Option func(req *job.SystemDeleteRequest)

func (api API) JobID(jobID string) Option {
	return func(req *job.SystemDeleteRequest) {
		req.JobID = jobID
	}
}
