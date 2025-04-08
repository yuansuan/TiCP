package batchadd

import (
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
	"github.com/yuansuan/ticp/common/project-root-api/common"
	moduleconfig "github.com/yuansuan/ticp/common/project-root-api/license/v1/module_config"
	"net/http"
)

type API func(options ...Option) (*moduleconfig.BatchAddModuleConfigResponse, error)

func New(hc *xhttp.Client) API {
	return func(options ...Option) (*moduleconfig.BatchAddModuleConfigResponse, error) {
		req, err := NewRequest(options)
		if err != nil {
			return nil, err
		}
		resolver := hc.Prepare(xhttp.NewPostRequestBuilder().
			URI("/admin/moduleConfigs/batch").
			Json(req))

		ret := new(moduleconfig.BatchAddModuleConfigResponse)
		err = resolver.Resolve(func(resp *http.Response) error {
			return common.ParseResp(resp, ret)
		})

		return ret, err
	}
}

func NewRequest(options []Option) (*moduleconfig.BatchAddModuleConfigRequest, error) {
	req := new(moduleconfig.BatchAddModuleConfigRequest)

	for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}

type Option func(req *moduleconfig.BatchAddModuleConfigRequest) error

func (api API) LicenseId(licenseID string) Option {
	return func(req *moduleconfig.BatchAddModuleConfigRequest) error {
		req.LicenseId = licenseID
		return nil
	}
}

func (api API) ModuleConfigs(moduleConfigs []*moduleconfig.AddModuleConfigRequest) Option {
	return func(req *moduleconfig.BatchAddModuleConfigRequest) error {
		req.ModuleConfigs = moduleConfigs
		return nil
	}
}
