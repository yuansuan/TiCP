package service

import (
	"context"

	"github.com/yuansuan/ticp/PSP/psp/internal/app/dto"
)

type AppService interface {
	ListApp(ctx context.Context, userId int64, computeType, state string, hasPermission, desktop bool) ([]*dto.App, error)
	AddApp(ctx context.Context, req *dto.AddAppServiceRequest) error
	UpdateApp(ctx context.Context, appDto *dto.App, baseName string) error
	DeleteApp(ctx context.Context, name, computeType string) error
	PublishApp(ctx context.Context, names []string, computeType, state string) error
	GetAppInfo(ctx context.Context, req *dto.GetAppInfoServiceRequest) (*dto.App, error)
	GetAppTotalNum(ctx context.Context) (int64, error)
	SyncAppContent(ctx context.Context, baseAppId string, syncAppIds []string) error

	ListZone(ctx context.Context) ([]string, error)
	ListQueue(ctx context.Context, appId string) ([]*dto.QueueInfo, error)
	ListLicense(ctx context.Context) ([]*dto.LicenseInfo, error)
	GetSchedulerResourceKey(ctx context.Context) ([]string, error)
	GetSchedulerResourceValue(ctx context.Context, appId, resourceType, resourceSubType string) ([]*dto.Item, error)
	CheckLicenseManagerIdUsed(ctx context.Context, licenseManagerId string) (bool, error)
}
