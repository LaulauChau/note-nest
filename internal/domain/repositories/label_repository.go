package repositories

import (
	"context"

	"github.com/LaulauChau/note-nest/internal/domain/entities"
)

type LabelRepository interface {
	Create(ctx context.Context, label *entities.Label) error

	GetByID(ctx context.Context, id string) (*entities.Label, error)
	GetByUserID(ctx context.Context, userID string) ([]*entities.Label, error)
	GetByName(ctx context.Context, userID, name string) (*entities.Label, error)

	Update(ctx context.Context, label *entities.Label) error

	Delete(ctx context.Context, id string) error

	// Note-Label relationship methods
	AddLabelToNote(ctx context.Context, noteID, labelID string) error

	RemoveLabelFromNote(ctx context.Context, noteID, labelID string) error

	GetLabelsForNote(ctx context.Context, noteID string) ([]*entities.Label, error)
	GetNotesForLabel(ctx context.Context, labelID string) ([]string, error) // Returns note IDs
}
