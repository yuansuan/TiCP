package models

import (
	"time"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
)

type CashVoucher struct {
	Id                 snowflake.ID `json:"id" xorm:"not null pk comment('主键') BIGINT(20)"`
	Name               string       `json:"name" xorm:"not null comment('代金券名称') VARCHAR(255)"`
	Amount             int64        `json:"amount" xorm:"not null default 0 comment('代金券金额') BIGINT(20)"`
	AvailabilityStatus int          `json:"availability_status" xorm:"not null default 0 comment('上下架状态 0:下架 1:上架') TINYINT(1)"`
	OptUserId          snowflake.ID `json:"opt_user_id" xorm:"not null comment('业务平台操作用户id') BIGINT(20)"`
	IsExpired          int          `json:"is_expired" xorm:"not null default 0 comment('是否过期 0:正常 1:过期') TINYINT(1)"`
	AbsExpiredTime     time.Time    `json:"abs_expired_time" xorm:"not null comment('绝对过期时间') DATETIME"`
	RelExpiredTime     int64        `json:"rel_expired_time" xorm:"not null default 0 comment('相对过期时间，以秒来计算') BIGINT(20)"`
	ExpiredType        int          `json:"expired_type" xorm:"not null default 0 comment('过期类型 1:绝对 2:相对') TINYINT(1)"`
	Comment            string       `json:"comment" xorm:"not null default '' comment('备注') VARCHAR(255)"`
	IsDeleted          int          `json:"is_deleted" xorm:"not null default 0 comment('删除标记') TINYINT(1)"`
	CreateTime         time.Time    `json:"create_time" xorm:"not null default CURRENT_TIMESTAMP comment('创建时间') DATETIME"`
	UpdateTime         time.Time    `json:"update_time" xorm:"not null default CURRENT_TIMESTAMP comment('更新时间') DATETIME"`
}

func (a *CashVoucher) TableName() string {
	return "cash_voucher"
}

type SelectCashVouchersReq struct {
	Id                 snowflake.ID
	Name               string
	IsExpired          int64
	AvailabilityStatus string
	StartTime          time.Time
	EndTime            time.Time
	Index              int64
	Size               int64
	OptUserId          snowflake.ID
}
