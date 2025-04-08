package conf_type

// Elasticsearch config type
type Elasticsearch struct {
	Startup   bool     `yaml:"_startup"`
	Addresses []string `yaml:"addresses"`
	Username  string   `yaml:"username"`
	Password  string   `yaml:"password"`
}

// Elasticsearches Elasticsearches
type Elasticsearches map[string]Elasticsearch
