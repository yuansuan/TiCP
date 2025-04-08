package collector

import (
	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/metrics"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/openapi"
	"github.com/yuansuan/ticp/PSP/psp/internal/monitor/consts"
	"github.com/yuansuan/ticp/PSP/psp/internal/monitor/service/impl"
)

type diskCollector struct {
	metric  GaugeDesc
	logger  log.Logger
	openapi *openapi.OpenAPI
}

func init() {
	Register(consts.DiskUsage, NewDiskCollector)
}

func NewDiskCollector(logger log.Logger) (Collector, error) {
	api, err := openapi.NewLocalHPCAPI()
	if err != nil {
		return nil, err
	}
	return &diskCollector{
		metric: GaugeDesc{
			Desc: NewDesc(consts.DiskUsage, "disk usage info", []string{"name", "mount_path"}),
		},
		logger:  logger,
		openapi: api,
	}, nil
}

func (c *diskCollector) UpdateMetrics(ch chan<- prometheus.Metric) error {
	//1.查询磁盘信息
	fields, diskMaps, err := impl.GetDiskData(nil, c.openapi)
	if err != nil {
		return err
	}

	//2.将磁盘信息写入到ch中
	for _, field := range fields {
		usedValue := 0.0
		unUsedValue := 0.0
		for _, diskMap := range diskMaps {
			if value, ok := diskMap[field]; ok {
				if diskMap[consts.Name].(string) == consts.Used {
					usedValue = float64(value.(int64))
				}
				if diskMap[consts.Name].(string) == consts.UnUsed {
					unUsedValue = float64(value.(int64))
				}
				ch <- c.metric.MustNewConstMetric(float64(value.(int64)), diskMap[consts.Name].(string), field)
			}
		}

		if usedValue == 0 {
			ch <- c.metric.MustNewConstMetric(0, metrics.DiskUsagePercent, field)
		} else if usedValue != 0 && unUsedValue == 0 {
			ch <- c.metric.MustNewConstMetric(100, metrics.DiskUsagePercent, field)
		} else {
			ch <- c.metric.MustNewConstMetric((usedValue/(usedValue+unUsedValue))*100, metrics.DiskUsagePercent, field)
		}
	}
	return nil
}
