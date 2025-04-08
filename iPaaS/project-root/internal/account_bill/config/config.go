package config

import (
	"context"
	"sync"

	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/env"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"gopkg.in/yaml.v2"
)

type CustomT struct {
}

// Init ...
func (c *CustomT) Init() error {
	return nil
}

// String String
func (c CustomT) String() string {
	bs, _ := yaml.Marshal(&c)

	return string(bs)
}

// Mut Mut
var Mut = &sync.Mutex{}

// Custom ...
var Custom CustomT

// InitConfig InitConfig
func InitConfig() error {
	logger := logging.Default()

	if boot.Config.App.LoadRemoteConfig {
		switchFunc := func(data interface{}) {
			Custom = *data.(*CustomT)
			err := Custom.Init()
			if err != nil {
				logger.Error("err_init_custom_config_failed", "error message:", err)
				return
			}
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

	if err := Custom.Init(); err != nil {
		logger.Error("err_init_custom_config_failed", "error message:", err)
		return err
	}
	return nil
}

// GetCustom ...
func GetCustom() CustomT {
	Mut.Lock()
	defer Mut.Unlock()
	return Custom
}
