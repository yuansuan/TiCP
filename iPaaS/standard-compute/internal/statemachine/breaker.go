package statemachine

import (
	"context"
	"sync"
)

// Breaker 管理任务的断路器
type Breaker struct {
	mu sync.Mutex
	m  map[int64]*breaker
}

// Break 关闭一个任务的断路器
func (b *Breaker) Break(id int64) (ok bool) {
	var brk *breaker

	b.mu.Lock()
	if brk, ok = b.m[id]; ok {
		brk.cancel()
		delete(b.m, id)
	}
	b.mu.Unlock()

	return
}

// Create 为任务创建一个断路器
func (b *Breaker) Create(ctx context.Context, id int64) context.Context {
	b.mu.Lock()
	defer b.mu.Unlock()

	if brk, ok := b.m[id]; ok {
		return brk.ctx
	}

	bCtx, cancel := context.WithCancel(ctx)
	brk := &breaker{ctx: bCtx, cancel: cancel}
	b.m[id] = brk

	return brk.ctx
}

// breaker 某一个具体任务的断路器
type breaker struct {
	ctx    context.Context
	cancel context.CancelFunc
}

// NewBreaker 创建任务断路器
func NewBreaker() (*Breaker, error) {
	return &Breaker{
		m: make(map[int64]*breaker),
	}, nil
}
