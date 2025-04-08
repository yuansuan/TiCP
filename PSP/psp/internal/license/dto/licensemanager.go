package dto

import (
	"time"
)

type LicenseManagerListRequest struct {
	LicenseType string `form:"license_type"` // 许可证类型
}

type LicenseManagerListResponse struct {
	LicenseManagers []*LicenseManagerResponse `json:"license_managers"`
	Total           int                       `json:"total"`
}

type LicenseManagerResponse struct {
	Id          string    `json:"id"`           // id
	LicenseType string    `json:"license_type"` // lic管理器的名称
	Os          int       `json:"os"`           // 操作系统 1:linux 2:win
	Desc        string    `json:"desc"`         // 描述
	ComputeRule string    `json:"compute_rule"` // license使用计算规则
	CreateTime  time.Time `json:"create_time"`  // 创建时间
	Status      int       `json:"Status"`       // 发布状态
}

type LicenseManagerInfoResponse struct {
	LicenseManagers *LicenseManagerData `json:"license_manager"`
}

type LicenseManagerData struct {
	Id           string                 `json:"id"`            // id
	AppType      string                 `json:"app_type"`      // lic管理器的名称
	Os           int                    `json:"os"`            // 操作系统 1:linux 2:win
	Desc         string                 `json:"desc"`          // 描述
	ComputeRule  string                 `json:"compute_rule"`  // license使用计算规则
	CreateTime   time.Time              `json:"create_time"`   // 创建时间
	LicenseInfos []*LicenseInfoResponse `json:"license_infos"` // license 服务器信息
}

type LicenseTypeListRequest struct {
	TypeName string `json:"type_name"` // 许可证类型
}

type LicenseTypeListResponse struct {
	LicenseTypeInfos []*LicenseTypeInfo `json:"license_type_infos"` // license类型列表
}

type LicenseTypeInfo struct {
	Id          string `json:"id"`           // id
	LicenseType string `json:"license_type"` // lic管理器的名称
}

type AddLicenseManagerRequest struct {
	AppType     string `json:"app_type"`     // lic管理器的名称
	Os          int    `json:"os"`           // 操作系统
	Desc        string `json:"desc"`         // 描述
	ComputeRule string `json:"compute_rule"` // license使用计算规则
}

type AddLicenseManagerResponse struct {
	Id string `json:"id"` // license manager id
}

type EditLicenseManagerRequest struct {
	Id          string `json:"id"`           // lic manager id
	AppType     string `json:"app_type"`     // lic管理器的名称
	Os          int    `json:"os"`           // 操作系统
	Desc        string `json:"desc"`         // 描述
	ComputeRule string `json:"compute_rule"` // license使用计算规则
}

type LicenseManagerRequest struct {
	Id           string                `json:"id"`            // id
	LicenseType  string                `json:"license_type"`  // 许可证类型
	Os           int                   `json:"os"`            // 操作系统 1:linux 2:win
	Desc         string                `json:"desc"`          // 描述
	ComputeRule  string                `json:"compute_rule"`  // license 使用计算规则
	LicenseInfos []*LicenseInfoRequest `json:"license_infos"` // license 服务器信息
}
type LicenseInfoRequest struct {
	Id                string                 `json:"id"`                  // id
	ManagerId         string                 `json:"manager_id"`          // 对应的manager_id
	LicenseName       string                 `json:"license_name"`        // 许可证名称
	MacAddr           string                 `json:"mac_addr"`            // mac地址
	ToolPath          string                 `json:"tool_path"`           // 从license server查询剩余信息的工具安装路径, 当license类型为外部时，将在对应的超算执行该tool获取信息
	LicenseUrl        string                 `json:"license_url"`         // 许可证服务器地址
	Port              int                    `json:"port"`                // 端口
	LicenseNum        string                 `json:"license_num"`         // license许可证序列号
	Weight            int                    `json:"weight"`              // 调度优先级
	BeginTime         string                 `json:"begin_time"`          // 使用有效期 开始
	EndTime           string                 `json:"end_time"`            // 使用有效期 结束
	Auth              bool                   `json:"auth"`                // 是否授权
	LicenseEnvVar     string                 `json:"license_env_var"`     // license环境变量名称
	ModuleConfigInfos []*ModuleConfigRequest `json:"module_config_infos"` // 模块配置信息
}
type ModuleConfigRequest struct {
	Id         string `json:"id"`          // id
	ModuleName string `json:"module_name"` // 模块名称
	Total      int    `json:"total"`       // 总数量
}
