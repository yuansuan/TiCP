package config

import (
	"sync"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// CustomConfig CustomConfig
type CustomConfig struct {
	FilterHideFileRegex string          `yaml:"filter_hide_file_regex"`
	WhitePathList       []string        `yaml:"white_path_list"`
	LocalRootPath       string          `yaml:"local_root_path"`
	HpcUploadConfig     hpcUploadConfig `yaml:"hpc_upload_config"`
	OnlyReadPathList    []string        `yaml:"only_read_path_list"`
	PublicFolderEnable  bool            `yaml:"public_folder_enable"`
}

type hpcUploadConfig struct {
	BlockSize        int `yaml:"block_size"`
	ConcurrencyLimit int `yaml:"concurrency_limit"`
	RetryCount       int `yaml:"retry_count"`
	RetryDelay       int `yaml:"retry_delay"`
	WaitResumeTime   int `yaml:"wait_resume_time"`
}

// Custom Custom
var custom CustomConfig

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
}

// InitConfig InitConfig
func InitConfig() {
	md := mapstructure.Metadata{}
	err := viper.Unmarshal(&custom, func(config *mapstructure.DecoderConfig) {
		config.TagName = "yaml"
		config.Metadata = &md
	})
	if err != nil {
		panic(errors.Wrap(err, "failed to init custom config"))
	}

}
