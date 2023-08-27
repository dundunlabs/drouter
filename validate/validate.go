package validate

import (
	"github.com/dundunlabs/prenn/validate/validation"
	"github.com/go-playground/validator/v10"
)

func New() *validator.Validate {
	v := validator.New()
	v.RegisterValidation("enum", validation.ValidateEnum)
	return v
}
