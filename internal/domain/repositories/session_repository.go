package repositories

import (
	"context"
	"time"

	"github.com/LaulauChau/note-nest/internal/domain/entities"
)

type SessionRepository interface {
	Create(ctx context.Context, session *entities.Session) error

	GetByID(ctx context.Context, id string) (*entities.Session, error)
	GetSessionWithUser(ctx context.Context, sessionID string) (*entities.SessionValidationResult, error)

	UpdateExpiresAt(ctx context.Context, sessionID string, expiresAt time.Time) error

	Delete(ctx context.Context, sessionID string) error
	DeleteAllByUserID(ctx context.Context, userID string) error
}
