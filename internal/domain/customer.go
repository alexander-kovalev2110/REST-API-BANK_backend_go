package domain

import (
	"context"
	"errors"
)

var (
	ErrCustomerAlreadyExists = errors.New("customer already exists")
	ErrCustomerNotFound      = errors.New("customer not found")
	ErrInvalidCredentials    = errors.New("invalid credentials")
)

type Customer struct {
	ID       int64  `json:"id" db:"id"`
	Name     string `json:"name" db:"name"`
	Password string `json:"-" db:"pw"`
}

type CustomerRepository interface {
	FindByName(ctx context.Context, name string) (*Customer, error)
	FindByID(ctx context.Context, id int64) (*Customer, error)
	Create(ctx context.Context, customer *Customer) error
}
