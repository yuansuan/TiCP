package dto

import (
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
	"github.com/yuansuan/ticp/PSP/psp/pkg/xtype"
)

// ProjectAddRequest 新增项目
type ProjectAddRequest struct {
	ProjectName  string   `json:"project_name"`  // 项目名称
	ProjectOwner string   `json:"project_owner"` // 项目管理员id
	StartTime    int64    `json:"start_time"`    // 项目开始时间
	EndTime      int64    `json:"end_time"`      // 项目结束时间
	Members      []string `json:"members"`       // 项目成员userIds
	Comment      string   `json:"comment"`       // 项目描述
}

// ProjectAddResponse 新增项目返回
type ProjectAddResponse struct {
	ID string `json:"id"` // 新增 projectId
}

// ProjectListRequest 项目请求列表
type ProjectListRequest struct {
	Page        *xtype.Page      `json:"page"`         // 分页信息
	OrderSort   *xtype.OrderSort `json:"order_sort"`   // 排序条件
	ProjectName string           `json:"project_name"` // 项目名称
	State       []string         `json:"state"`        // 项目状态, Init,Running,Terminated,Completed
	StartTime   int64            `json:"start_time"`   // 项目开始时间，以秒为单位
	EndTime     int64            `json:"end_time"`     // 项目结束时间，以秒为单位
	IsSysMenu   bool             `json:"is_sys_menu"`  // 请求是否来自系统管理菜单
}

// ProjectListInfo 项目列表信息
type ProjectListInfo struct {
	ID               string                `json:"id"`                 // 项目id
	ProjectName      string                `json:"project_name"`       // 项目名称
	State            string                `json:"state"`              // 项目状态
	ProjectOwnerName string                `json:"project_owner_name"` // 项目管理员名字
	ProjectOwnerID   string                `json:"project_owner_id"`   // 项目管理员id
	StartTime        string                `json:"start_time"`         // 项目开始时间
	EndTime          string                `json:"end_time"`           // 项目结束时间
	Comment          string                `json:"comment"`            // 项目描述
	IsProjectOwner   bool                  `json:"is_project_owner"`   // 是否当前项目管理员
	CreateTime       string                `json:"create_time"`        // 项目创建时间
	Members          []*ProjectMembersInfo `json:"members"`            // 项目成员
}

// ProjectListResponse 项目列表返回
type ProjectListResponse struct {
	ProjectList []*ProjectListInfo `json:"project_list"` // 项目列表
	Total       int64              `json:"total"`        // 总数
}

type CurrentProjectListRequest struct {
	State string `json:"state" form:"state" enums:"Init,Running,Terminated,Completed"` // 项目状态
}

type CurrentProjectInfo struct {
	Id   string `json:"id"`   // 项目 Id
	Name string `json:"name"` // 项目名称
}

type CurrentProjectListResponse struct {
	Projects []*CurrentProjectInfo `json:"projects"` // 项目列表
}

type CurrentProjectListForParamRequest struct {
	IsAdmin bool `json:"is_admin" form:"is_admin"` // 是否管理员
}

type CurrentProjectListForParamResponse struct {
	Projects []*CurrentProjectInfo `json:"projects"` // 项目列表
}

// ProjectDetailRequest 项目详情请求
type ProjectDetailRequest struct {
	ProjectID string `json:"project_id" form:"project_id"` // 项目id
}

// ProjectMembersInfo 项目成员信息
type ProjectMembersInfo struct {
	UserID   string `json:"user_id"`   // 用户id
	UserName string `json:"user_name"` // 用户名称
}

// ProjectDetailResponse 项目详情返回
type ProjectDetailResponse struct {
	*ProjectListInfo
}

// ProjectDeleteRequest 项目删除请求
type ProjectDeleteRequest struct {
	ProjectID string `json:"project_id"` // 项目id
}

// ProjectTerminatedRequest 项目终止请求
type ProjectTerminatedRequest struct {
	ProjectID string `json:"project_id"` // 项目id
}

// ProjectEditRequest 项目编辑
type ProjectEditRequest struct {
	ProjectID    string   `json:"project_id"`    // 项目id
	ProjectOwner string   `json:"project_owner"` // 项目管理员id
	StartTime    int64    `json:"start_time"`    // 项目开始时间
	EndTime      int64    `json:"end_time"`      // 项目结束时间
	Members      []string `json:"members"`       // 项目成员userIds
	Comment      string   `json:"comment"`       // 项目描述
}

// ProjectModifyOwnerRequest 批量修改管理员
type ProjectModifyOwnerRequest struct {
	ProjectIDs           []string `json:"project_ids"`             // 需要转移权限的项目id集合
	TargetProjectOwnerID string   `json:"target_project_owner_id"` // 目标项目管理员
}

type ProjectMemberPbResp struct {
	ProjectID   snowflake.ID `xorm:"id"`
	ProjectName string       `xorm:"project_name"`
	State       string       `xorm:"state"`
	LinkPath    string       `xorm:"link_path"`
}
