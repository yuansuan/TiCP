package dto

import (
	"github.com/yuansuan/ticp/PSP/psp/pkg/xtype"
)

// JobNumRequest ...
type JobNumRequest struct {
	Start int64 `form:"start"`
	End   int64 `form:"end"`
}

// AppJobResponse 应用作业数
type AppJobResponse struct {
	AppTotal int           `json:"app_total"`
	AppJobs  []*AppJobInfo `json:"app_jobs"`
}
type AppJobInfo struct {
	AppName string `json:"app_name"`
	Num     int64  `json:"num"`
}

type UserJobInfo struct {
	UserName string `json:"user_name"`
	Num      int64  `json:"num"`
}

// JobTerminateRequest 作业终止请求数据
type JobTerminateRequest struct {
	OutJobID    string `json:"out_job_id"`   // PaaS 作业 ID
	ComputeType string `json:"compute_type"` // 计算类型(local, cloud)
}

// JobListRequest 作业列表请求数据
type JobListRequest struct {
	Page      *xtype.Page      `json:"page"`       // 分页信息
	OrderSort *xtype.OrderSort `json:"order_sort"` // 排序条件
	Filter    *JobFilter       `json:"filter"`     // 过滤条件
}

// JobListResponse 作业列表返回数据
type JobListResponse struct {
	Page *xtype.PageResp `json:"page"` // 分页信息
	Jobs []*JobListInfo  `json:"jobs"` // 作业列表
}

// JobDetailRequest 作业详情请求数据
type JobDetailRequest struct {
	JobID string `form:"job_id"` // 作业 ID
}

// JobListInfo 作业列表信息
type JobListInfo struct {
	Id             string `json:"id"`              // 作业 ID
	ProjectId      string `json:"project_id"`      // 所属项目 ID
	AppId          string `json:"app_id"`          // 应用 ID
	UserId         string `json:"user_id"`         // 用户 ID
	JobSetId       string `json:"job_set_id"`      // 作业集 ID
	OutJobId       string `json:"out_job_id"`      // PaaS 作业 ID
	RealJobId      string `json:"real_job_id"`     // 调度器作业 ID
	Name           string `json:"name"`            // 作业名称
	State          string `json:"state"`           // 作业状态
	RawState       string `json:"raw_state"`       // 作业原始状态
	DataState      string `json:"data_state"`      // 数据状态
	Queue          string `json:"queue"`           // 队列名称
	ProjectName    string `json:"project_name"`    // 所属项目名称
	AppName        string `json:"app_name"`        // 应用名称
	UserName       string `json:"user_name"`       // 用户名称
	JobSetName     string `json:"job_set_name"`    // 作业集名称
	Type           string `json:"type"`            // 计算类型(local, cloud)
	Priority       string `json:"priority"`        // 优先级
	CpusAlloc      string `json:"cpus_alloc"`      // 已分配 CPU 核数
	MemAlloc       string `json:"mem_alloc"`       // 已分配内存
	ExecDuration   string `json:"exec_duration"`   // 执行时长
	EnableResidual bool   `json:"enable_residual"` // 启用残差图
	EnableSnapshot bool   `json:"enable_snapshot"` // 启用云图
	SubmitTime     string `json:"submit_time"`     // 提交时间
	PendTime       string `json:"pend_time"`       // 等待时间
	StartTime      string `json:"start_time"`      // 开始时间
	EndTime        string `json:"end_time"`        // 结束时间
	SuspendTime    string `json:"suspend_time"`    // 暂停时间
}

// JobDetailInfo 作业详情信息
type JobDetailInfo struct {
	Id             string         `json:"id"`               // 作业 ID
	AppId          string         `json:"app_id"`           // 应用 ID
	UserId         string         `json:"user_id"`          // 用户 ID
	JobSetId       string         `json:"job_set_id"`       // 作业集 ID
	OutJobId       string         `json:"out_job_id"`       // PaaS 作业 ID
	RealJobId      string         `json:"real_job_id"`      // 调度器作业 ID
	ProjectId      string         `json:"project_id"`       // 所属项目 ID
	Name           string         `json:"name"`             // 作业名称
	State          string         `json:"state"`            // 作业状态
	RawState       string         `json:"raw_state"`        // 作业原始状态
	DataState      string         `json:"data_state"`       // 数据状态
	Queue          string         `json:"queue"`            // 队列名称
	ExitCode       string         `json:"exit_code"`        // 退出码
	ProjectName    string         `json:"project_name"`     // 项目名称
	AppName        string         `json:"app_name"`         // 应用名称
	UserName       string         `json:"user_name"`        // 用户名称
	JobSetName     string         `json:"job_set_name"`     // 作业集名称
	ClusterName    string         `json:"cluster_name"`     // 集群名称
	WorkDir        string         `json:"work_dir"`         // 工作目录
	ExecHosts      string         `json:"exec_hosts"`       // 执行主机名称
	Type           string         `json:"type"`             // 计算类型(local, cloud)
	Priority       string         `json:"priority"`         // 优先级
	CpusAlloc      string         `json:"cpus_alloc"`       // 已分配 CPU 核数
	MemAlloc       string         `json:"mem_alloc"`        // 已分配内存
	ExecDuration   string         `json:"exec_duration"`    // 执行时长(s)
	CPUTime        string         `json:"cpu_time"`         // 核时(小时)
	ExecHostNum    string         `json:"exec_host_num"`    // 执行主机数量
	StateReason    string         `json:"state_reason"`     // 作业状态原因
	SubmitTime     string         `json:"submit_time"`      // 提交时间
	PendTime       string         `json:"pend_time"`        // 等待时间
	StartTime      string         `json:"start_time"`       // 开始时间
	EndTime        string         `json:"end_time"`         // 结束时间
	SuspendTime    string         `json:"suspend_time"`     // 暂停时间
	FileFilterRegs []string       `json:"file_filter_regs"` // 文件过滤正则
	Timelines      []*JobTimeLine `json:"timelines"`        // 时间线
}

