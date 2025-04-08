package config

import (
	"sync"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

type CustomConfig struct {
	Zone                     string                    `yaml:"zone"`
	SyncData                 *SyncData                 `yaml:"sync_data"`
	MountDirectory           *MountDirectory           `yaml:"mount_directory"`
	SessionNotificationCheck *SessionNotificationCheck `yaml:"session_notification_check"`
	Local                    bool                      `yaml:"local"`
}

type SessionNotificationCheck struct {
	Interval  int `yaml:"interval"`
	DailyTime int `yaml:"daily_time"`
	MinTime   int `yaml:"min_time"`
}

type SyncData struct {
	Enable         bool `yaml:"enable"`
	DataInterval   int  `yaml:"data_interval"`
	StatusInterval int  `yaml:"status_interval"`
}

type MountDirectory struct {
	LimitNum              int      `yaml:"limit_num"`
	DriveNames            []string `yaml:"drive_names"`
	LinuxMountRootPath    string   `yaml:"linux_mount_root_path"`
	EnablePublicDirectory bool     `yaml:"enable_public_directory"`
}

var Mut = &sync.Mutex{}

var custom CustomConfig

func GetConfig() CustomConfig {
	Mut.Lock()
	defer Mut.Unlock()
	return custom
}

func SetConfig(cfg CustomConfig) {
	Mut.Lock()
	defer Mut.Unlock()
	custom = cfg
}

func InitConfig() error {
	Mut.Lock()
	defer Mut.Unlock()
	md := mapstructure.Metadata{}
	err := viper.Unmarshal(&custom, func(config *mapstructure.DecoderConfig) {
		config.TagName = "yaml"
		config.Metadata = &md
	})
	if err != nil {
		return err
	}
	return nil
}

func GetZone() string {
	zone := GetConfig().Zone
	if zone == "" {
		return "az-shanghai"
	}
	return zone
}
