package util

import (
	"time"

	"github.com/rfyiamcool/backoff"
	"go.uber.org/zap"

	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/log"
)

type ReplayState int

const (
	StopReplay ReplayState = iota // 0
	HoldReplay                    // 1
)

// Replay 用于执行重试的实用工具类
type Replay struct {
	logger  *zap.SugaredLogger
	backoff *backoff.Backoff
}

// Retry 重试直到成功或次数超过限制
func (r *Replay) Retry(times int, action ReplayAction) error {
	b := *r.backoff

	var lastErr error
	for i := 0; i < times; i++ {
		log := r.logger.Debugw
		if 3*i > times {
			log = r.logger.Warnw
		} else if 3*i > 2*times {
			log = r.logger.Errorw
		}

		var state ReplayState
		if state, lastErr = action(i, lastErr, log); lastErr == nil || state == StopReplay {
			break
		}

		time.Sleep(b.Duration())
	}

	return lastErr
}

// ReplayAction 每次重试所需要执行的操作
type ReplayAction func(curr int, lastErr error, log WideLogger) (ReplayState, error)

// WideLogger 是支持输出键值对的一个日志函数
type WideLogger func(string, ...interface{})

// NewReplay 创建一个重放工具
func NewReplay(options ...ReplayOption) *Replay {
	reply := &Replay{
		logger:  log.GetLogger().Sugar(),
		backoff: backoff.NewBackOff(),
	}

	for _, opt := range options {
		opt(reply)
	}

	return reply
}

// ReplayOption 自动重试的配置选项
type ReplayOption func(r *Replay)

// WithLogger 配置日志记录器
func WithLogger(logger *zap.SugaredLogger) ReplayOption {
	return func(r *Replay) {
		r.logger = logger
	}
}

// WithBackoff 配置为幂级回退模式
func WithBackoff(backoff *backoff.Backoff) ReplayOption {
	return func(r *Replay) {
		r.backoff = backoff
	}
}

// WithFixedInterval 使用固定时间形势的回退
func WithFixedInterval(interval time.Duration) ReplayOption {
	return func(r *Replay) {
		r.backoff = backoff.NewBackOff(
			backoff.WithMinDelay(interval),
			backoff.WithMaxDelay(interval),
		)
	}
}

// RetryWithBackoff 以幂级回退模式重试任务
func RetryWithBackoff(times int, min, max time.Duration, action ReplayAction, options ...ReplayOption) error {
	return NewReplay(append([]ReplayOption{
		WithBackoff(backoff.NewBackOff(
			backoff.WithMinDelay(min),
			backoff.WithMaxDelay(max),
		)),
	}, options...)...).Retry(times, action)
}

// AutoStopRetry 根据错误类型决定是否需要停止重试
func AutoStopRetry(err error) (ReplayState, error) {
	if IsCanceledError(err) || err == nil {
		return StopReplay, err
	}
	return HoldReplay, err
}

// SecondToDuration 将数字秒转换为一个时间间隔
func SecondToDuration(sec int) time.Duration {
	return time.Duration(sec) * time.Second
}
