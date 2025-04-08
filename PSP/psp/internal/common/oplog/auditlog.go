package oplog

import (
	"context"
	"sync"

	"github.com/gin-gonic/gin"
	grpc_boot "github.com/yuansuan/ticp/common/go-kit/gin-boot/grpc-boot"
	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/approve"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/ginutil"
)

// LogService 作业服务依赖的rpc服务客户端
type LogService struct {
	Approve approve.AuditLogManagementClient `grpc_client_inject:"approve"`
}

var (
	once   sync.Once
	client *LogService
)

// GetInstance 获取客户端
func GetInstance() *LogService {
	once.Do(func() {
		client = &LogService{}
		grpc_boot.InjectAllClient(client)
	})

	return client
}

// SaveAuditLogInfo 保存操作日志
func (srv *LogService) SaveAuditLogInfo(ctx *gin.Context, opType approve.OperateTypeEnum, content string) error {
	userName := ginutil.GetUserName(ctx)
	saveAuditLogRequest := &approve.SaveAuditLogRequest{
		UserId:         snowflake.ID(ginutil.GetUserID(ctx)).String(),
		Username:       userName,
		IpAddress:      ginutil.GetRequestIP(ctx),
		OperateType:    opType,
		OperateContent: content,
	}
	if _, err := srv.Approve.SaveAuditLogInfo(ctx, saveAuditLogRequest); err != nil {
		logging.Default().Errorf("failed to save audit log info %+v, error: %v", saveAuditLogRequest, err)
		return err
	}
	return nil
}

// SaveAuditLogInfoGrpc 保存操作日志（grpc用）
func (srv *LogService) SaveAuditLogInfoGrpc(ctx context.Context, opType approve.OperateTypeEnum, userID snowflake.ID, userName, ipAddress, content string) {
	saveAuditLogRequest := &approve.SaveAuditLogRequest{
		UserId:         userID.String(),
		Username:       userName,
		IpAddress:      ipAddress,
		OperateType:    opType,
		OperateContent: content,
	}
	if _, err := srv.Approve.SaveAuditLogInfo(ctx, saveAuditLogRequest); err != nil {
		logging.Default().Errorf("failed to save audit log info %+v, error: %v", saveAuditLogRequest, err)
	}
}
