package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/alexander-kovalev2110/REST-API-BANK_backend_go/internal/domain"
	"github.com/alexander-kovalev2110/REST-API-BANK_backend_go/internal/middleware"
	"github.com/alexander-kovalev2110/REST-API-BANK_backend_go/internal/usecase"
	"github.com/go-chi/chi/v5"
)

type TransactionHandler struct {
	transactionUC *usecase.TransactionUseCase
}

func NewTransactionHandler(transactionUC *usecase.TransactionUseCase) *TransactionHandler {
	return &TransactionHandler{transactionUC: transactionUC}
}

type AmountTransactionRequest struct {
	Amount float64 `json:"amount"`
}

func (h *TransactionHandler) Create(w http.ResponseWriter, r *http.Request) {
	customerID, ok := middleware.GetCustomerID(r.Context())
	if !ok {
		middleware.WriteJSON(w, http.StatusUnauthorized, middleware.ErrorResponse{Error: "Unauthorized"})
		return
	}

	var req AmountTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.WriteValidationError(w, map[string]string{"amount": "invalid numeric amount"})
		return
	}

	if req.Amount <= 0 {
		middleware.WriteValidationError(w, map[string]string{"amount": "Amount must be greater than 0"})
		return
	}

	res, err := h.transactionUC.Create(r.Context(), customerID, req.Amount)
	if err != nil {
		middleware.WriteError(w, err)
		return
	}

	middleware.WriteJSON(w, http.StatusCreated, res)
}

func (h *TransactionHandler) GetList(w http.ResponseWriter, r *http.Request) {
	customerID, ok := middleware.GetCustomerID(r.Context())
	if !ok {
		middleware.WriteJSON(w, http.StatusUnauthorized, middleware.ErrorResponse{Error: "Unauthorized"})
		return
	}

	query := r.URL.Query()
	filter := domain.TransactionFilter{
		CustomerID: customerID,
		Page:       1,
		Limit:      10,
	}

	if p := query.Get("page"); p != "" {
		if page, err := strconv.Atoi(p); err == nil && page > 0 {
			filter.Page = page
		}
	}

	if l := query.Get("limit"); l != "" {
		if limit, err := strconv.Atoi(l); err == nil && limit > 0 {
			filter.Limit = limit
		}
	}

	if a := query.Get("amount"); a != "" {
		if amount, err := strconv.ParseFloat(a, 64); err == nil {
			filter.Amount = &amount
		}
	}

	if d := query.Get("date"); d != "" {
		if parsedDate, err := time.Parse("2006-01-02", d); err == nil {
			filter.Date = &parsedDate
		}
	}

	res, err := h.transactionUC.GetList(r.Context(), filter)
	if err != nil {
		middleware.WriteError(w, err)
		return
	}

	middleware.WriteJSON(w, http.StatusOK, res)
}

func (h *TransactionHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	customerID, ok := middleware.GetCustomerID(r.Context())
	if !ok {
		middleware.WriteJSON(w, http.StatusUnauthorized, middleware.ErrorResponse{Error: "Unauthorized"})
		return
	}

	transactionIDParam := chi.URLParam(r, "transactionId")
	transactionID, err := strconv.ParseInt(transactionIDParam, 10, 64)
	if err != nil {
		middleware.WriteValidationError(w, map[string]string{"transactionId": "invalid transaction ID"})
		return
	}

	res, err := h.transactionUC.GetByID(r.Context(), customerID, transactionID)
	if err != nil {
		middleware.WriteError(w, err)
		return
	}

	middleware.WriteJSON(w, http.StatusOK, res)
}

func (h *TransactionHandler) Update(w http.ResponseWriter, r *http.Request) {
	customerID, ok := middleware.GetCustomerID(r.Context())
	if !ok {
		middleware.WriteJSON(w, http.StatusUnauthorized, middleware.ErrorResponse{Error: "Unauthorized"})
		return
	}

	transactionIDParam := chi.URLParam(r, "transactionId")
	transactionID, err := strconv.ParseInt(transactionIDParam, 10, 64)
	if err != nil {
		middleware.WriteValidationError(w, map[string]string{"transactionId": "invalid transaction ID"})
		return
	}

	var req AmountTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.WriteValidationError(w, map[string]string{"amount": "invalid numeric amount"})
		return
	}

	if req.Amount <= 0 {
		middleware.WriteValidationError(w, map[string]string{"amount": "Amount must be greater than 0"})
		return
	}

	res, err := h.transactionUC.Update(r.Context(), customerID, transactionID, req.Amount)
	if err != nil {
		middleware.WriteError(w, err)
		return
	}

	middleware.WriteJSON(w, http.StatusOK, res)
}

func (h *TransactionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	customerID, ok := middleware.GetCustomerID(r.Context())
	if !ok {
		middleware.WriteJSON(w, http.StatusUnauthorized, middleware.ErrorResponse{Error: "Unauthorized"})
		return
	}

	transactionIDParam := chi.URLParam(r, "transactionId")
	transactionID, err := strconv.ParseInt(transactionIDParam, 10, 64)
	if err != nil {
		middleware.WriteValidationError(w, map[string]string{"transactionId": "invalid transaction ID"})
		return
	}

	if err := h.transactionUC.Delete(r.Context(), customerID, transactionID); err != nil {
		middleware.WriteError(w, err)
		return
	}

	middleware.WriteJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}
