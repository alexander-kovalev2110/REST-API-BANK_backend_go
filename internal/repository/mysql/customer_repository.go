package mysql

import (
	"context"
	"database/sql"
	"errors"

	"github.com/alexander-kovalev2110/REST-API-BANK_backend_go/internal/domain"
)

type CustomerRepository struct {
	db *sql.DB
}

func NewCustomerRepository(db *sql.DB) *CustomerRepository {
	return &CustomerRepository{db: db}
}

func (r *CustomerRepository) FindByName(ctx context.Context, name string) (*domain.Customer, error) {
	query := `SELECT id, name, pw FROM customers WHERE name = ?`
	row := r.db.QueryRowContext(ctx, query, name)

	var customer domain.Customer
	err := row.Scan(&customer.ID, &customer.Name, &customer.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &customer, nil
}

func (r *CustomerRepository) FindByID(ctx context.Context, id int64) (*domain.Customer, error) {
	query := `SELECT id, name, pw FROM customers WHERE id = ?`
	row := r.db.QueryRowContext(ctx, query, id)

	var customer domain.Customer
	err := row.Scan(&customer.ID, &customer.Name, &customer.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &customer, nil
}

func (r *CustomerRepository) Create(ctx context.Context, customer *domain.Customer) error {
	query := `INSERT INTO customers (name, pw) VALUES (?, ?)`
	res, err := r.db.ExecContext(ctx, query, customer.Name, customer.Password)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	customer.ID = id
	return nil
}
