package moduleconfig

import (
	"github.com/pkg/errors"
	moduleconfigedit "github.com/yuansuan/ticp/common/openapi-go/apiv1/license/moduleconfig/put"
	moduleconfig "github.com/yuansuan/ticp/common/project-root-api/license/v1/module_config"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi"
)

func Edit(api *openapi.OpenAPI, req *moduleconfig.PutModuleConfigRequest) (*moduleconfig.PutModuleConfigResponse, error) {
	options := []moduleconfigedit.Option{
		api.Client.License.PutModuleConfig.Id(req.Id),
		api.Client.License.PutModuleConfig.ModuleName(req.ModuleName),
		api.Client.License.PutModuleConfig.Total(req.Total),
	}

	resp, err := api.Client.License.PutModuleConfig(options...)
	if err != nil {
		return nil, errors.Wrap(err, "openapi edit license module config error")
	}

	return resp, err
}
