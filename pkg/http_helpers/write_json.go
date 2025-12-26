package http_helpers

import (
	"encoding/json"
	"net/http"
)

func WriteJson[T any](w http.ResponseWriter, data T, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}

func WriteSuccessJson[T any](w http.ResponseWriter, data T, status int) {
	WriteJson(
		w, SuccessResponse[T]{
			Status: status,
			Data:   data,
		}, status,
	)
}

func WriteErrorJson(w http.ResponseWriter, error string, status int) {
	WriteJson(
		w, ErrorResponse{
			Status: status,
			Error:  error,
		}, status,
	)
}

func WriteValidationErrorsJson(w http.ResponseWriter, errors map[string][]string) {
	WriteJson(
		w, ValidationErrorResponse{
			Status: http.StatusBadRequest,
			Errors: errors,
		}, http.StatusBadRequest,
	)
}
