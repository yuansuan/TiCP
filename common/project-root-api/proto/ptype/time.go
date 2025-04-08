package ptype

import (
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
)

// ModelToProtoTime ...
func ModelToProtoTime(time *time.Time) *timestamp.Timestamp {
	if time == nil {
		return nil
	}

	return &timestamp.Timestamp{
		Seconds: time.Unix(),
		Nanos:   int32(time.Nanosecond()),
	}
}
