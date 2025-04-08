package conf_type

import "github.com/go-pg/pg/v10"

// Pgsql Pgsql
type Pgsql struct {
	Builder func() *pg.DB
	Startup bool   `yaml:"_startup"`
	URL     string `yaml:"url"`
}

// Pgsqls Pgsqls
type Pgsqls map[string]Pgsql
