package pbspro

import (
	"bufio"
	"context"
	"fmt"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/backend/common"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"

	"github.com/yuansuan/ticp/iPaaS/standard-compute/config"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/backend/job"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/cmdhelp"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/dao/models"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/oshelp"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/util"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/pkg/fsutil/filemode"
)

const (
	scriptFileName  = "__script.sh"
	commandFileName = "__command.sh"
	envFileName     = "__env"
	stdoutFileName  = "stdout.log"
	stderrFileName  = "stderr.log"
	nullFile        = "/dev/null"
	submitEndWith   = "\"${script}\""
)

var singularityScriptTpl = `#!/bin/bash
set -euo pipefail

exec 1> {{ .Stdout }} 2> {{ .Stderr }}

cd {{ .Workspace }}

echo YS_NODELIST=$(paste -sd, $PBS_NODEFILE) >> ` + envFileName + `
echo YS_CPUS_ON_NODE=$NCPUS >> ` + envFileName + `
echo YS_SUBMIT_DIR=$PBS_O_WORKDIR >> ` + envFileName + `
echo YS_NUM_NODES=$(cat $PBS_NODEFILE | wc -w) >> ` + envFileName + `

singularity run --contain --bind ~/.ssh --bind {{.PreparedFilePath}} --bind {{ .Workspace }} --cleanenv --env-file ` + envFileName + ` --pwd {{ .Workspace }} {{ .AppPath }} bash -euo pipefail ` + commandFileName + `

`

var localAppScript = `#!/bin/bash
set -euo pipefail

exec 1> {{ .Stdout }} 2> {{ .Stderr }}

cd {{ .Workspace }}

echo YS_NODELIST=$(paste -sd, $PBS_NODEFILE) >> ` + envFileName + `
echo YS_CPUS_ON_NODE=$NCPUS >> ` + envFileName + `
echo YS_SUBMIT_DIR=$PBS_O_WORKDIR >> ` + envFileName + `
echo YS_NUM_NODES=$(cat $PBS_NODEFILE | wc -w) >> ` + envFileName + `

source ./__env

{{ .Command }}
`

type ExecFunc func(ctx context.Context, cmdStr string, opts ...cmdhelp.Option) (stdOut string, stdErr string, err error)

type Provider struct {
	commonCfg config.SchedulerCommon
	customCfg *config.PbsProBackendProvider
	execFunc  ExecFunc
}

func NewProvider(commonCfg config.SchedulerCommon, customCfg *config.PbsProBackendProvider, execFunc ExecFunc) *Provider {
	return &Provider{
		commonCfg: commonCfg,
		customCfg: customCfg,
		execFunc:  execFunc,
	}
}

func (p Provider) Submit(ctx context.Context, j *job.Job) (string, error) {
	j.Script = filepath.Join(j.Workspace, scriptFileName)
	j.Stdout = filepath.Join(j.Workspace, stdoutFileName)
	j.Stderr = filepath.Join(j.Workspace, stderrFileName)
	commandFilepath := filepath.Join(j.Workspace, commandFileName)
	envFilepath := filepath.Join(j.Workspace, envFileName)

	if err := os.MkdirAll(j.Workspace, filemode.Directory); err != nil {
		return "", err
	}

	scriptTpl := singularityScriptTpl
	if j.AppMode == models.LocalAppMode {
		scriptTpl = localAppScript
	}

	osHelpOpts := make([]oshelp.Option, 0)
	if p.commonCfg.SubmitSysUser != "" {
		osHelpOpts = append(osHelpOpts, oshelp.WithChown(p.commonCfg.SubmitSysUser))
	}

	// generate script file
	script, err := util.RenderScript(context.TODO(), j, scriptTpl, "pbs pro")
	if err != nil {
		return "", errors.Wrap(err, "pbs pro submit")
	}
	err = func() error {
		fd, err := os.Create(j.Script)
		if err != nil {
			return err
		}
		defer func() { _ = fd.Close() }()

		_, err = oshelp.Write(fd, script, osHelpOpts...)
		return err
	}()
	if err != nil {
		return "", errors.Wrap(err, "pbs pro submit")
	}

	// generate env file
	envContent, err := util.RenderEnvVars(ctx, j, "pbs pro")
	if err != nil {
		return "", errors.Wrap(err, "pbs pro submit")
	}
	err = func() error {
		fd, err := os.Create(envFilepath)
		if err != nil {
			return err
		}
		defer fd.Close()

		_, err = oshelp.Write(fd, envContent, osHelpOpts...)
		return err
	}()
	if err != nil {
		return "", errors.Wrap(err, "pbs pro submit")
	}

	// generate command file
	err = func() error {
		fd, err := os.Create(commandFilepath)
		if err != nil {
			return err
		}
		defer func() { _ = fd.Close() }()

		_, err = oshelp.Write(fd, []byte(j.Command), osHelpOpts...)
		return err
	}()
	if err != nil {
		return "", errors.Wrap(err, "pbs pro submit")
	}
	nTasksPerNode := util.EnsureNTaskPerNode(p.commonCfg, j)
	nodes := util.OccupiedNodesNum(int(j.RequestCores), nTasksPerNode)
	memPerNode := int(math.Ceil(float64(j.RequestMemory) / float64(nodes)))
	queue := p.commonCfg.DefaultQueue
	if j.Queue != "" {
		queue = j.Queue
	}

	envs := []string{
		util.EnvPair("cwd", j.Workspace),
		// pbs的job的标准输出只有在job结束的时候才flush信息, 所以在__script里把标准输出和标准错误重定向到job.Stdout和job.Stderr里了,
		util.EnvPair("out", nullFile),
		util.EnvPair("err", nullFile),
		util.EnvPair("nodes", nodes),
		util.EnvPair("number_of_cpu", nTasksPerNode),
		util.EnvPair("script", scriptFileName),
		util.EnvPair("memory_mb", memPerNode),
		util.EnvPair("queue", queue),
	}

	submitOpts := []cmdhelp.Option{
		cmdhelp.WithCmdEnv(envs),
		cmdhelp.WithCmdDir(j.Workspace),
	}
	if p.commonCfg.SubmitSysUser != "" {
		submitOpts = append(submitOpts, cmdhelp.WithCmdUser(p.commonCfg.SubmitSysUser))
	}

	cmdStr, err := cmdhelp.EnsureSubmitCommand(p.customCfg.Submit, j.SchedulerSubmitFlags, submitEndWith)
	if err != nil {
		return "", fmt.Errorf("ensure submit command failed, %w", err)
	}

	var stdout, stderr string
	stdout, stderr, err = p.execFunc(ctx, cmdStr, submitOpts...)

	if err != nil {
		return "", errors.Wrapf(err, "pbs pro submit: %v, %v", stdout, stderr)
	}
	stdoutSplit := strings.Split(stdout, ".")
	//判断是否以整数(job id)开头
	_, err = strconv.ParseInt(stdoutSplit[0], 10, 64)
	if err != nil || len(stdoutSplit) != 2 {
		return "", errors.Wrap(err, "pbs pro submit")
	} else {
		return stdoutSplit[0], nil
	}

}

