package model

import (
	"time"

	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
)

type Software struct {
	ID            snowflake.ID `xorm:"id BIGINT(20) pk"`
	OutSoftwareID string       `xorm:"out_software_id VARCHAR(255)"`
	Name          string       `xorm:"name VARCHAR(64)"`
	Desc          string       `xorm:"desc VARCHAR(255)"`
	Platform      string       `xorm:"platform VARCHAR(64)"`
	ImageID       string       `xorm:"image_id VARCHAR(255)"`
	State         string       `xorm:"state VARCHAR(64)"`
	InitScript    string       `xorm:"init_script TEXT"`
	Icon          string       `xorm:"icon MEDIUMTEXT"`
	GpuDesired    bool         `xorm:"gpu_desired TINYINT(4)"`
	Zone          string       `xorm:"zone VARCHAR(64)"`
	Deleted       time.Time    `xorm:"deleted DATETIME deleted"`
	CreateTime    time.Time    `xorm:"create_time DATETIME created"`
	UpdateTime    time.Time    `xorm:"update_time DATETIME updated"`
}

const SoftwareTableName = "visual_software"

func (s *Software) TableName() string {
	return SoftwareTableName
}

type SoftwarePreset struct {
	ID         snowflake.ID `xorm:"id BIGINT(20) pk"`
	SoftwareID snowflake.ID `xorm:"software_id BIGINT(20)"`
	HardwareID snowflake.ID `xorm:"hardware_id BIGINT(20)"`
	Defaulted  bool         `xorm:"defaulted TINYINT(4)"`
	CreateTime time.Time    `xorm:"create_time DATETIME created"`
	UpdateTime time.Time    `xorm:"update_time DATETIME updated"`
}

const SoftwarePresetTableName = "visual_software_preset"

func (s *SoftwarePreset) TableName() string {
	return SoftwarePresetTableName
}

type RemoteApp struct {
	ID             snowflake.ID `xorm:"id BIGINT(20) pk"`
	OutRemoteAppID string       `xorm:"out_remote_app_id VARCHAR(255)"`
	SoftwareID     snowflake.ID `xorm:"software_id BIGINT(20)"`
	OutSoftwareID  string       `xorm:"out_software_id VARCHAR(255)"`
	Name           string       `xorm:"name VARCHAR(64)"`
	Desc           string       `xorm:"desc VARCHAR(255)"`
	Dir            string       `xorm:"dir VARCHAR(255)"`
	Args           string       `xorm:"args VARCHAR(255)"`
	Logo           string       `xorm:"logo VARCHAR(255)"`
	DisableGfx     bool         `xorm:"disable_gfx TINYINT(4)"`
	Deleted        time.Time    `xorm:"deleted DATETIME deleted"`
	CreateTime     time.Time    `xorm:"create_time DATETIME created"`
	UpdateTime     time.Time    `xorm:"update_time DATETIME updated"`
}

const RemoteAppTableName = "visual_remote_app"

func (r *RemoteApp) TableName() string {
	return RemoteAppTableName
}
