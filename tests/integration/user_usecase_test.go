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

func TestUserUseCaseIntegration(t *testing.T) {
	// Set up test database
	ctx := context.Background()
	db, err := SetupTestDatabase(ctx)
	require.NoError(t, err)
	defer db.Close(ctx)

	// Initialize repositories
	queries := repositories.New(db.Pool)
	userRepo := repositories.NewUserRepository(queries)

	// Initialize services
	hashService := services.NewArgonHashService()

	// Initialize use case
	userUseCase := use_cases.NewUserUseCase(userRepo, hashService)

	t.Run("RegisterAndAuthenticateUser", func(t *testing.T) {
		// Register a new user
		email := "test@example.com"
		name := "Test User"
		password := "TestUs3r!P@ssw0rd"

		user, err := userUseCase.RegisterUser(ctx, email, name, password)
		require.NoError(t, err)
		require.NotNil(t, user)
		assert.Equal(t, email, user.Email)
		assert.Equal(t, name, user.Name)
		assert.Empty(t, user.Password) // Password should not be returned

		// Authenticate with correct credentials
		authenticatedUser, err := userUseCase.AuthenticateUser(ctx, email, password)
		require.NoError(t, err)
		require.NotNil(t, authenticatedUser)
		assert.Equal(t, user.ID, authenticatedUser.ID)
		assert.Equal(t, email, authenticatedUser.Email)
		assert.Empty(t, authenticatedUser.Password)

		// Try to authenticate with incorrect password
		invalidUser, err := userUseCase.AuthenticateUser(ctx, email, "Wr0ng!P@ssw0rd123")
		assert.Error(t, err)
		assert.Nil(t, invalidUser)

		// Try to authenticate with non-existent email
		nonExistentUser, err := userUseCase.AuthenticateUser(ctx, "nonexistent@example.com", password)
		assert.Error(t, err)
		assert.Nil(t, nonExistentUser)
	})

	t.Run("RegisterDuplicateEmail", func(t *testing.T) {
		// Register a user
		email := "duplicate@example.com"
		name := "First User"
		password := "FirstUs3r!P@ssw0rd"

		user, err := userUseCase.RegisterUser(ctx, email, name, password)
		require.NoError(t, err)
		require.NotNil(t, user)

		// Try to register another user with the same email
		duplicateUser, err := userUseCase.RegisterUser(ctx, email, "Second User", "Sec0ndUs3r!P@ssw0rd")
		assert.Error(t, err)
		assert.Nil(t, duplicateUser)
		assert.Equal(t, "email already taken", err.Error())
	})
}