func (p Provider) Kill(ctx context.Context, j *job.Job) error {
	envs := []string{
		util.EnvPair("job_id", j.OriginJobId),
	}
	stdout, stderr, err := p.execFunc(ctx, p.customCfg.Kill, cmdhelp.WithCmdEnv(envs))
	if err != nil {
		return errors.Wrapf(err, "pbspro kill: %v, %v", stdout, stderr)
	}
	return nil
}

func (p Provider) CheckAlive(ctx context.Context, j *job.Job) (*job.Job, error) {
	envs := []string{
		util.EnvPair("job_id", j.OriginJobId),
	}
	stdout, stderr, err := cmdhelp.ExecShellCmd(ctx, p.customCfg.CheckAlive, cmdhelp.WithCmdEnv(envs))
	if err != nil {
		return j, errors.Wrapf(err, "pbspro check-alive: %v, %v", stdout, stderr)
	}

	detail := GetJobFromQstat(stdout)
	updateByDetail(j, detail)
	util.UpdateByEnv(j)

	return j, nil
}

func GetJobFromQstat(out string) map[string]string {
	fields := strings.Split(out, "\n")
	ret := map[string]string{}

	for _, field := range fields[1:] {
		arr := strings.Split(field, "=")
		if len(arr) == 2 {
			key, value := strings.TrimSpace(arr[0]), strings.TrimSpace(arr[1])
			ret[key] = value
		}
	}
	return ret
}

func updateByDetail(j *job.Job, detail map[string]string) {
	j.AllocCores, _ = strconv.ParseInt(detail["resources_used.ncpus"], 10, 64)
	j.OriginState = detail["job_state"]
	// 和slurm保持统一的退出码格式
	if len(detail["Exit_status"]) > 0 {
		j.ExitCode = fmt.Sprintf("%s:0", detail["Exit_status"])
	}

	runningTime := convert2Time(detail["stime"])
	j.BackendJobState = mappingPbsJobState(j.OriginState, runningTime == nil) // pbs在作业running时，有一段时间无法查询到stime，此时将状态保持pending

	j.PendingTime = convert2Time(detail["etime"])
	switch {
	case j.BackendJobState == job.StatePending:
	case j.BackendJobState <= job.StateRunning:
		j.RunningTime = runningTime
		if j.RunningTime != nil {
			j.ExecutionDuration = int64(time.Now().Sub(*j.RunningTime) / time.Second)
		}
	case j.BackendJobState <= job.StateCompleted:
		j.RunningTime = runningTime
		j.CompletingTime = convert2Time(detail["mtime"])
		if j.RunningTime != nil && j.CompletingTime != nil {
			j.ExecutionDuration = int64(j.CompletingTime.Sub(*j.RunningTime) / time.Second)
		}
	}

	// TODO Priority
}

func convert2Time(date string) *time.Time {
	return convert2TimeWithLocation(date, time.Local)
}

func convert2TimeWithLocation(jobDate string, location *time.Location) *time.Time {
	const DateFormat = "Mon Jan 2 15:04:05 2006"
	dateTime, err := time.ParseInLocation(DateFormat, jobDate, location)
	if err != nil { // invalid value maybe "StartTime=Unknown EndTime=Unknown"
		return nil
	}
	return &dateTime
}

