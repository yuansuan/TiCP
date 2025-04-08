package models

import (
	"time"

	schema "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	zone "github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/config"
)

type Session struct {
	Id             snowflake.ID         `xorm:"pk 'id'"`
	Zone           zone.Zone            `xorm:"'zone' comment('可用区')"`
	UserId         snowflake.ID         `xorm:"'user_id' comment('创建会话用户Id')"`
	InstanceId     snowflake.ID         `xorm:"'instance_id' comment('实例Id')"`
	Status         schema.SessionStatus `xorm:"'status' comment('会话状态')"`
	DesktopUrl     string               `xorm:"'desktop_url' comment('桌面接入地址')"`
	StartTime      *time.Time           `xorm:"'start_time' comment('会话开始时间')"`
	EndTime        *time.Time           `xorm:"'end_time' comment('会话结束时间')"`
	CloseSignal    bool                 `xorm:"'close_signal' comment('用户关闭会话信号')"`
	UserCloseTime  time.Time            `xorm:"'user_close_time' comment('用户关闭会话时间')"`
	ExitReason     string               `xorm:"'exit_reason' comment('会话退出原因')"`
	Deleted        bool                 `xorm:"'deleted' comment('是否已删除')"`
	AccountId      snowflake.ID         `xorm:"'account_id' comment('计费账户Id')"`
	PayByAccountId snowflake.ID         `xorm:"'pay_by_account_id' comment('代计费账户Id')"`
	ChargeType     schema.ChargeType    `xorm:"'charge_type' comment(计费模式 PrePaid | PostPaid | '')"`
	IsPaidFinished bool                 `xorm:"'is_paid_finished' comment('是否完成计费')"`
	RoomId         snowflake.ID         `xorm:"'room_id' comment('webrtc双端通信唯一标志')"`
	CreateTime     time.Time            `xorm:"'create_time' comment('创建时间') created"`
	UpdateTime     time.Time            `xorm:"'update_time' comment('更新时间') updated"`
}

func (*Session) TableName() string {
	return "cloudapp_session"
}

func (s *Session) ToHTTPModel() *schema.Session {
	return &schema.Session{}
}

func SessionStatusExist(s schema.SessionStatus) bool {
	_, exist := schema.SessionStatusMap[s]
	return exist
}
