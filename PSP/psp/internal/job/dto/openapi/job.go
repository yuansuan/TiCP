package openapi

// CreateTempDirRequest 创建作业临时目录请求数据
type CreateTempDirRequest struct {
	ComputeType string `json:"compute_type" validate:"required,oneof=local cloud"` // 计算类型(local, cloud)
}

// CreateTempDirResponse 创建作业临时目录响应数据
type CreateTempDirResponse struct {
	Path string `json:"path"` // 创建作业的临时目录
}

// JobSubmitRequest 作业提交请求数据
type JobSubmitRequest struct {
	AppID     string   `json:"app_id" validate:"required"`     // 应用 ID
	ProjectID string   `json:"project_id"`                     // 项目 ID
	MainFiles []string `json:"main_files" validate:"required"` // 主文件
	WorkDir   *WorkDir `json:"work_dir" validate:"required"`   // 工作目录
	Fields    []*Field `json:"fields"`                         // 参数信息列表
}

// WorkDir 工作目录信息
type WorkDir struct {
	Path string `json:"path"` // 工作目录
}

// Field 作业提交表单的字段信息
type Field struct {
	ID     string   `json:"id"`     // 参数 ID
	Type   string   `json:"type"`   // 参数类型
	Value  string   `json:"value"`  // 单个参数值
	Values []string `json:"values"` // 多个参数值
}

// JobSubmitResponse 作业提交请求数据
type JobSubmitResponse struct {
	JobIDs []string `json:"jobIDs"`
}

// JobDetailRequest 作业详情请求数据
type JobDetailRequest struct {
	JobID string `form:"job_id" validate:"required"` // 作业 ID
}

// JobDetailInfo 作业详情信息
type JobDetailInfo struct {
	Id             string   `json:"id"`               // 作业 ID
	AppId          string   `json:"app_id"`           // 应用 ID
	JobSetId       string   `json:"job_set_id"`       // 作业集 ID
	ProjectId      string   `json:"project_id"`       // 所属项目 ID
	Name           string   `json:"name"`             // 作业名称
	State          string   `json:"state"`            // 作业状态
	DataState      string   `json:"data_state"`       // 数据状态
	Queue          string   `json:"queue"`            // 队列名称
	ExitCode       string   `json:"exit_code"`        // 退出码
	ProjectName    string   `json:"project_name"`     // 项目名称
	AppName        string   `json:"app_name"`         // 应用名称
	UserName       string   `json:"user_name"`        // 用户名称
	JobSetName     string   `json:"job_set_name"`     // 作业集名称
	ClusterName    string   `json:"cluster_name"`     // 集群名称
	WorkDir        string   `json:"work_dir"`         // 工作目录
	Type           string   `json:"type"`             // 计算类型(local, cloud)
	CpusAlloc      string   `json:"cpus_alloc"`       // 已分配 CPU 核数
	MemAlloc       string   `json:"mem_alloc"`        // 已分配内存
	ExecDuration   string   `json:"exec_duration"`    // 执行时长(s)
	CPUTime        string   `json:"cpu_time"`         // 核时(小时)
	SubmitTime     string   `json:"submit_time"`      // 提交时间
	PendTime       string   `json:"pend_time"`        // 等待时间
	StartTime      string   `json:"start_time"`       // 开始时间
	EndTime        string   `json:"end_time"`         // 结束时间
	SuspendTime    string   `json:"suspend_time"`     // 暂停时间
	FileFilterRegs []string `json:"file_filter_regs"` // 文件过滤正则
}

// JobTerminateRequest 作业终止请求数据
type JobTerminateRequest struct {
	JobID       string `json:"job_id" validate:"required"`                         // 作业id
	ComputeType string `json:"compute_type" validate:"required,oneof=local cloud"` // 计算类型(local, cloud)
}
