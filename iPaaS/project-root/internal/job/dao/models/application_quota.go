package models

import (
	"time"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
)

// ApplicationQuota 应用配额
type ApplicationQuota struct {
	ID snowflake.ID `xorm:"id not null pk BIGINT(20)"`
	// application_id 和 ys_id 联合索引
	ApplicationID snowflake.ID `xorm:"application_id not null comment('应用id') unique(idx_application_quota_unique)"`
	YsID          snowflake.ID `xorm:"ys_id not null comment('用户id') unique(idx_application_quota_unique)"`
	CreateTime    time.Time    `json:"create_time" xorm:"created"`
	UpdateTime    time.Time    `json:"update_time" xorm:"updated"`
}
