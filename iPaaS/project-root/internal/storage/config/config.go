package config

import (
	"context"
	"sync"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/env"
	"github.com/yuansuan/ticp/common/go-kit/logging"
)

var Mut = &sync.Mutex{}

type CustomT struct {
	StorageType          string         `yaml:"storage_type"`
	AuthEnabled          *bool          `yaml:"auth_enabled"`
	Local                Local          `yaml:"local"`
	YsId                 string         `yaml:"ys_id"`
	AccessKeyId          string         `yaml:"access_key_id"`
	AccessKeySecret      string         `yaml:"access_key_secret"`
	IamServerUrl         string         `yaml:"iam_server_url"`
	TmpFileCleanup       TmpFileCleanup `yaml:"tmp_file_cleanup"`
	Quota                Quota          `yaml:"quota"`
	OperationLog         OperationLog   `yaml:"operation_log"`
	SharedHost           string         `yaml:"shared_host"`
	ShareRegisterAddress string         `yaml:"share_register_address"`
	TimeoutCleanup       TimeoutCleanup `yaml:"timeout_cleanup"`
}

type Local struct {
	RootPath string `yaml:"root_path"`
	LinkType string `yaml:"link_type"`
}

type TimeoutCleanup struct {
	// 是否启动TimeoutCleanup
	Enabled bool `yaml:"enabled"`
	// direct indirect
	Mode string `yaml:"mode"`
	// 单位 分钟
	CleanupInterval int `yaml:"cleanup_interval"`
	// 单位 分钟
	TimeoutDuration int `yaml:"timeout_duration"`
	// 需要清理的目录
	CleanupPaths []string `yaml:"cleanup_paths"`
	// 移动到的目录 只有在indirect模式下才需要设置
	TmpPath string `yaml:"tmp_path"`
}

// TmpFileCleanup 单位为小时
type TmpFileCleanup struct {
	UploadingFileExpireDuration   int `yaml:"uploading_file_expire_duration"`
	UploadingFileCleanInterval    int `yaml:"uploading_file_clean_interval"`
	CompressingFileExpireDuration int `yaml:"compressing_file_expire_duration"`
	CompressingFileCleanInterval  int `yaml:"compressing_file_clean_interval"`
}

type Quota struct {
	// 单位为秒
	StorageUsageUpdateInterval int64 `yaml:"storage_usage_update_interval"`
	// 单位为GB
	DefaultUserStorageLimit int64 `yaml:"default_user_storage_limit"`
	// 单位为GB
	MaxUserStorageLimit int64 `yaml:"max_user_storage_limit"`
	// 单位为GB
	MaxSystemStorageLimit int64 `yaml:"max_system_storage_limit"`
	// 单位为GB
	SystemStorageLimitWarningBufferSize int64 `yaml:"system_storage_limit_warning_buffer_size"`
}

type OperationLog struct {
	// 单位为秒
	OperationLogCleanInterval int64 `yaml:"operation_log_clean_interval"`
	// 单位为条
	MaxUserOperationLogCount int64 `yaml:"max_user_operation_log_count"`
	// 保留天数
	RetentionPeriod int64 `yaml:"retention_period"`
}

var Custom CustomT

func GetConfig() CustomT {
	Mut.Lock()
	defer Mut.Unlock()
	return Custom
}

// InitConfig InitConfig
func InitConfig() error {
	logger := logging.Default()
	if boot.Config.App.LoadRemoteConfig {
		switchFunc := func(data interface{}) {
			Custom = *data.(*CustomT)
			logger.Info(Custom)
		}

		err := boot.MW.RegisterRemoteConfigDefault(context.TODO(), Mut, true, boot.Config.App.Name+"/"+env.ModeName(boot.Env.Mode)+"_custom.yaml", &CustomT{}, switchFunc)
		return err
	}

	Mut.Lock()
	defer Mut.Unlock()
	md := mapstructure.Metadata{}
	err := viper.Unmarshal(&Custom, func(config *mapstructure.DecoderConfig) {
		config.TagName = "yaml"
		config.Metadata = &md
	})

	if err != nil {
		return err
	}

	return nil
}
