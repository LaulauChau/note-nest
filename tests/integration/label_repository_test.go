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

func TestLabelRepository(t *testing.T) {
	// Set up test database
	ctx := context.Background()
	db, err := SetupTestDatabase(ctx)
	require.NoError(t, err)
	defer db.Close(ctx)

	// Create repositories
	queries := repositories.New(db.Pool)
	userRepo := repositories.NewUserRepository(queries)
	labelRepo := repositories.NewLabelRepository(queries)
	noteRepo := repositories.NewNoteRepository(queries) // Needed for note-label tests

	// Create a test user first
	now := time.Now()
	testUser := &entities.User{
		ID:        uuid.New().String(),
		Email:     "labeltestuser@example.com",
		Name:      "Label Test User",
		Password:  "hashedpassword",
		CreatedAt: now,
		UpdatedAt: now,
	}
	err = userRepo.Create(ctx, testUser)
	require.NoError(t, err)

	t.Run("CreateAndGetLabelByID", func(t *testing.T) {
		labelID := uuid.New().String()
		label := &entities.Label{
			ID:        labelID,
			UserID:    testUser.ID,
			Name:      "Test Label 1",
			Color:     "#FF0000",
			CreatedAt: now,
			UpdatedAt: now,
		}

		err := labelRepo.Create(ctx, label)
		require.NoError(t, err)

		retrievedLabel, err := labelRepo.GetByID(ctx, labelID)
		require.NoError(t, err)
		require.NotNil(t, retrievedLabel)

		assert.Equal(t, labelID, retrievedLabel.ID)
		assert.Equal(t, testUser.ID, retrievedLabel.UserID)
		assert.Equal(t, "Test Label 1", retrievedLabel.Name)
		assert.Equal(t, "#FF0000", retrievedLabel.Color)
	})

	t.Run("GetLabelByUserID", func(t *testing.T) {
		// Create a separate user for isolation
		userForThisTest := &entities.User{
			ID:        uuid.New().String(),
			Email:     "labelsbyuser@example.com",
			Name:      "Labels By User Test",
			Password:  "hashed",
			CreatedAt: now,
			UpdatedAt: now,
		}
		require.NoError(t, userRepo.Create(ctx, userForThisTest))

		// Create labels for this user
		label1 := &entities.Label{ID: uuid.New().String(), UserID: userForThisTest.ID, Name: "UserLabel 1", Color: "#111", CreatedAt: now, UpdatedAt: now}
		label2 := &entities.Label{ID: uuid.New().String(), UserID: userForThisTest.ID, Name: "UserLabel 2", Color: "#222", CreatedAt: now, UpdatedAt: now}
		// Create label for another user (should not be returned)
		labelOther := &entities.Label{ID: uuid.New().String(), UserID: testUser.ID, Name: "OtherUser Label", Color: "#000", CreatedAt: now, UpdatedAt: now}

		require.NoError(t, labelRepo.Create(ctx, label1))
		require.NoError(t, labelRepo.Create(ctx, label2))
		require.NoError(t, labelRepo.Create(ctx, labelOther))

		retrievedLabels, err := labelRepo.GetByUserID(ctx, userForThisTest.ID)
		require.NoError(t, err)
		assert.Len(t, retrievedLabels, 2)

		// Verify correct labels are returned
		foundIDs := map[string]bool{label1.ID: false, label2.ID: false}
		for _, l := range retrievedLabels {
			assert.Equal(t, userForThisTest.ID, l.UserID)
			if _, ok := foundIDs[l.ID]; ok {
				foundIDs[l.ID] = true
			}
		}
		assert.True(t, foundIDs[label1.ID], "Label 1 not found")
		assert.True(t, foundIDs[label2.ID], "Label 2 not found")
	})

	t.Run("GetLabelByName", func(t *testing.T) {
		labelName := "Unique Label Name"
		label := &entities.Label{
			ID:        uuid.New().String(),
			UserID:    testUser.ID,
			Name:      labelName,
			Color:     "#334455",
			CreatedAt: now,
			UpdatedAt: now,
		}
		require.NoError(t, labelRepo.Create(ctx, label))

		retrievedLabel, err := labelRepo.GetByName(ctx, testUser.ID, labelName)
		require.NoError(t, err)
		require.NotNil(t, retrievedLabel)
		assert.Equal(t, label.ID, retrievedLabel.ID)
		assert.Equal(t, labelName, retrievedLabel.Name)

		// Test not found
		retrievedLabelNotFound, err := labelRepo.GetByName(ctx, testUser.ID, "NonExistentName")
		require.NoError(t, err)
		assert.Nil(t, retrievedLabelNotFound)

		// Test wrong user
		retrievedLabelWrongUser, err := labelRepo.GetByName(ctx, uuid.New().String(), labelName)
		require.NoError(t, err)
		assert.Nil(t, retrievedLabelWrongUser)
	})

	t.Run("UpdateLabel", func(t *testing.T) {
		labelID := uuid.New().String()
		originalLabel := &entities.Label{
			ID:        labelID,
			UserID:    testUser.ID,
			Name:      "To Be Updated",
			Color:     "#ABCDEF",
			CreatedAt: now,
			UpdatedAt: now,
		}
		require.NoError(t, labelRepo.Create(ctx, originalLabel))

		updatedLabel := &entities.Label{
			ID:        labelID,
			UserID:    testUser.ID, // UserID should not change
			Name:      "Updated Name",
			Color:     "#FEDCBA",
			UpdatedAt: time.Now(), // This will be set by the repo method
		}

		time.Sleep(1 * time.Millisecond) // Ensure UpdatedAt changes measurably
		err := labelRepo.Update(ctx, updatedLabel)
		require.NoError(t, err)

		retrievedLabel, err := labelRepo.GetByID(ctx, labelID)
		require.NoError(t, err)
		require.NotNil(t, retrievedLabel)

		assert.Equal(t, "Updated Name", retrievedLabel.Name)
		assert.Equal(t, "#FEDCBA", retrievedLabel.Color)
		assert.Equal(t, originalLabel.UserID, retrievedLabel.UserID) // Check UserID unchanged
		assert.True(t, retrievedLabel.UpdatedAt.After(originalLabel.UpdatedAt))
	})

	t.Run("DeleteLabel", func(t *testing.T) {
		labelID := uuid.New().String()
		label := &entities.Label{
			ID:        labelID,
			UserID:    testUser.ID,
			Name:      "To Be Deleted",
			Color:     "#000000",
			CreatedAt: now,
			UpdatedAt: now,
		}
		require.NoError(t, labelRepo.Create(ctx, label))

		// Verify it exists
		_, err := labelRepo.GetByID(ctx, labelID)
		require.NoError(t, err)

		// Delete it
		err = labelRepo.Delete(ctx, labelID)
		require.NoError(t, err)

		// Verify it's gone
		retrievedLabel, err := labelRepo.GetByID(ctx, labelID)
		require.NoError(t, err)
		assert.Nil(t, retrievedLabel)
	})

	t.Run("NoteLabelOperations", func(t *testing.T) {
		// Create a user, a note, and two labels for this test
		opUserID := uuid.New().String()
		opUser := &entities.User{ID: opUserID, Email: "opuser@example.com", Name: "Op User", Password: "p", CreatedAt: now, UpdatedAt: now}
		require.NoError(t, userRepo.Create(ctx, opUser))

		opNoteID := uuid.New().String()
		opNote := &entities.Note{ID: opNoteID, UserID: opUserID, Title: "Note For Labels", Content: "c", CreatedAt: now, UpdatedAt: now}
		require.NoError(t, noteRepo.Create(ctx, opNote))

		opLabel1ID := uuid.New().String()
		opLabel1 := &entities.Label{ID: opLabel1ID, UserID: opUserID, Name: "OpLabel1", Color: "#op1", CreatedAt: now, UpdatedAt: now}
		require.NoError(t, labelRepo.Create(ctx, opLabel1))

		opLabel2ID := uuid.New().String()
		opLabel2 := &entities.Label{ID: opLabel2ID, UserID: opUserID, Name: "OpLabel2", Color: "#op2", CreatedAt: now, UpdatedAt: now}
		require.NoError(t, labelRepo.Create(ctx, opLabel2))

		// 1. Add labels to note
		err := labelRepo.AddLabelToNote(ctx, opNoteID, opLabel1ID)
		require.NoError(t, err)
		err = labelRepo.AddLabelToNote(ctx, opNoteID, opLabel2ID)
		require.NoError(t, err)

		// 2. Get labels for the note
		labelsForNote, err := labelRepo.GetLabelsForNote(ctx, opNoteID)
		require.NoError(t, err)
		assert.Len(t, labelsForNote, 2)
		foundLabels := map[string]bool{opLabel1ID: false, opLabel2ID: false}
		for _, l := range labelsForNote {
			if _, ok := foundLabels[l.ID]; ok {
				foundLabels[l.ID] = true
			}
		}
		assert.True(t, foundLabels[opLabel1ID], "Label 1 not found for note")
		assert.True(t, foundLabels[opLabel2ID], "Label 2 not found for note")

		// 3. Get notes for a label
		notesForLabel1, err := labelRepo.GetNotesForLabel(ctx, opLabel1ID)
		require.NoError(t, err)
		assert.Len(t, notesForLabel1, 1)
		assert.Equal(t, opNoteID, notesForLabel1[0])

		// 4. Remove one label from the note
		err = labelRepo.RemoveLabelFromNote(ctx, opNoteID, opLabel1ID)
		require.NoError(t, err)

		// 5. Verify removal
		labelsForNoteAfterRemove, err := labelRepo.GetLabelsForNote(ctx, opNoteID)
		require.NoError(t, err)
		assert.Len(t, labelsForNoteAfterRemove, 1)
		assert.Equal(t, opLabel2ID, labelsForNoteAfterRemove[0].ID)

		notesForLabel1AfterRemove, err := labelRepo.GetNotesForLabel(ctx, opLabel1ID)
		require.NoError(t, err)
		assert.Len(t, notesForLabel1AfterRemove, 0)
	})
}
