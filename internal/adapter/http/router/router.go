package router

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	httpMiddleware "github.com/LaulauChau/note-nest/internal/adapter/http/middleware"
	"github.com/LaulauChau/note-nest/internal/infrastructure/services"
)

func NewRouter() http.Handler {
	r := chi.NewRouter()

	// Initialize services
	csrfService := services.NewCSRFService(30 * time.Minute)

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(httpMiddleware.SecurityHeaders)

	// CSRF token endpoint
	r.Get("/api/csrf-token", func(w http.ResponseWriter, r *http.Request) {
		token := csrfService.GenerateToken(r)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		err := json.NewEncoder(w).Encode(map[string]string{
			"csrf_token": token,
		})
		if err != nil {
			log.Printf("Error encoding response: %v", err)
		}
	})

	return r
}
