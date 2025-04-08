package get

import "github.com/yuansuan/ticp/common/project-root-api/hpc/v1/job"

type Option func(req *job.SystemGetRequest)

func (api API) JobID(jobID string) Option {
	return func(req *job.SystemGetRequest) {
		req.JobID = jobID
	}
}
