package get

import (
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
	"github.com/yuansuan/ticp/common/project-root-api/common"
	app "github.com/yuansuan/ticp/common/project-root-api/job/v1/admin/app/get"
	"net/http"
)

type API func(options ...Option) (*app.Response, error)

func New(hc *xhttp.Client) API {
	return func(options ...Option) (*app.Response, error) {
		req, err := NewRequest(options)
		if err != nil {
			return nil, err
		}
		resolver := hc.Prepare(xhttp.NewRequestBuilder().
			URI("/admin/apps/" + req.AppID))

		ret := new(app.Response)
		err = resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

func NewRequest(options []Option) (*app.Request, error) {
	req := &app.Request{}

	for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}

type Option func(req *app.Request) error

func (api API) AppID(AppID string) Option {
	return func(req *app.Request) error {
		req.AppID = AppID
		return nil
	}
}
