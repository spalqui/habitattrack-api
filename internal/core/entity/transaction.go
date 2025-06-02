package entity

import "time"

// TransactionType defines the type of transaction (income or expense).
type TransactionType string

const (
	IncomeTransaction  TransactionType = "income"
	ExpenseTransaction TransactionType = "expense"
)

// Transaction represents a financial transaction.
type Transaction struct {
	ID              string          `json:"id" db:"id"`
	Amount          float64         `json:"amount" db:"amount"` // Consider using a dedicated money type for precision
	TransactionDate time.Time       `json:"transactionDate" db:"transaction_date"`
	Description     *string         `json:"description,omitempty" db:"description"` // Pointer for optional field
	Type            TransactionType `json:"type" db:"type"`
	CategoryID      string          `json:"categoryId" db:"category_id"`
	PropertyID      *string         `json:"propertyId,omitempty" db:"property_id"` // Pointer for optional field
	CreatedAt       time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at" db:"updated_at"`
}