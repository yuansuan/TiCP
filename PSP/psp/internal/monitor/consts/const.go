package consts

const (
	NodeName   = "NodeName"
	CPUTot     = "CPUTot"     // 总核数
	CPUAlloc   = "CPUAlloc"   // 已分配核数
	CPUIdle    = "CPUIdle"    // 空闲核数
	State      = "State"      // 状态
	Partitions = "Partitions" // 队列
	OS         = "OS"         // 操作系统
	CPU        = "cpu"        // CPU 相关指标
)

const (
	Mom             = "Mom"
	Arch            = "resources_available.arch"
	PbsProState     = "state"
	Assigned        = "resources_assigned.ncpus"
	ToTalCpus       = "resources_available.ncpus"
	Queue           = "queue"
	DefaultPlatform = "resources_available.platform"
	Platform        = "platform"
)
const (
	NodeClose = "node_close"
	NodeStart = "node_start"
)

const (
	Idle    = "idle"
	Downed  = "down"
	Drain   = "drain"
	OffLine = "offline"
)

// metric 上指标的label
const (
	// HostNameLabel 节点 label
	HostNameLabel = "host_name"
	// NameLabel 指标label
	NameLabel = "name"
	// MountPathLabel 挂载盘符label
	MountPathLabel = "mount_path"
	// AppTypeLabel license 使用app类型
	AppTypeLabel = "app_type"
	// FeatureName license 里的 feature_name label
	FeatureName = "feature_name"
	// ValueTypeLabel  license 里 value_type label
	ValueTypeLabel = "value_type"
	// LicenseID license 里的licenseID
	LicenseID = "license_id"
)

const (
	RangeOneDaySec    = 1 * 24 * 60 * 60
	RangeOneWeekSec   = 7 * 24 * 60 * 60
	RangeOneMonthSec  = 30 * 24 * 60 * 60
	RangeHalfAYearSec = 180 * 24 * 60 * 60
	RangeOneYearSec   = 360 * 24 * 60 * 60

	RangeFiveMinSec  = 5 * 60
	RangeTenMinSec   = 10 * 60
	Range30MinSec    = 30 * 60
	RangeFortyMinSec = 40 * 60
	RangeOneHourSec  = 60 * 60
)

const (
	// CPUUtAvg cpu 平均利用率
	CPUUtAvg = "CPU_UT_AVG"
	// MemUtAvg 内存平均利用率
	MemUtAvg = "MEM_UT_AVG"
	//TotalIoUtAvg 磁盘总吞吐率
	TotalIoUtAvg = "TOTAL_IO_UT_AVG"
	//ReadIoUtAvg 磁盘读吞吐率
	ReadIoUtAvg = "READ_IO_UT_AVG"
	//WriteIoUtAvg 磁盘写吞吐率
	WriteIoUtAvg = "WRITE_IO_UT_AVG"
	// DiskUtAvg 磁盘使用情况
	DiskUtAvg = "DISK_UT_AVG"
	// CPUTimeSum 核时使用情况
	CPUTimeSum = "CPU_TIME_SUM"
	// LicenseAppUsedUtAvg  app的license使用率
	LicenseAppUsedUtAvg = "LICENSE_APP_USED_UT_AVG"
	// LicenseAppModuleUsedUtAvg app 模块license使用率
	LicenseAppModuleUsedUtAvg = "LICENSE_APP_MODULE_USED_UT_AVG"
	// Node Availabel 节点可用资源
	NodeAvailabel = "NODE_AVAILABEL"
)

const (
	CPUTimeSumTopSize = 10
)

const (
	PercentSymbol   = "%"
	IoBandwidthUnit = "KB/s"
	CPUTimeSumUnit  = "核时"
	JobCountUnit    = "作业数"
)

const (
	JobDeliverUserCountTitle      = "提交作业的用户数"
	JobDeliverJobCountTitle       = "投递作业数"
	JobWaitTimeStatisticAvgName   = "作业平均等待时间"
	JobWaitTimeStatisticMaxName   = "作业最大等待时间"
	JobWaitTimeStatisticTotalName = "作业总等待时间"
	JobWaitNumStatisticTotalName  = "等待作业人次"
)

const (
	ResourceUtAvgAllNodesName = "所有节点"
	CPUUtAvgName              = "CPU平均利用率"
	MemUtAvgName              = "内存平均利用率"
	TotalIoUtAvgName          = "总速率"
	ReadIoUtAvgName           = "读速率"
	WriteIoUtAvgName          = "写速率"
	CPUTimeSumName            = "核时使用情况"
	DiskUtAvgName             = "磁盘使用率情况"
)
const (
	DiskUsage = "disk_usage"
	JobStatus = "status_job"
	Feature   = "feature"
	Scheduler = "scheduler"
)

const (
	Name   = "name"
	Used   = "已使用"
	UnUsed = "未使用"
)
const (
	DiskUsagePercent = "disk_usage_percent"
)

const (
	JobStatusNum = "psp_monitor_status_job"

	// JobStateCompleted 完成
	JobStateCompleted = "Completed"
	// JobStateFailed 失败
	JobStateFailed = "Failed"
	// JobStateRunning 运行
	JobStateRunning = "Running"
	// JobStatePending 等待
	JobStatePending = "Pending"
)

const (
	NodeNormal   = 1 //节点正常
	NodeAbnormal = 2 //节点异常
)

const (
	DateTypeSecond = "second"
	dateTypeMinute = "minute"
	DateTypeHour   = "hour"

	DateTypeDay   = "day"
	DateTypeMonth = "month"
	DateTypeYear  = "year"
)
