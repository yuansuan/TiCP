package xhttp

import (
	"net/http"
	"sync/atomic"
)

var (
	_Sequence int32
)

const (
	_DefaultMiddlewarePriority = 50
)

type MiddlewareHandler func(req *http.Request) (*http.Response, error)

type Middleware interface {
	Process(req *http.Request, next MiddlewareHandler) (*http.Response, error)
}

type MiddlewareFunc func(req *http.Request, next MiddlewareHandler) (*http.Response, error)

func (f MiddlewareFunc) Process(req *http.Request, next MiddlewareHandler) (*http.Response, error) {
	return f(req, next)
}

type withPriorityMiddleware struct {
	Middleware

	index    int
	priority int
}

func newPriorityMiddleware(middleware Middleware, priority int) *withPriorityMiddleware {
	return &withPriorityMiddleware{
		index:      int(atomic.AddInt32(&_Sequence, 1)),
		priority:   priority,
		Middleware: middleware,
	}
}

type Middlewares []*withPriorityMiddleware

func (m Middlewares) Len() int            { return len(m) }
func (m Middlewares) Swap(i, j int)       { m[i], m[j] = m[j], m[i] }
func (m *Middlewares) Pop() interface{}   { *m = (*m)[1:]; return (*m)[0] }
func (m *Middlewares) Push(x interface{}) { *m = append(*m, x.(*withPriorityMiddleware)) }
func (m Middlewares) Less(i, j int) bool {
	if m[i].priority == m[j].priority {
		return m[i].index < m[j].index
	}
	return m[i].priority < m[j].priority
}
