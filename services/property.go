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

func (s *PropertyService) GetPropertyByID(ctx context.Context, id string) (types.Property, error) {
	return s.repo.GetPropertyByID(ctx, id)
}

func (s *PropertyService) GetProperties(ctx context.Context, limit int, cursor string) ([]types.Property, string, error) {
	return s.repo.GetProperties(ctx, limit, cursor)
}

func (s *PropertyService) CreateProperty(ctx context.Context, property *types.Property) error {
	err := s.repo.CreateProperty(ctx, property)
	if err != nil {
		return err
	}
	return nil
}

func (s *PropertyService) UpdateProperty(ctx context.Context, id string, property *types.Property) error {
	return s.repo.UpdateProperty(ctx, id, property)
}

func (s *PropertyService) DeleteProperty(ctx context.Context, id string) error {
	return s.repo.DeleteProperty(ctx, id)
}
