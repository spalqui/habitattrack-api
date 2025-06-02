package transactions

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"

	"github.com/spalqui/habitattrack-api/internal/core/entity"
	"github.com/spalqui/habitattrack-api/internal/features/categories" // To check category existence and type
	"github.com/spalqui/habitattrack-api/internal/shared/apierrors"
	"github.com/spalqui/habitattrack-api/internal/shared/apitypes" // For PaginationInfo
)

const (
	defaultServicePageSize = 20
	maxServicePageSize     = 100
	defaultServicePageNum  = 1
)

// TransactionService defines the interface for transaction business logic.
type TransactionService interface {
	CreateTransaction(ctx context.Context, req CreateTransactionRequest) (*TransactionResponse, error)
	GetTransactionByID(ctx context.Context, id string) (*TransactionResponse, error)
	ListTransactions(
		ctx context.Context,
		propertyIDFilter *string,
		typeFilter *entity.TransactionType,
		categoryIDFilter *string,
		startDateFilter *time.Time,
		endDateFilter *time.Time,
		page, pageSize int,
	) (*PaginatedTransactionsResponse, error)
	UpdateTransaction(ctx context.Context, id string, req UpdateTransactionRequest) (*TransactionResponse, error)
	DeleteTransaction(ctx context.Context, id string) error
}

type transactionService struct {
	repo         TransactionRepository
	categoryRepo categories.TransactionCategoryRepository // Dependency to validate category
}

// NewTransactionService creates a new instance of TransactionService.
func NewTransactionService(repo TransactionRepository, categoryRepo categories.TransactionCategoryRepository) TransactionService {
	return &transactionService{
		repo:         repo,
		categoryRepo: categoryRepo,
	}
}

func (s *transactionService) CreateTransaction(ctx context.Context, req CreateTransactionRequest) (*TransactionResponse, error) {
	// Validate CategoryID
	category, err := s.categoryRepo.GetByID(ctx, req.CategoryID)
	if err != nil {
		// Assuming GetByID returns an error that can be identified as "not found".
		// If so, map to apierrors.ErrNotFound or a specific validation error.
		// For now, a validation error indicating the categoryId is problematic.
		return nil, apierrors.ErrValidation("Validation failed for categoryId.", map[string]string{"categoryId": "Category not found or invalid."})
	}

	// Validate if transaction type matches category classification
	if (req.Type == entity.IncomeTransaction && category.Classification != entity.IncomeClassification) ||
		(req.Type == entity.ExpenseTransaction && category.Classification != entity.ExpenseClassification) {
		return nil, apierrors.ErrValidation(
			fmt.Sprintf("Transaction type '%s' does not match category classification '%s'.", req.Type, category.Classification),
			map[string]string{"type": "Type mismatch with category classification."},
		)
	}

	// PropertyID validation (if exists) would go here.

	transaction := req.ToEntity()
	transaction.ID = uuid.NewString()
	transaction.CreatedAt = time.Now().UTC()
	transaction.UpdatedAt = transaction.CreatedAt

	_, err = s.repo.Create(ctx, transaction)
	if err != nil {
		return nil, apierrors.ErrInternal(fmt.Errorf("creating transaction: %w", err))
	}

	resp := ToTransactionResponse(transaction)
	return &resp, nil
}

func (s *transactionService) GetTransactionByID(ctx context.Context, id string) (*TransactionResponse, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, apierrors.ErrBadRequest(fmt.Sprintf("Invalid transaction ID format: %s", id))
	}

	transaction, err := s.repo.GetByID(ctx, id)
	if err != nil {
		// Assuming the repository returns an error that can be identified as "not found".
		return nil, apierrors.ErrNotFound("Transaction", id) // Simplified, needs robust check
	}

	resp := ToTransactionResponse(transaction)
	return &resp, nil
}

func (s *transactionService) ListTransactions(
	ctx context.Context,
	propertyIDFilter *string,
	typeFilter *entity.TransactionType,
	categoryIDFilter *string,
	startDateFilter *time.Time,
	endDateFilter *time.Time,
	page, pageSize int,
) (*PaginatedTransactionsResponse, error) {
	if page <= 0 {
		page = defaultServicePageNum
	}
	if pageSize <= 0 {
		pageSize = defaultServicePageSize
	}
	if pageSize > maxServicePageSize {
		pageSize = maxServicePageSize
	}
	offset := (page - 1) * pageSize

	transactions, totalItems, err := s.repo.List(ctx, propertyIDFilter, typeFilter, categoryIDFilter, startDateFilter, endDateFilter, pageSize, offset)
	if err != nil {
		return nil, apierrors.ErrInternal(fmt.Errorf("listing transactions: %w", err))
	}

	responseList := ToTransactionResponseList(transactions)
	totalPages := 0
	if totalItems > 0 {
		totalPages = int(math.Ceil(float64(totalItems) / float64(pageSize)))
	}

	return &PaginatedTransactionsResponse{
		Data: responseList,
		Pagination: apitypes.PaginationInfo{
			TotalItems:  int64(totalItems),
			TotalPages:  totalPages,
			CurrentPage: page,
			PageSize:    pageSize,
		},
	}, nil
}

