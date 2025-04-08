package validators

import (
	"gopkg.in/go-playground/validator.v9"
	"path/filepath"
)

const (
	absolutePathValidatorName = "absolutePath"
)

func init() {
	registerValidator(absolutePathValidatorName, absolutePath)
}

func absolutePath(fl validator.FieldLevel) bool {
	path := fl.Field().String()
	isAbsolute := filepath.IsAbs(path)
	return isAbsolute
}
