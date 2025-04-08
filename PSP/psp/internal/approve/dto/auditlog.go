package dto

import (
	"time"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/approve"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
	"github.com/yuansuan/ticp/PSP/psp/pkg/xtype"
)

type SaveLogRequest struct {
	UserId         snowflake.ID            `json:"user_id" `
	UserName       string                  `json:"user_name" `
	OperateType    approve.OperateTypeEnum `json:"operate_type" `
	OperateContent string                  `json:"operate_content" `
	IpAddress      string                  `json:"ip_address"`
}

type AuditLogListRequest struct {
	Page        *xtype.Page             `json:"page"`
	UserName    string                  `json:"user_name"`
	OperateType approve.OperateTypeEnum `json:"operate_type" `
	IpAddress   string                  `json:"ip_address"`
	StartTime   string                  `json:"start_time"`
	EndTime     string                  `json:"end_time"`
}

type AuditLogListAllRequest struct {
	Page            *xtype.Page             `json:"page"`
	UserName        string                  `json:"user_name"`
	OperateType     approve.OperateTypeEnum `json:"operate_type" `
	IpAddress       string                  `json:"ip_address"`
	StartTime       string                  `json:"start_time"`
	EndTime         string                  `json:"end_time"`
	OperateUserType OperateUserType         `json:"operate_user_type"`
}

type AuditLogListResponse struct {
	Page         *xtype.PageResp `json:"page"`
	AuditLogInfo []*AuditLogInfo `json:"list"`
}

type AuditLogListAllResponse struct {
	Page         *xtype.PageResp       `json:"page"`
	AuditLogInfo []*AuditLogExportInfo `json:"list"`
}

type ApproveLogListAllResponse struct {
	Page    *xtype.PageResp `json:"page"`
	LogInfo []*ApproveInfo  `json:"list"`
}

type ThreeStateRequest struct {
}

type ThreeStateResponse struct {
	State bool `json:"state"`
}

type AuditLogInfo struct {
	Id             string    `json:"id"`
	UserName       string    `json:"user_name"`
	OperateType    string    `json:"operate_type" `
	OperateContent string    `json:"operate_content" `
	IpAddress      string    `json:"ip_address"`
	OperateTime    time.Time `json:"operate_time"`
}

type AuditLogExportInfo struct {
	Id             string    `json:"id"`
	UserName       string    `json:"user_name"`
	OperateType    string    `json:"operate_type" `
	OperateContent string    `json:"operate_content" `
	IpAddress      string    `json:"ip_address"`
	OperateTime    time.Time `json:"operate_time"`
}

func OperateTypeString(e approve.OperateTypeEnum) string {
	switch e {
	case approve.OperateTypeEnum_FILE_MANAGER:
		return "文件管理"
	case approve.OperateTypeEnum_USER_MANAGER:
		return "用户管理"
	case approve.OperateTypeEnum_RBAC_MANAGER:
		return "权限管理"
	case approve.OperateTypeEnum_JOB_MANAGER:
		return "作业管理"
	case approve.OperateTypeEnum_APP_MANAGER:
		return "计算应用"
	case approve.OperateTypeEnum_NODE_MANAGER:
		return "集群管理"
	case approve.OperateTypeEnum_LICENSE_MANAGER:
		return "许可证管理"
	case approve.OperateTypeEnum_PROJECT_MANAGER:
		return "项目管理"
	case approve.OperateTypeEnum_VIS_MANAGER:
		return "3D云应用"
	case approve.OperateTypeEnum_SECURITY_APPROVAL:
		return "安全审批"
	default:
		return ""
	}
}

type OperateUserType int8

const (
	OperateUserTypeUser     OperateUserType = 1
	OperateUserTypeAdmin    OperateUserType = 2
	OperateUserTypeSecurity OperateUserType = 3
)
