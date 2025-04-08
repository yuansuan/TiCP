package dto

type App struct {
	ID                string         `json:"id" yaml:"id"`                                         // 应用 ID
	OutAppID          string         `json:"out_app_id" yaml:"out_app_id"`                         // 外部 ID
	CloudOutAppID     string         `json:"cloud_out_app_id" yaml:"-"`                            // 云应用外部 ID
	CloudOutAppName   string         `json:"cloud_out_app_name" yaml:"-"`                          // 云应用外部名称
	Name              string         `json:"name" yaml:"name"`                                     // 名称
	Type              string         `json:"type" yaml:"type"`                                     // 类型
	Version           string         `json:"version" yaml:"version"`                               // 版本
	ComputeType       string         `json:"compute_type" yaml:"compute_type" enums:"local,cloud"` // 计算类型
	Queues            []*QueueInfo   `json:"queues" yaml:"-"`                                      // 队列信息列表
	Licenses          []*LicenseInfo `json:"licenses" yaml:"licenses"`                             // 许可证列表
	State             string         `json:"state" yaml:"state" enums:"published,unpublished"`     // 状态
	Image             string         `json:"image" yaml:"image"`                                   // 镜像名称
	BinPath           []*KeyValue    `json:"bin_path" yaml:"bin_path"`                             // 执行命令路径
	SchedulerParam    []*KeyValue    `json:"scheduler_param" yaml:"scheduler_param"`               // 调度器参数
	EnableResidual    bool           `json:"enable_residual" yaml:"enable_residual"`               // 启用残差图
	EnableSnapshot    bool           `json:"enable_snapshot" yaml:"enable_snapshot"`               // 启用云图
	ResidualLogParser string         `json:"residual_log_parser" yaml:"residual_log_parser"`       // 残差图日志解析器
	Script            string         `json:"script" yaml:"script"`                                 // 脚本
	Icon              string         `json:"icon" yaml:"icon"`                                     // 图标
	Description       string         `json:"description" yaml:"description"`                       // 描述
	HelpDoc           *HelpDoc       `json:"help_doc" yaml:"help_doc"`                             // 帮助文档
	SubForm           *SubForm       `json:"sub_form" yaml:"sub_form"`                             // 参数表单
}

type KeyValue struct {
	Key   string `json:"key" form:"key"`     // 关键字
	Value string `json:"value" form:"value"` // 关键字对应的值
}

type QueueInfo struct {
	QueueName string `json:"queue_name"` // 队列名称
	CPUNumber int64  `json:"cpu_number"` // CPU 核数
	Select    bool   `json:"select"`     // 已选择
}

type LicenseInfo struct {
	Id           string `json:"id"`            // 许可证 ID
	Name         string `json:"name"`          // 许可证名称
	Select       bool   `json:"select"`        // 已选择
	LicenceValid bool   `json:"licence_valid"` //许可证是否有效
}

type HelpDoc struct {
	Type  string `json:"type" yaml:"type"`   // 类型
	Value string `json:"value" yaml:"value"` // 内容
}

type SubForm struct {
	Section []*Section `json:"section" yaml:"section"` // 参数区
}

type Section struct {
	Name  string   `json:"name" yaml:"name"`   // 参数区名称
	Field []*Field `json:"field" yaml:"field"` // 参数列表
}

type Field struct {
	ID                      string   `json:"id" yaml:"id"`                         // 参数 ID
	Label                   string   `json:"label" yaml:"label"`                   // 标签
	Help                    string   `json:"help" yaml:"help"`                     // 帮助信息
	Type                    string   `json:"type" yaml:"type"`                     // 类型
	Required                bool     `json:"required" yaml:"required"`             // 是否必须
	Hidden                  bool     `json:"hidden" yaml:"hidden"`                 // 是否隐藏
	DefaultValue            string   `json:"default_value" yaml:"default_value"`   // 默认值
	DefaultValues           []string `json:"default_values" yaml:"default_values"` // 默认列表值
	Value                   string   `json:"value" yaml:"value"`                   // 值
	Values                  []string `json:"values" yaml:"values"`                 // 列表值
	Action                  string   `json:"action" yaml:"action"`
	Options                 []string `json:"options" yaml:"options"`
	PostText                string   `json:"post_text" yaml:"post_text"`
	FileFromType            string   `json:"file_from_type" yaml:"file_from_type"`
	IsMasterSlave           bool     `json:"is_master_slave" yaml:"is_master_slave"`
	MasterIncludeKeywords   string   `json:"master_include_keywords" yaml:"master_include_keywords"`
	MasterIncludeExtensions string   `json:"master_include_extensions" yaml:"master_include_extensions"`
	MasterSlave             string   `json:"master_slave" yaml:"master_slave"`
	OptionsFrom             string   `json:"options_from" yaml:"options_from"`
	OptionsScript           string   `json:"options_script" yaml:"options_script"`
	CustomJSONValueString   string   `json:"custom_json_value_string" yaml:"custom_json_value_string"`
	IsSupportMaster         bool     `json:"is_support_master" yaml:"is_support_master"`   // 是否支持主文件
	MasterFile              string   `json:"master_file" yaml:"master_file"`               // 主文件路径
	IsSupportWorkdir        bool     `json:"is_support_workdir" yaml:"is_support_workdir"` // 是否支持工作空间
	Workdir                 string   `json:"workdir" yaml:"workdir"`                       // 工作空间路径
}

