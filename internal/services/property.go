package services

import (
	"context"
	"errors"
	"strings"

	"github.com/spalqui/habitattrack-api/internal/models"
	"github.com/spalqui/habitattrack-api/internal/repositories"
)

type PropertyService interface {
	CreateProperty(ctx context.Context, property *models.Property) error
	GetProperty(ctx context.Context, id string) (*models.Property, error)
	GetAllProperties(ctx context.Context) ([]*models.Property, error)
	UpdateProperty(ctx context.Context, property *models.Property) error
	DeleteProperty(ctx context.Context, id string) error
}

type propertyService struct {
	propertyRepo repositories.PropertyRepository
}

func NewPropertyService(propertyRepo repositories.PropertyRepository) PropertyService {
	return &propertyService{
		propertyRepo: propertyRepo,
	}
}

func (s *propertyService) CreateProperty(ctx context.Context, property *models.Property) error {
	if err := s.validateProperty(property); err != nil {
		return err
	}

	return s.propertyRepo.Create(ctx, property)
}

func (s *propertyService) GetProperty(ctx context.Context, id string) (*models.Property, error) {
	if strings.TrimSpace(id) == "" {
		return nil, errors.New("property ID is required")
	}

	return s.propertyRepo.GetByID(ctx, id)
}

func (s *propertyService) GetAllProperties(ctx context.Context) ([]*models.Property, error) {
	return s.propertyRepo.GetAll(ctx)
}

func (s *propertyService) UpdateProperty(ctx context.Context, property *models.Property) error {
	if err := s.validateProperty(property); err != nil {
		return err
	}

	if strings.TrimSpace(property.ID) == "" {
		return errors.New("property ID is required for update")
	}

	return s.propertyRepo.Update(ctx, property)
}

func (s *propertyService) DeleteProperty(ctx context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		return errors.New("property ID is required")
	}

	return s.propertyRepo.Delete(ctx, id)
}

func (s *propertyService) validateProperty(property *models.Property) error {
	if strings.TrimSpace(property.Address) == "" {
		return errors.New("address is required")
	}

	if strings.TrimSpace(property.Postcode) == "" {
		return errors.New("postcode is required")
	}

	return nil
}
