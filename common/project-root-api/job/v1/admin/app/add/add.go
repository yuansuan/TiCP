package add

import (
	schema "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
)

// Request 请求
type Request struct {
	Name               string `json:"Name,omitempty" xquery:"Name" form:"Name"`
	Type               string `json:"Type,omitempty" xquery:"Type" form:"Type"`
	Version            string `json:"Version,omitempty" xquery:"Version" form:"Version"`
	AppParamsVersion   int    `json:"AppParamsVersion,omitempty" xquery:"AppParamsVersion" form:"AppParamsVersion"`
	Image              string `json:"Image,omitempty" xquery:"Image" form:"Image"`
	Endpoint           string `json:"Endpoint,omitempty" xquery:"Endpoint" form:"Endpoint"`
	Command            string `json:"Command,omitempty" xquery:"Command" form:"Command"`
	Description        string `json:"Description,omitempty" xquery:"Description" form:"Description"`
	IconUrl            string `json:"IconUrl,omitempty" xquery:"IconUrl" form:"IconUrl"`
	CoresMaxLimit      int64  `json:"CoresMaxLimit,omitempty" xquery:"CoresMaxLimit" form:"CoresMaxLimit"`
	CoresPlaceholder   string `json:"CoresPlaceholder,omitempty" xquery:"CoresPlaceholder" form:"CoresPlaceholder"`
	FileFilterRule     string `json:"FileFilterRule,omitempty" xquery:"FileFilterRule" form:"FileFilterRule"`
	ResidualEnable     bool   `json:"ResidualEnable,omitempty" xquery:"ResidualEnable" form:"ResidualEnable"`             // 是否启用残差图
	ResidualLogRegexp  string `json:"ResidualLogRegexp,omitempty" xquery:"ResidualLogRegexp" form:"ResidualLogRegexp"`    // 残差图文件，默认为工作路径下的stdout.log
	ResidualLogParser  string `json:"ResidualLogParser,omitempty" xquery:"ResidualLogParser" form:"ResidualLogParser"`    // 残差图解析器类型，目前只支持以下枚举：["starccm","fluent"]
	MonitorChartEnable bool   `json:"MonitorChartEnable,omitempty" xquery:"MonitorChartEnable" form:"MonitorChartEnable"` // 是否启用监控图表
	MonitorChartRegexp string `json:"MonitorChartRegexp,omitempty" xquery:"MonitorChartRegexp" form:"MonitorChartRegexp"` // 监控图表文件规则，默认为'.*\.out'
	MonitorChartParser string `json:"MonitorChartParser,omitempty" xquery:"MonitorChartParser" form:"MonitorChartParser"` // 监控图表解析器类型，目前只支持以下枚举：["fluent","cfx"]
	LicenseVars        string `json:"LicenseVars,omitempty" xquery:"LicenseVars" form:"LicenseVars"`
	SnapshotEnable     bool   `json:"SnapshotEnable,omitempty" xquery:"SnapshotEnable" form:"SnapshotEnable"`
	BinPath            string `json:"BinPath,omitempty" xquery:"BinPath" form:"BinPath"`
	ExtentionParams    string `json:"ExtentionParams,omitempty" xquery:"ExtentionParams" form:"ExtentionParams"`
	// 如果LicManagerId为空代表则代表免费软件
	LicManagerId  string            `json:"LicManagerId,omitempty" xquery:"LicManagerId" form:"LicManagerId"`
	NeedLimitCore bool              `json:"NeedLimitCore,omitempty" xquery:"NeedLimitCore" form:"NeedLimitCore"`
	SpecifyQueue  map[string]string `json:"SpecifyQueue,omitempty" xquery:"SpecifyQueue" form:"SpecifyQueue"`
}

// Response 响应
type Response struct {
	Data            *Data `json:"Data,omitempty"`
	schema.Response `json:",inline"`
}

// Data 数据
type Data struct {
	AppID string `json:"AppID"`
}
