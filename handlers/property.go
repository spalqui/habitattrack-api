package handlers

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/spalqui/habitattrack-api/services"
	"github.com/spalqui/habitattrack-api/types"
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
	h.logger.Info("Property GetByID called")

	id := c.Param("id")

	property, err := h.propertyService.GetPropertyByID(id)
	if err != nil {
		h.logger.Error(fmt.Sprintf("failed to get property by id (id: %s): %v", id, err))
		c.JSON(http.StatusInternalServerError, types.InternalServerError{
			Message: "Failed to retrieve property",
		})
		return
	}

	c.JSON(200, property)
}
