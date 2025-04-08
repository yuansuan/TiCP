package realpath

import (
	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
	"github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/realpath"
	"net/http"
)

type API func(options ...Option) (*realpath.Response, error)

func New(hc *xhttp.Client) API {
	return func(options ...Option) (*realpath.Response, error) {
		req, err := NewRequest(options)
		if err != nil {
			return nil, err
		}
		resolver := hc.Prepare(xhttp.NewRequestBuilder().
			URI("/system/storage/realpath").
			AddQuery("RelativePath", utils.Stringify(req.RelativePath)),
		)

		ret := new(realpath.Response)
		err = resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

func NewRequest(options []Option) (*realpath.Request, error) {
	req := &realpath.Request{}

	for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}

type Option func(req *realpath.Request) error

func (api API) RelativePath(path string) Option {
	return func(req *realpath.Request) error {
		req.RelativePath = path
		return nil
	}
}
