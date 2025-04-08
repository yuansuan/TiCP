package conf_type

import (
	"database/sql"
	"time"
)

// Mysql Mysql
type Mysql struct {
	Builder   func() *sql.DB
	Startup   bool `yaml:"_startup"`
	Encrypt   bool
	Dsn       string
	HiddenSQL bool `yaml:"hidden_sql"`

	MaxIdleConnection int           `yaml:"max_idle_connection"`
	MaxOpenConnection int           `yaml:"max_open_connection"`
	MaxIdleTime       time.Duration `yaml:"max_idle_time"`
}

// Mysqls Mysqls
type Mysqls map[string]Mysql
