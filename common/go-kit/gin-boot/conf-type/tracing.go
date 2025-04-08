package conf_type

type Tracing struct {
	Startup bool `json:"startup" yaml:"startup"`
	Details struct {
		Enabled  bool `json:"enabled" yaml:"enabled"`
		Request  bool `json:"request" yaml:"request"`
		Response bool `json:"response" yaml:"response"`
	} `json:"details" yaml:"details"`
	Database struct {
		Enabled  bool `json:"enabled" yaml:"enabled"`
		Binding  bool `json:"binding" yaml:"binding"`
		Dangling bool `json:"dangling" yaml:"dangling"`
	} `json:"database" yaml:"database"`
	Http struct {
		Excludes []string `json:"excludes" yaml:"excludes"`
	} `json:"http" yaml:"http"`
	Jaeger struct {
		Endpoint string `json:"endpoint" yaml:"endpoint"`
	} `json:"jaeger" yaml:"jaeger"`
}
