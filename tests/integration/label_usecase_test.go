package integration

import (
	"context"
	"testing"

	"github.com/LaulauChau/note-nest/internal/application/use_cases"
	"github.com/LaulauChau/note-nest/internal/infrastructure/persistence/repositories"
	"github.com/LaulauChau/note-nest/internal/infrastructure/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLabelUseCaseIntegration(t *testing.T) {
	// Set up test database
	ctx := context.Background()
	db, err := SetupTestDatabase(ctx)
	require.NoError(t, err)
	defer db.Close(ctx)

	// Create repositories
	queries := repositories.New(db.Pool)
	userRepo := repositories.NewUserRepository(queries)
	noteRepo := repositories.NewNoteRepository(queries)
	labelRepo := repositories.NewLabelRepository(queries)

	// Initialize services
	hashService := services.NewArgonHashService()

	// Initialize use cases
	userUseCase := use_cases.NewUserUseCase(userRepo, hashService)
	noteUseCase := use_cases.NewNoteUseCase(noteRepo, userRepo, labelRepo)
	labelUseCase := use_cases.NewLabelUseCase(labelRepo, userRepo, noteRepo)

	// Create two test users
	email1 := "labeluser1@example.com"
	email2 := "labeluser2@example.com"
	password := "password123"

	user1, err := userUseCase.RegisterUser(ctx, email1, "Label User One", password)
	require.NoError(t, err)
	require.NotNil(t, user1)

	user2, err := userUseCase.RegisterUser(ctx, email2, "Label User Two", password)
	require.NoError(t, err)
	require.NotNil(t, user2)

	t.Run("CreateAndGetLabel", func(t *testing.T) {
		labelName := "My First Label"
		labelColor := "#1ABC9C"

		// User 1 creates a label
		label, err := labelUseCase.CreateLabel(ctx, user1.ID, labelName, labelColor)
		require.NoError(t, err)
		require.NotNil(t, label)
		assert.Equal(t, user1.ID, label.UserID)
		assert.Equal(t, labelName, label.Name)
		assert.Equal(t, labelColor, label.Color)

		// User 1 can retrieve their label
		retrievedLabel, err := labelUseCase.GetLabelByID(ctx, label.ID, user1.ID)
		require.NoError(t, err)
		require.NotNil(t, retrievedLabel)
		assert.Equal(t, label.ID, retrievedLabel.ID)

		// User 2 cannot retrieve User 1's label
		_, err = labelUseCase.GetLabelByID(ctx, label.ID, user2.ID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "label not found")

		// Cannot create label with duplicate name for the same user
		_, err = labelUseCase.CreateLabel(ctx, user1.ID, labelName, "#FFFFFF")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")

		// Can create label with same name for different user
		_, err = labelUseCase.CreateLabel(ctx, user2.ID, labelName, labelColor)
		assert.NoError(t, err)
	})

	t.Run("GetLabelsByUser", func(t *testing.T) {
		// User 1 creates multiple labels
		_, err := labelUseCase.CreateLabel(ctx, user1.ID, "User1 Label A", "#AAA")
		require.NoError(t, err)
		_, err = labelUseCase.CreateLabel(ctx, user1.ID, "User1 Label B", "#BBB")
		require.NoError(t, err)

		// User 2 creates a label
		_, err = labelUseCase.CreateLabel(ctx, user2.ID, "User2 Label C", "#CCC")
		require.NoError(t, err)

		// Get User 1's labels
		user1Labels, err := labelUseCase.GetLabelsByUser(ctx, user1.ID)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(user1Labels), 2)
		for _, l := range user1Labels {
			assert.Equal(t, user1.ID, l.UserID)
		}

		// Get User 2's labels
		user2Labels, err := labelUseCase.GetLabelsByUser(ctx, user2.ID)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(user2Labels), 1)
		for _, l := range user2Labels {
			assert.Equal(t, user2.ID, l.UserID)
		}
	})

	t.Run("UpdateLabel", func(t *testing.T) {
		// User 1 creates a label
		label, err := labelUseCase.CreateLabel(ctx, user1.ID, "Label To Update", "#123456")
		require.NoError(t, err)
		require.NotNil(t, label)

		// User 2 tries to update User 1's label
		_, err = labelUseCase.UpdateLabel(ctx, label.ID, user2.ID, "Updated by User 2", "#FFFFFF")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "label not found")

		// User 1 updates their own label
		updatedName := "Successfully Updated Label"
		updatedColor := "#654321"
		updatedLabel, err := labelUseCase.UpdateLabel(ctx, label.ID, user1.ID, updatedName, updatedColor)
		require.NoError(t, err)
		require.NotNil(t, updatedLabel)
		assert.Equal(t, updatedName, updatedLabel.Name)
		assert.Equal(t, updatedColor, updatedLabel.Color)
		assert.True(t, updatedLabel.UpdatedAt.After(label.UpdatedAt))

		// Create another label to test duplicate name on update
		_, err = labelUseCase.CreateLabel(ctx, user1.ID, "Existing Name", "#000")
		require.NoError(t, err)

		// Try to update the first label to have the same name as the other label
		_, err = labelUseCase.UpdateLabel(ctx, label.ID, user1.ID, "Existing Name", "#FFF")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})

	t.Run("DeleteLabel", func(t *testing.T) {
		// User 1 creates a label
		label, err := labelUseCase.CreateLabel(ctx, user1.ID, "Label To Delete", "#DELETE")
		require.NoError(t, err)
		require.NotNil(t, label)

		// User 2 tries to delete User 1's label
		err = labelUseCase.DeleteLabel(ctx, label.ID, user2.ID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "label not found")

		// Verify the label still exists
		_, err = labelUseCase.GetLabelByID(ctx, label.ID, user1.ID)
		require.NoError(t, err)

		// User 1 deletes their own label
		err = labelUseCase.DeleteLabel(ctx, label.ID, user1.ID)
		require.NoError(t, err)

		// Verify the label is deleted
		_, err = labelUseCase.GetLabelByID(ctx, label.ID, user1.ID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "label not found")
	})

	t.Run("NoteLabelAssociation", func(t *testing.T) {
		// User 1 creates a note and two labels
		note, err := noteUseCase.CreateNote(ctx, user1.ID, "Note for Association", "Content", "note-label")
		require.NoError(t, err)

		labelA, err := labelUseCase.CreateLabel(ctx, user1.ID, "Label A", "#AAA")
		require.NoError(t, err)
		labelB, err := labelUseCase.CreateLabel(ctx, user1.ID, "Label B", "#BBB")
		require.NoError(t, err)

		// User 2 creates a label (should not be associable by user 1)
		labelC_User2, err := labelUseCase.CreateLabel(ctx, user2.ID, "User 2 Label", "#CCC")
		require.NoError(t, err)

		// 1. Add labels to note (User 1)
		err = labelUseCase.AddLabelToNote(ctx, note.ID, labelA.ID, user1.ID)
		require.NoError(t, err)
		err = labelUseCase.AddLabelToNote(ctx, note.ID, labelB.ID, user1.ID)
		require.NoError(t, err)

		// 2. Try invalid associations
		// User 2 tries to add their label to User 1's note
		err = labelUseCase.AddLabelToNote(ctx, note.ID, labelC_User2.ID, user2.ID)
		assert.Error(t, err) // Fails because note doesn't belong to user 2
		assert.Contains(t, err.Error(), "note not found")

		// User 1 tries to add User 2's label to their note
		err = labelUseCase.AddLabelToNote(ctx, note.ID, labelC_User2.ID, user1.ID)
		assert.Error(t, err) // Fails because label doesn't belong to user 1
		assert.Contains(t, err.Error(), "label not found")

		// 3. Get labels for User 1's note
		labelsForNote, err := labelUseCase.GetLabelsForNote(ctx, note.ID, user1.ID)
		require.NoError(t, err)
		assert.Len(t, labelsForNote, 2)
		foundLabels := map[string]bool{labelA.ID: false, labelB.ID: false}
		for _, l := range labelsForNote {
			if _, ok := foundLabels[l.ID]; ok {
				foundLabels[l.ID] = true
			}
		}
		assert.True(t, foundLabels[labelA.ID], "Label A not found for note")
		assert.True(t, foundLabels[labelB.ID], "Label B not found for note")

		// 4. Get notes for Label A (should include User 1's note)
		notesForLabelA, err := labelUseCase.GetNotesForLabel(ctx, labelA.ID, user1.ID)
		require.NoError(t, err)
		assert.Len(t, notesForLabelA, 1)
		assert.Equal(t, note.ID, notesForLabelA[0].ID)

		// User 2 tries to get notes for User 1's label A
		_, err = labelUseCase.GetNotesForLabel(ctx, labelA.ID, user2.ID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "label not found")

		// 5. Remove Label A from the note (User 1)
		err = labelUseCase.RemoveLabelFromNote(ctx, note.ID, labelA.ID, user1.ID)
		require.NoError(t, err)

		// 6. Verify removal
		labelsForNoteAfterRemove, err := labelUseCase.GetLabelsForNote(ctx, note.ID, user1.ID)
		require.NoError(t, err)
		assert.Len(t, labelsForNoteAfterRemove, 1)
		assert.Equal(t, labelB.ID, labelsForNoteAfterRemove[0].ID)

		notesForLabelAAfterRemove, err := labelUseCase.GetNotesForLabel(ctx, labelA.ID, user1.ID)
		require.NoError(t, err)
		assert.Len(t, notesForLabelAAfterRemove, 0)
	})
}
