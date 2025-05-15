package repositories

import (
	"context"
	"errors"
	"fmt"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/spalqui/habitattrack-api/constants"
	"github.com/spalqui/habitattrack-api/types"
)

const (
	propertyCollection = "properties"
)

var (
	ErrEmptyProjectID   = errors.New("project ID cannot be empty")
	ErrEmptyPropertyID  = errors.New("property ID cannot be empty")
	ErrInvalidLimit     = errors.New("limit must be between 1 and 100")
	ErrNilProperty      = errors.New("property cannot be nil")
	ErrPropertyNotFound = errors.New("property not found")
)

// PropertyRepository defines the interface for property repository
type PropertyRepository interface {
	GetPropertyByID(ctx context.Context, id string) (types.Property, error)
	GetProperties(ctx context.Context, limit int, cursor string) ([]types.Property, string, error)
	CreateProperty(ctx context.Context, property *types.Property) error
	UpdateProperty(ctx context.Context, id string, property *types.Property) error
	DeleteProperty(ctx context.Context, id string) error
	Close() error
}

// FirestorePropertyRepository implements PropertyRepository using Firestore
type FirestorePropertyRepository struct {
	client *firestore.Client
}

func NewFirestorePropertyRepository(ctx context.Context, projectID string, databaseID string) (PropertyRepository, error) {
	if projectID == "" {
		return nil, ErrEmptyProjectID
	}

	ctx, cancel := context.WithTimeout(ctx, constants.DefaultTimeout)
	defer cancel()

	client, err := firestore.NewClientWithDatabase(ctx, projectID, databaseID)
	if err != nil {
		return nil, fmt.Errorf("creating firestore client: %w", err)
	}

	return &FirestorePropertyRepository{
		client: client,
	}, nil
}

func (r *FirestorePropertyRepository) Close() error {
	if r.client != nil {
		if err := r.client.Close(); err != nil {
			return fmt.Errorf("closing firestore client: %w", err)
		}
	}
	return nil
}

func (r *FirestorePropertyRepository) GetPropertyByID(ctx context.Context, id string) (types.Property, error) {
	if id == "" {
		return types.Property{}, ErrEmptyPropertyID
	}

	ctx, cancel := context.WithTimeout(ctx, constants.DefaultTimeout)
	defer cancel()

	doc, err := r.client.Collection(propertyCollection).Doc(id).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return types.Property{}, ErrPropertyNotFound
		}
		return types.Property{}, fmt.Errorf("getting property by ID: %w", err)
	}

	var property types.Property
	err = doc.DataTo(&property)
	if err != nil {
		return types.Property{}, fmt.Errorf("parsing property data: %w", err)
	}

	property.ID = doc.Ref.ID
	return property, nil
}

func (r *FirestorePropertyRepository) GetProperties(ctx context.Context, limit int, cursor string) ([]types.Property, string, error) {
	if limit < 1 || limit > constants.MaxPageSize {
		return nil, "", ErrInvalidLimit
	}

	ctx, cancel := context.WithTimeout(ctx, constants.DefaultTimeout)
	defer cancel()

	query := r.client.Collection(propertyCollection).
		//OrderBy("createdAt", firestore.Desc).
		Limit(limit)

	if cursor != "" {
		query = query.StartAfter(cursor)
	}

	iter := query.Documents(ctx)
	defer iter.Stop()

	var properties []types.Property
	for {
		doc, err := iter.Next()
		if err != nil {
			if errors.Is(err, iterator.Done) {
				break
			}
			if status.Code(err) == codes.NotFound {
				break
			}
			return nil, "", fmt.Errorf("iterating properties: %w", err)
		}

		var property types.Property
		err = doc.DataTo(&property)
		if err != nil {
			return nil, "", fmt.Errorf("parsing property data: %w", err)
		}

		property.ID = doc.Ref.ID

		properties = append(properties, property)
		cursor = doc.Ref.ID
	}

	return properties, cursor, nil
}

func (r *FirestorePropertyRepository) CreateProperty(ctx context.Context, property *types.Property) error {
	if *property == (types.Property{}) {
		return ErrNilProperty
	}

	ctx, cancel := context.WithTimeout(ctx, constants.DefaultTimeout)
	defer cancel()

	docRef, _, err := r.client.Collection(propertyCollection).Add(ctx, property)
	if err != nil {
		return fmt.Errorf("adding property: %w", err)
	}

	property.ID = docRef.ID
	return nil
}

func (r *FirestorePropertyRepository) UpdateProperty(ctx context.Context, id string, property *types.Property) error {
	if id == "" {
		return ErrEmptyPropertyID
	}
	if *property == (types.Property{}) {
		return ErrNilProperty
	}

	ctx, cancel := context.WithTimeout(ctx, constants.DefaultTimeout)
	defer cancel()

	// Ensure the property exists before updating
	docRef, err := r.client.Collection(propertyCollection).Doc(id).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return ErrPropertyNotFound
		}
		return fmt.Errorf("checking property exists: %w", err)
	}

	_, err = docRef.Ref.Set(ctx, property)
	if err != nil {
		return fmt.Errorf("updating property: %w", err)
	}

	property.ID = docRef.Ref.ID

	return nil
}

func (r *FirestorePropertyRepository) DeleteProperty(ctx context.Context, id string) error {
	if id == "" {
		return ErrEmptyPropertyID
	}

	ctx, cancel := context.WithTimeout(ctx, constants.DefaultTimeout)
	defer cancel()

	// Ensure the property exists before deleting
	docRef, err := r.client.Collection(propertyCollection).Doc(id).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil // Already deleted, return nil
		}
		return fmt.Errorf("checking property exists: %w", err)
	}

	_, err = docRef.Ref.Delete(ctx)
	if err != nil {
		return fmt.Errorf("deleting property: %w", err)
	}

	return nil
}
