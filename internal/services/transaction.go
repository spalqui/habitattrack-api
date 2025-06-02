package services

import (
	"context"
	"errors"
	"strings"

	"github.com/spalqui/habitattrack-api/internal/models"
	"github.com/spalqui/habitattrack-api/internal/repositories"
)

type TransactionService interface {
	CreateTransaction(ctx context.Context, transaction *models.Transaction) error
	GetTransaction(ctx context.Context, id string) (*models.Transaction, error)
	GetTransactionsByProperty(ctx context.Context, propertyID string) ([]*models.Transaction, error)
	GetAllTransactions(ctx context.Context) ([]*models.Transaction, error)
	UpdateTransaction(ctx context.Context, transaction *models.Transaction) error
	DeleteTransaction(ctx context.Context, id string) error
}

type transactionService struct {
	transactionRepo repositories.TransactionRepository
	categoryRepo    repositories.CategoryRepository
	propertyRepo    repositories.PropertyRepository
}

func NewTransactionService(
	transactionRepo repositories.TransactionRepository,
	categoryRepo repositories.CategoryRepository,
	propertyRepo repositories.PropertyRepository,
) TransactionService {
	return &transactionService{
		transactionRepo: transactionRepo,
		categoryRepo:    categoryRepo,
		propertyRepo:    propertyRepo,
	}
}

func (s *transactionService) CreateTransaction(ctx context.Context, transaction *models.Transaction) error {
	if err := s.validateTransaction(ctx, transaction); err != nil {
		return err
	}

	return s.transactionRepo.Create(ctx, transaction)
}

func (s *transactionService) GetTransaction(ctx context.Context, id string) (*models.Transaction, error) {
	if strings.TrimSpace(id) == "" {
		return nil, errors.New("transaction ID is required")
	}

	return s.transactionRepo.GetByID(ctx, id)
}

func (s *transactionService) GetTransactionsByProperty(ctx context.Context, propertyID string) ([]*models.Transaction, error) {
	if strings.TrimSpace(propertyID) == "" {
		return nil, errors.New("property ID is required")
	}

	return s.transactionRepo.GetByPropertyID(ctx, propertyID)
}

func (s *transactionService) GetAllTransactions(ctx context.Context) ([]*models.Transaction, error) {
	return s.transactionRepo.GetAll(ctx)
}

func (s *transactionService) UpdateTransaction(ctx context.Context, transaction *models.Transaction) error {
	if err := s.validateTransaction(ctx, transaction); err != nil {
		return err
	}

	if strings.TrimSpace(transaction.ID) == "" {
		return errors.New("transaction ID is required for update")
	}

	return s.transactionRepo.Update(ctx, transaction)
}

func (s *transactionService) DeleteTransaction(ctx context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		return errors.New("transaction ID is required")
	}

	return s.transactionRepo.Delete(ctx, id)
}

func (s *transactionService) validateTransaction(ctx context.Context, transaction *models.Transaction) error {
	if strings.TrimSpace(transaction.PropertyID) == "" {
		return errors.New("property ID is required")
	}

	if strings.TrimSpace(transaction.CategoryID) == "" {
		return errors.New("category ID is required")
	}

	if transaction.Amount <= 0 {
		return errors.New("amount must be greater than zero")
	}

	if transaction.Type != models.TransactionTypeIncome && transaction.Type != models.TransactionTypeExpense {
		return errors.New("invalid transaction type")
	}

	// Verify property exists
	if _, err := s.propertyRepo.GetByID(ctx, transaction.PropertyID); err != nil {
		return errors.New("property not found")
	}

	// Verify category exists and matches transaction type
	category, err := s.categoryRepo.GetByID(ctx, transaction.CategoryID)
	if err != nil {
		return errors.New("category not found")
	}

	if category.Type != transaction.Type {
		return errors.New("category type does not match transaction type")
	}

	return nil
}
