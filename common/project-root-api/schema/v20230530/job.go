package v20230530

// Progress 同步进度
type Progress struct {
	TotalSize int `json:"TotalSize"` // 总大小
	Progress  int `json:"Progress"`  // 进度
}

type JobInfo struct {
	// 基本信息
	ID            string         `json:"ID"`            //作业ID
	Name          string         `json:"Name"`          //作业名称
	JobState      string         `json:"JobState"`      //作业状态, InitiallySuspended, Pending（中央调度器等待、数据上传、超算调度器等待都在这个状态，具体）, Running, Failed（包括系统失败和用户程序失败，具体原因可以检查StateReason和ExitCode）, Terminating（中止操作中）, Terminated（已中止，无法通过resnum操作恢复）, Completed（成功结束）, Suspending（暂停操作中）, Suspended（已暂停，可通过resume操作恢复）
	FileSyncState string         `json:"FileSyncState"` //文件同步状态, 非终态：Waiting,Syncing,Pausing,Paused,Resuming,终态：Completed,Failed
	StateReason   string         `json:"StateReason"`   //作业状态原因 eg: "Uploading：data is transting"
	AllocResource *AllocResource `json:"AllocResource"` //作业分配的资源
	AllocType     string         `json:"AllocType"`     //作业CPU资源的分配方式：average or other
	ExecHostNum   int            `json:"ExecHostNum"`   //作业执行节点总数
	Zone          string         `json:"Zone"`          //作业区域
	Workdir       string         `json:"Workdir"`       //作业工作目录
	OutputDir     string         `json:"OutputDir"`     //盒子作业结果输出路径
	NoNeededPaths string         `json:"NoNeededPaths"` //正则表达式,符合规则的文件路径将不会进行回传
	NeededPaths   string         `json:"NeededPaths"`   //正则表达式,符合规则的文件路径将会进行回传
	Parameters    string         `json:"Parameters"`    //作业提交时的Parameters参数
	NoRound       bool           `json:"NoRound"`       //单节点是否不进行取整,仅限内部用户使用
	PreScheduleID string         `json:"PreScheduleID"` //作业预调度的ID

	// 时间信息
	PendingTime     string `json:"PendingTime"`     //作业等待时间 状态变成Pending的时间
	RunningTime     string `json:"RunningTime"`     //作业运行时间 状态变成Running的时间
	TerminatingTime string `json:"TerminatingTime"` //作业终止中时间 状态变成Terminating的时间
	SuspendingTime  string `json:"SuspendingTime"`  //作业挂起中时间(暂留)
	SuspendedTime   string `json:"SuspendedTime"`   //作业挂起时间(暂留)
	EndTime         string `json:"EndTime"`         //作业结束时间 状态变成终态的时间
	CreateTime      string `json:"CreateTime"`      //作业创建时间
	UpdateTime      string `json:"UpdateTime"`      //作业更新时间

	FileReadyTime    string `json:"FileReadyTime"`    //作业计算文件上传完成时间 hpc完全获取到所有计算文件的时间
	TransmittingTime string `json:"TransmittingTime"` //作业回传中时间 hpc最后一次回传开始的时间
	TransmittedTime  string `json:"TransmittedTime"`  //作业回传完成时间 hpc最后一次回传完成的时间, 即FileSyncState变成终态的时间

	// 文件同步信息
	DownloadProgress *DownloadProgress `json:"DownloadProgress"` //下载进度
	UploadProgress   *UploadProgress   `json:"UploadProgress"`   //上传进度

	// 完成信息
	ExecutionDuration int    `json:"ExecutionDuration"` //作业执行时长,求解器运行时间，单位秒。也是需要计量收费的时间。
	ExitCode          string `json:"ExitCode"`          //程序退出码，格式： code:signal 。第一字段是程序的退出码，第二个字段是程序收到的信号
	IsSystemFailed    bool   `json:"IsSystemFailed"`    //是否系统失败
	StdoutPath        string `json:"StdoutPath"`        //作业标准输出路径
	StderrPath        string `json:"StderrPath"`        //作业标准错误输出路径
}

type AdminJobInfo struct {
	JobInfo
	// admin 部分信息
	Queue       string `json:"Queue"`       // 作业实际运行的队列
	Priority    int    `json:"Priority"`    // 作业实际优先级(预留)
	OriginJobID string `json:"OriginJobID"` // 调度器作业ID
	ExecHosts   string `json:"ExecHosts"`   // 作业执行节点名称列表
	SubmitTime  string `json:"SubmitTime"`  // 作业提交时间
	UserID      string `json:"UserID"`      // 作业提交用户ID
	HPCJobID    string `json:"HPCJobID"`    // HPC作业ID
	IsDeleted   bool   `json:"IsDeleted"`   // 是否已删除
}

// DownloadProgress 下载进度
type DownloadProgress struct {
	*Progress `json:",inline"`
}

// UploadProgress 上传进度
type UploadProgress struct {
	*Progress `json:",inline"`
}

// AllocResource 作业分配的资源
type AllocResource struct {
	Cores  int `json:"Cores"`  //作业分配的核数
	Memory int `json:"Memory"` //作业分配的内存
}

// JobCpuUsage 作业CPU使用率
type JobCpuUsage struct {
	JobID           string             `json:"JobID"`
	AverageCpuUsage float64            `json:"AverageCpuUsage"`
	NodeUsages      map[string]float64 `json:"NodeUsages"`
}
