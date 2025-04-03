package repositories

import (
	"context"

	"github.com/LaulauChau/note-nest/internal/domain/entities"
)

type NoteRepository interface {
	Create(ctx context.Context, note *entities.Note) error

	GetByID(ctx context.Context, id string) (*entities.Note, error)
	GetByUserID(ctx context.Context, userID string) ([]*entities.Note, error)
	GetArchivedByUserID(ctx context.Context, userID string) ([]*entities.Note, error) // Get archived

	Update(ctx context.Context, note *entities.Note) error

	Delete(ctx context.Context, id string) error
}
