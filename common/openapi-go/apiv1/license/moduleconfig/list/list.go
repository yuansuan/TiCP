package list

import (
	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
	moduleconfig "github.com/yuansuan/ticp/common/project-root-api/license/v1/module_config"
	"net/http"
)

type API func(option ...Option) (*moduleconfig.ListModuleConfigResponse, error)

type Option func(req *moduleconfig.ListModuleConfigRequest) error

func New(hc *xhttp.Client) API {
	return func(option ...Option) (*moduleconfig.ListModuleConfigResponse, error) {
		req, err := NewRequest(option...)
		if err != nil {
			return nil, err
		}
		resolver := hc.Prepare(xhttp.NewRequestBuilder().
			URI("/admin/moduleConfigs").
			Json(req))

		ret := new(moduleconfig.ListModuleConfigResponse)
		err = resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

func NewRequest(option ...Option) (*moduleconfig.ListModuleConfigRequest, error) {
	req := new(moduleconfig.ListModuleConfigRequest)
	for _, o := range option {
		if err := o(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}

func (api API) LicenseId(id string) Option {
	return func(req *moduleconfig.ListModuleConfigRequest) error {
		req.LicenseId = id
		return nil
	}
}
