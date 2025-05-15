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
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	databaseID := os.Getenv("DATABASE_ID")

	c, err := config.New(
		config.WithPort(port),
		config.WithGoogleCloudProject(projectID),
		config.WithFirestoreDatabase(databaseID),
	)
	if err != nil {
		log.Fatalf("error creating config: %v", err)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	r := gin.Default()

	propertyRepo, err := repositories.NewFirestorePropertyRepository(context.Background(), c.ProjectID, c.DatabaseID)
	if err != nil {
		log.Fatalf("error creating firestore property repository: %v", err)
	}

	propertyService := services.NewPropertyService(propertyRepo)

	healthHandler := handlers.NewHealthHandler(logger)
	propertyHandler := handlers.NewPropertyHandler(logger, propertyService)

	propertyRouter := r.Group("/property")
	{
		propertyRouter.GET("/:id", propertyHandler.GetByID)
		propertyRouter.GET("/", propertyHandler.List)
		propertyRouter.POST("/", propertyHandler.Create)
		propertyRouter.PATCH("/:id", propertyHandler.Update)
		propertyRouter.DELETE("/:id", propertyHandler.Delete)
	}

	r.GET("/health", healthHandler.Check)

	err = r.Run(fmt.Sprintf(":%d", c.Port))
	if err != nil {
		log.Fatalf("failed to run: %v", err)
	}
}
