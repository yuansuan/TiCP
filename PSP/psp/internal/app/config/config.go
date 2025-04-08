package config

import (
	"sync"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

type CustomConfig struct {
	SchedulerResourcePath string `yaml:"scheduler_resource_path"`
}

var Mut = &sync.Mutex{}

var custom CustomConfig

func GetConfig() CustomConfig {
	Mut.Lock()
	defer Mut.Unlock()
	return custom
}

func SetConfig(cfg CustomConfig) {
	Mut.Lock()
	defer Mut.Unlock()
	custom = cfg
}

func InitConfig() error {
	Mut.Lock()
	defer Mut.Unlock()
	md := mapstructure.Metadata{}
	err := viper.Unmarshal(&custom, func(config *mapstructure.DecoderConfig) {
		config.TagName = "yaml"
		config.Metadata = &md
	})
	if err != nil {
		return err
	}
	return nil
}
