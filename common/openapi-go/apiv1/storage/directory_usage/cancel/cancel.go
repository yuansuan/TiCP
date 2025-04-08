package cancel

import (
	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
	"github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/directoryUsage/cancel"
	"net/http"
)

type API func(options ...Option) (*cancel.Response, error)

func New(hc *xhttp.Client) API {
	return func(options ...Option) (*cancel.Response, error) {
		req, err := NewRequest(options)
		if err != nil {
			return nil, err
		}
		resolver := hc.Prepare(xhttp.NewPostRequestBuilder().
			URI("/api/storage/directoryUsage/cancel").
			Json(req))

		ret := new(cancel.Response)
		err = resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

func NewRequest(options []Option) (*cancel.Request, error) {
	req := &cancel.Request{}

	for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}

type Option func(req *cancel.Request) error

func (api API) DirectoryUsageTaskID(directoryUsageTaskID string) Option {
	return func(req *cancel.Request) error {
		req.DirectoryUsageTaskID = directoryUsageTaskID
		return nil
	}
}
