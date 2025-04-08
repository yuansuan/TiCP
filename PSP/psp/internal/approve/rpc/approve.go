package rpc

import (
	"context"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/approve"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
)

func (s *GRPCService) CheckUnhandledApprove(ctx context.Context, in *approve.CheckUnhandledApproveRequest) (*approve.CheckUnhandledApproveResponse, error) {
	response := &approve.CheckUnhandledApproveResponse{}
	unhandled, err := s.ApproveService.CheckUnhandledApprove(ctx, snowflake.MustParseString(in.UserId).Int64())
	if err != nil {
		return response, err
	}

	response.Unhandled = unhandled
	return response, nil
}
