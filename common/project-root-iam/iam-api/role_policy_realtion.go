package iam_api

type AddRolePolicyRelationRequest struct {
	RoleName   string `json:"RoleName"`
	PolicyName string `json:"PolicyName"`
}

type DeleteRolePolicyRelationRequest struct {
	RoleName   string `json:"RoleName"`
	PolicyName string `json:"PolicyName"`
}

type AdminAddRolePolicyRelationRequest struct {
	AddRolePolicyRelationRequest
	UserId string `json:"UserId"`
}

type AdminDeleteRolePolicyRelationRequest struct {
	DeleteRolePolicyRelationRequest
	UserId string `json:"UserId"`
}

type ListPolicyByRoleNameRequest struct {
	RoleName string `json:"RoleName"`
}

type AdminListPolicyByRoleNameRequest struct {
	ListPolicyByRoleNameRequest
	UserId string `json:"UserId"`
}
