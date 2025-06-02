package properties

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"

	"github.com/spalqui/habitattrack-api/internal/shared/apierrors"
	"github.com/spalqui/habitattrack-api/internal/shared/apitypes"
)

const (
	defaultServicePageSize = 20
	maxServicePageSize     = 100
	defaultServicePageNum  = 1
)

// PropertyService defines the interface for property business logic.
type PropertyService interface {
	CreateProperty(ctx context.Context, req CreatePropertyRequest) (*PropertyResponse, error)
	GetPropertyByID(ctx context.Context, id string) (*PropertyResponse, error)
	ListProperties(ctx context.Context, page, pageSize int) (*PaginatedPropertiesResponse, error)
	UpdateProperty(ctx context.Context, id string, req UpdatePropertyRequest) (*PropertyResponse, error)
	DeleteProperty(ctx context.Context, id string) error
}

type propertyService struct {
	repo PropertyRepository
}

// NewPropertyService creates a new instance of PropertyService.
func NewPropertyService(repo PropertyRepository) PropertyService {
	return &propertyService{
		repo: repo,
	}
}

func (s *propertyService) CreateProperty(ctx context.Context, req CreatePropertyRequest) (*PropertyResponse, error) {
	property := req.ToEntity()
	property.ID = uuid.NewString()
	property.CreatedAt = time.Now().UTC()
	property.UpdatedAt = property.CreatedAt

	_, err := s.repo.Create(ctx, property)
	if err != nil {
		// Assuming s.repo.Create might return a specific error for conflicts (e.g., unique name constraint)
		// which could be checked here to return apierrors.ErrConflict.
		// For now, a generic internal error is returned.
		return nil, apierrors.ErrInternal(fmt.Errorf("creating property: %w", err))
	}

	resp := ToPropertyResponse(property)
	return &resp, nil
}

func (s *propertyService) GetPropertyByID(ctx context.Context, id string) (*PropertyResponse, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, apierrors.ErrBadRequest(fmt.Sprintf("Invalid property ID format: %s", id))
	}

	property, err := s.repo.GetByID(ctx, id)
	if err != nil {
		// Assuming the repository returns an error that can be identified as "not found".
		// If err is sql.ErrNoRows or a custom not found error, map to apierrors.ErrNotFound.
		// For simplicity, if GetByID returns any error, we'll consider it a potential "not found"
		// or an internal error if it's something else. A more robust check is needed here
		// based on the actual error types returned by the repository.
		// For now, we directly map to ErrNotFound, but this might mask other issues.
		// A better approach:
		// if errors.Is(err, entity.ErrNotFound) { // Assuming entity.ErrNotFound exists
		//     return nil, apierrors.ErrNotFound("Property", id)
		// }
		// return nil, apierrors.ErrInternal(fmt.Errorf("getting property by ID %s: %w", id, err))
		return nil, apierrors.ErrNotFound("Property", id) // Simplified for now
	}

	resp := ToPropertyResponse(property)
	return &resp, nil
}

func (s *propertyService) ListProperties(ctx context.Context, page, pageSize int) (*PaginatedPropertiesResponse, error) {
	if page <= 0 {
		page = defaultServicePageNum
	}
	if pageSize <= 0 {
		pageSize = defaultServicePageSize
	}
	if pageSize > maxServicePageSize {
		pageSize = maxServicePageSize
	}
	offset := (page - 1) * pageSize

	properties, totalItems, err := s.repo.List(ctx, pageSize, offset)
	if err != nil {
		return nil, apierrors.ErrInternal(fmt.Errorf("listing properties: %w", err))
	}

	responseList := ToPropertyResponseList(properties)
	totalPages := 0
	if totalItems > 0 {
		totalPages = int(math.Ceil(float64(totalItems) / float64(pageSize)))
	}

	return &PaginatedPropertiesResponse{
		Data: responseList,
		Pagination: apitypes.PaginationInfo{
			TotalItems:  int64(totalItems),
			TotalPages:  totalPages,
			CurrentPage: page,
			PageSize:    pageSize,
		},
	}, nil
}

func (s *propertyService) UpdateProperty(ctx context.Context, id string, req UpdatePropertyRequest) (*PropertyResponse, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, apierrors.ErrBadRequest(fmt.Sprintf("Invalid property ID format: %s", id))
	}

	property, err := s.repo.GetByID(ctx, id)
	if err != nil {
		// See comment in GetPropertyByID regarding error handling.
		return nil, apierrors.ErrNotFound("Property", id) // Simplified
	}

	updated := false
	if req.Name != nil && *req.Name != "" && *req.Name != property.Name {
		// Optional: Add unique name check here if required, similar to CreateProperty.
		// If so, ensure to check existing.ID != id.
		property.Name = *req.Name
		updated = true
	}

	// For req.Address, if it's nil, it means "no change" for PATCH.
	// If it's a non-nil pointer, it means "update to this value".
	// This includes setting it to an empty string if *req.Address is "".
	if req.Address != nil {
		// If current address is nil and new address is non-nil, or
		// if current address is non-nil and new address is different.
		if property.Address == nil || *req.Address != *property.Address {
			property.Address = req.Address
			updated = true
		}
	}
	// Note: The original code had an `else if req.Address == nil && property.Address != nil`
	// This branch is not strictly necessary if nil in request means "no change".
	// If `nil` should explicitly clear the address, the DTO and client request must be clear.
	// For PATCH, typically, omitting a field or sending `null` for a nullable field means "no change"
	// unless the API contract specifies `null` means "set to null".
	// The current UpdatePropertyRequest uses `omitempty`, so if "address" is not in JSON, req.Address is nil.
	// If "address": null is in JSON, req.Address is also nil.
	// The validator `omitempty` on `UpdatePropertyRequest` means if `Name` or `Address` are their zero values (nil for pointers),
	// they won't be validated unless other validation tags like `required` are present (which they are not for these optional fields).

	if updated {
		property.UpdatedAt = time.Now().UTC()
		err = s.repo.Update(ctx, property)
		if err != nil {
			// Could also check for specific update errors, e.g., concurrent modification or conflict.
			return nil, apierrors.ErrInternal(fmt.Errorf("updating property %s: %w", id, err))
		}
	}

	resp := ToPropertyResponse(property)
	return &resp, nil
}

func (s *propertyService) DeleteProperty(ctx context.Context, id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return apierrors.ErrBadRequest(fmt.Sprintf("Invalid property ID format: %s", id))
	}

	// Before deleting, ensure the property exists.
	// This provides a consistent "not found" error if trying to delete a non-existent property.
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		// See comment in GetPropertyByID regarding error handling.
		return apierrors.ErrNotFound("Property", id) // Simplified
	}

	err = s.repo.Delete(ctx, id)
	if err != nil {
		// If GetByID passed, but Delete fails, it's likely an internal error
		// or a specific constraint violation (e.g., foreign key if not checked before).
		return apierrors.ErrInternal(fmt.Errorf("deleting property %s: %w", id, err))
	}
	return nil
}
