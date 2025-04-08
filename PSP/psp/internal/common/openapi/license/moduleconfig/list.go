package moduleconfig

import (
	moduleconfigList "github.com/yuansuan/ticp/common/openapi-go/apiv1/license/moduleconfig/list"
	moduleconfig "github.com/yuansuan/ticp/common/project-root-api/license/v1/module_config"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi"
)

func ListModuleConfig(api *openapi.OpenAPI, licenseId string) (*moduleconfig.ListModuleConfigResponse, error) {

	options := []moduleconfigList.Option{
		api.Client.License.ListModuleConfig.LicenseId(licenseId),
	}

	response, err := api.Client.License.ListModuleConfig(options...)
	return response, err
}
