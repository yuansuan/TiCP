package v20230530

// SharedDirectory 共享目录
type SharedDirectory struct {
	Path       string `json:"Path"`
	UserName   string `json:"UserName"`
	Password   string `json:"Password"`
	SharedHost string `json:"SharedHost"`
	SharedSrc  string `json:"SharedSrc"`
}
