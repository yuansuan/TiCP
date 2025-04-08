package start

import (
	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
	"github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/directoryUsage/start"
	"net/http"
)

type API func(options ...Option) (*start.Response, error)

func New(hc *xhttp.Client) API {
	return func(options ...Option) (*start.Response, error) {
		req, err := NewRequest(options)
		if err != nil {
			return nil, err
		}
		resolver := hc.Prepare(xhttp.NewPostRequestBuilder().
			URI("/system/storage/directoryUsage/start").
			Json(req))

		ret := new(start.Response)
		err = resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

func NewRequest(options []Option) (*start.Request, error) {
	req := &start.Request{}

	for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}

type Option func(req *start.Request) error

func (api API) Path(path string) Option {
	return func(req *start.Request) error {
		req.Path = path
		return nil
	}
}
