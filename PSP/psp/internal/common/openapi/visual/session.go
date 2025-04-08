package openapivisual

import (
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/common/project-root-api/cloud_app/v1/session"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/strutil"
)

func ListSession(api *openapi.OpenAPI, sessionIDsStr string, pageIndex, pageSize int) (*session.ApiListResponse, error) {
	if strutil.IsEmpty(sessionIDsStr) {
		return nil, ErrSessionIDsStrEmpty
	}

	client := api.Client.CloudApp.Session.User
	response, err := client.List(
		client.List.SessionIds(sessionIDsStr),
		client.List.PageOffset(pageIndex),
		client.List.PageSize(pageSize),
	)

	if response != nil {
		logging.Default().Debugf("openapi list session request id: [%v], sessionIDsStr: [%v] pageIndex: [%v], "+
			"pageSize: [%v]", response.RequestID, sessionIDsStr, pageIndex, pageSize)
	}

	return response, err
}

func GetSession(api *openapi.OpenAPI, sessionID string) (*session.ApiGetResponse, error) {
	if strutil.IsEmpty(sessionID) {
		return nil, ErrSessionIDEmpty
	}

	client := api.Client.CloudApp.Session.User
	response, err := client.Get(
		client.Get.Id(sessionID),
	)

	if response != nil {
		logging.Default().Debugf("openapi get session request id: [%v], sessionID: [%v]", response.RequestID, sessionID)
	}

	return response, err
}

func StartSession(api *openapi.OpenAPI, hardwareID, softwareID string, mountPaths map[string]string) (*session.ApiPostResponse, error) {
	if strutil.IsEmpty(hardwareID) {
		return nil, ErrHardwareIDEmpty
	}
	if strutil.IsEmpty(softwareID) {
		return nil, ErrSoftwareIDEmpty
	}

	client := api.Client.CloudApp.Session.User
	response, err := client.Start(
		client.Start.HardwareId(hardwareID),
		client.Start.SoftwareId(softwareID),
		client.Start.MountPaths(mountPaths),
	)

	if response != nil {
		logging.Default().Debugf("openapi start session request id: [%v], hardwareID: [%v], softwareID: [%v], "+
			"mountPaths: [%v]", response.RequestID, hardwareID, softwareID, mountPaths)
	}

	return response, err
}

func UserCloseSession(api *openapi.OpenAPI, sessionID string) (*session.ApiCloseResponse, error) {
	if strutil.IsEmpty(sessionID) {
		return nil, ErrSessionIDEmpty
	}

	client := api.Client.CloudApp.Session.User
	response, err := client.Close(
		client.Close.Id(sessionID),
	)

	if response != nil {
		logging.Default().Debugf("openapi user close session request id: [%v], sessionID: [%v]", response.RequestID, sessionID)
	}

	return response, err
}

func AdminPowerOffSession(api *openapi.OpenAPI, sessionID string) (*session.PowerOffResponse, error) {
	if strutil.IsEmpty(sessionID) {
		return nil, ErrSessionIDEmpty
	}

	client := api.Client.CloudApp.Session.Admin
	response, err := client.PowerOff(
		client.PowerOff.Id(sessionID),
	)

	if response != nil {
		logging.Default().Debugf("openapi admin power off session request id: [%v], sessionID: [%v]", response.RequestID, sessionID)
	}

	return response, err
}

func AdminPowerOnSession(api *openapi.OpenAPI, sessionID string) (*session.PowerOnResponse, error) {
	if strutil.IsEmpty(sessionID) {
		return nil, ErrSessionIDEmpty
	}

	client := api.Client.CloudApp.Session.Admin
	response, err := client.PowerOn(
		client.PowerOn.Id(sessionID),
	)

	if response != nil {
		logging.Default().Debugf("openapi admin power on session request id: [%v], sessionID: [%v]", response.RequestID, sessionID)
	}

	return response, err
}

func AdminCloseSession(api *openapi.OpenAPI, sessionID, existReason string) (*session.AdminCloseResponse, error) {
	if strutil.IsEmpty(sessionID) {
		return nil, ErrSessionIDEmpty
	}
	if strutil.IsEmpty(existReason) {
		return nil, ErrExistReasonEmpty
	}

	client := api.Client.CloudApp.Session.Admin
	response, err := client.Close(
		client.Close.Id(sessionID),
		client.Close.Reason(existReason),
	)

	if response != nil {
		logging.Default().Debugf("openapi admin close session request id: [%v], sessionID: [%v], existReason: [%v]",
			response.RequestID, sessionID, existReason)
	}

	return response, err
}

func AdminRebootSession(api *openapi.OpenAPI, sessionID string) (*session.RebootResponse, error) {
	if strutil.IsEmpty(sessionID) {
		return nil, ErrSessionIDEmpty
	}

	client := api.Client.CloudApp.Session.Admin
	response, err := client.Reboot(
		client.Reboot.Id(sessionID),
	)

	if response != nil {
		logging.Default().Debugf("openapi admin reboot session request id: [%v], sessionID: [%v],",
			response.RequestID, sessionID)
	}

	return response, err
}

func ReadySession(api *openapi.OpenAPI, sessionID string) (*session.ApiReadyResponse, error) {
	if strutil.IsEmpty(sessionID) {
		return nil, ErrSessionIDEmpty
	}

	client := api.Client.CloudApp.Session.User
	response, err := client.Ready(
		client.Ready.Id(sessionID),
	)

	if response != nil {
		logging.Default().Debugf("openapi ready session request id: [%v], sessionID: [%v]", response.RequestID, sessionID)
	}

	return response, err
}

func GetRemoteAppURL(api *openapi.OpenAPI, sessionID string) (*session.ApiGetResponse, error) {
	if strutil.IsEmpty(sessionID) {
		return nil, ErrSessionIDEmpty
	}

	client := api.Client.CloudApp.Session.User
	response, err := client.Get(
		client.Get.Id(sessionID),
	)

	if response != nil {
		logging.Default().Debugf("openapi get remote app url request id: [%v], sessionID: [%v]", response.RequestID, sessionID)
	}

	return response, err
}
