package conf_type

// Etcd config type
type Etcd struct {
	Startup        bool     `yaml:"startup"`
	Endpoints      []string `yaml:"endpoints"`
	TLS            bool     `yaml:"tls"`
	CertFile       string   `yaml:"cert_file"`
	KeyFile        string   `yaml:"key_file"`
	CAFile         string   `yaml:"ca_file"`
	DialTimeoutSec int      `yaml:"dail_timeout_sec"`
	Polling        bool     `yaml:"polling"`
}
