package models

import (
	// Native packages
	"fmt"

	// 3rd party packages
	"gopkg.in/go-playground/validator.v9"
)

/**
 *	Validates struct using validator, {@link https://github.com/go-playground/validator}.
 *
 *	@param model mixed
 *
 *	@return error, []string
 */
func Validate(model interface{}) (error, []string) {
	var validationErrors []string

	modelValidator := validator.New()
	err := modelValidator.Struct(model)

	if err != nil {
		for _, fieldErr := range err.(validator.ValidationErrors) {
			validationErrorMessage := fmt.Sprintf("Validation failed on '%s' for '%s'", fieldErr.Tag(), fieldErr.StructField())
			validationErrors = append(validationErrors, validationErrorMessage)
		}
	}

	return err, validationErrors
}
