package services

import (
	"context"

	"github.com/spalqui/habitattrack-api/repositories"
	"github.com/spalqui/habitattrack-api/types"
)

type PropertyService struct {
	repo repositories.PropertyRepository
}

func NewPropertyService(repo repositories.PropertyRepository) *PropertyService {
	return &PropertyService{
		repo: repo,
	}
}

func (s *PropertyService) GetPropertyByID(id string) (types.Property, error) {
	return s.repo.GetPropertyByID(context.Background(), id)
}