func (s *transactionService) UpdateTransaction(ctx context.Context, id string, req UpdateTransactionRequest) (*TransactionResponse, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, apierrors.ErrBadRequest(fmt.Sprintf("Invalid transaction ID format: %s", id))
	}

	transaction, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, apierrors.ErrNotFound("Transaction", id) // Simplified, needs robust check
	}

	updated := false
	if req.Amount != nil && *req.Amount != transaction.Amount {
		transaction.Amount = *req.Amount
		updated = true
	}
	if req.TransactionDate != nil && !req.TransactionDate.IsZero() && !req.TransactionDate.Equal(transaction.TransactionDate) {
		transaction.TransactionDate = *req.TransactionDate
		updated = true
	}

	// Handle Description update (pointer to string)
	if req.Description != nil {
		if transaction.Description == nil || *req.Description != *transaction.Description {
			transaction.Description = req.Description
			updated = true
		}
	}
	// Note: If req.Description is nil (field not provided in PATCH), no change occurs.
	// If "description": null is sent, and ShouldBindJSON correctly unmarshals it to a nil *string,
	// then the above logic correctly does not update if transaction.Description was already nil.
	// If transaction.Description was non-nil, and "description":null is sent, it will be set to nil.

	if req.Type != nil && *req.Type != transaction.Type {
		transaction.Type = *req.Type
		updated = true
	}
	if req.CategoryID != nil && *req.CategoryID != transaction.CategoryID {
		transaction.CategoryID = *req.CategoryID
		updated = true
	}

	// Handle PropertyID update (pointer to string)
	if req.PropertyID != nil { // If PropertyID is provided in the request
		if transaction.PropertyID == nil || *req.PropertyID != *transaction.PropertyID {
			transaction.PropertyID = req.PropertyID // Update to new value (could be pointer to empty string)
			updated = true
		}
	}
	// If req.PropertyID is nil (e.g. field not in PATCH request), no change to transaction.PropertyID.
	// If the request explicitly sends "propertyId": null, and it unmarshals to req.PropertyID being a nil *string,
	// this logic correctly does not change transaction.PropertyID.
	// To explicitly clear PropertyID, the DTO might need a `ClearPropertyID: true` or similar,
	// or the client sends `propertyId: "" ` if an empty string means "cleared" vs `null` meaning "no change".
	// The current DTO `validate:"omitempty,uuid"` means an empty string for PropertyID would fail UUID validation.
	// So, `null` is the only way to "not provide" or "not change" via PATCH for optional UUIDs.

	if updated {
		// Re-validate category and type if they changed or if category changed (type might need re-check)
		if req.CategoryID != nil || req.Type != nil {
			currentCategoryID := transaction.CategoryID // Use the potentially updated category ID
			currentType := transaction.Type             // Use the potentially updated type

			category, err := s.categoryRepo.GetByID(ctx, currentCategoryID)
			if err != nil {
				return nil, apierrors.ErrValidation("Validation failed for categoryId.", map[string]string{"categoryId": "Category not found or invalid."})
			}
			if (currentType == entity.IncomeTransaction && category.Classification != entity.IncomeClassification) ||
				(currentType == entity.ExpenseTransaction && category.Classification != entity.ExpenseClassification) {
				return nil, apierrors.ErrValidation(
					fmt.Sprintf("Transaction type '%s' does not match category classification '%s'.", currentType, category.Classification),
					map[string]string{"type": "Type mismatch with category classification."},
				)
			}
		}
		// PropertyID validation (if exists and changed) would go here.

		transaction.UpdatedAt = time.Now().UTC()
		err = s.repo.Update(ctx, transaction)
		if err != nil {
			return nil, apierrors.ErrInternal(fmt.Errorf("updating transaction %s: %w", id, err))
		}
	}

	resp := ToTransactionResponse(transaction)
	return &resp, nil
}

func (s *transactionService) DeleteTransaction(ctx context.Context, id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return apierrors.ErrBadRequest(fmt.Sprintf("Invalid transaction ID format: %s", id))
	}

	// Check if transaction exists before attempting to delete
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return apierrors.ErrNotFound("Transaction", id) // Simplified, needs robust check
	}

	err = s.repo.Delete(ctx, id)
	if err != nil {
		// If GetByID passed but Delete fails, it's an internal error.
		return apierrors.ErrInternal(fmt.Errorf("deleting transaction %s: %w", id, err))
	}
	return nil
}
