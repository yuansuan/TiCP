package get

import (
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
	"github.com/yuansuan/ticp/common/project-root-api/common"
	"github.com/yuansuan/ticp/common/project-root-api/job/v1/admin/app/delete"
	"net/http"
)

type API func(options ...Option) (*delete.Response, error)

func New(hc *xhttp.Client) API {
	return func(options ...Option) (*delete.Response, error) {
		req, err := NewRequest(options)
		if err != nil {
			return nil, err
		}
		resolver := hc.Prepare(xhttp.NewDeleteRequestBuilder().
			URI("/admin/apps/" + req.AppID))

		ret := new(delete.Response)
		err = resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

func NewRequest(options []Option) (*delete.Request, error) {
	req := &delete.Request{}

	for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}

type Option func(req *delete.Request) error

func (api API) AppID(AppID string) Option {
	return func(req *delete.Request) error {
		req.AppID = AppID
		return nil
	}
}
