package dto

type LicenseInfoResponse struct {
	Id                    string                  `json:"id"`                      // id
	ManagerId             string                  `json:"manager_id"`              // 对应的manager_id
	LicenseName           string                  `json:"license_name"`            // 提供者
	MacAddr               string                  `json:"mac_addr"`                // mac地址
	ToolPath              string                  `json:"tool_path"`               // 从license server查询剩余信息的工具安装路径, 当license类型为外部时，将在对应的超算执行该tool获取信息
	LicenseUrl            string                  `json:"license_url"`             // 许可证服务器地址
	Port                  int                     `json:"port"`                    // 端口
	LicenseNum            string                  `json:"license_num"`             // license许可证序列号
	Weight                int                     `json:"weight"`                  // 调度优先级
	BeginTime             string                  `json:"begin_time"`              // 使用有效期 开始
	EndTime               string                  `json:"end_time"`                // 使用有效期 结束
	Auth                  bool                    `json:"auth"`                    // 是否授权
	LicenseEnvVar         string                  `json:"license_env_var"`         // license环境变量名称
	AllowableHpcEndpoints []string                `json:"allowable_hpc_endpoints"` // 支持的HpcEndpoint地址
	CollectorType         string                  `json:"collector_type"`          // license 类型（flex、lsdyna、altair）
	ModuleConfigInfos     []*ModuleConfigResponse `json:"module_config_infos"`     // 模块配置信息
}

type LicenseInfoAddRequest struct {
	ManagerId     string `json:"manager_id"  binding:"required"`                                                         // 对应的manager_id
	LicenseName   string `json:"license_name" binding:"max=64"`                                                          // 许可证名称
	LicenseEnvVar string `json:"license_env_var"  binding:"max=255"`                                                     // license环境变量名称
	MacAddr       string `json:"mac_addr" binding:"omitempty,mac" error:"MacAddr format error, e.g.: 00:00:00:00:00:00"` // mac地址
	LicenseUrl    string `json:"license_url" binding:"max=255"`                                                          // 许可证服务器地址
	Port          int    `json:"port" binding:"min=0,max=65535" error:"Port number ranges from 0 to 65535"`              // 端口
	LicenseNum    string `json:"license_num" binding:"omitempty,max=255"`                                                // license许可证序列号
	Weight        int    `json:"weight" binding:"min=0,max=100" error:"Weight ranges from 0 to 100"`                     // 调度优先级
	StartTime     string `json:"start_time" binding:"required"`                                                          // 使用有效期 开始
	EndTime       string `json:"end_time" binding:"required"`                                                            // 使用有效期 结束
	Auth          bool   `json:"auth"`                                                                                   // 是否授权
	ToolPath      string `json:"tool_path"`                                                                              // 从license server查询剩余信息的工具安装路径, 当license类型为外部时，将在对应的超算执行该tool获取信息
	CollectorType string `json:"collector_type" binding:"required"`                                                      // license 类型
}

type LicenseInfoAddResponse struct {
	Id string `json:"id"` // license info id
}

type LicenseInfoEditRequest struct {
	Id string `json:"id" binding:"required"` // license info id
	LicenseInfoAddRequest
}
