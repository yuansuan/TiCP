package delete

import (
	"fmt"
	"net/http"

	"github.com/yuansuan/ticp/common/project-root-api/hpc/v1/job"

	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
)

type API func(option ...Option) (*job.SystemDeleteResponse, error)

func New(hc *xhttp.Client) API {
	return func(opts ...Option) (*job.SystemDeleteResponse, error) {
		req := NewRequest(opts)
		resolver := hc.Prepare(xhttp.NewDeleteRequestBuilder().
			URI(fmt.Sprintf("/system/jobs/%s", req.JobID)))

		ret := new(job.SystemDeleteResponse)
		err := resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

func NewRequest(opts []Option) *job.SystemDeleteRequest {
	req := new(job.SystemDeleteRequest)
	for _, opt := range opts {
		opt(req)
	}

	return req
}
