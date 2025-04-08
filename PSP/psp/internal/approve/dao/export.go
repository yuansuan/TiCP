package dao

import (
	"context"

	"github.com/yuansuan/ticp/PSP/psp/internal/approve/dao/model"
	"github.com/yuansuan/ticp/PSP/psp/internal/approve/dto"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
	"github.com/yuansuan/ticp/PSP/psp/pkg/xtype"
)

type AuditLogDao interface {
	Add(ctx context.Context, log *model.AuditLog) error
	List(ctx context.Context, page *xtype.Page, userId int64, userName, ipAddress string, operateType string, startTime, endTime string) (logList []*model.AuditLog, total int64, err error)
	ListAll(ctx context.Context, page *xtype.Page, userId int64, userName, ipAddress string, operateType string, startTime, endTime string, operateUserType dto.OperateUserType) (logList []*model.AuditLog, total int64, err error)
}

type ApproveDao interface {
	// CancelApprove 取消审批
	CancelApprove(ctx context.Context, id int64) error
	// ApproveList 审批列表
	ApproveList(ctx context.Context, page *xtype.Page, condition *dto.ApproveListCondition) (list []*model.ApproveUserWithRecord, total int64, err error)
	// AddApproveRecord 添加审批记录
	AddApproveRecord(ctx context.Context, record *model.ApproveRecord) error
	// AddApproveUser 添加审批用户关联
	AddApproveUser(ctx context.Context, record *model.ApproveUser) error
	// CheckSign 检查签名是否重复
	CheckSign(ctx context.Context, sign string) (bool, error)
	// ApplicationList 申请列表
	ApplicationList(ctx context.Context, page *xtype.Page, condition *dto.ApplicationListCondition) (list []*model.ApproveUserWithRecord, total int64, err error)
	// UpdateApproveUser 更新审批用户关联
	UpdateApproveUser(ctx context.Context, approveUser *model.ApproveUser) error
	// AllApproved 判断是否还有未审批的记录 true-都审批了
	AllApproved(ctx context.Context, recordID snowflake.ID) (bool, error)
	// UpdateApproveRecord 更新审批记录
	UpdateApproveRecord(ctx context.Context, record *model.ApproveRecord) error
	// GetRecord 获取审批记录
	GetRecord(ctx context.Context, ID snowflake.ID) (*model.ApproveRecord, error)
	// CheckUnhandledApprove 检查用户是否还存在待处理的审批(自己发起的审批或者需要自己审批的审批)
	CheckUnhandledApprove(ctx context.Context, id int64) (bool, error)
}
