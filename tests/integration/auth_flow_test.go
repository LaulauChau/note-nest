package integration

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/LaulauChau/note-nest/internal/application/use_cases"
	"github.com/LaulauChau/note-nest/internal/infrastructure/persistence/repositories"
	"github.com/LaulauChau/note-nest/internal/infrastructure/services"
)

func TestAuthFlowIntegration(t *testing.T) {
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

		// 2. Authenticate the user
		authenticatedUser, err := userUseCase.AuthenticateUser(ctx, email, password)
		require.NoError(t, err)
		require.NotNil(t, authenticatedUser)
		assert.Equal(t, user.ID, authenticatedUser.ID)

		// 3. Create a session for the authenticated user
		token, err := sessionUseCase.GenerateSessionToken(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, token)

		session, err := sessionUseCase.CreateSession(ctx, token, authenticatedUser.ID)
		require.NoError(t, err)
		require.NotNil(t, session)

		// 4. Validate the session (simulates subsequent requests)
		validationResult, err := sessionUseCase.ValidateSessionToken(ctx, token)
		require.NoError(t, err)
		require.NotNil(t, validationResult)
		require.NotNil(t, validationResult.Session)
		require.NotNil(t, validationResult.User)
		assert.Equal(t, authenticatedUser.ID, validationResult.User.ID)

		// 5. Logout (invalidate the session)
		err = sessionUseCase.InvalidateSession(ctx, session.ID)
		require.NoError(t, err)

		// 6. Try to use the invalidated session
		invalidResult, err := sessionUseCase.ValidateSessionToken(ctx, token)
		require.NoError(t, err)
		assert.Nil(t, invalidResult.Session)
		assert.Nil(t, invalidResult.User)

		// 7. Login again and create a new session
		authenticatedUser, err = userUseCase.AuthenticateUser(ctx, email, password)
		require.NoError(t, err)

		newToken, err := sessionUseCase.GenerateSessionToken(ctx)
		require.NoError(t, err)

		newSession, err := sessionUseCase.CreateSession(ctx, newToken, authenticatedUser.ID)
		require.NoError(t, err)
		require.NotNil(t, newSession)

		// 8. Verify new session is valid
		newValidationResult, err := sessionUseCase.ValidateSessionToken(ctx, newToken)
		require.NoError(t, err)
		require.NotNil(t, newValidationResult.Session)
		require.NotNil(t, newValidationResult.User)

		// 9. Test session invalidation by changing password (security feature)
		// This would normally be in a separate use case, but for testing we'll mock it
		// by directly invalidating all sessions
		err = sessionUseCase.InvalidateAllSessions(ctx, authenticatedUser.ID)
		require.NoError(t, err)

		// 10. Verify the session is no longer valid
		invalidatedResult, err := sessionUseCase.ValidateSessionToken(ctx, newToken)
		require.NoError(t, err)
		assert.Nil(t, invalidatedResult.Session)
		assert.Nil(t, invalidatedResult.User)
	})
}
