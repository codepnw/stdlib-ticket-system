package utils

import "github.com/go-playground/validator/v10"

var v = validator.New()

func Validate(data any) error {
	return v.Struct(data)
}