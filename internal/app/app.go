package app

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

// App represents the application with all its dependencies
type App struct {
	config *config.Config
	db     *database.DB
	router *mux.Router
}

// New creates a new application instance
func New() (*App, error) {
	app := &App{}
	
	if err := app.loadConfig(); err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	
	if err := app.setupDatabase(); err != nil {
		return nil, fmt.Errorf("failed to setup database: %w", err)
	}
	
	if err := app.setupRoutes(); err != nil {
		return nil, fmt.Errorf("failed to setup routes: %w", err)
	}
	
	return app, nil
}

// loadConfig loads the application configuration
func (a *App) loadConfig() error {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config.yaml"
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		return err
	}
	
	a.config = cfg
	return nil
}

// setupDatabase initializes the database connection and runs migrations
func (a *App) setupDatabase() error {
	// Connect to database
	db, err := database.NewConnection(&a.config.Database)
	if err != nil {
		return err
	}
	a.db = db

	// Run migrations
	migrator := database.NewMigrator(db)
	if err := migrator.RunMigrations("migrations"); err != nil {
		return err
	}

	return nil
}

// setupRoutes configures all the HTTP routes
func (a *App) setupRoutes() error {
	// Initialize repository and services
	packSizeRepo := database.NewPackSizeRepository(a.db)
	packingService := service.NewPackingService(packSizeRepo)

	// Initialize handlers
	apiHandler := handlers.NewAPIHandler(packingService, packSizeRepo)
	webHandler, err := handlers.NewWebHandler(packingService, packSizeRepo)
	if err != nil {
		return err
	}

	// Setup router
	router := mux.NewRouter()
	
	// Web UI routes
	router.HandleFunc("/", webHandler.HomePage).Methods("GET", "POST")

	// API routes
	api := router.PathPrefix("/api/v1").Subrouter()
	api.Use(middleware.CORS)
	
	// Calculation routes
	api.HandleFunc("/calculate", apiHandler.Calculate).Methods("POST", "OPTIONS")
	api.HandleFunc("/config", apiHandler.GetConfig).Methods("GET", "OPTIONS")
	
	// Pack size management routes
	api.HandleFunc("/pack-sizes", apiHandler.ListPackSizes).Methods("GET", "OPTIONS")
	api.HandleFunc("/pack-sizes", apiHandler.CreatePackSize).Methods("POST", "OPTIONS")
	api.HandleFunc("/pack-sizes/{id}", apiHandler.GetPackSize).Methods("GET", "OPTIONS")
	api.HandleFunc("/pack-sizes/{id}", apiHandler.UpdatePackSize).Methods("PUT", "OPTIONS")
	api.HandleFunc("/pack-sizes/{id}", apiHandler.DeletePackSize).Methods("DELETE", "OPTIONS")

	// Health check
	router.HandleFunc("/health", a.healthCheck).Methods("GET")

	a.router = router
	return nil
}

// healthCheck handles the health check endpoint
func (a *App) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// Run starts the HTTP server
func (a *App) Run() error {
	addr := fmt.Sprintf("%s:%d", a.config.Server.Host, a.config.Server.Port)
	
	log.Printf("Starting server on %s", addr)
	log.Printf("Web UI available at http://localhost:%d", a.config.Server.Port)
	log.Printf("API available at http://localhost:%d/api/v1", a.config.Server.Port)

	return http.ListenAndServe(addr, a.router)
}

// Close cleans up resources
func (a *App) Close() error {
	if a.db != nil {
		return a.db.Close()
	}
	return nil
}
