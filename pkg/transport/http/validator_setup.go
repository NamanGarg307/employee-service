package http

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func init() {
	Validate = validator.New()
	err := Validate.RegisterValidation("trimspace", trimSpaceValidator)
	if err != nil {
		// Slack Alert
		return
	}
}

func trimSpaceValidator(fl validator.FieldLevel) bool {
	if fl.Field().String() == "" {
		return true
	}
	return strings.TrimSpace(fl.Field().String()) != ""
}
