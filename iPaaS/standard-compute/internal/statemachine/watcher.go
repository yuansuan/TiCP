package statemachine

import (
	"context"
	"fmt"

	"github.com/yuansuan/ticp/common/project-root-api/hpc/jobstate"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/backend/job"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/pkg/xtime"
)

// watcherState 监听器的状态
type watcherState int

const (
	// stopWatcher 停止监听
	stopWatcher watcherState = iota // 0
	holdWatcher                     // 1

	checkAliveToleranceErrorTimes = 20
)

// stateWatcher 根据当前作业的状态决定是否需要继续轮询
type stateWatcher func(j *job.Job) (watcherState, error)

// watchState 监听任务状态知道发生错误或者返回 stopWatcher 表示停止监听
func (m *StateMachine) watchState(ctx context.Context, j *job.Job, watcher stateWatcher) error {
	checkAliveErrorTimes := 0

	for {
		js := j.Job.State
		// 每次查询都会休眠固定的时间
		if err := xtime.Sleep(ctx, m.watcherInterval); err != nil {
			return err
		}

		// 从数据库中获取最新的作业信息, 用于检查是否被用户取消等问题
		nj, err := m.dao.GetJobWithError(ctx, j.Id)
		if err != nil {
			return err
		}
		j.Job = nj // 直接替换模型
		if j.Job.State != js && j.Job.State != jobstate.Failed {
			// 可能是由于状态转换时数据库更新失败
			j.Job.State = js
		}

		// 如果作业被取消则直接返回错误，让外层处理关闭
		if j.ControlBitTerminate {
			return ErrJobCanceled
		}

		// 通过调度器检查作业的状态
		j, err = m.backend.CheckAlive(ctx, j)
		if err != nil {
			j.TraceLogger.Warnf("check alive failed, %v", err)
			// 超算环境有小概率出现执行check alive失败的异常情况，这里加上容错，连续失败一定次数才认为check alive失败，且打error日志告警
			checkAliveErrorTimes++
			if checkAliveErrorTimes < checkAliveToleranceErrorTimes {
				continue
			}

			err = fmt.Errorf("check alive failed more than %d times continuously, should check it manually! err detail: %w", checkAliveToleranceErrorTimes, err)
			j.TraceLogger.Error(err)
			return err
		}
		checkAliveErrorTimes = 0

		// 更新并推送数据
		// @TODO 这里如果可以做个diff最好，没有变更就不需要推送了
		if err = m.UpdateAndPushJob(ctx, j); err != nil {
			return err
		}

		if j.ExecutionDuration < 0 {
			// 可能会由于时间同步问题导致执行时间为负数，会使上面的UpdateAndPushJob更新数据库失败
			// 为不影响后续作业状态更新，这里做一下修正
			j.ExecutionDuration = 0
		}

		if state, wErr := watcher(j); wErr != nil {
			// 在发生错误时，可能会修改作业的状态，这里需要再推送一次
			// 这里直接忽略推送返回的需要终止错误
			_ = m.UpdateAndPushJob(ctx, j)
			return wErr
		} else if state == stopWatcher {
			return nil
		}
	}
}
