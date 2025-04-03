package integration

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/LaulauChau/note-nest/internal/domain/entities"
	"github.com/LaulauChau/note-nest/internal/infrastructure/persistence/repositories"
)

func TestSessionRepository(t *testing.T) {
	// Set up test database
	ctx := context.Background()
	db, err := SetupTestDatabase(ctx)
	require.NoError(t, err)
	defer db.Close(ctx)

	// Create repositories
	queries := repositories.New(db.Pool)
	userRepo := repositories.NewUserRepository(queries)
	sessionRepo := repositories.NewSessionRepository(queries)

	// Create a test user first
	now := time.Now()
	user := &entities.User{
		ID:        uuid.New().String(),
		Email:     "test@example.com",
		Name:      "Test User",
		Password:  "hashedpassword",
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Save the test user
	err = userRepo.Create(ctx, user)
	require.NoError(t, err)

	t.Run("CreateAndGetSession", func(t *testing.T) {
		// Create a session
		sessionID := "sessionhashid123456789"
		expiresAt := time.Now().Add(24 * time.Hour)
		session := &entities.Session{
			ID:        sessionID,
			UserID:    user.ID,
			ExpiresAt: expiresAt,
			CreatedAt: now,
		}

		// Save the session
		err := sessionRepo.Create(ctx, session)
		require.NoError(t, err)

		// Retrieve the session
		retrievedSession, err := sessionRepo.GetByID(ctx, sessionID)
		require.NoError(t, err)
		require.NotNil(t, retrievedSession)

		assert.Equal(t, sessionID, retrievedSession.ID)
		assert.Equal(t, user.ID, retrievedSession.UserID)
		assert.WithinDuration(t, expiresAt, retrievedSession.ExpiresAt, time.Second)
	})

	t.Run("GetSessionWithUser", func(t *testing.T) {
		// Create a session
		sessionID := "anothersessionhash987654321"
		expiresAt := time.Now().Add(24 * time.Hour)
		session := &entities.Session{
			ID:        sessionID,
			UserID:    user.ID,
			ExpiresAt: expiresAt,
			CreatedAt: now,
		}

		// Save the session
		err := sessionRepo.Create(ctx, session)
		require.NoError(t, err)

		// Get session with user
		result, err := sessionRepo.GetSessionWithUser(ctx, sessionID)
		require.NoError(t, err)
		require.NotNil(t, result)
		require.NotNil(t, result.Session)
		require.NotNil(t, result.User)

		assert.Equal(t, sessionID, result.Session.ID)
		assert.Equal(t, user.ID, result.Session.UserID)
		assert.Equal(t, user.ID, result.User.ID)
		assert.Equal(t, user.Email, result.User.Email)
		assert.Equal(t, user.Name, result.User.Name)
	})

	t.Run("UpdateExpiresAt", func(t *testing.T) {
		// Create a session
		sessionID := "updatesessionhash123456789"
		expiresAt := time.Now().Add(24 * time.Hour)
		session := &entities.Session{
			ID:        sessionID,
			UserID:    user.ID,
			ExpiresAt: expiresAt,
			CreatedAt: now,
		}

		// Save the session
		err := sessionRepo.Create(ctx, session)
		require.NoError(t, err)

		// Update expires_at
		newExpiresAt := time.Now().Add(48 * time.Hour)
		err = sessionRepo.UpdateExpiresAt(ctx, sessionID, newExpiresAt)
		require.NoError(t, err)

		// Get the session to verify the update
		updatedSession, err := sessionRepo.GetByID(ctx, sessionID)
		require.NoError(t, err)
		require.NotNil(t, updatedSession)

		assert.WithinDuration(t, newExpiresAt, updatedSession.ExpiresAt, time.Second)
	})

	t.Run("DeleteSession", func(t *testing.T) {
		// Create a session
		sessionID := "deletesessionhash123456789"
		expiresAt := time.Now().Add(24 * time.Hour)
		session := &entities.Session{
			ID:        sessionID,
			UserID:    user.ID,
			ExpiresAt: expiresAt,
			CreatedAt: now,
		}

		// Save the session
		err := sessionRepo.Create(ctx, session)
		require.NoError(t, err)

		// Delete the session
		err = sessionRepo.Delete(ctx, sessionID)
		require.NoError(t, err)

		// Try to get the session - should be nil
		deletedSession, err := sessionRepo.GetByID(ctx, sessionID)
		require.NoError(t, err)
		assert.Nil(t, deletedSession)
	})

	t.Run("DeleteAllSessionsByUserID", func(t *testing.T) {
		// Create multiple sessions for the user
		sessionID1 := "usersession1hash123456789"
		sessionID2 := "usersession2hash123456789"
		expiresAt := time.Now().Add(24 * time.Hour)

		session1 := &entities.Session{
			ID:        sessionID1,
			UserID:    user.ID,
			ExpiresAt: expiresAt,
			CreatedAt: now,
		}

		session2 := &entities.Session{
			ID:        sessionID2,
			UserID:    user.ID,
			ExpiresAt: expiresAt,
			CreatedAt: now,
		}

		// Save the sessions
		err := sessionRepo.Create(ctx, session1)
		require.NoError(t, err)

		err = sessionRepo.Create(ctx, session2)
		require.NoError(t, err)

		// Delete all sessions for the user
		err = sessionRepo.DeleteAllByUserID(ctx, user.ID)
		require.NoError(t, err)

		// Verify sessions are gone
		s1, err := sessionRepo.GetByID(ctx, sessionID1)
		require.NoError(t, err)
		assert.Nil(t, s1)

		s2, err := sessionRepo.GetByID(ctx, sessionID2)
		require.NoError(t, err)
		assert.Nil(t, s2)
	})
}
