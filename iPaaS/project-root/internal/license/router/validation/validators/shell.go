package validators

import (
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"gopkg.in/go-playground/validator.v9"
	"os/exec"
)

const (
	shellValidatorName = "shell"
)

func init() {
	registerValidator(shellValidatorName, shell)
}

func shell(fl validator.FieldLevel) bool {
	shell := fl.Field().String()

	cmd := exec.Command("/bin/sh", "-c", shell)
	_, err := cmd.CombinedOutput()
	if err != nil {
		logging.Default().Infof("this is not a correct shell: %s, err: %v", shell, err)
		return false
	}

	return true
}
