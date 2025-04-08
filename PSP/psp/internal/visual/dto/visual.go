package dto

import (
	"time"
)

type Session struct {
	ID          string    `json:"id" form:"id"`                     // 会话 ID
	OutAppID    string    `json:"out_app_id" form:"out_app_id"`     // 外部 ID
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

type Hardware struct {
	ID             string    `json:"id" form:"id"`                           // 硬件 ID
	Name           string    `json:"name" form:"name"`                       // 名称
	Desc           string    `json:"desc" form:"desc"`                       // 描述
	Network        int       `json:"network" form:"network"`                 // 网络带宽
	CPU            int       `json:"cpu" form:"cpu"`                         // CPU 核数
	Mem            int       `json:"mem" form:"mem"`                         // 内存大小
	GPU            int       `json:"gpu" form:"gpu"`                         // GPU 核数
	CPUModel       string    `json:"cpu_model" form:"cpu_model"`             // CPU 型号
	GPUModel       string    `json:"gpu_model" form:"gpu_model"`             // GPU 型号
	InstanceType   string    `json:"instance_type" form:"instance_type"`     // 实例类型
	InstanceFamily string    `json:"instance_family" form:"instance_family"` // 实例族
	CreateTime     time.Time `json:"create_time" form:"create_time"`         // 创建时间
	UpdateTime     time.Time `json:"update_time" form:"update_time"`         // 更新时间
	DefaultPreset  bool      `json:"default_preset" form:"default_preset"`   // 默认预设
}

type Software struct {
	ID         string       `json:"id" form:"id"`                   // 软件 ID
	Name       string       `json:"name" form:"name"`               // 名称
	Desc       string       `json:"desc" form:"desc"`               // 描述
	Platform   string       `json:"platform" form:"platform"`       // 平台
	ImageID    string       `json:"image_id" form:"image_id"`       // 镜像 ID
	State      string       `json:"state" form:"state"`             // 状态
	InitScript string       `json:"init_script" form:"init_script"` // 初始化脚本
	Icon       string       `json:"icon" form:"icon"`               // 图标
	GPUDesired bool         `json:"gpu_desired" form:"gpu_desired"` // GPU 是否必须
	Presets    []*Hardware  `json:"presets" form:"presets"`         // 软件预设
	RemoteApps []*RemoteApp `json:"remote_apps" form:"remote_apps"` // 远程应用
	CreateTime time.Time    `json:"create_time" form:"create_time"` // 创建时间
	UpdateTime time.Time    `json:"update_time" form:"update_time"` // 更新时间
}

type RemoteApp struct {
	ID         string    `json:"id" form:"id"`                   // 远程应用 ID
	Name       string    `json:"name" form:"name"`               // 名称
	Desc       string    `json:"desc" form:"desc"`               // 描述
	BaseURL    string    `json:"base_url" form:"base_url"`       // 基础地址
	Dir        string    `json:"dir" form:"dir"`                 // 目录地址
	Args       string    `json:"args" form:"args"`               // 参数
	Logo       string    `json:"logo" form:"logo"`               // 图标
	DisableGfx bool      `json:"disable_gfx" form:"disable_gfx"` // 禁用图形
	CreateTime time.Time `json:"create_time" form:"create_time"` // 创建时间
	UpdateTime time.Time `json:"update_time" form:"update_time"` // 更新时间
}

type SoftwarePreset struct {
	HardwareID string `json:"hardware_id" form:"hardware_id"` // 硬件 ID
	Default    bool   `json:"default" form:"default"`         // 默认预设
}

type DurationStatistic struct {
	AppID    string `json:"app_id" form:"app_id"`     // 软件 ID
	AppName  string `json:"app_name" form:"app_name"` // 软件名称
	Duration string `json:"duration" form:"duration"` // 时长
}

type HistoryDuration struct {
	ID        string    `json:"id" form:"id"`                 // 会话 ID
	AppName   string    `json:"app_name" form:"app_name"`     // 软件名称
	Platform  string    `json:"platform" form:"platform"`     // 平台
	Duration  string    `json:"duration" form:"duration"`     // 时长
	StartTime time.Time `json:"start_time" form:"start_time"` // 开始时间
	EndTime   time.Time `json:"end_time" form:"end_time"`     // 结束时间
}

type SoftwareUsingStatus struct {
	Id        string `json:"id" form:"id"`                 // 软件 ID
	Name      string `json:"name" form:"name"`             // 软件名称
	Icon      string `json:"icon" form:"icon"`             // 软件图标
	SessionId string `json:"session_id" form:"session_id"` // 会话 ID
	Status    string `json:"status" form:"status"`         // 会话状态
	StreamURL string `json:"stream_url" form:"stream_url"` // 流地址
}

type ListSessionsRequest struct {
	HardwareIDs []string `json:"hardware_ids" form:"hardware_ids"` // 硬件 ID 列表
	SoftwareIDs []string `json:"software_ids" form:"software_ids"` // 软件 ID 列表
	UserName    string   `json:"user_name" form:"user_name"`       // 用户名
	ProjectIDs  []string `json:"project_ids" form:"project_ids"`   // 所属项目 ID 列表
	Statuses    []string `json:"statuses" form:"statuses"`         // 状态列表
	IsAdmin     bool     `json:"is_admin" form:"is_admin"`         // 是否管理员
	PageIndex   int      `json:"page_index" form:"page_index"`     // 分页索引
	PageSize    int      `json:"page_size" form:"page_size"`       // 分页大小
}

type ListSessionsResponse struct {
	Sessions []*Session `json:"sessions" form:"sessions"` // 会话列表
	Total    int64      `json:"total" form:"total"`       // 总数
}

type StartSessionRequest struct {
	ProjectID  string   `json:"project_id" form:"project_id"`   // 项目 ID
	HardwareID string   `json:"hardware_id" form:"hardware_id"` // 硬件 ID
	SoftwareID string   `json:"software_id" form:"software_id"` // 软件 ID
	Mounts     []string `json:"mounts" form:"mounts"`           // 挂载项目 IDs
}

type StartSessionResponse struct {
	SessionID string `json:"session_id" form:"session_id"` // 会话 ID
}

type GetMountInfoRequest struct {
	ProjectID string `json:"project_id" form:"project_id"` // 项目 ID
}

type MountInfo struct {
	ID   string `json:"id"`   // 挂载 ID
	Name string `json:"name"` // 挂载名称
}

type GetMountInfoResponse struct {
	DefaultMounts []*MountInfo `json:"default_mounts"` // 默认挂载项目目录
	SelectMounts  []*MountInfo `json:"select_mount"`   // 可选挂载项目目录
	SelectLimit   int          `json:"select_limit"`   // 可选挂载目录数限制(等于默认挂载项目目录数量+可选挂载项目目录数量)
}

type CloseSessionRequest struct {
	SessionID  string `json:"session_id" form:"session_id"`   // 会话 ID
	ExitReason string `json:"exit_reason" form:"exit_reason"` // 退出原因
	Admin      bool   `json:"admin" form:"admin"`             // 是否管理员
}

type CloseSessionResponse struct {
	Success bool `json:"success" form:"success"` // 是否成功
}

type PowerOffSessionRequest struct {
	SessionID string `json:"session_id" form:"session_id"` // 会话 ID
}

type PowerOffSessionResponse struct {
	Success bool `json:"success" form:"success"` // 是否成功
}

type PowerOnSessionRequest struct {
	SessionID string `json:"session_id" form:"session_id"` // 会话 ID
}

type PowerOnSessionResponse struct {
	Success bool `json:"success" form:"success"` // 是否成功
}

type RebootSessionRequest struct {
	SessionID string `json:"session_id" form:"session_id"` // 会话 ID
}

type RebootSessionResponse struct {
	Status bool `json:"status" form:"status"` // 重启是否准备重启
}

type ReadySessionRequest struct {
	SessionID string `json:"session_id" form:"session_id"` // 会话 ID
}

type ReadySessionResponse struct {
	Ready bool `json:"ready" form:"ready"` // 是否准备就绪
}

type GetRemoteAppURLRequest struct {
	SessionID     string `json:"session_id" form:"session_id"`           // 会话 ID
	RemoteAppName string `json:"remote_app_name" form:"remote_app_name"` // 远程应用名称
}

type GetRemoteAppURLResponse struct {
	URL string `json:"url" form:"url"` // URL
}

type ListUsedProjectNamesRequest struct {
	HasUsed bool `json:"has_used" form:"has_used"` // 用户已使用
}

type ListUsedProjectNamesResponse struct {
	Names []string `json:"names" form:"names"` // 已使用项目名称列表
}

type ExportSessionInfoRequest struct {
	Start int64 `json:"start" form:"start"` // 开始时间
	End   int64 `json:"end" form:"end"`     // 结束时间
}

type ExportSessionInfoResponse struct{}

type ListHardwareRequest struct {
	Name      string `json:"name" form:"name"`             // 名称
	CPU       int    `json:"cpu" form:"cpu"`               // CPU 核数
	Mem       int    `json:"mem" form:"mem"`               // 内存大小
	GPU       int    `json:"gpu" form:"gpu"`               // GPU 核数
	HasUsed   bool   `json:"has_used" form:"has_used"`     // 用户已使用
	IsAdmin   bool   `json:"is_admin" form:"is_admin"`     // 是否管理员
	PageIndex int    `json:"page_index" form:"page_index"` // 分页索引
	PageSize  int    `json:"page_size" form:"page_size"`   // 分页大小
}

type ListHardwareResponse struct {
	Hardwares []*Hardware `json:"hardwares" form:"hardwares"` // 硬件列表
	Total     int64       `json:"total" form:"total"`         // 总数
}

type AddHardwareRequest struct {
	Name           string `json:"name" form:"name"`                       // 名称
	Desc           string `json:"desc" form:"desc"`                       // 描述
	Network        int    `json:"network" form:"network"`                 // 网络带宽
	CPU            int    `json:"cpu" form:"cpu"`                         // CPU 核数
	Mem            int    `json:"mem" form:"mem"`                         // 内存大小
	GPU            int    `json:"gpu" form:"gpu"`                         // GPU 核数
	CPUModel       string `json:"cpu_model" form:"cpu_model"`             // CPU 型号
	GPUModel       string `json:"gpu_model" form:"gpu_model"`             // GPU 型号
	InstanceType   string `json:"instance_type" form:"instance_type"`     // 实例类型
	InstanceFamily string `json:"instance_family" form:"instance_family"` // 实例族
}

type AddHardwareResponse struct {
	ID string `json:"id" form:"id"` // 硬件 ID
}

type UpdateHardwareRequest struct {
	ID             string `json:"id" form:"id"`                           // 硬件 ID
	Name           string `json:"name" form:"name"`                       // 名称
	Desc           string `json:"desc" form:"desc"`                       // 描述
	Network        int    `json:"network" form:"network"`                 // 网络带宽
	CPU            int    `json:"cpu" form:"cpu"`                         // CPU 核数
	Mem            int    `json:"mem" form:"mem"`                         // 内存大小
	GPU            int    `json:"gpu" form:"gpu"`                         // GPU 核数
	CPUModel       string `json:"cpu_model" form:"cpu_model"`             // CPU 型号
	GPUModel       string `json:"gpu_model" form:"gpu_model"`             // GPU 型号
	InstanceType   string `json:"instance_type" form:"instance_type"`     // 实例类型
	InstanceFamily string `json:"instance_family" form:"instance_family"` // 实例族
}

type UpdateHardwareResponse struct{}

type DeleteHardwareRequest struct {
	ID string `json:"id" form:"id"` // 硬件 ID
}

type DeleteHardwareResponse struct{}

type ListSoftwareRequest struct {
	Name          string `json:"name" form:"name"`                                 // 名称
	Platform      string `json:"platform" form:"platform"`                         // 平台
	State         string `json:"state" form:"state" enums:"published,unpublished"` // 状态
	HasPermission bool   `json:"has_permission" form:"has_permission"`             // 用户有权限
	HasUsed       bool   `json:"has_used" form:"has_used"`                         // 用户已使用
	IsAdmin       bool   `json:"is_admin" form:"is_admin"`                         // 是否管理员
	PageIndex     int    `json:"page_index" form:"page_index"`                     // 分页索引
	PageSize      int    `json:"page_size" form:"page_size"`                       // 分页大小
}

type ListSoftwareResponse struct {
	Softwares []*Software `json:"softwares" form:"softwares"` // 软件列表
	Total     int64       `json:"total" form:"total"`         // 总数
}

type ListSoftwareUsingStatusesRequest struct{}

type ListSoftwareUsingStatusesResponse struct {
	UsingStatuses []*SoftwareUsingStatus `json:"using_statuses" form:"using_statuses"` // 软件使用情况
}

type AddSoftwareRequest struct {
	Name       string `json:"name" form:"name"`               // 名称
	Desc       string `json:"desc" form:"desc"`               // 描述
	Platform   string `json:"platform" form:"platform"`       // 平台
	ImageID    string `json:"image_id" form:"image_id"`       // 镜像 ID
	InitScript string `json:"init_script" form:"init_script"` // 初始化脚本
	Icon       string `json:"icon" form:"icon"`               // 图标
	GPUDesired bool   `json:"gpu_desired" form:"gpu_desired"` // GPU 是否必须
}

type AddSoftwareResponse struct {
	ID string `json:"id" form:"id"` // 软件 ID
}

type UpdateSoftwareRequest struct {
	ID         string `json:"id" form:"id"`                   // 软件 ID
	Name       string `json:"name" form:"name"`               // 名称
	Desc       string `json:"desc" form:"desc"`               // 描述
	Platform   string `json:"platform" form:"platform"`       // 平台
	ImageID    string `json:"image_id" form:"image_id"`       // 镜像 ID
	InitScript string `json:"init_script" form:"init_script"` // 初始化脚本
	Icon       string `json:"icon" form:"icon"`               // 图标
	GPUDesired bool   `json:"gpu_desired" form:"gpu_desired"` // GPU 是否必须
}

type UpdateSoftwareResponse struct{}

type DeleteSoftwareRequest struct {
	ID string `json:"id" form:"id"` // 软件 ID
}

type DeleteSoftwareResponse struct{}

type PublishSoftwareRequest struct {
	Id    string `json:"id" form:"id"`                                     // 软件 ID
	State string `json:"state" form:"state" enums:"published,unpublished"` // 状态
}

type PublishSoftwareResponse struct{}

type GetSoftwarePresetsRequest struct {
	SoftwareID string `json:"software_id" form:"software_id"` // 软件 ID
}

type GetSoftwarePresetsResponse struct {
	Presets []*Hardware `json:"presets" form:"presets"` // 软件预设
}

type SetSoftwarePresetsRequest struct {
	SoftwareID string            `json:"software_id" form:"software_id"` // 软件 ID
	Presets    []*SoftwarePreset `json:"presets" form:"presets"`         // 软件预设
}

type SetSoftwarePresetsResponse struct{}

type AddRemoteAppRequest struct {
	SoftwareID string `json:"software_id" form:"software_id"` // 软件 ID
	Name       string `json:"name" form:"name"`               // 名称
	Desc       string `json:"desc" form:"desc"`               // 描述
	BaseURL    string `json:"base_url" form:"base_url"`       // 基础 URL
	Dir        string `json:"dir" form:"dir"`                 // 目录地址
	Args       string `json:"args" form:"args"`               // 参数
	Logo       string `json:"logo" form:"logo"`               // 图标
	DisableGFX bool   `json:"disable_gfx" form:"disable_gfx"` // 禁用图形
}

type AddRemoteAppResponse struct {
	ID string `json:"id" form:"id"` // 远程应用 ID
}

type UpdateRemoteAppRequest struct {
	ID         string `json:"id" form:"id"`                   // 远程应用 ID
	SoftwareID string `json:"software_id" form:"software_id"` // 软件 ID
	Name       string `json:"name" form:"name"`               // 名称
	Desc       string `json:"desc" form:"desc"`               // 描述
	BaseURL    string `json:"base_url" form:"base_url"`       // 基础 URL
	Dir        string `json:"dir" form:"dir"`                 // 目录地址
	Args       string `json:"args" form:"args"`               // 参数
	Logo       string `json:"logo" form:"logo"`               // 图标
	DisableGFX bool   `json:"disable_gfx" form:"disable_gfx"` // 禁用图形
}

type UpdateRemoteAppResponse struct{}

type DeleteRemoteAppRequest struct {
	ID string `json:"id" form:"id"` // 远程应用 ID
}

type DeleteRemoteAppResponse struct{}

type DurationStatisticRequest struct {
	AppIDs    []string `json:"app_ids" form:"app_ids"`       // 软件 IDs
	StartTime string   `json:"start_time" form:"start_time"` // 开始时间
	EndTime   string   `json:"end_time" form:"end_time"`     // 结束时间
}

type DurationStatisticResponse struct {
	Statistics []*DurationStatistic `json:"statistics" form:"statistics"` // 统计数据
}

type ListHistoryDurationRequest struct {
	AppIDs    []string `json:"app_id" form:"app_id"`         // 软件 IDs
	StartTime string   `json:"start_time" form:"start_time"` // 开始时间
	EndTime   string   `json:"end_time" form:"end_time"`     // 结束时间
	PageIndex int      `json:"page_index" form:"page_index"` // 分页索引
	PageSize  int      `json:"page_size" form:"page_size"`   // 分页大小
}

type ListHistoryDurationResponse struct {
	Statistics []*HistoryDuration `json:"statistics" form:"statistics"` // 统计数据
	Total      int64              `json:"total" form:"total"`           // 总数
}

type StatisticItem struct {
	Key   string  `json:"key"`   // 键
	Value float64 `json:"value"` // 值
}

type OriginStatisticData struct {
	Name         string           `json:"name"`          // 名称
	OriginalData []*StatisticItem `json:"original_data"` // 原始数据
}

type SessionUsageDurationStatisticRequest struct {
	Start int64 `json:"start" form:"start"` // 开始时间
	End   int64 `json:"end" form:"end"`     // 结束时间
}

type SessionUsageDurationStatisticResponse struct {
	UsageDurationBySoftwre *OriginStatisticData `json:"usage_duration_by_software"` // 软件使用时长统计数据按照软件
	UsageDurationByUser    *OriginStatisticData `json:"usage_duration_by_user"`     // 软件使用时长统计数据按照用户
}

type ExportUsageDurationStatisticRequest struct {
	Start int64 `json:"start" form:"start"` // 开始时间
	End   int64 `json:"end" form:"end"`     // 结束时间
}

type ExportUsageDurationStatisticResponse struct{}

type SessionCreateNumberStatisticRequest struct {
	Start int64 `json:"start" form:"start"` // 开始时间
	End   int64 `json:"end" form:"end"`     // 结束时间
}

type SessionCreateNumberStatisticResponse struct {
	CreateNumberBySoftwre *OriginStatisticData `json:"create_number_by_software"` // 软件创建数量统计数据按照软件
	CreateNumberByUser    *OriginStatisticData `json:"create_number_by_user"`     // 软件创建数量统计数据按照用户
}

type StatisticItems struct {
	Key   int64   `json:"t"` // 键
	Value float64 `json:"v"` // 值
}

type OriginStatisticDatas struct {
	Name         string            `json:"n"` // 名称
	OriginalData []*StatisticItems `json:"d"` // 原始数据
}

type SessionNumberStatusStatisticRequest struct {
	Start int64 `json:"start" form:"start"` // 开始时间
	End   int64 `json:"end" form:"end"`     // 结束时间
}

type SessionNumberStatusStatisticResponse struct {
	NumberStatus []*OriginStatisticDatas `json:"number_status"` // 会话数量统计数据
}
