package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type APIResponse map[string]any

func WriteJSON(w http.ResponseWriter, statusCode int, data APIResponse, logger *log.Logger) {
	jsonData, err := json.MarshalIndent(data, "", " ")

	if err != nil {
		logger.Printf("WriteJSON: failed to marshal data: %v. ", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"internal server error during response generation"}`))
		return
	}

	jsonData = append(jsonData, '\n')

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, err = w.Write(jsonData)

	if err != nil {
		logger.Printf("WriteJSON: failed to write jsonData to ResponseWriter: %v. Status code was: %d", err, statusCode)
	}
}

type ValidationErrorResponse struct {
	Errors []string `json:"errors"`
}

func WriteValidationErrorResponse(w http.ResponseWriter, _ int, initialData APIResponse, validationAttemptError error, logger *log.Logger) {
	if fieldErrors, ok := validationAttemptError.(validator.ValidationErrors); ok {

		formattedErrors := ValidationErrorResponse{
			Errors: make([]string, len(fieldErrors)),
		}
		for i, err := range fieldErrors {
			field := err.Field()
			param := err.Param()
			switch err.Tag() {
			case "required":
				formattedErrors.Errors[i] = fmt.Sprintf("%s is a required field", field)
			case "max":
				formattedErrors.Errors[i] = fmt.Sprintf("%s must be a maximum of %s in length", field, param)
			case "min":
				formattedErrors.Errors[i] = fmt.Sprintf("%s must be a minimum of %s in length", field, param)
			case "email":
				formattedErrors.Errors[i] = fmt.Sprintf("%s must be a valid email address", field)
			case "url":
				formattedErrors.Errors[i] = fmt.Sprintf("%s must be a valid URL", field)
			default:
				formattedErrors.Errors[i] = fmt.Sprintf("field %s: validation failed for rule '%s'", field, err.Tag())
			}
		}

		finalData := APIResponse{"error": "validation failed", "details": formattedErrors}
		WriteJSON(w, http.StatusBadRequest, finalData, logger)
	} else {

		logger.Printf("WriteValidationErrorResponse called with a non-validator.ValidationErrors type, or validator internal error: %v", validationAttemptError)

		WriteJSON(w, http.StatusInternalServerError, APIResponse{"error": "an unexpected error occurred during input validation"}, logger)
	}
}
