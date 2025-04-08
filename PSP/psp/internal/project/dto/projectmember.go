package dto

import "github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"

type ProjectMemberRequest struct {
	ProjectId string   `json:"project_id"` // 项目id
	UserIds   []string `json:"user_ids"`   // 成员id
}

type ProjectMemberResponse struct {
	ID string `json:"project_id"` // 项目id
}

type ProjectMemberCount struct {
	ProjectId snowflake.ID `json:"project_id" xorm:"project_id"` // 项目id
	Count     int64        `json:"count" xorm:"count"`           // 成员数量
}
