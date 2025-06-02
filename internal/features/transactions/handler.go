package transactions

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/spalqui/habitattrack-api/internal/core/entity"
	"github.com/spalqui/habitattrack-api/internal/shared/apierrors"
)

const (
	defaultPage     = 1
	defaultPageSize = 20
	queryPage       = "page"
	queryLimit      = "limit"
	queryPropertyID = "propertyId"
	queryType       = "type"
	queryCategoryID = "categoryId"
	queryStartDate  = "startDate"
	queryEndDate    = "endDate"
	rfc3339Format   = time.RFC3339 // OpenAPI date-time format
)

// TransactionHandler handles HTTP requests for transactions.
type TransactionHandler struct {
	service  TransactionService
	validate *validator.Validate
}

// NewTransactionHandler creates a new instance of TransactionHandler.
func NewTransactionHandler(service TransactionService) *TransactionHandler {
	return &TransactionHandler{
		service:  service,
		validate: validator.New(),
	}
}

// RegisterRoutes sets up the routes for transaction operations.
func (h *TransactionHandler) RegisterRoutes(router *gin.Engine) {
	transactionsGroup := router.Group("/transactions")
	{
		transactionsGroup.POST("", h.createTransaction)
		transactionsGroup.GET("", h.listTransactions)
		transactionsGroup.GET("/:transactionId", h.getTransactionByID)
		transactionsGroup.PUT("/:transactionId", h.updateTransaction)
		transactionsGroup.PATCH("/:transactionId", h.patchTransaction)
		transactionsGroup.DELETE("/:transactionId", h.deleteTransaction)
	}
}

// bindAndValidate binds the request JSON to the given struct and validates it.
// It returns true if binding or validation fails and the response has been written.
func (h *TransactionHandler) bindAndValidate(c *gin.Context, req interface{}) bool {
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

// createTransaction handles POST /transactions
func (h *TransactionHandler) createTransaction(c *gin.Context) {
	var req CreateTransactionRequest
	if h.bindAndValidate(c, &req) {
		return
	}

	transactionResp, err := h.service.CreateTransaction(c.Request.Context(), req)
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, transactionResp)
}

// listTransactions handles GET /transactions
func (h *TransactionHandler) listTransactions(c *gin.Context) {
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

	var propertyIDFilter *string
	if pID, ok := c.GetQuery(queryPropertyID); ok {
		propertyIDFilter = &pID
	}

	var typeFilter *entity.TransactionType
	if t, ok := c.GetQuery(queryType); ok {
		tt := entity.TransactionType(t)
		if tt == entity.IncomeTransaction || tt == entity.ExpenseTransaction {
			typeFilter = &tt
		} else {
			appErr := apierrors.ErrBadRequest("Invalid 'type' query parameter. Must be 'income' or 'expense'.")
			c.JSON(http.StatusBadRequest, appErr)
			return
		}
	}

	var categoryIDFilter *string
	if catID, ok := c.GetQuery(queryCategoryID); ok {
		categoryIDFilter = &catID
	}

	var startDateFilter *time.Time
	if sdStr, ok := c.GetQuery(queryStartDate); ok {
		sd, err := time.Parse(rfc3339Format, sdStr)
		if err != nil {
			appErr := apierrors.ErrBadRequest("Invalid 'startDate' format. Use RFC3339 (YYYY-MM-DDTHH:mm:ssZ).")
			c.JSON(http.StatusBadRequest, appErr)
			return
		}
		startDateFilter = &sd
	}

	var endDateFilter *time.Time
	if edStr, ok := c.GetQuery(queryEndDate); ok {
		ed, err := time.Parse(rfc3339Format, edStr)
		if err != nil {
			appErr := apierrors.ErrBadRequest("Invalid 'endDate' format. Use RFC3339 (YYYY-MM-DDTHH:mm:ssZ).")
			c.JSON(http.StatusBadRequest, appErr)
			return
		}
		endDateFilter = &ed
	}

	transactionsResp, err := h.service.ListTransactions(
		c.Request.Context(),
		propertyIDFilter,
		typeFilter,
		categoryIDFilter,
		startDateFilter,
		endDateFilter,
		page,
		pageSize,
	)
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, transactionsResp)
}

