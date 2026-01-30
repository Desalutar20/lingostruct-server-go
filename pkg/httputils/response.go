package httputils

import (
	"encoding/json"
	"net/http"
)

type successResponse[T any] struct {
	Data T `json:"data"`
}

type errorResponse struct {
	Error string `json:"error"`
}

type validationErrorResponse struct {
	Errors map[string]string `json:"errors"`
}

func writeResponse(w http.ResponseWriter, data any, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "something went wrong", http.StatusInternalServerError)
	}
}

func SuccessResponse[T any](w http.ResponseWriter, data T, statusCode int) {
	response := successResponse[T]{Data: data}

	writeResponse(w, response, statusCode)
}

func ErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	response := errorResponse{Error: message}

	writeResponse(w, response, statusCode)
}

func ValidationErrorResponse(w http.ResponseWriter, errors map[string]string) {
	response := validationErrorResponse{Errors: errors}

	writeResponse(w, response, http.StatusBadRequest)
}
