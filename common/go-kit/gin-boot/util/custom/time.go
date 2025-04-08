// Copyright (C) 2019 LambdaCal Inc.

package custom

import (
	"fmt"
	"time"
)

// JSONTime JSONTime
type JSONTime struct {
	time.Time
}

// MarshalJSON MarshalJSON
func (t JSONTime) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("%d", t.Unix())
	return []byte(stamp), nil
}

// YAMLDuration is duration type for yaml
type YAMLDuration struct {
	Duration time.Duration
}

// UnmarshalYAML unmarshal duration in yaml
// duration is golang format duration, like 1s, 2m, 3h, 3h2m1s
// Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
func (p *YAMLDuration) UnmarshalYAML(unmarshal func(interface{}) error) error {
	durationString := ""
	err := unmarshal(&durationString)
	if err != nil {
		return err
	}
	dur, err := time.ParseDuration(durationString)
	if err != nil {
		return err
	}

	p.Duration = dur
	return nil
}
