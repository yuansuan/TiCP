package openapi

import "github.com/yuansuan/ticp/PSP/psp/pkg/xtype"

// ProjectListRequest 项目请求列表
type ProjectListRequest struct {
	Page        *xtype.Page `json:"page"`         // 分页信息
	ProjectName string      `json:"project_name"` // 项目名称
	State       []string    `json:"state"`        // 项目状态, Init,Running,Terminated,Completed
	StartTime   int64       `json:"start_time"`   // 项目开始时间，以秒为单位
	EndTime     int64       `json:"end_time"`     // 项目结束时间，以秒为单位
}

// ProjectListResponse 项目列表返回
type ProjectListResponse struct {
	ProjectList []*ProjectListInfo `json:"project_list"` // 项目列表
	Total       int64              `json:"total"`        // 总数
}

// ProjectListInfo 项目列表信息
type ProjectListInfo struct {
	ID          string `json:"id"`           // 项目id
	ProjectName string `json:"project_name"` // 项目名称
	State       string `json:"state"`        // 项目状态
	StartTime   string `json:"start_time"`   // 项目开始时间
	EndTime     string `json:"end_time"`     // 项目结束时间
	Comment     string `json:"comment"`      // 项目描述
}

// ProjectMembersInfo 项目成员信息
type ProjectMembersInfo struct {
	UserID   string `json:"user_id"`   // 用户id
	UserName string `json:"user_name"` // 用户名称
}
