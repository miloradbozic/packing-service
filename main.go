package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/miloradbozic/packing-service/internal/config"
	"github.com/miloradbozic/packing-service/internal/database"
	"github.com/miloradbozic/packing-service/internal/handlers"
	"github.com/miloradbozic/packing-service/internal/middleware"
	"github.com/miloradbozic/packing-service/internal/service"
)

func main() {
	// Load configuration
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config.yaml"
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database
	db, err := database.NewConnection(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run migrations
	migrator := database.NewMigrator(db)
	if err := migrator.RunMigrations("migrations"); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize repository
	packSizeRepo := database.NewPackSizeRepository(db)

	// Pack sizes are now managed via the API - no need for config migration

	// Initialize service
	packingService := service.NewPackingService(packSizeRepo)

	// Initialize handlers
	apiHandler := handlers.NewAPIHandler(packingService, packSizeRepo)
	webHandler, err := handlers.NewWebHandler(packingService, packSizeRepo)
	if err != nil {
		log.Fatalf("Failed to initialize web handler: %v", err)
	}

	// Setup routes
	router := mux.NewRouter()

	// Web UI routes
	router.HandleFunc("/", webHandler.HomePage).Methods("GET", "POST")

	// API routes
	api := router.PathPrefix("/api/v1").Subrouter()
	api.Use(middleware.CORS)
	api.HandleFunc("/calculate", apiHandler.Calculate).Methods("POST", "OPTIONS")
	api.HandleFunc("/config", apiHandler.GetConfig).Methods("GET", "OPTIONS")
	
	// Pack size management routes
	api.HandleFunc("/pack-sizes", apiHandler.ListPackSizes).Methods("GET", "OPTIONS")
	api.HandleFunc("/pack-sizes", apiHandler.CreatePackSize).Methods("POST", "OPTIONS")
	api.HandleFunc("/pack-sizes/{id}", apiHandler.GetPackSize).Methods("GET", "OPTIONS")
	api.HandleFunc("/pack-sizes/{id}", apiHandler.UpdatePackSize).Methods("PUT", "OPTIONS")
	api.HandleFunc("/pack-sizes/{id}", apiHandler.DeletePackSize).Methods("DELETE", "OPTIONS")

	// Health check
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	// Start server
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Starting server on %s", addr)
	log.Printf("Web UI available at http://localhost:%d", cfg.Server.Port)
	log.Printf("API available at http://localhost:%d/api/v1", cfg.Server.Port)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
