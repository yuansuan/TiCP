package model

import (
	"time"

	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
)

type ProjectMember struct {
	Id         snowflake.ID
	ProjectId  snowflake.ID
	UserId     snowflake.ID
	IsDelete   int
	LinkPath   string
	CreateTime time.Time `xorm:"create_time DATETIME created"`
	UpdateTime time.Time `xorm:"update_time DATETIME updated"`
}

const ProjectMemberTableName = "project_member"

func (p *ProjectMember) TableName() string {
	return ProjectMemberTableName
}
