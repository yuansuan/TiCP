package mv

import (
	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
	"github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/mv"
	"net/http"
)

type API func(options ...Option) (*mv.Response, error)

func New(hc *xhttp.Client) API {
	return func(options ...Option) (*mv.Response, error) {
		req, err := NewRequest(options)
		if err != nil {
			return nil, err
		}
		resolver := hc.Prepare(xhttp.NewPostRequestBuilder().
			URI("/api/storage/mv").
			Json(req))

		ret := new(mv.Response)
		err = resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

func NewRequest(options []Option) (*mv.Request, error) {
	req := &mv.Request{}

	for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}

type Option func(req *mv.Request) error

func (api API) Src(path string) Option {
	return func(req *mv.Request) error {
		req.SrcPath = path
		return nil
	}
}

func (api API) Dest(path string) Option {
	return func(req *mv.Request) error {
		req.DestPath = path
		return nil
	}
}
