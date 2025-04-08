package errorcode

const (
	InvalidArgument             = "InvalidArgument"
	InvalidCustomStateRule      = "InvalidCustomStateRule"
	InvalidJobID                = "InvalidJobID"
	InvalidPageOffset           = "InvalidPageOffset"
	InvalidPageSize             = "InvalidPageSize"
	InvalidArgumentCoresPerNode = "InvalidArgument.CoresPerNode"

	InternalServerError         = "InternalServerError"
	DatabaseInternalServerError = "DatabaseInternalServerError"

	WrongCPUUsage = "WrongCPUUsage"

	JobNotFound = "JobNotFound"

	CancelJobForbidden = "CancelJobForbidden"
	DeleteJobForbidden = "DeleteJobForbidden"
	JobNotRunning      = "JobNotRunning"
)
