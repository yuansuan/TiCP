package common

import (
	jsoniter "github.com/json-iterator/go"
)

func MustString(v interface{}) string {
	if v == nil {
		return ""
	}

	s, _ := jsoniter.MarshalToString(v)
	return s
}
