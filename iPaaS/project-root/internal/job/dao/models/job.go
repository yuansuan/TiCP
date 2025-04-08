package models

import (
	"time"

	jobcreate "github.com/yuansuan/ticp/common/project-root-api/job/v1/jobcreate"
	v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/consts"
)

// Job 作业
type Job struct {
	// 作业基本信息
	ID        snowflake.ID `json:"id" xorm:"pk id comment('作业ID') BIGINT(20)"`
	Name      string       `json:"name" xorm:"not null default '' comment('作业名称') VARCHAR(255)"`
	Comment   string       `json:"comment" xorm:"comment('作业备注') TEXT"`
	UserID    snowflake.ID `json:"user_id" xorm:"user_id not null comment('用户 ID') BIGINT(20)"`
	JobSource string       `json:"job_source" xorm:"not null comment('job source') VARCHAR(20)"` // 预留

	// 状态信息
	State         int    `json:"state" xorm:"not null comment('作业状态') INT(11)"`
	SubState      int    `json:"sub_state" xorm:"not null comment('作业子状态') INT(11)"`
	StateReason   string `json:"state_reason" xorm:"comment('作业等待或者中间其他原因') TEXT"`
	ExitCode      string `json:"exit_code" xorm:"comment('作业退出码') TEXT"`
	FileSyncState string `json:"file_sync_state" xorm:"not null default '' comment('文件同步状态, 非终态：Waiting,Syncing,Pausing,Paused,Resuming,终态：Completed,Failed') VARCHAR(32)"`

	// 参数信息
	Params                      string `json:"params" xorm:"comment('用户参数') TEXT"`
	UserZone                    string `json:"user_zone" xorm:"not null default '' comment('用户选择的分区 若空为未选择') VARCHAR(64)"`
	Timeout                     int64  `json:"timeout" xorm:"not null default 0 comment('超时时间') BIGINT(20)"`
	FileClassifier              string `json:"file_classifier" xorm:"comment('文件分类器') TEXT"` // 预留
	ResourceUsageCpus           int64  `json:"resource_usage_cpus" xorm:"not null default 1 comment('用户选择使用核数') BIGINT(20)"`
	ResourceUsageMemory         int64  `json:"resource_usage_memory" xorm:"not null default 1 comment('用户选择使用内存') BIGINT(20)"`
	AllocType                   string `json:"alloc_type" xorm:"comment('分配核数的方式') VARCHAR(50)"`
	CustomStateRuleKeyStatement string `json:"custom_state_rule_key_statement" xorm:"comment('自定义状态规则key语句') TEXT"`
	CustomStateRuleResultState  string `json:"custom_state_rule_result_state" xorm:"default '' comment('自定义状态规则结果状态') VARCHAR(20)"`
	NoRound                     bool   `json:"no_round" xorm:"not null default 0 comment('单节点是否不进行取整,仅限内部用户使用') TINYINT(4)"`
	PreScheduleID               string `json:"pre_schedule_id" xorm:"pre_schedule_id comment('预调度ID') VARCHAR(64)"`

	// 作业运行信息
	HPCJobID             string `json:"hpc_job_id" xorm:"hpc_job_id comment('HPC作业ID') VARCHAR(64)"`
	Zone                 string `json:"zone" xorm:"not null default '' comment('实际运行的分区') VARCHAR(64)"`
	ResourceAssignCpus   int64  `json:"resource_assign_cpus" xorm:"not null default 0 comment('实际分配使用核数') BIGINT(20)"`
	ResourceAssignMemory int64  `json:"resource_assign_memory" xorm:"not null default 0 comment('实际分配使用内存') BIGINT(20)"`
	Command              string `json:"command" xorm:"not null comment('作业实际执行命令') TEXT"`
	WorkDir              string `json:"work_dir" xorm:"not null default '' comment('工作目录') VARCHAR(255)"`  // 实际的工作路径（hpc返回的）
	OriginJobID          string `json:"origin_job_id" xorm:"origin_job_id comment('调度器作业ID') VARCHAR(32)"` // 调度器返回的作业id，属于作业原始信息
	Queue                string `json:"queue" xorm:"comment('作业实际运行的队列') VARCHAR(32)"`
	Priority             int    `json:"priority" xorm:"not null default 0 comment('作业实际优先级') BIGINT(20)"`
	ExecHosts            string `json:"exec_hosts" xorm:"comment('作业执行节点名称列表') VARCHAR(256)"`
	ExecHostNum          int    `json:"exec_host_num" xorm:"comment('作业执行节点总数') INT(11)"`
	ExecutionDuration    int64  `json:"execution_duration" xorm:"not null default 0 comment('作业文件执行时间') INT(11)"`

	// 文件信息
	InputType               string `json:"input_type" xorm:"not null default '' comment('输入数据类型为超算存储或者远算云盒子') VARCHAR(64)"`
	InputDir                string `json:"input_dir" xorm:"not null default '' comment('盒子上的作业输入文件目录') VARCHAR(255)"`
	Destination             string `json:"destination" xorm:"not null default '' comment('输入文件的目标路径') VARCHAR(255)"` // 也是工作路径(用户指定的路径)
	OutputType              string `json:"output_type" xorm:"not null default '' comment('输出数据类型为超算存储或者远算云盒子') VARCHAR(64)"`
	OutputDir               string `json:"output_dir" xorm:"not null default '' comment('盒子上的作业输出文件目录') VARCHAR(255)"`
	NoNeededPaths           string `json:"no_needed_paths" xorm:"comment('正则表达式,符合规则的文件路径将不会进行回传') TEXT"`
	NeededPaths             string `json:"needed_paths" xorm:"comment('正则表达式,符合规则的文件路径将会进行回传') TEXT"`
	FileInputStorageZone    string `json:"file_input_storage_zone" xorm:"not null default '' comment('输入文件区域') VARCHAR(10)"`
	FileOutputStorageZone   string `json:"file_output_storage_zone" xorm:"not null default '' comment('输出文件区域') VARCHAR(10)"`
	DownloadFileSizeTotal   int64  `json:"download_file_size_total" xorm:"null default 0 comment('下载文件总大小') BIGINT(20)"`
	DownloadFileSizeCurrent int64  `json:"download_file_size_current" xorm:"null default 0 comment('下载文件当前大小') BIGINT(20)"`
	UploadFileSizeTotal     int64  `json:"upload_file_size_total" xorm:"null default 0 comment('上传文件总大小') BIGINT(20)"`
	UploadFileSizeCurrent   int64  `json:"upload_file_size_current" xorm:"null default 0 comment('上传文件当前大小') BIGINT(20)"`

	// 应用信息
	AppID   snowflake.ID `json:"app_id" xorm:"app_id not null comment('app ID') BIGINT(20)"`
	AppName string       `json:"app_name" xorm:"not null comment('app Name') VARCHAR(255)"`

	// 标志信息
	UserCancel       int `json:"user_cancel" xorm:"not null default 0 comment('用户取消标记') TINYINT(4)"`
	IsFileReady      int `json:"is_file_ready" xorm:"not null default 0 comment('作业文件是否准备完成') TINYINT(4)"`
	DownloadFinished int `json:"download_finished" xorm:"not null default 0 comment('作业下载是否完成') TINYINT(4)"`
	IsSystemFailed   int `json:"is_system_failed" xorm:"not null default 0 comment('是否系统失败') TINYINT(4)"`
	IsDeleted        int `json:"is_deleted" xorm:"not null default 0 comment('标识作业是否已被删除, 0 - 未删除, 1 - 已删除') TINYINT(4)"`

	// 计费信息
	AccountID      snowflake.ID         `json:"account_id" xorm:"account_id not null default '0' comment('账户ID') BIGINT(20)"`
	PayByAccountID snowflake.ID         `json:"pay_by_account_id" xorm:"pay_by_account_id default '0' comment('代支付账户ID') BIGINT(20)"`
	ChargeType     v20230530.ChargeType `json:"charge_type" xorm:"not null default '' comment('计费类型 [PrePaid | PostPaid]')"`
	IsPaidFinished bool                 `json:"is_paid_finished" xorm:"not null default 0 comment('是否完成付费')"`

	// 时间信息
	UploadTime       time.Time `json:"upload_time" xorm:"not null default CURRENT_TIMESTAMP comment('作业计算文件上传完成时间 hpc完全获取到所有计算文件的时间') DATETIME"`
	DownloadTime     time.Time `json:"download_time" xorm:"not null default CURRENT_TIMESTAMP comment('作业回传完成时间 hpc最后一次回传完成的时间, 即FileSyncState变成终态的时间') DATETIME"`
	PendingTime      time.Time `json:"pending_time" xorm:"not null default CURRENT_TIMESTAMP comment('变成pending 状态的时间') DATETIME"`
	RunningTime      time.Time `json:"running_time" xorm:"not null default CURRENT_TIMESTAMP comment('变成running 状态的时间') DATETIME"`
	TerminatingTime  time.Time `json:"terminating_time" xorm:"not null default CURRENT_TIMESTAMP comment('状态变成terminating 的时间') DATETIME"`
	TransmittingTime time.Time `json:"transmitting_time" xorm:"not null default CURRENT_TIMESTAMP comment('作业回传中时间 hpc最后一次回传开始的时间') DATETIME"`
	SuspendingTime   time.Time `json:"suspending_time" xorm:"not null default CURRENT_TIMESTAMP comment('状态变成suspending 的时间') DATETIME"`
	SuspendedTime    time.Time `json:"suspended_time" xorm:"not null default CURRENT_TIMESTAMP comment('状态变成suspended 的时间') DATETIME"`
	SubmitTime       time.Time `json:"submit_time" xorm:"not null default CURRENT_TIMESTAMP comment('提交给hpc的时间') DATETIME"`
	EndTime          time.Time `json:"end_time" xorm:"not null default CURRENT_TIMESTAMP comment('作业结束时间') DATETIME"`
	CreateTime       time.Time `json:"create_time" xorm:"not null default CURRENT_TIMESTAMP comment('创建时间') DATETIME"`
	UpdateTime       time.Time `json:"update_time" xorm:"not null default CURRENT_TIMESTAMP comment('更新时间') DATETIME"`
}

// TableName job
func (t *Job) TableName() string {
	return consts.TableJob
}

// AdminParams 作业参数
// swagger:model JobCreateParams
type AdminParams struct {
	jobcreate.Params        `json:",inline"`
	JobSchedulerSubmitFlags map[string]string `json:"JobSchedulerSubmitFlags,omitempty"` //自定义调度器提交参数  {"-o": "stdout.log", "-e", "stderr.log", "-l": "select=xxx"}
}

// IsLongRunning 核数*运行时间 是否超过阈值
func (t *Job) IsLongRunning(threshold int64) bool {
	if t.ResourceAssignCpus > 0 {
		return t.ResourceAssignCpus*t.ExecutionDuration > threshold
	}
	return false
}
