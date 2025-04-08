package models

import (
	"time"
)

// SsoExternalUser ...
type SsoExternalUser struct {
	Ysid       int64     `json:"ysid" xorm:"pk default 0 comment('主键ID') BIGINT(20)"`
	UserId     int64     `json:"user_id" xorm:"not null default 0 BIGINT(20)"`
	UserName   string    `json:"user_name" xorm:"not null index(uniq_user_name)  VARCHAR(100)"`
	CreateTime time.Time `json:"create_time" xorm:"not null default CURRENT_TIMESTAMP comment('创建时间') DATETIME"`
}
