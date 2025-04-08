package collector

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/metrics"
)

var (
	hostname = ""

	once      = sync.Once{}
	mutex     = sync.Mutex{}
	factories = make(map[string]func(logger log.Logger) (Collector, error))

	scrapeDuration = GaugeDesc{
		Desc: NewDesc("collector_scrape_duration", "collector scrape duration, unit: seconds", []string{"name"}),
	}
	scrapeSuccess = GaugeDesc{
		Desc: NewDesc("collector_scrape_success", "whether collector scrape success", []string{"name"}),
	}
)

type GaugeDesc struct {
	Desc *prometheus.Desc
}

type CounterDesc struct {
	Desc *prometheus.Desc
}

type Collector interface {
	UpdateMetrics(ch chan<- prometheus.Metric) error
}

func NewCollector(logger log.Logger) (*MonitorCollector, error) {
	collectors := make(map[string]Collector, len(factories))

	mutex.Lock()
	defer mutex.Unlock()
	for key, _ := range factories {
		collector, err := factories[key](logger)
		if err != nil {
			return nil, err
		}
		collectors[key] = collector
	}
	return &MonitorCollector{
		Collectors: collectors,
		logger:     logger,
	}, nil
}

func NewDesc(fqName, help string, variableLabels []string) *prometheus.Desc {
	if variableLabels != nil {
		variableLabels = append(variableLabels, metrics.HostName)
	}

	name := prometheus.BuildFQName(metrics.Namespace, metrics.MonitorSystem, fqName)
	return prometheus.NewDesc(name, help, variableLabels, nil)
}

func Register(collectorName string, collector func(logger log.Logger) (Collector, error)) {
	mutex.Lock()
	defer mutex.Unlock()
	factories[collectorName] = collector
}

func (d *GaugeDesc) MustNewConstMetric(value float64, labels ...string) prometheus.Metric {
	if labels != nil {
		labels = append(labels, getHostname())
	}

	return mustNewConstMetric(d.Desc, prometheus.GaugeValue, value, labels...)
}

func (d *CounterDesc) MustNewConstMetric(value float64, labels ...string) prometheus.Metric {
	return mustNewConstMetric(d.Desc, prometheus.GaugeValue, value, labels...)
}

type MonitorCollector struct {
	Collectors map[string]Collector
	logger     log.Logger
}

func (n *MonitorCollector) Describe(ch chan<- *prometheus.Desc) {
}

func (n *MonitorCollector) Collect(ch chan<- prometheus.Metric) {
	wg := sync.WaitGroup{}
	wg.Add(len(n.Collectors))

	for name, c := range n.Collectors {
		go func(name string, c Collector) {
			executeCollector(name, c, ch, n.logger)
			wg.Done()
		}(name, c)
	}

	wg.Wait()
}

func executeCollector(collectorName string, collector Collector, ch chan<- prometheus.Metric, logger log.Logger) {
	startTime := time.Now()
	err := collector.UpdateMetrics(ch)
	duration := time.Since(startTime)

	var success int
	if err != nil {
		level.Error(logger).Log("msg", "collector failed", "name", collectorName, "duration_seconds", duration.Seconds(), "err", err)
	} else {
		level.Debug(logger).Log("msg", "collector succeeded", "name", collectorName, "duration_seconds", duration.Seconds())
		success = 1
	}

	ch <- scrapeSuccess.MustNewConstMetric(float64(success), collectorName)
	ch <- scrapeDuration.MustNewConstMetric(duration.Seconds(), collectorName)
}

func mustNewConstMetric(desc *prometheus.Desc, valueType prometheus.ValueType, value float64, labels ...string) prometheus.Metric {
	metric, err := prometheus.NewConstMetric(desc, valueType, value, labels...)
	if err != nil {
		println(fmt.Sprintf("err: %v", err))
	}

	// TODO: 处理nil的情况
	return metric
}

func getHostname() string {
	var err error
	once.Do(func() {
		hostname, err = os.Hostname()
		if err != nil {
			//logger.Errorf("Unable to get hostname:%v", err)
		}
	})

	return hostname
}
