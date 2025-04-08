package consts

const (
	JobConfig      = "JobConfig"
	JobBurstConfig = "JobBurstConfig"
)

const (
	RBACDefaultRoleId   = "DefaultRoleId"
	RBACDefaultSafeUser = "DefaultSafeUser"
)

const (
	NodeBreakdown  = "nodebreakdown"
	AgentBreakdown = "agentbreakdown"
	DiskUsage      = "diskusage"
	JobFailNum     = "jobfailnum"

	Enable = "enable"
)
const (
	// FileType is yaml
	FileType = "yml"
	// AlertManagerName is alertmanager
	AlertManagerName = "alertmanager"
)
const (
	KeyNodeBreakdown  = "node_breakdown"
	KeyDiskUsage      = "disk_usage"
	KeyAgentBreakdown = "agent_breakdown"
	KeyJobFailNum     = "job_fail_num"

	KeyHost      = "host"
	KeyPort      = "port"
	KeyUseTLS    = "use_tls"
	KeyPassword  = "password"
	KeyUsername  = "user_name"
	KeyFrom      = "from"
	KeyAdminAddr = "admin_addr"

	KeyCommon = "common"
)

const (
	AlertManagerType = "AlertManager"
	GlobalEmailType  = "GlobalEmail"
)

const (
	TestEmailSubject = "测试发送邮件"
	TestEmailBody    = "您好！邮件发送成功!"
)
