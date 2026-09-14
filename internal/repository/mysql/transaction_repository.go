package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/alexander-kovalev2110/REST-API-BANK_backend_go/internal/domain"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Create(ctx context.Context, t *domain.Transaction) error {
	query := `INSERT INTO transactions (customer_id, amount, date) VALUES (?, ?, ?)`
	res, err := r.db.ExecContext(ctx, query, t.CustomerID, t.Amount, t.Date)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	t.ID = id
	return nil
}

func (r *TransactionRepository) GetByID(ctx context.Context, id int64) (*domain.Transaction, error) {
	query := `SELECT id, customer_id, amount, date FROM transactions WHERE id = ?`
	row := r.db.QueryRowContext(ctx, query, id)

	var t domain.Transaction
	err := row.Scan(&t.ID, &t.CustomerID, &t.Amount, &t.Date)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &t, nil
}

func (r *TransactionRepository) GetList(ctx context.Context, filter domain.TransactionFilter) ([]domain.Transaction, int, error) {
	whereClause := "WHERE customer_id = ?"
	args := []interface{}{filter.CustomerID}

	if filter.Amount != nil {
		whereClause += " AND amount = ?"
		args = append(args, *filter.Amount)
	}

	if filter.Date != nil {
		year, month, day := filter.Date.Date()
		start := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
		end := time.Date(year, month, day, 23, 59, 59, 0, time.UTC)
		whereClause += " AND date BETWEEN ? AND ?"
		args = append(args, start, end)
	}

	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM transactions %s", whereClause)
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 10
	}
	offset := (filter.Page - 1) * filter.Limit

	query := fmt.Sprintf("SELECT id, customer_id, amount, date FROM transactions %s ORDER BY date ASC LIMIT ? OFFSET ?", whereClause)
	queryArgs := append(args, filter.Limit, offset)

	rows, err := r.db.QueryContext(ctx, query, queryArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var transactions []domain.Transaction
	for rows.Next() {
		var t domain.Transaction
		if err := rows.Scan(&t.ID, &t.CustomerID, &t.Amount, &t.Date); err != nil {
			return nil, 0, err
		}
		transactions = append(transactions, t)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}

func (r *TransactionRepository) Update(ctx context.Context, t *domain.Transaction) error {
	query := `UPDATE transactions SET amount = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, t.Amount, t.ID)
	return err
}

func (r *TransactionRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM transactions WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
