package config

import (
	"sync"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// CustomConfig 自定义配置
type CustomConfig struct {
	TimeOut            int                 `yaml:"timeout"`
	Scheduler          *Scheduler          `yaml:"scheduler"`
	UnavailableStatus  string              `yaml:"unavailable_status"`
	HiddenNode         string              `yaml:"hidden_node"`
	SyncData           *SyncData           `yaml:"sync_data"`
	NodeClassification *NodeClassification `yaml:"node_classification"`
	HostNameMapping    *HostNameMapping    `yaml:"hostname_mapping"`
}

type Scheduler struct {
	Type                 string `yaml:"type"`
	CmdPath              string `yaml:"cmd_path"`
	MountPath            string `yaml:"mount_path"`
	ConfigPath           string `yaml:"conf_path"`
	DefaultQueue         string `yaml:"default_queue"`
	HiddenQueue          string `yaml:"hidden_queue"`
	ResAvailablePlatform string `yaml:"res_available_platform"`
}

type SyncData struct {
	Enable   bool `yaml:"enable"`
	Interval int  `yaml:"interval"`
}

type HostNameMapping struct {
	Enable bool   `yaml:"enable"`
	Path   string `yaml:"path"`
}

type NodeClassification struct {
	ClassificationRule string  `yaml:"classification_rule"`
	Nodes              []*Node `yaml:"nodes"`
}

type Node struct {
	ClassifyTag string `yaml:"classify_tag"`
	Label       string `yaml:"label"`
	Type        string `yaml:"type"`
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
