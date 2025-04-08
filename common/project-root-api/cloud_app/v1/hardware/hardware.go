package hardware

import v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"

type APIGetRequest struct {
	HardwareId *string `uri:"HardwareId" required:"true"`
}

type APIGetResponse struct {
	v20230530.Response `json:",inline"`

	Data *APIGetResponseData `json:"Data,omitempty"`
}

type APIGetResponseData struct {
	v20230530.Hardware
}

type APIListRequest struct {
	Zone       *string `query:"Zone" required:"false"`
	Name       *string `query:"Name" required:"false"`
	Cpu        *int    `query:"Cpu" required:"false"`
	Mem        *int    `query:"Mem" required:"false"`
	Gpu        *int    `query:"Gpu" required:"false"`
	PageSize   *int    `query:"PageSize" required:"false"`
	PageOffset *int    `query:"PageOffset" required:"false"`
}

type APIListResponse struct {
	v20230530.Response

	Data *APIListResponseData `json:"Data,omitempty"`
}

type APIListResponseData struct {
	Hardware []*v20230530.Hardware `json:"Hardware,omitempty"`

	Offset     int `json:"Offset"`
	Size       int `json:"Size"`
	Total      int `json:"Total"`
	NextMarker int `json:"NextMarker"`
}

type AdminGetRequest struct {
	HardwareId *string `uri:"HardwareId" required:"true"`
}

type AdminGetResponse struct {
	v20230530.Response `json:",inline"`

	Data *AdminGetResponseData `json:"Data,omitempty"`
}

type AdminGetResponseData struct {
	v20230530.Hardware
}

type AdminListRequest struct {
	UserId     *string `query:"UserId" required:"false"`
	Zone       *string `query:"Zone" required:"false"`
	Name       *string `query:"Name" required:"false"`
	Cpu        *int    `query:"Cpu" required:"false"`
	Mem        *int    `query:"Mem" required:"false"`
	Gpu        *int    `query:"Gpu" required:"false"`
	PageSize   *int    `query:"PageSize" required:"false"`
	PageOffset *int    `query:"PageOffset" required:"false"`
}

type AdminListResponse struct {
	v20230530.Response

	Data *AdminListResponseData `json:"Data,omitempty"`
}

type AdminListResponseData struct {
	Hardware []*v20230530.Hardware `json:"Hardware,omitempty"`

	Offset     int `json:"Offset"`
	Size       int `json:"Size"`
	Total      int `json:"Total"`
	NextMarker int `json:"NextMarker"`
}

type AdminPostRequest struct {
	Zone           *string `json:"Zone,omitempty" required:"true"`
	Name           *string `json:"Name,omitempty" required:"true"`
	Desc           *string `json:"Desc,omitempty" required:"false"`
	InstanceType   *string `json:"InstanceType,omitempty" required:"true"`
	InstanceFamily *string `json:"InstanceFamily,omitempty" required:"false"`
	Network        *int    `json:"Network,omitempty" required:"false"`
	Cpu            *int    `json:"Cpu,omitempty" required:"true"`
	CpuModel       *string `json:"CpuModel,omitempty" required:"false"`
	Mem            *int    `json:"Mem,omitempty" required:"true"`
	Gpu            *int    `json:"Gpu,omitempty" required:"false"`
	GpuModel       *string `json:"GpuModel,omitempty" required:"false"`
}

type AdminPostResponse struct {
	v20230530.Response

	Data *AdminPostResponseData `json:"Data,omitempty"`
}

type AdminPostResponseData struct {
	HardwareId string `json:"HardwareId"`
}

type AdminPutRequest struct {
	HardwareId *string `uri:"HardwareId" required:"true"`
	AdminPostRequest
}

type AdminPutResponse struct {
	v20230530.Response
}

type AdminPatchRequest struct {
	HardwareId     *string `uri:"HardwareId" required:"true"`
	Zone           *string `json:"Zone,omitempty" required:"false"`
	Name           *string `json:"Name,omitempty" required:"false"`
	Desc           *string `json:"Desc,omitempty" required:"false"`
	InstanceType   *string `json:"InstanceType,omitempty" required:"false"`
	InstanceFamily *string `json:"InstanceFamily,omitempty" required:"false"`
	Network        *int    `json:"Network,omitempty" required:"false"`
	Cpu            *int    `json:"Cpu,omitempty" required:"false"`
	CpuModel       *string `json:"CpuModel,omitempty" required:"false"`
	Mem            *int    `json:"Mem,omitempty" required:"false"`
	Gpu            *int    `json:"Gpu,omitempty" required:"false"`
	GpuModel       *string `json:"GpuModel,omitempty" required:"false"`
}

type AdminPatchResponse struct {
	v20230530.Response
}

type AdminDeleteRequest struct {
	HardwareId *string `uri:"HardwareId" required:"true"`
}

type AdminDeleteResponse struct {
	v20230530.Response
}

type AdminPostUsersRequest struct {
	Users     []string `json:"Users,omitempty" required:"false"`
	Hardwares []string `json:"Hardwares,omitempty" required:"false"`
}

type AdminPostUsersResponse struct {
	v20230530.Response
}

type AdminDeleteUsersRequest struct {
	Users     []string `json:"Users,omitempty" required:"false"`
	Hardwares []string `json:"Hardwares,omitempty" required:"false"`
}

type AdminDeleteUsersResponse struct {
	v20230530.Response
}
