package services

import (
	"context"
	"errors"
	"strings"

	"github.com/spalqui/habitattrack-api/internal/models"
	"github.com/spalqui/habitattrack-api/internal/repositories"
)

type CategoryService interface {
	CreateCategory(ctx context.Context, category *models.Category) error
	GetCategory(ctx context.Context, id string) (*models.Category, error)
	GetAllCategories(ctx context.Context) ([]*models.Category, error)
	GetCategoriesByType(ctx context.Context, transactionType models.TransactionType) ([]*models.Category, error)
	UpdateCategory(ctx context.Context, category *models.Category) error
	DeleteCategory(ctx context.Context, id string) error
}

type categoryService struct {
	categoryRepo repositories.CategoryRepository
}

func NewCategoryService(categoryRepo repositories.CategoryRepository) CategoryService {
	return &categoryService{
		categoryRepo: categoryRepo,
	}
}

func (s *categoryService) CreateCategory(ctx context.Context, category *models.Category) error {
	if err := s.validateCategory(category); err != nil {
		return err
	}

	return s.categoryRepo.Create(ctx, category)
}

func (s *categoryService) GetCategory(ctx context.Context, id string) (*models.Category, error) {
	if strings.TrimSpace(id) == "" {
		return nil, errors.New("category ID is required")
	}

	return s.categoryRepo.GetByID(ctx, id)
}

func (s *categoryService) GetAllCategories(ctx context.Context) ([]*models.Category, error) {
	return s.categoryRepo.GetAll(ctx)
}

func (s *categoryService) GetCategoriesByType(ctx context.Context, transactionType models.TransactionType) ([]*models.Category, error) {
	if transactionType != models.TransactionTypeIncome && transactionType != models.TransactionTypeExpense {
		return nil, errors.New("invalid transaction type")
	}

	return s.categoryRepo.GetByType(ctx, transactionType)
}

func (s *categoryService) UpdateCategory(ctx context.Context, category *models.Category) error {
	if err := s.validateCategory(category); err != nil {
		return err
	}

	if strings.TrimSpace(category.ID) == "" {
		return errors.New("category ID is required for update")
	}

	return s.categoryRepo.Update(ctx, category)
}

func (s *categoryService) DeleteCategory(ctx context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		return errors.New("category ID is required")
	}

	return s.categoryRepo.Delete(ctx, id)
}

func (s *categoryService) validateCategory(category *models.Category) error {
	if strings.TrimSpace(category.Name) == "" {
		return errors.New("category name is required")
	}

	if category.Type != models.TransactionTypeIncome && category.Type != models.TransactionTypeExpense {
		return errors.New("invalid transaction type")
	}

	return nil
}
