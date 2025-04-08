package dto

import (
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/rbac"
	"github.com/yuansuan/ticp/PSP/psp/pkg/xtype"
)

type AddRole struct {
	Name    string  `json:"name"`    // 角色名
	Comment string  `json:"comment"` // 描述
	Perms   []int64 `json:"perms"`   // 权限id
}

type ListQueryRequest struct {
	NameFilter string      `json:"name_filter"` // 过滤参数，支持传入名称
	Page       *xtype.Page `json:"page"`        // 分页参数
	Desc       bool        `json:"desc"`        // 排序方式
	OrderBy    string      `json:"order_by"`    // 排序字段，支持["name", "id"]
}

type ListQueryResponse struct {
	Role  []*RoleInfo `json:"roles"` // 角色列表
	Total int64       `json:"total"` // 总数
}

type RoleInfo struct {
	Id         int64   `json:"id"`          // 角色id
	Name       string  `json:"name"`        // 角色名
	Comment    string  `json:"comment"`     // 描述
	Type       int32   `json:"type"`        // 类型
	IsInternal bool    `json:"is_internal"` // 是否内部角色
	IsDefault  bool    `json:"is_default"`  // 是否为ldap用户登录默认角色
	Perms      []int64 `json:"perms"`       // 权限
}

type RoleDetail struct {
	Role      *RoleInfo        `json:"role"`     // 角色信息
	Resources []*Resource      `json:"has_perm"` // 权限
	Objects   []*rbac.ObjectID `json:"objects"`  // 关联角色的用户
	Perm      *Perm            `json:"perm"`
}

type Resources struct {
	Perms []*Resource `json:"perms"` // 权限
	Total int64       `json:"total"` // 总数
}

type Resource struct {
	Id           int64  `json:"id"`            // 权限id
	DisplayName  string `json:"display_name"`  // 展示字段
	Action       string `json:"action"`        // 资源操作方式 GET、POST、PUT、DELETE、NONE 默认NONE
	ResourceType string `json:"resource_type"` // 类型 system-菜单 job_sub_app-求解应用 remote_app-远程应用 api-接口 internal-内部
	ResourceName string `json:"resource_name"` // 名称
	Custom       int32  `json:"custom"`        // 1-可自定义 0-不可自定义
	ExternalId   int64  `json:"external_id"`   // 外部id
	ParentId     int64  `json:"parent_id"`     // 父级权限
}

type Perm struct {
	LocalApp       []*CustomPerm `json:"local_app"`
	CloudApp       []*CustomPerm `json:"cloud_app"`
	VisualSoftware []*CustomPerm `json:"visual_software"`
	System         []*CustomPerm `json:"system"`
}

type CustomPerm struct {
	Id         int64  `json:"id"`
	Name       string `json:"name"`
	Key        string `json:"key"`
	ExternalId string `json:"external_id"`
	Has        bool   `json:"has"`
}

type LdapDefRoleRequest struct {
	ID int64 `json:"id"` // 角色id
}
