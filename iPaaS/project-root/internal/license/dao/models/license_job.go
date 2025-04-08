package models

import (
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"time"
)

type LicenseJob struct {
	Id         snowflake.ID `json:"id" xorm:"pk BIGINT(20)"`
	ModuleId   snowflake.ID `json:"module_id" xorm:"not null comment('module id') BIGINT(20)"`
	JobId      snowflake.ID `json:"job_id" xorm:"not null comment('作业ID') index BIGINT(20)"`
	Licenses   int64        `json:"licenses" xorm:"not null comment('licenses') BIGINT(20)"`
	Used       int          `json:"used" xorm:"not null default 1 comment('1-使用中 2-使用完成') TINYINT(1)"`
	CreateTime time.Time    `json:"create_time" xorm:"created"`
	UpdateTime time.Time    `json:"update_time" xorm:"updated"`
	LicenseId  snowflake.ID `json:"license_id" xorm:"not null comment('license id') BIGINT(20)"`
}

func (*LicenseJob) TableName() string {
	return "license_job"
}

type LicenseUsedResult struct {
	Id snowflake.ID `json:"id"`
	// 使用数
	Licenses int64 `json:"licenses"`
	// 作业id
	JobId snowflake.ID `json:"job_id"`
	// 作业名称
	JobName string `json:"job_name"`
	// 开始时间
	CreateTime time.Time `json:"create_time"`
	// 应用id
	AppId snowflake.ID `json:"app_id"`
}
