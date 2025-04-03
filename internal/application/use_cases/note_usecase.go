package use_cases

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/LaulauChau/note-nest/internal/domain/entities"
	"github.com/LaulauChau/note-nest/internal/domain/repositories"
)

type NoteUseCase struct {
	noteRepo  repositories.NoteRepository
	userRepo  repositories.UserRepository
	labelRepo repositories.LabelRepository
}

func NewNoteUseCase(
	noteRepo repositories.NoteRepository,
	userRepo repositories.UserRepository,
	labelRepo repositories.LabelRepository,
) *NoteUseCase {
	return &NoteUseCase{
		noteRepo:  noteRepo,
		userRepo:  userRepo,
		labelRepo: labelRepo,
	}
}

func (uc *NoteUseCase) CreateNote(ctx context.Context, userID, title, content, label string) (*entities.Note, error) {
	// Verify the user exists
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	// Create a new note
	now := time.Now()
	note := &entities.Note{
		ID:         uuid.New().String(),
		UserID:     userID,
		Title:      title,
		Content:    content,
		IsArchived: false,
		Label:      label,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	// Save the note
	if err := uc.noteRepo.Create(ctx, note); err != nil {
		return nil, err
	}

	return note, nil
}

func (uc *NoteUseCase) GetNoteByID(ctx context.Context, noteID, userID string) (*entities.Note, error) {
	// Get the note
	note, err := uc.noteRepo.GetByID(ctx, noteID)
	if err != nil {
		return nil, err
	}

	// If note not found or doesn't belong to the user, return nil
	if note == nil || note.UserID != userID {
		return nil, errors.New("note not found")
	}

	return note, nil
}

func (uc *NoteUseCase) GetActiveNotes(ctx context.Context, userID string) ([]*entities.Note, error) {
	// Verify the user exists
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	// Get active notes for the user
	return uc.noteRepo.GetByUserID(ctx, userID)
}

func (uc *NoteUseCase) GetArchivedNotes(ctx context.Context, userID string) ([]*entities.Note, error) {
	// Verify the user exists
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	// Get archived notes for the user
	return uc.noteRepo.GetArchivedByUserID(ctx, userID)
}

func (uc *NoteUseCase) UpdateNote(ctx context.Context, noteID, userID, title, content, label string, isArchived bool) (*entities.Note, error) {
	// Get the note
	note, err := uc.noteRepo.GetByID(ctx, noteID)
	if err != nil {
		return nil, err
	}

	// If note not found or doesn't belong to the user, return error
	if note == nil || note.UserID != userID {
		return nil, errors.New("note not found")
	}

	// Update the note fields
	note.Title = title
	note.Content = content
	note.Label = label
	note.IsArchived = isArchived
	note.UpdatedAt = time.Now() // Make sure this line is present

	// Save the updated note
	if err := uc.noteRepo.Update(ctx, note); err != nil {
		return nil, err
	}

	return note, nil
}

func (uc *NoteUseCase) DeleteNote(ctx context.Context, noteID, userID string) error {
	// Get the note
	note, err := uc.noteRepo.GetByID(ctx, noteID)
	if err != nil {
		return err
	}

	// If note not found or doesn't belong to the user, return error
	if note == nil || note.UserID != userID {
		return errors.New("note not found")
	}

	// Delete the note
	return uc.noteRepo.Delete(ctx, noteID)
}

func (uc *NoteUseCase) GetNoteWithLabels(ctx context.Context, noteID, userID string) (*entities.Note, []*entities.Label, error) {
	// Get the note
	note, err := uc.GetNoteByID(ctx, noteID, userID)
	if err != nil {
		return nil, nil, err
	}

	// Get labels for the note
	labels, err := uc.labelRepo.GetLabelsForNote(ctx, noteID)
	if err != nil {
		return note, nil, err
	}

	return note, labels, nil
}

func (uc *NoteUseCase) CreateNoteWithLabels(ctx context.Context, userID, title, content, label string, labelIDs []string) (*entities.Note, error) {
	// Create the note
	note, err := uc.CreateNote(ctx, userID, title, content, label)
	if err != nil {
		return nil, err
	}

	// Associate labels with the note
	for _, labelID := range labelIDs {
		// Verify label exists and belongs to the user
		label, err := uc.labelRepo.GetByID(ctx, labelID)
		if err != nil || label == nil || label.UserID != userID {
			continue // Skip invalid labels
		}

		// Associate label with note
		err = uc.labelRepo.AddLabelToNote(ctx, note.ID, labelID)
		if err != nil {
			return nil, err
		}
	}

	return note, nil
}

func (uc *NoteUseCase) UpdateNoteWithLabels(ctx context.Context, noteID, userID, title, content, label string, isArchived bool, labelIDs []string) (*entities.Note, error) {
	// Update the note
	note, err := uc.UpdateNote(ctx, noteID, userID, title, content, label, isArchived)
	if err != nil {
		return nil, err
	}

	// Get current labels for the note
	currentLabels, err := uc.labelRepo.GetLabelsForNote(ctx, noteID)
	if err != nil {
		return note, err
	}

	// Create a map of current label IDs for easy lookup
	currentLabelMap := make(map[string]bool)
	for _, label := range currentLabels {
		currentLabelMap[label.ID] = true
	}

	// Create a map of new label IDs for easy lookup
	newLabelMap := make(map[string]bool)
	for _, labelID := range labelIDs {
		newLabelMap[labelID] = true
	}

	// Add new labels that weren't previously associated
	for _, labelID := range labelIDs {
		if !currentLabelMap[labelID] {
			// Verify label exists and belongs to the user
			label, err := uc.labelRepo.GetByID(ctx, labelID)
			if err != nil || label == nil || label.UserID != userID {
				continue // Skip invalid labels
			}

			// Associate label with note
			err = uc.labelRepo.AddLabelToNote(ctx, noteID, labelID)
			if err != nil {
				return nil, err
			}
		}
	}

	// Remove labels that are no longer associated
	for _, label := range currentLabels {
		if !newLabelMap[label.ID] {
			err = uc.labelRepo.RemoveLabelFromNote(ctx, noteID, label.ID)
			if err != nil {
				return nil, err
			}
		}
	}

	return note, nil
}
