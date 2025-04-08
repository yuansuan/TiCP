package script

import "fmt"

type ExecResult struct {
	ExitCode int
	Stdout   string
	Stderr   string
}

type Runner interface {
	Exec(scriptPath string, asynchronous bool, onFinishedHandlers ...OnExecFinished) (*ExecResult, error)
}

func NewRunner(scriptType string) (Runner, error) {
	switch scriptType {
	case "powershell":
		return newPowershellRunner(), nil
	default:
		return nil, fmt.Errorf("unknow script runner type %s", scriptType)
	}
}
