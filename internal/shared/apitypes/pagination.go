package apitypes

// PaginationInfo defines the structure for pagination metadata in API responses.
// It maps to PaginationInfo in OpenAPI.
type PaginationInfo struct {
	TotalItems  int64  `json:"totalItems"`
	TotalPages  int    `json:"totalPages"`
	CurrentPage int    `json:"currentPage"`
	PageSize    int    `json:"pageSize"`
	NextCursor  string `json:"nextCursor,omitempty"` // Used for cursor-based pagination
	PrevCursor  string `json:"prevCursor,omitempty"` // Used for cursor-based pagination
}