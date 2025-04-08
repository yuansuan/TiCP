package conf_type

// Temporal Server Configure
type Temporal struct {
	TemporalStartUp bool   `yaml:"_startup"`
	TemporalHost    string `yaml:"host"`
	NameSpace       string `yaml:"namespace"`
}
