package repositories

import (
	"context"

	"github.com/spalqui/habitattrack-api/internal/models"
)

type PropertyRepository interface {
	Create(ctx context.Context, property *models.Property) error
	GetByID(ctx context.Context, id string) (*models.Property, error)
	GetAll(ctx context.Context) ([]*models.Property, error)
	Update(ctx context.Context, property *models.Property) error
	Delete(ctx context.Context, id string) error
}
