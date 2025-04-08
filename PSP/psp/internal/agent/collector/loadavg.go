package collector

import (
	"fmt"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/load"

	"github.com/yuansuan/ticp/PSP/psp/internal/agent/consts"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/metrics"
)

type loadavgCollector struct {
	metric GaugeDesc
	logger log.Logger
}

func init() {
	Register(consts.LoadavgMetric, NewLoadavgCollector)
}

func NewLoadavgCollector(logger log.Logger) (Collector, error) {
	return &loadavgCollector{
		metric: GaugeDesc{
			Desc: NewDesc(consts.LoadavgMetric, "load average info", []string{"name"}),
		},
		logger: logger,
	}, nil
}

func (c *loadavgCollector) UpdateMetrics(ch chan<- prometheus.Metric) error {
	// 获取负载信息
	loadAvg := c.getLoadAverage()

	ch <- c.metric.MustNewConstMetric(loadAvg.Load1, metrics.Load1m)
	ch <- c.metric.MustNewConstMetric(loadAvg.Load5, metrics.Load5m)
	ch <- c.metric.MustNewConstMetric(loadAvg.Load15, metrics.Load15m)

	return nil
}

// 获取负载信息
func (c *loadavgCollector) getLoadAverage() *load.AvgStat {

	loadAvg, err := load.Avg()
	if err != nil {
		level.Error(c.logger).Log("msg", fmt.Sprintf("failed to obtain load information:%v", err))
		return nil
	}
	return loadAvg
}
