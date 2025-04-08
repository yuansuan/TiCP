package moduleconfig

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	moduleconfigdelete "github.com/yuansuan/ticp/common/openapi-go/apiv1/license/moduleconfig/delete"
	moduleconfig "github.com/yuansuan/ticp/common/project-root-api/license/v1/module_config"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi"
	"github.com/yuansuan/ticp/PSP/psp/pkg/tracelog"
)

func Delete(api *openapi.OpenAPI, moduleConfigId string) (*moduleconfig.DeleteModuleConfigResponse, error) {
	options := []moduleconfigdelete.Option{
		api.Client.License.DeleteModuleConfig.Id(moduleConfigId),
	}

	resp, err := api.Client.License.DeleteModuleConfig(options...)
	if err != nil {
		return nil, errors.Wrap(err, "openapi edit license manager error")
	}
	tracelog.Info(context.Background(), fmt.Sprintf("openapi edit license manager error, moduleConfigId:[%v]", moduleConfigId))

	return resp, err
}
