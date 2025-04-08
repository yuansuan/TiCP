package licmanager

import (
	licenseinfo "github.com/yuansuan/ticp/common/project-root-api/license/v1/license_info"
	v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
	"time"
)

type AddLicManagerRequest struct {
	// lic管理器的名称
	AppType string `json:"AppType" binding:"min=1,max=64"`
	// 操作系统
	Os int `json:"Os" binding:"oneof=1 2" error:"Os only 1, 2 are supported"`
	// 描述
	Desc string `json:"Desc" binding:"max=512"`
	// license使用计算规则
	ComputeRule string `json:"ComputeRule" binding:"shell" error:"ComputeRule can't run"`
}

type AddLicManagerResponseData struct {
	Id string `json:"Id"`
}

type AddLicManagerResponse struct {
	v20230530.Response
	Data *AddLicManagerResponseData `json:"Data"`
}

type PutLicManagerRequest struct {
	Id                   string `json:"Id" binding:"required"`
	Status               int    `json:"Status" binding:"oneof=1 2" error:"Status only 1, 2 are supported"`
	AddLicManagerRequest `json:",inline"`
}

type PutLicManagerResponse v20230530.Response

type GetLicManagerRequest string

type GetLicManagerResponse struct {
	Data *GetLicManagerResponseData `json:"Data"`
	v20230530.Response
}

type GetLicManagerResponseData struct {
	Id string `json:"Id"`
	// 创建时间
	CreateTime time.Time `json:"CreateTime"`
	// 发布状态
	Status int `json:"Status"`
	AddLicManagerRequest
	LicenseInfos []*licenseinfo.GetLicenseInfoResponseData `json:"LicenseInfos"`
}

type DeleteLicManagerRequest string
type DeleteLicManagerResponse v20230530.Response

type ListLicManagerResponse struct {
	Data *ListLicManagerResponseData `json:"Data"`
	v20230530.Response
}

type ListLicManagerResponseData struct {
	Items []*GetLicManagerResponseData `json:"Items"`
	Total int                          `json:"Total"`
}
