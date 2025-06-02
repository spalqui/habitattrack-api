package transactions

import (
	"context"
	"time"

	"github.com/spalqui/habitattrack-api/internal/core/entity"
)

// TransactionRepository defines the interface for interacting with transaction storage.
type TransactionRepository interface {
	Create(ctx context.Context, transaction *entity.Transaction) (string, error)
	GetByID(ctx context.Context, id string) (*entity.Transaction, error)
	List(
		ctx context.Context,
		propertyID *string,
		transactionType *entity.TransactionType,
		categoryID *string,
		startDate *time.Time,
		endDate *time.Time,
		limit, offset int,
	) ([]*entity.Transaction, int, error) // Returns transactions and total count
	Update(ctx context.Context, transaction *entity.Transaction) error
	Delete(ctx context.Context, id string) error
}