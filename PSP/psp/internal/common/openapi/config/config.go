package config

import (
	"path/filepath"
	"runtime"
	"sync"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/env"

	"github.com/yuansuan/ticp/PSP/psp/internal/common"
)

var custom CustomConfig

type CustomConfig struct {
	Local *Local `yaml:"local"`
}

type Local struct {
	Settings *Settings `yaml:"settings"`
}

type Settings struct {
	AppKey      string `yaml:"app_key"`
	AppSecret   string `yaml:"app_secret"`
	Endpoint    string `yaml:"api_endpoint"`
	HPCEndpoint string `yaml:"hpc_endpoint"`
	UserId      string `yaml:"user_id"`
	Zone        string `yaml:"zone"`
}

var Mut = &sync.Mutex{}

func GetConfig() CustomConfig {
	Mut.Lock()
	defer Mut.Unlock()
	return custom
}

func SetConfig(cfg CustomConfig) {
	if custom.Local != nil && custom.Local.Settings.Endpoint != "" {
		return
	}

	Mut.Lock()
	defer Mut.Unlock()
	if custom.Local != nil && custom.Local.Settings.Endpoint != "" {
		return
	}
	custom = cfg
}

func InitConfig() {
	// 如果是线上方式启动,return
	if env.ModeName(boot.Env.Mode) == common.SysEnvProd {
		return
	}

	_, filename, _, _ := runtime.Caller(0)
	pkgPath := filepath.Dir(filename)

	viper.SetConfigType("yaml")
	viper.SetConfigFile(pkgPath + "/dev_custom.yml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	var cfg CustomConfig

	md := mapstructure.Metadata{}
	err = viper.Unmarshal(&cfg, func(config *mapstructure.DecoderConfig) {
		config.TagName = "yaml"
		config.Metadata = &md
	})
	if err != nil {
		panic(err)
	}

	SetConfig(cfg)
}
