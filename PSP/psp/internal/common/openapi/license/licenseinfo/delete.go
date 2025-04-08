package licenseinfo

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	licinfodelete "github.com/yuansuan/ticp/common/openapi-go/apiv1/license/licenseinfo/delete"
	licinfo "github.com/yuansuan/ticp/common/project-root-api/license/v1/license_info"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi"
	"github.com/yuansuan/ticp/PSP/psp/pkg/tracelog"
)

func Delete(api *openapi.OpenAPI, licManagerId string) (*licinfo.DeleteLicenseInfoResponse, error) {
	options := []licinfodelete.Option{
		api.Client.License.DeleteLicenseInfo.Id(licManagerId),
	}

	resp, err := api.Client.License.DeleteLicenseInfo(options...)
	if err != nil {
		return nil, errors.Wrap(err, "openapi delete license info error")
	}
	tracelog.Info(context.Background(), fmt.Sprintf("openapi delete license info error, managerId:[%v]", licManagerId))

	return resp, err
}
