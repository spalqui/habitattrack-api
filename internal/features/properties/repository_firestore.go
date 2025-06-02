package properties

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
	propertyCollection = "properties"
)

var (
	ErrEmptyProjectID   = errors.New("project ID cannot be empty")
	ErrPropertyNotFound = errors.New("property not found")
	ErrNilProperty      = errors.New("property cannot be nil")
	ErrEmptyPropertyID  = errors.New("property ID cannot be empty")
)

type FirestorePropertyRepository struct {
	client *firestore.Client
}

func NewFirestorePropertyRepository(ctx context.Context, projectID string, databaseID string) (*FirestorePropertyRepository, error) {
	if projectID == "" {
		return nil, ErrEmptyProjectID
	}

	ctx, cancel := context.WithTimeout(ctx, constants.DefaultTimeout)
	defer cancel()

	client, err := firestore.NewClientWithDatabase(ctx, projectID, databaseID)
	if err != nil {
		return nil, fmt.Errorf("creating firestore client for properties: %w", err)
	}

	return &FirestorePropertyRepository{
		client: client,
	}, nil
}

func (r *FirestorePropertyRepository) CloseClient() error {
	if r.client != nil {
		if err := r.client.Close(); err != nil {
			return fmt.Errorf("closing firestore client for properties: %w", err)
		}
	}
	return nil
}

func (r *FirestorePropertyRepository) Create(ctx context.Context, property *entity.Property) (string, error) {
	if property == nil {
		return "", ErrNilProperty
	}

	ctx, cancel := context.WithTimeout(ctx, constants.DefaultTimeout)
	defer cancel()

	now := time.Now().UTC()
	if property.CreatedAt.IsZero() {
		property.CreatedAt = now
	}
	if property.UpdatedAt.IsZero() {
		property.UpdatedAt = now
	}

	docRef := r.client.Collection(propertyCollection).Doc(property.ID) // Assuming ID is pre-set by service
	_, err := docRef.Set(ctx, property)
	if err != nil {
		return "", fmt.Errorf("creating property %s: %w", property.ID, err)
	}
	return property.ID, nil
}

func (r *FirestorePropertyRepository) GetByID(ctx context.Context, id string) (*entity.Property, error) {
	if id == "" {
		return nil, ErrEmptyPropertyID
	}

	ctx, cancel := context.WithTimeout(ctx, constants.DefaultTimeout)
	defer cancel()

	doc, err := r.client.Collection(propertyCollection).Doc(id).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, ErrPropertyNotFound
		}
		return nil, fmt.Errorf("getting property by ID %s: %w", id, err)
	}

	var property entity.Property
	if err := doc.DataTo(&property); err != nil {
		return nil, fmt.Errorf("parsing property data for ID %s: %w", id, err)
	}
	property.ID = doc.Ref.ID
	return &property, nil
}

func (r *FirestorePropertyRepository) List(ctx context.Context, limit, offset int) ([]*entity.Property, int, error) {
	ctx, cancel := context.WithTimeout(ctx, constants.DefaultTimeout)
	defer cancel()

	query := r.client.Collection(propertyCollection).Query

	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, 0, fmt.Errorf("counting properties: %w", err)
	}
	totalItems := len(docs)

	dataQuery := query.OrderBy("Name", firestore.Asc).Limit(limit).Offset(offset)
	iter := dataQuery.Documents(ctx)
	defer iter.Stop()

	var properties []*entity.Property
	for {
		doc, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return nil, 0, fmt.Errorf("iterating properties for list: %w", err)
		}
		var prop entity.Property
		if err := doc.DataTo(&prop); err != nil {
			return nil, 0, fmt.Errorf("parsing property data for list: %w", err)
		}
		prop.ID = doc.Ref.ID
		properties = append(properties, &prop)
	}
	return properties, totalItems, nil
}

func (r *FirestorePropertyRepository) Update(ctx context.Context, property *entity.Property) error {
	if property == nil || property.ID == "" {
		return ErrNilProperty
	}

	ctx, cancel := context.WithTimeout(ctx, constants.DefaultTimeout)
	defer cancel()

	docRef := r.client.Collection(propertyCollection).Doc(property.ID)
	property.UpdatedAt = time.Now().UTC()

	_, err := docRef.Set(ctx, property)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return ErrPropertyNotFound
		}
		return fmt.Errorf("updating property %s: %w", property.ID, err)
	}
	return nil
}

func (r *FirestorePropertyRepository) Delete(ctx context.Context, id string) error {
	if id == "" {
		return ErrEmptyPropertyID
	}

	ctx, cancel := context.WithTimeout(ctx, constants.DefaultTimeout)
	defer cancel()

	docRef := r.client.Collection(propertyCollection).Doc(id)
	_, err := docRef.Delete(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return ErrPropertyNotFound
		}
		return fmt.Errorf("deleting property %s: %w", id, err)
	}
	return nil
}

func (r *FirestorePropertyRepository) FindByName(ctx context.Context, name string) (*entity.Property, error) {
	if name == "" {
		return nil, errors.New("property name cannot be empty for FindByName")
	}

	ctx, cancel := context.WithTimeout(ctx, constants.DefaultTimeout)
	defer cancel()

	query := r.client.Collection(propertyCollection).Where("Name", "==", name).Limit(1)
	iter := query.Documents(ctx)
	defer iter.Stop()

	doc, err := iter.Next()
	if err == iterator.Done {
		return nil, ErrPropertyNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("querying property by name '%s': %w", name, err)
	}

	var property entity.Property
	if err := doc.DataTo(&property); err != nil {
		return nil, fmt.Errorf("parsing property data for FindByName: %w", err)
	}
	property.ID = doc.Ref.ID
	return &property, nil
}
