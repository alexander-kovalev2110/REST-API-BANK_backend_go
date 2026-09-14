package usecase

import (
	"context"
	"time"

	"github.com/alexander-kovalev2110/REST-API-BANK_backend_go/internal/domain"
)

type TransactionUseCase struct {
	repo domain.TransactionRepository
}

func NewTransactionUseCase(repo domain.TransactionRepository) *TransactionUseCase {
	return &TransactionUseCase{repo: repo}
}

func (uc *TransactionUseCase) Create(ctx context.Context, customerID int64, amount float64) (*domain.TransactionListView, error) {
	t := &domain.Transaction{
		CustomerID: customerID,
		Amount:     amount,
		Date:       time.Now().UTC(),
	}

	if err := uc.repo.Create(ctx, t); err != nil {
		return nil, err
	}

	view := domain.TransactionView{
		ID:         t.ID,
		Amount:     t.Amount,
		Date:       t.Date.Format("2006-01-02T15:04:05-07:00"),
		CustomerID: t.CustomerID,
	}

	return &domain.TransactionListView{
		Transactions: []domain.TransactionView{view},
		Total:        1,
	}, nil
}

func (uc *TransactionUseCase) GetList(ctx context.Context, filter domain.TransactionFilter) (*domain.TransactionListView, error) {
	items, total, err := uc.repo.GetList(ctx, filter)
	if err != nil {
		return nil, err
	}

	views := make([]domain.TransactionView, 0, len(items))
	for _, item := range items {
		views = append(views, domain.TransactionView{
			ID:         item.ID,
			Amount:     item.Amount,
			Date:       item.Date.Format("2006-01-02T15:04:05-07:00"),
			CustomerID: item.CustomerID,
		})
	}

	return &domain.TransactionListView{
		Transactions: views,
		Total:        total,
	}, nil
}

func (uc *TransactionUseCase) GetByID(ctx context.Context, customerID int64, transactionID int64) (*domain.TransactionListView, error) {
	t, err := uc.repo.GetByID(ctx, transactionID)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, domain.ErrTransactionNotFound
	}
	if t.CustomerID != customerID {
		return nil, domain.ErrTransactionNotFound // match PHP behavior (returns 404 or domain exception if not found for user)
	}

	view := domain.TransactionView{
		ID:         t.ID,
		Amount:     t.Amount,
		Date:       t.Date.Format("2006-01-02T15:04:05-07:00"),
		CustomerID: t.CustomerID,
	}

	return &domain.TransactionListView{
		Transactions: []domain.TransactionView{view},
		Total:        1,
	}, nil
}

func (uc *TransactionUseCase) Update(ctx context.Context, customerID int64, transactionID int64, amount float64) (*domain.TransactionListView, error) {
	t, err := uc.repo.GetByID(ctx, transactionID)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, domain.ErrTransactionNotFound
	}
	if t.CustomerID != customerID {
		return nil, domain.ErrTransactionNotFound
	}

	t.Amount = amount
	if err := uc.repo.Update(ctx, t); err != nil {
		return nil, err
	}

	view := domain.TransactionView{
		ID:         t.ID,
		Amount:     t.Amount,
		Date:       t.Date.Format("2006-01-02T15:04:05-07:00"),
		CustomerID: t.CustomerID,
	}

	return &domain.TransactionListView{
		Transactions: []domain.TransactionView{view},
		Total:        1,
	}, nil
}

func (uc *TransactionUseCase) Delete(ctx context.Context, customerID int64, transactionID int64) error {
	t, err := uc.repo.GetByID(ctx, transactionID)
	if err != nil {
		return err
	}
	if t == nil {
		return domain.ErrTransactionNotFound
	}
	if t.CustomerID != customerID {
		return domain.ErrTransactionNotFound
	}

	return uc.repo.Delete(ctx, transactionID)
}
