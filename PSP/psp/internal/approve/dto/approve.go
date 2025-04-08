package dto

import (
	"github.com/yuansuan/ticp/PSP/psp/pkg/xtype"
)

type ApplyApproveRequest struct {
	ApproveType            ApproveType            `json:"approve_type"`              // 审批类型
	ApproveUserID          string                 `json:"approve_user_id"`           // 审批人id，必须是具有安全审批权限的用户
	ApproveUserName        string                 `json:"approve_user_name"`         // 审批人名称
	UserApproveInfoRequest UserApproveInfoRequest `json:"user_approve_info_request"` // 审批详情，用户相关审批专用
	RoleApproveInfoRequest RoleApproveInfoRequest `json:"role_approve_Info_request"` // 审批详情，角色相关审批专用
}

type CancelApproveRequest struct {
	Id string // apply_record 的id
}
type CancelApproveResponse struct {
	State bool
}

type GetApproveListRequest struct {
	Page      *xtype.Page `json:"page"`
	Type      int8        `json:"type"`
	Result    int8        `json:"result"`
	Status    int8        `json:"status"`
	StartTime int64       `json:"start_time"`
	EndTime   int64       `json:"end_time"`
}

type GetApprovePendingRequest struct {
	Page            *xtype.Page `json:"page"`
	Type            int8        `json:"type"`
	ApplicationName string      `json:"application_name"`
	StartTime       int64       `json:"start_time"`
	EndTime         int64       `json:"end_time"`
}

type GetApproveCompleteRequest struct {
	Page            *xtype.Page `json:"page"`
	Type            int8        `json:"type"`
	ApplicationName string      `json:"application_name"`
	StartTime       int64       `json:"start_time"`
	EndTime         int64       `json:"end_time"`
	Status          int8        `json:"status"`
}

type HandleApproveRequest struct {
	ID          string      `json:"id"`           // approve_user表id
	RecordID    string      `json:"record_id"`    // approve_user表id
	ApproveType ApproveType `json:"approve_type"` // 审批类型
	Suggest     string      `json:"suggest"`      // 意见
}

type RoleApproveInfoRequest struct {
	Id      int64   `json:"id"`      // 角色id
	Name    string  `json:"name"`    // 角色名
	Comment string  `json:"comment"` // 描述
	Perms   []int64 `json:"perms"`   // 权限id
}

type UserApproveInfoRequest struct {
	Id            string  `json:"id"`             // 用户id
	Name          string  `json:"name"`           // 用户名
	Password      string  `json:"password"`       // 密码
	Email         string  `json:"email"`          // 邮箱
	Mobile        string  `json:"mobile"`         // 手机号
	RealName      string  `json:"real_name"`      // 真实姓名
	Roles         []int64 `json:"roles"`          // 所属角色信息
	EnableOpenapi bool    `json:"enable_openapi"` // 是否启动openapi
}

type ApproveInfo struct {
	Id                string `json:"id"`
	RecordID          string `json:"record_id"`
	ApplicationName   string `json:"application_name"`
	ApproveCreateTime string `json:"create_time"`
	Type              int8   `json:"type"`
	Content           string `json:"content"`
	ApproveUserName   string `json:"approve_user_name"`
	ApproveTime       string `json:"approve_time"`
	Status            int8   `json:"status"`
	Suggest           string `json:"suggest"`
}

type ApproveListCondition struct {
	ApplyId    int64
	RecordType int8
	StartTime  int64
	EndTime    int64
	Status     int8
}

type ApplicationListCondition struct {
	UserId     int64
	ApplyName  string
	RecordType int8
	StartTime  int64
	EndTime    int64
	Status     []int8
}

type ApproveStatue int8

const (
	ApproveStatusWaiting = ApproveStatue(1)
	ApproveStatusPass    = ApproveStatue(2)
	ApproveStatusRefuse  = ApproveStatue(3)
	ApproveStatusCancel  = ApproveStatue(4)
	ApproveStatusFailed  = ApproveStatue(5)
)

type ApproveType int8

const (
	ApproveTypeAddUser     = ApproveType(1) // 新增用户
	ApproveTypeDelUser     = ApproveType(2) // 删除用户
	ApproveTypeEditUser    = ApproveType(3) // 编辑用户
	ApproveTypeEnableUser  = ApproveType(4) // 启用用户
	ApproveTypeDisableUser = ApproveType(5) // 禁用用户

	ApproveTypeAddRole        = ApproveType(6) // 新增角色
	ApproveTypeDelRole        = ApproveType(7) // 删除角色
	ApproveTypeEditRole       = ApproveType(8) // 编辑角色
	ApproveTypeSetLdapDefRole = ApproveType(9) // 设置LDAP默认角色
)

type ApproveResult int8

const (
	ApproveResultDefault = ApproveStatue(0) // 默认状态
	ApproveResultPass    = ApproveStatue(1) // 通过
	ApproveResultRefuse  = ApproveStatue(2) // 拒绝
)
