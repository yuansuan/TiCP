package model

import (
	"time"
)

// User user object
type User struct {
	Id            int64     `json:"id" xorm:"pk BIGINT(20)"`
	Name          string    `json:"name" xorm:"not null VARCHAR(128)"`
	Password      string    `json:"password" xorm:"VARCHAR(64)"`
	Email         string    `json:"email" xorm:"VARCHAR(128)"`
	Mobile        string    `json:"mobile" xorm:"VARCHAR(16)"`
	ApproveStatus int8      `json:"approve_status" xorm:"default '-1' TINYINT"`
	IsInternal    bool      `json:"is_internal" xorm:"default '0' TINYINT(1)"`
	Enabled       bool      `json:"enabled" xorm:"not null default '0' TINYINT(4)"`
	CreatedAt     time.Time `json:"created_at" xorm:"not null default 'CURRENT_TIMESTAMP' DATETIME"`
	UpdatedAt     time.Time `json:"updated_at" xorm:"not null default 'CURRENT_TIMESTAMP' DATETIME"`
	IsDeleted     time.Time `json:"is_deleted" xorm:"deleted DATETIME"`
	RealName      string    `json:"real_name" xorm:"not null VARCHAR(50)"`
	EnableOpenapi bool      `json:"enable_openapi"  xorm:"default '0' TINYINT(1)"`
}

// UserOrg org object
type UserOrg struct {
	Id     int64 `json:"id" xorm:"pk autoincr BIGINT(20)"`
	UserId int64 `json:"user_id" xorm:"BIGINT(20)"`
	OrgId  int64 `json:"org_id" xorm:"BIGINT(20)"`
}
