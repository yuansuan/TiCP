package openapi

import (
	"time"
)

type ListSoftwareRequest struct {
	PageIndex int `json:"page_index" form:"page_index"` // 分页索引
	PageSize  int `json:"page_size" form:"page_size"`   // 分页大小
}

type ListSoftwareResponse struct {
	Softwares []*Software `json:"softwares" form:"softwares"` // 软件列表
	Total     int64       `json:"total" form:"total"`         // 总数
}

type Software struct {
	ID         string    `json:"id" form:"id"`                   // 软件 ID
	Name       string    `json:"name" form:"name"`               // 名称
	Desc       string    `json:"desc" form:"desc"`               // 描述
	Platform   string    `json:"platform" form:"platform"`       // 平台
	State      string    `json:"state" form:"state"`             // 状态
	GPUDesired bool      `json:"gpu_desired" form:"gpu_desired"` // GPU 是否必须
	CreateTime time.Time `json:"create_time" form:"create_time"` // 创建时间
	UpdateTime time.Time `json:"update_time" form:"update_time"` // 更新时间
}

type Hardware struct {
	ID         string    `json:"id" form:"id"`                   // 硬件 ID
	Name       string    `json:"name" form:"name"`               // 名称
	Desc       string    `json:"desc" form:"desc"`               // 描述
	Network    int       `json:"network" form:"network"`         // 网络带宽
	CPU        int       `json:"cpu" form:"cpu"`                 // CPU 核数
	Mem        int       `json:"mem" form:"mem"`                 // 内存大小
	GPU        int       `json:"gpu" form:"gpu"`                 // GPU 核数
	CPUModel   string    `json:"cpu_model" form:"cpu_model"`     // CPU 型号
	GPUModel   string    `json:"gpu_model" form:"gpu_model"`     // GPU 型号
	CreateTime time.Time `json:"create_time" form:"create_time"` // 创建时间
	UpdateTime time.Time `json:"update_time" form:"update_time"` // 更新时间
}

type ListHardwareRequest struct {
	PageIndex int `json:"page_index" form:"page_index"` // 分页索引
	PageSize  int `json:"page_size" form:"page_size"`   // 分页大小
}

type ListHardwareResponse struct {
	Hardwares []*Hardware `json:"hardwares" form:"hardwares"` // 硬件列表
	Total     int64       `json:"total" form:"total"`         // 总数
}

type CloseSessionRequest struct {
	SessionID  string `json:"session_id" form:"session_id" validate:"required"`   // 会话 ID
	ExitReason string `json:"exit_reason" form:"exit_reason" validate:"required"` // 退出原因
}

type CloseSessionResponse struct {
	Success bool `json:"success" form:"success"` // 是否成功
}

type RebootSessionRequest struct {
	SessionID string `json:"session_id" form:"session_id" validate:"required"` // 会话 ID
}

type RebootSessionResponse struct {
	Status bool `json:"status" form:"status"` // 重启是否准备重启
}

type StartSessionRequest struct {
	ProjectID  string   `json:"project_id" form:"project_id"`                       // 项目 ID
	HardwareID string   `json:"hardware_id" form:"hardware_id" validate:"required"` // 硬件 ID
	SoftwareID string   `json:"software_id" form:"software_id" validate:"required"` // 软件 ID
	Mounts     []string `json:"mounts" form:"mounts"`                               // 挂载项目 IDs
}

type StartSessionResponse struct {
	SessionID string `json:"session_id" form:"session_id"` // 会话 ID
}

type ListSessionsRequest struct {
	HardwareIDs []string `json:"hardware_ids" form:"hardware_ids"` // 硬件 ID 列表
	SoftwareIDs []string `json:"software_ids" form:"software_ids"` // 软件 ID 列表
	ProjectIDs  []string `json:"project_ids" form:"project_ids"`   // 所属项目 ID 列表
	Statuses    []string `json:"statuses" form:"statuses"`         // 状态列表
	PageIndex   int      `json:"page_index" form:"page_index"`     // 分页索引
	PageSize    int      `json:"page_size" form:"page_size"`       // 分页大小
}

type ListSessionsResponse struct {
	Sessions []*Session `json:"sessions" form:"sessions"` // 会话列表
	Total    int64      `json:"total" form:"total"`       // 总数
}

type Session struct {
	ID          string    `json:"id" form:"id"`                     // 会话 ID
	UserName    string    `json:"user_name" form:"user_name"`       // 用户名
	ProjectName string    `json:"project_name" form:"project_name"` // 项目名称
	Status      string    `json:"status" form:"status"`             // 状态
	StreamURL   string    `json:"stream_url" form:"stream_url"`     // 流地址
	ExitReason  string    `json:"exit_reason" form:"exit_reason"`   // 退出原因
	Duration    string    `json:"duration" form:"duration"`         // 使用时长
	Hardware    *Hardware `json:"hardware" form:"hardware"`         // 硬件
	Software    *Software `json:"software" form:"software"`         // 软件
	StartTime   time.Time `json:"start_time" form:"start_time"`     // 开始时间
	EndTime     string    `json:"end_time" form:"end_time"`         // 结束时间
	CreateTime  time.Time `json:"create_time" form:"create_time"`   // 创建时间
	UpdateTime  time.Time `json:"update_time" form:"update_time"`   // 更新时间
}

type SessionInfoRequest struct {
	SessionID string `json:"session_id" form:"session_id" validate:"required"` // 会话 ID
}
