package v20230530

// MonitorChart 监控图表
type MonitorChart struct {
	Key   string
	Items []*MonitorChartItem
}

// MonitorChartItem 监控图表项
type MonitorChartItem struct {
	Kv []float64
}
