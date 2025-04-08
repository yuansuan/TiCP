package collector

import (
	"fmt"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/cpu"

	"github.com/yuansuan/ticp/PSP/psp/internal/agent/consts"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/metrics"
)

const (
	CpuIdleThreshold = 10
)

type cpuCollector struct {
	metric GaugeDesc
	logger log.Logger
}

func init() {
	Register(consts.CPUMetric, NewCpuCollector)
}

func NewCpuCollector(logger log.Logger) (Collector, error) {
	return &cpuCollector{
		metric: GaugeDesc{
			Desc: NewDesc(consts.CPUMetric, "cpu info", []string{"name"}),
		},
		logger: logger,
	}, nil
}

func (c *cpuCollector) UpdateMetrics(ch chan<- prometheus.Metric) error {
	// 逻辑核数
	cores := c.getLogicalCores()

	// 空闲逻辑核数
	idleCores := c.getIdleCores()

	// CPU使用率
	percentage := c.getCPUUsage()

	// CPU空闲时间
	idleTime := c.getIdleTime()

	ch <- c.metric.MustNewConstMetric(float64(cores), metrics.CPUCore)
	ch <- c.metric.MustNewConstMetric(float64(idleCores), metrics.CPUIdleCore)
	ch <- c.metric.MustNewConstMetric(percentage, metrics.CPUPercent)
	ch <- c.metric.MustNewConstMetric(idleTime, metrics.CPUIdleTime)

	//节点状态
	ch <- c.metric.MustNewConstMetric(metrics.Normal, metrics.NodeStatus)
	return nil
}

// 获取逻辑核数
func (c *cpuCollector) getLogicalCores() int {
	cores, err := cpu.Counts(true)
	if err != nil {
		level.Error(c.logger).Log("msg", fmt.Sprintf("unable to obtain logical kernel:%v", err))
		return 0
	}
	return cores
}

// 获取 CPU 空闲核数
func (c *cpuCollector) getIdleCores() int {
	cpuPercent, err := cpu.Percent(0, true) // 获取每个逻辑核心的使用率
	if err != nil {
		level.Error(c.logger).Log("msg", fmt.Sprintf("failed to obtain the number of idle CPU cores:%v", err))
		return 0
	}
	idleCores := 0

	for _, usage := range cpuPercent {
		if usage < CpuIdleThreshold {
			idleCores++
		}
	}
	return idleCores
}

// CPU使用率
func (c *cpuCollector) getCPUUsage() float64 {
	percentage, err := cpu.Percent(time.Second, false)
	if err != nil {
		level.Error(c.logger).Log("msg", fmt.Sprintf("failed to obtain CPU usage:%v", err))
		return 0
	}
	return percentage[0]
}

// 获取 CPU 空闲时间
func (c *cpuCollector) getIdleTime() float64 {
	idleTimes, err := cpu.Times(false)
	if err != nil {
		level.Error(c.logger).Log("msg", fmt.Sprintf("failed to get CPU times:%v", err))
		return -1
	}
	if len(idleTimes) == 0 {
		return -1
	}

	idleTime := idleTimes[0].Idle
	return idleTime
}
