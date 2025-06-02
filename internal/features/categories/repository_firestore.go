package categories

import (
	"context"
	"errors"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/spalqui/habitattrack-api/constants"
	"github.com/spalqui/habitattrack-api/internal/core/entity"
)

const (
	categoryCollection = "transaction_categories"
)

var (
	ErrEmptyProjectID   = errors.New("project ID cannot be empty")
	ErrCategoryNotFound = errors.New("transaction category not found")
	ErrNilCategory      = errors.New("transaction category cannot be nil")
	ErrEmptyCategoryID  = errors.New("transaction category ID cannot be empty")
)

// FirestoreCategoryRepository implements the categories.TransactionCategoryRepository interface.
type FirestoreCategoryRepository struct {
	client *firestore.Client
}

// NewFirestoreCategoryRepository creates a new FirestoreCategoryRepository.
func NewFirestoreCategoryRepository(ctx context.Context, projectID string, databaseID string) (*FirestoreCategoryRepository, error) {
	if projectID == "" {
		return nil, ErrEmptyProjectID // Reusing from repository_firestore.go
	}

	ctx, cancel := context.WithTimeout(ctx, constants.DefaultTimeout)
	defer cancel()

	client, err := firestore.NewClientWithDatabase(ctx, projectID, databaseID)
	if err != nil {
		return nil, fmt.Errorf("creating firestore client for categories: %w", err)
	}

	return &FirestoreCategoryRepository{
		client: client,
	}, nil
}

// CloseClient closes the Firestore client.
func (r *FirestoreCategoryRepository) CloseClient() error {
	if r.client != nil {
		if err := r.client.Close(); err != nil {
			return fmt.Errorf("closing firestore client for categories: %w", err)
		}
	}
	return nil
}

// Create implements categories.TransactionCategoryRepository
func (r *FirestoreCategoryRepository) Create(ctx context.Context, category *entity.TransactionCategory) (string, error) {
	if category == nil {
		return "", ErrNilCategory
	}

	ctx, cancel := context.WithTimeout(ctx, constants.DefaultTimeout)
	defer cancel()

	now := time.Now().UTC()
	if category.CreatedAt.IsZero() {
		category.CreatedAt = now
	}
	if category.UpdatedAt.IsZero() {
		category.UpdatedAt = now
	}

	docRef := r.client.Collection(categoryCollection).Doc(category.ID) // Assuming ID is pre-set by service
	_, err := docRef.Set(ctx, category)
	if err != nil {
		return "", fmt.Errorf("creating category %s: %w", category.ID, err)
	}
	return category.ID, nil
}

// GetByID implements categories.TransactionCategoryRepository
func (r *FirestoreCategoryRepository) GetByID(ctx context.Context, id string) (*entity.TransactionCategory, error) {
	if id == "" {
		return nil, ErrEmptyCategoryID
	}

	ctx, cancel := context.WithTimeout(ctx, constants.DefaultTimeout)
	defer cancel()

	doc, err := r.client.Collection(categoryCollection).Doc(id).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, ErrCategoryNotFound
		}
		return nil, fmt.Errorf("getting category by ID %s: %w", id, err)
	}

	var category entity.TransactionCategory
	if err := doc.DataTo(&category); err != nil {
		return nil, fmt.Errorf("parsing category data for ID %s: %w", id, err)
	}
	category.ID = doc.Ref.ID
	return &category, nil
}

// List implements categories.TransactionCategoryRepository
func (r *FirestoreCategoryRepository) List(
	ctx context.Context,
	classification *entity.ClassificationType,
	limit, offset int,
) ([]*entity.TransactionCategory, int, error) {
	ctx, cancel := context.WithTimeout(ctx, constants.DefaultTimeout)
	defer cancel()

	query := r.client.Collection(categoryCollection).Query

	if classification != nil && *classification != "" {
		query = query.Where("Classification", "==", *classification)
	}

	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, 0, fmt.Errorf("counting categories: %w", err)
	}
	totalItems := len(docs)

	dataQuery := query.OrderBy("Name", firestore.Asc).Limit(limit).Offset(offset)
	iter := dataQuery.Documents(ctx)
	defer iter.Stop()

	var categories []*entity.TransactionCategory
	for {
		doc, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return nil, 0, fmt.Errorf("iterating categories for list: %w", err)
		}
		var cat entity.TransactionCategory
		if err := doc.DataTo(&cat); err != nil {
			return nil, 0, fmt.Errorf("parsing category data for list: %w", err)
		}
		cat.ID = doc.Ref.ID
		categories = append(categories, &cat)
	}
	return categories, totalItems, nil
}

// Update implements categories.TransactionCategoryRepository
func (r *FirestoreCategoryRepository) Update(ctx context.Context, category *entity.TransactionCategory) error {
	if category == nil || category.ID == "" {
		return ErrNilCategory // Or ErrEmptyCategoryID
	}

	ctx, cancel := context.WithTimeout(ctx, constants.DefaultTimeout)
	defer cancel()

	docRef := r.client.Collection(categoryCollection).Doc(category.ID)
	category.UpdatedAt = time.Now().UTC()

	_, err := docRef.Set(ctx, category)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return ErrCategoryNotFound
		}
		return fmt.Errorf("updating category %s: %w", category.ID, err)
	}
	return nil
}

// Delete implements categories.TransactionCategoryRepository
func (r *FirestoreCategoryRepository) Delete(ctx context.Context, id string) error {
	if id == "" {
		return ErrEmptyCategoryID
	}

	ctx, cancel := context.WithTimeout(ctx, constants.DefaultTimeout)
	defer cancel()

	docRef := r.client.Collection(categoryCollection).Doc(id)
	_, err := docRef.Delete(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return ErrCategoryNotFound // Or return nil
		}
		return fmt.Errorf("deleting category %s: %w", id, err)
	}
	return nil
}

// FindByNameAndClassification implements categories.TransactionCategoryRepository
func (r *FirestoreCategoryRepository) FindByNameAndClassification(ctx context.Context, name string, classification entity.ClassificationType) (*entity.TransactionCategory, error) {
	if name == "" {
		return nil, errors.New("category name cannot be empty for FindByNameAndClassification")
	}
	if classification == "" {
		return nil, errors.New("classification cannot be empty for FindByNameAndClassification")
	}

	ctx, cancel := context.WithTimeout(ctx, constants.DefaultTimeout)
	defer cancel()

	query := r.client.Collection(categoryCollection).
		Where("Name", "==", name).
		Where("Classification", "==", classification).
		Limit(1)
	iter := query.Documents(ctx)
	defer iter.Stop()

	doc, err := iter.Next()
	if err == iterator.Done {
		return nil, ErrCategoryNotFound // Using specific error
	}
	if err != nil {
		return nil, fmt.Errorf("querying category by name '%s' and classification '%s': %w", name, classification, err)
	}

	var category entity.TransactionCategory
	if err := doc.DataTo(&category); err != nil {
		return nil, fmt.Errorf("parsing category data for FindByNameAndClassification: %w", err)
	}
	category.ID = doc.Ref.ID
	return &category, nil
}
