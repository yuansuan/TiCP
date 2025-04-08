package config

import (
	"gopkg.in/yaml.v2"

	"github.com/yuansuan/ticp/iPaaS/standard-compute/config"
)

// LoadConfig 根据环境变量加载配置
func LoadConfig() (*config.Config, error) {
	f, err := OpenConfig("config.yml")
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()

	cfg := new(config.Config)
	if err = yaml.NewDecoder(f).Decode(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
