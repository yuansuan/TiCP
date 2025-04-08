package xtime

import (
	"strconv"
	"time"
)

// CurrentTimestamp returns the timestamp of the current timezone
// of the system as a string
func CurrentTimestamp() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}

// CurrentMilliTimestamp returns the millisecond timestamp of
// the current timezone of the system as a string
func CurrentMilliTimestamp() string {
	return strconv.FormatInt(time.Now().UnixMilli(), 10)
}
