package xsignal

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// WithExit 监听退出相关的信号
func WithExit(parent context.Context) (context.Context, context.CancelFunc) {
	return With(parent, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
}

// With 基于系统信号创建相对应的上下文对象
func With(parent context.Context, signals ...os.Signal) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(parent)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, signals...)
	go func() {
		defer func() {
			signal.Stop(sigs)
			close(sigs)
			cancel()
		}()

		for {
			select {
			case <-parent.Done():
				return
			case <-sigs:
				return
			}
		}
	}()

	return ctx, cancel
}
