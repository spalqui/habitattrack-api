package categories

import (
	"context"

	"github.com/spalqui/habitattrack-api/internal/core/entity"
)

// TransactionCategoryRepository defines the interface for interacting with transaction category storage.
// This interface belongs to the Application Layer, and its implementation will be in the Infrastructure Layer.
type TransactionCategoryRepository interface {
	Create(ctx context.Context, category *entity.TransactionCategory) (string, error)
	GetByID(ctx context.Context, id string) (*entity.TransactionCategory, error)
	List(ctx context.Context, classification *entity.ClassificationType, limit, offset int) ([]*entity.TransactionCategory, int, error) // Returns categories and total count
	Update(ctx context.Context, category *entity.TransactionCategory) error
	Delete(ctx context.Context, id string) error
	FindByNameAndClassification(ctx context.Context, name string, classification entity.ClassificationType) (*entity.TransactionCategory, error)
}