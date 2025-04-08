package utils

import "github.com/yuansuan/ticp/common/go-kit/gin-boot/monitor"

// Metrics 云服务指标
type Metrics struct {
	name string
}

// CallAPI 调用API计数
func (m *Metrics) CallAPI(api string) {
	_ = monitor.AddCounter("cloud_api_calls", 1, []*monitor.Label{
		{
			Name:  "platform",
			Value: m.name,
		},
		{
			Name:  "api",
			Value: api,
		},
	})
}

// NewMetrics 创建一个信的云服务指标
func NewMetrics(name string) *Metrics {
	return &Metrics{name: name}
}
