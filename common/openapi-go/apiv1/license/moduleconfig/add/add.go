package add

import (
	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
	moduleconfig "github.com/yuansuan/ticp/common/project-root-api/license/v1/module_config"
	"net/http"
)

type API func(options ...Option) (*moduleconfig.AddModuleConfigResponse, error)

func New(hc *xhttp.Client) API {
	return func(options ...Option) (*moduleconfig.AddModuleConfigResponse, error) {
		req, err := NewRequest(options)
		if err != nil {
			return nil, err
		}
		resolver := hc.Prepare(xhttp.NewPostRequestBuilder().
			URI("/admin/moduleConfigs").
			Json(req))

		ret := new(moduleconfig.AddModuleConfigResponse)
		err = resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}
func NewRequest(options []Option) (*moduleconfig.AddModuleConfigRequest, error) {
	req := new(moduleconfig.AddModuleConfigRequest)

	for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}

type Option func(req *moduleconfig.AddModuleConfigRequest) error

func (api API) LicenseId(id string) Option {
	return func(req *moduleconfig.AddModuleConfigRequest) error {
		req.LicenseId = id
		return nil
	}
}

func (api API) ModuleName(name string) Option {
	return func(req *moduleconfig.AddModuleConfigRequest) error {
		req.ModuleName = name
		return nil
	}
}

func (api API) Total(total int) Option {
	return func(req *moduleconfig.AddModuleConfigRequest) error {
		req.Total = total
		return nil
	}
}
