package init

import (
	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
	uploadInit "github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/upload/init"
	"net/http"
)

type API func(options ...Option) (*uploadInit.Response, error)

func New(hc *xhttp.Client) API {
	return func(options ...Option) (*uploadInit.Response, error) {
		req, err := NewRequest(options)
		if err != nil {
			return nil, err
		}
		resolver := hc.Prepare(xhttp.NewPostRequestBuilder().
			URI("/system/storage/upload/init").
			Json(req))

		ret := new(uploadInit.Response)
		err = resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

func NewRequest(options []Option) (*uploadInit.Request, error) {
	req := &uploadInit.Request{}

	for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}

type Option func(req *uploadInit.Request) error

func (api API) Path(path string) Option {
	return func(req *uploadInit.Request) error {
		req.Path = path
		return nil
	}
}

func (api API) Size(size int64) Option {
	return func(req *uploadInit.Request) error {
		req.Size = size
		return nil
	}
}

func (api API) Overwrite(overwrite bool) Option {
	return func(req *uploadInit.Request) error {
		req.Overwrite = overwrite
		return nil
	}
}
