package service

import (
	"context"

	"github.com/yuansuan/ticp/PSP/psp/internal/sysconfig/dto"
)

type SysConfigService interface {
	GetGlobalSysConfig(ctx context.Context) (*dto.GetGlobalSysConfigResponse, error)

	GetJobConfig(ctx context.Context) (*dto.GetJobConfigResponse, error)
	SetJobConfig(ctx context.Context, req *dto.SetJobConfigRequest) error
	GetJobBurstConfig(ctx context.Context) (*dto.GetJobBurstConfigResponse, error)
	SetJobBurstConfig(ctx context.Context, req *dto.SetJobBurstConfigRequest) error

	GetRBACDefaultRoleId(ctx context.Context) (int64, error)
	SetRBACDefaultRoleId(ctx context.Context, defaultRoleId int64) error

	SetEmailConfig(ctx context.Context, in *dto.SetEmailConfigReq) error
	GetEmailConfig(ctx context.Context) (*dto.GetEmailConfigRes, error)

	SetGlobalEmail(ctx context.Context, in *dto.EmailConfig) error
	GetGlobalEmail(ctx context.Context) (*dto.EmailConfig, error)
	SendEmail(ctx context.Context) error
	GetThreePersonManagementConfig(ctx context.Context) (*dto.GetThreePersonConfigResponse, error)
	SetThreePersonManagementConfig(ctx context.Context, req *dto.SetThreePersonConfigRequest) error
}
