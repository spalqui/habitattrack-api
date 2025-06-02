package entity

import "time"

// Property represents a real estate asset.
type Property struct {
	ID        string    `json:"id" db:"id"` // Using db tag for sqlx or similar ORM
	Name      string    `json:"name" db:"name"`
	Address   *string   `json:"address,omitempty" db:"address"` // Pointer for optional field
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}