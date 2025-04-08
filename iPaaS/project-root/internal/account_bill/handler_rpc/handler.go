package handler_rpc

import (
	"context"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/http"
	"sync"
)

var (
	quit bool
	lock sync.RWMutex
)

var (
	ctx    context.Context
	cancel context.CancelFunc
)

func init() {
	ctx, cancel = context.WithCancel(context.TODO())
}

// OnShutdown OnShutdown
func OnShutdown(drv *http.Driver) {
	lock.Lock()
	defer lock.Unlock()
	quit = true

	cancel()
}

// InitGRPCClient InitGRPCClient
func InitGRPCClient(drv *http.Driver) {
	//payment.SchedulerStart(ctx)
}
