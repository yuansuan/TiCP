package slurm

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/config"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/cmdhelp"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/backend/job"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/dao/models"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/util"
)

func TestRenderScript(t *testing.T) {
	j := &job.Job{
		Job: &models.Job{
			Command:   "echo 123",
			AppPath:   "/apps/dummy_image.sif",
			Workspace: "/tmp/workspace/job-id-1",
		},
		EnvVars: []string{
			"ENV_A=a",
			"ENV_B=b",
			"ENV_C=c",
			"ENV_D=d",
		},
	}
	out, err := util.RenderScript(context.TODO(), j, singularityScriptTpl, "slurm")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(out)
}

func TestSubmit(t *testing.T) {
	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "slurm_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	commonCfg := config.SchedulerCommon{
		DefaultQueue: "default",
		Workspace:    tempDir,
	}
	customCfg := &config.SlurmBackendProvider{
		Submit:     `sbatch ${script} --wrap "/bin/bash ${script}"`,
		JobIdRegex: `Submitted batch job (\d+)`,
	}

	provider := NewProvider(commonCfg, customCfg, fake_ExecShellCmd)

	t.Run("ValidSubmission", func(t *testing.T) {
		j := &job.Job{
			Job: &models.Job{
				Command:       "echo 'Hello, World!'",
				Workspace:     filepath.Join(tempDir, "job-1"),
				RequestCores:  2,
				RequestMemory: 1000,
			},
		}

		jobID, err := provider.Submit(context.Background(), j)

		assert.NoError(t, err)
		assert.Equal(t, "12345", jobID)
	})
}

func TestKill(t *testing.T) {
	j := &job.Job{
		Job: &models.Job{
			OriginJobId: "1",
		},
	}
	conf := &config.SlurmBackendProvider{
		Kill: "scancel ${job_id}",
	}
	p := NewProvider(config.SchedulerCommon{}, conf, fake_ExecShellCmd)

	err := p.Kill(context.TODO(), j)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCheckAlive(t *testing.T) {
	// 测试 scontrol 情况
	t.Run("TestScontrol", func(t *testing.T) {
		SetUseSacct(false)
		runCheckAliveTest(t, "129", 8)
	})

	// 测试 sacct 情况
	t.Run("TestSacct", func(t *testing.T) {
		SetUseSacct(true)
		runCheckAliveTest(t, "19411344", 560)
	})
}

func runCheckAliveTest(t *testing.T, jobID string, expectedCores int64) {
	// Initialize logger
	logConfig := config.Log{
		Level:        "debug",
		ReleaseLevel: "development",
		UseConsole:   true,
	}
	if err := log.InitLogger(logConfig); err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}

	j := &job.Job{
		Job: &models.Job{
			OriginJobId: jobID,
		},
	}

	mockProvider := &config.SlurmBackendProvider{
		CheckAlive:   "scontrol show job ${job_id}",
		CheckHistory: "sacct -j ${job_id}",
	}

	p := NewProvider(config.SchedulerCommon{}, mockProvider, fake_ExecShellCmd)

	result, err := p.CheckAlive(context.Background(), j)
	if err != nil {
		t.Fatalf("CheckAlive failed: %v", err)
	}
	// 检查 JobId
	if result.OriginJobId != jobID {
		t.Errorf("Unexpected JobId. Got: %s, Want: %s", result.OriginJobId, jobID)
	}
	// 检查分配的 CPU
	if result.AllocCores != expectedCores {
		t.Errorf("Unexpected AllocCPUs. Got: %d, Want: %d", result.AllocCores, expectedCores)
	}
	// 检查 ExitCode
	if result.ExitCode != "0:0" {
		t.Errorf("Unexpected ExitCode. Got: %s, Want: 0:0", result.ExitCode)
	}
	// 检查 Priority
	if result.Priority != 1 {
		t.Errorf("Unexpected Priority. Got: %d, Want: 1", result.Priority)
	}
}

func TestGetFreeResource(t *testing.T) {
	commonCfg := config.SchedulerCommon{
		DefaultQueue: "default",
	}
	customCfg := &config.SlurmBackendProvider{
		GetResource: "sinfo -N --Format=StateCompact,CPUS,FreeMem,Memory -p ${queue}",
	}
	p := NewProvider(commonCfg, customCfg, fake_ExecShellCmd)
	result, err := p.GetFreeResource(context.Background(), []string{"default"})
	assert.NoError(t, err, "GetFreeResource should not return an error")
	assert.NotNil(t, result, "Result should not be nil")
	assert.Len(t, result, 1, "Result should contain exactly one queue")

	defaultResource, exists := result["default"]
	assert.True(t, exists, "Result should contain 'default' queue")

	if defaultResource != nil {
		assert.Equal(t, int64(8), defaultResource.Cpu, "Unexpected Cpu value")
		assert.Equal(t, int64(28), defaultResource.TotalCpu, "Unexpected TotalCpu value")
		assert.Equal(t, int64(2), defaultResource.IdleNodeNum, "Unexpected IdleNodeNum value")
		assert.Equal(t, int64(2), defaultResource.AllocNodeNum, "Unexpected AllocNodeNum value")
		assert.Equal(t, int64(5), defaultResource.TotalNodeNum, "Unexpected TotalNodeNum value")
		assert.Equal(t, int64(18600), defaultResource.Memory, "Unexpected Memory value")
		assert.Equal(t, int64(110000), defaultResource.TotalMemory, "Unexpected TotalMemory value")
		assert.True(t, defaultResource.IsDefault, "Default queue should have IsDefault set to true")
	}

	/*	// 验证结果
		expectedResult := &v20230530.Resource{
			Cpu:          8,      // 2 + 6 (只计算 idle 状态的 CPU)
			TotalCpu:     28,     // 2 + 6 + 4 + 8 + 8
			IdleNodeNum:  2,      // 两个 idle 节点
			AllocNodeNum: 2,      // 一个 alloc 节点和一个 mix 节点
			TotalNodeNum: 5,      // 总共 5 个节点
			Memory:       18600,  // 210 + 15390 + 1000 + 2000 (所有节点的 FREE_MEM)
			TotalMemory:  110000, // 15000 + 15000 + 16000 + 32000 + 32000
		}
	*/
}

