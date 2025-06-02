package entity

import "time"

// ClassificationType defines the type of category (income or expense).
type ClassificationType string

const (
	IncomeClassification  ClassificationType = "income"
	ExpenseClassification ClassificationType = "expense"
)

// TransactionCategory represents a category for financial transactions.
type TransactionCategory struct {
	ID             string             `json:"id" db:"id"`
	Name           string             `json:"name" db:"name"`
	Classification ClassificationType `json:"classification" db:"classification"`
	CreatedAt      time.Time          `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at" db:"updated_at"`
}