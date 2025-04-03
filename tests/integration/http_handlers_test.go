package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/LaulauChau/note-nest/internal/adapter/http/controller"
	"github.com/LaulauChau/note-nest/internal/adapter/http/router"
	"github.com/LaulauChau/note-nest/internal/application/use_cases"
	"github.com/LaulauChau/note-nest/internal/infrastructure/persistence/repositories"
	"github.com/LaulauChau/note-nest/internal/infrastructure/services"
)

func TestHTTPHandlers(t *testing.T) {
	// Set up test database
	ctx := context.Background()
	db, err := SetupTestDatabase(ctx)
	require.NoError(t, err)
	defer db.Close(ctx)

	// Initialize repositories
	queries := repositories.New(db.Pool)
	userRepo := repositories.NewUserRepository(queries)
	sessionRepo := repositories.NewSessionRepository(queries)
	noteRepo := repositories.NewNoteRepository(queries)
	labelRepo := repositories.NewLabelRepository(queries)

	// Initialize services
	tokenService := services.NewTokenService()
	hashService := services.NewArgonHashService()

	// Initialize use cases
	userUseCase := use_cases.NewUserUseCase(userRepo, hashService)
	sessionUseCase := use_cases.NewSessionUseCase(sessionRepo, userRepo, tokenService)
	noteUseCase := use_cases.NewNoteUseCase(noteRepo, userRepo, labelRepo)
	labelUseCase := use_cases.NewLabelUseCase(labelRepo, userRepo, noteRepo)

	// Initialize controllers
	userController := controller.NewUserController(userUseCase, sessionUseCase)
	sessionController := controller.NewSessionController(sessionUseCase)
	noteController := controller.NewNoteController(noteUseCase, labelUseCase)
	labelController := controller.NewLabelController(labelUseCase)

	// Initialize router
	r := router.NewRouter(userController, sessionController, noteController, labelController)

	t.Run("RegisterUser", func(t *testing.T) {
		// Create request payload
		registerPayload := map[string]string{
			"email":    "registertest@example.com",
			"name":     "Register Test",
			"password": "SecureP@ssw0rd",
		}
		jsonPayload, err := json.Marshal(registerPayload)
		require.NoError(t, err)

		// Create request
		req := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewBuffer(jsonPayload))
		req.Header.Set("Content-Type", "application/json")

		// Perform request
		recorder := httptest.NewRecorder()
		r.ServeHTTP(recorder, req)

		// Check response
		require.Equal(t, http.StatusCreated, recorder.Code)

		var response map[string]interface{}
		err = json.Unmarshal(recorder.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.NotNil(t, response["id"])
		assert.Equal(t, registerPayload["email"], response["email"])
		assert.Equal(t, registerPayload["name"], response["name"])
	})

	t.Run("RegisterUserInvalidEmail", func(t *testing.T) {
		// Create request payload with invalid email
		registerPayload := map[string]string{
			"email":    "invalid-email",
			"name":     "Invalid Email Test",
			"password": "SecureP@ssw0rd123",
		}
		jsonPayload, err := json.Marshal(registerPayload)
		require.NoError(t, err)

		// Create request
		req := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewBuffer(jsonPayload))
		req.Header.Set("Content-Type", "application/json")

		// Perform request
		recorder := httptest.NewRecorder()
		r.ServeHTTP(recorder, req)

		// Check response
		require.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("RegisterUserPasswordTooShort", func(t *testing.T) {
		// Create request payload with short password
		registerPayload := map[string]string{
			"email":    "short@example.com",
			"name":     "Short Password Test",
			"password": "Short1@", // Less than 12 characters
		}
		jsonPayload, err := json.Marshal(registerPayload)
		require.NoError(t, err)

		// Create request
		req := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewBuffer(jsonPayload))
		req.Header.Set("Content-Type", "application/json")

		// Perform request
		recorder := httptest.NewRecorder()
		r.ServeHTTP(recorder, req)

		// Check response
		require.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("RegisterUserPasswordNoUppercase", func(t *testing.T) {
		// Create request payload with password missing uppercase
		registerPayload := map[string]string{
			"email":    "noupper@example.com",
			"name":     "No Uppercase Test",
			"password": "nouppercase1@#", // Missing uppercase
		}
		jsonPayload, err := json.Marshal(registerPayload)
		require.NoError(t, err)

		// Create request
		req := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewBuffer(jsonPayload))
		req.Header.Set("Content-Type", "application/json")

		// Perform request
		recorder := httptest.NewRecorder()
		r.ServeHTTP(recorder, req)

		// Check response
		require.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("LoginUserInvalidEmail", func(t *testing.T) {
		// Create login payload with invalid email
		loginPayload := map[string]string{
			"email":    "invalid-email",
			"password": "L0ginP@ssword123",
		}
		jsonPayload, err := json.Marshal(loginPayload)
		require.NoError(t, err)

		// Create login request
		req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewBuffer(jsonPayload))
		req.Header.Set("Content-Type", "application/json")

		// Perform request
		recorder := httptest.NewRecorder()
		r.ServeHTTP(recorder, req)

		// Check response
		require.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("LoginUser", func(t *testing.T) {
		// First register a user
		email := "logintest@example.com"
		name := "Login Test"
		password := "L0ginP@ssword123"

		user, err := userUseCase.RegisterUser(ctx, email, name, password)
		require.NoError(t, err)
		require.NotNil(t, user)

		// Create login payload
		loginPayload := map[string]string{
			"email":    email,
			"password": password,
		}
		jsonPayload, err := json.Marshal(loginPayload)
		require.NoError(t, err)

		// Create login request
		req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewBuffer(jsonPayload))
		req.Header.Set("Content-Type", "application/json")

		// Perform request
		recorder := httptest.NewRecorder()
		r.ServeHTTP(recorder, req)

		// Check response
		require.Equal(t, http.StatusOK, recorder.Code)

		var response map[string]interface{}
		err = json.Unmarshal(recorder.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, user.ID, response["user_id"])

		// Check for session cookie
		cookies := recorder.Result().Cookies()
		var sessionCookie *http.Cookie
		for _, cookie := range cookies {
			if cookie.Name == "session" {
				sessionCookie = cookie
				break
			}
		}
		require.NotNil(t, sessionCookie, "Session cookie should be set")
		assert.True(t, sessionCookie.HttpOnly, "Session cookie should be HttpOnly")
		assert.Equal(t, http.SameSiteStrictMode, sessionCookie.SameSite, "Session cookie should have SameSite=Strict")
	})

	t.Run("ProtectedRoute", func(t *testing.T) {
		// First register and login a user
		email := "protectedtest@example.com"
		name := "Protected Test"
		password := "Pr0tected!P@ssw0rd"

		user, err := userUseCase.RegisterUser(ctx, email, name, password)
		require.NoError(t, err)

		// Generate session token and create session
		token, err := sessionUseCase.GenerateSessionToken(ctx)
		require.NoError(t, err)
		session, err := sessionUseCase.CreateSession(ctx, token, user.ID)
		require.NoError(t, err)
		require.NotNil(t, session)

		// Create request to protected route
		req := httptest.NewRequest(http.MethodGet, "/api/me", nil)

		// Add session cookie
		req.AddCookie(&http.Cookie{
			Name:  "session",
			Value: token,
		})

		// Perform request
		recorder := httptest.NewRecorder()
		r.ServeHTTP(recorder, req)

		// Check response
		require.Equal(t, http.StatusOK, recorder.Code)

		var response map[string]interface{}
		err = json.Unmarshal(recorder.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, user.ID, response["id"])
		assert.Equal(t, user.Email, response["email"])
		assert.Equal(t, user.Name, response["name"])
	})

	t.Run("Logout", func(t *testing.T) {
		// First register and login a user
		email := "logouttest@example.com"
		name := "Logout Test"
		password := "L0g0ut!P@ssword123"

		user, err := userUseCase.RegisterUser(ctx, email, name, password)
		require.NoError(t, err)

		// Generate session token and create session
		token, err := sessionUseCase.GenerateSessionToken(ctx)
		require.NoError(t, err)
		session, err := sessionUseCase.CreateSession(ctx, token, user.ID)
		require.NoError(t, err)
		require.NotNil(t, session)

		// Create logout request
		req := httptest.NewRequest(http.MethodPost, "/api/logout", nil)

		// Add session cookie
		req.AddCookie(&http.Cookie{
			Name:  "session",
			Value: token,
		})

		// Perform request
		recorder := httptest.NewRecorder()
		r.ServeHTTP(recorder, req)

		// Check response
		require.Equal(t, http.StatusNoContent, recorder.Code)

		// Check that session cookie is cleared
		cookies := recorder.Result().Cookies()
		var sessionCookie *http.Cookie
		for _, cookie := range cookies {
			if cookie.Name == "session" {
				sessionCookie = cookie
				break
			}
		}
		require.NotNil(t, sessionCookie, "Session cookie should be present")
		assert.Equal(t, "", sessionCookie.Value, "Session cookie value should be empty")
		assert.True(t, sessionCookie.MaxAge < 0, "Session cookie should be expired")

		// Try to access protected route after logout
		protectedReq := httptest.NewRequest(http.MethodGet, "/api/me", nil)
		protectedReq.AddCookie(&http.Cookie{
			Name:  "session",
			Value: token,
		})

		protectedRecorder := httptest.NewRecorder()
		r.ServeHTTP(protectedRecorder, protectedReq)

		// Should return unauthorized
		assert.Equal(t, http.StatusUnauthorized, protectedRecorder.Code)
	})
}
