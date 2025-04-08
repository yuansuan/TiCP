package validation

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin/binding"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/license/router/validation/validators"
	"gopkg.in/go-playground/validator.v9"
)

func Initialize() error {
	validate, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		errorMessage := "init binging validator error"
		logging.Default().Errorf(errorMessage)
		return errors.New(errorMessage)
	}

	for k, v := range validators.GetValidators() {
		if err := validate.RegisterValidation(k, v); err != nil {
			errorMessage := fmt.Sprintf("register validator error, err: %v", err)
			logging.Default().Errorf(errorMessage)
			return errors.New(errorMessage)
		}
	}
	return nil
}
