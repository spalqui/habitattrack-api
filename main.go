package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/spalqui/habitattrack-api/config"
	"github.com/spalqui/habitattrack-api/internal/features/categories"
	"github.com/spalqui/habitattrack-api/internal/features/properties"
	"github.com/spalqui/habitattrack-api/internal/features/transactions"
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

	propertyRepo, err := properties.NewFirestorePropertyRepository(context.Background(), c.ProjectID, c.DatabaseID)
	if err != nil {
		logger.Error("error creating firestore property repository", slog.Any("error", err))
		log.Fatalf("error creating firestore property repository: %v", err)
	}
	// defer propertyRepo.CloseClient() // Consider client close on shutdown

	transactionRepo, err := transactions.NewFirestoreTransactionRepository(context.Background(), c.ProjectID, c.DatabaseID)
	if err != nil {
		logger.Error("error creating firestore transaction repository", slog.Any("error", err))
		log.Fatalf("error creating firestore transaction repository: %v", err)
	}
	// defer transactionRepo.CloseClient()

	categoryRepo, err := categories.NewFirestoreCategoryRepository(context.Background(), c.ProjectID, c.DatabaseID)
	if err != nil {
		logger.Error("error creating firestore category repository", slog.Any("error", err))
		log.Fatalf("error creating firestore category repository: %v", err)
	}
	// defer categoryRepo.CloseClient()

	// Feature-sliced services and handlers
	propertyService := properties.NewPropertyService(propertyRepo)
	propertyHandler := properties.NewPropertyHandler(propertyService)

	transactionsService := transactions.NewTransactionService(transactionRepo, categoryRepo)
	transactionsHandler := transactions.NewTransactionHandler(transactionsService)

	categoriesService := categories.NewCategoryService(categoryRepo)
	categoriesHandler := categories.NewCategoryHandler(categoriesService)

	// Register feature routes
	propertyHandler.RegisterRoutes(r)
	transactionsHandler.RegisterRoutes(r)
	categoriesHandler.RegisterRoutes(r)

	logger.Info("Server starting", slog.Int("port", c.Port))
	err = r.Run(fmt.Sprintf(":%d", c.Port))
	if err != nil {
		log.Fatalf("failed to run: %v", err)
	}
}
