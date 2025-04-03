package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/LaulauChau/note-nest/internal/adapter/http/controller"
	httpMiddleware "github.com/LaulauChau/note-nest/internal/adapter/http/middleware"
)

func NewRouter(userController *controller.UserController, sessionController *controller.SessionController) http.Handler {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(httpMiddleware.SecurityHeaders)

	// Public routes
	r.Group(func(r chi.Router) {
		r.Post("/api/register", userController.Register)
		r.Post("/api/login", userController.Login)
	})

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(sessionController.AuthMiddleware)

		r.Post("/api/logout", sessionController.Logout)
		r.Get("/api/me", userController.GetCurrentUser)
	})

	return r
}
