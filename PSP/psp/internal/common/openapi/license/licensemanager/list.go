package licensemanager

import (
	"github.com/pkg/errors"
	licmanager "github.com/yuansuan/ticp/common/project-root-api/license/v1/license_manager"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi"
)

func ListLicenseManage(api *openapi.OpenAPI) (*licmanager.ListLicManagerResponse, error) {
	licenseManagerResp, err := api.Client.License.ListLicenseManager()

	if err != nil {
		return nil, errors.Wrap(err, "openapi get admin job info err")
	}

	return licenseManagerResp, nil
}
