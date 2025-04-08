package monitor

import (
	"fmt"
	"strconv"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/middleware"
	"github.com/yuansuan/ticp/common/go-kit/logging"
)

var usage = `
prometheus usage:
	prom.Use(cfg *prom.Config)
`

var prometheusRegister = prometheus.DefaultRegisterer
var internalMonitor *Monitor
var defaultMetricPath = "/metrics"

var (
	requestDurationHistogram = &Metric{
		Name: "request_durations_histogram_seconds",
		Help: "request latency distribution",
		Type: "histogram_vec",
		Args: []string{"uri", "method", "code"},
	}

	requestTotal = &Metric{
		Name: "request_total",
		Help: "request total",
		Type: "counter_vec",
		Args: []string{"uri", "method", "code"},
	}
)

var standardMetrics = []*Metric{
	requestDurationHistogram,
	requestTotal,
}

// Config Config
type Config struct {
	Server     *gin.Engine
	ListenAddr string
	MetricPath string
}

// Metric Metric
type Metric struct {
	MetricCollector prometheus.Collector
	Name            string
	Help            string
	Type            string
	Args            []string
}

// Monitor Monitor
type Monitor struct {
	reqCounter       *prometheus.CounterVec
	requestHistogram *prometheus.HistogramVec
	listenAddress    string
	MetricsList      []*Metric
	MetricsPath      string
	Server           *gin.Engine
}

func newMonitor(cfg *Config) *Monitor {

	metricPath := defaultMetricPath
	if cfg.MetricPath != "" {
		metricPath = cfg.MetricPath
	}

	internalMonitor = &Monitor{
		MetricsList:   standardMetrics,
		MetricsPath:   metricPath,
		listenAddress: cfg.ListenAddr,
		Server:        cfg.Server,
	}

	internalMonitor.registerMetrics()

	return internalMonitor
}

func (p *Monitor) registerMetrics() {

	for _, metricDef := range p.MetricsList {
		metric := newMetric(metricDef, "yuansuan")
		if err := prometheusRegister.Register(metric); err != nil {
			panic(fmt.Sprintf("promtheus error due to %v\n", err))
		}

		switch metricDef {
		case requestTotal:
			p.reqCounter = metric.(*prometheus.CounterVec)
		case requestDurationHistogram:
			p.requestHistogram = metric.(*prometheus.HistogramVec)
		}
		metricDef.MetricCollector = metric
	}

}

// Use Use
func Use(cfg *Config) {
	p := newMonitor(cfg)                   // new Monitor with cfg
	cfg.Server.Use(p.collectHandlerFunc()) // Use collect middleware
	p.runMetricsServer()                   // run metric http server
}

func (p *Monitor) runMetricsServer() {
	if p.listenAddress != "" {
		router := gin.New()
		router.Use(middleware.GinLogger("/debug/pprof/cmdline", "/metrics"), logging.GinRecovery())
		router.GET(p.MetricsPath, prometheusHandler())
		go router.Run(p.listenAddress)
	} else {
		p.Server.GET(p.MetricsPath, prometheusHandler())
	}
}

func prometheusHandler() gin.HandlerFunc {
	return gin.WrapH(promhttp.Handler())
}

func (p *Monitor) collectHandlerFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.String() == p.MetricsPath {
			c.Next()
			return
		}

		requestUrl := c.Request.URL.Path
		requestMethod := c.Request.Method
		start := time.Now()

		c.Next()

		status := strconv.Itoa(c.Writer.Status())
		elapsed := float64(time.Since(start)) / float64(time.Second)

		labels := []string{requestUrl, requestMethod, status}
		for _, val := range labels {
			if !utf8.ValidString(val) {
				logging.Default().Warnf("label value %q is not valid UTF-8", val)
				return
			}
		}

		p.requestHistogram.WithLabelValues(labels...).Observe(elapsed)
		p.reqCounter.WithLabelValues(labels...).Inc()
	}
}

func newMetric(m *Metric, subsystem string) prometheus.Collector {
	var metric prometheus.Collector
	switch m.Type {
	case "counter_vec":
		metric = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Subsystem: subsystem,
				Name:      m.Name,
				Help:      m.Help,
			},
			m.Args,
		)
	case "counter":
		metric = prometheus.NewCounter(
			prometheus.CounterOpts{
				Subsystem: subsystem,
				Name:      m.Name,
				Help:      m.Help,
			},
		)
	case "gauge_vec":
		metric = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Subsystem: subsystem,
				Name:      m.Name,
				Help:      m.Help,
			},
			m.Args,
		)
	case "gauge":
		metric = prometheus.NewGauge(
			prometheus.GaugeOpts{
				Subsystem: subsystem,
				Name:      m.Name,
				Help:      m.Help,
			},
		)
	case "histogram_vec":
		metric = prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Subsystem: subsystem,
				Name:      m.Name,
				Help:      m.Help,
			},
			m.Args,
		)
	case "histogram":
		metric = prometheus.NewHistogram(
			prometheus.HistogramOpts{
				Subsystem: subsystem,
				Name:      m.Name,
				Help:      m.Help,
			},
		)
	case "summary_vec":
		metric = prometheus.NewSummaryVec(
			prometheus.SummaryOpts{
				Subsystem: subsystem,
				Name:      m.Name,
				Help:      m.Help,
			},
			m.Args,
		)
	case "summary":
		metric = prometheus.NewSummary(
			prometheus.SummaryOpts{
				Subsystem: subsystem,
				Name:      m.Name,
				Help:      m.Help,
			},
		)
	}
	return metric
}
