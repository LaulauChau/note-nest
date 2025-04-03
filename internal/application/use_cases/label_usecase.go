package use_cases

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/LaulauChau/note-nest/internal/domain/entities"
	"github.com/LaulauChau/note-nest/internal/domain/repositories"
)

type LabelUseCase struct {
	labelRepo repositories.LabelRepository
	userRepo  repositories.UserRepository
	noteRepo  repositories.NoteRepository
}

func NewLabelUseCase(
	labelRepo repositories.LabelRepository,
	userRepo repositories.UserRepository,
	noteRepo repositories.NoteRepository,
) *LabelUseCase {
	return &LabelUseCase{
		labelRepo: labelRepo,
		userRepo:  userRepo,
		noteRepo:  noteRepo,
	}
}

func (uc *LabelUseCase) CreateLabel(ctx context.Context, userID, name, color string) (*entities.Label, error) {
	// Verify the user exists
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	// Check if label with same name already exists for this user
	existingLabel, err := uc.labelRepo.GetByName(ctx, userID, name)
	if err != nil {
		return nil, err
	}
	if existingLabel != nil {
		return nil, errors.New("label with this name already exists")
	}

	// Create a new label
	now := time.Now()
	label := &entities.Label{
		ID:        uuid.New().String(),
		UserID:    userID,
		Name:      name,
		Color:     color,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Save the label
	if err := uc.labelRepo.Create(ctx, label); err != nil {
		return nil, err
	}

	return label, nil
}

func (uc *LabelUseCase) GetLabelByID(ctx context.Context, labelID, userID string) (*entities.Label, error) {
	// Get the label
	label, err := uc.labelRepo.GetByID(ctx, labelID)
	if err != nil {
		return nil, err
	}

	// If label not found or doesn't belong to the user, return nil
	if label == nil || label.UserID != userID {
		return nil, errors.New("label not found")
	}

	return label, nil
}

func (uc *LabelUseCase) GetLabelsByUser(ctx context.Context, userID string) ([]*entities.Label, error) {
	// Verify the user exists
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	// Get labels for the user
	return uc.labelRepo.GetByUserID(ctx, userID)
}

func (uc *LabelUseCase) GetLabelsForNote(ctx context.Context, noteID, userID string) ([]*entities.Label, error) {
	// Verify the note exists and belongs to the user
	note, err := uc.noteRepo.GetByID(ctx, noteID)
	if err != nil {
		return nil, err
	}
	if note == nil || note.UserID != userID {
		return nil, errors.New("note not found")
	}

	// Get labels for the note
	return uc.labelRepo.GetLabelsForNote(ctx, noteID)
}

func (uc *LabelUseCase) GetNotesForLabel(ctx context.Context, labelID, userID string) ([]*entities.Note, error) {
	// Verify the label exists and belongs to the user
	label, err := uc.GetLabelByID(ctx, labelID, userID)
	if err != nil {
		return nil, err
	}
	if label == nil {
		return nil, errors.New("label not found")
	}

	// Get note IDs for the label
	noteIDs, err := uc.labelRepo.GetNotesForLabel(ctx, labelID)
	if err != nil {
		return nil, err
	}

	// Get the notes
	notes := make([]*entities.Note, 0, len(noteIDs))
	for _, noteID := range noteIDs {
		note, err := uc.noteRepo.GetByID(ctx, noteID)
		if err != nil {
			continue // Skip notes with errors
		}

		// Only include notes that belong to the user
		if note != nil && note.UserID == userID {
			notes = append(notes, note)
		}
	}

	return notes, nil
}

func (uc *LabelUseCase) UpdateLabel(ctx context.Context, labelID, userID, name, color string) (*entities.Label, error) {
	// Get the label
	label, err := uc.labelRepo.GetByID(ctx, labelID)
	if err != nil {
		return nil, err
	}

	// If label not found or doesn't belong to the user, return error
	if label == nil || label.UserID != userID {
		return nil, errors.New("label not found")
	}

	// Check if another label with the same name already exists for this user
	if name != label.Name {
		existingLabel, err := uc.labelRepo.GetByName(ctx, userID, name)
		if err != nil {
			return nil, err
		}
		if existingLabel != nil && existingLabel.ID != labelID {
			return nil, errors.New("label with this name already exists")
		}
	}

	// Update the label fields
	label.Name = name
	label.Color = color
	label.UpdatedAt = time.Now()

	// Save the updated label
	if err := uc.labelRepo.Update(ctx, label); err != nil {
		return nil, err
	}

	return label, nil
}

func (uc *LabelUseCase) DeleteLabel(ctx context.Context, labelID, userID string) error {
	// Get the label
	label, err := uc.labelRepo.GetByID(ctx, labelID)
	if err != nil {
		return err
	}

	// If label not found or doesn't belong to the user, return error
	if label == nil || label.UserID != userID {
		return errors.New("label not found")
	}

	// Delete the label
	return uc.labelRepo.Delete(ctx, labelID)
}

func (uc *LabelUseCase) AddLabelToNote(ctx context.Context, noteID, labelID, userID string) error {
	// Verify the note exists and belongs to the user
	note, err := uc.noteRepo.GetByID(ctx, noteID)
	if err != nil {
		return err
	}
	if note == nil || note.UserID != userID {
		return errors.New("note not found")
	}

	// Verify the label exists and belongs to the user
	label, err := uc.labelRepo.GetByID(ctx, labelID)
	if err != nil {
		return err
	}
	if label == nil || label.UserID != userID {
		return errors.New("label not found")
	}

	// Associate the label with the note
	return uc.labelRepo.AddLabelToNote(ctx, noteID, labelID)
}

func (uc *LabelUseCase) RemoveLabelFromNote(ctx context.Context, noteID, labelID, userID string) error {
	// Verify the note exists and belongs to the user
	note, err := uc.noteRepo.GetByID(ctx, noteID)
	if err != nil {
		return err
	}
	if note == nil || note.UserID != userID {
		return errors.New("note not found")
	}

	// Verify the label exists and belongs to the user
	label, err := uc.labelRepo.GetByID(ctx, labelID)
	if err != nil {
		return err
	}
	if label == nil || label.UserID != userID {
		return errors.New("label not found")
	}

	// Disassociate the label from the note
	return uc.labelRepo.RemoveLabelFromNote(ctx, noteID, labelID)
}
