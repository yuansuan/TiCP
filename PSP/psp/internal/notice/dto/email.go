package dto

type SystemConfig struct {
	Email *EmailConfig `json:"email" yaml:"email"` // 邮件配置
}

type EmailConfig struct {
	Enable   bool        `json:"enable" yaml:"enable"`     // 是否开启
	Setting  *Setting    `json:"setting" yaml:"setting"`   // 邮件配置信息
	Template []*Template `json:"template" yaml:"template"` // 邮件模版信息
}

type Setting struct {
	Host      string `json:"host" yaml:"host"`             // 邮件服务心地址
	Port      int    `json:"port" yaml:"port"`             // 邮件服务器端口
	TLS       bool   `json:"tls" yaml:"tls"`               // 是否使用 TLS
	SendEmail string `json:"send_email" yaml:"send_email"` // 发送方邮箱地址
	Password  string `json:"password" yaml:"password"`     // 发送方邮箱密码
}

type Template struct {
	Type    string `json:"type" yaml:"type"`       // 模版类型
	Subject string `json:"subject" yaml:"subject"` // 模版邮件主题
	Content string `json:"content" yaml:"content"` // 模版邮件内容
}

type GetEmailRequest struct{}

type GetEmailResponse struct {
	Email *EmailConfig `json:"email"` // 邮件配置
}

type SetEmailRequest struct {
	Email *EmailConfig `json:"email"` // 邮件配置
}

type SetEmailResponse struct{}

type SendEmailRequest struct {
	Receiver      string `json:"receiver"`       // 邮件接收者
	EmailTemplate string `json:"email_template"` // 邮件模版类型
	JsonData      string `json:"json_data"`      // 邮件模版参数(json格式)
}

type SendEmailResponse struct{}

type TestSendEmailRequest struct {
	Receiver string `json:"receiver"` // 邮件接收者
}

type TestSendEmailResponse struct{}
