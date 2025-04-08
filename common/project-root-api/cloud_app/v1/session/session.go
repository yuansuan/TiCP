package session

import (
	v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
)

type ApiGetRequest struct {
	SessionId *string `uri:"SessionId" required:"true"`
}

type ApiGetResponse struct {
	v20230530.Response

	Data *ApiGetResponseData `json:"Data,omitempty"`
}

type ApiGetResponseData struct {
	v20230530.Session
}

type ApiListRequest struct {
	PageOffset *int    `query:"PageOffset" required:"false"`
	PageSize   *int    `query:"PageSize" required:"false"`
	Status     *string `query:"Status" required:"false"`
	SessionIds *string `query:"SessionIds" required:"false"`
	Zone       *string `query:"Zone" required:"false"`
}

type ApiListResponse struct {
	v20230530.Response

	Data *ApiListResponseData `json:"Data,omitempty"`
}

type ApiListResponseData struct {
	Sessions []*v20230530.Session `json:"Sessions,omitempty"`

	Offset     int `json:"Offset"`
	Size       int `json:"Size"`
	Total      int `json:"Total"`
	NextMarker int `json:"NextMarker"`
}

type ApiPostRequest struct {
	HardwareId   *string                 `json:"HardwareId,omitempty" required:"true"`
	SoftwareId   *string                 `json:"SoftwareId,omitempty" required:"true"`
	MountPaths   *map[string]string      `json:"MountPaths,omitempty" required:"false"`
	ChargeParams *v20230530.ChargeParams `json:"ChargeParams,omitempty" required:"false"`
	PayBy        *string                 `json:"PayBy,omitempty" required:"false"`
}

type ApiPostResponse struct {
	v20230530.Response

	Data *ApiPostResponseData `json:"Data,omitempty"`
}

type ApiPostResponseData struct {
	v20230530.Session
}

type ApiCloseRequest struct {
	SessionId *string `uri:"SessionId" required:"true"`
}

type ApiCloseResponse struct {
	v20230530.Response
}

type ApiReadyRequest struct {
	SessionId *string `uri:"SessionId" required:"true"`
}

type ApiReadyResponse struct {
	v20230530.Response

	Data *ApiReadyResponseData `json:"Data,omitempty"`
}

type ApiReadyResponseData struct {
	Ready bool `json:"Ready"`
}

type ApiDeleteRequest struct {
	SessionId *string `uri:"SessionId" required:"true"`
}

type ApiDeleteResponse struct {
	v20230530.Response
}

type AdminListRequest struct {
	PageOffset  *int    `query:"PageOffset" required:"false"`
	PageSize    *int    `query:"PageSize" required:"false"`
	Status      *string `query:"Status" required:"false"`
	SessionIds  *string `query:"SessionIds" required:"false"`
	Zone        *string `query:"Zone" required:"false"`
	UserIds     *string `query:"UserIds" required:"false"`
	WithDeleted bool    `query:"WithDeleted" required:"false"`
}

type AdminListResponse struct {
	v20230530.Response

	Data *AdminListResponseData `json:"Data,omitempty"`
}

type AdminListResponseData struct {
	Sessions []*v20230530.Session `json:"Sessions,omitempty"`

	Offset     int `json:"Offset"`
	Size       int `json:"Size"`
	Total      int `json:"Total"`
	NextMarker int `json:"NextMarker"`
}

type AdminCloseRequest struct {
	SessionId *string `uri:"SessionId" required:"true"`
	Reason    *string `json:"Reason,omitempty" required:"true"`
}

type AdminCloseResponse struct {
	v20230530.Response
}

type PowerOnRequest struct {
	SessionId *string `uri:"SessionId" required:"true"`
}

type PowerOnResponse struct {
	v20230530.Response
}

type PowerOffRequest struct {
	SessionId *string `uri:"SessionId" required:"true"`
}

type PowerOffResponse struct {
	v20230530.Response
}

type RebootRequest struct {
	SessionId *string `uri:"SessionId" required:"true"`
}

type RebootResponse struct {
	v20230530.Response
}

type ApiRestoreRequest struct {
	SessionId *string `uri:"SessionId" required:"true"`
}

type ApiRestoreResponse struct {
	v20230530.Response

	Data *ApiRestoreResponseData `json:"Data,omitempty"`
}

type ApiRestoreResponseData struct {
	v20230530.Session
}

type AdminRestoreRequest struct {
	UserId    *string `json:"UserId,omitempty" required:"true"`
	SessionId *string `uri:"SessionId" required:"true"`
}

type AdminRestoreResponse struct {
	v20230530.Response

	Data *AdminRestoreResponseData `json:"Data,omitempty"`
}

type AdminRestoreResponseData struct {
	v20230530.Session
}

type ExecScriptRequest struct {
	SessionId     *string `uri:"SessionId" required:"true"`
	ScriptRunner  *string `json:"ScriptRunner" required:"false"` // powershell, default is powershell
	ScriptContent *string `json:"ScriptContent" required:"true"` // Base64 encoded. Restricted to 65535 bytes
	WaitTillEnd   *bool   `json:"WaitTillEnd" required:"false"`
}

type ExecScriptResponse struct {
	v20230530.Response

	Data *ExecScriptResponseData `json:"Data,omitempty"` // got value only WaitTillEnd = true
}

type ExecScriptResponseData struct {
	ExitCode int    `json:"ExitCode"`
	Stdout   string `json:"Stdout"`
	Stderr   string `json:"Stderr"`
}

type MountRequest struct {
	SessionId      *string `uri:"SessionId" required:"true"`
	ShareDirectory *string `json:"ShareDirectory" required:"false"` // example: dir1/dir2 , ""
	MountPoint     *string `json:"MountPoint" required:"true"`      // example: X:   /mnt/data_example
}

type MountResponse struct {
	v20230530.Response
}

type UmountRequest struct {
	SessionId  *string `uri:"SessionId" required:"true"`
	MountPoint *string `json:"MountPoint" required:"true"` // example: X:   /mnt/data_example
}

type UmountResponse struct {
	v20230530.Response
}
