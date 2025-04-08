package module

import "time"

const (
	DefaultRetryTimes    = 3
	DefaultRetryInterval = 1 * time.Second
	DefaultTimeout       = 5 * time.Second
)
