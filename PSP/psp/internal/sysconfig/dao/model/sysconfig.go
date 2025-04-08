package model

import (
	"time"

	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
)

type SysConfig struct {
	Id         snowflake.ID `xorm:"id BIGINT(20) pk"`
	Key        string       `xorm:"key VARCHAR(255)"`
	Value      string       `xorm:"value TEXT"`
	CreateTime time.Time    `xorm:"create_time DATETIME created"`
	UpdateTime time.Time    `xorm:"update_time DATETIME updated"`
}

const SysConfigTableName = "sys_config"

func (SysConfig) TableName() string {
	return SysConfigTableName
}
