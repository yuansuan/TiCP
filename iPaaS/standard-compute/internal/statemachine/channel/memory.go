package channel

import "context"

type MemoryChannel struct {
	ch chan int64
}

func NewMemoryChannel() *MemoryChannel {
	return &MemoryChannel{
		ch: make(chan int64, 1000),
	}
}

func (c *MemoryChannel) SendMessage(ctx context.Context, jobID int64) (err error) {
	c.ch <- jobID
	return nil
}

func (c *MemoryChannel) RecvMessage(ctx context.Context) <-chan int64 {
	return c.ch
}

func (c *MemoryChannel) Close() error {
	close(c.ch)
	return nil
}
