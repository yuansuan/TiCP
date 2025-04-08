package rpc

import (
	"context"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/sysconfig"
	sysconfigdto "github.com/yuansuan/ticp/PSP/psp/internal/sysconfig/dto"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/strutil"
)

func (s *GRPCService) GetJobConfig(ctx context.Context, in *sysconfig.GetJobConfigRequest) (*sysconfig.GetJobConfigResponse, error) {
	jobConfig, err := s.SysConfigService.GetJobConfig(ctx)
	if err != nil {
		return nil, err
	}

	return &sysconfig.GetJobConfigResponse{
		Queue: jobConfig.Queue,
	}, nil
}

func (s *GRPCService) SetJobConfig(ctx context.Context, in *sysconfig.SetJobConfigRequest) (*sysconfig.SetJobConfigResponse, error) {
	err := s.SysConfigService.SetJobConfig(ctx, &sysconfigdto.SetJobConfigRequest{
		Queue: in.Queue,
	})
	if err != nil {
		return nil, err
	}

	return &sysconfig.SetJobConfigResponse{}, nil
}

func (s *GRPCService) GetJobBurstConfig(ctx context.Context, in *sysconfig.GetJobBurstConfigRequest) (*sysconfig.GetJobBurstConfigResponse, error) {
	jobBurstConfig, err := s.SysConfigService.GetJobBurstConfig(ctx)
	if err != nil {
		return nil, err
	}

	return &sysconfig.GetJobBurstConfigResponse{
		Enable:    jobBurstConfig.Enable,
		Threshold: jobBurstConfig.Threshold,
	}, nil
}

func (s *GRPCService) SetJobBurstConfig(ctx context.Context, in *sysconfig.SetJobBurstConfigRequest) (*sysconfig.SetJobBurstConfigResponse, error) {
	err := s.SysConfigService.SetJobBurstConfig(ctx, &sysconfigdto.SetJobBurstConfigRequest{
		Enable:    in.Enable,
		Threshold: in.Threshold,
	})
	if err != nil {
		return nil, err
	}

	return &sysconfig.SetJobBurstConfigResponse{}, nil
}

func (s *GRPCService) GetRBACDefaultRoleId(ctx context.Context, in *sysconfig.GetRBACDefaultRoleIdRequest) (*sysconfig.GetRBACDefaultRoleIdResponse, error) {
	defaultRoleId, err := s.SysConfigService.GetRBACDefaultRoleId(ctx)
	if err != nil {
		return nil, err
	}

	return &sysconfig.GetRBACDefaultRoleIdResponse{RoleId: defaultRoleId}, nil
}

func (s *GRPCService) SetRBACDefaultRoleId(ctx context.Context, in *sysconfig.SetRBACDefaultRoleIdRequest) (*sysconfig.SetRBACDefaultRoleIdResponse, error) {
	err := s.SysConfigService.SetRBACDefaultRoleId(ctx, in.RoleId)
	if err != nil {
		return nil, err
	}

	return &sysconfig.SetRBACDefaultRoleIdResponse{}, nil
}

func (s *GRPCService) GetThreePersonDefaultUserId(ctx context.Context, in *sysconfig.GetThreePersonDefaultUserIdRequest) (*sysconfig.GetThreePersonDefaultUserIdResponse, error) {
	threePersonConfig, err := s.SysConfigService.GetThreePersonManagementConfig(ctx)
	if err != nil {
		return nil, err
	}

	if strutil.IsEmpty(threePersonConfig.DefSafeUserID) {
		return &sysconfig.GetThreePersonDefaultUserIdResponse{}, nil
	}

	return &sysconfig.GetThreePersonDefaultUserIdResponse{UserId: snowflake.MustParseString(threePersonConfig.DefSafeUserID).Int64()}, nil
}
