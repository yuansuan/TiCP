package conf_type

// Logger Logger
type Logger struct {
	LogDir     string `yaml:"log_dir"`
	UseFile    bool   `yaml:"use_file"`
	MaxSize    int    `yaml:"max_size"` // MB
	MaxAge     int    `yaml:"max_age"`  // DAY
	MaxBackups int    `yaml:"max_backups"`
}
