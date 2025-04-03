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
	noteRepo repositories.NoteRepository
	userRepo repositories.UserRepository
}

func NewNoteUseCase(
	noteRepo repositories.NoteRepository,
	userRepo repositories.UserRepository,
) *NoteUseCase {
	return &NoteUseCase{
		noteRepo: noteRepo,
		userRepo: userRepo,
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
