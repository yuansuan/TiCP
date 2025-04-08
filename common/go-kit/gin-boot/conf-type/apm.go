package conf_type

// ELastic APM Server Configure
type APM struct {
	APMStartUp   bool   `yaml:"startup"`
	APMServerURL string `yaml:"apm_server_url"`
}
