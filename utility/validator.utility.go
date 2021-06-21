package utility

import "github.com/go-playground/validator"

var validate *validator.Validate

func ValidateStruct(s interface{}) error {
	if validate == nil {
		validate = validator.New()
	}
	return validate.Struct(s)
}
