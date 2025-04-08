package stat

import (
	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"

	"github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/stat"
	"net/http"
)

type API func(options ...Option) (*stat.Response, error)

func New(hc *xhttp.Client) API {
	return func(options ...Option) (*stat.Response, error) {
		req, err := NewRequest(options)
		if err != nil {
			return nil, err
		}
		resolver := hc.Prepare(xhttp.NewRequestBuilder().
			URI("/api/storage/stat").
			AddQuery("Path", utils.Stringify(req.Path)))

		ret := new(stat.Response)
		err = resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

func NewRequest(options []Option) (*stat.Request, error) {
	req := &stat.Request{}

	for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}

type Option func(req *stat.Request) error

func (api API) Path(path string) Option {
	return func(req *stat.Request) error {
		req.Path = path
		return nil
	}
}
