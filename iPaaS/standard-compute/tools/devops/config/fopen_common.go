//go:build !embedded
// +build !embedded

package config

import (
	"io"
	"os"
	"path/filepath"
)

// GetEnv 获取环境变量
func GetEnv() string {
	env := "local"
	if val, ok := os.LookupEnv("YS_MODE"); ok && len(val) != 0 {
		env = val
	} else if val, ok = os.LookupEnv("YS_ENV"); ok && len(val) != 0 {
		env = val
	}
	return env
}

// OpenConfig 打开配置文件
func OpenConfig(name string) (io.ReadSeekCloser, error) {
	return os.Open(filepath.Join("config", name))
}
