package dto

import (
	"time"
)

// NodeRequest 节点详情请求数据
type NodeRequest struct {
	NodeName string `form:"node_name"` // 节点名称
}

// NodeDetailResponse 节点详情信息
type NodeDetailResponse struct {
	*NodeDetail `json:",inline"`
}

// NodeDetail 节点详情信息
type NodeDetail struct {
	NodeName        string `json:"node_name"`        // 节点名称
	NodeStatus      string `json:"node_status"`      // 节点状态
	SchedulerStatus string `json:"scheduler_status"` // 调度器状态（原始的）
	Status          string `json:"status"`           // 调度器状态（调整后的）

	CpuPercent  float64 `json:"cpu_percent"`   // CPU使用率
	CPUIdleCore int64   `json:"cpu_idle_core"` // CPU空闲核数
	UsedCore    int64   `json:"used_core"`     // 已使用核数
	NCore       int64   `json:"n_core"`        // 总核数
	CpuIdleTime float64 `json:"cpu_idle_time"` // CPU空闲时间

	R15m float64 `json:"r15m"` // 15分钟平均负载
	R5m  float64 `json:"r5m"`  // 5分钟平均负载
	R1m  float64 `json:"r1m"`  // 1分钟平均负载

	NDisk               int64   `json:"n_disk"`           // 磁盘总数
	DiskReadThroughput  float64 `json:"read_throughput"`  // 磁盘读吞吐率
	DiskWriteThroughput float64 `json:"write_throughput"` // 磁盘写吞吐率

	MaxMem       float64 `json:"max_mem"`       // 最大可用内存
	AvailableMem float64 `json:"available_mem"` // 可用内存
	UsedMem      float64 `json:"used_mem"`      // 已用内存
	FreeMem      float64 `json:"free_mem"`      // 空闲内存
	MaxSwap      float64 `json:"max_swap"`      // 最大可用交换空间
	FreeSwap     float64 `json:"free_swap"`     // 空闲交换空间
	FreeTmp      float64 `json:"free_tmp"`      // 空闲tmp空间
	MaxTmp       float64 `json:"max_tmp"`       // 最大可用tmp空间
}

// NodeListRequest 节点列表请求数据
type NodeListRequest struct {
	NodeName  string `form:"node_name"`  // 节点名称
	PageIndex int64  `form:"page_index"` // 页码
	PageSize  int64  `form:"page_size"`  // 每页数量
}

// NodeInfo 节点信息
type NodeInfo struct {
	Id              string    `json:"id"`               // 节点ID
	NodeName        string    `json:"node_name"`        // 节点名称
	NodeType        string    `json:"node_type"`        // 节点类型
	SchedulerStatus string    `json:"scheduler_status"` // 调度器状态（原始的）
	Status          string    `json:"status"`           // 调度器状态（调整后的）
	QueueName       string    `json:"queue_name"`       // 队列名称
	TotalCoreNum    int       `json:"total_core_num"`   // 总核数
	UsedCoreNum     int       `json:"used_core_num"`    // 已使用核数
	FreeCoreNum     int       `json:"free_core_num"`    // 空闲核数
	TotalMem        int       `json:"total_mem"`        // 总内存
	UsedMem         int       `json:"used_mem"`         // 已使用内存
	FreeMem         int       `json:"free_mem"`         // 空闲内存
	AvailableMem    int       `json:"available_mem"`    // 可用内存
	CreateTime      time.Time `json:"create_time"`      // 创建时间
}

type NodeListResponse struct {
	NodeInfoList []*NodeInfo `json:"list"`  // 节点列表
	Total        int64       `json:"total"` // 总数
}

// NodeOperateRequest 节点操作请求数据
type NodeOperateRequest struct {
	NodeNames []string `json:"node_names"` // 节点名称
	Operation string   `json:"operation"`  // 操作（拒绝标识：node_close、接受表示：node_start）
}

type CoreStatisticsResponse struct {
	*CoreStatistics `json:",inline"`
}

type CoreStatistics struct {
	TotalNum int `json:"total_num"` // 总核数
	FreeNum  int `json:"free_num"`  // 空闲核数
}

type NodeHpcRunningInfo struct {
	NodeName string
	State    string
	CPUAlloc int
	CPUTotal int
	CPUIdle  int
}
