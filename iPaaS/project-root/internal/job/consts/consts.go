package consts

import v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"

// Job list filters
const (
	TableJob = "job"
)

// FileType 输入输出文件的类型
type FileType string

// HpcStorage 超算存储
// CloudStorage 远算云盒子
const (
	HpcStorage   FileType = "hpc_storage"
	CloudStorage FileType = "cloud_storage"
)

// String string
func (ft FileType) String() string {
	return string(ft)
}

// ToHPCFileType 转成hpcApi需要的文件类型
func (ft FileType) ToHPCFileType() v20230530.StorageType {
	switch ft {
	case HpcStorage:
		return v20230530.HPCStorageType
	case CloudStorage:
		return v20230530.CloudStorageType
	default:
		return v20230530.StorageType(ft)
	}
}

// EnableStorageType 是否是超算存储或者云盒子存储
func (ft FileType) EnableStorageType() bool {
	enabledStorageType := map[FileType]bool{
		HpcStorage:   true,
		CloudStorage: true,
	}
	return enabledStorageType[ft]
}

// EnableStorageTypeString func return all enable storage type string
func EnableStorageTypeString() []string {
	return []string{string(HpcStorage), string(CloudStorage)}
}

// Zone 分区
type Zone string

// ZoneUnknown 未知分区
const (
	ZoneUnknown Zone = "unknown"
	ZoneWuxi    Zone = "az-wuxi"
)

// String string
func (z Zone) String() string {
	return string(z)
}

const (
	// AppImagePrefix 镜像应用前缀
	AppImagePrefix = "image:"
	// LocalImagePrefix 本地应用前缀
	LocalImagePrefix = "local:"
)

// TmpWorkdirPrefix
const (
	// TmpWorkdirPrefix = "tmp_workdir_"
	TmpWorkdirPrefix = "" // 先直接拿jobid做tmp
)

const (
	// PreScheduleDir 预调度目录
	PreScheduleDir = "pre_schedule"
)

// UserCancel 用户取消
const (
	UserCancel = 1
)

const (
	// DefaultPageSize 默认分页大小
	DefaultPageSize int64 = 100
	// DefaultPageOffset 默认分页偏移
	DefaultPageOffset int64 = 0
)

const (
	// MaxNameLength 作业名称最大长度
	MaxNameLength = 255
	// MaxCommentLength 作业备注最大长度
	MaxCommentLength = 256
	// MaxCommandLength 作业命令最大长度
	MaxCommandLength = 1024 * 32
	// MaxCustomStateRuleLength 作业自定义状态规则最大长度
	MaxCustomStateRuleLength = 64
	// DefaultCoreNum 默认核数
	DefaultCoreNum = 1
	// DefaultMemory 默认内存
	DefaultMemory = 0
	// MinCoreNum 最小核数
	MinCoreNum = 1
	// MinMemory 最小内存
	MinMemory = 0
	// MaxBatchGetJobIDs 批量获取作业最大数量
	MaxBatchGetJobIDs = 100
)

var (
	InvalidChargeParams = v20230530.ChargeParams{}
)

// PublishStatus 应用发布状态
type PublishStatus string

// 全部 未发布 已发布
const (
	PublishStatusAll         PublishStatus = "all"
	PublishStatusUnpublished PublishStatus = "unpublished"
	PublishStatusPublished   PublishStatus = "published"
)

// DefaultResidualReg 默认的残差图解析文件名
const DefaultResidualReg = "stdout.log"

// DefaultMonitorChartRegexp 默认的监控图表解析文件名
const DefaultMonitorChartRegexp = ".*\\.out"

const DefaultStderr = "stderr.log"
const DefaultStdout = "stdout.log"

const (
	AppPreparedFlag = "#YS_COMMAND_PREPARED"
)
