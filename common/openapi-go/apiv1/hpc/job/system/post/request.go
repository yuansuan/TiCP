package post

import (
	"net/http"

	"github.com/yuansuan/ticp/common/project-root-api/hpc/v1/job"

	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
)

type API func(option ...Option) (*job.SystemPostResponse, error)

func New(hc *xhttp.Client) API {
	return func(opts ...Option) (*job.SystemPostResponse, error) {
		req := NewRequest(opts)
		builder := xhttp.NewPostRequestBuilder().
			URI("/system/jobs").
			Json(req)
		if req.IdempotentID != "" {
			builder.AddQuery("IdempotentID", req.IdempotentID)
		}
		resolver := hc.Prepare(builder)
		ret := new(job.SystemPostResponse)
		err := resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

func NewRequest(opts []Option) *job.SystemPostRequest {
	req := new(job.SystemPostRequest)
	for _, opt := range opts {
		opt(req)
	}

	return req
}