func (p Provider) NewWorkspace() string {
	return filepath.Join(p.commonCfg.Workspace, uuid.NewString())
}

const (
	pbsJobStateInvalid = ""
	//Array job has at least one subjob running
	pbsJobStateBegin = "B"
	//Job is exiting after having run
	pbsJobStateExist = "E"
	//Job is finished. Job has completed execution, job failed during execution, or job was deleted
	pbsJobStateFinished = "F"
	//Job is held. A job is put into a held state by the server or by a user or administrator.
	//A job stays in a held state until it is released by a user or administrator.
	pbsJobStateHeld = "H"
	//Job was moved to another server
	pbsJobStateMoved = "M"
	//Job is queued, eligible to run or be routed
	pbsJobStateQueued = "Q"
	//Job is running
	pbsJobStateRunning = "R"
	//Job is suspended by scheduler. A job is put into the suspended state when a higher priority job
	//needs the resources.sub-state of Running
	pbsJobStateSuspended = "S"
	//Job is being moved to new location
	pbsJobStateTransition = "T"
	//Job is suspended due to workstation becoming busy.sub-state of Running
	pbsJobStateU = "U"
	//Job is waiting for its requested execution time to be reached or job specified a stage in request
	//which failed for some reason.
	pbsJobStateWait = "W"
	//Subjobs only; subjob is finished (expired.)
	pbsJobStateX = "X"
)

func mappingPbsJobState(jobState string, runningTimeNil bool) job.State {
	switch jobState {
	case pbsJobStateQueued, pbsJobStateWait, pbsJobStateHeld, pbsJobStateMoved, pbsJobStateTransition:
		return job.StatePending
	case pbsJobStateRunning, pbsJobStateSuspended, pbsJobStateU, pbsJobStateBegin:
		if runningTimeNil { // pbs在作业running时，有一段时间无法查询到stime，此时将状态保持pending
			return job.StatePending
		}
		return job.StateRunning
	case pbsJobStateExist, pbsJobStateFinished, pbsJobStateX:
		return job.StateCompleted
	default:
		return job.StateCompleted
	}
}

func (p Provider) GetFreeResource(ctx context.Context, queues []string) (map[string]*v20230530.Resource, error) {
	stdout, stderr, err := p.execFunc(ctx, p.customCfg.GetResource)
	if err != nil {
		return nil, errors.Wrapf(err, "ExecmdFail, Stdout: %s, Stderr: %s", stdout, stderr)
	}
	m := parseResource(stdout, p.commonCfg.DefaultQueue)

	for k, v := range m {
		if k == p.commonCfg.DefaultQueue {
			v.IsDefault = true
		}
	}

	return m, nil
}

// ref resourceInfoTest.txt

func parseResource(stdout, defaultQueue string) map[string]*v20230530.Resource {

	queueMap := make(map[string]*v20230530.Resource)

	var nodeName, queueName string
	var mem, ncpus int
	scanner := bufio.NewScanner(strings.NewReader(stdout))

	for scanner.Scan() {
		line := scanner.Text()

		if len(line) > 0 && line[0] != ' ' {
			// This is a node name
			nodeName = strings.TrimSpace(line)
			queueName = defaultQueue // Reset queueName for each new node
		} else if strings.Contains(line, "queue") {
			queueName = strings.Split(line, "=")[1]
			queueName = strings.TrimSpace(queueName)
		} else if strings.Contains(line, "resources_available.mem") {
			memStr := strings.Split(line, "=")[1]
			memStr = strings.Replace(memStr, "kb", "", -1)
			memStr = strings.TrimSpace(memStr)
			mem, _ = strconv.Atoi(memStr)
		} else if strings.Contains(line, "resources_available.ncpus") {
			ncpusStr := strings.Split(line, "=")[1]
			ncpusStr = strings.TrimSpace(ncpusStr)
			ncpus, _ = strconv.Atoi(ncpusStr)
		} else if line == "" {
			// This is the end of a node block
			if info, ok := queueMap[queueName]; ok {
				info.Memory += int64(mem)
				info.Cpu += int64(ncpus)
				queueMap[queueName] = info
			} else {
				queueMap[queueName] = &v20230530.Resource{Memory: int64(mem), Cpu: int64(ncpus)}
			}
			nodeName = ""
			mem = 0
			ncpus = 0
		}
	}

	// Check if there's any node info left to add to the map
	if nodeName != "" {
		if info, ok := queueMap[queueName]; ok {
			info.Memory += int64(mem)
			info.Cpu += int64(ncpus)
			queueMap[queueName] = info
		} else {
			queueMap[queueName] = &v20230530.Resource{Memory: int64(mem), Cpu: int64(ncpus)}
		}
	}

	return queueMap
}
func (p *Provider) GetCpuUsage(ctx context.Context, j *models.Job, nodes []string, adjustFactor float64) (*v20230530.CpuUsage, error) {
	return common.GetCpuUsage(ctx, j, nodes, adjustFactor)
}
