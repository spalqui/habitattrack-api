package handlers

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	logger *slog.Logger
}

func NewHealthHandler(logger *slog.Logger) *HealthHandler {
	return &HealthHandler{
		logger: logger,
	}
}

type HealthCheckResponse struct {
	Status string `json:"status"`
}

func (h *HealthHandler) Check(c *gin.Context) {
	h.logger.Info("Health Check called")

	c.JSON(http.StatusOK, HealthCheckResponse{
		Status: http.StatusText(http.StatusOK),
	})
}
