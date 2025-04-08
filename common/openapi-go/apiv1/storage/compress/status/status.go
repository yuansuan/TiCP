package status

import (
	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
	"github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/compress/status"
	"net/http"
)

type API func(options ...Option) (*status.Response, error)

func New(hc *xhttp.Client) API {
	return func(options ...Option) (*status.Response, error) {
		req, err := NewRequest(options)
		if err != nil {
			return nil, err
		}
		resolver := hc.Prepare(xhttp.NewRequestBuilder().
			URI("/api/storage/compress/status").
			AddQuery("CompressID", utils.Stringify(req.CompressID)))

		ret := new(status.Response)
		err = resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

func NewRequest(options []Option) (*status.Request, error) {
	req := &status.Request{}

	for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}

type Option func(req *status.Request) error

func (api API) CompressID(compressID string) Option {
	return func(req *status.Request) error {
		req.CompressID = compressID
		return nil
	}
}
