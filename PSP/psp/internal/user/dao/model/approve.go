package model

import (
	"time"

	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
)

// Approve approve object
type Approve struct {
	Id             snowflake.ID `json:"id" xorm:"pk BIGINT(20)"`
	ApproverName   string       `json:"approver_name" xorm:"VARCHAR(64)"`
	SubmitterName  string       `json:"submitter_name" xorm:"VARCHAR(64)"`
	Content        string       `json:"content" xorm:"VARCHAR(128)"`
	Status         int8         `json:"status" xorm:"default '-1' TINYINT"`
	OperateType    string       `json:"operate_type" xorm:"VARCHAR(64)"`
	OperateContent string       `json:"operate_content" xorm:"TEXT"`
	OperateTarget  string       `json:"operate_target" xorm:"VARCHAR(64)"`
	TargetType     int8         `json:"target_type" xorm:"default '0' TINYINT"`
	ApiData        string       `json:"api_data" xorm:"TEXT"`
	SubmitTime     time.Time    `json:"submit_time" xorm:"not null default 'CURRENT_TIMESTAMP' DATETIME"`
	ApprovedTime   time.Time    `json:"approved_time" xorm:"not null default 'CURRENT_TIMESTAMP' DATETIME"`
}
