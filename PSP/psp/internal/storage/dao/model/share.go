package model

import (
	"time"

	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
)

type ShareFileRecord struct {
	Id         snowflake.ID `json:"id" xorm:"pk BIGINT(20)"`
	FilePath   string       `json:"file_path" xorm:"not null VARCHAR(255)"`
	Owner      string       `json:"owner" xorm:"not null VARCHAR(128)"`
	Type       int8         `json:"type" xorm:"not null TINYINT(1)"`
	ExpireTime time.Time    `json:"expire_time" xorm:"not null DATETIME"`
	CreateTime time.Time    `json:"create_time" xorm:"not null created"`
	UpdateTime time.Time    `json:"update_time" xorm:"not null updated"`
}

type ShareFileUser struct {
	ShareRecordId snowflake.ID `json:"share_record_id"  xorm:"not null BIGINT(20)"`
	UserId        snowflake.ID `json:"user_id" xorm:"not null BIGINT(20)"`
	State         int8         `json:"state" xorm:"not null TINYINT(1)"`
}
