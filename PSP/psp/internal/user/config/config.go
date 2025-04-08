package config

import (
	"sync"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

// CustomConfig CustomConfig
type CustomConfig struct {
	LdapConf Ldap    `yaml:"ldap"`
	Openapi  Openapi `yaml:"openapi"`
}

type Ldap struct {
	Enable            bool   `yaml:"enable"`
	Server            string `yaml:"server"`
	BaseDN            string `yaml:"base_dn"`
	AdminBindDn       string `yaml:"admin_bind_dn"`
	AdminBindPassword string `yaml:"admin_bind_password"`
	UID               string `yaml:"uid"`
	UserFilter        string `yaml:"userfilter"`
	Encryption        string `yaml:"encryption"`
	//ExtraAttrs ExtraAttrs `yaml:"extra_attrs"`
}

type Openapi struct {
	Enable bool `yaml:"enable"`
}

//type ExtraAttrs struct {
//	UIDKey      string `yaml:"uid_key"`
//	EmailKey    string `yaml:"email_key"`
//	RealNameKey string `yaml:"real_name_key"`
//	MobileKey   string `yaml:"mobile_key"`
//}

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

// InitConfig InitConfig
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