type ListAppRequest struct {
	ComputeType   string `json:"compute_type" form:"compute_type" enums:"local,cloud"` // 计算类型
	State         string `json:"state" form:"state" enums:"published,unpublished"`     // 状态
	HasPermission bool   `json:"has_permission" form:"has_permission"`                 // 检测是否有权限
	Desktop       bool   `json:"desktop" form:"desktop"`                               // 是否桌面展示(桌面展示时不包括 PAAS 方面数据)
}

type ListAppResponse struct {
	Apps []*App `json:"apps"` // 应用列表
}

type ListTemplateRequest struct {
	ComputeType string `json:"compute_type" form:"compute_type" enums:"local,cloud"` // 计算类型
	State       string `json:"state" form:"state" enums:"published,unpublished"`     // 状态
}

type ListTemplateResponse struct {
	Apps []*App `json:"apps"` // 应用列表
}

type GetAppInfoRequest struct {
	Name        string `json:"name" form:"name"`                                     // 应用名称
	ComputeType string `json:"compute_type" form:"compute_type" enums:"local,cloud"` // 计算类型
}

type GetAppInfoResponse struct {
	App *App `json:"app"` // 应用
}

type AddAppRequest struct {
	NewType           string         `json:"new_type"`                         // 新的类型
	NewVersion        string         `json:"new_version"`                      // 新的版本
	ComputeType       string         `json:"compute_type" enums:"local,cloud"` // 计算类型
	Queues            []*QueueInfo   `json:"queues"`                           // 队列信息列表
	Licenses          []*LicenseInfo `json:"licenses"`                         // 许可证信息
	Image             string         `json:"image"`                            // 镜像名称
	BinPath           []*KeyValue    `json:"bin_path"`                         // 执行命令路径
	SchedulerParam    []*KeyValue    `json:"scheduler_param"`                  // 调度器参数
	EnableResidual    bool           `json:"enable_residual"`                  // 启用残差图
	EnableSnapshot    bool           `json:"enable_snapshot"`                  // 启用云图
	ResidualLogParser string         `json:"residual_log_parser"`              // 残差图日志解析器
	Description       string         `json:"description"`                      // 描述
	Icon              string         `json:"icon"`                             // 图标
	BaseName          string         `json:"base_name"`                        // 基础名称
	CloudOutAppID     string         `json:"cloud_out_app_id"`                 // 云应用外部 ID
}

type AddAppResponse struct{}

type UpdateAppRequest struct {
	App      *App   `json:"app"`       // 应用
	BaseName string `json:"base_name"` // 基础名称
}

type UpdateAppResponse struct{}

type DeleteAppRequest struct {
	Name        string `json:"name" form:"name"`                                     // 名称
	ComputeType string `json:"compute_type" form:"compute_type" enums:"local,cloud"` // 计算类型
}

type DeleteAppResponse struct{}

type PublishAppRequest struct {
	Names       []string `json:"names" form:"names"`                                   // 名称列表
	ComputeType string   `json:"compute_type" form:"compute_type" enums:"local,cloud"` // 计算类型
	State       string   `json:"state" form:"state" enums:"published,unpublished"`     // 状态
}

type PublishAppResponse struct{}

type SyncAppContentRequest struct {
	BaseAppId  string   `json:"base_app_id" form:"base_app_id"`   // 基于模版 ID
	SyncAppIds []string `json:"sync_app_ids" form:"sync_app_ids"` // 要同步的模版 ID 数组
}

type SyncAppContentResponse struct{}

type ListQueueRequest struct {
	AppId string `json:"app_id" form:"app_id"` // 应用 ID, 不传参默认查询所有
}

type ListQueueResponse struct {
	Queues []*QueueInfo `json:"queues"` // 队列信息列表
}

type ListLicenseRequest struct{}

type ListLicenseResponse struct {
	Licenses []*LicenseInfo `json:"licenses"` // 许可证列表
}

type GetSchedulerResourceKeyRequest struct{}

type GetSchedulerResourceKeyResponse struct {
	Keys []string `json:"keys"` // 调度器资源键列表
}

type GetSchedulerResourceValueRequest struct {
	AppId           string `json:"app_id" form:"app_id"`                       // 应用 ID
	ResourceType    string `json:"resource_type" form:"resource_type"`         // 调度器资源类型
	ResourceSubType string `json:"resource_sub_type" form:"resource_sub_type"` // 调度器资源子类型
}

type Item struct {
	Value  string `json:"value"`  // 数值
	Suffix string `json:"suffix"` // 后缀
}

type GetSchedulerResourceValueResponse struct {
	Items []*Item `json:"items"` // 调度器资源列表
}
