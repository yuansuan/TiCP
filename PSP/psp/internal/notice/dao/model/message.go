package model

import (
	"time"

	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
)

type Message struct {
	Id         snowflake.ID
	UserId     snowflake.ID
	State      int
	Type       string
	Content    string
	CreateTime time.Time `xorm:"create_time DATETIME created"`
	UpdateTime time.Time `xorm:"update_time DATETIME updated"`
}

func (m *Message) TableName() string {
	return "notice_message"
}
