package model

import (
	"time"

	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
)

type AlertNotification struct {
	Id         snowflake.ID `xorm:"id BIGINT(20) pk"`
	Key        string       `xorm:"key VARCHAR(255)"`
	Value      string       `xorm:"value VARCHAR(255)"`
	Type       string       `xorm:"type VARCHAR(16)"`
	CreateTime time.Time    `xorm:"create_time DATETIME created"`
	UpdateTime time.Time    `xorm:"update_time DATETIME updated"`
}

const AlertNotificationTableName = "alert_notification"

func (AlertNotification) TableName() string {
	return AlertNotificationTableName
}
