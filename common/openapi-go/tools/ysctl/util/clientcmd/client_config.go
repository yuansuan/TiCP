package clientcmd

import "io"

type Config struct {
	AccessKeyID     string `yaml:"access_key_id" mapstructure:"access_key_id, omitempty"`
	AccessKeySecret string `yaml:"access_key_secret" mapstructure:"access_key_secret, omitempty"`
	Endpoint        string `yaml:"endpoint" mapstructure:"endpoint, omitempty"`
	StorageEndpoint string `yaml:"storage_endpoint" mapstructure:"storage_endpoint, omitempty"`
}

func NewConfig() *Config {
	return &Config{}
}

type IOStreams struct {
	In     io.Reader
	Out    io.Writer
	ErrOut io.Writer
}
