package config

import (
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// CustomConfig CustomConfig
type CustomConfig struct {
	RBACConfigPath     string `yaml:"rbac_config_path"`
	EnableApiAuthorize bool   `yaml:"enable_api_authorize"`
}

// Custom Custom
var Custom CustomConfig

// InitConfig InitConfig
func InitConfig() {
	md := mapstructure.Metadata{}
	err := viper.Unmarshal(&Custom, func(config *mapstructure.DecoderConfig) {
		config.TagName = "yaml"
		config.Metadata = &md
	})
	if err != nil {
		panic(errors.Wrap(err, "failed to init custom config"))
	}

}
