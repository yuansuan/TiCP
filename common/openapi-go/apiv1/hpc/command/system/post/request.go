package post

import (
	"net/http"

	"github.com/yuansuan/ticp/common/project-root-api/hpc/v1/command"

	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
)

type API func(option ...Option) (*command.SystemPostResponse, error)

func New(hc *xhttp.Client) API {
	return func(opts ...Option) (*command.SystemPostResponse, error) {
		req := NewRequest(opts)
		resolver := hc.Prepare(xhttp.NewPostRequestBuilder().
			URI("/system/command").
			Json(req))
		ret := new(command.SystemPostResponse)
		err := resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

func NewRequest(opts []Option) *command.SystemPostRequest {
	req := new(command.SystemPostRequest)
	for _, opt := range opts {
		opt(req)
	}

	return req
}
