package licensemanager

import (
	"github.com/pkg/errors"
	licmanageradd "github.com/yuansuan/ticp/common/openapi-go/apiv1/license/licensemanager/add"
	licmanager "github.com/yuansuan/ticp/common/project-root-api/license/v1/license_manager"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi"
)

func Add(api *openapi.OpenAPI, req *licmanager.AddLicManagerRequest) (*licmanager.AddLicManagerResponse, error) {
	options := []licmanageradd.Option{
		api.Client.License.AddLicenseManager.AppType(req.AppType),
		api.Client.License.AddLicenseManager.Os(req.Os),
		api.Client.License.AddLicenseManager.Desc(req.Desc),
		api.Client.License.AddLicenseManager.ComputeRule(req.ComputeRule),
	}

	resp, err := api.Client.License.AddLicenseManager(options...)
	if err != nil {
		return nil, errors.Wrap(err, "openapi add license manager error")
	}

	return resp, err
}
