package models

import (
	"time"

	v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	zone "github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/config"
)

type Platform string

const (
	Linux   Platform = "LINUX"
	Windows Platform = "WINDOWS"
)

func (p Platform) String() string {
	return string(p)
}

type Software struct {
	Id         snowflake.ID `xorm:"pk 'id'"`
	Zone       zone.Zone    `xorm:"'zone' comment('可用区')"`
	Name       string       `xorm:"'name' comment('软件方案名字')"`
	Desc       string       `xorm:"'desc' comment('软件方案描述')"`
	Icon       string       `xorm:"'icon' comment('软件方案图标')"`
	Platform   Platform     `xorm:"'platform' comment('软件平台：LINUX， WINDOWS')"`
	ImageId    string       `xorm:"'image_id' comment('腾讯云镜像Id')"`
	InitScript string       `xorm:"'init_script' comment('初始化脚本内容')"`
	GpuDesired *bool        `xorm:"'gpu_desired' comment('是否需要GPU支持')"`
	CreateTime time.Time    `xorm:"'create_time' comment('创建时间') created"`
	UpdateTime time.Time    `xorm:"'update_time' comment('更新时间') updated"`
}

func (*Software) TableName() string {
	return "cloudapp_software"
}

func (s *Software) ToHTTPModel() *v20230530.Software {
	return &v20230530.Software{
		SoftwareId: s.Id.String(),
		Zone:       s.Zone.String(),
		Name:       s.Name,
		Desc:       s.Desc,
		Icon:       s.Icon,
		Platform:   s.Platform.String(),
		ImageId:    s.ImageId,
		InitScript: s.InitScript,
		GpuDesired: s.GpuDesired,
	}
}

type SoftwareWithUser struct {
	Software     `xorm:"extends"`
	SoftwareUser `xorm:"extends"`
}

type SoftwareUser struct {
	Id         int64        `xorm:"pk 'id'"`
	SoftwareId snowflake.ID `xorm:"'software_id'"` // SoftwareId + UserId 唯一性索引
	UserId     snowflake.ID `xorm:"'user_id'"`
}

func (*SoftwareUser) TableName() string {
	return "cloudapp_software_user"
}
