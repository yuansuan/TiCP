package validators

import "gopkg.in/go-playground/validator.v9"

var validators = make(map[string]validator.Func)

func GetValidators() map[string]validator.Func {
	return validators
}

func registerValidator(name string, fn validator.Func) {
	validators[name] = fn
}
