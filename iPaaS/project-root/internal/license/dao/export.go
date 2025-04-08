package dao

import (
	"context"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/common/project-root-api/proto/idgen"
	dbModels "github.com/yuansuan/ticp/iPaaS/project-root/internal/license/dao/models"
)

//go:generate mockgen -self_package=github.com/yuansuan/ticp/iPaaS/project-root/internal/license/dao -destination mock_license_dao.go -package dao github.com/yuansuan/ticp/iPaaS/project-root/internal/license/dao LicenseManagerDao

// LicenseManagerDao dao
type LicenseManagerDao interface {
	ListLicenseManagers(ctx context.Context, params *ListAllLicenses) (list []*dbModels.JoinEntity, total int64, err error)
	AddLicenseManager(ctx context.Context, models *dbModels.LicenseManager) (err error)
	UpdateLicenseManager(ctx context.Context, models *dbModels.LicenseManager) (suc bool, err error)
	DeleteLicenseManager(ctx context.Context, lmId snowflake.ID) (suc bool, canDelete bool, err error)
	GetLicenseManager(ctx context.Context, lmId snowflake.ID) (lm *dbModels.LicenseManagerExt, err error)

	AddLicenseInfo(ctx context.Context, licInfo *dbModels.LicenseInfo) error
	DeleteLicenseInfo(ctx context.Context, licId snowflake.ID) (suc bool, canDelete bool, err error)
	UpdateLicenseInfo(ctx context.Context, licInfo *dbModels.LicenseInfo) (suc bool, err error)
	GetLicenseInfoByID(ctx context.Context, licId snowflake.ID) (existed bool, lic *dbModels.LicenseInfo, err error)
	SetLicenseServerStatus(ctx context.Context, licId snowflake.ID, status string) (suc bool, err error)

	ListModuleConfig(ctx context.Context, licId snowflake.ID) (list []*dbModels.ModuleConfig, err error)
	AddModuleConfig(ctx context.Context, cfg *dbModels.ModuleConfig) (suc bool, err error)
	BatchAddModuleConfigs(ctx context.Context, moduleConfigs []*dbModels.ModuleConfig) error
	DeleteModuleConfig(ctx context.Context, id snowflake.ID) (suc bool, candelete bool, err error)
	UpdateModuleConfigTotal(ctx context.Context, cfg *dbModels.ModuleConfig) (suc bool, err error)
	UpdateModuleConfigActual(ctx context.Context, cfg *dbModels.ModuleConfig) (suc bool, err error)

	IsAppUsed(ctx context.Context, appID int64) (bool, error)
	IsJobUsed(ctx context.Context, jobID snowflake.ID) (exit bool, licenseJob []*dbModels.LicenseJob, err error)
	IsLicenseUsed(ctx context.Context, moduleID snowflake.ID) (licenseJobs []*dbModels.LicenseJob, err error)
	AcquireLicense(ctx context.Context, jobID snowflake.ID, idGenClient idgen.IdGenClient, lic *dbModels.LicenseInfoExt, required map[string]int) (err error)
	ReleaseLicense(ctx context.Context, jobID snowflake.ID, licenseJob []*dbModels.LicenseJob) (err error)
	SelectByManagerID(ctx context.Context, managerID snowflake.ID) (fromDB []*dbModels.LicenseInfo, err error)
	SelectByInfoIDs(ctx context.Context, ids []snowflake.ID) (fromDB []*dbModels.ModuleConfig, err error)
	LicenseUsed(ctx context.Context, managerId, size, index int64) (results []*dbModels.LicenseUsedResult, total int64, err error)
}
