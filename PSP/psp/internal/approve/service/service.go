package service

import (
	"context"

	"github.com/gin-gonic/gin"

	"github.com/yuansuan/ticp/PSP/psp/internal/approve/dto"
)

type AuditLogService interface {
	SaveLog(ctx context.Context, req *dto.SaveLogRequest) error
	List(ctx *gin.Context, req *dto.AuditLogListRequest) (*dto.AuditLogListResponse, error)
	ListAll(ctx *gin.Context, req *dto.AuditLogListAllRequest) (*dto.AuditLogListAllResponse, error)
	Export(ctx *gin.Context, req *dto.AuditLogListRequest) error
	ExportAll(ctx *gin.Context, req *dto.AuditLogListAllRequest) error
}

type ApproveService interface {
	ApplyApprove(ctx *gin.Context, req *dto.ApplyApproveRequest) error
	CancelApprove(ctx *gin.Context, reqId int64) error
	GetApproveList(ctx context.Context, UserId int64, Request *dto.GetApproveListRequest) (*dto.ApproveLogListAllResponse, error)
	GetApprovePendingList(ctx context.Context, UserId int64, Request *dto.GetApprovePendingRequest) (*dto.ApproveLogListAllResponse, error)
	GetApprovedList(ctx context.Context, UserId int64, Request *dto.GetApproveCompleteRequest) (*dto.ApproveLogListAllResponse, error)
	Pass(ctx *gin.Context, req *dto.HandleApproveRequest) error
	Refuse(ctx *gin.Context, req *dto.HandleApproveRequest) error
	CheckUnhandledApprove(ctx context.Context, userId int64) (bool, error)
}
