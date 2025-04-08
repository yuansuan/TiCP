package utils

import (
	"time"
)

func PNow() *time.Time {
	now := time.Now()
	return &now
}

func PTime(t time.Time) *time.Time {
	return &t
}

var emptyTime = time.Time{}

func EmptyTime(t time.Time) bool {
	return t == emptyTime
}
