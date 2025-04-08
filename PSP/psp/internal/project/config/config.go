package config

import (
	"sync"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// CustomConfig 自定义配置
type CustomConfig struct {
	ProjectCheck      *ProjectEndCheck   `yaml:"project_check"`
	SelectorListLimit *SelectorListLimit `yaml:"selector_list_limit"`
}

type ProjectEndCheck struct {
	Enable    bool `yaml:"enable"`
	Interval  int  `yaml:"interval"`
	DailyTime int  `yaml:"daily_time"`
	MinTime   int  `yaml:"min_time"`
}

type SelectorListLimit struct {
	Enable    bool `yaml:"enable"`
	MaxMonths int  `yaml:"max_months"`
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
