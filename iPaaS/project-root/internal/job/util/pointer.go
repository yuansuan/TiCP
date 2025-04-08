package util

import (
	"time"
)

func PString(s string) *string {
	return &s
}

func PFloat64(f float64) *float64 {
	return &f
}

func PBool(b bool) *bool {
	return &b
}

func PTime(t time.Time) *time.Time {
	return &t
}
