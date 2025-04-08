package scheduler

import (
	"context"
	"errors"

	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/license/collector"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/license/dao"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/license/dao/models"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/license/handler_rpc/impl"
)

const (
	// licenseServerStatusErrorThreshold license服务状态异常阈值
	licenseServerStatusErrorThreshold = 3
)

var (
	// license服务状态异常计数器
	licenseServerStatusErrorCounter = make(map[snowflake.ID]int)
)

type Scheduler struct {
}

func NewScheduler() *Scheduler {
	return &Scheduler{}
}

func (Scheduler *Scheduler) Run(ctx context.Context) {
	logging.GetLogger(ctx).Info("StartUpdateOtherOwnedLicense")
	licImpl := dao.NewLicenseImpl(boot.MW.DefaultORMEngine())
	//查询所有已发布的license
	publishedLic, err := licImpl.LicenseInfoPublished(ctx)
	if err != nil {
		logging.Default().Errorf("LicenseInfoPublished, Error: %s", err.Error())
		return
	}
	for _, lic := range publishedLic {
		// 初始化license服务状态
		initLicenseServerStatus(lic)
		if lic.Auth != 1 {
			continue
		}
		if lic.HpcEndpoint != "" && lic.CollectorType != "" && lic.ToolPath != "" {
			components, err := impl.GetRemainLicense(lic)
			checkLicenseServerStatus(ctx, lic, err)
			if err != nil {
				logging.Default().Warnf("GetRemainLicense, Error: %s, LicenseId: %d", err.Error(), lic.Id)
				continue
			}
			for name, num := range components {
				moduleCfg := models.ModuleConfig{
					LicenseId:   lic.Id,
					ModuleName:  name,
					ActualTotal: int(num.Total),
					ActualUsed:  int(num.Used),
				}
				succ, err := licImpl.UpdateModuleConfigActual(ctx, &moduleCfg)
				if err != nil {
					logging.Default().Warnf("UpdateModuleConfig, Error: %s, LicenseId: %d, ModuleName: %s",
						err.Error(), lic.Id, name)
				}
				if !succ {
					logging.Default().Infof("UpdateModuleConfig, LicenseId: %s, ModuleConfigName: %s",
						lic.Id, name)
				} else {
					logging.GetLogger(ctx).Infof("UpdateModuleConfig, ModuleName: %s, Num: %f, LicenseId: %s",
						name, num, lic.Id)
				}
			}
		}
	}
}

func initLicenseServerStatus(lic *models.LicenseInfo) {
	if _, ok := licenseServerStatusErrorCounter[lic.Id]; !ok {
		licenseServerStatusErrorCounter[lic.Id] = 0
	}
}

func checkLicenseServerStatus(ctx context.Context, licInfo *models.LicenseInfo, collectorErr error) {
	licImpl := dao.NewLicenseImpl(boot.MW.DefaultORMEngine())
	// 正常之后，如果之前是异常状态，则恢复正常
	if collectorErr == nil {
		if licInfo.LicenseServerStatus == models.LicenseServerStatusAbnormal {
			if _, err := licImpl.SetLicenseServerStatus(ctx, licInfo.Id, models.LicenseServerStatusNormal); err != nil {
				logging.Default().Warnf("set server status fail, error: %v, licinfo: %v", err, licInfo)
			}
			logging.Default().Infof("license server status changes to a normal, licinfo: %v", licInfo)
			licenseServerStatusErrorCounter[licInfo.Id] = 0
		}
		return
	}
	// 不是CollectorRuntimeErr，无需处理
	if !errors.Is(collectorErr, collector.CollectorRuntimeErr) {
		return
	}
	// 异常状态, 并且之前是正常状态
	licenseServerStatusErrorCounter[licInfo.Id]++
	if licInfo.LicenseServerStatus == models.LicenseServerStatusNormal {
		if licenseServerStatusErrorCounter[licInfo.Id] >= licenseServerStatusErrorThreshold {
			if _, err := licImpl.SetLicenseServerStatus(ctx, licInfo.Id, models.LicenseServerStatusAbnormal); err != nil {
				logging.Default().Warnf("set server status fail, error: %v, licinfo: %v", err, licInfo)
			}
			logging.Default().Infof("license server status changes to a abnormal, licinfo: %v", licInfo)
		}
	}
}
