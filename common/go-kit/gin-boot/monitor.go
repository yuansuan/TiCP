package boot

import _monitor "github.com/yuansuan/ticp/common/go-kit/gin-boot/monitor"

type monitorType struct {
}

var (
	Monitor = &monitorType{}
)

// AddCounter AddCounter
func (m *monitorType) AddCounter(name string, val float64, labels []*_monitor.Label) error {
	return _monitor.AddCounter(name, val, labels)
}

// Add Add
func (m *monitorType) Add(name string, val float64, labels []*_monitor.Label) error {
	return _monitor.Add(name, val, labels)
}

// Set Set
func (m *monitorType) Set(name string, val float64, labels []*_monitor.Label) error {
	return _monitor.Set(name, val, labels)
}

// Reset Reset
func (m *monitorType) Reset(name string) error {
	return _monitor.Reset(name)
}

// Observe Observe
func (m *monitorType) Observe(name string, val float64, objectives _monitor.Objectives, labels []*_monitor.Label) error {
	return _monitor.Observe(name, val, objectives, labels)
}
