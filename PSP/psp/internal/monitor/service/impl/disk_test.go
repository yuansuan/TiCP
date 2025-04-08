package impl

import (
	"testing"
	"time"
)

func TestName(t *testing.T) {
	var start = time.Date(2023, 8, 16, 2, 10, 0, 0, time.UTC) // 获取时间戳（Unix 时间戳，秒）
	var end = time.Date(2023, 8, 16, 13, 10, 0, 0, time.UTC)  // 获取时间戳（Unix 时间戳，秒）
	t.Log(start.UnixMilli())
	t.Log(end.UnixMilli())
}
