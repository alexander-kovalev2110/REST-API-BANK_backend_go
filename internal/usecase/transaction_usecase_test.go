package usecase_test

import (
	"context"
	"testing"

	"github.com/alexander-kovalev2110/REST-API-BANK_backend_go/internal/domain"
	"github.com/alexander-kovalev2110/REST-API-BANK_backend_go/internal/usecase"
)

type MockTransactionRepo struct {
	transactions map[int64]*domain.Transaction
	lastID       int64
}

func NewMockTransactionRepo() *MockTransactionRepo {
	return &MockTransactionRepo{
		transactions: make(map[int64]*domain.Transaction),
	}
}

func (m *MockTransactionRepo) Create(ctx context.Context, t *domain.Transaction) error {
	m.lastID++
	t.ID = m.lastID
	m.transactions[t.ID] = t
	return nil
}

func (m *MockTransactionRepo) GetByID(ctx context.Context, id int64) (*domain.Transaction, error) {
	if t, ok := m.transactions[id]; ok {
		return t, nil
	}
	return nil, nil
}

func (m *MockTransactionRepo) GetList(ctx context.Context, filter domain.TransactionFilter) ([]domain.Transaction, int, error) {
	var res []domain.Transaction
	for _, t := range m.transactions {
		if t.CustomerID == filter.CustomerID {
			res = append(res, *t)
		}
	}
	return res, len(res), nil
}

func (m *MockTransactionRepo) Update(ctx context.Context, t *domain.Transaction) error {
	m.transactions[t.ID] = t
	return nil
}

func (m *MockTransactionRepo) Delete(ctx context.Context, id int64) error {
	delete(m.transactions, id)
	return nil
}

func TestTransactionUseCase(t *testing.T) {
	repo := NewMockTransactionRepo()
	uc := usecase.NewTransactionUseCase(repo)
	ctx := context.Background()

	// 1. Create transaction
	createRes, err := uc.Create(ctx, 1, 250.50)
	if err != nil {
		t.Fatalf("Create transaction failed: %v", err)
	}
	if len(createRes.Transactions) == 0 || createRes.Transactions[0].Amount != 250.50 {
		t.Errorf("Expected amount 250.50, got %v", createRes)
	}
	txID := createRes.Transactions[0].ID

	// 2. Get transaction by ID
	tx, err := uc.GetByID(ctx, 1, txID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if len(tx.Transactions) == 0 || tx.Transactions[0].ID != txID {
		t.Errorf("Expected ID %d, got %v", txID, tx)
	}

	// 3. Update transaction
	updated, err := uc.Update(ctx, 1, txID, 300.00)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	if len(updated.Transactions) == 0 || updated.Transactions[0].Amount != 300.00 {
		t.Errorf("Expected updated amount 300.00, got %v", updated)
	}

	// 4. List transactions
	list, err := uc.GetList(ctx, domain.TransactionFilter{CustomerID: 1})
	if err != nil {
		t.Fatalf("GetList failed: %v", err)
	}
	if list.Total != 1 {
		t.Errorf("Expected total 1, got %d", list.Total)
	}

	// 5. Delete transaction
	if err := uc.Delete(ctx, 1, txID); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// 6. Get deleted transaction should fail with ErrTransactionNotFound
	_, err = uc.GetByID(ctx, 1, txID)
	if err != domain.ErrTransactionNotFound {
		t.Errorf("Expected ErrTransactionNotFound after delete, got %v", err)
	}
}
