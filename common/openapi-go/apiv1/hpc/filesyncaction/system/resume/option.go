package resume

import (
	"github.com/yuansuan/ticp/common/project-root-api/hpc/v1/filesyncaction"
)

type Option func(req *filesyncaction.SystemResumeFileSyncRequest)

func (api API) JobID(jobId string) Option {
	return func(req *filesyncaction.SystemResumeFileSyncRequest) {
		req.JobID = jobId
	}
}
