package conf_type

// SqLite 数据库相关配置
type SqLite struct {
	Startup   bool `yaml:"_startup"`
	Encrypt   bool
	Dsn       string
	HiddenSQL bool `yaml:"hidden_sql"`

	MaxIdleConnection int `yaml:"max_idle_connection"`
	MaxOpenConnection int `yaml:"max_open_connection"`
}

// Sqlites 数据库组配置别名
type Sqlites map[string]SqLite
