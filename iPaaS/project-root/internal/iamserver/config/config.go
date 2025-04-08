package config

import (
	"context"
	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"sync"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

// Custom Custom
var Custom CustomT

type CustomT struct {
	AdminAKID     string     `yaml:"admin_ak_id"`
	AdminAKSECRET string     `yaml:"admin_ak_secret"`
	ApiServerT    ApiServerT `yaml:"ys_api_server"`
	Database      Database   `yaml:"database"`
	AllowAddUsers []string   `yaml:"allow_add_users"`
	// 角色扮演最长时间，单位小时，默认4小时
	MaxAssumeRoleTime int64 `yaml:"max_assume_role_time"`
	// 清理审计日志的时间， 单位天，默认3天
	CleanAuditLogInterval int64 `yaml:"clean_audit_log_interval"`
}

type Database struct {
	Host     string `yaml:"host"`
	UserName string `yaml:"username"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

type ApiServerT struct {
	Job         string `yaml:"job"`
	AccountBill string `yaml:"account_bill"`
	CloudApp    string `yaml:"cloud_app"`
	LicManager  string `yaml:"lic_manager"`
	Merchandise string `yaml:"merchandise"`
}

// Mut Mut
var Mut = &sync.Mutex{}

// GetConfig return the config
func GetConfig() CustomT {
	Mut.Lock()
	defer Mut.Unlock()
	return Custom
}

// InitConfig InitConfig
func InitConfig() error {
	logger := logging.Default()
	if boot.Config.App.LoadRemoteConfig {
		switchFunc := func(data interface{}) {
			Custom = *data.(*CustomT)
			logger.Info(Custom)
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
