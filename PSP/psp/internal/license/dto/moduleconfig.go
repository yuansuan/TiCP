package dto

type ModuleConfigResponse struct {
	Id         string `json:"id"`          // id
	ModuleName string `json:"module_name"` // 模块名称
	Total      int    `json:"total"`       // 总数量
	UsedNum    int    `json:"used_num"`    // 已使用数量
	FreeNum    int    `json:"free_num"`    // 空闲数量
}

type ModuleConfigListRequest struct {
	LicenseId string `form:"license_id"` // license id
}

type ModuleConfigListResponse struct {
	UsedPercent       string                  `json:"used_percent"`        // 使用百分比
	ModuleConfigInfos []*ModuleConfigResponse `json:"module_config_infos"` // 模块配置信息
}

type AddModuleConfigRequest struct {
	LicenseId  string `json:"license_id"`  // license info id
	Total      int    `json:"total" `      // license 总数
	ModuleName string `json:"module_name"` // 模块配置名字
}

type AddModuleConfigResponse struct {
	Id string `json:"id"` // 保存后的 module config id
}

type EditModuleConfigRequest struct {
	LicenseId  string `json:"license_id"`  // license info id
	Id         string `json:"id"`          // module config id
	ModuleName string `json:"module_name"` // 模块名称
	Total      int    `json:"total"`       // license 总数
}
