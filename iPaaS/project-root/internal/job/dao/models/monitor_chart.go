package models

import (
	"time"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
)

// MonitorChart 监控图表
type MonitorChart struct {
	ID                 snowflake.ID `json:"id" xorm:"pk autoincr id comment('监控图表ID') BIGINT(20)"`
	JobID              snowflake.ID `json:"job_id" xorm:"job_id not null index comment('作业ID') BIGINT(20)"`
	Content            string       `json:"content" xorm:"comment('监控图表内容') LONGTEXT"`
	Finished           bool         `json:"finished" xorm:"finished not null default 0 index comment('是否完成') TINYINT(1)"`
	MonitorChartRegexp string       `json:"monitor_chart_regexp" xorm:"monitor_chart_regexp not null default '.*\\.out' comment('监控图表文件规则') VARCHAR(255)"`
	MonitorChartParser string       `json:"monitor_chart_parser" xorm:"monitor_chart_parser not null default '' comment('监控图表解析器') VARCHAR(255)"`
	FailedReason       string       `json:"failed_reason" xorm:"failed_reason comment('失败原因') VARCHAR(512)"`
	CreateTime         time.Time    `json:"create_time" xorm:"created"`
	UpdateTime         time.Time    `json:"update_time" xorm:"updated"`
}