var (
	useSacct bool
	mu       sync.Mutex
)

// SetUseSacct 设置是否使用 sacct
func SetUseSacct(use bool) {
	mu.Lock()
	defer mu.Unlock()
	useSacct = use
}

func fake_ExecShellCmd(ctx context.Context, cmdStr string, opts ...cmdhelp.Option) (stdOut string, stdErr string, err error) {
	// cmdhelp.ExecShellCmd函数需要真实的slurm环境 执行命令
	// 此处fake函数 用来模拟调slurm命令的结果返回

	if strings.Contains(cmdStr, "scontrol") {
		if useSacct {
			// 没查到的情况
			return "", "scontrol: error: Invalid job id specified", fmt.Errorf("scontrol failed")
		}
		// 正常的 scontrol 输出
		return `JobId=129 JobName=hello.sh
   UserId=yuansuan(1002) GroupId=yuansuan(1002) MCS_label=N/A
   Priority=1 Nice=0 Account=(null) QOS=normal
   JobState=COMPLETED Reason=None Dependency=(null)
   Requeue=1 Restarts=0 BatchFlag=1 Reboot=0 ExitCode=0:0
   RunTime=00:00:01 TimeLimit=UNLIMITED TimeMin=N/A
   SubmitTime=2024-09-14T03:40:55 EligibleTime=2024-09-14T03:40:55
   AccrueTime=2024-09-14T03:40:55
   StartTime=2024-09-14T03:40:55 EndTime=2024-09-14T03:40:56 Deadline=N/A
   SuspendTime=None SecsPreSuspend=0 LastSchedEval=2024-09-14T03:40:55 Scheduler=Main
   Partition=compute AllocNode:Sid=project-root-dev10:1564307
   ReqNodeList=(null) ExcNodeList=(null)
   NodeList=project-root-dev11.novalocal
   BatchHost=project-root-dev11.novalocal
   NumNodes=1 NumCPUs=8 NumTasks=1 CPUs/Task=1 ReqB:S:C:T=0:0:*:*
   ReqTRES=cpu=1,mem=15731M,node=1,billing=1
   AllocTRES=cpu=8,node=1,billing=8
   Socks/Node=* NtasksPerN:B:S:C=0:0:*:* CoreSpec=*
   MinCPUsNode=1 MinMemoryNode=0 MinTmpDiskNode=0
   Features=(null) DelayBoot=00:00:00
   OverSubscribe=NO Contiguous=0 Licenses=(null) Network=(null)
   Command=/home/yuansuan/tmp_job/hello.sh
   WorkDir=/home/yuansuan/tmp_job
   StdErr=/home/yuansuan/tmp_job/slurm-129.out
   StdIn=/dev/null
   StdOut=/home/yuansuan/tmp_job/slurm-129.out`, "", nil
	} else if strings.Contains(cmdStr, "sacct") {
		return `JobID|JobName|AllocCPUS|State|ExitCode|Submit|Start|End|Priority
19411344|wrap|560|RUNNING|0:0|2024-08-16T13:21:57|2024-08-16T13:21:57|Unknown|1
19411344.batch|batch|56|RUNNING|0:0|2024-08-16T13:21:57|2024-08-16T13:21:57|Unknown|
19411344.extern|extern|560|RUNNING|0:0|2024-08-16T13:21:57|2024-08-16T13:21:57|Unknown|
19411344.0|hostname|560|COMPLETED|0:0|2024-08-16T13:21:57|2024-08-16T13:21:57|2024-08-16T13:21:58|`, "", nil
	} else if strings.Contains(cmdStr, "sbatch") {
		// 模拟 sbatch 命令
		jobID := "12345" // 模拟生成的作业 ID
		return fmt.Sprintf("Submitted batch job %s", jobID), "", nil
	} else if strings.Contains(cmdStr, "scancel") {
		// 模拟 scancel 命令
		return "scancel 1", "", nil
	} else if strings.Contains(cmdStr, "sinfo") {
		return `STATE               CPUS                FREE_MEM            MEMORY
idle                2                   210                 15000
idle                6                   15390               15000
alloc               4                   1000                16000
mix                 8                   2000                32000
down                8                   0                   32000`, "", nil
	}
	return "", "", fmt.Errorf("unexpected command: %s", cmdStr)
}
