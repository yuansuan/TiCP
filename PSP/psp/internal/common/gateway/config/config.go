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

// CustomConfig ...
type CustomConfig struct {
	WhiteUrlList []string `yaml:"white_url_list"`
	TokenExpire  int64    `yaml:"token_expire"`
	WhiteUrlMap  map[string]bool
}

var Mut = &sync.Mutex{}

// GetConfig GetConfig
func GetConfig() CustomConfig {
	Mut.Lock()
	defer Mut.Unlock()
	return custom
}

// SetConfig SetConfig
func SetConfig(cfg CustomConfig) {
	Mut.Lock()
	defer Mut.Unlock()

	custom = cfg
	if len(custom.WhiteUrlList) > 0 {
		whiteUrlMap := make(map[string]bool, len(custom.WhiteUrlList))
		for _, url := range custom.WhiteUrlList {
			whiteUrlMap[url] = true
		}
		custom.WhiteUrlMap = whiteUrlMap
	}
}

// InitConfig 本地调试使用
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
