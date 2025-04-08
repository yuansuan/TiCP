package licenseinfo

import (
	moduleconfig "github.com/yuansuan/ticp/common/project-root-api/license/v1/module_config"
	v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
)

type AddLicenseInfoRequest struct {
	// manager id
	ManagerId string `json:"ManagerId" binding:"required"`
	// 提供者
	Provider string `json:"Provider" binding:"max=64"`
	// mac地址
	MacAddr string `json:"MacAddr" binding:"mac" error:"MacAddr format error, e.g.: 00:00:00:00:00:00"`
	// 从license server查询剩余信息的工具安装路径, 当license类型为外部时，将在对应的超算执行该tool获取信息
	ToolPath string `json:"ToolPath" binding:"omitempty,absolutePath" error:"ToolPath must be absolute path"`
	// 许可证服务器地址
	LicenseUrl string `json:"LicenseUrl" binding:"max=255"`
	// 端口
	Port int `json:"Port" binding:"min=0,max=65535" error:"Port number ranges from 0 to 65535"`
	// license许可证服务器地址[key:hpcEndpoint->value:LicenseProxies]
	LicenseProxies map[string]LicenseProxy `json:"LicenseProxies" binding:"dive"`
	// license许可证序列号
	LicenseNum string `json:"LicenseNum" binding:"required,max=255"`
	// 调度优先级
	Weight int `json:"Weight" binding:"min=0,max=100" error:"Weight ranges from 0 to 100"`
	// 使用有效期 开始
	BeginTime string `json:"BeginTime" binding:"required"`
	// 使用有效期 结束
	EndTime string `json:"EndTime" binding:"required"`
	// 是否授权
	Auth bool `json:"Auth"`
	// 是否开启模块预调度
	EnableScheduling bool `json:"EnableScheduling"`
	// license环境变量名称
	LicenseEnvVar string `json:"LicenseEnvVar" binding:"max=255"`
	// license类型（1：自有，2：外部， 3： 寄售）
	LicenseType int `json:"LicenseType" binding:"oneof=1 2 3" error:"LicenseType only 1, 2, 3 are supported"`
	// 当license类型是外部时有效
	HpcEndpoint string `json:"HpcEndpoint" binding:"url" error:"HpcEndpoint format error"`
	// 支持的HpcEndpoint地址
	AllowableHpcEndpoints []string `json:"AllowableHpcEndpoints" binding:"nonemptyArray" error:"AllowableHpcEndpoints can't be nil or an empty array"`
	// license type, 四种
	CollectorType string `json:"CollectorType" binding:"oneof=flex lsdyna altair dsli" error:"CollectorType only flex, lsdyna, altair, dsli are supported"`
	// license server状态，Normal-正常，Abnormal-异常
	LicenseServerStatus string `json:"LicenseServerStatus"`
}

type AddLicenseInfoResponseData struct {
	Id string `json:"Id"`
}

type AddLicenseInfoResponse struct {
	v20230530.Response `json:",inline"`
	Data               *AddLicenseInfoResponseData `json:"Data"`
}

type PutLicenseInfoRequest struct {
	// license ID
	Id                    string `json:"Id" binding:"required"`
	AddLicenseInfoRequest `json:",inline"`
}

type PutLicenseInfoResponse v20230530.Response

type GetLicenseInfoRequest string

type GetLicenseInfoResponse struct {
	Data *GetLicenseInfoResponseData `json:"Data"`
	v20230530.Response
}

type GetLicenseInfoResponseData struct {
	Id string `json:"Id"`
	AddLicenseInfoRequest
	ModuleConfigs []*moduleconfig.GetModuleConfigResponseData `json:"ModuleConfigs"`
}

type DeleteLicenseInfoRequest string
type DeleteLicenseInfoResponse v20230530.Response

type LicenseProxy struct {
	// 许可证服务器地址
	Url string `json:"Url" binding:"max=255"`
	// 端口
	Port int `json:"Port" binding:"min=0,max=65535" error:"Port number ranges from 0 to 65535"`
}
