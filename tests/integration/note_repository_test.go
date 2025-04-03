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

func TestNoteRepository(t *testing.T) {
	// Set up test database
	ctx := context.Background()
	db, err := SetupTestDatabase(ctx)
	require.NoError(t, err)
	defer db.Close(ctx)

	// Create repositories
	queries := repositories.New(db.Pool)
	userRepo := repositories.NewUserRepository(queries)
	noteRepo := repositories.NewNoteRepository(queries)

	// Create a test user first
	now := time.Now()
	user := &entities.User{
		ID:        uuid.New().String(),
		Email:     "notetest@example.com",
		Name:      "Note Test User",
		Password:  "hashedpassword",
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Save the test user
	err = userRepo.Create(ctx, user)
	require.NoError(t, err)

	t.Run("CreateAndGetNote", func(t *testing.T) {
		// Create a note
		noteID := uuid.New().String()
		note := &entities.Note{
			ID:         noteID,
			UserID:     user.ID,
			Title:      "Test Note",
			Content:    "This is a test note content.",
			IsArchived: false,
			Label:      "test-label",
			CreatedAt:  now,
			UpdatedAt:  now,
		}

		// Save the note
		err := noteRepo.Create(ctx, note)
		require.NoError(t, err)

		// Retrieve the note
		retrievedNote, err := noteRepo.GetByID(ctx, noteID)
		require.NoError(t, err)
		require.NotNil(t, retrievedNote)

		assert.Equal(t, noteID, retrievedNote.ID)
		assert.Equal(t, user.ID, retrievedNote.UserID)
		assert.Equal(t, "Test Note", retrievedNote.Title)
		assert.Equal(t, "This is a test note content.", retrievedNote.Content)
		assert.Equal(t, "test-label", retrievedNote.Label)
		assert.False(t, retrievedNote.IsArchived)
	})
}
