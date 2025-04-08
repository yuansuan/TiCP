package monitor

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

// GaugeVec GaugeVec
type GaugeVec struct {
	sync.Mutex
	gaugeMap sync.Map
}

func (self *GaugeVec) add(name string, val float64, labels []*Label) error {
	labelNames, labelValues := getLabelNameAndValue(labels)

	prometheusGauge, err := self.get(name, labelNames)
	if err != nil {
		return err
	}
	prometheusGauge.WithLabelValues(labelValues...).Add(val)

	return nil
}

func (self *GaugeVec) set(name string, val float64, labels []*Label) error {
	labelNames, labelValues := getLabelNameAndValue(labels)

	gaugeObject, err := self.get(name, labelNames)
	if err != nil {
		return err
	}
	gaugeObject.WithLabelValues(labelValues...).Set(val)
	return nil
}

func (self *GaugeVec) reset(name string) error {
	if v, ok := self.gaugeMap.Load(name); ok {
		gaugeObject := v.(*prometheus.GaugeVec)
		gaugeObject.Reset()
	}
	return nil
}

func (self *GaugeVec) get(name string, labelNames []string) (*prometheus.GaugeVec, error) {
	if v, ok := self.gaugeMap.Load(name); ok {
		return v.(*prometheus.GaugeVec), nil
	}

	self.Lock()
	defer self.Unlock()
	if v, ok := self.gaugeMap.Load(name); ok {
		return v.(*prometheus.GaugeVec), nil
	}
	newGaugeVec, err := self.create(name, labelNames)
	if err != nil {
		return nil, err
	}

	return newGaugeVec, nil
}

func (self *GaugeVec) create(name string, labelNames []string) (*prometheus.GaugeVec, error) {
	newGaugeVec := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: name,
		},
		labelNames,
	)

	if err := prometheusRegister.Register(newGaugeVec); err != nil {
		return nil, err
	}

	self.gaugeMap.Store(name, newGaugeVec)

	return newGaugeVec, nil
}
