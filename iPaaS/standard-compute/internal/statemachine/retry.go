package statemachine

import (
	"time"

	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/log"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/util"
)

// withRetry 自动重试
func (m *StateMachine) withRetry(times int, action util.ReplayAction) error {
	return util.RetryWithBackoff(times, time.Second, time.Minute, action, util.WithLogger(log.GetLogger().Sugar()))
}

// withRetry 快速自动重试
func (m *StateMachine) withFastRetry(times int, action util.ReplayAction) error {
	return util.RetryWithBackoff(times, time.Second, 3*time.Second, action, util.WithLogger(log.GetLogger().Sugar()))
}
