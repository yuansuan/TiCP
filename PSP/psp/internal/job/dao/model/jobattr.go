package model

import (
	"time"

	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
)

// JobAttr 作业相关属性
type JobAttr struct {
	JobId      snowflake.ID
	Key        string
	Value      string
	CreateTime time.Time `xorm:"create_time DATETIME created"`
	UpdateTime time.Time `xorm:"update_time DATETIME updated"`
}

func (j *JobAttr) TableName() string {
	return "job_attr"
}
