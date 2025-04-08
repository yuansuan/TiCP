package slurm

import (
	"context"
	"fmt"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/backend/common"
	"os"
	"path/filepath"
	"regexp"
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
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/log"
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
	// 与配置文件中提交命令的尾部需对应
	submitEndWith = "--wrap \"/bin/bash ${script}\""
)

var singularityScriptTpl = `#!/bin/bash
set -euo pipefail

echo "YS_NODELIST=$(scontrol show hostname $SLURM_NODELIST | paste -d, -s)" >> ` + envFileName + `
echo "YS_CPUS_ON_NODE=$SLURM_CPUS_ON_NODE" >> ` + envFileName + `
echo "YS_SUBMIT_DIR=$SLURM_SUBMIT_DIR" >> ` + envFileName + `
echo "YS_NUM_NODES=$SLURM_JOB_NUM_NODES" >> ` + envFileName + `

singularity run --contain --bind ~/.ssh --bind {{.PreparedFilePath}} --bind {{ .Workspace }} --bind /tmp --cleanenv --env-file ` + envFileName + ` --pwd {{ .Workspace }} {{ .AppPath }} bash -euo pipefail ` + commandFileName + `

`

var localAppScript = `#!/bin/bash
set -euo pipefail

echo "YS_NODELIST=$(scontrol show hostname $SLURM_NODELIST | paste -d, -s)" >> ` + envFileName + `
echo "YS_CPUS_ON_NODE=$SLURM_CPUS_ON_NODE" >> ` + envFileName + `
echo "YS_SUBMIT_DIR=$SLURM_SUBMIT_DIR" >> ` + envFileName + `
echo "YS_NUM_NODES=$SLURM_JOB_NUM_NODES" >> ` + envFileName + `

source ./__env

{{ .Command }}
`

type ExecFunc func(ctx context.Context, cmdStr string, opts ...cmdhelp.Option) (stdOut string, stdErr string, err error)

type Provider struct {
	commonCfg config.SchedulerCommon
	customCfg *config.SlurmBackendProvider
	execFunc  ExecFunc
}

func NewProvider(commonCfg config.SchedulerCommon, customCfg *config.SlurmBackendProvider, execFunc ExecFunc) *Provider {
	return &Provider{
		commonCfg: commonCfg,
		customCfg: customCfg,
		execFunc:  execFunc,
	}
}

func (p *Provider) Submit(ctx context.Context, j *job.Job) (string, error) {
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
	script, err := util.RenderScript(context.TODO(), j, scriptTpl, "slurm")
	if err != nil {
		return "", errors.Wrap(err, "slurm submit")
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
		return "", errors.Wrap(err, "slurm submit")
	}

	// generate env file
	envContent, err := util.RenderEnvVars(ctx, j, "slurm")
	if err != nil {
		return "", errors.Wrap(err, "slurm submit")
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
		return "", errors.Wrap(err, "slurm submit")
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
		return "", errors.Wrap(err, "slurm submit")
	}

	var submitCmd string
	var envs []string

	queue := p.commonCfg.DefaultQueue
	if j.Queue != "" {
		queue = j.Queue
	}

	commonEnvs := []string{
		util.EnvPair("cwd", j.Workspace),
		util.EnvPair("out", j.Stdout),
		util.EnvPair("err", j.Stderr),
		util.EnvPair("script", scriptFileName),
		util.EnvPair("memory_mb", j.RequestMemory),
		util.EnvPair("queue", queue),
	}

	// Choose submit command and add specific environment variables based on allocation type
	if j.AllocType == "average" {
		submitCmd = p.customCfg.SubmitAverage
		envs = append(commonEnvs, util.EnvPair("ntasks", j.RequestCores))
	} else {
		submitCmd = p.customCfg.Submit
		nTasksPerNode := util.EnsureNTaskPerNode(p.commonCfg, j)
		nodes := util.OccupiedNodesNum(int(j.RequestCores), nTasksPerNode)
		envs = append(commonEnvs,
			util.EnvPair("nodes", nodes),
			util.EnvPair("ntasks_per_node", nTasksPerNode),
		)
	}

	submitOpts := []cmdhelp.Option{
		cmdhelp.WithCmdEnv(envs),
		cmdhelp.WithCmdDir(j.Workspace),
	}
	if p.commonCfg.SubmitSysUser != "" {
		submitOpts = append(submitOpts, cmdhelp.WithCmdUser(p.commonCfg.SubmitSysUser))
	}

	cmdStr, err := cmdhelp.EnsureSubmitCommand(submitCmd, j.SchedulerSubmitFlags, submitEndWith)
	log.Infof("SLURM submit command: %s", cmdStr) //执行的啥命令打印出来

	if err != nil {
		return "", fmt.Errorf("ensure submit command failed, %w", err)
	}
	stdout, stderr, err := p.execFunc(ctx, cmdStr, submitOpts...)

	if err != nil {
		return "", errors.Wrapf(err, "slurm submit: %v, %v", stdout, stderr)
	}

	r, err := regexp.Compile(p.customCfg.JobIdRegex)
	if err != nil {
		return "", errors.Wrap(err, "slurm submit")
	}

	matchs := r.FindStringSubmatch(stdout)
	if len(matchs) > 1 {
		return matchs[1], nil
	}
	return "", errors.Errorf("slurm submit: not found job id: %v", stdout)
}

