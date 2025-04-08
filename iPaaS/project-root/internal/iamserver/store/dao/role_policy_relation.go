package dao

import (
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/pkg/common/snowflake"
)

type RolePolicyRelation struct {
	ID       snowflake.ID `json:"Id,omitempty" gorm:"primary_key;column:id"`
	RoleId   snowflake.ID `json:"RoleId" gorm:"column:roleId;uniqueIndex:idx_name_role_policy"`
	PolicyId snowflake.ID `json:"PolicyId" gorm:"column:policyId;uniqueIndex:idx_name_role_policy"`
}

// TableName maps to dao table name.
func (r *RolePolicyRelation) TableName() string {
	return "role_policy_relation"
}
