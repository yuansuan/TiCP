/*
 * // Copyright (C) 2018 LambdaCal Inc.
 *
 */

package collection

import (
	"fmt"
	"strings"
)

// JoinInt64Array concats a int64 array into a single string
func JoinInt64Array(arr []int64, delim string) string {
	return strings.Trim(strings.Join(strings.Fields(fmt.Sprint(arr)), delim), "[]")
}

// ExistInInt64 return true if str exist in array
func ExistInInt64(arr []int64, i int64) bool {
	for _, v := range arr {
		if v == i {
			return true
		}
	}
	return false
}

// UniqueInt64Array unique int64 slice, not sort
func UniqueInt64Array(r []int64) []int64 {
	mp := map[int64]struct{}{}
	for _, ele := range r {
		mp[ele] = struct{}{}
	}

	result := make([]int64, 0, len(r))
	for id := range mp {
		result = append(result, id)
	}
	return result
}
