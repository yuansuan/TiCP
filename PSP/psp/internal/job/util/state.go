package util

import (
	"fmt"
	"time"

	"github.com/yuansuan/ticp/PSP/psp/internal/common"
	"github.com/yuansuan/ticp/PSP/psp/internal/job/consts"
	"github.com/yuansuan/ticp/PSP/psp/internal/job/dao/model"
)

var (
	// finishedJobStateMap 已结束的状态列表
	finishedJobStateMap = make(map[string]struct{}, 4)
	// UnFinishedStates 未结束的状态列表
	UnFinishedStates = []string{consts.JobStateSubmitted, consts.JobStatePending, consts.JobStateRunning, consts.JobStateSuspended}
	// MonitorJobStates 监控作业状态列表
	MonitorJobStates = []string{consts.JobStatePending, consts.JobStateRunning, consts.JobStateCompleted, consts.JobStateFailed}
)

func init() {
	finishedJobStateMap[consts.JobStateFailed] = struct{}{}
	finishedJobStateMap[consts.JobStateCompleted] = struct{}{}
	finishedJobStateMap[consts.JobStateTerminated] = struct{}{}
	finishedJobStateMap[consts.JobStateBurstFailed] = struct{}{}
}

// ConvertJobState 作业状态转换
func ConvertJobState(jobState string) string {
	var state string

	switch jobState {
	case consts.APIJobStatePending:
		state = consts.JobStatePending
	case consts.APIJobStateRunning:
		state = consts.JobStateRunning
	case consts.APIJobStateTerminated:
		state = consts.JobStateTerminated
	case consts.APIJobStateSuspended, consts.APIJobStateInitiallySuspended:
		state = consts.JobStateSuspended
	case consts.APIJobStateCompleted:
		state = consts.JobStateCompleted
	case consts.APIJobStateFailed:
		state = consts.JobStateFailed
	}

	return state
}

// ConvertJobStateMsg 作业状态描述转换
func ConvertJobStateMsg(jobState string) string {
	var msg string

	switch jobState {
	case consts.JobStatePending:
		msg = "等待中"
	case consts.JobStateRunning:
		msg = "运行中"
	case consts.JobStateTerminated:
		msg = "已终止"
	case consts.JobStateSuspended:
		msg = "已暂停"
	case consts.JobStateCompleted:
		msg = "已完成"
	case consts.JobStateFailed:
		msg = "已失败"
	}

	return msg
}

// ConvertJobTimelineEvent 作业时间线事件名称转换
func ConvertJobTimelineEvent(job *model.Job) (string, time.Time) {
	var eventName string
	var eventTime time.Time

	switch job.State {
	case consts.JobStateSubmitted:
		eventTime = job.SubmitTime
	case consts.JobStatePending:
		eventTime = job.PendTime
	case consts.JobStateRunning:
		eventTime = job.StartTime
	case consts.JobStateTerminated:
		eventTime = job.TerminateTime
	case consts.JobStateSuspended:
		eventTime = job.SuspendTime
	case consts.JobStateCompleted:
		eventTime = job.EndTime
	case consts.JobStateFailed:
		eventTime = job.EndTime
	}

	if job.Type == common.Local {
		eventName = fmt.Sprintf("%v%v", consts.TimelineLocalPrefix, job.State)
	}

	return eventName, eventTime
}
