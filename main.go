package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/spalqui/habitattrack-api/handlers"
	"github.com/spalqui/habitattrack-api/repositories"
	"github.com/spalqui/habitattrack-api/services"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	r := gin.Default()

	propertyRepo := repositories.NewFirestorePropertyRepository()

	propertyService := services.NewPropertyService(propertyRepo)

	healthHandler := handlers.NewHealthHandler(logger)
	propertyHandler := handlers.NewPropertyHandler(logger, propertyService)

	r.GET("/health", healthHandler.Check)

	r.GET("/property/:id", propertyHandler.GetByID)

	err := r.Run()
	if err != nil {
		log.Fatalf("failed to run: %v", err)
	}
}
