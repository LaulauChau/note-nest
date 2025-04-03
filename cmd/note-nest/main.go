package main

import (
	"context"
	"log"

	"github.com/LaulauChau/note-nest/internal/adapter/http"
	"github.com/LaulauChau/note-nest/internal/adapter/http/controller"
	"github.com/LaulauChau/note-nest/internal/adapter/http/router"
	"github.com/LaulauChau/note-nest/internal/application/use_cases"
	"github.com/LaulauChau/note-nest/internal/config"
	"github.com/LaulauChau/note-nest/internal/infrastructure/persistence/repositories"
	"github.com/LaulauChau/note-nest/internal/infrastructure/services"

	"github.com/jackc/pgx/v5"
)

func main() {
	// Load environment variables
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Initialize database
	db, err := pgx.Connect(context.Background(), config.DATABASE.URL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer func() {
		if err := db.Close(context.Background()); err != nil {
			log.Printf("error closing database connection: %v", err)
		}
	}()

	// Initialize the SQLC queries struct
	queries := repositories.New(db)

	// Initialize repository implementations
	userRepo := repositories.NewUserRepository(queries)
	sessionRepo := repositories.NewSessionRepository(queries)
	noteRepo := repositories.NewNoteRepository(queries)

	// Initialize services
	tokenService := services.NewTokenService()
	hashService := services.NewArgonHashService()

	// Initialize use cases
	userUseCase := use_cases.NewUserUseCase(userRepo, hashService)
	sessionUseCase := use_cases.NewSessionUseCase(sessionRepo, userRepo, tokenService)
	noteUseCase := use_cases.NewNoteUseCase(noteRepo, userRepo)

	// Initialize controllers
	userController := controller.NewUserController(userUseCase, sessionUseCase)
	sessionController := controller.NewSessionController(sessionUseCase)
	noteController := controller.NewNoteController(noteUseCase)

	// Initialize router
	r := router.NewRouter(userController, sessionController, noteController)

	// Initialize and start server
	server := http.NewServer(r, config.Server.Port)
	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
