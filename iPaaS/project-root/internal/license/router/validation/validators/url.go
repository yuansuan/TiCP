package validators

import (
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"gopkg.in/go-playground/validator.v9"
	"regexp"
)

const (
	urlValidatorName = "url"
	urlPattern       = "^(http|https):\\/\\/([a-zA-Z0-9-]+\\.)*[a-zA-Z0-9-]+(\\.[a-zA-Z]{2,})?(:\\d+)?(\\/[^\\s]*)?$"
)

func init() {
	registerValidator(urlValidatorName, url)
}

func url(fl validator.FieldLevel) bool {
	url := fl.Field().String()
	match, err := regexp.MatchString(urlPattern, url)
	if err != nil {
		logging.Default().Infof("the regular expression matching fail. url: %s, err: %v", mac, err)
		return false
	}
	return match
}
