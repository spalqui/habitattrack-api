package repositories

import (
	"github.com/spalqui/habitattrack-api/types"
)

type PropertyRepository interface {
	GetPropertyByID(id string) (types.Property, error)
}

type FirestorePropertyRepository struct{}

func NewFirestorePropertyRepository() PropertyRepository {
	return &FirestorePropertyRepository{}
}

func (r *FirestorePropertyRepository) GetPropertyByID(id string) (types.Property, error) {
	return types.Property{}, nil
}
