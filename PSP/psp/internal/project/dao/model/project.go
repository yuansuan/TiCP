package model

import (
	"time"

	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
)

type Project struct {
	Id           snowflake.ID
	ProjectName  string
	ProjectOwner snowflake.ID
	State        string
	StartTime    time.Time
	EndTime      time.Time
	Comment      string
	FilePath     string
	IsDelete     int
	CreateTime   time.Time `xorm:"create_time DATETIME created"`
	UpdateTime   time.Time `xorm:"update_time DATETIME updated"`
}

const ProjectTableName = "project"

func (p *Project) TableName() string {
	return ProjectTableName
}
