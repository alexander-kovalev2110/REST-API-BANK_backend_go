package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/alexander-kovalev2110/REST-API-BANK_backend_go/internal/middleware"
	"github.com/alexander-kovalev2110/REST-API-BANK_backend_go/internal/usecase"
)

type AuthHandler struct {
	customerUC *usecase.CustomerUseCase
}

func NewAuthHandler(customerUC *usecase.CustomerUseCase) *AuthHandler {
	return &AuthHandler{customerUC: customerUC}
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.WriteValidationError(w, map[string]string{"body": "invalid JSON payload"})
		return
	}

	errors := make(map[string]string)
	if strings.TrimSpace(req.Name) == "" {
		errors["name"] = "This value should not be blank."
	}
	if strings.TrimSpace(req.Password) == "" {
		errors["password"] = "This value should not be blank."
	}
	if len(errors) > 0 {
		middleware.WriteValidationError(w, errors)
		return
	}

	res, err := h.customerUC.Register(r.Context(), req.Name, req.Password)
	if err != nil {
		middleware.WriteError(w, err)
		return
	}

	middleware.WriteJSON(w, http.StatusCreated, res)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.WriteValidationError(w, map[string]string{"body": "invalid JSON payload"})
		return
	}

	errors := make(map[string]string)
	if strings.TrimSpace(req.Name) == "" {
		errors["name"] = "This value should not be blank."
	}
	if strings.TrimSpace(req.Password) == "" {
		errors["password"] = "This value should not be blank."
	}
	if len(errors) > 0 {
		middleware.WriteValidationError(w, errors)
		return
	}

	res, err := h.customerUC.Login(r.Context(), req.Name, req.Password)
	if err != nil {
		middleware.WriteError(w, err)
		return
	}

	middleware.WriteJSON(w, http.StatusOK, res)
}
