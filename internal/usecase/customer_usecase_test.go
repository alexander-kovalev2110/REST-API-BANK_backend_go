package usecase_test

import (
	"context"
	"testing"

	"github.com/alexander-kovalev2110/REST-API-BANK_backend_go/internal/domain"
	"github.com/alexander-kovalev2110/REST-API-BANK_backend_go/internal/usecase"
)

type MockCustomerRepo struct {
	customers map[string]*domain.Customer
	lastID    int64
}

func NewMockCustomerRepo() *MockCustomerRepo {
	return &MockCustomerRepo{
		customers: make(map[string]*domain.Customer),
	}
}

func (m *MockCustomerRepo) FindByName(ctx context.Context, name string) (*domain.Customer, error) {
	if c, ok := m.customers[name]; ok {
		return c, nil
	}
	return nil, nil
}

func (m *MockCustomerRepo) FindByID(ctx context.Context, id int64) (*domain.Customer, error) {
	for _, c := range m.customers {
		if c.ID == id {
			return c, nil
		}
	}
	return nil, nil
}

func (m *MockCustomerRepo) Create(ctx context.Context, customer *domain.Customer) error {
	m.lastID++
	customer.ID = m.lastID
	m.customers[customer.Name] = customer
	return nil
}

func TestRegisterAndLogin(t *testing.T) {
	repo := NewMockCustomerRepo()
	uc := usecase.NewCustomerUseCase(repo, "test-secret")
	ctx := context.Background()

	// 1. Register new customer
	res, err := uc.Register(ctx, "john_doe", "secret123")
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}
	if res.Token == "" {
		t.Errorf("Expected token, got empty string")
	}

	// 2. Register duplicate customer should fail
	_, err = uc.Register(ctx, "john_doe", "secret123")
	if err != domain.ErrCustomerAlreadyExists {
		t.Errorf("Expected ErrCustomerAlreadyExists, got %v", err)
	}

	// 3. Login with correct password
	loginRes, err := uc.Login(ctx, "john_doe", "secret123")
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}
	if loginRes.Token == "" {
		t.Errorf("Expected token on login, got empty string")
	}

	// 4. Login with invalid password
	_, err = uc.Login(ctx, "john_doe", "wrongpassword")
	if err != domain.ErrInvalidCredentials {
		t.Errorf("Expected ErrInvalidCredentials, got %v", err)
	}
}
