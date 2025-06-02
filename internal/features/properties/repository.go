package properties

import (
	"context"

	"github.com/spalqui/habitattrack-api/internal/core/entity"
)

// PropertyRepository defines the interface for interacting with property storage.
type PropertyRepository interface {
	Create(ctx context.Context, property *entity.Property) (string, error)
	GetByID(ctx context.Context, id string) (*entity.Property, error)
	List(ctx context.Context, limit, offset int) ([]*entity.Property, int, error) // Returns properties and total count
	Update(ctx context.Context, property *entity.Property) error
	Delete(ctx context.Context, id string) error
	FindByName(ctx context.Context, name string) (*entity.Property, error) // For checking uniqueness if needed
}