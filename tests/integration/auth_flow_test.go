package integration

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/LaulauChau/note-nest/internal/application/use_cases"
	"github.com/LaulauChau/note-nest/internal/infrastructure/persistence/repositories"
	"github.com/LaulauChau/note-nest/internal/infrastructure/services"
)

func TestAuthenticationFlow(t *testing.T) {
	// Set up test database
	ctx := context.Background()
	db, err := SetupTestDatabase(ctx)
	require.NoError(t, err)
	defer db.Close(ctx)

	// Initialize repositories
	queries := repositories.New(db.Pool)
	userRepo := repositories.NewUserRepository(queries)
	sessionRepo := repositories.NewSessionRepository(queries)

	// Initialize services
	tokenService := services.NewTokenService()
	hashService := services.NewArgonHashService()

	// Initialize use cases
	userUseCase := use_cases.NewUserUseCase(userRepo, hashService)
	sessionUseCase := use_cases.NewSessionUseCase(sessionRepo, userRepo, tokenService)

	// Test data
	email := "authflow@example.com"
	name := "Auth Flow User"
	password := "securepassword"

	t.Run("CompleteAuthFlow", func(t *testing.T) {
		// 1. Register a user
		user, err := userUseCase.RegisterUser(ctx, email, name, password)
		require.NoError(t, err)
		require.NotNil(t, user)
		assert.Equal(t, email, user.Email)
		assert.Equal(t, name, user.Name)
		assert.NotEmpty(t, user.ID)

		// 2. Authenticate the user with correct credentials
		authenticatedUser, err := userUseCase.AuthenticateUser(ctx, email, password)
		require.NoError(t, err)
		require.NotNil(t, authenticatedUser)
		assert.Equal(t, user.ID, authenticatedUser.ID)

		// 3. Try to authenticate with incorrect password
		invalidUser, err := userUseCase.AuthenticateUser(ctx, email, "wrongpassword")
		assert.Error(t, err)
		assert.Nil(t, invalidUser)

		// 4. Create a session for the authenticated user
		token, err := sessionUseCase.GenerateSessionToken(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, token)

		session, err := sessionUseCase.CreateSession(ctx, token, authenticatedUser.ID)
		require.NoError(t, err)
		require.NotNil(t, session)
		assert.Equal(t, authenticatedUser.ID, session.UserID)

		// 5. Validate the session token (simulates subsequent requests)
		validationResult, err := sessionUseCase.ValidateSessionToken(ctx, token)
		require.NoError(t, err)
		require.NotNil(t, validationResult)
		require.NotNil(t, validationResult.Session)
		require.NotNil(t, validationResult.User)
		assert.Equal(t, authenticatedUser.ID, validationResult.User.ID)
		assert.Equal(t, authenticatedUser.Email, validationResult.User.Email)

		// 6. Logout (invalidate the session)
		err = sessionUseCase.InvalidateSession(ctx, session.ID)
		require.NoError(t, err)

		// 7. Try to use the invalidated session
		invalidResult, err := sessionUseCase.ValidateSessionToken(ctx, token)
		require.NoError(t, err)
		assert.Nil(t, invalidResult.Session)
		assert.Nil(t, invalidResult.User)

		// 8. Login again and create a new session
		authenticatedUser, err = userUseCase.AuthenticateUser(ctx, email, password)
		require.NoError(t, err)

		newToken, err := sessionUseCase.GenerateSessionToken(ctx)
		require.NoError(t, err)

		newSession, err := sessionUseCase.CreateSession(ctx, newToken, authenticatedUser.ID)
		require.NoError(t, err)
		require.NotNil(t, newSession)

		// 9. Test session renewal (when close to expiry)
		// To test this properly, we'd need to create a session with a short expiry
		// For now, we'll just directly modify the session to simulate this
		shortExpiryTime := time.Now().Add(10 * 24 * time.Hour) // 10 days from now
		err = sessionRepo.UpdateExpiresAt(ctx, newSession.ID, shortExpiryTime)
		require.NoError(t, err)

		// Validate session - should trigger renewal since it's within 15 days of expiry
		renewalResult, err := sessionUseCase.ValidateSessionToken(ctx, newToken)
		require.NoError(t, err)
		require.NotNil(t, renewalResult.Session)

		// Should be renewed for 30 days
		expectedExpiryMin := time.Now().Add(29 * 24 * time.Hour)
		assert.True(t, renewalResult.Session.ExpiresAt.After(expectedExpiryMin),
			"Session should be renewed to approximately 30 days from now")

		// 10. Test invalidation of all sessions (e.g., after password change)
		err = sessionUseCase.InvalidateAllSessions(ctx, authenticatedUser.ID)
		require.NoError(t, err)

		// Verify the session is invalid now
		afterInvalidateResult, err := sessionUseCase.ValidateSessionToken(ctx, newToken)
		require.NoError(t, err)
		assert.Nil(t, afterInvalidateResult.Session)
		assert.Nil(t, afterInvalidateResult.User)
	})
}
