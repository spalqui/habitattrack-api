package firestore

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"

	"github.com/spalqui/habitattrack-api/internal/models"
	"github.com/spalqui/habitattrack-api/internal/repositories"
)

type categoryRepository struct {
	client     *firestore.Client
	collection string
}

func NewCategoryRepository(client *firestore.Client) repositories.CategoryRepository {
	return &categoryRepository{
		client:     client,
		collection: "categories",
	}
}

func (r *categoryRepository) Create(ctx context.Context, category *models.Category) error {
	category.CreatedAt = time.Now()
	category.UpdatedAt = time.Now()

	docRef, _, err := r.client.Collection(r.collection).Add(ctx, category)
	if err != nil {
		return err
	}

	category.ID = docRef.ID
	return nil
}

func (r *categoryRepository) GetByID(ctx context.Context, id string) (*models.Category, error) {
	doc, err := r.client.Collection(r.collection).Doc(id).Get(ctx)
	if err != nil {
		return nil, err
	}

	var category models.Category
	if err := doc.DataTo(&category); err != nil {
		return nil, err
	}

	category.ID = doc.Ref.ID
	return &category, nil
}

func (r *categoryRepository) GetAll(ctx context.Context) ([]*models.Category, error) {
	docs, err := r.client.Collection(r.collection).Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	categories := make([]*models.Category, len(docs))
	for i, doc := range docs {
		var category models.Category
		if err := doc.DataTo(&category); err != nil {
			return nil, err
		}
		category.ID = doc.Ref.ID
		categories[i] = &category
	}

	return categories, nil
}

func (r *categoryRepository) GetByType(ctx context.Context, transactionType models.TransactionType) ([]*models.Category, error) {
	docs, err := r.client.Collection(r.collection).Where("type", "==", string(transactionType)).Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	categories := make([]*models.Category, len(docs))
	for i, doc := range docs {
		var category models.Category
		if err := doc.DataTo(&category); err != nil {
			return nil, err
		}
		category.ID = doc.Ref.ID
		categories[i] = &category
	}

	return categories, nil
}

func (r *categoryRepository) Update(ctx context.Context, category *models.Category) error {
	category.UpdatedAt = time.Now()
	_, err := r.client.Collection(r.collection).Doc(category.ID).Set(ctx, category)
	return err
}

func (r *categoryRepository) Delete(ctx context.Context, id string) error {
	_, err := r.client.Collection(r.collection).Doc(id).Delete(ctx)
	return err
}
