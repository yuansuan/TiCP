package config

import (
	"context"
	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/env"
	"sync"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

// CustomT CustomT
type CustomT struct {
	NodeID int `yaml:"node_id"`
}

// switch storage
var switches map[string]bool

// Mut Mut
var Mut = &sync.Mutex{}

// Custom Custom
var Custom CustomT

// GetConfig GetConfig
func GetConfig() CustomT {
	Mut.Lock()
	defer Mut.Unlock()
	return Custom
}

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
		err := boot.MW.RegisterRemoteConfigDefault(context.TODO(), Mut, true, boot.Config.App.Name+"/"+env.ModeName(boot.Env.Mode)+"_custom.yaml", &Custom, nil)
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
