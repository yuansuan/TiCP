package models

import (
	"time"

	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/pkg/snowflake"
)

// UserRole ...
type UserRole struct {
	Id         snowflake.ID `json:"id" xorm:"pk default 0 comment('主键id') BIGINT(20)"`
	Name       string       `json:"name" xorm:"default '' comment('角色名称') unique(uni_name_company_id) VARCHAR(255)"`
	CompanyId  snowflake.ID `json:"company_id" xorm:"not null default 0 comment('企业ID') unique(uni_name_company_id) BIGINT(20)"`
	Type       int          `json:"type" xorm:"not null default 0 comment('1：超级 2：内置 3：自定义') INT(1)"`
	Status     int          `json:"status" xorm:"not null default 0 comment('状态（1启用 2停用）') TINYINT(1)"`
	CreateUid  snowflake.ID `json:"create_uid" xorm:"not null default 0 comment('创建者uid') BIGINT(20)"`
	CreateName string       `json:"create_name" xorm:"not null default '' comment('创建者姓名') VARCHAR(50)"`
	ModifyUid  snowflake.ID `json:"modify_uid" xorm:"not null default 0 comment('修改者uid') BIGINT(20)"`
	ModifyName string       `json:"modify_name" xorm:"not null default '' comment('修改者姓名') VARCHAR(50)"`
	UpdateTime time.Time    `json:"update_time" xorm:"default '1970-01-01 00:00:00' comment('更新时间') DATETIME"`
	CreateTime time.Time    `json:"create_time" xorm:"default CURRENT_TIMESTAMP comment('创建时间') DATETIME"`
}

// const define
const (
	// 角色状态
	// 正常
	RoleStatusNormal = 1
	// 删除
	RoleStatusDeleted = 2

	// 角色类型
	// 系统自带
	RoleTypeIsSystem = 1
	// 自定义
	RoleTypeIsCustom = 2

	// 系统自带角色名称,临时使用，未来企业可自行配置时， 可去掉
	SystemRoleNameIsSuperAdmin = "超级管理员"
	SystemRoleNameIsAdmin      = "管理员"
	SystemRoleNameIsNormal     = "普通用户"
)
