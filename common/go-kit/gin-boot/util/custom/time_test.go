// Copyright (C) 2019 LambdaCal Inc.

package custom

import (
	"encoding/json"
	"testing"
	"time"
)

type testTime struct {
	Tm  time.Time `json:"tm"`
	Ctm JSONTime  `json:"ctm"`
}

func TestJSONTime_MarshalJSON(t *testing.T) {
	const sec = 1546843884
	a := testTime{
		Tm:  time.Unix(sec, 0),
		Ctm: JSONTime{time.Unix(sec, 0)},
	}
	cont, err := json.Marshal(a)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(cont))
	if string(cont) != `{"tm":"2019-01-07T14:51:24+08:00","ctm":1546843884}` {
	}
}
