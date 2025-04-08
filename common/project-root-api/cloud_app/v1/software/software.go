package software

import v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"

type APIGetRequest struct {
	SoftwareId *string `uri:"SoftwareId" required:"true"`
}

type APIGetResponse struct {
	v20230530.Response `json:",inline"`

	Data *APIGetResponseData `json:"Data,omitempty"`
}

type APIGetResponseData struct {
	v20230530.Software
}

type APIListRequest struct {
	Zone       *string `query:"Zone" required:"false"`
	Name       *string `query:"Name" required:"false"`
	Platform   *string `query:"Platform" required:"false"`
	PageSize   *int    `query:"PageSize" required:"false"`
	PageOffset *int    `query:"PageOffset" required:"false"`
}

type APIListResponse struct {
	v20230530.Response

	Data *APIListResponseData `json:"Data,omitempty"`
}

type APIListResponseData struct {
	Software []*v20230530.Software `json:"Software"`

	Offset     int `json:"Offset"`
	Size       int `json:"Size"`
	Total      int `json:"Total"`
	NextMarker int `json:"NextMarker"`
}

type AdminGetRequest struct {
	SoftwareId *string `uri:"SoftwareId" required:"true"`
}

type AdminGetResponse struct {
	v20230530.Response `json:",inline"`

	Data *AdminGetResponseData `json:"Data,omitempty"`
}

type AdminGetResponseData struct {
	v20230530.Software
}

type AdminListRequest struct {
	UserId     *string `query:"UserId" required:"false"`
	Zone       *string `query:"Zone" required:"false"`
	Name       *string `query:"Name" required:"false"`
	Platform   *string `query:"Platform" required:"false"`
	PageSize   *int    `query:"PageSize" required:"false"`
	PageOffset *int    `query:"PageOffset" required:"false"`
}

type AdminListResponse struct {
	v20230530.Response

	Data *AdminListResponseData `json:"Data,omitempty"`
}

type AdminListResponseData struct {
	Software []*v20230530.Software `json:"Software,omitempty"`

	Offset     int `json:"Offset"`
	Size       int `json:"Size"`
	Total      int `json:"Total"`
	NextMarker int `json:"NextMarker"`
}

type AdminPostRequest struct {
	Zone       *string `json:"Zone,omitempty" required:"true"`
	Name       *string `json:"Name,omitempty" required:"true"`
	Desc       *string `json:"Desc,omitempty" required:"false"`
	Icon       *string `json:"Icon,omitempty" required:"false"`
	Platform   *string `json:"Platform,omitempty"  required:"true"`
	ImageId    *string `json:"ImageId,omitempty" required:"true"`
	InitScript *string `json:"InitScript,omitempty" required:"false"`
	GpuDesired *bool   `json:"GpuDesired,omitempty" required:"false"`
}

type AdminPostResponse struct {
	v20230530.Response

	Data *AdminPostResponseData `json:"Data,omitempty"`
}

type AdminPostResponseData struct {
	SoftwareId string `json:"SoftwareId"`
}

type AdminPutRequest struct {
	SoftwareId *string `uri:"SoftwareId" required:"true"`
	AdminPostRequest
}

type AdminPutResponse struct {
	v20230530.Response
}

type AdminPatchRequest struct {
	SoftwareId *string `uri:"SoftwareId" required:"true"`
	Zone       *string `json:"Zone,omitempty" required:"false"`
	Name       *string `json:"Name,omitempty" required:"false"`
	Desc       *string `json:"Desc,omitempty" required:"false"`
	Icon       *string `json:"Icon,omitempty" required:"false"`
	Platform   *string `json:"Platform,omitempty"  required:"false"`
	ImageId    *string `json:"ImageId,omitempty" required:"false"`
	InitScript *string `json:"InitScript,omitempty" required:"false"`
	GpuDesired *bool   `json:"GpuDesired,omitempty" required:"false"`
}

type AdminPatchResponse struct {
	v20230530.Response
}

type AdminDeleteRequest struct {
	SoftwareId *string `uri:"SoftwareId" required:"true"`
}

type AdminDeleteResponse struct {
	v20230530.Response
}

type AdminPostUsersRequest struct {
	Users     []string `json:"Users,omitempty" required:"false"`
	Softwares []string `json:"Softwares,omitempty" required:"false"`
}

type AdminPostUsersResponse struct {
	v20230530.Response
}

type AdminDeleteUsersRequest struct {
	Users     []string `json:"Users,omitempty" required:"false"`
	Softwares []string `json:"Softwares,omitempty" required:"false"`
}

type AdminDeleteUsersResponse struct {
	v20230530.Response
}
