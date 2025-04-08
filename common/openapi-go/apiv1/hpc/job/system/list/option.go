package list

import (
	"github.com/yuansuan/ticp/common/project-root-api/hpc/v1/job"
)

type Option func(req *job.SystemListRequest)

func (api API) PageOffset(pageOffset int) Option {
	return func(req *job.SystemListRequest) {
		req.PageOffset = pageOffset
	}
}

func (api API) PageSize(pageSize int) Option {
	return func(req *job.SystemListRequest) {
		req.PageSize = pageSize
	}
}

func (api API) Status(status string) Option {
	return func(req *job.SystemListRequest) {
		req.Status = status
	}
}

func (api API) IDs(ids []string) Option {
	return func(req *job.SystemListRequest) {
		req.JobIDs = ids
	}
}
