package domain

import (
	"context"
	"errors"
	"time"
)

var (
	ErrTransactionNotFound = errors.New("transaction not found")
	ErrForbiddenTransaction = errors.New("access denied to transaction")
)

type Transaction struct {
	ID         int64     `json:"id" db:"id"`
	Amount     float64   `json:"amount" db:"amount"`
	Date       time.Time `json:"date" db:"date"`
	CustomerID int64     `json:"customer_id" db:"customer_id"`
}

type TransactionView struct {
	ID         int64   `json:"transactionId"`
	Amount     float64 `json:"amount"`
	Date       string  `json:"date"`
	CustomerID int64   `json:"customerId"`
}

type TransactionListView struct {
	Transactions []TransactionView `json:"transactions"`
	Total        int               `json:"total"`
}

type TransactionFilter struct {
	CustomerID int64
	Amount     *float64
	Date       *time.Time
	Page       int
	Limit      int
}

type TransactionRepository interface {
	Create(ctx context.Context, transaction *Transaction) error
	GetByID(ctx context.Context, id int64) (*Transaction, error)
	GetList(ctx context.Context, filter TransactionFilter) ([]Transaction, int, error)
	Update(ctx context.Context, transaction *Transaction) error
	Delete(ctx context.Context, id int64) error
}
