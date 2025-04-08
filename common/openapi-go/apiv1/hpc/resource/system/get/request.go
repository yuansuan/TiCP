package get

import (
	"net/http"

	"github.com/yuansuan/ticp/common/project-root-api/hpc/v1/resource"

	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
)

type API func(option ...Option) (*resource.SystemGetResponse, error)

func New(hc *xhttp.Client) API {
	return func(opts ...Option) (*resource.SystemGetResponse, error) {
		req := NewRequest(opts)
		builder := xhttp.NewRequestBuilder().URI("/system/resource")
		if req.Queue != "" {
			builder.AddQuery("Queue", req.Queue)
		}

		resolver := hc.Prepare(builder)
		ret := new(resource.SystemGetResponse)
		err := resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

func NewRequest(opts []Option) *resource.SystemGetRequest {
	req := new(resource.SystemGetRequest)
	for _, opt := range opts {
		opt(req)
	}

	return req
}
