package model

import (
	"time"

	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
)

// JobTimeline 作业时间线
type JobTimeline struct {
	JobId     snowflake.ID
	EventName string
	EventTime time.Time
}

func (j *JobTimeline) TableName() string {
	return "job_timeline"
}
