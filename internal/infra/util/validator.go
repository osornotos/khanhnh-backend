package util

import (
	"github.com/go-playground/validator/v10"
	"reflect"
	"strings"
)

var Validator *validator.Validate

func init() {
	Validator = NewValidator()
}

func NewValidator() *validator.Validate {
	var v = validator.New()

	v.RegisterTagNameFunc(extractJsonTag)

	return v
}

// with `json:"name" validate:"required"`
// err = err.(validator.ValidationErrors)
// err.Field() return "Name" instead of "name"
func extractJsonTag(fld reflect.StructField) string {
	name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
	if name == "-" {
		return ""
	}
	return name
}
