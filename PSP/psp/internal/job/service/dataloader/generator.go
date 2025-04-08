package dataloader

import (
	"context"
	"reflect"

	"github.com/yuansuan/ticp/common/go-kit/logging"
	schema "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"

	"github.com/yuansuan/ticp/PSP/psp/internal/job/dao/model"
	"github.com/yuansuan/ticp/PSP/psp/internal/job/util"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/strutil"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/timeutil"
)

// generateJobUpdateInfo 生成作业同步更新的数据
func (loader *JobLoader) generateJobUpdateInfo(ctx context.Context, openapiJob *schema.AdminJobInfo, job *model.Job) (*model.Job, []string, string) {
	logger := logging.GetLogger(ctx)

	if openapiJob == nil {
		return nil, nil, ""
	}

	fieldNum := reflect.TypeOf(*job).NumField()
	cols := make([]string, 0, fieldNum)

	originJobId := openapiJob.OriginJobID
	if job.RealJobId != originJobId {
		job.RealJobId = originJobId
		cols = append(cols, "real_job_id")
	}

	queueName := openapiJob.Queue
	if !strutil.IsEmpty(queueName) && job.Queue != queueName {
		job.Queue = queueName
		cols = append(cols, "queue")
	}

	var preJobState string
	jobState := openapiJob.JobState
	if !strutil.IsEmpty(jobState) {
		if job.RawState != jobState {
			job.RawState = jobState
			cols = append(cols, "raw_state")
		}

		state := util.ConvertJobState(jobState)
		if strutil.IsEmpty(state) {
			state = job.State
		}

		if job.State != state {
			preJobState = job.State
			job.State = state
			cols = append(cols, "state")
		}
	}

	exitCode := openapiJob.ExitCode
	if job.ExitCode != exitCode {
		job.ExitCode = exitCode
		cols = append(cols, "exit_code")
	}

	stateReason := openapiJob.StateReason
	if job.Reason != stateReason {
		job.Reason = stateReason
		cols = append(cols, "reason")
	}

	workDir := openapiJob.Workdir
	if job.WorkDir != workDir {
		job.WorkDir = workDir
		cols = append(cols, "work_dir")
	}

	zoneName := openapiJob.Zone
	if job.ClusterName != zoneName {
		job.ClusterName = zoneName
		cols = append(cols, "cluster_name")
	}

	execHosts := openapiJob.ExecHosts
	if job.ExecHosts != execHosts {
		job.ExecHosts = execHosts
		cols = append(cols, "exec_hosts")
	}

	execHostNum := openapiJob.ExecHostNum
	if job.ExecHostNum != execHostNum {
		job.ExecHostNum = execHostNum
		cols = append(cols, "exec_host_num")
	}

	execDuration := openapiJob.ExecutionDuration
	if job.ExecDuration != execDuration {
		job.ExecDuration = execDuration
		cols = append(cols, "exec_duration")
	}

	allocResource := openapiJob.AllocResource
	if allocResource != nil {
		if job.CpusAlloc != allocResource.Cores {
			job.CpusAlloc = allocResource.Cores
			cols = append(cols, "cpus_alloc")
		}

		if job.MemAlloc != allocResource.Memory {
			job.MemAlloc = allocResource.Memory
			cols = append(cols, "mem_alloc")
		}
	}

	priority := openapiJob.Priority
	if job.Priority != priority {
		job.Priority = priority
		cols = append(cols, "priority")
	}

	createTime, err := timeutil.ParseJsonTime(openapiJob.CreateTime)
	if err != nil {
		logger.Errorf("parse openapi create time [%v] err: %v", openapiJob.CreateTime, err)
	} else {
		if !job.SubmitTime.Equal(createTime) {
			job.SubmitTime = createTime
			cols = append(cols, "submit_time")
		}
	}

	pendingTime, err := timeutil.ParseJsonTime(openapiJob.PendingTime)
	if err != nil {
		logger.Errorf("parse openapi create time [%v] err: %v", openapiJob.CreateTime, err)
	} else {
		if !job.PendTime.Equal(pendingTime) {
			job.PendTime = pendingTime
			cols = append(cols, "pend_time")
		}
	}

	runningTime, err := timeutil.ParseJsonTime(openapiJob.RunningTime)
	if err != nil {
		logger.Errorf("parse openapi create time [%v] err: %v", openapiJob.RunningTime, err)
	} else {
		if !job.StartTime.Equal(runningTime) {
			job.StartTime = runningTime
			cols = append(cols, "start_time")
		}
	}

	endTime, err := timeutil.ParseJsonTime(openapiJob.EndTime)
	if err != nil {
		logger.Errorf("parse openapi create time [%v] err: %v", openapiJob.RunningTime, err)
	} else {
		if !job.EndTime.Equal(endTime) {
			job.EndTime = endTime
			cols = append(cols, "end_time")
		}
	}

	terminateTime, err := timeutil.ParseJsonTime(openapiJob.TerminatingTime)
	if err != nil {
		logger.Errorf("parse openapi create time [%v] err: %v", openapiJob.RunningTime, err)
	} else {
		if !job.TerminateTime.Equal(terminateTime) {
			job.TerminateTime = terminateTime
			cols = append(cols, "terminate_time")
		}
	}

	suspendTime, err := timeutil.ParseJsonTime(openapiJob.SuspendedTime)
	if err != nil {
		logger.Errorf("parse openapi create time [%v] err: %v", openapiJob.RunningTime, err)
	} else {
		if !job.SuspendTime.Equal(suspendTime) {
			job.SuspendTime = suspendTime
			cols = append(cols, "suspend_time")
		}
	}

	return job, cols, preJobState
}
