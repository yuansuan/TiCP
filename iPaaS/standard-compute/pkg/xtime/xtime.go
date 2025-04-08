package xtime

import (
	"context"
	"time"
)

// MaxDuration 返回给定时间最长的那个
func MaxDuration(a, b time.Duration) time.Duration {
	if a > b {
		return a
	}
	return b
}

// Sleep 睡眠一段时间同时检查 ctx 是否被取消
func Sleep(ctx context.Context, duration time.Duration) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(duration):
		return nil
	}
}
