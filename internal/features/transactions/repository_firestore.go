package transactions

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
	transactionCollection = "transactions"
)

var (
	ErrEmptyProjectID      = errors.New("project ID cannot be empty")
	ErrTransactionNotFound = errors.New("transaction not found")
	ErrNilTransaction      = errors.New("transaction cannot be nil")
	ErrEmptyTransactionID  = errors.New("transaction ID cannot be empty")
)

// FirestoreTransactionRepository implements the transactions.TransactionRepository interface.
type FirestoreTransactionRepository struct {
	client *firestore.Client
}

// NewFirestoreTransactionRepository creates a new FirestoreTransactionRepository.
func NewFirestoreTransactionRepository(ctx context.Context, projectID string, databaseID string) (*FirestoreTransactionRepository, error) {
	if projectID == "" {
		return nil, ErrEmptyProjectID // Reusing from repository_firestore.go
	}

	ctx, cancel := context.WithTimeout(ctx, constants.DefaultTimeout)
	defer cancel()

	client, err := firestore.NewClientWithDatabase(ctx, projectID, databaseID)
	if err != nil {
		return nil, fmt.Errorf("creating firestore client for transactions: %w", err)
	}

	return &FirestoreTransactionRepository{
		client: client,
	}, nil
}

// CloseClient closes the Firestore client.
func (r *FirestoreTransactionRepository) CloseClient() error {
	if r.client != nil {
		if err := r.client.Close(); err != nil {
			return fmt.Errorf("closing firestore client for transactions: %w", err)
		}
	}
	return nil
}

// Create implements transactions.TransactionRepository
func (r *FirestoreTransactionRepository) Create(ctx context.Context, transaction *entity.Transaction) (string, error) {
	if transaction == nil {
		return "", ErrNilTransaction
	}

	ctx, cancel := context.WithTimeout(ctx, constants.DefaultTimeout)
	defer cancel()

	now := time.Now().UTC()
	if transaction.CreatedAt.IsZero() {
		transaction.CreatedAt = now
	}
	if transaction.UpdatedAt.IsZero() {
		transaction.UpdatedAt = now
	}
	// ID is usually set by the service layer before calling Create,
	// or Firestore can auto-generate one if Add is used without a specific Doc ID.
	// The interface expects the ID to be returned. If service sets it, this is fine.
	// If Firestore generates it, Add() returns a DocumentRef from which ID can be obtained.

	docRef := r.client.Collection(transactionCollection).Doc(transaction.ID) // Assuming ID is pre-set by service
	_, err := docRef.Set(ctx, transaction)
	if err != nil {
		return "", fmt.Errorf("creating transaction %s: %w", transaction.ID, err)
	}
	return transaction.ID, nil
}

// GetByID implements transactions.TransactionRepository
func (r *FirestoreTransactionRepository) GetByID(ctx context.Context, id string) (*entity.Transaction, error) {
	if id == "" {
		return nil, ErrEmptyTransactionID
	}

	ctx, cancel := context.WithTimeout(ctx, constants.DefaultTimeout)
	defer cancel()

	doc, err := r.client.Collection(transactionCollection).Doc(id).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, ErrTransactionNotFound
		}
		return nil, fmt.Errorf("getting transaction by ID %s: %w", id, err)
	}

	var transaction entity.Transaction
	if err := doc.DataTo(&transaction); err != nil {
		return nil, fmt.Errorf("parsing transaction data for ID %s: %w", id, err)
	}
	transaction.ID = doc.Ref.ID
	return &transaction, nil
}

// List implements transactions.TransactionRepository
func (r *FirestoreTransactionRepository) List(
	ctx context.Context,
	propertyID *string,
	transactionType *entity.TransactionType,
	categoryID *string,
	startDate *time.Time,
	endDate *time.Time,
	limit, offset int,
) ([]*entity.Transaction, int, error) {
	ctx, cancel := context.WithTimeout(ctx, constants.DefaultTimeout)
	defer cancel()

	query := r.client.Collection(transactionCollection).Query

	if propertyID != nil && *propertyID != "" {
		query = query.Where("PropertyID", "==", *propertyID)
	}
	if transactionType != nil && *transactionType != "" {
		query = query.Where("Type", "==", *transactionType)
	}
	if categoryID != nil && *categoryID != "" {
		query = query.Where("CategoryID", "==", *categoryID)
	}
	if startDate != nil {
		query = query.Where("TransactionDate", ">=", *startDate)
	}
	if endDate != nil {
		query = query.Where("TransactionDate", "<=", *endDate)
	}

	// For total count with filters, a separate count query is needed.
	// Firestore's Count() aggregation can be applied to the filtered query.
	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, 0, fmt.Errorf("counting categories: %w", err)
	}
	totalItems := len(docs)

	// Apply ordering, limit, and offset for the data retrieval query
	// Order by TransactionDate descending by default, then by CreatedAt for tie-breaking
	dataQuery := query.OrderBy("TransactionDate", firestore.Desc).
		OrderBy("CreatedAt", firestore.Desc).
		Limit(limit).
		Offset(offset)

	iter := dataQuery.Documents(ctx)
	defer iter.Stop()

	var transactions []*entity.Transaction
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, 0, fmt.Errorf("iterating transactions for list: %w", err)
		}
		var t entity.Transaction
		if err := doc.DataTo(&t); err != nil {
			return nil, 0, fmt.Errorf("parsing transaction data for list: %w", err)
		}
		t.ID = doc.Ref.ID
		transactions = append(transactions, &t)
	}

	return transactions, totalItems, nil
}

// Update implements transactions.TransactionRepository
func (r *FirestoreTransactionRepository) Update(ctx context.Context, transaction *entity.Transaction) error {
	if transaction == nil || transaction.ID == "" {
		return ErrNilTransaction // Or ErrEmptyTransactionID
	}

	ctx, cancel := context.WithTimeout(ctx, constants.DefaultTimeout)
	defer cancel()

	docRef := r.client.Collection(transactionCollection).Doc(transaction.ID)
	transaction.UpdatedAt = time.Now().UTC()

	_, err := docRef.Set(ctx, transaction) // Set overwrites the document
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return ErrTransactionNotFound
		}
		return fmt.Errorf("updating transaction %s: %w", transaction.ID, err)
	}
	return nil
}

// Delete implements transactions.TransactionRepository
func (r *FirestoreTransactionRepository) Delete(ctx context.Context, id string) error {
	if id == "" {
		return ErrEmptyTransactionID
	}

	ctx, cancel := context.WithTimeout(ctx, constants.DefaultTimeout)
	defer cancel()

	docRef := r.client.Collection(transactionCollection).Doc(id)
	_, err := docRef.Delete(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return ErrTransactionNotFound // Or return nil if "already deleted" is acceptable
		}
		return fmt.Errorf("deleting transaction %s: %w", id, err)
	}
	return nil
}
