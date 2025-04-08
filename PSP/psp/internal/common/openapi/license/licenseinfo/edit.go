package licenseinfo

import (
	"github.com/pkg/errors"
	licenseinfoput "github.com/yuansuan/ticp/common/openapi-go/apiv1/license/licenseinfo/put"
	licenseinfo "github.com/yuansuan/ticp/common/project-root-api/license/v1/license_info"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi"
)

func Edit(api *openapi.OpenAPI, req *licenseinfo.PutLicenseInfoRequest) (*licenseinfo.PutLicenseInfoResponse, error) {
	options := []licenseinfoput.Option{
		api.Client.License.PutLicenseInfo.Id(req.Id),
		api.Client.License.PutLicenseInfo.ManagerId(req.ManagerId),
		api.Client.License.PutLicenseInfo.Provider(req.Provider),
		api.Client.License.PutLicenseInfo.MacAddr(req.MacAddr),
		api.Client.License.PutLicenseInfo.ToolPath(req.ToolPath),
		api.Client.License.PutLicenseInfo.LicenseUrl(req.LicenseUrl),
		api.Client.License.PutLicenseInfo.Port(req.Port),
		api.Client.License.PutLicenseInfo.LicenseNum(req.LicenseNum),
		api.Client.License.PutLicenseInfo.Weight(req.Weight),
		api.Client.License.PutLicenseInfo.BeginTime(req.BeginTime),
		api.Client.License.PutLicenseInfo.EndTime(req.EndTime),
		api.Client.License.PutLicenseInfo.Auth(req.Auth),
		api.Client.License.PutLicenseInfo.LicenseEnvVar(req.LicenseEnvVar),
		api.Client.License.PutLicenseInfo.LicenseType(req.LicenseType),
		api.Client.License.PutLicenseInfo.HpcEndpoint(req.HpcEndpoint),
		api.Client.License.PutLicenseInfo.AllowableHpcEndpoints(req.AllowableHpcEndpoints),
		api.Client.License.PutLicenseInfo.CollectorType(req.CollectorType),
	}

	resp, err := api.Client.License.PutLicenseInfo(options...)
	if err != nil {
		return nil, errors.Wrap(err, "openapi edit license info error")
	}

	return resp, err

}
