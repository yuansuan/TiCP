package model

import (
	"time"

	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
)

type Hardware struct {
	ID             snowflake.ID `xorm:"id BIGINT(20) pk"`
	OutHardwareID  string       `xorm:"out_hardware_id VARCHAR(255)"`
	Name           string       `xorm:"name VARCHAR(64)"`
	Desc           string       `xorm:"desc VARCHAR(255)"`
	Network        int          `xorm:"network INT(11)"`
	CPU            int          `xorm:"cpu INT(11)"`
	Mem            int          `xorm:"mem INT(11)"`
	Gpu            int          `xorm:"gpu INT(11)"`
	CPUModel       string       `xorm:"cpu_model VARCHAR(64)"`
	GpuModel       string       `xorm:"gpu_model VARCHAR(64)"`
	InstanceType   string       `xorm:"instance_type VARCHAR(64)"`
	InstanceFamily string       `xorm:"instance_family VARCHAR(64)"`
	Zone           string       `xorm:"zone VARCHAR(64)"`
	Deleted        time.Time    `xorm:"deleted DATETIME deleted"`
	CreateTime     time.Time    `xorm:"create_time DATETIME created"`
	UpdateTime     time.Time    `xorm:"update_time DATETIME updated"`
}

const HardwareTableName = "visual_hardware"

func (h *Hardware) TableName() string {
	return HardwareTableName
}
