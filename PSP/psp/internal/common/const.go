package common

import "github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"

// 系统常量
const (
	SysEnvProd = "prod"
)

// 通用常量
const (
	StorageTypeHPC = "hpc_storage"

	Published   = "published"   //发布
	Unpublished = "unpublished" //未发布

	Local = "local" //本地
	Cloud = "cloud" //云端

	SystemFileNameMaxLength = 180 // 注意: 操作系统允许的文件(目录)最大长度是 255
	// CSVExportNumber csv 导出文件每次查询数据库数据数量
	CSVExportNumber = 30000

	UserName = "userName"

	// DecimalPlaces 小数位数
	DecimalPlaces = 5

	StringParamLengthLimit255 = 255

	TraceId = "trace_id"
)

// 符号常量
const (
	Empty            = ""
	Dot              = "."
	Colon            = ":"
	Tab              = "\t"
	Blank            = " "
	Underline        = "_"
	LeftParentheses  = "("
	RightParentheses = ")"
	Bar              = "-"
	Slash            = "/"
	Png              = "png"
	Yaml             = "yaml"
	Json             = "json"
)

const (
	Deleted = 1
	Normal  = 0
)
const (
	FloatPrecision0 = 0
	FloatPrecision2 = 2
	FloatPrecision6 = 6
)

// 时间常量
const (
	DefaultEmptyTime        = "1970-01-01T08:00:00+08:00"
	DatetimeFormat          = "2006-01-02 15:04:05"
	DatetimeFormatToHour    = "2006-01-02 15"
	DateOnly                = "2006-01-02"
	TimeOnly                = "15:04:05"
	StartTimeFormat         = "00:00:00"
	EndTimeFormat           = "23:59:59"
	DateUnderLineFormatOnly = "2006_01_02"
	YearMonthDayFormat      = "20060102"

	OneDayToSecond = 24 * 60 * 60
)

// 消息常量
const (
	NoticeWebsocketTopic = "notice_websocket_topic"
	NoticeWebsocketGroup = "notice_websocket_group"
	NoticeWebsocketKey   = "notice_websocket_key"
)

// 权限常量
const (
	PermissionResourceTypeLocalApp       = "local_app"
	PermissionResourceTypeAppCloudApp    = "cloud_app"
	PermissionResourceTypeVisualHardware = "visual_hardware"
	PermissionResourceTypeVisualSoftware = "visual_software"

	// PermissionResourceTypeInternal resource type - Internal resource
	PermissionResourceTypeInternal = "internal"

	// PermissionResourceTypeSystem resource type - System resource
	PermissionResourceTypeSystem = "system"

	// PermissionResourceTypeApi resource type - HttpApi resource
	PermissionResourceTypeApi = "api"

	ResourceActionNONE   = "NONE"
	ResourceActionGET    = "GET"
	ResourceActionPOST   = "POST"
	ResourceActionPUT    = "PUT"
	ResourceActionDELETE = "DELETE"

	ResourceProjectName         = "project_manager"
	ResourceSysManagerName      = "sys_manager"
	ResourceSecurityManagerName = "security_approval"
	ResourceNormalName          = "NormalName"

	ENABLE_CUSTOM  = int32(1)
	DISABLE_CUSTOM = int32(-1)
)

// 消息类型
const (
	JobEventType       = "job_event_type"
	ShareFileEventType = "share_file_event_type"
	ProjectEventType   = "project_event_type"
	SessionEventType   = "session_event_type"
	ApproveEventType   = "approve_event_type"
)

// 分页常量
const (
	DefaultPageOffset  = 0
	DefaultPageIndex   = 1
	DefaultMaxPageSize = 1000
)

const (
	// Terminated 终止
	Terminated = "Terminated"
)

// 缓存键常量
const (
	HpcUploadTaskPreKey = "hpc_upload_task"
	CopyFilePreKey      = "copy_file"
	CompressPreKey      = "compress"
	StorageModule       = "storage"
)

// paas层openapi错误码
const (
	// PaasErrorCodePathExists 文件已经存在
	PaasErrorCodePathExists = "PathExists"
)

// 通用文件夹
const (
	// PublicFolderPath 共享目录
	PublicFolderPath = "public"
	// WorkspaceFolderPath 工作目录
	WorkspaceFolderPath = "workspace"
	// PersonalFolderPath 个人工作目录
	PersonalFolderPath = "workspace/personal"
	// ProjectFolderPath 项目目录
	ProjectFolderPath = "project"
)

// 日志常量
const (
	LOG_FILE_MANAGER = "文件管理"
)

// 项目常量
const (
	PersonalProjectID   = snowflake.ID(1687026933658816512)
	PersonalProjectName = "personal"

	// ProjectInit 初始化
	ProjectInit = "Init"
	// ProjectRunning 运行中
	ProjectRunning = "Running"
	// ProjectTerminated 已终止
	ProjectTerminated = "Terminated"
	// ProjectCompleted 已完成
	ProjectCompleted = "Completed"
)

// DeployMode 部署模式
const (
	DeployModeLocal  = "local"
	DeployModeCloud  = "cloud"
	DeployModeHybrid = "hybrid"
)

const (
	HttpHeaderOpenapiCertificate = "Openapi-Certificate"
	HttpOpenapiUrl               = "/api/v1/openapi/"
)