func (p *Provider) Kill(ctx context.Context, j *job.Job) error {
	envs := []string{
		util.EnvPair("job_id", j.OriginJobId),
	}
	stdout, stderr, err := p.execFunc(ctx, p.customCfg.Kill, cmdhelp.WithCmdEnv(envs))
	if err != nil {
		return errors.Wrapf(err, "slurm kill: %v, %v", stdout, stderr)
	}
	return nil
}

func (p *Provider) CheckAlive(ctx context.Context, j *job.Job) (*job.Job, error) {
	envs := []string{
		util.EnvPair("job_id", j.OriginJobId),
	}
	stdout, stderr, err := p.execFunc(ctx, p.customCfg.CheckAlive, cmdhelp.WithCmdEnv(envs))
	// 打印 CheckAlive 的原始输出
	log.Debugf("CheckAlive Raw Output - Stdout:\n%s", stdout)
	log.Debugf("CheckAlive Raw Output - Stderr:\n%s", stderr)

	if err == nil {
		detail := GetJobFromScontrol(stdout)
		updateByDetail(j, detail)
		util.UpdateByEnv(j)
		return j, nil
	}
	log.Infof("scontrol check-alive failed: %v, %v", stdout, stderr)

	// If scontrol fails, try sacct
	historyStdout, historyStderr, historyErr := p.execFunc(ctx, p.customCfg.CheckHistory, cmdhelp.WithCmdEnv(envs))
	log.Debugf("CheckHistory Raw Output - Stdout:\n%s", historyStdout)
	log.Debugf("CheckHistory Raw Output - Stderr:\n%s", historyStderr)

	if historyErr == nil {
		detail := GetJobFromSacct(historyStdout)
		updateByDetail(j, detail)
		util.UpdateByEnv(j)
		return j, nil
	}
	log.Warnf("sacct check-history failed: %v, %v", historyStdout, historyStderr)

	return j, errors.New("failed to check job status with both scontrol and sacct")
}

func (p *Provider) NewWorkspace() string {
	return filepath.Join(p.commonCfg.Workspace, uuid.NewString())
}

func convert2Time(date string) *time.Time {
	return convert2TimeWithLocation(date, time.Local)
}

func convert2TimeWithLocation(jobDate string, location *time.Location) *time.Time {
	const DateFormat = "2006-01-02T15:04:05"
	dateTime, err := time.ParseInLocation(DateFormat, jobDate, location)
	if err != nil { // invalid value maybe "StartTime=Unknown EndTime=Unknown"
		return nil
	}
	return &dateTime
}

