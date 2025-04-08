package licensemanager

import (
	"github.com/pkg/errors"
	licmanagerget "github.com/yuansuan/ticp/common/openapi-go/apiv1/license/licensemanager/get"
	licmanager "github.com/yuansuan/ticp/common/project-root-api/license/v1/license_manager"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi"
)

func Get(api *openapi.OpenAPI, licManagerId string) (*licmanager.GetLicManagerResponse, error) {
	options := []licmanagerget.Option{
		api.Client.License.GetLicenseManager.Id(licManagerId),
	}

	resp, err := api.Client.License.GetLicenseManager(options...)
	if err != nil {
		return nil, errors.Wrap(err, "openapi delete license manager error")
	}

	return resp, err
}
