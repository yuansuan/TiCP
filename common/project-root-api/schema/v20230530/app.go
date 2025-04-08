package v20230530

// Application 应用
type Application struct {
	AppID              string            `json:"AppID"`   // 应用ID
	Name               string            `json:"Name"`    // 应用名称
	Type               string            `json:"Type"`    // 应用类型
	Version            string            `json:"Version"` // 版本
	AppParamsVersion   int               `json:"AppParamsVersion,omitempty"`
	Image              string            `json:"Image"` // 镜像id
	Endpoint           string            `json:"Endpoint,omitempty"`
	Command            string            `json:"Command"`       // 命令
	PublishStatus      string            `json:"PublishStatus"` // 发布状态
	Description        string            `json:"Description"`   // 描述
	IconUrl            string            `json:"IconUrl"`       // 图标url
	CoresMaxLimit      int64             `json:"CoresMaxLimit,omitempty"`
	CoresPlaceholder   string            `json:"CoresPlaceholder,omitempty"`
	FileFilterRule     string            `json:"FileFilterRule,omitempty"`
	ResidualEnable     bool              `json:"ResidualEnable"`               // 是否启用残差图
	ResidualLogRegexp  string            `json:"ResidualLogRegexp,omitempty"`  // 残差图文件，默认为工作路径下的stdout.log
	ResidualLogParser  string            `json:"ResidualLogParser,omitempty"`  // 残差图解析器类型，目前只支持以下枚举：["starccm","fluent"]
	MonitorChartEnable bool              `json:"MonitorChartEnable"`           // 是否启用监控图表
	MonitorChartRegexp string            `json:"MonitorChartRegexp,omitempty"` // 监控图表文件规则，默认为'.*\.out'
	MonitorChartParser string            `json:"MonitorChartParser,omitempty"` // 监控图表解析器类型，目前只支持以下枚举：["fluent","cfx"]
	LicenseVars        string            `json:"LicenseVars,omitempty"`
	SnapshotEnable     bool              `json:"SnapshotEnable"`  // 是否启用云图
	BinPath            string            `json:"BinPath"`         // 应用路径，map字符串，对应多个超算的路径
	ExtentionParams    string            `json:"ExtentionParams"` // 扩展参数，map字符串
	LicManagerId       string            `json:"LicManagerId"`    // 许可管理器ID
	NeedLimitCore      bool              `json:"NeedLimitCore"`   // 是否需要限制核数
	SpecifyQueue       map[string]string `json:"SpecifyQueue"`    // 指定队列
	CreateTime         string            `json:"CreateTime"`      // 创建时间
	UpdateTime         string            `json:"UpdateTime"`      // 更新时间
}

// ExtentionParams 扩展参数map
type ExtentionParams map[string]ExtentionParam

// AllowableValues 允许的值
type AllowableValues []string

// ExtentionParam 扩展参数
type ExtentionParam struct {
	Type         string          `json:"Type,omitempty"`
	ReadableName string          `json:"ReadableName,omitempty"`
	Values       AllowableValues `json:"Values,omitempty"`
	Must         bool            `json:"Must,omitempty"`
}

// Contains 是否包含
func (a AllowableValues) Contains(value string) bool {
	for _, v := range a {
		if v == value {
			return true
		}
	}
	return false
}

// 参数类型
const (
	// ExtentionParamTypeString 字符串
	ExtentionParamTypeString = "String"
	// ExtentionParamTypeStringList 参数值必须在Values中的字符串列表
	ExtentionParamTypeStringList = "StringList"
)

// ApplicationQuota 应用配额
type ApplicationQuota struct {
	ID     string `json:"ID,omitempty"`
	AppID  string `json:"AppID,omitempty"`
	UserID string `json:"UserID,omitempty"`
}

// ApplicationAllow 应用白名单
type ApplicationAllow struct {
	ID     string `json:"ID,omitempty"`
	AppID  string `json:"AppID,omitempty"`
}


// ResidualLogParserType 残差图解析器类型
const (
	// ResidualLogParserTypeStarccm starccm
	ResidualLogParserTypeStarccm = "starccm"
	// ResidualLogParserTypeFluent fluent
	ResidualLogParserTypeFluent = "fluent"
)

// MonitorChartParserType 监控图表解析器类型
const (
	// MonitorChartParserTypeFluent fluent
	MonitorChartParserTypeFluent = "fluent"
	// MonitorChartParserTypeCfx cfx
	MonitorChartParserTypeCfx = "cfx"
)
