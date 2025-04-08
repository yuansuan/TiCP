package metrics

const (
	Up   = "Up"
	Down = "Down"
)

const (
	Namespace     = "psp"
	HostName      = "host_name"
	AgentSystem   = "agent"
	MonitorSystem = "monitor"
	LicenseSystem = "license"
)

// node监控指标
const (
	CPUMetrics = "psp_agent_cpu"
	// CPUCore 内核数
	CPUCore = "core"
	// CPUIdleCore 空闲内核数
	CPUIdleCore = "idle_core"
	// CPUPercent CPU使用率
	CPUPercent = "cpu_percent"
	// CPUIdleTime CPU空闲时间
	CPUIdleTime = "cpu_idle_time"

	// NodeStatus agent服务状态
	NodeStatus = "node_status"
)
const (
	AvgLoadMetrics = "psp_agent_loadavg"

	// Load1m 1m负载
	Load1m = "load_1m"
	// Load5m 5m负载
	Load5m = "load_5m"
	// Load15m 15m负载
	Load15m = "load_15m"
)
const (
	Normal = 1
)

const (
	MemoryMetrics = "psp_agent_memory"

	// MemoryTotal 最大可用内存
	MemoryTotal = "total"
	// MemoryAvailable 空闲内存
	MemoryAvailable = "available"
	// MemoryUsed 已用内存
	MemoryUsed = "used"
	// MemoryPercent 内存使用率
	MemoryPercent = "memory_percent"
	// MemoryFree 空闲内存
	MemoryFree = "free"
	// MemoryTmpTotal 最大可用tmp空间
	MemoryTmpTotal = "tmp_total"
	// MemoryTmpFree 空闲tmp空间
	MemoryTmpFree = "tmp_free"
	// MemorySwapTotal 最大可用交换空间
	MemorySwapTotal = "swap_total"
	// MemorySwapFree 空闲交换空间
	MemorySwapFree = "swap_free"
)

const (
	DiskMetrics = "psp_agent_disk"

	// DiskWriteThroughput 磁盘写吞吐率
	DiskWriteThroughput = "write_throughput"
	// DiskReadThroughput 磁盘读吞吐率
	DiskReadThroughput = "read_throughput"
	// DiskTotalThroughput 磁盘总吞吐率
	DiskTotalThroughput = "total_throughput"
)
const (
	DiskUsage = "psp_monitor_disk_usage"
	// DiskUsagePercent 磁盘使用率
	DiskUsagePercent = "disk_monitor_percent"
)

const (
	LicenseMonitorFeature  = "psp_monitor_feature"
	CollectorScrapeSuccess = "psp_agent_collector_scrape_success"

	// FeatureUsage 模块license使用数量
	FeatureUsage = "feature_usage"
	// FeatureTotal 模块license总数量
	FeatureTotal = "feature_total"
	// FeatureAvailable 模块license可用数量
	FeatureAvailable = "feature_available"
	// FeatureUsagePercent 模块license 使用率
	FeatureUsagePercent = "feature_usage_percent"
)

const (
	SchedulerMonitorFeature = "psp_monitor_scheduler"

	// SchedulerStatus 调度器状态
	SchedulerStatus = "scheduler_status"
)
