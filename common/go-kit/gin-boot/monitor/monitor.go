package monitor

import (
	"github.com/pkg/errors"
)

var (
	globalGaugeVec   = &GaugeVec{}
	globalSummaryVec = &SummaryVec{}
	globalCounterVec = &CounterVec{}
)

// Label Label
type Label struct {
	Name  string
	Value string
}

// AddCounter AddCounter
func AddCounter(name string, val float64, labels []*Label) error {
	if name == "" {
		return errors.New("monitor Add name is empty")
	}

	if labels == nil {
		labels = []*Label{}
	}

	if err := globalCounterVec.add(name, val, labels); err != nil {
		return err
	}

	return nil
}

// Add Add
func Add(name string, val float64, labels []*Label) error {
	if name == "" {
		return errors.New("monitor Add name is empty")
	}

	if labels == nil {
		labels = []*Label{}
	}

	if err := globalGaugeVec.add(name, val, labels); err != nil {
		return err
	}

	return nil
}

// Set Set
func Set(name string, val float64, labels []*Label) error {
	if name == "" {
		return errors.New("monitor Set key is empty")
	}

	if labels == nil {
		labels = []*Label{}
	}

	if err := globalGaugeVec.set(name, val, labels); err != nil {
		return err
	}

	return nil
}

// Reset Reset
func Reset(name string) error {
	if name == "" {
		return errors.New("monitor Set key is empty")
	}

	if err := globalGaugeVec.reset(name); err != nil {
		return err
	}

	return nil
}

// Observe Observe
func Observe(name string, val float64, objectives Objectives, labels []*Label) error {
	if name == "" {
		return errors.New("monitor Observe name is empty")
	}

	if labels == nil {
		labels = []*Label{}
	}

	if err := globalSummaryVec.observe(name, val, objectives, labels); err != nil {
		return err
	}

	return nil
}

func getLabelNameAndValue(labels []*Label) (names, values []string) {
	names = []string{}
	values = []string{}
	for _, label := range labels {
		names = append(names, label.Name)
		values = append(values, label.Value)
	}

	return
}
