package firestore

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"

	"github.com/spalqui/habitattrack-api/internal/models"
	"github.com/spalqui/habitattrack-api/internal/repositories"
)

type propertyRepository struct {
	client     *firestore.Client
	collection string
}

func NewPropertyRepository(client *firestore.Client) repositories.PropertyRepository {
	return &propertyRepository{
		client:     client,
		collection: "properties",
	}
}

func (r *propertyRepository) Create(ctx context.Context, property *models.Property) error {
	property.CreatedAt = time.Now()
	property.UpdatedAt = time.Now()

	docRef, _, err := r.client.Collection(r.collection).Add(ctx, property)
	if err != nil {
		return err
	}

	property.ID = docRef.ID
	return nil
}

func (r *propertyRepository) GetByID(ctx context.Context, id string) (*models.Property, error) {
	doc, err := r.client.Collection(r.collection).Doc(id).Get(ctx)
	if err != nil {
		return nil, err
	}

	var property models.Property
	if err := doc.DataTo(&property); err != nil {
		return nil, err
	}

	property.ID = doc.Ref.ID
	return &property, nil
}

func (r *propertyRepository) GetAll(ctx context.Context) ([]*models.Property, error) {
	docs, err := r.client.Collection(r.collection).Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	properties := make([]*models.Property, len(docs))
	for i, doc := range docs {
		var property models.Property
		if err := doc.DataTo(&property); err != nil {
			return nil, err
		}
		property.ID = doc.Ref.ID
		properties[i] = &property
	}

	return properties, nil
}

func (r *propertyRepository) Update(ctx context.Context, property *models.Property) error {
	property.UpdatedAt = time.Now()
	_, err := r.client.Collection(r.collection).Doc(property.ID).Set(ctx, property)
	return err
}

func (r *propertyRepository) Delete(ctx context.Context, id string) error {
	_, err := r.client.Collection(r.collection).Doc(id).Delete(ctx)
	return err
}
