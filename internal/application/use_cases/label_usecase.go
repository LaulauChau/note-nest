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
