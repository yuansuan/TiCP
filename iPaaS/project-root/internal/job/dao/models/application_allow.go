package models

import (
	"time"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
)

// ApplicationAllow 应用配额
type ApplicationAllow struct {
	ID            snowflake.ID `xorm:"id not null pk BIGINT(20)"`
	ApplicationID snowflake.ID `xorm:"application_id not null comment('应用id') unique(idx_application_allow_unique)"`
	CreateTime    time.Time    `json:"create_time" xorm:"created"`
	UpdateTime    time.Time    `json:"update_time" xorm:"updated"`
}
