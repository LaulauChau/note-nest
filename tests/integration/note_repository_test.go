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
		// Create a separate test user for this test
		createUserID := uuid.New().String()
		createUser := &entities.User{
			ID:        createUserID,
			Email:     "createuser@example.com",
			Name:      "Create Test User",
			Password:  "hashedpassword",
			CreatedAt: now,
			UpdatedAt: now,
		}
		// Save the test user
		err = userRepo.Create(ctx, createUser)
		require.NoError(t, err)

		// Create a note
		noteID := uuid.New().String()
		note := &entities.Note{
			ID:         noteID,
			UserID:     createUserID,
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
		assert.Equal(t, createUserID, retrievedNote.UserID)
		assert.Equal(t, "Test Note", retrievedNote.Title)
		assert.Equal(t, "This is a test note content.", retrievedNote.Content)
		assert.Equal(t, "test-label", retrievedNote.Label)
		assert.False(t, retrievedNote.IsArchived)
	})

	t.Run("GetNotesByUserID", func(t *testing.T) {
		// Create a separate test user for this test
		testUserID := uuid.New().String()
		testUser := &entities.User{
			ID:        testUserID,
			Email:     "notesuser@example.com",
			Name:      "Notes Test User",
			Password:  "hashedpassword",
			CreatedAt: now,
			UpdatedAt: now,
		}
		// Save the test user
		err = userRepo.Create(ctx, testUser)
		require.NoError(t, err)

		// Create multiple notes (both archived and non-archived)
		note1 := &entities.Note{
			ID:         uuid.New().String(),
			UserID:     testUserID,
			Title:      "Active Note 1",
			Content:    "Active content 1",
			IsArchived: false,
			CreatedAt:  now,
			UpdatedAt:  now,
		}

		note2 := &entities.Note{
			ID:         uuid.New().String(),
			UserID:     testUserID,
			Title:      "Active Note 2",
			Content:    "Active content 2",
			IsArchived: false,
			CreatedAt:  now,
			UpdatedAt:  now,
		}

		note3 := &entities.Note{
			ID:         uuid.New().String(),
			UserID:     testUserID,
			Title:      "Archived Note",
			Content:    "Archived content",
			IsArchived: true,
			CreatedAt:  now,
			UpdatedAt:  now,
		}

		// Save all notes
		require.NoError(t, noteRepo.Create(ctx, note1))
		require.NoError(t, noteRepo.Create(ctx, note2))
		require.NoError(t, noteRepo.Create(ctx, note3))

		// Get active notes
		activeNotes, err := noteRepo.GetByUserID(ctx, testUserID)
		require.NoError(t, err)

		// Should only return the 2 active notes
		assert.Len(t, activeNotes, 2)

		// Check that both active notes are returned
		activeNoteIDs := map[string]bool{
			note1.ID: false,
			note2.ID: false,
		}

		for _, note := range activeNotes {
			activeNoteIDs[note.ID] = true
			assert.False(t, note.IsArchived)
		}

		// Verify we found both active notes
		for id, found := range activeNoteIDs {
			assert.True(t, found, "Active note %s was not returned", id)
		}
	})

	t.Run("GetArchivedNotesByUserID", func(t *testing.T) {
		// Create a separate test user for this test
		archivedUserID := uuid.New().String()
		archivedUser := &entities.User{
			ID:        archivedUserID,
			Email:     "archiveduser@example.com",
			Name:      "Archived Notes User",
			Password:  "hashedpassword",
			CreatedAt: now,
			UpdatedAt: now,
		}
		// Save the test user
		err = userRepo.Create(ctx, archivedUser)
		require.NoError(t, err)

		// Create archived notes
		note3 := &entities.Note{
			ID:         uuid.New().String(),
			UserID:     archivedUserID,
			Title:      "Archived Note 1",
			Content:    "Archived content 1",
			IsArchived: true,
			CreatedAt:  now,
			UpdatedAt:  now,
		}

		note4 := &entities.Note{
			ID:         uuid.New().String(),
			UserID:     archivedUserID,
			Title:      "Archived Note 2",
			Content:    "Archived content 2",
			IsArchived: true,
			CreatedAt:  now,
			UpdatedAt:  now,
		}

		// Save the notes
		require.NoError(t, noteRepo.Create(ctx, note3))
		require.NoError(t, noteRepo.Create(ctx, note4))

		// Get archived notes
		archivedNotes, err := noteRepo.GetArchivedByUserID(ctx, archivedUserID)
		require.NoError(t, err)

		// Should return exactly 2 archived notes
		assert.Len(t, archivedNotes, 2)

		// Check that all returned notes are archived
		for _, note := range archivedNotes {
			assert.True(t, note.IsArchived)
		}
	})

	t.Run("GetNoteByID", func(t *testing.T) {
		// Create a separate test user for this test
		idUserID := uuid.New().String()
		idUser := &entities.User{
			ID:        idUserID,
			Email:     "iduser@example.com",
			Name:      "ID Test User",
			Password:  "hashedpassword",
			CreatedAt: now,
			UpdatedAt: now,
		}
		// Save the test user
		err = userRepo.Create(ctx, idUser)
		require.NoError(t, err)

		// Create a note
		noteID := uuid.New().String()
		note := &entities.Note{
			ID:         noteID,
			UserID:     idUserID,
			Title:      "Get By ID Note",
			Content:    "Get by ID content",
			IsArchived: false,
			Label:      "get-by-id",
			CreatedAt:  now,
			UpdatedAt:  now,
		}

		// Save the note
		require.NoError(t, noteRepo.Create(ctx, note))

		// Get the note by ID
		retrievedNote, err := noteRepo.GetByID(ctx, noteID)
		require.NoError(t, err)
		require.NotNil(t, retrievedNote)

		// Verify note details
		assert.Equal(t, noteID, retrievedNote.ID)
		assert.Equal(t, idUserID, retrievedNote.UserID)
		assert.Equal(t, "Get By ID Note", retrievedNote.Title)
		assert.Equal(t, "Get by ID content", retrievedNote.Content)
		assert.Equal(t, "get-by-id", retrievedNote.Label)
		assert.False(t, retrievedNote.IsArchived)
	})

	t.Run("GetNoteByID_NotFound", func(t *testing.T) {
		// Try to get a non-existent note
		nonExistentID := uuid.New().String()
		retrievedNote, err := noteRepo.GetByID(ctx, nonExistentID)

		// Should not return an error, just nil
		require.NoError(t, err)
		assert.Nil(t, retrievedNote)
	})
}
