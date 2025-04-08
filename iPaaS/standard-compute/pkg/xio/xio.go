package xio

import (
	"context"
	"fmt"
	"io"
)

// SeekerLen 返回数据的大小
func SeekerLen(r io.ReadSeeker) (int64, error) {
	curr, err := r.Seek(0, io.SeekCurrent)
	if err != nil {
		return 0, err
	}

	sz, err := r.Seek(0, io.SeekEnd)
	if err != nil {
		return 0, err
	}

	if _, err = r.Seek(curr, io.SeekStart); err != nil {
		return 0, err
	}

	return sz, nil
}

// ctxReader 一个支持上下文对象的Reader
type ctxReader struct {
	io.Reader
	ctx context.Context
}

// Read 读取数据时检查上下文对象是否已经被取消
func (r *ctxReader) Read(p []byte) (n int, err error) {
	select {
	case <-r.ctx.Done():
		return 0, r.ctx.Err()
	default:
		return r.Reader.Read(p)
	}
}

// Copy 使用基于上下文的方式拷贝数据
func Copy(ctx context.Context, dst io.Writer, src io.Reader) (int64, error) {
	return io.Copy(dst, &ctxReader{Reader: src, ctx: ctx})
}

// ReadSeeker 合并Reader和Seeker
type ReadSeeker struct {
	io.Reader
	io.Seeker
}

type LogReader struct {
	io.Reader
}

func (r *LogReader) Read(p []byte) (n int, err error) {
	n, err = r.Reader.Read(p)
	fmt.Printf("read %d bytes (err = %v)\n", n, err)
	return n, err
}
