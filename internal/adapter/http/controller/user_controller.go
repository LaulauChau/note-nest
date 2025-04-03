package controller

import (
	"encoding/json"
	"net/http"

	"github.com/LaulauChau/note-nest/internal/application/use_cases"
	"github.com/LaulauChau/note-nest/internal/domain/entities"
)

type UserController struct {
	userUseCase    *use_cases.UserUseCase
	sessionUseCase *use_cases.SessionUseCase
}

func NewUserController(
	userUseCase *use_cases.UserUseCase,
	sessionUseCase *use_cases.SessionUseCase,
) *UserController {
	return &UserController{
		userUseCase:    userUseCase,
		sessionUseCase: sessionUseCase,
	}
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

func (c *UserController) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse the request body
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input
	if req.Email == "" || req.Name == "" || req.Password == "" {
		http.Error(w, "Email, name, and password are required", http.StatusBadRequest)
		return
	}

	// Register the user
	user, err := c.userUseCase.RegisterUser(ctx, req.Email, req.Name, req.Password)
	if err != nil {
		if err.Error() == "email already taken" {
			http.Error(w, "Email already taken", http.StatusConflict)
			return
		}
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	// Return the user
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(RegisterResponse{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
	}); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (c *UserController) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse the request body
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Authenticate the user
	user, err := c.userUseCase.AuthenticateUser(ctx, req.Email, req.Password)
	if err != nil {
		http.Error(w, "Authentication failed", http.StatusInternalServerError)
		return
	}

	// If authentication failed, return unauthorized
	if user == nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate a new session token
	token, err := c.sessionUseCase.GenerateSessionToken(ctx)
	if err != nil {
		http.Error(w, "Failed to generate session token", http.StatusInternalServerError)
		return
	}

	// Create a new session
	session, err := c.sessionUseCase.CreateSession(ctx, token, user.ID)
	if err != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	// Set the session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  session.ExpiresAt,
	})

	// Return the user
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(LoginResponse{
		UserID: user.ID,
	}); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

type UserResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

func (c *UserController) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	// Get user from context (added by the auth middleware)
	user, ok := r.Context().Value(UserContextKey).(*entities.User)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Return the user
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(UserResponse{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
	}); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