func updateByDetail(j *job.Job, detail map[string]string) {
	readMin := func(input string) string {
		arr := strings.Split(input, "-")
		return arr[0]
	}

	tresMap := map[string]string{}
	for _, v := range strings.Split(detail["TRES"], ",") {
		vv := strings.SplitN(v, "=", 2)
		if len(vv) == 2 {
			tresMap[vv[0]] = vv[1]
		}
	}
	var nCpus int64
	if cpu, ok := tresMap["cpu"]; ok {
		nCpus, _ = strconv.ParseInt(cpu, 10, 64)
	}

	if nCpus == 0 {
		nCpus, _ = strconv.ParseInt(readMin(detail["NumCPUs"]), 10, 64)
	}

	j.AllocCores = nCpus
	j.OriginState = detail["JobState"]
	j.ExitCode = detail["ExitCode"]

	j.BackendJobState = mappingSlurmJobState(j.OriginState)

	j.PendingTime = convert2Time(detail["SubmitTime"])
	switch {
	case j.BackendJobState == job.StatePending:
	case j.BackendJobState <= job.StateRunning:
		j.RunningTime = convert2Time(detail["StartTime"])
		if j.RunningTime != nil {
			j.ExecutionDuration = int64(time.Now().Sub(*j.RunningTime) / time.Second)
		}
	case j.BackendJobState <= job.StateCompleted:
		j.RunningTime = convert2Time(detail["StartTime"])
		j.CompletingTime = convert2Time(detail["EndTime"])
		if j.RunningTime != nil && j.CompletingTime != nil {
			j.ExecutionDuration = int64(j.CompletingTime.Sub(*j.RunningTime) / time.Second)
		}
	}
	j.Priority, _ = strconv.ParseInt(detail["Priority"], 10, 64)
}

func GetJobFromScontrol(out string) map[string]string {
	fields := strings.Fields(out)
	ret := map[string]string{}

	for _, field := range fields {
		arr := strings.SplitN(field, "=", 2)
		if len(arr) == 2 {
			key, value := arr[0], arr[1]
			ret[key] = value
		}
	}
	return ret
}

func GetJobFromSacct(out string) map[string]string {
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) < 2 {
		return nil
	}

	keys := strings.Split(lines[0], "|") // 读取第一行作为map的key

	values := strings.Split(lines[1], "|") // 读取第二行作为map的value

	result := make(map[string]string)
	for i, key := range keys {
		// 确保key和value数量一致
		if i < len(values) {
			result[key] = values[i]
		} else {
			result[key] = ""
		}
	}
	// Map sacct fields to scontrol fields
	if allocCPUs, exists := result["AllocCPUS"]; exists {
		result["NumCPUs"] = allocCPUs
	}
	if state, exists := result["State"]; exists {
		result["JobState"] = state
	}
	if submit, exists := result["Submit"]; exists {
		result["SubmitTime"] = submit
	}
	if start, exists := result["Start"]; exists {
		result["StartTime"] = start
	}
	if end, exists := result["End"]; exists {
		result["EndTime"] = end
	}

	return result
}

const (
	slurmJobStateInvalid = ""

	// Job terminated due to launch failure, typically due to a hardware failure
	// (e.g. unable to boot the node or block and the job can not be requeued).
	slurmJobStateBootFail = "boot_fail"

	// Job was explicitly cancelled by the user or system administrator.
	// The job may or may not have been initiated.
	slurmJobStateCancelled = "cancelled"

	// Job terminated on deadline.
	slurmJobStateDeadline = "deadline"

	// Job terminated with non-zero exit code or other failure condition.
	slurmJobStateFailed = "failed"

	// Job terminated due to failure of one or more allocated nodes.
	// configure: Stipulates whether a job should be requeued after a node failure: 0 for no, 1 for yes.
	// reference: https://slurm.schedmd.com/scontrol.html
	slurmJobStateNodeFail = "node_fail"

	// Job experienced out of memory error.
	slurmJobStateOutOfMemory = "out_of_memory"

	// Job terminated due to preemption.
	slurmJobStatePreempted = "preempted"

	// Job is about to change size.
	slurmJobStateResizing = "resizing"

	// Sibling was removed from cluster due to other cluster starting the job.
	slurmJobStateRevoked = "revoked"

	// Job terminated upon reaching its time limit.
	slurmJobStateTimeout = "timeout"

	// https://slurm.schedmd.com/sacct.html 中未找到
	// slurmJobStateSuccess   = "success"

	slurmJobStateCompleted = "completed"
	slurmJobStatePending   = "pending"
	slurmJobStateRunning   = "running"
	slurmJobStateRequeued  = "requeued"
	slurmJobStateSuspended = "suspended"

	slurmNodeIdleStatus  = "idle"
	slurmNodeAllocStatus = "alloc"
	slurmNodeMixStatus   = "mix"
)

