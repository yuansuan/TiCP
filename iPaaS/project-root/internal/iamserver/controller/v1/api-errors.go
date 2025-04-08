package v1

// APIErrorCode type of error status.
type APIErrorCode int

const (
	InvalidArgument = "InvalidArgument"

	PermissionDenied = "PermissionDenied"

	AlreadyExists = "AlreadyExists"

	AccessKeyNotFound = "AccessKeyNotFound"

	UserNotFound = "UserNotFound"

	RoleNotFound = "RoleNotFound"

	PolicyNotFound = "PolicyNotFound"

	RelationNotFound    = "RelationNotFound"
	InternalServerError = "InternalServerError"
)
