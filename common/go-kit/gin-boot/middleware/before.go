package middleware

import (
	"context"
	"sync"

	"google.golang.org/grpc"
)

// BeforeFunc is func of before
type BeforeFunc func(ctx context.Context) error

var before = map[string]BeforeFunc{}
var beforeLock sync.RWMutex

// SetBefore set before func for method,
// methodName is like: /rbac.RoleManager/AddRole
func SetBefore(methodName string, f BeforeFunc) {
	beforeLock.Lock()
	defer beforeLock.Unlock()

	before[methodName] = f
}

// BeforeInterceptor BeforeInterceptor
func BeforeInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	beforeLock.RLock()
	defer beforeLock.RUnlock()

	if f := before[info.FullMethod]; f != nil {
		err := f(ctx)
		if err != nil {
			return nil, err
		}
	}

	return handler(ctx, req)
}
