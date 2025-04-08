package collector

import (
	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/metrics"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi"
	"github.com/yuansuan/ticp/PSP/psp/internal/monitor/consts"
)

type featureCollector struct {
	metric  GaugeDesc
	logger  log.Logger
	openapi *openapi.OpenAPI
}

func init() {
	Register(consts.Feature, NewFeatureCollector)
}

func NewFeatureCollector(logger log.Logger) (Collector, error) {
	api, err := openapi.NewLocalAPI()
	if err != nil {
		return nil, err
	}
	return &featureCollector{
		metric: GaugeDesc{
			Desc: NewDesc(consts.Feature, "feature usage info", []string{"app_type", "value_type", "license_id", "feature_name"}),
		},
		logger:  logger,
		openapi: api,
	}, nil
}

func (c *featureCollector) UpdateMetrics(ch chan<- prometheus.Metric) error {
	licenseManager, err := c.openapi.Client.License.ListLicenseManager()
	if err != nil {
		logging.Default().Errorf("err: %v", err)
		return err
	}
	for _, item := range licenseManager.Data.Items {
		for _, info := range item.LicenseInfos {
			for _, config := range info.ModuleConfigs {
				ch <- c.metric.MustNewConstMetric(float64(config.UsedNum), item.AppType, metrics.FeatureUsage, info.Id, config.ModuleName)
				ch <- c.metric.MustNewConstMetric(float64(config.Total), item.AppType, metrics.FeatureTotal, info.Id, config.ModuleName)
				ch <- c.metric.MustNewConstMetric(float64(config.Total-config.UsedNum), item.AppType, metrics.FeatureAvailable, info.Id, config.ModuleName)

				if config.UsedNum == 0 || config.Total == 0 {
					ch <- c.metric.MustNewConstMetric(0, item.AppType, metrics.FeatureUsagePercent, info.Id, config.ModuleName)
				} else if config.UsedNum != 0 && config.Total == 0 {
					ch <- c.metric.MustNewConstMetric(100, item.AppType, metrics.FeatureUsagePercent, info.Id, config.ModuleName)
				} else {
					ch <- c.metric.MustNewConstMetric((float64(config.UsedNum)/float64(config.Total))*100, item.AppType, metrics.FeatureUsagePercent, info.Id, config.ModuleName)
				}
			}
		}
	}
	return nil
}
