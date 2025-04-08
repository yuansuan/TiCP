package locker

import (
	"context"
	"strings"
	"time"

	"github.com/yuansuan/ticp/iPaaS/standard-compute/pkg/xoss"
)

// S3Locker 一个基于S3轮询的锁
type S3Locker struct {
	oss  *xoss.ObjectStorageService
	name string
}

// Lock 通过创建文件进行上锁
func (l *S3Locker) Lock(timeout time.Duration) (bool, error) {
	ctx, cancel := withContext(context.Background(), timeout)
	defer cancel()

	_, err := l.oss.HeadObject(ctx, l.name)
	if err != nil && xoss.IsObjectNotExists(err) {
		_, err = l.oss.PutObject(ctx, l.name, strings.NewReader(l.name))
		return err == nil, err
	}
	return false, err
}

// Unlock 删除文件进行解锁
func (l *S3Locker) Unlock() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := l.oss.HeadObject(ctx, l.name)
	if err == nil {
		_, err = l.oss.DeleteObject(ctx, l.name)
		return err
	}
	return err
}

// NewS3Locker 创建一个基于S3轮询的锁
func NewS3Locker(oss *xoss.ObjectStorageService, name string) *S3Locker {
	return &S3Locker{oss: oss, name: name}
}

// withContext 创建上下文对象
func withContext(parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	if timeout <= 0 {
		return context.WithCancel(parent)
	}
	return context.WithTimeout(parent, timeout)
}
