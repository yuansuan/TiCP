package validation

import (
	"github.com/gin-gonic/gin/binding"
	"gopkg.in/go-playground/validator.v9"
)

// Validate Validate
var Validate = validator.New()

func init() {
	binding.Validator = new(defaultValidator)
}
