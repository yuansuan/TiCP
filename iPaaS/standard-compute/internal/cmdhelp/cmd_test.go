package cmdhelp

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecShellCmd(t *testing.T) {
	// test ok
	stdout, stderr, err := ExecShellCmd(context.TODO(), "echo 1234abc")
	if stdout != "1234abc\n" || err != nil || len(stderr) != 0 {
		t.Log(stdout, stderr)
		t.Fatal("BadResult")
	}
	// test fail
	stdout, stderr, err = ExecShellCmd(context.TODO(), "badshellcmd")
	if err == nil {
		t.Log(stdout, stderr)
		t.Fatal("NeedError")
	}
	// test env
	stdout, stderr, err = ExecShellCmd(context.TODO(), "unset SC_TEST_ENV_VAIL")
	if err != nil {
		t.Log(stdout, stderr)
		t.Fatal("BadError")
	}
	stdout, stderr, err = ExecShellCmd(context.TODO(), "env | grep SC_TEST_ENV_VAIL", WithCmdEnv([]string{"SC_TEST_ENV_VAIL=1122334455"}))
	if err != nil {
		t.Log(stdout, stderr)
		t.Fatal("BadError")
	}
	if stdout != "SC_TEST_ENV_VAIL=1122334455\n" {
		t.Log(stdout, stderr)
		t.Fatal(stdout)
	}
	ExecShellCmd(context.TODO(), "unset SC_TEST_ENV_VAIL")

	// test dir
	pwd, _ := os.Getwd()
	stdout, _, err = ExecShellCmd(context.TODO(), "pwd")
	if pwd != strings.TrimSpace(stdout) && err != nil {
		t.Fatalf(stdout, pwd)
	}

	destDir := filepath.Dir(pwd)
	stdout, _, err = ExecShellCmd(context.TODO(), "pwd", WithCmdDir(destDir))
	if destDir != strings.TrimSpace(stdout) && err != nil {
		t.Fatalf(destDir, stdout, err)
	}
	os.Remove(destDir)
}

func TestEnsureSubmitCommandSlurm(t *testing.T) {
	rawCmd := "sbatch -D \"${cwd}\" -o \"${out}\" -e \"${err}\" --nodes \"${nodes}\" --ntasks-per-node \"${ntasks_per_node}\" --mem \"${memory_mb}\" -p \"${queue}\" --wrap \"/bin/bash ${script}\""

	cmd, err := EnsureSubmitCommand(rawCmd, map[string]string{
		"-p": "test1",
	}, "--wrap \"/bin/bash ${script}\"")
	assert.NoError(t, err)
	t.Log(cmd)
	assert.Contains(t, cmd, "test1")
	assert.NotContains(t, cmd, "\"${queue}\"")

	rawCmd = "/opt/pbs/bin/qsub -o \"${out}\" -e \"${err}\" -q \"${queue}\" -l select=\"${nodes}\":ncpus=\"${number_of_cpu}\":mem=\"${memory_mb}\"mb \"${script}\""

	cmd, err = EnsureSubmitCommand(rawCmd, map[string]string{
		"-l": "select=1:ncpus=2:mem",
	}, "\"${script}\"")
	assert.NoError(t, err)
	t.Log(cmd)
	assert.Contains(t, cmd, `select=1:ncpus=2:mem`)
	assert.NotContains(t, cmd, "select=\"${nodes}\":ncpus=\"${number_of_cpu}\":mem=\"${memory_mb}\"mb")
}
