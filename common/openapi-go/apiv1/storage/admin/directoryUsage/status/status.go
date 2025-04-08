package status

import (
	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
	"github.com/yuansuan/ticp/common/project-root-api/storage/v20230530/directoryUsage/status"
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
			URI("/system/storage/directoryUsage/status").
			AddQuery("DirectoryUsageTaskID", utils.Stringify(req.DirectoryUsageTaskID)))

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

func (api API) DirectoryUsageTaskID(directoryUsageTaskID string) Option {
	return func(req *status.Request) error {
		req.DirectoryUsageTaskID = directoryUsageTaskID
		return nil
	}
}
