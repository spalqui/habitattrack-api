package properties

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/spalqui/habitattrack-api/internal/shared/apierrors"
)

const (
	defaultPage     = 1
	defaultPageSize = 20
	queryPage       = "page"
	queryLimit      = "limit" // Corresponds to 'limit' in OpenAPI for page size
)

// PropertyHandler handles HTTP requests for properties.
type PropertyHandler struct {
	service  PropertyService
	validate *validator.Validate
}

// NewPropertyHandler creates a new instance of PropertyHandler.
func NewPropertyHandler(service PropertyService) *PropertyHandler {
	return &PropertyHandler{
		service:  service,
		validate: validator.New(),
	}
}

// RegisterRoutes sets up the routes for property operations.
func (h *PropertyHandler) RegisterRoutes(router *gin.Engine) {
	propertiesGroup := router.Group("/properties")
	{
		propertiesGroup.POST("", h.createProperty)
		propertiesGroup.GET("", h.listProperties)
		propertiesGroup.GET("/:propertyId", h.getPropertyByID)
		propertiesGroup.PUT("/:propertyId", h.updateProperty)
		propertiesGroup.PATCH("/:propertyId", h.patchProperty)
		propertiesGroup.DELETE("/:propertyId", h.deleteProperty)
	}
}

// bindAndValidate binds the request JSON to the given struct and validates it.
// It returns true if binding or validation fails and the response has been written.
func (h *PropertyHandler) bindAndValidate(c *gin.Context, req interface{}) bool {
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

// createProperty handles POST /properties
func (h *PropertyHandler) createProperty(c *gin.Context) {
	var req CreatePropertyRequest
	if h.bindAndValidate(c, &req) {
		return
	}

	propertyResp, err := h.service.CreateProperty(c.Request.Context(), req)
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, propertyResp)
}

// listProperties handles GET /properties
func (h *PropertyHandler) listProperties(c *gin.Context) {
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

	propertiesResp, err := h.service.ListProperties(c.Request.Context(), page, pageSize)
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, propertiesResp)
}

// getPropertyByID handles GET /properties/:propertyId
func (h *PropertyHandler) getPropertyByID(c *gin.Context) {
	propertyID := c.Param("propertyId")
	propertyResp, err := h.service.GetPropertyByID(c.Request.Context(), propertyID)
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, propertyResp)
}

// updateProperty handles PUT /properties/:propertyId
func (h *PropertyHandler) updateProperty(c *gin.Context) {
	propertyID := c.Param("propertyId")
	// For PUT, OpenAPI implies all required fields of PropertyCreateRequest
	var req CreatePropertyRequest // Use CreatePropertyRequest as PUT implies all fields are present
	if h.bindAndValidate(c, &req) {
		return
	}

	// Convert CreatePropertyRequest to UpdatePropertyRequest for the service
	// as the service layer might expect partial updates via UpdatePropertyRequest.
	// However, for PUT, all fields are typically required.
	// If the service's UpdateProperty strictly expects UpdatePropertyRequest (with pointers for optionality),
	// this conversion is necessary.
	updateReq := UpdatePropertyRequest{
		Name:    &req.Name,   // Assuming Name is required in CreatePropertyRequest
		Address: req.Address, // Assuming Address is *string in CreatePropertyRequest or can be nil
		// Map other fields from CreatePropertyRequest to *UpdatePropertyRequest
	}

	propertyResp, err := h.service.UpdateProperty(c.Request.Context(), propertyID, updateReq)
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, propertyResp)
}

// patchProperty handles PATCH /properties/:propertyId
func (h *PropertyHandler) patchProperty(c *gin.Context) {
	propertyID := c.Param("propertyId")
	var req UpdatePropertyRequest // For PATCH, use UpdatePropertyRequest with pointers for partial updates
	if h.bindAndValidate(c, &req) {
		return
	}

	propertyResp, err := h.service.UpdateProperty(c.Request.Context(), propertyID, req)
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, propertyResp)
}

// deleteProperty handles DELETE /properties/:propertyId
func (h *PropertyHandler) deleteProperty(c *gin.Context) {
	propertyID := c.Param("propertyId")
	err := h.service.DeleteProperty(c.Request.Context(), propertyID)
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// handleError centralizes error handling for the handler.
func (h *PropertyHandler) handleError(c *gin.Context, err error) {
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
	c.JSON(http.StatusInternalServerError, apierrors.ErrInternal(err))
}

// formatValidationErrors converts validator.ValidationErrors to a map.
func (h *PropertyHandler) formatValidationErrors(err error) map[string]string {
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
			default:
				errors[fieldName] = fieldName + " is not valid (" + fieldErr.Tag() + ")."
			}
		}
	}
	return errors
}
