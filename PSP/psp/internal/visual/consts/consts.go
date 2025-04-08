package consts

import "github.com/yuansuan/ticp/PSP/psp/internal/common"

const (
	DefaultSyncNumber = 300 // 经测试, 上限值: 370个 ID

	UnixMilliToHour = 1000 * 60 * 60 * 1.0

	DefaultSyncDataInterval   = 10
	DefaultSyncStatusInterval = 5
)

// 会话状态枚举
const (
	SessionStatusPending     = "PENDING"
	SessionStatusStarting    = "STARTING"
	SessionStatusStarted     = "STARTED"
	SessionStatusClosing     = "CLOSING"
	SessionStatusClosed      = "CLOSED"
	SessionStatusUnAvailable = "UNAVAILABLE"
	SessionStatusRebooting   = "REBOOTING"
	SessionStatusPoweringOff = "POWERING OFF"
	SessionStatusPowerOff    = "POWER OFF"
	SessionStatusPoweringOn  = "POWERING ON"
)

const (
	NumberEqualZeroMark       = -1
	NumberGreaterThanZeroMark = -2
)

const (
	DefaultUserExistProjecReason = "when user exist project or the project end, auto close the sessions relation with the project"
)

var PublishMap = map[string]string{
	common.Unpublished: "取消发布",
	common.Published:   "发布",
}

const (
	EndSessionContent = "请将项目会话需要保留的数据提前备份，以免数据丢失。请悉知。"
)

const (
	PlatformTypeLinux = "LINUX"
)

const (
	ReportTypeSessionUsageDuration = "session_usage_duration"
	ReportTypeSessionCreateNumber  = "session_create_number"
)

const (
	DimensionTypeSoftware = "software"
	DimensionTypeUser     = "user"
)
