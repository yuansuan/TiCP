package leader

import (
	"time"
)

// RegisterOption ...
type RegisterOption struct {
	leaderKeyExpire time.Duration
	interval        time.Duration
	redis           string
}

func defaultOption() *RegisterOption {
	return &RegisterOption{
		interval:        1 * time.Minute,
		leaderKeyExpire: 2 * time.Minute,
		redis:           "default",
	}
}

// An Option configures a Register.
type Option interface {
	Apply(*RegisterOption)
}

// OptionFunc is a function that configures a Register.
type OptionFunc func(*RegisterOption)

// Apply calls f(mutex)
func (f OptionFunc) Apply(mutex *RegisterOption) {
	f(mutex)
}

// SetLeaderKeyExpire set the expiry of a leader key to the given value.
func SetLeaderKeyExpire(expire time.Duration) Option {
	return OptionFunc(func(o *RegisterOption) {
		o.leaderKeyExpire = expire
	})
}

// SetInterval set the interval of func to the given value.
func SetInterval(interval time.Duration) Option {
	return OptionFunc(func(o *RegisterOption) {
		o.interval = interval
	})
}

// SetRedis set the redis name
func SetRedis(name string) Option {
	return OptionFunc(func(o *RegisterOption) {
		o.redis = name
	})
}
