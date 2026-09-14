package usecase

import (
	"context"
	"time"

	"github.com/alexander-kovalev2110/REST-API-BANK_backend_go/internal/domain"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type CustomerUseCase struct {
	repo      domain.CustomerRepository
	jwtSecret string
}

func NewCustomerUseCase(repo domain.CustomerRepository, jwtSecret string) *CustomerUseCase {
	return &CustomerUseCase{
		repo:      repo,
		jwtSecret: jwtSecret,
	}
}

type AuthResponse struct {
	Token string `json:"token"`
}

func (uc *CustomerUseCase) Register(ctx context.Context, name, password string) (*AuthResponse, error) {
	existing, err := uc.repo.FindByName(ctx, name)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, domain.ErrCustomerAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	customer := &domain.Customer{
		Name:     name,
		Password: string(hashedPassword),
	}

	if err := uc.repo.Create(ctx, customer); err != nil {
		return nil, err
	}

	token, err := uc.generateJWT(customer)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{Token: token}, nil
}

func (uc *CustomerUseCase) Login(ctx context.Context, name, password string) (*AuthResponse, error) {
	customer, err := uc.repo.FindByName(ctx, name)
	if err != nil {
		return nil, err
	}
	if customer == nil {
		return nil, domain.ErrCustomerNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(customer.Password), []byte(password)); err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	token, err := uc.generateJWT(customer)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{Token: token}, nil
}

func (uc *CustomerUseCase) generateJWT(customer *domain.Customer) (string, error) {
	claims := jwt.MapClaims{
		"customer_id": customer.ID,
		"username":    customer.Name,
		"exp":         time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(uc.jwtSecret))
}
