package moduleconfig

import (
	"github.com/pkg/errors"
	moduleconfigadd "github.com/yuansuan/ticp/common/openapi-go/apiv1/license/moduleconfig/add"
	moduleconfig "github.com/yuansuan/ticp/common/project-root-api/license/v1/module_config"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi"
)

func Add(api *openapi.OpenAPI, req *moduleconfig.AddModuleConfigRequest) (*moduleconfig.AddModuleConfigResponse, error) {
	options := []moduleconfigadd.Option{
		api.Client.License.AddModuleConfig.LicenseId(req.LicenseId),
		api.Client.License.AddModuleConfig.ModuleName(req.ModuleName),
		api.Client.License.AddModuleConfig.Total(req.Total),
	}

	resp, err := api.Client.License.AddModuleConfig(options...)
	if err != nil {
		return nil, errors.Wrap(err, "openapi add license manager error")
	}

	return resp, err
}
