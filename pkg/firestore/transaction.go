package firestore

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"

	"github.com/spalqui/habitattrack-api/internal/models"
	"github.com/spalqui/habitattrack-api/internal/repositories"
)

type transactionRepository struct {
	client     *firestore.Client
	collection string
}

func NewTransactionRepository(client *firestore.Client) repositories.TransactionRepository {
	return &transactionRepository{
		client:     client,
		collection: "transactions",
	}
}

func (r *transactionRepository) Create(ctx context.Context, transaction *models.Transaction) error {
	transaction.CreatedAt = time.Now()
	transaction.UpdatedAt = time.Now()

	docRef, _, err := r.client.Collection(r.collection).Add(ctx, transaction)
	if err != nil {
		return err
	}

	transaction.ID = docRef.ID
	return nil
}

func (r *transactionRepository) GetByID(ctx context.Context, id string) (*models.Transaction, error) {
	doc, err := r.client.Collection(r.collection).Doc(id).Get(ctx)
	if err != nil {
		return nil, err
	}

	var transaction models.Transaction
	if err := doc.DataTo(&transaction); err != nil {
		return nil, err
	}

	transaction.ID = doc.Ref.ID
	return &transaction, nil
}

func (r *transactionRepository) GetByPropertyID(ctx context.Context, propertyID string) ([]*models.Transaction, error) {
	docs, err := r.client.Collection(r.collection).Where("propertyId", "==", propertyID).Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	transactions := make([]*models.Transaction, len(docs))
	for i, doc := range docs {
		var transaction models.Transaction
		if err := doc.DataTo(&transaction); err != nil {
			return nil, err
		}
		transaction.ID = doc.Ref.ID
		transactions[i] = &transaction
	}

	return transactions, nil
}

func (r *transactionRepository) GetAll(ctx context.Context) ([]*models.Transaction, error) {
	docs, err := r.client.Collection(r.collection).Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	transactions := make([]*models.Transaction, len(docs))
	for i, doc := range docs {
		var transaction models.Transaction
		if err := doc.DataTo(&transaction); err != nil {
			return nil, err
		}
		transaction.ID = doc.Ref.ID
		transactions[i] = &transaction
	}

	return transactions, nil
}

func (r *transactionRepository) Update(ctx context.Context, transaction *models.Transaction) error {
	transaction.UpdatedAt = time.Now()
	_, err := r.client.Collection(r.collection).Doc(transaction.ID).Set(ctx, transaction)
	return err
}

func (r *transactionRepository) Delete(ctx context.Context, id string) error {
	_, err := r.client.Collection(r.collection).Doc(id).Delete(ctx)
	return err
}