type JobSetDetailRequest struct {
	JobSetID string `form:"job_set_id"` // 作业集 ID
}

type JobSetInfo struct {
	ProjectId    string `json:"project_id"`    // 所属项目 ID
	ProjectName  string `json:"project_name"`  // 所属项目名称
	JobSetId     string `json:"job_set_id"`    // 作业集 ID
	JobSetName   string `json:"job_set_name"`  // 作业集名称
	JobType      string `json:"job_type"`      // 作业类型
	AppId        string `json:"app_id"`        // 应用 ID
	AppName      string `json:"app_name"`      // 应用名称
	UserId       string `json:"user_id"`       // 用户 ID
	UserName     string `json:"user_name"`     // 用户名称
	ExecDuration string `json:"exec_duration"` // 总执行时长(s)
	JobCount     int64  `json:"job_count"`     // 作业数量
	SuccessCount int64  `json:"success_count"` // 成功数量
	FailureCount int64  `json:"failure_count"` // 失败数量
	StartTime    string `json:"start_time"`    // 开始时间
	EndTime      string `json:"end_time"`      // 结束时间
}

type JobSetDetailResponse struct {
	JobSetInfo *JobSetInfo    `json:"job_set_info"` // 作业集信息
	JobList    []*JobListInfo `json:"job_list"`     // 作业列表
}

// JobFilterInfo 作业过滤器信息
type JobFilterInfo struct {
	AppNames  []string `json:"app_names"`  // 应用名称列表
	UserNames []string `json:"user_names"` // 用户名称列表
	Queues    []string `json:"queues"`     // 队列名称列表
}

// JobFilter 作业过滤器
type JobFilter struct {
	JobID       string   `json:"job_id"`        // 作业 ID
	JobName     string   `json:"job_name"`      // 作业名称
	JobSetID    string   `json:"job_set_id"`    // 作业集 ID
	ProjectIDs  []string `json:"project_ids"`   // 所属项目 ID 列表
	JobSetNames []string `json:"job_set_names"` // 作业集名称列表
	JobTypes    []string `json:"job_types"`     // 作业类型：local|cloud
	AppNames    []string `json:"app_names"`     // 应用名称列表
	UserNames   []string `json:"user_names"`    // 用户名称列表
	Queues      []string `json:"queues"`        // 队列名称列表
	States      []string `json:"states"`        // 作业状态列表
	StarTime    int      `json:"start_time"`    // 开始时间
	EndTime     int      `json:"end_time"`      // 结束时间
}

// JobMetricFiler 指标查询条件
type JobMetricFiler struct {
	TopSize   int
	StartTime int64
	EndTime   int64
}

// JobCPUTimeMetric 核时指标返回值
type JobCPUTimeMetric struct {
	AppMetrics  []*JobCPUTimeQueryMetric
	UserMetrics []*JobCPUTimeQueryMetric
}

type JobCPUTimeQueryMetric struct {
	GroupCol string  `xorm:"group_col"`
	CPUTime  float64 `xorm:"cpu_time"`
}

type JobQueryResultMetric struct {
	Item  string  `xorm:"item"`
	Count float64 `xorm:"count"`
}

// JobCountMetric 作业提交数量指标返回值
type JobCountMetric struct {
	AppCountMetrics  []*JobQueryResultMetric
	UserCountMetrics []*JobQueryResultMetric
}

// JobWaitStatistic 作业等待指标
type JobWaitStatistic struct {
	JobWaitTimeStatisticAvg   []*JobQueryResultMetric
	JobWaitTimeStatisticMax   []*JobQueryResultMetric
	JobWaitTimeStatisticTotal []*JobQueryResultMetric
	JobWaitNumStatistic       []*JobQueryResultMetric
}

type JobStatus struct {
	State string `json:"state"`
	Num   int64  `json:"num"`
}

type JobResidualRequest struct {
	JobID string `form:"job_id"` // 作业 ID
}

type VarItem struct {
	Name   string    `json:"name"`
	Values []float64 `json:"values"`
}

type JobResidualResponse struct {
	Vars          []*VarItem `json:"vars"`
	AvailableXvar []string   `json:"available_xvar"`
}

type JobSnapshotListRequest struct {
	JobID string `form:"job_id"` // 作业 ID
}

type JobSnapshotListResponse struct {
	Snapshots map[string][]string `json:"snapshots"` // 云图集列表
}

type JobSnapshotRequest struct {
	JobID string `form:"job_id"` // 作业 ID
	Path  string `form:"path"`   // 云图路径
}

type JobSnapshotResponse struct {
	Snapshot string `json:"snapshot"` // 云图数据
}
