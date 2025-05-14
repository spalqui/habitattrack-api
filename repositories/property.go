package repositories

import (
	"context"
	"errors"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/spalqui/habitattrack-api/types"
)

const (
	propertyCollection = "properties"

	defaultTimeout = 30 * time.Second
)

var (
	ErrEmptyProjectID  = errors.New("invalid project ID")
	ErrEmptyPropertyID = errors.New("invalid property ID")
)

// PropertyRepository defines the interface for property repository
type PropertyRepository interface {
	GetPropertyByID(ctx context.Context, id string) (types.Property, error)
	Close() error
}

// FirestorePropertyRepository implements PropertyRepository using Firestore
type FirestorePropertyRepository struct {
	client *firestore.Client
}

func NewFirestorePropertyRepository(ctx context.Context, projectID string) (PropertyRepository, error) {
	if projectID == "" {
		return nil, ErrEmptyProjectID
	}

	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	client, err := firestore.NewClient(ctx, projectID)
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

	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	doc, err := r.client.Collection(propertyCollection).Doc(id).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return types.Property{}, nil
		}
		return types.Property{}, fmt.Errorf("getting property by ID: %w", err)
	}

	var property types.Property
	err = doc.DataTo(&property)
	if err != nil {
		return types.Property{}, fmt.Errorf("parsing property data: %w", err)
	}

	return types.Property{}, nil
}
