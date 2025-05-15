package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

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

	property, err := h.propertyService.GetPropertyByID(id)
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
