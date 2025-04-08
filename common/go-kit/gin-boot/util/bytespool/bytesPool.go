package bytespool

import (
	"sync"
)

// Bytes is wrap for bytes.Buffer
type Bytes struct {
	bs []byte
}

// Bytes return bytes of bytes
func (p *Bytes) Bytes() []byte {
	return p.bs
}

// BytesPool is pool of bytes.Buffer
type BytesPool struct {
	pool sync.Pool
}

// NewBytesPool return a new bytes pool
func NewBytesPool(bytesLen int64) *BytesPool {
	return &BytesPool{
		pool: sync.Pool{
			New: func() interface{} {
				return &Bytes{
					make([]byte, bytesLen),
				}
			},
		},
	}
}

// Get : will return a len=0 bytes.Buffer
func (p *BytesPool) Get() *Bytes {
	return p.pool.Get().(*Bytes)
}

// Put put back to pool
func (p *BytesPool) Put(bs *Bytes) {
	p.pool.Put(bs)
}
