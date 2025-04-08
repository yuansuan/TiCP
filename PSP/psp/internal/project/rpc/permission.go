package rpc

import (
	"context"

	pb "github.com/yuansuan/ticp/PSP/psp/internal/common/proto/project"
	"github.com/yuansuan/ticp/PSP/psp/internal/project/util"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
)

// ExistsProjectMember 判断项目成员是否存在
func (s *GRPCService) CheckUserOperatorProjectsPermission(ctx context.Context, in *pb.CheckUserOperatorProjectsPermissionRequest) (*pb.CheckUserOperatorProjectsPermissionResponse, error) {
	userId := snowflake.MustParseString(in.UserId)
	pass, err := util.CheckProjectAdminRole(ctx, userId, true)
	if err != nil {
		return nil, err
	}

	return &pb.CheckUserOperatorProjectsPermissionResponse{
		Pass: pass,
	}, nil
}
