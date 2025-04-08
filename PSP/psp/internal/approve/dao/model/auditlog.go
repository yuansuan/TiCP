package model

import (
	"time"

	"github.com/yuansuan/ticp/PSP/psp/internal/approve/dto"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
)

type AuditLog struct {
	Id              snowflake.ID        `json:"id" xorm:"pk 'id' bigint"`
	UserId          snowflake.ID        `json:"user_id" xorm:"'user_id' bigint"`
	UserName        string              `json:"user_name" xorm:"'user_name' varchar"`
	IpAddress       string              `json:"ip_address" xorm:"'ip_address' varchar"`
	OperateType     string              `json:"operate_type" xorm:"'operate_type' varchar"`
	OperateContent  string              `json:"operate_content" xorm:"'operate_content' text"`
	OperateUserType dto.OperateUserType `json:"operate_user_type" xorm:"'operate_user_type' tinyint"`
	OperateTime     time.Time           `json:"operate_time" xorm:"'operate_time' datetime"`
}
