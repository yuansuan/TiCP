// Copyright (C) 2019 LambdaCal Inc.

package util

// MaxInt64 return the max(a, b)
func MaxInt64(a, b int64) (c int64) {
	c = a
	if c < b {
		c = b
	}
	return
}

// Abs returns abs int value
func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
