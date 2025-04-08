package openapivisual

import (
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/common/project-root-api/cloud_app/v1/remoteapp"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi"
	visualdto "github.com/yuansuan/ticp/PSP/psp/internal/visual/dto"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/strutil"
)

func AddRemoteApp(api *openapi.OpenAPI, softwareID string, remoteApp *visualdto.RemoteApp) (*remoteapp.AdminPostResponse, error) {
	if strutil.IsEmpty(softwareID) {
		return nil, ErrSoftwareIDEmpty
	}
	if strutil.IsEmpty(remoteApp.Name) {
		return nil, ErrRemoteAppNameEmpty
	}
	if strutil.IsEmpty(remoteApp.BaseURL) {
		return nil, ErrRemoteAppBaseURLEmpty
	}

	client := api.Client.CloudApp.RemoteApp.Admin
	response, err := client.Add(
		client.Add.SoftwareId(softwareID),
		client.Add.Name(remoteApp.Name),
		client.Add.Desc(remoteApp.Desc),
		client.Add.Dir(remoteApp.Dir),
		client.Add.Args(remoteApp.Args),
		client.Add.Logo(remoteApp.Logo),
		client.Add.DisableGfx(remoteApp.DisableGfx),
	)

	if response != nil {
		logging.Default().Debugf("openapi add remote app request id: [%v], softwareID: [%v] req: [%+v]", response.RequestID, softwareID, remoteApp)
	}

	return response, err
}

func UpdateRemoteApp(api *openapi.OpenAPI, remoteAppID, softwareID string, remoteApp *visualdto.RemoteApp) (*remoteapp.AdminPutResponse, error) {
	if strutil.IsEmpty(remoteAppID) {
		return nil, ErrRemoteAppIDEmpty
	}
	if strutil.IsEmpty(remoteApp.Name) {
		return nil, ErrRemoteAppNameEmpty
	}
	if strutil.IsEmpty(remoteApp.BaseURL) {
		return nil, ErrRemoteAppBaseURLEmpty
	}

	client := api.Client.CloudApp.RemoteApp.Admin
	response, err := client.Put(
		client.Put.Id(remoteAppID),
		client.Put.SoftwareId(softwareID),
		client.Put.Name(remoteApp.Name),
		client.Put.Desc(remoteApp.Desc),
		client.Put.Dir(remoteApp.Dir),
		client.Put.Args(remoteApp.Args),
		client.Put.Logo(remoteApp.Logo),
		client.Put.DisableGfx(remoteApp.DisableGfx),
	)

	if response != nil {
		logging.Default().Debugf("openapi update remote app request id: [%v], remoteAppID: [%v], softwareID: [%v] "+
			"req: [%+v]", response.RequestID, remoteAppID, softwareID, remoteApp)
	}

	return response, err
}

func DeleteRemoteApp(api *openapi.OpenAPI, remoteAppID string) (*remoteapp.AdminDeleteResponse, error) {
	if strutil.IsEmpty(remoteAppID) {
		return nil, ErrRemoteAppIDEmpty
	}

	client := api.Client.CloudApp.RemoteApp.Admin
	response, err := client.Delete(
		client.Delete.Id(remoteAppID),
	)

	if response != nil {
		logging.Default().Debugf("openapi delete remote app request id: [%v], remoteAppID: [%v]", response.RequestID, remoteAppID)
	}

	return response, err
}
