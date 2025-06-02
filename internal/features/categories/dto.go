package categories

import (
	"time"

	"github.com/spalqui/habitattrack-api/internal/core/entity"
	"github.com/spalqui/habitattrack-api/internal/shared/apitypes"
)

// CreateCategoryRequest defines the structure for creating a new transaction category.
// It maps to TransactionCategoryCreateRequest in OpenAPI.
type CreateCategoryRequest struct {
	Name           string                      `json:"name" validate:"required,min=1,max=100"`
	Classification entity.ClassificationType `json:"classification" validate:"required,oneof=income expense"`
}

// UpdateCategoryRequest defines the structure for updating an existing transaction category.
// It maps to TransactionCategoryUpdateRequest in OpenAPI.
// Uses pointers to distinguish between a field not being provided and a field being provided with an empty/zero value.
type UpdateCategoryRequest struct {
	Name           *string                     `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	Classification *entity.ClassificationType `json:"classification,omitempty" validate:"omitempty,oneof=income expense"`
}

// CategoryResponse defines the structure for a transaction category API response.
// It maps to TransactionCategory in OpenAPI.
type CategoryResponse struct {
	ID             string                      `json:"id"`
	Name           string                      `json:"name"`
	Classification entity.ClassificationType `json:"classification"`
	CreatedAt      time.Time                   `json:"created_at"`
	UpdatedAt      time.Time                   `json:"updated_at"`
}

// PaginatedCategoriesResponse defines the structure for a paginated list of categories.
// It maps to PaginatedCategoriesResponse in OpenAPI.
type PaginatedCategoriesResponse struct {
	Data       []CategoryResponse          `json:"data"`
	Pagination apitypes.PaginationInfo `json:"pagination"`
}

// ToEntity converts a CreateCategoryRequest DTO to an entity.TransactionCategory.
func (cr *CreateCategoryRequest) ToEntity() *entity.TransactionCategory {
	return &entity.TransactionCategory{
		Name:           cr.Name,
		Classification: cr.Classification,
	}
}

// ToCategoryResponse converts an entity.TransactionCategory to a CategoryResponse DTO.
func ToCategoryResponse(cat *entity.TransactionCategory) CategoryResponse {
	return CategoryResponse{
		ID:             cat.ID,
		Name:           cat.Name,
		Classification: cat.Classification,
		CreatedAt:      cat.CreatedAt,
		UpdatedAt:      cat.UpdatedAt,
	}
}

// ToCategoryResponseList converts a slice of entity.TransactionCategory to a slice of CategoryResponse DTOs.
func ToCategoryResponseList(categories []*entity.TransactionCategory) []CategoryResponse {
	res := make([]CategoryResponse, len(categories))
	for i, cat := range categories {
		res[i] = ToCategoryResponse(cat)
	}
	return res
}