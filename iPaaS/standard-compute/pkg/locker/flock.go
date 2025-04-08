package locker

import (
	"os"
	"syscall"
	"time"

	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/log"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/pkg/xtime"
)

// FileLocker 是一个基于文件的互斥锁
type FileLocker struct {
	cleanup func()
	open    FileOpener
	f       *os.File
}

// Lock 尝试对文件进行上锁, 直到获取锁或者超时
func (l *FileLocker) Lock(timeout time.Duration) (_ bool, err error) {
	if l.f != nil {
		return false, ErrDuplicateLock
	}

	l.f, l.cleanup, err = l.open()
	if err != nil {
		return false, err
	}

	deadline := time.Now().Add(timeout)
	for err = l.acquire(); err != nil; err = l.acquire() {
		if err != nil {
			log.Warn(err)
		}

		if timeout == FastFail {
			return false, nil
		}
		if timeout != UntilAcquire && deadline.Before(time.Now()) {
			return false, ErrLockTimeout
		}
		time.Sleep(xtime.MaxDuration(timeout/_AcquireTries, _MinimalSleepTime))
	}

	return true, nil
}

// Unlock 释放文件上的锁
func (l *FileLocker) Unlock() error {
	if l.f != nil {
		defer func() {
			l.cleanup()
			l.f = nil
			l.cleanup = nil
		}()
		return l.release()
	}
	return nil
}

// acquire 通过调用系统API对文件进行上锁
func (l *FileLocker) acquire() error {
	// locking with exclusive and non-blocking
	return syscall.Flock(int(l.f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
}

// release 释放获取到的文件锁
func (l *FileLocker) release() error {
	return syscall.Flock(int(l.f.Fd()), syscall.LOCK_UN)
}

// FileOpener 打开一个文件并返回文件对象或者一个错误
type FileOpener func() (*os.File, func(), error)

// NewFileLocker 创建一个延迟打开的文件锁
func NewFileLocker(open FileOpener) *FileLocker {
	return &FileLocker{open: open}
}
