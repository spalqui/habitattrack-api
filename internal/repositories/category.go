package repositories

import (
	"context"

	"github.com/spalqui/habitattrack-api/internal/models"
)

type CategoryRepository interface {
	Create(ctx context.Context, category *models.Category) error
	GetByID(ctx context.Context, id string) (*models.Category, error)
	GetAll(ctx context.Context) ([]*models.Category, error)
	GetByType(ctx context.Context, transactionType models.TransactionType) ([]*models.Category, error)
	Update(ctx context.Context, category *models.Category) error
	Delete(ctx context.Context, id string) error
}
