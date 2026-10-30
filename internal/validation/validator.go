package validation

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	validator *validator.Validate
}

func New() *Validator {
	instance := validator.New()
	instance.RegisterTagNameFunc(func(field reflect.StructField) string {
		name, _, _ := strings.Cut(field.Tag.Get("json"), ",")
		if name == "-" {
			return ""
		}
		return name
	})
	return &Validator{validator: instance}
}

func (v *Validator) Validate(value interface{}) error {
	return v.validator.Struct(value)
}
