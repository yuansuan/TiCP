package models

import (
	"encoding/json"
	"time"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/common/go-kit/logging"
)

type AccountCashVoucherRelation struct {
	Id                snowflake.ID `json:"id" xorm:"not null pk comment('主键') BIGINT(20)"`
	AccountId         snowflake.ID `json:"account_id" xorm:"not null comment('账户ID') BIGINT(20)"`
	CashVoucherId     snowflake.ID `json:"cash_voucher_id" xorm:"not null comment('代金券ID') BIGINT(20)"`
	CashVoucherAmount int64        `json:"cash_voucher_amount" xorm:"not null default 0 comment('代金券原始总金额') BIGINT(20)"`
	UsedAmount        int64        `json:"used_amount" xorm:"not null default 0 comment('已使用金额') BIGINT(20)"`
	RemainingAmount   int64        `json:"remaining_amount" xorm:"not null default 0 comment('剩余金额') BIGINT(20)"`
	Status            int          `json:"status" xorm:"not null default 1 comment('账户代金券状态: 1:正常，2:禁用') TINYINT(20)"`
	ExpiredTime       time.Time    `json:"expired_time" xorm:"not null comment('过期时间') DATETIME"`
	IsExpired         int          `json:"is_expired" xorm:"not null default 0 comment('是否过期  0:正常 1:过期') TINYINT(1)"`
	IsDeleted         int          `json:"is_deleted" xorm:"not null default 0 comment('删除标记 0:正常  1:删除') TINYINT(1)"`
	OptUserId         snowflake.ID `json:"opt_user_id" xorm:"not null default 0 comment('操作用户id') BIGINT(20)"`
	CreateTime        time.Time    `json:"create_time" xorm:"not null default CURRENT_TIMESTAMP comment('创建时间') DATETIME"`
	UpdateTime        time.Time    `json:"update_time" xorm:"updated not null default CURRENT_TIMESTAMP comment('更新时间') DATETIME"`
}

func (a *AccountCashVoucherRelation) TableName() string {
	return "account_cash_voucher_relation"
}

func (a *AccountCashVoucherRelation) String() string {
	bytes, err := json.Marshal(a)
	if err != nil {
		logging.Default().Errorf("json marshal err: %v", err)
		return ""
	}

	return string(bytes)
}
