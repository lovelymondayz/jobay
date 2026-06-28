package utils

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

// FormatValidationError mengubah error validator menjadi pesan kalimat
func FormatValidationError(err error) string {
	var errors []string

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validationErrors {
			switch fieldError.Tag() {
			case "required":
				errors = append(errors, fieldError.Field()+" is required")
			case "email":
				errors = append(errors, fieldError.Field()+" must be a valid email")
			case "min":
				errors = append(errors, fieldError.Field()+" must be at least "+fieldError.Param()+" characters")
			case "max":
				errors = append(errors, fieldError.Field()+" must be at most "+fieldError.Param()+" characters")
			default:
				errors = append(errors, fieldError.Field()+" is invalid")
			}
		}
	}
	return strings.Join(errors, ", ")
}
