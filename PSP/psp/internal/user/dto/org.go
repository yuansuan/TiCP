package dto

type CreateOrgRequest struct {
	Name     string   `json:"name"`      // 组织结构名称
	ParentID string   `json:"parent_id"` // 父组织结构id
	Comment  string   `json:"comment"`   // 备注
	UserList []string `json:"user_list"` // 组织结构成员
}

type UpdateOrgRequest struct {
	ID       string   `json:"org_id"`    // 组织架构id
	Name     string   `json:"name"`      // 组织结构名称
	ParentID string   `json:"parent_id"` // 父组织结构id
	Comment  string   `json:"comment"`   // 备注
	UserList []string `json:"user_list"` // 组织结构成员
}

type AddOrgMemberRequest struct {
	OrgID    string   `json:"org_id"`    // 组织id
	UserList []string `json:"user_list"` // 组织结构成员
}

type DeleteOrgMemberRequest struct {
	IDs []string `json:"ids" form:"ids"` // 组织结构成员
}

type UpdateOrgMemberRequest struct {
	OrgID    string   `json:"org_id"`    // 组织id
	UserList []string `json:"user_list"` // 组织结构成员
}

type ListMemberResponse struct {
	UserID   string `json:"user_id"`
	UserName string `json:"user_name"`
}
