package config

import (
	"sync"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

type CustomConfig struct {
	AlertManager AlertManager `yaml:"alert_manager"`
}

type AlertManager struct {
	AlertManagerConfigPath string `yaml:"alert_manager_config_path"`
	AlertManagerUrl        string `yaml:"alert_manager_url"`
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
