package models

import "time"

type TransactionType string

const (
	TransactionTypeIncome  TransactionType = "income"
	TransactionTypeExpense TransactionType = "expense"
)

type Transaction struct {
	ID          string          `json:"id,omitempty" firestore:"-"`
	PropertyID  string          `json:"property_id" firestore:"propertyId"`
	Type        TransactionType `json:"type" firestore:"type"`
	CategoryID  string          `json:"category_id" firestore:"categoryId"`
	Amount      float64         `json:"amount" firestore:"amount"`
	Description string          `json:"description,omitempty" firestore:"description,omitempty"`
	Date        time.Time       `json:"date" firestore:"date"`
	CreatedAt   time.Time       `json:"created_at" firestore:"createdAt"`
	UpdatedAt   time.Time       `json:"updated_at" firestore:"updatedAt"`
}
