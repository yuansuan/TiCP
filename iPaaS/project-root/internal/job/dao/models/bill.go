package models

import (
	"time"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
)

// Bill
// 预付费：作业开始时扣一次费，写表即可
// 后付费：运行中的作业按按需周期性的扣费，每次扣费更新BillTime
type Bill struct {
	Id             int          `xorm:"pk 'id'"`
	JobId          snowflake.ID `xorm:"'job_id'"`   // 唯一性索引
	OrderId        snowflake.ID `xorm:"'order_id'"` // 唯一性索引
	AppId          snowflake.ID `xorm:"'app_id'"`
	BilledDuration int64        `xorm:"'billed_duration'"` // 已经被收费的duration，单位s，对应job表中的ExecutionDuration
	BillTime       time.Time    `xorm:"'bill_time'"`       // 扣上一个BillDuration的时间点
}

func (*Bill) TableName() string {
	return "job_bill"
}
