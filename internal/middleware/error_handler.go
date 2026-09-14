package middleware

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/alexander-kovalev2110/REST-API-BANK_backend_go/internal/domain"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type ValidationErrorResponse struct {
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors"`
}

func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func WriteError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrCustomerAlreadyExists):
		WriteJSON(w, http.StatusConflict, ErrorResponse{Error: err.Error()})
	case errors.Is(err, domain.ErrCustomerNotFound):
		WriteJSON(w, http.StatusNotFound, ErrorResponse{Error: err.Error()})
	case errors.Is(err, domain.ErrInvalidCredentials):
		WriteJSON(w, http.StatusUnauthorized, ErrorResponse{Error: err.Error()})
	case errors.Is(err, domain.ErrTransactionNotFound):
		WriteJSON(w, http.StatusNotFound, ErrorResponse{Error: err.Error()})
	case errors.Is(err, domain.ErrForbiddenTransaction):
		WriteJSON(w, http.StatusForbidden, ErrorResponse{Error: err.Error()})
	default:
		WriteJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "internal server error"})
	}
}

func WriteValidationError(w http.ResponseWriter, fieldErrors map[string]string) {
	WriteJSON(w, http.StatusBadRequest, ValidationErrorResponse{
		Message: "Validation failed",
		Errors:  fieldErrors,
	})
}
