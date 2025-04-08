package xtime

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCurrentTimestamp(t *testing.T) {
	if ts, err := strconv.ParseInt(CurrentTimestamp(), 10, 64); assert.NoError(t, err) {
		assert.LessOrEqual(t, ts, time.Now().Unix())
	}
}

func TestCurrentMilliTimestamp(t *testing.T) {
	if mts, err := strconv.ParseInt(CurrentMilliTimestamp(), 10, 64); assert.NoError(t, err) {
		assert.LessOrEqual(t, mts, time.Now().UnixMilli())
	}
}
