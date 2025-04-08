package locker

import (
	"errors"
	"time"
)

const (
	// UntilAcquire 阻塞直到获取到锁
	UntilAcquire time.Duration = -1

	// FastFail 快速失败
	FastFail time.Duration = 0

	// _MinimalSleepTime 指定最小的睡眠时间
	_MinimalSleepTime = 20 * time.Millisecond

	// _AcquireTries 在超时时间内获取锁的重试次数
	_AcquireTries = 25
)

var (
	// ErrLockTimeout 表示获取锁超时了
	ErrLockTimeout = errors.New("lock: timeout")
	// ErrDuplicateLock 表示重复调用了Lock方法
	ErrDuplicateLock = errors.New("lock: duplicated")
)

// Locker 是一个支持并发锁的接口
type Locker interface {
	// Lock 尝试进行上锁, 直到获取锁或者超时
	Lock(timeout time.Duration) (ok bool, err error)

	// Unlock 释放锁
	Unlock() error
}
