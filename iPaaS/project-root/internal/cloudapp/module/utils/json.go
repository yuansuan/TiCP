package utils

import "encoding/json"

func MustMarshalJson(v interface{}) string {
	bs, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return string(bs)
}
