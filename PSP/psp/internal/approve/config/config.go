package config

import (
	"sync"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

type CustomConfig struct {
	ThreePersonManagement bool `yaml:"three_person_management_enable"`
}

var (
	mutex        sync.Mutex
	customConfig CustomConfig
)

func GetConfig() CustomConfig {
	mutex.Lock()
	defer mutex.Unlock()
	return customConfig
}

func SetConfig(config CustomConfig) {
	mutex.Lock()
	defer mutex.Unlock()
	customConfig = config
}

func InitConfig() error {
	md := mapstructure.Metadata{}
	err := viper.Unmarshal(&customConfig, func(config *mapstructure.DecoderConfig) {
		config.TagName = "yaml"
		config.Metadata = &md
	})
	if err != nil {
		return err
	}

	return nil
}
