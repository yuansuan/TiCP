package config

import (
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type Server struct {
	Addr string `json:"addr" yaml:"addr"`
}

type Logger struct {
	Level string `json:"level" yaml:"level"`
}

type TurnServer struct {
	Uri    string `json:"uri" yaml:"uri"`
	Secret string `json:"secret" yaml:"secret"`
	Expire int64  `json:"expire" yaml:"expire"`
}

type Config struct {
	Server     Server        `json:"server" yaml:"server"`
	Logger     Logger        `json:"logger" yaml:"logger"`
	TurnServer []*TurnServer `json:"turn_server" yaml:"turn_server"`
}

func Read(filename string) (*Config, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, errors.Wrap(err, "config")
	}

	var cfg Config
	if err = yaml.Unmarshal(content, &cfg); err != nil {
		return nil, errors.Wrap(err, "config")
	}
	return override(&cfg), nil
}

func override(cfg *Config) *Config {
	if addr, ok := os.LookupEnv("SERVER_ADDR"); ok && len(addr) != 0 {
		cfg.Server.Addr = addr
	}
	if level, ok := os.LookupEnv("LOGGER_LEVEL"); ok && len(level) != 0 {
		cfg.Logger.Level = level
	}

	return cfg
}
