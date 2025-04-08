package cancel

import (
	"github.com/yuansuan/ticp/common/project-root-api/hpc/v1/job"
)

type Option func(req *job.SystemCancelRequest)

func (api API) JobID(jobID string) Option {
	return func(req *job.SystemCancelRequest) {
		req.JobID = jobID
	}
}
