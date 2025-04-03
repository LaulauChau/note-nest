package integration

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/LaulauChau/note-nest/internal/application/use_cases"
	"github.com/LaulauChau/note-nest/internal/domain/entities"
	"github.com/LaulauChau/note-nest/internal/infrastructure/persistence/repositories"
	"github.com/LaulauChau/note-nest/internal/infrastructure/services"
)

func TestSessionUseCaseIntegration(t *testing.T) {
	// Set up test database
	ctx := context.Background()
	db, err := SetupTestDatabase(ctx)
	require.NoError(t, err)
	defer db.Close(ctx)

	// Initialize the SQLC queries struct
	queries := repositories.New(db.Pool)

	// Initialize repositories
	userRepo := repositories.NewUserRepository(queries)
	sessionRepo := repositories.NewSessionRepository(queries)

	// Initialize services
	tokenService := services.NewTokenService()
	hashService := services.NewArgonHashService()

	// Initialize use cases
	userUseCase := use_cases.NewUserUseCase(userRepo, hashService)
	sessionUseCase := use_cases.NewSessionUseCase(sessionRepo, userRepo, tokenService)

	// Test user
	email := "test@example.com"
	name := "Test User"
	password := "S3ssionT3st!P@ss123"

	// Register a test user
	user, err := userUseCase.RegisterUser(ctx, email, name, password)
	require.NoError(t, err)
	require.NotNil(t, user)

	t.Run("CreateAndValidateSession", func(t *testing.T) {
		// Generate a session token
		token, err := sessionUseCase.GenerateSessionToken(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, token)

		// Create a session
		session, err := sessionUseCase.CreateSession(ctx, token, user.ID)
		require.NoError(t, err)
		require.NotNil(t, session)

		// Validate the session token
		result, err := sessionUseCase.ValidateSessionToken(ctx, token)
		require.NoError(t, err)
		require.NotNil(t, result)
		require.NotNil(t, result.Session)
		require.NotNil(t, result.User)
		assert.Equal(t, user.ID, result.User.ID)
	})

	t.Run("InvalidateSession", func(t *testing.T) {
		// Generate a session token
		token, err := sessionUseCase.GenerateSessionToken(ctx)
		require.NoError(t, err)

		// Create a session
		session, err := sessionUseCase.CreateSession(ctx, token, user.ID)
		require.NoError(t, err)
		require.NotNil(t, session)

		// Invalidate the session
		err = sessionUseCase.InvalidateSession(ctx, session.ID)
		require.NoError(t, err)

		// Try to validate the invalidated session
		result, err := sessionUseCase.ValidateSessionToken(ctx, token)
		require.NoError(t, err)
		assert.Nil(t, result.Session)
		assert.Nil(t, result.User)
	})

	t.Run("ExpiredSession", func(t *testing.T) {
		// Create a session that expires immediately
		shortExpiryToken, err := sessionUseCase.GenerateSessionToken(ctx)
		require.NoError(t, err)

		// Store the token mapping in a way that ValidateSessionToken will find it
		// This is a bit of a hack for testing, but it's necessary
		hashedToken, err := tokenService.HashToken(ctx, shortExpiryToken)
		require.NoError(t, err)

		// Create a new session with the hashed token as ID instead of trying to update
		// In a real system, the hashed token would be the session ID from creation
		expiredSession := &entities.Session{
			ID:        hashedToken,
			UserID:    user.ID,
			ExpiresAt: time.Now().Add(-1 * time.Hour),
		}
		err = sessionRepo.Create(ctx, expiredSession)
		require.NoError(t, err)

		// Validate the session token - should fail due to expiration
		result, err := sessionUseCase.ValidateSessionToken(ctx, shortExpiryToken)
		require.NoError(t, err)
		assert.Nil(t, result.Session)
		assert.Nil(t, result.User)
	})

	t.Run("InvalidateAllSessions", func(t *testing.T) {
		// Create multiple sessions for the user
		token1, _ := sessionUseCase.GenerateSessionToken(ctx)
		token2, _ := sessionUseCase.GenerateSessionToken(ctx)

		session1, _ := sessionUseCase.CreateSession(ctx, token1, user.ID)
		session2, _ := sessionUseCase.CreateSession(ctx, token2, user.ID)

		require.NotNil(t, session1)
		require.NotNil(t, session2)

		// Invalidate all sessions for the user
		err = sessionUseCase.InvalidateAllSessions(ctx, user.ID)
		require.NoError(t, err)

		// Try to validate both sessions
		result1, _ := sessionUseCase.ValidateSessionToken(ctx, token1)
		result2, _ := sessionUseCase.ValidateSessionToken(ctx, token2)

		assert.Nil(t, result1.Session)
		assert.Nil(t, result2.Session)
	})
}
