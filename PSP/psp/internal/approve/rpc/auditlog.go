package rpc

import (
	"context"

	"github.com/yuansuan/ticp/PSP/psp/internal/approve/dto"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/approve"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
)

func (s *GRPCService) SaveAuditLogInfo(ctx context.Context, in *approve.SaveAuditLogRequest) (*approve.AuditLogEmptyReply, error) {
	go s.AuditLogService.SaveLog(context.Background(), &dto.SaveLogRequest{
		UserId:         snowflake.MustParseString(in.UserId),
		UserName:       in.Username,
		OperateType:    in.OperateType,
		OperateContent: in.OperateContent,
		IpAddress:      in.IpAddress,
	})

	return &approve.AuditLogEmptyReply{}, nil
}
