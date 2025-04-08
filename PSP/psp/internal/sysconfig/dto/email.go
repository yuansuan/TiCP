package dto

type GetEmailConfigRes struct {
	// Notification is job event config
	Notification Notification `json:"notification,omitempty"`
}

type SetEmailConfigReq struct {
	// Notification is job event config
	Notification Notification `json:"notification,omitempty"`
}

type SetEmailConfigRes struct {
}

type SendEmailTestResponse struct {
}

type Notification struct {
	// NodeBreakdown event
	NodeBreakdown bool `json:"node_breakdown"`

	// DiskUsage event
	DiskUsage bool `json:"disk_usage"`

	// AgentBreakdown event
	AgentBreakdown bool `json:"agent_breakdown"`

	// JobFailNum event
	JobFailNum bool `json:"job_fail_num"`
}

type SetGlobalEmailRequest struct {
	EmailConfig *EmailConfig `json:"email_config"` // 邮件配置
}
type SetGlobalEmailResponse struct {
}

type GetGlobalEmailRes struct {
	// Notification is job event config
	EmailConfig *EmailConfig `json:"email_config,omitempty"`
}

// EmailConfig is the email server config
type EmailConfig struct {

	// Host is the email server name
	Host string `json:"host"`

	// Port is the mail server port
	Port int `json:"port"`

	// UseTLS means whether or not use TLS connection
	UseTLS bool `json:"use_tls"`

	// UserName is the email address
	UserName string `json:"user_name"`

	// Password is the email password
	Password string `json:"password"`

	// From is the email from value
	From string `json:"from"`

	// AdminAddr means the admin email address
	AdminAddr string `json:"admin_addr"`
}
