// Copyright (C) 2019 LambdaCal Inc.

package collection

import "reflect"

// Contain Contain
func Contain(items interface{}, item interface{}) bool {
	targetValue := reflect.ValueOf(items)

	switch reflect.TypeOf(items).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == item {
				return true
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(item)).IsValid() {
			return true
		}
	}

	return false
}

// ContainInt64 ContainInt64
func ContainInt64(items []int64, item int64) bool {
	for _, v := range items {
		if item == v {
			return true
		}
	}
	return false
}

// RemoveInt64 RemoveInt64
func RemoveInt64(items []int64, item int64) []int64 {
	result := make([]int64, 0, len(items))
	for _, v := range items {
		if v != item {
			result = append(result, v)
		}
	}
	return result
}
