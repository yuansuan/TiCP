package models

import (
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"time"
)

type ModuleConfig struct {
	Id          snowflake.ID `json:"id" db:"id" xorm:"pk BIGINT(20)"`
	LicenseId   snowflake.ID `json:"license_id" db:"license_id" xorm:"not null comment('license id') BIGINT(20) unique(license_module_name)"`
	ModuleName  string       `json:"module_name" db:"module_name" xorm:"not null comment('模块名称') VARCHAR(255) unique(license_module_name)"`
	Total       int          `json:"total" db:"total" xorm:"not null default 0 comment('licenses数量') INT(11)"`
	Used        int          `json:"used" db:"used" xorm:"not null default 0 comment('licenses已使用数量') INT(11)"`
	ActualTotal int          `json:"actual_total" db:"actual_total" xorm:"not null default 0 comment('实时总数量') INT(11)"`
	ActualUsed  int          `json:"actual_used" db:"actual_used" xorm:"not null default 0 comment('实时已使用数量') INT(11)"`
	CreateTime  time.Time    `json:"create_time" xorm:"created"`
	UpdateTime  time.Time    `json:"update_time" xorm:"updated"`
}
