package model

import (
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"time"
)

type StorageQuota struct {
	UserId       snowflake.ID `json:"user_id" xorm:"not null comment('用户 ID') BIGINT(20)"`
	StorageUsage float64      `json:"storage_usage" xorm:"not null default 0 comment('存储空间用量') FLOAT(10,2) "`
	StorageLimit float64      `json:"storage_limit" xorm:"not null default 0 comment('存储上限') FLOAT(10,2)"`
	CreateTime   time.Time    `json:"create_time" xorm:"not null default CURRENT_TIMESTAMP comment('创建时间') DATETIME"`
	UpdateTime   time.Time    `json:"update_time" xorm:"not null default CURRENT_TIMESTAMP comment('更新时间') DATETIME"`
}
