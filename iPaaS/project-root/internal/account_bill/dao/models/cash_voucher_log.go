package models

import (
	"time"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
)

type AccountCashVoucherLog struct {
	Id                   snowflake.ID `json:"id" xorm:"not null pk comment('主键') BIGINT(20)"`
	AccountId            snowflake.ID `json:"account_id" xorm:"not null comment('账户ID') BIGINT(20)"`
	CashVoucherId        snowflake.ID `json:"cash_voucher_id" xorm:"not null comment('代金券id') BIGINT(20)"`
	AccountCashVoucherId snowflake.ID `json:"account_cash_voucher_id" xorm:"not null comment('用户代金券id') BIGINT(20)"`
	SignType             int          `json:"sign_type" xorm:"not null default 0 comment('使用标记 1:消费 2:过期') TINYINT(1)"`
	SourceInfo           string       `json:"source_info" xorm:"not null comment('账户代金券修改前信息') TEXT"`
	TargetInfo           string       `json:"target_info" xorm:"not null comment('账号代金券修改后信息') TEXT"`
	Comment              string       `json:"comment" xorm:"not null default '' comment('使用备注') VARCHAR(255)"`
	AccountBillId        snowflake.ID `json:"account_bill_id " xorm:"null default 0 comment('账单id') BIGINT(20)"`
	OptUserId            snowflake.ID `json:"opt_user_id" xorm:"not null default 0 comment('操作用户id') BIGINT(20)"`
	CreateTime           time.Time    `json:"create_time" xorm:"not null default CURRENT_TIMESTAMP comment('创建时间') DATETIME"`
	UpdateTime           time.Time    `json:"update_time" xorm:"not null default CURRENT_TIMESTAMP comment('更新时间') DATETIME"`
}

func (a *AccountCashVoucherLog) TableName() string {
	return "account_cash_voucher_log"
}
