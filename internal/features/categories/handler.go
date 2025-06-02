package categories

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/spalqui/habitattrack-api/internal/core/entity"
	"github.com/spalqui/habitattrack-api/internal/shared/apierrors"
)

const (
	defaultPage         = 1
	defaultPageSize     = 20
	queryPage           = "page"
	queryLimit          = "limit"
	queryClassification = "classification"
)

// CategoryHandler handles HTTP requests for transaction categories.
type CategoryHandler struct {
	service  TransactionCategoryService
	validate *validator.Validate // Validator instance
}

// NewCategoryHandler creates a new instance of CategoryHandler.
func NewCategoryHandler(service TransactionCategoryService) *CategoryHandler {
	return &CategoryHandler{
		service:  service,
		validate: validator.New(), // Initialize validator
	}
}

// RegisterRoutes sets up the routes for category operations.
// It's a common pattern to have a method like this to be called from main.go or a router setup file.
func (h *CategoryHandler) RegisterRoutes(router *gin.Engine) {
	categoriesGroup := router.Group("/transaction-categories")
	{
		categoriesGroup.POST("", h.createCategory)
		categoriesGroup.GET("", h.listCategories)
		categoriesGroup.GET("/:categoryId", h.getCategoryByID)
		categoriesGroup.PUT("/:categoryId", h.updateCategory)  // Handles full update
		categoriesGroup.PATCH("/:categoryId", h.patchCategory) // Handles partial update
		categoriesGroup.DELETE("/:categoryId", h.deleteCategory)
	}
}

// bindAndValidate binds the request JSON to the given struct and validates it.
// It returns true if binding or validation fails and the response has been written.
func (h *CategoryHandler) bindAndValidate(c *gin.Context, req interface{}) bool {
	if err := c.ShouldBindJSON(req); err != nil {
		appErr := apierrors.ErrBadRequest("Invalid request payload: " + err.Error())
		c.JSON(http.StatusBadRequest, appErr)
		return true
	}

	if err := h.validate.Struct(req); err != nil {
		appErr := apierrors.ErrValidation("Validation failed", h.formatValidationErrors(err))
		c.JSON(http.StatusBadRequest, appErr)
		return true
	}
	return false
}

// createCategory handles POST /transaction-categories
func (h *CategoryHandler) createCategory(c *gin.Context) {
	var req CreateCategoryRequest
	if h.bindAndValidate(c, &req) {
		return
	}

	categoryResp, err := h.service.CreateCategory(c.Request.Context(), req)
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, categoryResp)
}

// listCategories handles GET /transaction-categories
func (h *CategoryHandler) listCategories(c *gin.Context) {
	pageStr := c.DefaultQuery(queryPage, strconv.Itoa(defaultPage))
	pageSizeStr := c.DefaultQuery(queryLimit, strconv.Itoa(defaultPageSize))

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = defaultPage
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = defaultPageSize
	}

	var classificationFilter *entity.ClassificationType
	if classificationQuery, ok := c.GetQuery(queryClassification); ok {
		cf := entity.ClassificationType(classificationQuery)
		if cf == entity.IncomeClassification || cf == entity.ExpenseClassification {
			classificationFilter = &cf
		} else {
			appErr := apierrors.ErrBadRequest("Invalid 'classification' query parameter. Must be 'income' or 'expense'.")
			c.JSON(http.StatusBadRequest, appErr)
			return
		}
	}

	categoriesResp, err := h.service.ListCategories(c.Request.Context(), classificationFilter, page, pageSize)
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, categoriesResp)
}

// getCategoryByID handles GET /transaction-categories/:categoryId
func (h *CategoryHandler) getCategoryByID(c *gin.Context) {
	categoryID := c.Param("categoryId")

	categoryResp, err := h.service.GetCategoryByID(c.Request.Context(), categoryID)
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, categoryResp)
}

// updateCategory handles PUT /transaction-categories/:categoryId (full update)
func (h *CategoryHandler) updateCategory(c *gin.Context) {
	categoryID := c.Param("categoryId")
	var req CreateCategoryRequest // For PUT, expect all fields
	if h.bindAndValidate(c, &req) {
		return
	}

	// Convert CreateCategoryRequest to UpdateCategoryRequest for the service.
	// This assumes the service's UpdateCategory method expects UpdateCategoryRequest
	// for consistency, even if PUT implies all fields are present.
	updateReq := UpdateCategoryRequest{
		Name:           &req.Name,
		Classification: &req.Classification,
	}

	categoryResp, err := h.service.UpdateCategory(c.Request.Context(), categoryID, updateReq)
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, categoryResp)
}

// patchCategory handles PATCH /transaction-categories/:categoryId (partial update)
func (h *CategoryHandler) patchCategory(c *gin.Context) {
	categoryID := c.Param("categoryId")
	var req UpdateCategoryRequest // For PATCH, use UpdateCategoryRequest with pointers for partial updates
	if h.bindAndValidate(c, &req) {
		return
	}

	categoryResp, err := h.service.UpdateCategory(c.Request.Context(), categoryID, req)
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, categoryResp)
}

// deleteCategory handles DELETE /transaction-categories/:categoryId
func (h *CategoryHandler) deleteCategory(c *gin.Context) {
	categoryID := c.Param("categoryId")

	err := h.service.DeleteCategory(c.Request.Context(), categoryID)
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// handleError centralizes error handling for the handler.
func (h *CategoryHandler) handleError(c *gin.Context, err error) {
	if appErr, ok := err.(*apierrors.Error); ok {
		switch appErr.Code {
		case apierrors.CodeNotFound:
			c.JSON(http.StatusNotFound, appErr)
		case apierrors.CodeConflict:
			c.JSON(http.StatusConflict, appErr)
		case apierrors.CodeValidationError, apierrors.CodeBadRequest, apierrors.CodeUnprocessableEntity:
			if appErr.Code == apierrors.CodeUnprocessableEntity {
				c.JSON(http.StatusUnprocessableEntity, appErr)
			} else {
				c.JSON(http.StatusBadRequest, appErr)
			}
		case apierrors.CodeUnauthorized:
			c.JSON(http.StatusUnauthorized, appErr)
		case apierrors.CodeForbidden:
			c.JSON(http.StatusForbidden, appErr)
		default:
			c.JSON(http.StatusInternalServerError, appErr)
		}
		return
	}
	// Fallback for unexpected errors
	c.JSON(http.StatusInternalServerError, apierrors.ErrInternal(err))
}

// formatValidationErrors converts validator.ValidationErrors to a map.
func (h *CategoryHandler) formatValidationErrors(err error) map[string]string {
	errors := make(map[string]string)
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldErr := range validationErrors {
			fieldName := fieldErr.Field()
			switch fieldErr.Tag() {
			case "required":
				errors[fieldName] = fieldName + " is required."
			case "min":
				errors[fieldName] = fieldName + " must be at least " + fieldErr.Param() + " characters."
			case "max":
				errors[fieldName] = fieldName + " must not exceed " + fieldErr.Param() + " characters."
			case "oneof":
				errors[fieldName] = fieldName + " must be one of [" + fieldErr.Param() + "]."
			case "uuid":
				errors[fieldName] = fieldName + " must be a valid UUID."
			default:
				errors[fieldName] = fieldName + " is not valid (" + fieldErr.Tag() + ")."
			}
		}
	}
	return errors
}
