package model

import (
	"time"

	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
)

type Session struct {
	ID            snowflake.ID `xorm:"id BIGINT(20) pk"`
	OutSessionID  string       `xorm:"out_session_id VARCHAR(255)"`
	HardwareID    snowflake.ID `xorm:"hardware_id BIGINT(20)"`
	OutHardwareID string       `xorm:"out_hardware_id VARCHAR(255)"`
	SoftwareID    snowflake.ID `xorm:"software_id BIGINT(20)"`
	OutSoftwareID string       `xorm:"out_software_id VARCHAR(255)"`
	UserID        snowflake.ID `xorm:"user_id BIGINT(20)"`
	UserName      string       `xorm:"user_name VARCHAR(64)"`
	ProjectID     snowflake.ID `xorm:"project_id BIGINT(20)"`
	ProjectName   string       `xorm:"project_name VARCHAR(255)"`
	RawStatus     string       `xorm:"raw_status VARCHAR(64)"`
	Status        string       `xorm:"status VARCHAR(64)"`
	StreamURL     string       `xorm:"stream_url TEXT"`
	ExitReason    string       `xorm:"exit_reason VARCHAR(255)"`
	Duration      int64        `xorm:"duration BIGINT(20)"`
	Zone          string       `xorm:"zone VARCHAR(64)"`
	IsAutoClose   bool         `xorm:"is_auto_close TINYINT(4)"`
	Deleted       time.Time    `xorm:"deleted DATETIME deleted"`
	StartTime     time.Time    `xorm:"start_time DATETIME"`
	EndTime       time.Time    `xorm:"end_time DATETIME"`
	CreateTime    time.Time    `xorm:"create_time DATETIME created"`
	UpdateTime    time.Time    `xorm:"update_time DATETIME updated"`
}

const SessionTableName = "visual_session"

func (s *Session) TableName() string {
	return SessionTableName
}
