package get

import (
	"github.com/yuansuan/ticp/common/project-root-api/hpc/v1/resource"
)

type Option func(req *resource.SystemGetRequest)

func (api API) Queue(queue string) Option {
	return func(req *resource.SystemGetRequest) {
		req.Queue = queue
	}
}
