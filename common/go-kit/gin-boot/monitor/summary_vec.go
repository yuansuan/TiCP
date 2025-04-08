package monitor

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

// SummaryVec SummaryVec
type SummaryVec struct {
	sync.Mutex
	summaryMap sync.Map
}

// Objectives Objectives
type Objectives map[float64]float64

func (self *SummaryVec) observe(name string, val float64, objectives Objectives, labels []*Label) error {
	labelNames, labelValues := getLabelNameAndValue(labels)

	s, err := self.get(name, objectives, labelNames)
	if err != nil {
		return err
	}
	s.WithLabelValues(labelValues...).Observe(val)

	return nil
}

func (self *SummaryVec) get(name string, objectives Objectives, labelNames []string) (*prometheus.SummaryVec, error) {

	if v, ok := self.summaryMap.Load(name); ok {
		return v.(*prometheus.SummaryVec), nil
	}

	self.Lock()
	defer self.Unlock()
	if v, ok := self.summaryMap.Load(name); ok {
		return v.(*prometheus.SummaryVec), nil
	}
	s, err := self.create(name, objectives, labelNames)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (self *SummaryVec) create(name string, objectives Objectives, labelNames []string) (*prometheus.SummaryVec, error) {
	newSummaryVec := prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       name,
			Objectives: objectives,
		},
		labelNames,
	)

	if err := prometheusRegister.Register(newSummaryVec); err != nil {
		return nil, err
	}
	self.summaryMap.Store(name, newSummaryVec)

	return newSummaryVec, nil
}
