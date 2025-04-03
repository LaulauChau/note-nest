package main

import (
	"log"

	"github.com/LaulauChau/note-nest/internal/adapter/http"
	"github.com/LaulauChau/note-nest/internal/adapter/http/router"
	"github.com/LaulauChau/note-nest/internal/config"
)

func main() {
	// Load environment variables
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Initialize router
	r := router.NewRouter()

	// Initialize and start server
	server := http.NewServer(r, config.Server.Port)
	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
