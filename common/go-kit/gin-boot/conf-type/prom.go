package conf_type

// Monitor Monitor
type Monitor struct {
	ListenAddr               string `yaml:"listen"`
	StartUp                  bool   `yaml:"_startup"`
	MetricPath               string `yaml:"metric"`
	PrometheusServerEndpoint string `yaml:"prometheus_server_endpoint"`
}
