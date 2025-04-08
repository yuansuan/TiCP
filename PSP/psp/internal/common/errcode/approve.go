package errcode

import (
	"google.golang.org/grpc/codes"
)

// App 错误码范围: 21001 ~ 22000
// 命名格式：Err + 服务名 + 具体错误
const (
	ErrApproveAuditLogAdd     codes.Code = 21001
	ErrApproveAuditLogList    codes.Code = 21002
	ErrApproveApplicationList codes.Code = 21003
	ErrApprovePendingList     codes.Code = 21004
	ErrApprovedList           codes.Code = 21005
)

const (
	ErrApproveInvalidType            codes.Code = 21100
	ErrApproveHasConflict            codes.Code = 21101
	ErrApproveRecordNotExist         codes.Code = 21102
	ErrApprovedApply                 codes.Code = 21103
	ErrApprovedPass                  codes.Code = 21104
	ErrApprovedRefuse                codes.Code = 21105
	ErrApproveStatusEnd              codes.Code = 21106
	ErrApproveNecessaryRoleNotExist  codes.Code = 21107
	ErrApproveNecessaryPermNotExist  codes.Code = 21108
	ErrApproveNecessaryRoleNameExist codes.Code = 21109
	ErrApproveUnhandledExist         codes.Code = 21110
	ErrUserApproveSelf               codes.Code = 21111
)

var ApproveCodeMsg = map[codes.Code]string{
	ErrApproveAuditLogAdd:            "日志记录失败",
	ErrApproveAuditLogList:           "获取日志失败",
	ErrApproveApplicationList:        "获取申请列表失败",
	ErrApprovePendingList:            "获取待审批列表失败",
	ErrApprovedList:                  "获取已审批列表失败",
	ErrApproveInvalidType:            "不支持的审批类型",
	ErrApproveHasConflict:            "与正在进行中的审批存在冲突，请等待审批完成",
	ErrApproveRecordNotExist:         "审批记录不存在",
	ErrApproveNecessaryRoleNotExist:  "给用户分配的角色已不存在，此条审批将失效",
	ErrApproveNecessaryPermNotExist:  "给角色分配的权限已不存在，此条审批将失效",
	ErrApproveNecessaryRoleNameExist: "角色名称与现有其他角色重复，此条审批将失效",
	ErrApprovedApply:                 "发起审批失败",
	ErrApprovedPass:                  "通过审批失败",
	ErrApprovedRefuse:                "拒绝审批失败",
	ErrApproveStatusEnd:              "此审批已经结束",
	ErrApproveUnhandledExist:         "该用户存在未处理的审批，请等待审批结束再进行操作",
	ErrUserApproveSelf:               "审批人不能审批自身用户",
}
