package copyRange

import (
	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
	"github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/copyRange"
	"net/http"
)

type API func(options ...Option) (*copyRange.Response, error)

func New(hc *xhttp.Client) API {
	return func(options ...Option) (*copyRange.Response, error) {
		req, err := NewRequest(options)
		if err != nil {
			return nil, err
		}
		resolver := hc.Prepare(xhttp.NewPostRequestBuilder().
			URI("/api/storage/copyRange").
			Json(req))

		ret := new(copyRange.Response)
		err = resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

func NewRequest(options []Option) (*copyRange.Request, error) {
	req := &copyRange.Request{}

	for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}

type Option func(req *copyRange.Request) error

func (api API) Src(path string) Option {
	return func(req *copyRange.Request) error {
		req.SrcPath = path
		return nil
	}
}

func (api API) Dest(path string) Option {
	return func(req *copyRange.Request) error {
		req.DestPath = path
		return nil
	}
}

func (api API) SrcOffset(srcOffset int64) Option {
	return func(req *copyRange.Request) error {
		req.SrcOffset = srcOffset
		return nil
	}
}

func (api API) DestOffset(destOffset int64) Option {
	return func(req *copyRange.Request) error {
		req.DestOffset = destOffset
		return nil
	}
}

func (api API) Length(length int64) Option {
	return func(req *copyRange.Request) error {
		req.Length = length
		return nil
	}
}
