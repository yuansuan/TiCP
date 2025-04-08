package copy

import (
	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
	copy2 "github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/copy"
	"net/http"
)

type API func(options ...Option) (*copy2.Response, error)

func New(hc *xhttp.Client) API {
	return func(options ...Option) (*copy2.Response, error) {
		req, err := NewRequest(options)
		if err != nil {
			return nil, err
		}
		resolver := hc.Prepare(xhttp.NewPostRequestBuilder().
			URI("/system/storage/copy").
			Json(req))

		ret := new(copy2.Response)
		err = resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

func NewRequest(options []Option) (*copy2.Request, error) {
	req := &copy2.Request{}

	for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}

type Option func(req *copy2.Request) error

func (api API) Src(path string) Option {
	return func(req *copy2.Request) error {
		req.SrcPath = path
		return nil
	}
}

func (api API) Dest(path string) Option {
	return func(req *copy2.Request) error {
		req.DestPath = path
		return nil
	}
}
