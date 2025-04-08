package pbspro

import (
	"context"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/config"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/cmdhelp"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/dao/models"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"

	v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/backend/job"
)

func TestSubmit(t *testing.T) {
	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "pbspro_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	commonCfg := config.SchedulerCommon{
		DefaultQueue: "workq",
		Workspace:    tempDir,
	}
	customCfg := &config.PbsProBackendProvider{
		Submit: `qsub -o "${out}" -e "${err}" -q "${queue}" -l select="${nodes}":ncpus="${number_of_cpu}":mem="${memory_mb}"mb "${script}"`,
	}

	p := NewProvider(commonCfg, customCfg, fake_ExecShellCmd)

	t.Run("ValidSubmission", func(t *testing.T) {
		j := &job.Job{
			Job: &models.Job{
				Command:       "echo 'Hello, World!'",
				Workspace:     filepath.Join(tempDir, "job-1"),
				RequestCores:  2,
				RequestMemory: 100,
			},
		}

		jobID, err := p.Submit(context.Background(), j)

		assert.NoError(t, err)
		assert.Equal(t, "45", jobID)
	})
}

func TestKill(t *testing.T) {
	j := &job.Job{
		Job: &models.Job{
			OriginJobId: "45",
		},
	}

	conf := &config.PbsProBackendProvider{
		Kill: "qdel -x ${job_id}",
	}

	p := NewProvider(config.SchedulerCommon{}, conf, fake_ExecShellCmd)

	err := p.Kill(context.TODO(), j)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetFreeResource(t *testing.T) {
	commonCfg := config.SchedulerCommon{
		DefaultQueue: "workq",
	}
	customCfg := &config.PbsProBackendProvider{
		GetResource: "pbsnodes -av",
	}
	p := NewProvider(commonCfg, customCfg, fake_ExecShellCmd)
	resources, err := p.GetFreeResource(context.Background(), []string{"workq"})
	if err != nil {
		t.Fatalf("GetFreeResource failed: %v", err)
	}
	// 验证结果
	expectedTotalCPU := int64(16)       // 8 + 8
	expectedTotalMem := int64(32530712) // 16265356 * 2

	if len(resources) != 1 {
		t.Fatalf("Expected 1 queue, got %d", len(resources))
	}
	workqResource, ok := resources["workq"]
	if !ok {
		t.Fatalf("Expected 'workq' queue, not found")
	}
	if workqResource.Cpu != expectedTotalCPU {
		t.Errorf("Expected %d CPUs, got %d", expectedTotalCPU, workqResource.Cpu)
	}
	if workqResource.Memory != expectedTotalMem {
		t.Errorf("Expected %d GB memory, got %d", expectedTotalMem, workqResource.Memory)
	}
	if !workqResource.IsDefault {
		t.Errorf("Expected 'workq' to be the default queue")
	}
}

func TestParserResourceInfo(t *testing.T) {
	f, err := ioutil.ReadFile("resourceInfoTest.txt")
	if err != nil {
		fmt.Println("read fail", err)
	}
	var mem int64
	var cpu int64
	lines := strings.Split(strings.TrimSpace(string(f)), "\n")
	for _, line := range lines {
		if strings.Contains(line, "resources_available.mem") {
			memNumLine := strings.Split(strings.TrimSpace(line), "=")
			if len(memNumLine) == 2 {
				//去除末位的kb
				len := len(memNumLine[1])
				memNum, _ := strconv.Atoi(strings.TrimSpace(memNumLine[1][:len-2]))
				mem += int64(memNum)
			} else {
				_ = errors.New(fmt.Sprintf("UnexpectedResourceInfoLine, info: %s", memNumLine))
				break
			}

		}
		if strings.Contains(line, "resources_available.ncpus") {
			cpuNumLine := strings.Split(strings.TrimSpace(line), "=")
			if len(cpuNumLine) == 2 {
				cpuNum, _ := strconv.Atoi(strings.TrimSpace(cpuNumLine[1]))
				cpu += int64(cpuNum)
			} else {
				_ = errors.New(fmt.Sprintf("UnexpectedResourceInfoLine, info: %s", cpuNumLine))
				break
			}

		}
	}
	println(mem / 1024)
	println(cpu)
}

//	mem       ncpus   nmics   ngpus
//
// vnode           state           njobs   run   susp      f/t        f/t     f/t     f/t   jobs
// --------------- --------------- ------ ----- ------ ------------ ------- ------- ------- -------
// pbspro          free                 0     0      0      8gb/8gb     4/4     0/0     0/0 --
// pbspro2         free                 0     0      0      2gb/4gb     2/4     0/0     0/0 --
func TestParserResourceInfo2(t *testing.T) {
	str :=
		"                                                       mem       ncpus   nmics   ngpus\n" +
			"vnode           state           njobs   run   susp      f/t        f/t     f/t     f/t   jobs\n" +
			"--------------- --------------- ------ ----- ------ ------------ ------- ------- ------- -------\n" +
			"pbspro          free                 0     0      0      8gb/4gb     4/4     0/0     0/0 --\n" +
			"pbspro2         free                 0     0      0      2gb/4gb     3/4     0/0     0/0 --"
	var mem int64
	var cpu int64
	regexStr := "(\\w+)\\s+(\\w+)\\s+(\\d+)\\s+(\\d+)\\s+(\\d+)\\s+(\\w+\\/\\w+)\\s+(\\w+\\/\\w+)\\s+(\\w+\\/+\\w+)\\s+(\\w+\\/\\w+)\\s+"
	r, _ := regexp.Compile(regexStr)
	sList := strings.Split(strings.TrimSpace(str), "\n")
	if 4 > len(sList) {
		err := errors.New(fmt.Sprintf("UnexpectedResourceInfo, %s", sList))
		println(err.Error())
		return
	}
	for _, line := range sList[3:] {
		matchs := r.FindStringSubmatch(strings.TrimSpace(line))
		if len(matchs) == 10 && matchs[2] == "free" {
			//去掉单位gb
			memString := strings.Split(matchs[6], "/")
			cpuString := strings.Split(matchs[7], "/")
			len := len(memString[0])
			memNum, _ := strconv.Atoi(memString[0][:len-2])
			cpuNum, _ := strconv.Atoi(cpuString[0])
			cpu += int64(cpuNum)
			mem += int64(memNum)
		}
	}
	println(mem)
	println(cpu)
}
func TestParse(t *testing.T) {

	file, _ := os.Open("resourceInfoTest.txt")

	defer file.Close()

	expected := map[string]*v20230530.Resource{
		"workq":   &v20230530.Resource{Memory: 700000, Cpu: 5},
		"test123": &v20230530.Resource{Memory: 200000, Cpu: 4},
	}

	stdout, _ := io.ReadAll(file)
	result := parseResource(string(stdout), "workq")

	if len(result) != len(expected) {
		t.Errorf("Unexpected number of queues. Expected: %d, Got: %d", len(expected), len(result))
	}

	for queueName, expectedResource := range expected {
		resultResource, ok := result[queueName]
		if !ok {
			t.Errorf("Queue '%s' not found in the result", queueName)
			continue
		}

		if resultResource.Memory != expectedResource.Memory {
			t.Errorf("Unexpected memory for queue '%s'. Expected: %d, Got: %d", queueName, expectedResource.Memory, resultResource.Memory)
		}

		if resultResource.Cpu != expectedResource.Cpu {
			t.Errorf("Unexpected CPU for queue '%s'. Expected: %d, Got: %d", queueName, expectedResource.Cpu, resultResource.Cpu)
		}
	}
}

func fake_ExecShellCmd(ctx context.Context, cmdStr string, opts ...cmdhelp.Option) (stdOut string, stdErr string, err error) {
	if strings.Contains(cmdStr, "qsub") {
		// 模拟 qsub 命令
		return "45.dev8", "", nil
	} else if strings.Contains(cmdStr, "qdel") {
		// 模拟 qdel 命令
		return "", "", nil
	} else if strings.Contains(cmdStr, "qstat") {
		// 模拟 qstat 命令
		return `Job Id: 45.dev8
    Job_Name = hello.sh
    Job_Owner = yuansuan@project-root-dev8
    resources_used.cpupercent = 0
    resources_used.cput = 00:00:00
    resources_used.mem = 0kb
    resources_used.ncpus = 1
    resources_used.vmem = 0kb
    resources_used.walltime = 00:00:00
    job_state = F
    queue = workq
    server = dev8
    Checkpoint = u
    ctime = Sat Sep 14 08:37:45 2024
    Error_Path = project-root-dev8:/home/yuansuan/tmp_job/hello.sh.e45
    exec_host = project-root-dev8/0
    exec_vnode = (project-root-dev8:ncpus=1)
    Hold_Types = n
    Join_Path = n
    Keep_Files = n
    Mail_Points = a
    mtime = Sat Sep 14 08:37:55 2024
    Output_Path = project-root-dev8:/home/yuansuan/tmp_job/hello.sh.o45
    Priority = 0
    qtime = Sat Sep 14 08:37:45 2024
    Rerunable = True
    Resource_List.ncpus = 1
    Resource_List.nodect = 1
    Resource_List.place = pack
    Resource_List.select = 1:ncpus=1
    stime = Sat Sep 14 08:37:55 2024
    session_id = 14948
    jobdir = /home/yuansuan
    substate = 92
    Variable_List = PBS_O_HOME=/home/yuansuan,PBS_O_LANG=en_US.UTF-8,
        PBS_O_LOGNAME=yuansuan,
        PBS_O_PATH=/shared/singularity/bin/:/usr/local/bin:/bin:/usr/bin:/usr/
        local/sbin:/usr/sbin:/opt/pbs/bin:/home/yuansuan/.local/bin:/home/yuans
        uan/bin,PBS_O_MAIL=/var/spool/mail/yuansuan,PBS_O_SHELL=/bin/bash,
        PBS_O_WORKDIR=/home/yuansuan/tmp_job,PBS_O_SYSTEM=Linux,
        PBS_O_QUEUE=workq,PBS_O_HOST=project-root-dev8
    comment = Job run at Sat Sep 14 at 08:37 on (project-root-dev8:ncpus=1) and
         finished
    etime = Sat Sep 14 08:37:45 2024
    run_count = 1
    Stageout_status = 1
    Exit_status = 0
    Submit_arguments = hello.sh
    history_timestamp = 1726303075
    project = _pbs_project_default
    Submit_Host = project-root-dev8`, "", nil
	} else if strings.Contains(cmdStr, "pbsnodes") {
		// 模拟 pbsnodes 命令
		return `project-root-dev8
     Mom = project-root-dev8
     ntype = PBS
     state = free
     pcpus = 8
     resources_available.arch = linux
     resources_available.host = project-root-dev8
     resources_available.mem = 16265356kb
     resources_available.ncpus = 8
     resources_available.vnode = project-root-dev8
     resources_assigned.accelerator_memory = 0kb
     resources_assigned.hbmem = 0kb
     resources_assigned.mem = 0kb
     resources_assigned.naccelerators = 0
     resources_assigned.ncpus = 0
     resources_assigned.vmem = 0kb
     resv_enable = True
     sharing = default_shared
     last_state_change_time = Tue Aug 27 12:02:25 2024
     last_used_time = Sat Sep 14 08:37:55 2024

node1
     Mom = project-root-dev9
     ntype = PBS
     state = free
     pcpus = 8
     resources_available.arch = linux
     resources_available.host = project-root-dev9
     resources_available.mem = 16265356kb
     resources_available.ncpus = 8
     resources_available.vnode = node1
     resources_assigned.accelerator_memory = 0kb
     resources_assigned.hbmem = 0kb
     resources_assigned.mem = 0kb
     resources_assigned.naccelerators = 0
     resources_assigned.ncpus = 0
     resources_assigned.vmem = 0kb
     resv_enable = True
     sharing = default_shared
     last_state_change_time = Fri Aug 16 10:22:03 2024
     last_used_time = Fri Aug 16 10:22:03 2024

master
     Mom = project-root-dev8
     ntype = PBS
     state = free
     resources_available.host = project-root-dev8
     resources_available.vnode = master
     resources_assigned.accelerator_memory = 0kb
     resources_assigned.hbmem = 0kb
     resources_assigned.mem = 0kb
     resources_assigned.naccelerators = 0
     resources_assigned.ncpus = 0
     resources_assigned.vmem = 0kb
     resv_enable = True
     sharing = default_shared
     last_state_change_time = Mon Aug  5 07:01:13 2024`, "", nil
	}
	return "", "", fmt.Errorf("unexpected command: %s", cmdStr)
}
