package validators

import (
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"gopkg.in/go-playground/validator.v9"
	"regexp"
)

const (
	macValidatorName = "mac"
	macPattern       = "^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$"
)

func init() {
	registerValidator(macValidatorName, mac)
}

func mac(fl validator.FieldLevel) bool {
	mac := fl.Field().String()
	match, err := regexp.MatchString(macPattern, mac)
	if err != nil {
		logging.Default().Infof("the regular expression matching fail. mac: %s, err: %v", mac, err)
		return false
	}
	return match
}
