package create

import (
	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
	"github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/create"
	"net/http"
)

type API func(options ...Option) (*create.Response, error)

func New(hc *xhttp.Client) API {
	return func(options ...Option) (*create.Response, error) {
		req, err := NewRequest(options)
		if err != nil {
			return nil, err
		}
		resolver := hc.Prepare(xhttp.NewPostRequestBuilder().
			URI("/api/storage/create").
			Json(req))

		ret := new(create.Response)
		err = resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

func NewRequest(options []Option) (*create.Request, error) {
	req := &create.Request{}

	for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}

type Option func(req *create.Request) error

func (api API) Path(path string) Option {
	return func(req *create.Request) error {
		req.Path = path
		return nil
	}
}

func (api API) Size(size int64) Option {
	return func(req *create.Request) error {
		req.Size = size
		return nil
	}
}

func (api API) Overwrite(overwrite bool) Option {
	return func(req *create.Request) error {
		req.Overwrite = overwrite
		return nil
	}
}
