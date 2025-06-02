package models

import "time"

type Property struct {
	ID          string    `json:"id,omitempty" firestore:"-"`
	Address     string    `json:"address" firestore:"address"`
	Postcode    string    `json:"postcode" firestore:"postcode"`
	Description string    `json:"description,omitempty" firestore:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at" firestore:"createdAt"`
	UpdatedAt   time.Time `json:"updated_at" firestore:"updatedAt"`
}
