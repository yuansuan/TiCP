package openapivisual

import (
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/common/project-root-api/cloud_app/v1/software"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi"
	visualdto "github.com/yuansuan/ticp/PSP/psp/internal/visual/dto"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/strutil"
)

func AddSoftware(api *openapi.OpenAPI, software *visualdto.Software, zone string) (*software.AdminPostResponse, error) {
	if strutil.IsEmpty(software.Name) {
		return nil, ErrSoftwareNameEmpty
	}
	if strutil.IsEmpty(software.Platform) {
		return nil, ErrSoftwarePlatformEmpty
	}
	if strutil.IsEmpty(software.ImageID) {
		return nil, ErrSoftwareImageIDEmpty
	}

	client := api.Client.CloudApp.Software.Admin
	response, err := api.Client.CloudApp.Software.Admin.Add(
		client.Add.Name(software.Name),
		client.Add.Desc(software.Desc),
		client.Add.Platform(software.Platform),
		client.Add.ImageId(software.ImageID),
		client.Add.InitScript(software.InitScript),
		//client.Add.Icon(software.Icon),
		client.Add.GpuDesired(software.GPUDesired),
		client.Add.Zone(zone),
	)

	if response != nil {
		logging.Default().Debugf("openapi add software request id: [%v], req: [%+v], zone: [%v]", response.RequestID, software, zone)
	}

	return response, err
}

func UpdateSoftware(api *openapi.OpenAPI, softwareID string, software *visualdto.Software, zone string) (*software.AdminPutResponse, error) {
	if strutil.IsEmpty(softwareID) {
		return nil, ErrSoftwareIDEmpty
	}
	if strutil.IsEmpty(software.Name) {
		return nil, ErrSoftwareNameEmpty
	}
	if strutil.IsEmpty(software.Platform) {
		return nil, ErrSoftwarePlatformEmpty
	}
	if strutil.IsEmpty(software.ImageID) {
		return nil, ErrSoftwareImageIDEmpty
	}

	client := api.Client.CloudApp.Software.Admin
	response, err := api.Client.CloudApp.Software.Admin.Put(
		client.Put.Id(softwareID),
		client.Put.Name(software.Name),
		client.Put.Desc(software.Desc),
		client.Put.Platform(software.Platform),
		client.Put.ImageId(software.ImageID),
		client.Put.InitScript(software.InitScript),
		//client.Put.Icon(software.Icon),
		client.Put.GpuDesired(software.GPUDesired),
		client.Put.Zone(zone),
	)

	if response != nil {
		logging.Default().Debugf("openapi update software request id: [%v], softwareID: [%v], req: [%+v], zone: [%v]",
			response.RequestID, softwareID, software, zone)
	}

	return response, err
}

func DeleteSoftware(api *openapi.OpenAPI, softwareID string) (*software.AdminDeleteResponse, error) {
	if strutil.IsEmpty(softwareID) {
		return nil, ErrSoftwareIDEmpty
	}

	client := api.Client.CloudApp.Software.Admin
	response, err := api.Client.CloudApp.Software.Admin.Delete(
		client.Delete.Id(softwareID),
	)

	if response != nil {
		logging.Default().Debugf("openapi delete software request id: [%v], softwareID: [%v]", response.RequestID, softwareID)
	}

	return response, err
}
