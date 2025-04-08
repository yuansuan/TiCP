package model

import (
	"time"

	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
)

type ApproveRecord struct {
	Id            snowflake.ID `json:"id" xorm:"pk 'id' bigint"`
	Type          int8         `json:"type" xorm:"'type' tinyint"`
	ApproveInfo   string       `json:"approve_info" xorm:"'approve_info' varchar"`
	Status        int8         `json:"status" xorm:"'status' tinyint"`
	ApplyUserId   snowflake.ID `json:"apply_user_id" xorm:"'apply_user_id' bigint"`
	ApplyUserName string       `json:"apply_user_name" xorm:"'apply_user_name' varchar"`
	Sign          string       `json:"sign" xorm:"'sign' varchar"`
	Content       string       `json:"content" xorm:"'content' text"`
	CreateTime    time.Time    `json:"create_time" xorm:"'create_time' datetime"`
	UpdateTime    time.Time    `json:"update_time" xorm:"'update_time' datetime"`
	ApproveTime   time.Time    `json:"approve_time" xorm:"'approve_time' datetime"`
}

type ApproveUser struct {
	Id              snowflake.ID `json:"id" xorm:"pk 'id' bigint"`
	ApproveRecordId snowflake.ID `json:"approve_record_id " xorm:"'approve_record_id' bigint"`
	ApproveUserId   snowflake.ID `json:"approve_user_id " xorm:"'approve_user_id' bigint"`
	ApproveUserName string       `json:"approve_user_name " xorm:"'approve_user_name' varchar"`
	Result          int8         `json:"result" xorm:"'result' tinyint"`
	Suggest         string       `json:"suggest" xorm:"'suggest' text"`
	CreateTime      time.Time    `json:"create_time" xorm:"'create_time' datetime"`
	UpdateTime      time.Time    `json:"update_time" xorm:"'update_time' datetime"`
	ApproveTime     time.Time    `json:"approve_time" xorm:"'approve_time' datetime"`
}

func (a *ApproveRecord) TableName() string {
	return "approve_record"
}

const (
	ApproveRecordName = "approve_record"
	ApproveUserName   = "approve_user"
)

type ApproveUserWithRecord struct {
	ApproveUser   `xorm:"extends"`
	ApproveRecord `xorm:"extends"`
}
