package script

import (
	"bytes"
	"os/exec"
)

type powershellRunner struct{}

func newPowershellRunner() *powershellRunner {
	return new(powershellRunner)
}

func (r *powershellRunner) Exec(scriptPath string, asynchronous bool, onFinishedHandlers ...OnExecFinished) (*ExecResult, error) {
	if asynchronous {
		go r.exec(scriptPath, onFinishedHandlers...)
		return new(ExecResult), nil
	} else {
		return r.exec(scriptPath, onFinishedHandlers...)
	}
}

type OnExecFinished func(res *ExecResult)

func (r *powershellRunner) exec(scriptPath string, onFinishHandlers ...OnExecFinished) (*ExecResult, error) {
	cmd := exec.Command("powershell.exe", "-File", scriptPath)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	exitCode := 0
	err := cmd.Run()
	if err != nil {
		// 获取退出状态
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode = exitError.ExitCode()
		}
	}

	res := &ExecResult{
		ExitCode: exitCode,
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
	}

	if len(onFinishHandlers) > 0 {
		for _, onFinish := range onFinishHandlers {
			onFinish(res)
		}
	}

	return res, err
}
