package models

import (
	"time"

	v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	zone "github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/config"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/module/utils"
)

type Instance struct {
	Id             snowflake.ID   `xorm:"pk 'id'"`
	Zone           zone.Zone      `xorm:"'zone' comment('可用区')"`
	HardwareId     snowflake.ID   `xorm:"'hardware_id' comment('绑定的硬件Id')"`
	SoftwareId     snowflake.ID   `xorm:"'software_id' comment('绑定的软件Id')"`
	InitScript     string         `xorm:"'init_script' comment('实例初始化脚本模板')"`
	UserParams     string         `xorm:"'user_params' comment('实例初始化用户传入参数')"`
	UserScript     string         `xorm:"'user_script' comment('实例实际执行的脚本内容')"`
	InstanceId     string         `xorm:"'instance_id' comment('实例Id')"`
	InstanceData   string         `xorm:"'instance_data' comment('实例数据')"`
	SshPassword    string         `xorm:"'ssh_password' comment('登录密码')"`
	InstanceStatus InstanceStatus `xorm:"'instance_status' comment('腾讯云实例状态')"`
	BootVolumeId   string         `xorm:"'boot_volume_id' comment('启动盘Id')"`
	StartTime      *time.Time     `xorm:"'start_time' comment('实例开始时间')"`
	EndTime        *time.Time     `xorm:"'end_time' comment('实例结束时间')"`
	CreateTime     time.Time      `xorm:"'create_time' comment('创建时间') created"`
	UpdateTime     time.Time      `xorm:"'update_time' comment('更新时间') updated"`
}

func (*Instance) TableName() string {
	return "cloudapp_instance"
}

type InstanceStatus string

const (
	InstanceInvalid      InstanceStatus = "INVALID"
	InstancePending      InstanceStatus = "PENDING"
	InstanceLaunchFailed InstanceStatus = "LAUNCH FAILED"
	InstanceCreated      InstanceStatus = "CREATED"
	InstanceRunning      InstanceStatus = "RUNNING"
	InstanceStarting     InstanceStatus = "STARTING"
	InstanceStopping     InstanceStatus = "STOPPING"
	InstanceStopped      InstanceStatus = "STOPPED"
	InstanceRebooting    InstanceStatus = "REBOOTING"
	InstanceShutdown     InstanceStatus = "SHUTDOWN"
	InstanceTerminating  InstanceStatus = "TERMINATING"
	InstanceTerminated   InstanceStatus = "TERMINATED"
)

func (i InstanceStatus) String() string {
	return string(i)
}

type SessionWithDetail struct {
	Session  `xorm:"extends"`
	Instance `xorm:"extends"`
	Hardware `xorm:"extends"`
	Software `xorm:"extends"`
}

func (sd *SessionWithDetail) ToDetailHTTPModel() *v20230530.Session {
	return &v20230530.Session{
		Id:              sd.Session.Id.String(),
		Zone:            sd.Session.Zone.String(),
		Status:          sd.Session.Status.String(),
		StreamUrl:       sd.Session.DesktopUrl,
		CreateTime:      utils.PTime(sd.Session.CreateTime),
		StartTime:       sd.Session.StartTime,
		EndTime:         sd.Session.EndTime,
		MachinePassword: sd.Instance.SshPassword,
		ExitReason:      sd.Session.ExitReason,
		Software:        sd.Software.ToHTTPModel(),
		Hardware:        sd.Hardware.ToHTTPModel(),
		RemoteApps:      make([]*v20230530.RemoteApp, 0),
	}
}

func (sd *SessionWithDetail) ToAdminDetailHTTPModel() *v20230530.Session {
	sess := sd.ToDetailHTTPModel()
	sess.UserId = sd.UserId.String()

	return sess
}

type SessionWithInstance struct {
	Session  `xorm:"extends"`
	Instance `xorm:"extends"`
}
