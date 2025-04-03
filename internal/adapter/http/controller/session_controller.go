package controller

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/LaulauChau/note-nest/internal/application/use_cases"
)

// ContextKey is used to identify values in the context
type ContextKey string

const (
	// UserContextKey is the key used to store the user in the context
	UserContextKey ContextKey = "user"
)

type SessionController struct {
	sessionUseCase *use_cases.SessionUseCase
}

func NewSessionController(sessionUseCase *use_cases.SessionUseCase) *SessionController {
	return &SessionController{
		sessionUseCase: sessionUseCase,
	}
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	UserID string `json:"user_id"`
}

func (c *SessionController) Login(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// We need to inject a UserUseCase dependency to authenticate the user
	// For now, let's implement this in a separate controller
	http.Error(w, "Not implemented", http.StatusInternalServerError)
}

func (c *SessionController) Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get the session token from the cookie
	cookie, err := r.Cookie("session")
	if err != nil {
		http.Error(w, "No session found", http.StatusBadRequest)
		return
	}

	// Validate the session token
	token := cookie.Value
	result, err := c.sessionUseCase.ValidateSessionToken(ctx, token)
	if err != nil {
		http.Error(w, "Failed to validate session", http.StatusInternalServerError)
		return
	}

	// If session is valid, invalidate it
	if result.Session != nil {
		if err := c.sessionUseCase.InvalidateSession(ctx, result.Session.ID); err != nil {
			http.Error(w, "Failed to invalidate session", http.StatusInternalServerError)
			return
		}
	}

	// Clear the session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		MaxAge:   -1,
	})

	w.WriteHeader(http.StatusNoContent)
}

func (c *SessionController) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Get the session token from the cookie
		cookie, err := r.Cookie("session")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Validate the session token
		token := cookie.Value
		result, err := c.sessionUseCase.ValidateSessionToken(ctx, token)
		if err != nil {
			http.Error(w, "Failed to validate session", http.StatusInternalServerError)
			return
		}

		// If session is invalid, return unauthorized
		if result.Session == nil || result.User == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Add the user to the context
		ctx = context.WithValue(ctx, UserContextKey, result.User)

		// Call the next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
