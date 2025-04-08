package licensemanager

import (
	"github.com/pkg/errors"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	licmanagerput "github.com/yuansuan/ticp/common/openapi-go/apiv1/license/licensemanager/put"
	licmanager "github.com/yuansuan/ticp/common/project-root-api/license/v1/license_manager"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi"
)

func Edit(api *openapi.OpenAPI, req *licmanager.PutLicManagerRequest) (*licmanager.PutLicManagerResponse, error) {
	options := []licmanagerput.Option{
		api.Client.License.PutLicenseManager.Id(req.Id),
		api.Client.License.PutLicenseManager.AppType(req.AppType),
		api.Client.License.PutLicenseManager.Os(req.Os),
		api.Client.License.PutLicenseManager.Desc(req.Desc),
		api.Client.License.PutLicenseManager.Status(req.Status),
		api.Client.License.PutLicenseManager.ComputeRule(req.ComputeRule),
	}

	resp, err := api.Client.License.PutLicenseManager(options...)
	if err != nil {
		return nil, errors.Wrap(err, "openapi edit license manager error")
	}

	if resp != nil {
		logging.Default().Debugf("openapi edit license manager resp: %v", resp)
	}

	return resp, err
}
