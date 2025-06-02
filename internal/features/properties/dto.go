package properties

import (
	"github.com/spalqui/habitattrack-api/internal/core/entity"
	"github.com/spalqui/habitattrack-api/internal/shared/apitypes"
	"time"
)

// CreatePropertyRequest defines the structure for creating a new property.
// Maps to PropertyCreateRequest in OpenAPI.
type CreatePropertyRequest struct {
	Name    string  `json:"name" validate:"required,min=1,max=255"`
	Address *string `json:"address,omitempty" validate:"omitempty,max=500"`
}

// UpdatePropertyRequest defines the structure for updating an existing property.
// Maps to PropertyUpdateRequest in OpenAPI.
type UpdatePropertyRequest struct {
	Name    *string `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	Address *string `json:"address,omitempty" validate:"omitempty,max=500"`
}

// PropertyResponse defines the structure for a property API response.
// Maps to Property in OpenAPI.
type PropertyResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Address   *string   `json:"address,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// PaginatedPropertiesResponse defines the structure for a paginated list of properties.
// Maps to PaginatedPropertiesResponse in OpenAPI.
type PaginatedPropertiesResponse struct {
	Data       []PropertyResponse          `json:"data"`
	Pagination apitypes.PaginationInfo `json:"pagination"`
}

// ToEntity converts a CreatePropertyRequest DTO to an entity.Property.
func (cr *CreatePropertyRequest) ToEntity() *entity.Property {
	return &entity.Property{
		Name:    cr.Name,
		Address: cr.Address,
	}
}

// ToPropertyResponse converts an entity.Property to a PropertyResponse DTO.
func ToPropertyResponse(p *entity.Property) PropertyResponse {
	return PropertyResponse{
		ID:        p.ID,
		Name:      p.Name,
		Address:   p.Address,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

// ToPropertyResponseList converts a slice of entity.Property to a slice of PropertyResponse DTOs.
func ToPropertyResponseList(properties []*entity.Property) []PropertyResponse {
	res := make([]PropertyResponse, len(properties))
	for i, p := range properties {
		res[i] = ToPropertyResponse(p)
	}
	return res
}