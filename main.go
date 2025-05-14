package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/spalqui/habitattrack-api/config"
	"github.com/spalqui/habitattrack-api/handlers"
	"github.com/spalqui/habitattrack-api/repositories"
	"github.com/spalqui/habitattrack-api/services"
)

func main() {
	port := os.Getenv("PORT")
	googleCloudProject := os.Getenv("GOOGLE_CLOUD_PROJECT")

	c, err := config.New(
		config.WithPort(port),
		config.WithGoogleCloudProject(googleCloudProject),
	)
	if err != nil {
		log.Fatalf("error creating config: %v", err)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	r := gin.Default()

	propertyRepo, err := repositories.NewFirestorePropertyRepository(context.Background(), c.GoogleCloudProject)
	if err != nil {
		log.Fatalf("error creating firestore property repository: %v", err)
	}

	propertyService := services.NewPropertyService(propertyRepo)

	healthHandler := handlers.NewHealthHandler(logger)
	propertyHandler := handlers.NewPropertyHandler(logger, propertyService)

	r.GET("/health", healthHandler.Check)

	r.GET("/property/:id", propertyHandler.GetByID)

	err = r.Run(fmt.Sprintf(":%d", c.Port))
	if err != nil {
		log.Fatalf("failed to run: %v", err)
	}
}
