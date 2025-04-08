package registry

import (
	"io"
	"sync/atomic"
)

// AtomicProgress 是一个并发安全的进度控制
type AtomicProgress struct {
	io.Reader

	total   int64
	current int64
}

// Read 读取数据并增加当前已读取的字节数
func (ap *AtomicProgress) Read(p []byte) (n int, err error) {
	n, err = ap.Reader.Read(p)
	atomic.AddInt64(&ap.current, int64(n))
	return n, err
}

// TotalBytes 需要传输的总字节数
func (ap *AtomicProgress) TotalBytes() int64 {
	return ap.total
}

// TransferredBytes 当前已发送的字节数
func (ap *AtomicProgress) TransferredBytes() int64 {
	return atomic.LoadInt64(&ap.current)
}

// NewAtomicProgress 创建一个自定义的进度
func NewAtomicProgress(sz int64, r io.Reader) *AtomicProgress {
	return &AtomicProgress{Reader: r, total: sz}
}
