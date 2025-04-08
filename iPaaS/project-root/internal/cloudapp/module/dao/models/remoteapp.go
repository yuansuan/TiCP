package models

import (
	"time"

	v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/module/utils"
)

type RemoteApp struct {
	Id         snowflake.ID `xorm:"'id' pk"`
	SoftwareId snowflake.ID `xorm:"'software_id' comment('关联软件id')"`
	Desc       string       `xorm:"'desc' comment('描述')"`
	Name       string       `xorm:"'name' comment('RemoteApp名称')"`
	Dir        string       `xorm:"'dir' comment('RemoteApp路径')"`
	Args       string       `xorm:"'args' comment('RemoteApp参数')"`
	Logo       string       `xorm:"'logo' comment('RemoteApp logo 前端显示')"`
	DisableGfx bool         `xorm:"'disable_gfx' comment('是否禁用gfx true 不使用视频 false 使用视频')"`
	LoginUser  string       `xorm:"'login_user' comment('登陆的系统用户名')"`
	CreateTime time.Time    `xorm:"'create_time' comment('创建时间') created"`
	UpdateTime time.Time    `xorm:"'update_time' comment('更新时间') updated"`
}

func (*RemoteApp) TableName() string {
	return "cloudapp_remote_app"
}

func (r *RemoteApp) ToHTTPModel() *v20230530.RemoteApp {
	return &v20230530.RemoteApp{
		Id:         r.Id.String(),
		SoftwareId: utils.PString(r.SoftwareId.String()),
		Desc:       utils.PString(r.Desc),
		Name:       utils.PString(r.Name),
		Dir:        utils.PString(r.Dir),
		Args:       utils.PString(r.Args),
		Logo:       utils.PString(r.Logo),
		DisableGfx: utils.PBool(r.DisableGfx),
		LoginUser:  utils.PString(r.LoginUser),
	}
}

type RemoteAppUserPass struct {
	Id            int64        `xorm:"'id' pk"`
	SessionId     snowflake.ID `xorm:"'session_id' comment('会话Id')"` // SessionId + RemoteAppName unique index
	RemoteAppName string       `xorm:"'remote_app_name' comment('远程应用名称')"`
	Username      string       `xorm:"'username' comment('系统用户')"`
	Password      string       `xorm:"'password' comment('系统用户密码')"`
}

func (*RemoteAppUserPass) TableName() string {
	return "cloudapp_remote_app_user_pass"
}