// mappingSlurmJobState maps sc states to gw states
func mappingSlurmJobState(jobState string) job.State {
	switch strings.ToLower(strings.Fields(jobState)[0]) {
	case slurmJobStatePending:
		return job.StatePending
	case slurmJobStateRunning:
		return job.StateRunning
	case slurmJobStateCompleted:
		return job.StateCompleted
	case slurmJobStateCancelled:
		return job.StateCompleted
	case slurmJobStateFailed:
		return job.StateCompleted
	case slurmJobStateTimeout:
		return job.StateCompleted
	case slurmJobStateNodeFail:
		return job.StateCompleted
	case slurmJobStateBootFail:
		return job.StateCompleted
	case slurmJobStateDeadline:
		return job.StateCompleted
	case slurmJobStateOutOfMemory:
		return job.StateCompleted
	case slurmJobStatePreempted:
		return job.StateCompleted
	case slurmJobStateRequeued:
		return job.StatePending
	default:
		return job.StateCompleted
	}
}

func (p *Provider) getFreeResource(ctx context.Context, queue string) (*v20230530.Resource, error) {
	env := []string{fmt.Sprintf("queue=%s", queue)}
	stdout, stderr, err := p.execFunc(ctx, p.customCfg.GetResource, cmdhelp.WithCmdEnv(env))
	if err != nil {
		return nil, errors.Wrapf(err, "slurm ExecmdFail, Stdout: %s, Stderr: %s", stdout, stderr)
	}
	return parserResourceInfo(stdout)
}

func (p *Provider) GetFreeResource(ctx context.Context, queues []string) (map[string]*v20230530.Resource, error) {
	res := make(map[string]*v20230530.Resource)
	for _, queue := range queues {
		r, err := p.getFreeResource(ctx, queue)
		if err != nil {
			return nil, err
		}
		if queue == p.commonCfg.DefaultQueue {
			r.IsDefault = true
		}
		res[queue] = r
	}
	return res, nil
}

func parserResourceInfo(s string) (*v20230530.Resource, error) {
	var err error
	// cmd stdout like below
	//STATE               CPUS                FREE_MEM            MEMORY
	//idle                2                   210                 15000
	//idle                6                   15390               15000
	if 2 > len(s) {
		return nil, errors.New(fmt.Sprintf("UnexpectedResourceInfo, %s", s))
	}

	regexStr := "(\\S+)\\s+(\\d+)\\s+(\\S+)\\s+(\\d+)"
	r, _ := regexp.Compile(regexStr)
	sList := strings.Split(strings.TrimSpace(s), "\n")
	res := new(v20230530.Resource)
	res.TotalNodeNum = int64(len(sList[1:]))

	for _, line := range sList[1:] {
		matchs := r.FindStringSubmatch(strings.TrimSpace(line))
		if len(matchs) == 5 {
			cpuNum, _ := strconv.Atoi(matchs[2])
			idleMemNum, err := strconv.Atoi(matchs[3])
			if err != nil {
				log.Warnf("parserResourceInfo, strconv.Atoi(%s) failed, err: %v", matchs[3], err)
				continue
			}
			totalMemNum, _ := strconv.Atoi(matchs[4])
			switch matchs[1] {
			case slurmNodeIdleStatus:
				res.Cpu += int64(cpuNum)
				res.IdleNodeNum++
			case slurmNodeAllocStatus, slurmNodeMixStatus:
				res.AllocNodeNum++
			}
			res.TotalCpu += int64(cpuNum)
			res.Memory += int64(idleMemNum)
			res.TotalMemory += int64(totalMemNum)
		} else {
			err = errors.New(fmt.Sprintf("UnexpectedResourceInfoLine, info: %s", line))
			break
		}
	}
	return res, err
}
func (p *Provider) GetCpuUsage(ctx context.Context, j *models.Job, nodes []string, adjustFactor float64) (*v20230530.CpuUsage, error) {
	return common.GetCpuUsage(ctx, j, nodes, adjustFactor)
}
