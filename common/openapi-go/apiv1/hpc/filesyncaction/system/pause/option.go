package pause

import (
	"github.com/yuansuan/ticp/common/project-root-api/hpc/v1/filesyncaction"
)

type Option func(req *filesyncaction.SystemPauseFileSyncRequest)

func (api API) JobID(jobId string) Option {
	return func(req *filesyncaction.SystemPauseFileSyncRequest) {
		req.JobID = jobId
	}
}
