//go:build embedded
// +build embedded

package config

import (
	"io/fs"

	"github.com/yuansuan/ticp/iPaaS/standard-compute/config"
)

var env = "dev"

func GetEnv() string {
	return env
}

// OpenConfig 打开配置文件
func OpenConfig(name string) (fs.File, error) {
	return config.FS.Open(name)
}
