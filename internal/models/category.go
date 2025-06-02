package models

import "time"

type Category struct {
	ID          string          `json:"id,omitempty" firestore:"-"`
	Name        string          `json:"name" firestore:"name"`
	Type        TransactionType `json:"type" firestore:"type"`
	Description string          `json:"description,omitempty" firestore:"description,omitempty"`
	CreatedAt   time.Time       `json:"created_at" firestore:"createdAt"`
	UpdatedAt   time.Time       `json:"updated_at" firestore:"updatedAt"`
}
