package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func validateRequest(writer http.ResponseWriter, req interface{}) bool {
	err := validate.Struct(req)
	if err != nil {
		var errors []string
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, fmt.Sprintf(
				"field '%s' failed on '%s'",
				err.Field(),
				err.Tag(),
			))
		}
		http.Error(writer, fmt.Sprintf("validation failed: %v", errors), http.StatusBadRequest)
		return false
	}
	return true
}
