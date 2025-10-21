package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/miloradbozic/packing-service/internal/config"
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

	// Initialize service
	packingService := service.NewPackingService(cfg.Packs.Sizes)

	// Initialize handlers
	apiHandler := handlers.NewAPIHandler(packingService)
	webHandler, err := handlers.NewWebHandler(packingService)
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
