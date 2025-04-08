package collector

import (
	"fmt"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"

	"github.com/yuansuan/ticp/PSP/psp/internal/agent/consts"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/metrics"
)

type memoryCollector struct {
	metric GaugeDesc
	logger log.Logger
}

func init() {
	Register(consts.MemoryMetric, NewMemoryCollector)
}

func NewMemoryCollector(logger log.Logger) (Collector, error) {
	return &memoryCollector{
		metric: GaugeDesc{
			Desc: NewDesc(consts.MemoryMetric, "memory info", []string{"name"}),
		},
		logger: logger,
	}, nil
}

func (c *memoryCollector) UpdateMetrics(ch chan<- prometheus.Metric) error {
	// 获取内存信息（空闲/最大）
	memInfo := c.getMemoryInfo()

	// 交换空间信息（空闲/最大）
	swapInfo := c.getSwapInfo()

	// tmp空间（空闲/最大）
	fsStat := c.getDiskUsage()

	ch <- c.metric.MustNewConstMetric(float64(memInfo.Available/(1024*1024)), metrics.MemoryAvailable)
	ch <- c.metric.MustNewConstMetric(float64(memInfo.Total/(1024*1024)), metrics.MemoryTotal)
	ch <- c.metric.MustNewConstMetric(float64(memInfo.Free/(1024*1024)), metrics.MemoryFree)
	ch <- c.metric.MustNewConstMetric(float64(memInfo.Used/(1024*1024)), metrics.MemoryUsed)
	ch <- c.metric.MustNewConstMetric((float64((memInfo.Total-memInfo.Available)/(1024*1024))/float64(memInfo.Total/(1024*1024)))*100, metrics.MemoryPercent)
	ch <- c.metric.MustNewConstMetric(float64(swapInfo.Free/(1024*1024)), metrics.MemorySwapFree)
	ch <- c.metric.MustNewConstMetric(float64(swapInfo.Total/(1024*1024)), metrics.MemorySwapTotal)
	ch <- c.metric.MustNewConstMetric(float64(fsStat.Free/(1024*1024)), metrics.MemoryTmpFree)
	ch <- c.metric.MustNewConstMetric(float64(fsStat.Total/(1024*1024)), metrics.MemoryTmpTotal)

	return nil
}

// 获取内存信息（空闲/最大）
func (c *memoryCollector) getMemoryInfo() *mem.VirtualMemoryStat {
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		level.Error(c.logger).Log("msg", fmt.Sprintf("unable to get memory information (Free/Max):%v", err))
		return nil
	}
	return memInfo
}

// 获取交换空间信息
func (c *memoryCollector) getSwapInfo() *mem.SwapMemoryStat {
	swapInfo, err := mem.SwapMemory()
	if err != nil {
		level.Error(c.logger).Log("msg", fmt.Sprintf("unable to get exchange space information:%v", err))
		return nil
	}
	return swapInfo
}

// 获取磁盘空间使用情况
func (c *memoryCollector) getDiskUsage() *disk.UsageStat {
	path := "/tmp"
	fsStat, err := disk.Usage(path)
	if err != nil {
		level.Error(c.logger).Log("msg", fmt.Sprintf("unable to get the usage of the %s directory: %s", err))
		return nil
	}
	return fsStat
}
