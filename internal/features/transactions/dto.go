package transactions

import (
	"time"

	"github.com/spalqui/habitattrack-api/internal/core/entity"
	"github.com/spalqui/habitattrack-api/internal/shared/apitypes"
)

// CreateTransactionRequest defines the structure for creating a new transaction.
// Maps to TransactionCreateRequest in OpenAPI.
type CreateTransactionRequest struct {
	Amount          float64                `json:"amount" validate:"required,gt=0"`
	TransactionDate time.Time              `json:"transactionDate" validate:"required"`
	Description     *string                `json:"description,omitempty" validate:"omitempty,max=500"`
	Type            entity.TransactionType `json:"type" validate:"required,oneof=income expense"`
	CategoryID      string                 `json:"categoryId" validate:"required,uuid"`
	PropertyID      *string                `json:"propertyId,omitempty" validate:"omitempty,uuid"`
}

// UpdateTransactionRequest defines the structure for updating an existing transaction.
// Maps to TransactionUpdateRequest in OpenAPI.
// Uses pointers to distinguish between a field not being provided and a field being provided with an empty/zero value.
type UpdateTransactionRequest struct {
	Amount          *float64               `json:"amount,omitempty" validate:"omitempty,gt=0"`
	TransactionDate *time.Time             `json:"transactionDate,omitempty"`
	Description     *string                `json:"description,omitempty" validate:"omitempty,max=500"`
	Type            *entity.TransactionType `json:"type,omitempty" validate:"omitempty,oneof=income expense"`
	CategoryID      *string                `json:"categoryId,omitempty" validate:"omitempty,uuid"`
	PropertyID      *string                `json:"propertyId,omitempty" validate:"omitempty,uuid"`
}

// TransactionResponse defines the structure for a transaction API response.
// Maps to Transaction in OpenAPI.
type TransactionResponse struct {
	ID              string                 `json:"id"`
	Amount          float64                `json:"amount"`
	TransactionDate time.Time              `json:"transactionDate"`
	Description     *string                `json:"description,omitempty"`
	Type            entity.TransactionType `json:"type"`
	CategoryID      string                 `json:"categoryId"`
	PropertyID      *string                `json:"propertyId,omitempty"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
}

// PaginatedTransactionsResponse defines the structure for a paginated list of transactions.
// Maps to PaginatedTransactionsResponse in OpenAPI.
type PaginatedTransactionsResponse struct {
	Data       []TransactionResponse       `json:"data"`
	Pagination apitypes.PaginationInfo `json:"pagination"`
}

// ToEntity converts a CreateTransactionRequest DTO to an entity.Transaction.
func (cr *CreateTransactionRequest) ToEntity() *entity.Transaction {
	return &entity.Transaction{
		Amount:          cr.Amount,
		TransactionDate: cr.TransactionDate,
		Description:     cr.Description,
		Type:            cr.Type,
		CategoryID:      cr.CategoryID,
		PropertyID:      cr.PropertyID,
	}
}

// ToTransactionResponse converts an entity.Transaction to a TransactionResponse DTO.
func ToTransactionResponse(t *entity.Transaction) TransactionResponse {
	return TransactionResponse{
		ID:              t.ID,
		Amount:          t.Amount,
		TransactionDate: t.TransactionDate,
		Description:     t.Description,
		Type:            t.Type,
		CategoryID:      t.CategoryID,
		PropertyID:      t.PropertyID,
		CreatedAt:       t.CreatedAt,
		UpdatedAt:       t.UpdatedAt,
	}
}

// ToTransactionResponseList converts a slice of entity.Transaction to a slice of TransactionResponse DTOs.
func ToTransactionResponseList(transactions []*entity.Transaction) []TransactionResponse {
	res := make([]TransactionResponse, len(transactions))
	for i, t := range transactions {
		res[i] = ToTransactionResponse(t)
	}
	return res
}