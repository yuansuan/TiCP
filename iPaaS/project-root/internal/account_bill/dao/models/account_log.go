package models

import (
	"time"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
)

type AccountLog struct {
	Id          snowflake.ID `json:"id" xorm:"pk default 0 comment('主键ID') BIGINT(20)"`
	AccountId   snowflake.ID `json:"account_id" xorm:"not null default 0 comment('账户ID') BIGINT(20)"`
	OperatorUid snowflake.ID `json:"operator_uid" xorm:"not null default 0 comment('操作人ID') BIGINT(20)"`
	Params      string       `json:"params" xorm:"comment('参数') TEXT"`
	Old         string       `json:"old" xorm:"comment('修改前数据') TEXT"`
	Updated     string       `json:"updated" xorm:"comment('修改后数据') TEXT"`
	CreateTime  time.Time    `json:"create_time" xorm:"not null default CURRENT_TIMESTAMP comment('创建时间') DATETIME"`
}

func (a *AccountLog) TableName() string {
	return "account_log"
}
