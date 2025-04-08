package conf_type

// HTTP 其中保存有关HTTP服务的配置选型
type HTTP struct {
	Logger struct {
		// 不需要进行日志记录的路径列表
		Excludes []string `json:"excludes" yaml:"excludes"`
	} `json:"logger" yaml:"logger"`
}
