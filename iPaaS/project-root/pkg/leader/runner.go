package leader

import (
	"context"
	"time"

	"github.com/yuansuan/ticp/common/go-kit/logging"
	"gopkg.in/redsync.v1"
)

// RunnerFunc ...
type RunnerFunc func(ctx context.Context)

// Runner ...
func Runner(ctx context.Context, key string, f RunnerFunc, options ...Option) {
	option := defaultOption()

	for _, v := range options {
		v.Apply(option)
	}

	rs := redsync.New(newPools(option.redis))

	mut := rs.NewMutex("leader_runner:"+key, redsync.SetExpiry(option.leaderKeyExpire))

	var isLeader = false
	preRun := func() bool {
		// is leader
		if isLeader {
			isLeader = mut.Extend()
			return isLeader
		}

		// not a leader
		if err := mut.Lock(); err != nil {
			if err != redsync.ErrFailed {
				logging.GetLogger(ctx).Warnf("leader prerun: %v", err)
			}
			return false
		}
		isLeader = true
		return true
	}

	go func() {
		ticker := time.NewTicker(option.interval)
		defer ticker.Stop()

		for {

			select {
			case <-ticker.C:
				if preRun() {
					logging.GetLogger(ctx).Infof("leader runner %v start", key)
					f(ctx)
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}
