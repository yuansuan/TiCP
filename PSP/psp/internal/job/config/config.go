package config

import (
	"sync"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// CustomConfig 自定义配置
type CustomConfig struct {
	SyncData       *SyncData `yaml:"sync_data"`
	TempUpload     string    `yaml:"tmp_upload"`
	CloudQueue     string    `yaml:"cloud_queue"`
	WorkDir        *WorkDir  `yaml:"work_dir"`
	PlatformRegexp string    `yaml:"platform_regexp"`
}

type SyncData struct {
	Enable   bool `yaml:"enable"`
	Interval int  `yaml:"interval"`
}

type WorkDir struct {
	Type      string `yaml:"type"`
	Format    string `yaml:"format"`
	Workspace string `yaml:"workspace"`
}

var (
	mutex        sync.Mutex
	customConfig CustomConfig
)

// GetConfig 获取配置
func GetConfig() CustomConfig {
	mutex.Lock()
	defer mutex.Unlock()
	return customConfig
}

// SetConfig 设置配置
func SetConfig(config CustomConfig) {
	mutex.Lock()
	defer mutex.Unlock()
	customConfig = config
}

// InitConfig 初始化配置信息
func InitConfig() {
	mutex.Lock()
	defer mutex.Unlock()
	md := mapstructure.Metadata{}
	err := viper.Unmarshal(&customConfig, func(config *mapstructure.DecoderConfig) {
		config.TagName = "yaml"
		config.Metadata = &md
	})
	if err != nil {
		panic(errors.Wrap(err, "failed to init custom config"))
	}
}
