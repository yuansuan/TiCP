package moduleconfig

import v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"

type ListModuleConfigRequest struct {
	LicenseId string `json:"LicenseId" binding:"required"`
}

type ListModuleConfigResponse struct {
	v20230530.Response `json:",inline"`
	Data               ListModuleConfigResponseData `json:"Data"`
}

type ListModuleConfigResponseData struct {
	ModuleConfigs []*GetModuleConfigResponseData `json:"ModuleConfigs"`
}

type BatchAddModuleConfigRequest struct {
	LicenseId     string                    `json:"LicenseId" binding:"required"`
	ModuleConfigs []*AddModuleConfigRequest `json:"ModuleConfigs" binding:"dive" `
}

type BatchAddModuleConfigResponseData struct {
	Ids []string `json:"Ids"`
}

type BatchAddModuleConfigResponse struct {
	v20230530.Response `json:",inline"`
	Data               BatchAddModuleConfigResponseData `json:"Data"`
}

type AddModuleConfigRequest struct {
	LicenseId  string `json:"LicenseId" binding:"required"`
	Total      int    `json:"Total" binding:"min=0,max=10000"`
	ModuleName string `json:"ModuleName" binding:"required,max=64"`
}

type AddModuleConfigResponseData struct {
	Id string `json:"Id"`
}

type AddModuleConfigResponse struct {
	v20230530.Response `json:",inline"`
	Data               AddModuleConfigResponseData `json:"Data"`
}

type PutModuleConfigRequest struct {
	Id         string `json:"Id"`
	ModuleName string `json:"ModuleName" binding:"omitempty,max=64"`
	Total      int    `json:"Total" binding:"min=0,max=10000"`
}

type PutModuleConfigResponse v20230530.Response

type GetModuleConfigRequest string

type GetModuleConfigResponse struct {
	Data *GetModuleConfigResponseData `json:"Data"`
	v20230530.Response
}

type GetModuleConfigResponseData struct {
	Id        string `json:"Id"`
	LicenseId string `json:"LicenseId"`
	// 模块名称
	ModuleName string `json:"ModuleName"`
	// 总数量
	Total int `json:"Total"`
	// 已使用数量
	UsedNum int `json:"UsedNum"`
	// 实时总数量（监控统计的）
	ActualTotal int `json:"ActualTotal"`
	// 实时已使用数量（监控统计的）
	ActualUsed int `json:"ActualUsed"`
}

type DeleteModuleConfigRequest string
type DeleteModuleConfigResponse v20230530.Response
