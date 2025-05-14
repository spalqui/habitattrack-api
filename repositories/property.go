package repositories

import (
	"context"

	"cloud.google.com/go/firestore"

	"github.com/spalqui/habitattrack-api/types"
)

type PropertyRepository interface {
	GetPropertyByID(id string) (types.Property, error)
}

type FirestorePropertyRepository struct {
	client *firestore.Client
}

func NewFirestorePropertyRepository(googleCloudProjectID string) (PropertyRepository, error) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, googleCloudProjectID)
	if err != nil {
		return &FirestorePropertyRepository{}, err
	}
	return &FirestorePropertyRepository{
		client: client,
	}, nil
}

func (r *FirestorePropertyRepository) GetPropertyByID(id string) (types.Property, error) {
	return types.Property{}, nil
}
