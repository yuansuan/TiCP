package channel

import (
	"context"

	"github.com/yuansuan/ticp/iPaaS/standard-compute/config"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/log"
)

type Channel interface {
	SendMessage(context.Context, int64) error
	RecvMessage(context.Context) <-chan int64

	Close() error
}

// NewChannel 创建一个通道
func NewChannel(cfg *config.StateMachine) (ch Channel) {
	switch cfg.Channel {
	case "memory":
		return NewMemoryChannel()
	default:
		log.Infof("unknown channel type %q, using default memory channel", cfg.Channel)
		return NewMemoryChannel()
	}
}
