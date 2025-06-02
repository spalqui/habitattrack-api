package main

import (
	"context"
	"log"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/gorilla/mux"
	"google.golang.org/api/option"

	"github.com/spalqui/habitattrack-api/internal/config"
	"github.com/spalqui/habitattrack-api/internal/handlers"
	"github.com/spalqui/habitattrack-api/internal/services"
	firestoreRepo "github.com/spalqui/habitattrack-api/pkg/firestore"
	"github.com/spalqui/habitattrack-api/pkg/middleware"
)

func main() {
	cfg := config.Load()

	// Initialize Firestore client
	ctx := context.Background()
	var client *firestore.Client
	var err error

	if cfg.FirestoreKeyPath != "" {
		client, err = firestore.NewClientWithDatabase(ctx, cfg.GoogleProject, "habitattrack", option.WithCredentialsFile(cfg.FirestoreKeyPath))
	} else {
		client, err = firestore.NewClientWithDatabase(ctx, cfg.GoogleProject, "habitattrack")
	}
	if err != nil {
		log.Fatalf("Failed to create Firestore client: %v", err)
	}
	defer client.Close()

	// Initialize repositories
	propertyRepo := firestoreRepo.NewPropertyRepository(client)
	transactionRepo := firestoreRepo.NewTransactionRepository(client)
	categoryRepo := firestoreRepo.NewCategoryRepository(client)

	// Initialize services
	propertyService := services.NewPropertyService(propertyRepo)
	transactionService := services.NewTransactionService(transactionRepo, categoryRepo, propertyRepo)
	categoryService := services.NewCategoryService(categoryRepo)

	// Initialize handlers
	propertyHandler := handlers.NewPropertyHandler(propertyService)
	transactionHandler := handlers.NewTransactionHandler(transactionService)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	// Setup routes
	router := setupRoutes(propertyHandler, transactionHandler, categoryHandler)

	log.Printf("Server starting on port %s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, router))
}

func setupRoutes(propertyHandler *handlers.PropertyHandler, transactionHandler *handlers.TransactionHandler, categoryHandler *handlers.CategoryHandler) *mux.Router {
	router := mux.NewRouter()

	// Add middleware
	router.Use(middleware.CORS)
	router.Use(middleware.JSONContentType)
	router.Use(middleware.Logging)

	// Property routes
	router.HandleFunc("/properties", propertyHandler.CreateProperty).Methods("POST")
	router.HandleFunc("/properties", propertyHandler.GetAllProperties).Methods("GET")
	router.HandleFunc("/properties/{id}", propertyHandler.GetProperty).Methods("GET")
	router.HandleFunc("/properties/{id}", propertyHandler.UpdateProperty).Methods("PUT")
	router.HandleFunc("/properties/{id}", propertyHandler.DeleteProperty).Methods("DELETE")

	// Transaction routes
	router.HandleFunc("/transactions", transactionHandler.CreateTransaction).Methods("POST")
	router.HandleFunc("/transactions", transactionHandler.GetAllTransactions).Methods("GET")
	router.HandleFunc("/transactions/{id}", transactionHandler.GetTransaction).Methods("GET")
	router.HandleFunc("/transactions/{id}", transactionHandler.UpdateTransaction).Methods("PUT")
	router.HandleFunc("/transactions/{id}", transactionHandler.DeleteTransaction).Methods("DELETE")
	router.HandleFunc("/properties/{propertyId}/transactions", transactionHandler.GetTransactionsByProperty).Methods("GET")

	// Category routes
	router.HandleFunc("/categories", categoryHandler.CreateCategory).Methods("POST")
	router.HandleFunc("/categories", categoryHandler.GetAllCategories).Methods("GET")
	router.HandleFunc("/categories/{id}", categoryHandler.GetCategory).Methods("GET")
	router.HandleFunc("/categories/{id}", categoryHandler.UpdateCategory).Methods("PUT")
	router.HandleFunc("/categories/{id}", categoryHandler.DeleteCategory).Methods("DELETE")
	router.HandleFunc("/categories/type/{type}", categoryHandler.GetCategoriesByType).Methods("GET")

	// Health check
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	return router
}
