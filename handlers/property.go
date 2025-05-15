package handlers

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/spalqui/habitattrack-api/constants"
	"github.com/spalqui/habitattrack-api/repositories"
	"github.com/spalqui/habitattrack-api/services"
)

type PropertyHandler struct {
	logger          *slog.Logger
	propertyService *services.PropertyService
}

func NewPropertyHandler(
	logger *slog.Logger,
	propertyService *services.PropertyService,
) *PropertyHandler {
	return &PropertyHandler{
		logger:          logger,
		propertyService: propertyService,
	}
}

func (h *PropertyHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		h.logger.Error("property ID is missing")
		returnError(c, http.StatusBadRequest, "property ID is missing")
		return
	}

	property, err := h.propertyService.GetPropertyByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repositories.ErrPropertyNotFound) {
			h.logger.Warn("property not found", slog.String("id", id))
			returnError(c, http.StatusNotFound, "property not found")
			return
		}
		h.logger.Error("failed to get property", slog.String("id", id), slog.Any("error", err))
		returnError(c, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	returnOK(c, property, nil)
}

func (h *PropertyHandler) List(c *gin.Context) {
	limit := constants.MaxPageSize
	cursor := c.Query("cursor")

	if limitStr := c.Query("limit"); limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err != nil {
			h.logger.Error("invalid limit parameter", slog.Any("error", err))
			returnError(c, http.StatusBadRequest, "invalid limit parameter")
			return
		}
		if parsedLimit < 1 || parsedLimit > constants.MaxPageSize {
			h.logger.Error("limit out of range", slog.Int("limit", parsedLimit))
			returnError(c, http.StatusBadRequest, "limit must be between 1 and 100")
			return
		}
		limit = parsedLimit
	}

	properties, nextCursor, err := h.propertyService.GetProperties(c.Request.Context(), limit, cursor)
	if err != nil {
		h.logger.Error("failed to get properties", slog.Any("error", err))
		returnError(c, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	returnOK(c, properties, &Metadata{
		NextCursor: nextCursor,
		Limit:      limit,
		Total:      len(properties),
	})
}

func (h *PropertyHandler) Create(c *gin.Context) {
	var req CreatePropertyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("invalid request body", slog.Any("error", err))
		returnError(c, http.StatusBadRequest, "invalid request body")
		return
	}

	property := req.ToProperty()
	err := h.propertyService.CreateProperty(c.Request.Context(), property)
	if err != nil {
		h.logger.Error("failed to create property", slog.Any("error", err))
		returnError(c, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	c.Header("Location", fmt.Sprintf("/property/%s", property.ID))
	returnOK(c, property, nil)
}

func (h *PropertyHandler) Update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		h.logger.Error("property ID is missing")
		returnError(c, http.StatusBadRequest, "property ID is missing")
		return
	}

	var req UpdatePropertyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("invalid request body", slog.Any("error", err))
		returnError(c, http.StatusBadRequest, "invalid request body")
		return
	}

	property, err := h.propertyService.GetPropertyByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repositories.ErrPropertyNotFound) {
			h.logger.Warn("property not found", slog.String("id", id))
			returnError(c, http.StatusNotFound, "property not found")
			return
		}
		h.logger.Error("failed to get property", slog.String("id", id), slog.Any("error", err))
		returnError(c, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	req.ApplyUpdates(&property)
	if err := h.propertyService.UpdateProperty(c.Request.Context(), id, &property); err != nil {
		h.logger.Error("failed to update property", slog.String("id", id), slog.Any("error", err))
		returnError(c, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	returnOK(c, property, nil)
}

func (h *PropertyHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		h.logger.Error("property ID is missing")
		returnError(c, http.StatusBadRequest, "property ID is missing")
		return
	}

	if err := h.propertyService.DeleteProperty(c.Request.Context(), id); err != nil {
		if !errors.Is(err, repositories.ErrPropertyNotFound) {
			h.logger.Error("failed to delete property", slog.String("id", id), slog.Any("error", err))
			returnError(c, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
	}

	returnOK(c, nil, nil)
}
