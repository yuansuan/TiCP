package leader

import (
	"gopkg.in/redsync.v1"
)

// WithLock a global entry lock
func WithLock(key string, f func() error, options ...Option) error {
	option := defaultOption()
	for _, v := range options {
		v.Apply(option)
	}

	rs := redsync.New(newPools(option.redis))
	lock := rs.NewMutex("leader_lock:"+key, redsync.SetExpiry(option.leaderKeyExpire))
	if err := lock.Lock(); err != nil {
		return err
	}
	defer lock.Unlock()

	return f()
}
