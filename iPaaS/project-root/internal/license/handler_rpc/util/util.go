package util

import (
	"github.com/golang/protobuf/ptypes/timestamp"
	"time"
)

// ModelToProtoTime convert time.Time to proto timestamp
func ModelToProtoTime(time *time.Time) *timestamp.Timestamp {
	if time.IsZero() {
		return nil
	}
	return &timestamp.Timestamp{
		Seconds: time.Unix(),
		Nanos:   int32(time.Nanosecond()),
	}
}

// ProtoTimeToTime convert proto timestamp to time.Time
func ProtoTimeToTime(t *timestamp.Timestamp) time.Time {
	return time.Unix(t.Seconds, int64(t.Nanos))
}
