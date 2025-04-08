package monitor

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

// CounterVec CounterVec
type CounterVec struct {
	sync.Mutex
	counterMap sync.Map
}

func (c *CounterVec) add(name string, val float64, labels []*Label) error {
	labelNames, labelValues := getLabelNameAndValue(labels)

	prometheusCounter, err := c.get(name, labelNames)
	if err != nil {
		return err
	}
	prometheusCounter.WithLabelValues(labelValues...).Add(val)
	return nil
}

func (c *CounterVec) get(name string, labelNames []string) (*prometheus.CounterVec, error) {
	if v, ok := c.counterMap.Load(name); ok {
		return v.(*prometheus.CounterVec), nil
	}

	c.Lock()
	defer c.Unlock()
	if v, ok := c.counterMap.Load(name); ok {
		return v.(*prometheus.CounterVec), nil
	}

	newCounterVec, err := c.create(name, labelNames)
	if err != nil {
		return nil, err
	}

	return newCounterVec, nil
}

func (c *CounterVec) create(name string, labelNames []string) (*prometheus.CounterVec, error) {
	newCounterVec := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: name,
		},
		labelNames,
	)

	if err := prometheusRegister.Register(newCounterVec); err != nil {
		return nil, err
	}

	c.counterMap.Store(name, newCounterVec)

	return newCounterVec, nil
}
