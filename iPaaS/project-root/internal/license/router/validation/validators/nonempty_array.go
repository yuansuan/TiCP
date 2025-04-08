package validators

import (
	"gopkg.in/go-playground/validator.v9"
)

const (
	nonemptyArrayValidatorName = "nonemptyArray"
)

func init() {
	registerValidator(nonemptyArrayValidatorName, nonemptyArray)
}

func nonemptyArray(fl validator.FieldLevel) bool {
	return fl.Field().Len() > 0
}
