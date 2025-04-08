package util

import "encoding/json"

func MustParseJSON(v interface{}) (s string) {
	if bs, err := json.Marshal(v); err == nil {
		s = string(bs)
	}
	return
}
