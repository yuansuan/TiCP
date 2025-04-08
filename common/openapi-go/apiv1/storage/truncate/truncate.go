package truncate

import (
	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
	"github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/truncate"
	"net/http"
)

type API func(options ...Option) (*truncate.Response, error)

func New(hc *xhttp.Client) API {
	return func(options ...Option) (*truncate.Response, error) {
		req, err := NewRequest(options)
		if err != nil {
			return nil, err
		}
		resolver := hc.Prepare(xhttp.NewPostRequestBuilder().
			URI("/api/storage/truncate").
			Json(req))

		ret := new(truncate.Response)
		err = resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

func NewRequest(options []Option) (*truncate.Request, error) {
	req := &truncate.Request{}

	for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}

type Option func(req *truncate.Request) error

func (api API) Path(path string) Option {
	return func(req *truncate.Request) error {
		req.Path = path
		return nil
	}
}

func (api API) Size(size int64) Option {
	return func(req *truncate.Request) error {
		req.Size = size
		return nil
	}
}

func (api API) CreateIfNotExists(createIfNotExists bool) Option {
	return func(req *truncate.Request) error {
		req.CreateIfNotExists = createIfNotExists
		return nil
	}
}
