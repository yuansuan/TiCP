package licenseinfo

import (
	"github.com/pkg/errors"
	licenseinfoadd "github.com/yuansuan/ticp/common/openapi-go/apiv1/license/licenseinfo/add"
	licenseinfo "github.com/yuansuan/ticp/common/project-root-api/license/v1/license_info"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi"
)

func Add(api *openapi.OpenAPI, req *licenseinfo.AddLicenseInfoRequest) (*licenseinfo.AddLicenseInfoResponse, error) {
	options := []licenseinfoadd.Option{
		api.Client.License.AddLicenseInfo.ManagerId(req.ManagerId),
		api.Client.License.AddLicenseInfo.Provider(req.Provider),
		api.Client.License.AddLicenseInfo.MacAddr(req.MacAddr),
		api.Client.License.AddLicenseInfo.ToolPath(req.ToolPath),
		api.Client.License.AddLicenseInfo.LicenseUrl(req.LicenseUrl),
		api.Client.License.AddLicenseInfo.Port(req.Port),
		api.Client.License.AddLicenseInfo.LicenseNum(req.LicenseNum),
		api.Client.License.AddLicenseInfo.Weight(req.Weight),
		api.Client.License.AddLicenseInfo.BeginTime(req.BeginTime),
		api.Client.License.AddLicenseInfo.EndTime(req.EndTime),
		api.Client.License.AddLicenseInfo.Auth(req.Auth),
		api.Client.License.AddLicenseInfo.LicenseEnvVar(req.LicenseEnvVar),
		api.Client.License.AddLicenseInfo.LicenseType(req.LicenseType),
		api.Client.License.AddLicenseInfo.HpcEndpoint(req.HpcEndpoint),
		api.Client.License.AddLicenseInfo.AllowableHpcEndpoints(req.AllowableHpcEndpoints),
		api.Client.License.AddLicenseInfo.CollectorType(req.CollectorType),
	}

	resp, err := api.Client.License.AddLicenseInfo(options...)
	if err != nil {
		return nil, errors.Wrap(err, "openapi add license info error")
	}

	return resp, err

}
