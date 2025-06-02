package repositories

import (
	"context"

	"github.com/spalqui/habitattrack-api/internal/models"
)

type TransactionRepository interface {
	Create(ctx context.Context, transaction *models.Transaction) error
	GetByID(ctx context.Context, id string) (*models.Transaction, error)
	GetByPropertyID(ctx context.Context, propertyID string) ([]*models.Transaction, error)
	GetAll(ctx context.Context) ([]*models.Transaction, error)
	Update(ctx context.Context, transaction *models.Transaction) error
	Delete(ctx context.Context, id string) error
}
