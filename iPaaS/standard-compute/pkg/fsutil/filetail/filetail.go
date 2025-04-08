package filetail

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
	"go.uber.org/multierr"

	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/log"
)

// Tail 从文件尾部读取数据
type Tail struct {
	f *os.File
	w *fsnotify.Watcher

	filesize int64

	filename  string
	follow    bool
	bytes     int64
	chunkSize int64
}

// Stream 从文件中获取数据
func (t *Tail) Stream(ctx context.Context) (<-chan []byte, error) {
	if _, err := t.f.Seek(t.bytes, io.SeekStart); err != nil {
		return nil, err
	}

	ch := make(chan []byte, 32)
	go t.followStream(ctx, ch)

	return ch, nil
}

// followStream 监听文件并在发生变化时读取数据
func (t *Tail) followStream(ctx context.Context, ch chan<- []byte) {
	var i int
	var firstEOF bool
	buf := make([]byte, t.chunkSize)
	for {
		n, err := io.ReadFull(t.f, buf[i:])
		if err != nil && err != io.ErrUnexpectedEOF && err != io.EOF {
			log.Warnf("reading from file failed: %s", err)
			break
		}
		isEOF := err == io.ErrUnexpectedEOF || err == io.EOF

		i += n
		log.Infof("reads data len = %d (curr = %d)", i, n)
		if int64(i) == t.chunkSize || !t.follow || !firstEOF {
			log.Info("writes data into channel (0)")
			writeBufferToChannel(ch, buf[:i])

			i = 0
			if isEOF && !t.follow {
				break
			}
		}

		firstEOF = firstEOF || isEOF
		if !isEOF {
			i = 0
			continue
		}

		ok, err := t.waiting(ctx, 3*time.Second)
		if err != nil {
			log.Infof("failed to waiting data: %s", err)
			break
		}

		if ok && i != 0 {
			log.Info("writes data into channel (1)")
			writeBufferToChannel(ch, buf[:i])
			i = 0
			continue
		}
	}

	close(ch)
}

// waiting 等待新的写入事件
func (t *Tail) waiting(ctx context.Context, timeout time.Duration) (bool, error) {
	select {
	case ev := <-t.w.Events:
		return ev.Op&fsnotify.Write == fsnotify.Write, nil
	case err := <-t.w.Errors:
		return false, err
	case <-time.After(timeout):
		return true, nil
	case <-ctx.Done():
		return false, ctx.Err()
	}
}

// Close 关闭文件以及监听程序
func (t *Tail) Close() (err error) {
	if t.w != nil {
		err = multierr.Append(err, t.w.Close())
	}
	return multierr.Append(err, t.f.Close())
}

// TailOption 用于指定额外的配置参数
type TailOption func(t *Tail) error

// WithFollow 为文件创建监听器
func WithFollow(follow bool) TailOption {
	return func(t *Tail) (err error) {
		if t.follow = follow; t.follow {
			if t.w, err = fsnotify.NewWatcher(); err == nil {
				err = t.w.Add(t.filename)
			}
		}
		return
	}
}

// WithBytes 配置需要预览的字节数
func WithBytes(bytes int64) TailOption {
	return func(t *Tail) error {
		if bytes > t.filesize || bytes < 0 {
			t.bytes = t.filesize
		}
		t.bytes = bytes
		return nil
	}
}

// WithChunkSize 每一块的大小
func WithChunkSize(sz int64) TailOption {
	return func(t *Tail) error {
		if sz > 32*1024 { // 32K
			sz = 32 * 1024
		}
		t.chunkSize = sz
		return nil
	}
}

// New 从一个文件名创建Tail
func New(filename string, options ...TailOption) (*Tail, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}

	t := &Tail{f: f, filename: filename, filesize: stat.Size()}
	for _, option := range options {
		if err = option(t); err != nil {
			return nil, err
		}
	}

	return t, nil
}

// writeBufferToChannel 将字节切片写入通道中
func writeBufferToChannel(ch chan<- []byte, buf []byte) {
	nbf := make([]byte, len(buf))
	for i, b := range buf {
		nbf[i] = b
	}
	ch <- nbf
}
