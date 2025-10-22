package main

import (
	"log"

	"github.com/miloradbozic/packing-service/internal/app"
)

func main() {
	// Create and initialize the application
	application, err := app.New()
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}
	defer application.Close()

	// Start the server
	if err := application.Run(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
