package openapi

type ListAppRequest struct {
	ComputeType string `json:"compute_type" form:"compute_type" validate:"required,oneof=local cloud"` // 计算类型
}

type ListAppResponse struct {
	Apps []*App `json:"apps"` // 应用列表
}

type App struct {
	ID          string       `json:"id" yaml:"id"`                                         // 应用 ID
	Name        string       `json:"name" yaml:"name"`                                     // 名称
	Type        string       `json:"type" yaml:"type"`                                     // 类型
	Version     string       `json:"version" yaml:"version"`                               // 版本
	ComputeType string       `json:"compute_type" yaml:"compute_type" enums:"local,cloud"` // 计算类型
	Queues      []*QueueInfo `json:"queues" yaml:"-"`                                      // 队列信息列表
	State       string       `json:"state" yaml:"state" enums:"published,unpublished"`     // 状态
	Description string       `json:"description" yaml:"description"`                       // 描述
	SubForm     *SubForm     `json:"sub_form" yaml:"sub_form"`                             // 参数表单
}

type QueueInfo struct {
	QueueName string `json:"queue_name"` // 队列名称
	CPUNumber int64  `json:"cpu_number"` // CPU 核数
}

type SubForm struct {
	Section []*Section `json:"section" yaml:"section"` // 参数区
}

type Section struct {
	Name  string   `json:"name" yaml:"name"`   // 参数区名称
	Field []*Field `json:"field" yaml:"field"` // 参数列表
}

type Field struct {
	ID            string   `json:"id" yaml:"id"`                         // 参数 ID
	Label         string   `json:"label" yaml:"label"`                   // 标签
	Type          string   `json:"type" yaml:"type"`                     // 类型
	Required      bool     `json:"required" yaml:"required"`             // 是否必须
	DefaultValue  string   `json:"default_value" yaml:"default_value"`   // 默认值
	DefaultValues []string `json:"default_values" yaml:"default_values"` // 默认列表值
	Options       []string `json:"options" yaml:"options"`               // 选项
}

type KeyValue struct {
	Key   string `json:"key" form:"key"`     // 关键字
	Value string `json:"value" form:"value"` // 关键字对应的值
}
