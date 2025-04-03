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

func TestNoteUseCaseIntegration(t *testing.T) {
	// Set up test database
	ctx := context.Background()
	db, err := SetupTestDatabase(ctx)
	require.NoError(t, err)
	defer db.Close(ctx)

	// Create repositories
	queries := repositories.New(db.Pool)
	userRepo := repositories.NewUserRepository(queries)
	noteRepo := repositories.NewNoteRepository(queries)

	// Initialize services
	hashService := services.NewArgonHashService()

	// Initialize use cases
	userUseCase := use_cases.NewUserUseCase(userRepo, hashService)
	noteUseCase := use_cases.NewNoteUseCase(noteRepo, userRepo)

	// Create two test users
	email1 := "user1@example.com"
	email2 := "user2@example.com"
	password := "password123"

	user1, err := userUseCase.RegisterUser(ctx, email1, "User One", password)
	require.NoError(t, err)
	require.NotNil(t, user1)

	user2, err := userUseCase.RegisterUser(ctx, email2, "User Two", password)
	require.NoError(t, err)
	require.NotNil(t, user2)

	t.Run("CreateAndRetrieveNote", func(t *testing.T) {
		// User 1 creates a note
		note, err := noteUseCase.CreateNote(ctx, user1.ID, "User 1's Note", "This belongs to user 1", "user1-label")
		require.NoError(t, err)
		require.NotNil(t, note)
		assert.Equal(t, user1.ID, note.UserID)

		// User 1 can retrieve their note
		retrievedNote, err := noteUseCase.GetNoteByID(ctx, note.ID, user1.ID)
		require.NoError(t, err)
		require.NotNil(t, retrievedNote)
		assert.Equal(t, note.ID, retrievedNote.ID)

		// User 2 cannot retrieve User 1's note
		_, err = noteUseCase.GetNoteByID(ctx, note.ID, user2.ID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "note not found")
	})

	t.Run("UpdateNote", func(t *testing.T) {
		// User 1 creates a note
		note, err := noteUseCase.CreateNote(ctx, user1.ID, "Note to Update", "Original content", "update-test")
		require.NoError(t, err)
		require.NotNil(t, note)

		// User 2 tries to update User 1's note
		_, err = noteUseCase.UpdateNote(ctx, note.ID, user2.ID, "Updated by User 2", "This should fail", "user2-label", true)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "note not found")

		// User 1 updates their own note
		updatedNote, err := noteUseCase.UpdateNote(ctx, note.ID, user1.ID, "Updated by User 1", "This should work", "user1-label", true)
		require.NoError(t, err)
		require.NotNil(t, updatedNote)
		assert.Equal(t, "Updated by User 1", updatedNote.Title)
		assert.Equal(t, "This should work", updatedNote.Content)
		assert.Equal(t, "user1-label", updatedNote.Label)
		assert.True(t, updatedNote.IsArchived)
	})

	t.Run("DeleteNote", func(t *testing.T) {
		// User 1 creates a note
		note, err := noteUseCase.CreateNote(ctx, user1.ID, "Note to Delete", "Delete content", "delete-test")
		require.NoError(t, err)
		require.NotNil(t, note)

		// User 2 tries to delete User 1's note
		err = noteUseCase.DeleteNote(ctx, note.ID, user2.ID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "note not found")

		// Verify the note still exists
		retrievedNote, err := noteUseCase.GetNoteByID(ctx, note.ID, user1.ID)
		require.NoError(t, err)
		require.NotNil(t, retrievedNote)

		// User 1 deletes their own note
		err = noteUseCase.DeleteNote(ctx, note.ID, user1.ID)
		require.NoError(t, err)

		// Verify the note is deleted
		_, err = noteUseCase.GetNoteByID(ctx, note.ID, user1.ID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "note not found")
	})

	t.Run("GetNotesByUserID", func(t *testing.T) {
		// User 1 creates multiple notes
		_, err := noteUseCase.CreateNote(ctx, user1.ID, "User 1 Note 1", "Content 1", "label1")
		require.NoError(t, err)

		_, err = noteUseCase.CreateNote(ctx, user1.ID, "User 1 Note 2", "Content 2", "label2")
		require.NoError(t, err)

		// User 2 creates a note
		_, err = noteUseCase.CreateNote(ctx, user2.ID, "User 2 Note", "User 2 Content", "user2-label")
		require.NoError(t, err)

		// Get User 1's notes
		user1Notes, err := noteUseCase.GetActiveNotes(ctx, user1.ID)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(user1Notes), 2)

		// Verify all notes belong to User 1
		for _, note := range user1Notes {
			assert.Equal(t, user1.ID, note.UserID)
		}

		// Get User 2's notes
		user2Notes, err := noteUseCase.GetActiveNotes(ctx, user2.ID)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(user2Notes), 1)

		// Verify all notes belong to User 2
		for _, note := range user2Notes {
			assert.Equal(t, user2.ID, note.UserID)
		}
	})
}
