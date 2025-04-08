package delete

import (
	"github.com/yuansuan/ticp/common/openapi-go/common"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
	moduleconfig "github.com/yuansuan/ticp/common/project-root-api/license/v1/module_config"
	"net/http"
)

type API func(options ...Option) (*moduleconfig.DeleteModuleConfigResponse, error)

func New(hc *xhttp.Client) API {
	return func(options ...Option) (*moduleconfig.DeleteModuleConfigResponse, error) {
		req, err := NewRequest(options)
		if err != nil {
			return nil, err
		}
		resolver := hc.Prepare(xhttp.NewDeleteRequestBuilder().
			URI("/admin/moduleConfigs/" + string(*req)))

		ret := new(moduleconfig.DeleteModuleConfigResponse)
		err = resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

func NewRequest(options []Option) (*moduleconfig.DeleteModuleConfigRequest, error) {
	req := new(moduleconfig.DeleteModuleConfigRequest)

	for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}

type Option func(req *moduleconfig.DeleteModuleConfigRequest) error

func (api API) Id(id string) Option {
	return func(req *moduleconfig.DeleteModuleConfigRequest) error {
		*req = moduleconfig.DeleteModuleConfigRequest(id)
		return nil
	}
}