// getTransactionByID handles GET /transactions/:transactionId
func (h *TransactionHandler) getTransactionByID(c *gin.Context) {
	transactionID := c.Param("transactionId")
	transactionResp, err := h.service.GetTransactionByID(c.Request.Context(), transactionID)
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, transactionResp)
}

// updateTransaction handles PUT /transactions/:transactionId
func (h *TransactionHandler) updateTransaction(c *gin.Context) {
	transactionID := c.Param("transactionId")
	var req CreateTransactionRequest // For PUT, all fields are typically required
	if h.bindAndValidate(c, &req) {
		return
	}

	// Convert CreateTransactionRequest to UpdateTransactionRequest for the service.
	// This assumes the service's UpdateTransaction method expects UpdateTransactionRequest
	// for consistency, even if PUT implies all fields.
	updateReq := UpdateTransactionRequest{
		Amount:          &req.Amount,
		TransactionDate: &req.TransactionDate,
		Description:     req.Description, // Assumes Description is *string in CreateTransactionRequest or can be nil
		Type:            &req.Type,
		CategoryID:      &req.CategoryID,
		PropertyID:      req.PropertyID, // Assumes PropertyID is *string in CreateTransactionRequest or can be nil
	}

	transactionResp, err := h.service.UpdateTransaction(c.Request.Context(), transactionID, updateReq)
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, transactionResp)
}

// patchTransaction handles PATCH /transactions/:transactionId
func (h *TransactionHandler) patchTransaction(c *gin.Context) {
	transactionID := c.Param("transactionId")
	var req UpdateTransactionRequest // For PATCH, use UpdateTransactionRequest with pointers for partial updates
	if h.bindAndValidate(c, &req) {
		return
	}

	transactionResp, err := h.service.UpdateTransaction(c.Request.Context(), transactionID, req)
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, transactionResp)
}

// deleteTransaction handles DELETE /transactions/:transactionId
func (h *TransactionHandler) deleteTransaction(c *gin.Context) {
	transactionID := c.Param("transactionId")
	err := h.service.DeleteTransaction(c.Request.Context(), transactionID)
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// handleError centralizes error handling for the handler.
func (h *TransactionHandler) handleError(c *gin.Context, err error) {
	if appErr, ok := err.(*apierrors.Error); ok {
		switch appErr.Code {
		case apierrors.CodeNotFound:
			c.JSON(http.StatusNotFound, appErr)
		case apierrors.CodeConflict:
			c.JSON(http.StatusConflict, appErr)
		case apierrors.CodeValidationError, apierrors.CodeBadRequest, apierrors.CodeUnprocessableEntity:
			// For CodeUnprocessableEntity, http.StatusUnprocessableEntity (422) might be more appropriate
			// but for simplicity, grouping with Bad Request for now.
			if appErr.Code == apierrors.CodeUnprocessableEntity {
				c.JSON(http.StatusUnprocessableEntity, appErr)
			} else {
				c.JSON(http.StatusBadRequest, appErr)
			}
		case apierrors.CodeUnauthorized:
			c.JSON(http.StatusUnauthorized, appErr)
		case apierrors.CodeForbidden:
			c.JSON(http.StatusForbidden, appErr)
		default: // Includes CodeInternalError
			c.JSON(http.StatusInternalServerError, appErr)
		}
		return
	}
	// Fallback for unexpected errors not of type *apierrors.Error
	c.JSON(http.StatusInternalServerError, apierrors.ErrInternal(err))
}

// formatValidationErrors converts validator.ValidationErrors to a map.
func (h *TransactionHandler) formatValidationErrors(err error) map[string]string {
	errors := make(map[string]string)
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldErr := range validationErrors {
			// Provide a more user-friendly message based on the tag
			fieldName := fieldErr.Field()
			switch fieldErr.Tag() {
			case "required":
				errors[fieldName] = fieldName + " is required."
			case "uuid":
				errors[fieldName] = fieldName + " must be a valid UUID."
			case "oneof":
				errors[fieldName] = fieldName + " must be one of [" + fieldErr.Param() + "]."
			case "gt":
				errors[fieldName] = fieldName + " must be greater than " + fieldErr.Param() + "."
			case "max":
				errors[fieldName] = fieldName + " must not exceed " + fieldErr.Param() + " characters."
			default:
				errors[fieldName] = fieldName + " is not valid (" + fieldErr.Tag() + ")."
			}
		}
	}
	return errors
}
