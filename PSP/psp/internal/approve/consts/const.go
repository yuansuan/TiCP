package consts

const (
	JobConfig      = "JobConfig"
	JobBurstConfig = "JobBurstConfig"
)

const (
	RBACDefaultRoleId = "DefaultRoleId"
)

type approveStatus int8

const (
	Application = approveStatus(1)
	Pending     = approveStatus(2)
	Approved    = approveStatus(3)
)
