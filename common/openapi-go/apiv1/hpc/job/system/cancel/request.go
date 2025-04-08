package cancel

import (
	"fmt"
	"net/http"

	"github.com/yuansuan/ticp/common/project-root-api/hpc/v1/job"

	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
)

type API func(option ...Option) (*job.SystemCancelResponse, error)

func New(hc *xhttp.Client) API {
	return func(opts ...Option) (*job.SystemCancelResponse, error) {
		req := NewRequest(opts)
		resolver := hc.Prepare(xhttp.NewPostRequestBuilder().
			URI(fmt.Sprintf("/system/jobs/%s/cancel", req.JobID)).
			Json(req))

		ret := new(job.SystemCancelResponse)
		err := resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

func NewRequest(opts []Option) *job.SystemCancelRequest {
	req := new(job.SystemCancelRequest)
	for _, opt := range opts {
		opt(req)
	}

	return req
}
