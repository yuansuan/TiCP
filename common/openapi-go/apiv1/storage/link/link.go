package link

import (
	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
	"github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/link"
	"net/http"
)

type API func(options ...Option) (*link.Response, error)

func New(hc *xhttp.Client) API {
	return func(options ...Option) (*link.Response, error) {
		req, err := NewRequest(options)
		if err != nil {
			return nil, err
		}
		resolver := hc.Prepare(xhttp.NewPostRequestBuilder().
			URI("/api/storage/link").
			Json(req))

		ret := new(link.Response)
		err = resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

func NewRequest(options []Option) (*link.Request, error) {
	req := &link.Request{}

	for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}

type Option func(req *link.Request) error

func (api API) Src(path string) Option {
	return func(req *link.Request) error {
		req.SrcPath = path
		return nil
	}
}

func (api API) Dest(path string) Option {
	return func(req *link.Request) error {
		req.DestPath = path
		return nil
	}
}
