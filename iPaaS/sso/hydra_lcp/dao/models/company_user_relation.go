package models

import (
	"time"

	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/pkg/snowflake"
)

// CompanyUserRelation ...
type CompanyUserRelation struct {
	Id         snowflake.ID `json:"id" xorm:"pk default 0 comment('主键id') BIGINT(20)"`
	CompanyId  snowflake.ID `json:"company_id" xorm:"not null pk default 0 comment('企业ID') unique(uni_company_user_id) BIGINT(20)"`
	UserId     snowflake.ID `json:"user_id" xorm:"not null pk default 0 comment('sso用户id') unique(uni_company_user_id) BIGINT(20)"`
	Status     int          `json:"status" xorm:"default 1 comment('用户状态：1正常； 2删除') TINYINT(1)"`
	UpdateTime time.Time    `json:"update_time" xorm:"default '1970-01-01 00:00:00' comment('更新时间') DATETIME"`
	CreateTime time.Time    `json:"create_time" xorm:"default CURRENT_TIMESTAMP comment('创建时间') DATETIME"`
}

// const define
const (
	CompanyUserRelationStatusNormal  = 1
	CompanyUserRelationStatusDeleted = 2
)

// CompanyUserModel 这是一个虚的的企业用户Model, UserModel 和 CompanyUserRelation的合并
type CompanyUserModel struct {
	UserId        snowflake.ID `json:"user_id" `         // sso用户id
	CompanyId     snowflake.ID `json:"company_id"  `     // 企业ID
	Phone         string       `json:"phone" `           // 用户电话
	UserName      string       `json:"User_name"`        // 用户名
	RealName      string       `json:"real_name"`        // 用户姓名
	Email         string       `json:"email" `           // email
	RoleList      []*UserRole  `json:"role_list"`        // 用户角色信息
	UpdateTime    time.Time    `json:"update_time" `     // 更新时间
	CreateTime    time.Time    `json:"create_time" `     // 创建时间
	Status        int32        `json:"status" `          // 用户状态：1正常; 2删除;
	Role          []*UserRole  `json:"role"`             // 用户角色
	LastLoginTime time.Time    `json:"last_login_time" ` // 最后登录时间
}

type CompanyUserConfig struct {
	Id         snowflake.ID `json:"id" xorm:"pk default 0 comment('主键id') BIGINT(20)"`
	RelationId snowflake.ID `json:"relation_id" xorm:"not null default 0 comment('关系ID') BIGINT(20)"`
	Key        string       `json:"key" xorm:"not null default '' comment('key') VARCHAR(256)"`
	Value      string       `json:"value" xorm:"not null default '' comment('value') VARCHAR(256)"`
	IsDeleted  bool         `json:"is_deleted" xorm:"not null default 0 comment('是否已被删除,') TINYINT(4)"`
	UpdateTime time.Time    `json:"update_time" xorm:"default '1970-01-01 00:00:00' comment('更新时间') DATETIME"`
	CreateTime time.Time    `json:"create_time" xorm:"default CURRENT_TIMESTAMP comment('创建时间') DATETIME"`
}
