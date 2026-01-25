package utils

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Validator wraps the validator instance
var validate *validator.Validate

func init() {
	validate = validator.New()

	// Register custom tag name function to use json tags
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

// ValidateStruct validates a struct and returns field errors
func ValidateStruct(s interface{}) []FieldError {
	var errors []FieldError

	err := validate.Struct(s)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, FieldError{
				Field:   err.Field(),
				Message: getErrorMessage(err),
			})
		}
	}

	return errors
}

// getErrorMessage returns a human-readable error message for a validation error
func getErrorMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		if err.Type().Kind() == reflect.String {
			return "Must be at least " + err.Param() + " characters"
		}
		return "Must be at least " + err.Param()
	case "max":
		if err.Type().Kind() == reflect.String {
			return "Must be at most " + err.Param() + " characters"
		}
		return "Must be at most " + err.Param()
	case "gte":
		return "Must be greater than or equal to " + err.Param()
	case "gt":
		return "Must be greater than " + err.Param()
	case "lte":
		return "Must be less than or equal to " + err.Param()
	case "lt":
		return "Must be less than " + err.Param()
	case "oneof":
		return "Must be one of: " + err.Param()
	case "uuid":
		return "Must be a valid UUID"
	case "alphanum":
		return "Must contain only alphanumeric characters"
	case "numeric":
		return "Must be a numeric value"
	case "url":
		return "Must be a valid URL"
	default:
		return "Invalid value"
	}
}

// Validate is a shorthand for validating and returning errors in gin handler
func Validate(s interface{}) ([]FieldError, bool) {
	errors := ValidateStruct(s)
	return errors, len(errors) == 0
}
