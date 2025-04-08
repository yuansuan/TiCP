package shutdown

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

type listener struct {
	shutdownCh chan os.Signal
}

func NewListener() *listener {
	l := &listener{
		shutdownCh: make(chan os.Signal, 1),
	}

	signal.Notify(l.shutdownCh,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGKILL)
	return l
}

func (l *listener) WaitWithCancel(cancel context.CancelFunc) {
	<-l.shutdownCh

	cancel()
}
