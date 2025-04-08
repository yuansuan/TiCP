package iam_api

type AddRoleRequest struct {
	RoleName    string      `json:"RoleName"`
	Description string      `json:"Description"`
	TrustPolicy *RolePolicy `json:"TrustPolicy"`
}

type AddRoleResponse struct{}

type DeleteRoleRequest struct {
	RoleName string `json:"RoleName"`
}

type GetRoleRequest struct {
	RoleName string `json:"RoleName"`
}

type ListRoleResponse struct {
	Roles []*Role `json:"Roles"`
}

type Role struct {
	RoleName    string      `json:"RoleName"`
	Description string      `json:"Description"`
	TrustPolicy *RolePolicy `json:"TrustPolicy"`
}

type RolePolicy struct {
	Actions    []string `json:"Actions"`
	Effect     string   `json:"Effect"`
	Principals []string `json:"Principals"`
	Resources  []string `json:"Resources"`
}

type UpdateRoleRequest struct {
	Role *Role `json:"Role"`
}

type GetRoleResponse AddRoleRequest

type RoleStatement struct {
	Effect string `json:"Effect"`
	// like: "RAM": ["acs:ram::1904037259021276:root"]
	Principal map[string][]string `json:"Principal"`
	Action    string              `json:"Action"`
}

type TrustPolicy struct {
	Version   string           `json:"Version"`
	Statement []*RoleStatement `json:"Statement"`
}

type AdminAddRoleRequest struct {
	AddRoleRequest
	UserId string `json:"UserId"`
}

type AdminGetRoleRequest struct {
	UserId   string `json:"UserId"`
	RoleName string `json:"RoleName"`
}

type AdminGetRoleResponse struct {
	Role *Role `json:"Role"`
}

type AdminListRoleRequest struct {
	UserId string `json:"UserId"`
}

type AdminListRoleResponse struct {
	Roles []*Role `json:"Roles"`
}

type AdminDeleteRoleRequest struct {
	UserId   string `json:"UserId"`
	RoleName string `json:"RoleName"`
}

type AdminUpdateRoleRequest struct {
	Role   *Role  `json:"Role"`
	UserId string `json:"UserId"`
}
