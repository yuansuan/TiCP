package service

import (
	"context"

	"github.com/yuansuan/ticp/PSP/psp/internal/license/dto"
)

type LicenseManagerService interface {
	// LicenseManagerList  列表
	LicenseManagerList(ctx context.Context, licenseType string) (*dto.LicenseManagerListResponse, error)
	// LicenseManagerInfo  详情
	LicenseManagerInfo(ctx context.Context, managerID string) (*dto.LicenseManagerData, error)
	// LicenseTypeList license类型列表
	LicenseTypeList(ctx context.Context) (*dto.LicenseTypeListResponse, error)
	// AddLicenseManager 新增license manager
	AddLicenseManager(ctx context.Context, req *dto.AddLicenseManagerRequest) (*dto.AddLicenseManagerResponse, error)
	// EditLicenseManager 编辑 license manager
	EditLicenseManager(ctx context.Context, req *dto.EditLicenseManagerRequest) error
	// DeleteLicenseManager 删除 license manager
	DeleteLicenseManager(ctx context.Context, managerID string) error
}

type LicenseInfoService interface {
	// AddLicenseInfo 新增license Info
	AddLicenseInfo(ctx context.Context, req *dto.LicenseInfoAddRequest) (*dto.LicenseInfoAddResponse, error)
	// EditLicenseInfo 修改license Info
	EditLicenseInfo(ctx context.Context, req *dto.LicenseInfoEditRequest) error
	// DeleteLicenseInfo 删除license Info
	DeleteLicenseInfo(ctx context.Context, licenseID string) error
}

type ModuleConfigService interface {
	// ModuleConfigList 模块使用数量详情列表
	ModuleConfigList(ctx context.Context, licenseID string) (*dto.ModuleConfigListResponse, error)
	// AddModuleConfig 新增 module config
	AddModuleConfig(ctx context.Context, req *dto.AddModuleConfigRequest) (*dto.AddModuleConfigResponse, error)
	// EditModuleConfig 编辑 module config
	EditModuleConfig(ctx context.Context, req *dto.EditModuleConfigRequest) error
	// DeleteModuleConfig 删除 module config
	DeleteModuleConfig(ctx context.Context, moduleConfigID string) error
}
