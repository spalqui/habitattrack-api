package categories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"

	"github.com/spalqui/habitattrack-api/internal/core/entity"
	"github.com/spalqui/habitattrack-api/internal/shared/apierrors"
	"github.com/spalqui/habitattrack-api/internal/shared/apitypes"
)

const (
	defaultServicePageSize = 20
	maxServicePageSize     = 100
	defaultServicePageNum  = 1
)

// TransactionCategoryService defines the interface for category business logic.
// This is what handlers will interact with.
type TransactionCategoryService interface {
	CreateCategory(ctx context.Context, req CreateCategoryRequest) (*CategoryResponse, error)
	GetCategoryByID(ctx context.Context, id string) (*CategoryResponse, error)
	ListCategories(ctx context.Context, classificationFilter *entity.ClassificationType, page, pageSize int) (*PaginatedCategoriesResponse, error)
	UpdateCategory(ctx context.Context, id string, req UpdateCategoryRequest) (*CategoryResponse, error)
	DeleteCategory(ctx context.Context, id string) error
}

// categoryService implements TransactionCategoryService.
type categoryService struct {
	repo TransactionCategoryRepository
}

// NewCategoryService creates a new instance of TransactionCategoryService.
func NewCategoryService(repo TransactionCategoryRepository) TransactionCategoryService {
	return &categoryService{
		repo: repo,
	}
}

// CreateCategory handles the business logic for creating a new transaction category.
func (s *categoryService) CreateCategory(ctx context.Context, req CreateCategoryRequest) (*CategoryResponse, error) {
	// Handler performs struct validation. Specific business rule validations are here.
	if req.Name == "" { // This should ideally be caught by `validate:"required"` on DTO
		return nil, apierrors.ErrValidation("Validation failed", map[string]string{"name": "Name is required."})
	}
	if req.Classification != entity.IncomeClassification && req.Classification != entity.ExpenseClassification {
		// This should be caught by `validate:"oneof=income expense"` on DTO
		return nil, apierrors.ErrValidation("Validation failed", map[string]string{"classification": "Classification must be 'income' or 'expense'."})
	}

	// Check for uniqueness: name + classification
	existing, err := s.repo.FindByNameAndClassification(ctx, req.Name, req.Classification)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		// Any error other than "not found" is an internal server error.
		return nil, apierrors.ErrInternal(fmt.Errorf("checking for existing category: %w", err))
	}
	if existing != nil {
		return nil, apierrors.ErrConflict(fmt.Sprintf("Category with name '%s' and classification '%s' already exists.", req.Name, req.Classification))
	}

	category := req.ToEntity()
	category.ID = uuid.NewString()
	category.CreatedAt = time.Now().UTC()
	category.UpdatedAt = category.CreatedAt

	_, err = s.repo.Create(ctx, category)
	if err != nil {
		return nil, apierrors.ErrInternal(fmt.Errorf("creating category: %w", err))
	}

	resp := ToCategoryResponse(category)
	return &resp, nil
}

// GetCategoryByID retrieves a category by its ID.
func (s *categoryService) GetCategoryByID(ctx context.Context, id string) (*CategoryResponse, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, apierrors.ErrBadRequest(fmt.Sprintf("Invalid category ID format: %s", id))
	}

	category, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) { // Assuming repository returns sql.ErrNoRows for not found
			return nil, apierrors.ErrNotFound("Category", id)
		}
		return nil, apierrors.ErrInternal(fmt.Errorf("getting category by ID %s: %w", id, err))
	}

	resp := ToCategoryResponse(category)
	return &resp, nil
}

// ListCategories retrieves a paginated list of categories, optionally filtered by classification.
func (s *categoryService) ListCategories(ctx context.Context, classificationFilter *entity.ClassificationType, page, pageSize int) (*PaginatedCategoriesResponse, error) {
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

	categories, totalItems, err := s.repo.List(ctx, classificationFilter, pageSize, offset)
	if err != nil {
		return nil, apierrors.ErrInternal(fmt.Errorf("listing categories: %w", err))
	}

	responseList := ToCategoryResponseList(categories)
	totalPages := 0
	if totalItems > 0 {
		totalPages = int(math.Ceil(float64(totalItems) / float64(pageSize)))
	}

	return &PaginatedCategoriesResponse{
		Data: responseList,
		Pagination: apitypes.PaginationInfo{
			TotalItems:  int64(totalItems),
			TotalPages:  totalPages,
			CurrentPage: page,
			PageSize:    pageSize,
		},
	}, nil
}

// UpdateCategory handles updating an existing category.
func (s *categoryService) UpdateCategory(ctx context.Context, id string, req UpdateCategoryRequest) (*CategoryResponse, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, apierrors.ErrBadRequest(fmt.Sprintf("Invalid category ID format: %s", id))
	}

	category, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apierrors.ErrNotFound("Category", id)
		}
		return nil, apierrors.ErrInternal(fmt.Errorf("getting category by ID %s for update: %w", id, err))
	}

	updated := false
	originalName := category.Name
	originalClassification := category.Classification

	if req.Name != nil && *req.Name != "" && *req.Name != category.Name {
		category.Name = *req.Name
		updated = true
	}
	if req.Classification != nil && *req.Classification != "" && *req.Classification != category.Classification {
		if *req.Classification != entity.IncomeClassification && *req.Classification != entity.ExpenseClassification {
			return nil, apierrors.ErrValidation("Validation failed", map[string]string{"classification": "Classification must be 'income' or 'expense'."})
		}
		category.Classification = *req.Classification
		updated = true
	}

	if updated {
		// Check for uniqueness only if name or classification actually changed to new values
		if category.Name != originalName || category.Classification != originalClassification {
			existing, err := s.repo.FindByNameAndClassification(ctx, category.Name, category.Classification)
			if err != nil && !errors.Is(err, sql.ErrNoRows) {
				return nil, apierrors.ErrInternal(fmt.Errorf("checking for existing category during update: %w", err))
			}
			if existing != nil && existing.ID != id { // If it exists and it's not the same category
				return nil, apierrors.ErrConflict(fmt.Sprintf("Another category with name '%s' and classification '%s' already exists.", category.Name, category.Classification))
			}
		}

		category.UpdatedAt = time.Now().UTC()
		err = s.repo.Update(ctx, category)
		if err != nil {
			return nil, apierrors.ErrInternal(fmt.Errorf("updating category %s: %w", id, err))
		}
	}

	resp := ToCategoryResponse(category)
	return &resp, nil
}

// DeleteCategory handles deleting a category by its ID.
func (s *categoryService) DeleteCategory(ctx context.Context, id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return apierrors.ErrBadRequest(fmt.Sprintf("Invalid category ID format: %s", id))
	}

	// Check if category exists before attempting to delete
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apierrors.ErrNotFound("Category", id)
		}
		return apierrors.ErrInternal(fmt.Errorf("checking category %s before delete: %w", id, err))
	}

	// Add check for category in use by transactions here if that logic is implemented
	// e.g., count, err := s.transactionRepo.CountByCategoryID(ctx, id)
	// if count > 0 { return apierrors.ErrConflict("Category is in use and cannot be deleted.") }

	err = s.repo.Delete(ctx, id)
	if err != nil {
		// If GetByID passed but Delete fails, it might be a concurrent deletion or other DB error.
		// Or, if a foreign key constraint exists and wasn't checked above, it could be apierrors.ErrConflict.
		return apierrors.ErrInternal(fmt.Errorf("deleting category %s: %w", id, err))
	}
	return nil
}
