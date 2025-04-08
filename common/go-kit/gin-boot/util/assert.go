// Copyright (C) 2019 LambdaCal Inc.

package util

import (
	"log"
)

// Assert b, if not, panic and log
func Assert(b bool, v ...interface{}) {
	if !b {
		log.Panicf("PANIC! reason: %v", v)
	}
}

// Assertf b, if not, panic and log
func Assertf(b bool, format string, v ...interface{}) {
	if !b {
		log.Panicf("PANIC! reason: %v", v)
	}
}
