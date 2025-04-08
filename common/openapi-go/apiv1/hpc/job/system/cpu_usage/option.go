package cpu_usage

import (
	"github.com/yuansuan/ticp/common/project-root-api/hpc/v1/cpuusage"
)

type Option func(req *cpuusage.SystemGetRequest)

func (api API) JobID(jobID string) Option {
	return func(req *cpuusage.SystemGetRequest) {
		req.JobID = jobID
	}
}
