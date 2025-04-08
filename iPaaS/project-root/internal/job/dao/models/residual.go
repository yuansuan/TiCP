package models

import (
	"time"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
)

// Residual 残差图
type Residual struct {
	ID                snowflake.ID `json:"id" xorm:"pk autoincr id comment('残差图ID') BIGINT(20)" bson:"id"`
	JobID             snowflake.ID `json:"job_id" xorm:"job_id not null index comment('作业ID') BIGINT(20)" bson:"job_id"`
	Content           string       `json:"content" xorm:"comment('残差图内容,经过base64') LONGTEXT" bson:"content"`
	Finished          bool         `json:"finished" xorm:"finished not null default 0 index comment('是否完成') TINYINT(1)" bson:"finished"`
	ResidualLogRegexp string       `json:"residual_log_regexp" xorm:"residual_log_regexp not null default 'stdout.log' comment('残差图文件') VARCHAR(255)" bson:"residual_log_regexp"`
	ResidualLogParser string       `json:"residual_log_parser" xorm:"residual_log_parser not null default '' comment('残差图解析器类型') VARCHAR(255)" bson:"residual_log_parser"`
	FailedReason      string       `json:"failed_reason" xorm:"failed_reason comment('失败原因') VARCHAR(512)" bson:"failed_reason"`
	CreateTime        time.Time    `json:"create_time" xorm:"created" bson:"create_time"`
	UpdateTime        time.Time    `json:"update_time" xorm:"updated" bson:"update_time"`
}

func (r Residual) TableName() string {
	return "residual"
}
