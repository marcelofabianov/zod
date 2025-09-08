package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	validate *validator.Validate
}

func New() *Validator {
	return &Validator{
		validate: validator.New(),
	}
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (v *Validator) Validate(data any) []ValidationError {
	var validationErrors []ValidationError

	errs := v.validate.Struct(data)
	if errs == nil {
		return nil
	}

	for _, err := range errs.(validator.ValidationErrors) {
		field := strings.ToLower(err.Field())
		message := fmt.Sprintf("field validation for '%s' failed on the '%s' tag", field, err.Tag())
		validationErrors = append(validationErrors, ValidationError{Field: field, Message: message})
	}

	return validationErrors
}
