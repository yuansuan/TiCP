package collector

import (
	"fmt"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/disk"

	"github.com/yuansuan/ticp/PSP/psp/internal/agent/consts"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/metrics"
)

const (
	DiskThroughputTime = 1
)

type diskCollector struct {
	metric GaugeDesc
	logger log.Logger
}

func init() {
	Register(consts.DiskMetric, NewDiskCollector)
}

func NewDiskCollector(logger log.Logger) (Collector, error) {
	return &diskCollector{
		metric: GaugeDesc{
			Desc: NewDesc(consts.DiskMetric, "disk info", []string{"name"}),
		},
		logger: logger,
	}, nil
}

func (c *diskCollector) UpdateMetrics(ch chan<- prometheus.Metric) error {
	logger := c.logger

	lastStats, err := disk.IOCounters()
	if err != nil {
		level.Error(logger).Log("msg", fmt.Sprintf("could not get initial IO counter data:%v", err))
		return err
	}

	// 等待一段时间，再次获取 IO 计数器数据
	time.Sleep(time.Duration(DiskThroughputTime) * time.Second)

	currentStats, err := disk.IOCounters()
	if err != nil {
		level.Error(logger).Log("msg", fmt.Sprintf("unable to get current IO counter data:%v", err))
		return err
	}

	duration := time.Duration(DiskThroughputTime) * time.Second

	var totalReadBytes, totalWriteBytes uint64
	for key, stat := range lastStats {
		totalReadBytes += currentStats[key].ReadBytes - stat.ReadBytes
		totalWriteBytes += currentStats[key].WriteBytes - stat.WriteBytes
	}

	readThroughput := float64(totalReadBytes) / duration.Seconds()
	writeThroughput := float64(totalWriteBytes) / duration.Seconds()

	ch <- c.metric.MustNewConstMetric(readThroughput/1024, metrics.DiskReadThroughput)
	ch <- c.metric.MustNewConstMetric(writeThroughput/1024, metrics.DiskWriteThroughput)
	ch <- c.metric.MustNewConstMetric((writeThroughput+readThroughput)/1024, metrics.DiskTotalThroughput)
	return nil
}
