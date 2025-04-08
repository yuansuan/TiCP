package mkdir

import (
	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
	"github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/mkdir"
	"net/http"
)

type API func(options ...Option) (*mkdir.Response, error)

func New(hc *xhttp.Client) API {
	return func(options ...Option) (*mkdir.Response, error) {
		req, err := NewRequest(options)
		if err != nil {
			return nil, err
		}
		resolver := hc.Prepare(xhttp.NewPostRequestBuilder().
			URI("/system/storage/mkdir").
			Json(req))

		ret := new(mkdir.Response)
		err = resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

func NewRequest(options []Option) (*mkdir.Request, error) {
	req := &mkdir.Request{}

	for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}

type Option func(req *mkdir.Request) error

func (api API) Path(path string) Option {
	return func(req *mkdir.Request) error {
		req.Path = path
		return nil
	}
}

func (api API) IgnoreExist(ignoreExist bool) Option {
	return func(req *mkdir.Request) error {
		req.IgnoreExist = ignoreExist
		return nil
	}
}
