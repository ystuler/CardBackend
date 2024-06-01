package util

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"strings"
)

type Validator struct {
	validate *validator.Validate
}

func NewValidator() *Validator {
	return &Validator{
		validate: validator.New(),
	}
}

func (v *Validator) Validate(i interface{}) error {
	return v.validate.Struct(i)
}

func (v *Validator) ValidateWithDetailedErrors(i interface{}) error {
	err := v.validate.Struct(i)
	if err == nil {
		return nil
	}

	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		var errorMessages []string
		for _, validationError := range validationErrors {
			errorMessages = append(errorMessages, fmt.Sprintf("Field '%s' %s", validationError.Field(), validationError.ActualTag()))
		}
		return fmt.Errorf("invalid input: %s", strings.Join(errorMessages, ", "))
	}

	return err
}
