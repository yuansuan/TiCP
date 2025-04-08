package licensemanager

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	licmanagerdelete "github.com/yuansuan/ticp/common/openapi-go/apiv1/license/licensemanager/delete"
	licmanager "github.com/yuansuan/ticp/common/project-root-api/license/v1/license_manager"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi"
	"github.com/yuansuan/ticp/PSP/psp/pkg/tracelog"
)

func Delete(api *openapi.OpenAPI, licManagerId string) (*licmanager.DeleteLicManagerResponse, error) {
	options := []licmanagerdelete.Option{
		api.Client.License.DeleteLicenseManager.Id(licManagerId),
	}

	resp, err := api.Client.License.DeleteLicenseManager(options...)
	if err != nil {
		return nil, errors.Wrap(err, "openapi delete license manager error")
	}
	tracelog.Info(context.Background(), fmt.Sprintf("openapi delete license manager error, managerId:[%v]", licManagerId))

	return resp, err
}
