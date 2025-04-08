package remoteapp

import (
	v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
)

type ApiGetRequest struct {
	SessionId     *string `uri:"SessionId" required:"true"`
	RemoteAppName *string `uri:"RemoteAppName" required:"true"`
}

type ApiGetResponse struct {
	v20230530.Response

	Data *ApiGetResponseData `json:"Data,omitempty"`
}

type ApiGetResponseData struct {
	Url string `json:"Url,omitempty"`
}

type AdminPostRequest struct {
	SoftwareId *string `json:"SoftwareId,omitempty" required:"true"`
	Desc       *string `json:"Desc,omitempty" required:"false"`
	Name       *string `json:"Name,omitempty" required:"true"`
	Dir        *string `json:"Dir,omitempty" required:"false"`
	Args       *string `json:"Args,omitempty" required:"false"`
	Logo       *string `json:"Logo,omitempty" required:"false"`
	DisableGfx *bool   `json:"DisableGfx,omitempty" required:"false"`
	LoginUser  *string `json:"LoginUser,omitempty" required:"false"` // 远程应用登陆用户
}

type AdminPostResponse struct {
	v20230530.Response

	Data *AdminPostResponseData `json:"Data,omitempty"`
}

type AdminPostResponseData struct {
	Id string `json:"Id"`
}

type AdminPutRequest struct {
	RemoteAppId *string `uri:"RemoteAppId" require:"true"`
	AdminPostRequest
}

type AdminPutResponse struct {
	v20230530.Response
}

type AdminPatchRequest struct {
	RemoteAppId *string `uri:"RemoteAppId" required:"true"`
	SoftwareId  *string `json:"SoftwareId,omitempty" required:"false"`
	Desc        *string `json:"Desc,omitempty" required:"false"`
	Name        *string `json:"Name,omitempty" required:"false"`
	Dir         *string `json:"Dir,omitempty" required:"false"`
	Args        *string `json:"Args,omitempty" required:"false"`
	Logo        *string `json:"Logo,omitempty" required:"false"`
	DisableGfx  *bool   `json:"DisableGfx,omitempty" required:"false"`
	LoginUser   *string `json:"LoginUser,omitempty" required:"false"`
}

type AdminPatchResponse struct {
	v20230530.Response
}

type AdminDeleteRequest struct {
	RemoteAppId *string `uri:"RemoteAppId" required:"true"`
}

type AdminDeleteResponse struct {
	v20230530.Response
}
