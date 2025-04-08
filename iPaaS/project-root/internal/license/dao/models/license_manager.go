package models

import (
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/common/project-root-api/license/v1/license_manager/os"
	"github.com/yuansuan/ticp/common/project-root-api/license/v1/license_manager/publish"
	"time"
)

type LicenseManager struct {
	Id          snowflake.ID   `json:"id" xorm:"pk BIGINT(20)"`
	AppType     string         `json:"app_type" xorm:"not null default '' comment('求解器软件类型') VARCHAR(255)"`
	Os          os.OS          `json:"os" xorm:"not null default 1 comment('操作系统 1-linux 2-win') TINYINT(1)"`
	Status      publish.Status `json:"status" xorm:"not null default 2 comment('发布状态 1-已发布 2-未发布') TINYINT(1)"`
	Description string         `json:"description" xorm:"not null default '' comment('描述') VARCHAR(1024)"`
	ComputeRule string         `json:"compute_rule" xorm:"not null default '' comment('license使用计算规则') VARCHAR(1024)"`
	PublishTime time.Time      `json:"publish_time" xorm:"comment('发布时间') DATETIME"`
	CreateTime  time.Time      `json:"create_time" xorm:"created"`
	UpdateTime  time.Time      `json:"update_time" xorm:"updated"`
}

func (*LicenseManager) TableName() string {
	return "license_manager"
}
