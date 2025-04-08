package models

import (
	"time"

	v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	zone "github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/config"
)

type Hardware struct {
	Id             snowflake.ID `xorm:"pk 'id'"`
	Zone           zone.Zone    `xorm:"'zone' comment('可用区')"`
	Name           string       `xorm:"'name' comment('硬件方案名字')"`
	Desc           string       `xorm:"'desc' comment('硬件方案描述')"`
	InstanceType   string       `xorm:"'instance_type' comment('实例类型名字')"`
	InstanceFamily string       `xorm:"'instance_family' comment('实例机型系列')"`
	Network        int64        `xorm:"'network' comment('实例的最大内网带宽')"`
	Cpu            int64        `xorm:"'cpu' comment('CPU核数')"`
	CpuModel       string       `xorm:"'cpu_model' comment('CPU型号')"`
	Mem            int64        `xorm:"'mem' comment('内存容量，单位G')"`
	Gpu            int64        `xorm:"'gpu' comment('GPU数量')"`
	GpuModel       string       `xorm:"'gpu_model' comment('GPU型号')"`
	CreateTime     time.Time    `xorm:"'create_time' comment('创建时间') created"`
	UpdateTime     time.Time    `xorm:"'update_time' comment('更新时间') updated"`
}

func (*Hardware) TableName() string {
	return "cloudapp_hardware"
}

func (h *Hardware) ToHTTPModel() *v20230530.Hardware {
	return &v20230530.Hardware{
		HardwareId:     h.Id.String(),
		Zone:           h.Zone.String(),
		Name:           h.Name,
		Desc:           h.Desc,
		InstanceType:   h.InstanceType,
		InstanceFamily: h.InstanceFamily,
		Network:        int(h.Network),
		Cpu:            int(h.Cpu),
		CpuModel:       h.CpuModel,
		Mem:            int(h.Mem),
		Gpu:            int(h.Gpu),
		GpuModel:       h.GpuModel,
	}
}

type HardwareWithUser struct {
	Hardware     `xorm:"extends"`
	HardwareUser `xorm:"extends"`
}

type HardwareUser struct {
	Id         int64        `xorm:"pk 'id'"`
	HardwareId snowflake.ID `xorm:"'hardware_id'"` // HardwareId + UserId 唯一性索引
	UserId     snowflake.ID `xorm:"'user_id'"`
}

func (*HardwareUser) TableName() string {
	return "cloudapp_hardware_user"
}
