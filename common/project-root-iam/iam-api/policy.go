package iam_api

type AddPolicyRequest struct {
	PolicyName string   `json:"PolicyName"`
	Version    string   `json:"Version"`
	Effect     string   `json:"Effect"`
	Resources  []string `json:"Resources"`
	Actions    []string `json:"Actions"`
}

type UpdatePolicyRequest struct {
	Policy AddPolicyRequest `json:"Policy"`
}

type PatchPolicyRequest struct {
	RoleName   string `json:"RoleName"`
	PolicyName string `json:"PolicyName"`
}

type DetachPolicyRequest struct {
	RoleName   string `json:"RoleName"`
	PolicyName string `json:"PolicyName"`
}

type GetPolicyRequest struct {
	PolicyName string `json:"PolicyName"`
}

type GetPolicyResponse struct {
	PolicyName string   `json:"PolicyName"`
	Effect     string   `json:"Effect"`
	Resources  []string `json:"Resources"`
	Actions    []string `json:"Actions"`
}

type ListPolicyResponse struct {
	Policies []*GetPolicyResponse
}

type DeletePolicyRequest struct {
	PolicyName string `json:"PolicyName"`
}

type AdminGetPolicyRequest struct {
	PolicyName string `json:"PolicyName"`
	UserId     string `json:"UserId"`
}

type AdminGetPolicyResponse GetPolicyResponse

type AdminAddPolicyRequest struct {
	AddPolicyRequest
	UserId string `json:"UserId"`
}

type AdminUpdatePolicyRequest struct {
	Policy AddPolicyRequest `json:"Policy"`
	UserId string           `json:"UserId"`
}

type AdminListPolicyRequest struct {
	UserId string `json:"UserId"`
}

type AdminListPolicyResponse struct {
	Policies []*GetPolicyResponse `json:"Policies"`
}

type AdminDeletePolicyRequest struct {
	PolicyName string `json:"PolicyName"`
	UserId     string `json:"UserId"`
}

type AdminPatchPolicyRequest struct {
	PatchPolicyRequest
	UserId string `json:"UserId"`
}
type AdminDetachPolicyRequest struct {
	DetachPolicyRequest
	UserId string `json:"UserId"`
}
