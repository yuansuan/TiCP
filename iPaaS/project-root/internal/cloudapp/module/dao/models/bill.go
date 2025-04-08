package models

import (
	"time"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
)

// Bill
// 后付费：按顺序发送后付费订单详情，发送一次更新一下BillTime
// 预付费：创建会话的时候按资源收费，记录一次即可
type Bill struct {
	Id           int64        `xorm:"pk 'id'"`
	SessionId    snowflake.ID `xorm:"'session_id' comment('会话Id')"` // SessionId + ResourceId 唯一性索引
	OrderId      snowflake.ID `xorm:"'order_id'"`                   // 唯一性索引
	ResourceId   snowflake.ID `xorm:"'resource_id'"`
	ResourceType string       `xorm:"'resource_type'"`
	BillTime     time.Time    `xorm:"'bill_time'"`
}

func (*Bill) TableName() string {
	return "cloudapp_bill"
}
