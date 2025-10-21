package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/miloradbozic/packing-service/internal/config"
	"github.com/miloradbozic/packing-service/internal/database"
)

func main() {
	var configPath = flag.String("config", "config.yaml", "Path to configuration file")
	var migrationsPath = flag.String("migrations", "migrations", "Path to migrations directory")
	flag.Parse()

	// Load configuration
	cfg, err := config.Load(*configPath)
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
	if err := migrator.RunMigrations(*migrationsPath); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	fmt.Println("Migrations completed successfully!")
}
