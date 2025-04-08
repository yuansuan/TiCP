package config

import (
	"context"
	"sync"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
)

// CustomT CustomT
type CustomT struct {
	AccessKeyId     string `yaml:"access_key_id"`
	AccessKeySecret string `yaml:"access_key_secret"`
}

// switch storage
var switches map[string]bool

// Mut Mut
var Mut = &sync.Mutex{}

// Custom Custom
var Custom CustomT

// CheckSwitch CheckSwitch
func CheckSwitch(key string) bool {
	Mut.Lock()
	defer Mut.Unlock()
	if v, ok := switches[key]; v && ok {
		return true
	}
	return false
}

// InitConfig InitConfig
func InitConfig() error {
	if boot.Config.App.LoadRemoteConfig {
		switchFunc := func(data interface{}) {
			Custom = *data.(*CustomT)
			logging.Default().Info(Custom)
		}

		err := boot.MW.RegisterRemoteConfigDefault(context.TODO(), Mut, true,
			common.GetRemoteConfigPath(), &CustomT{}, switchFunc)
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

// GetCustom ...
func GetCustom() CustomT {
	Mut.Lock()
	defer Mut.Unlock()
	return Custom
}
