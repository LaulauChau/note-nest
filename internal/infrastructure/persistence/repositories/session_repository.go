package repositories

import (
	"context"
	"time"

	"github.com/LaulauChau/note-nest/internal/domain/entities"
	"github.com/LaulauChau/note-nest/internal/domain/repositories"
	"github.com/google/uuid"
)

type SessionRepositoryImpl struct {
	q *Queries
}

func NewSessionRepository(q *Queries) repositories.SessionRepository {
	return &SessionRepositoryImpl{q: q}
}

func (r *SessionRepositoryImpl) Create(ctx context.Context, session *entities.Session) error {
	// Parse the user ID (which should be a UUID)
	userID, err := uuid.Parse(session.UserID)
	if err != nil {
		return err
	}

	// Use manual query - session ID is a string (hash) not a UUID
	_, err = r.q.db.Exec(ctx,
		"INSERT INTO sessions (id, user_id, expires_at, created_at) VALUES ($1, $2, $3, $4)",
		session.ID, userID, session.ExpiresAt, session.CreatedAt)

	return err
}

func (r *SessionRepositoryImpl) GetByID(ctx context.Context, id string) (*entities.Session, error) {
	// Use manual query - the sessionID is a string (hash) not a UUID
	var session struct {
		ID        string    `db:"id"`
		UserID    uuid.UUID `db:"user_id"`
		ExpiresAt time.Time `db:"expires_at"`
		CreatedAt time.Time `db:"created_at"`
	}

	err := r.q.db.QueryRow(ctx, "SELECT * FROM sessions WHERE id = $1", id).Scan(
		&session.ID,
		&session.UserID,
		&session.ExpiresAt,
		&session.CreatedAt,
	)

	if err != nil {
		// Handle "no rows in result set" error gracefully
		if err.Error() == "no rows in result set" {
			return nil, nil
		}
		return nil, err
	}

	return &entities.Session{
		ID:        session.ID,
		UserID:    session.UserID.String(),
		ExpiresAt: session.ExpiresAt,
		CreatedAt: session.CreatedAt,
	}, nil
}

func (r *SessionRepositoryImpl) GetSessionWithUser(ctx context.Context, sessionID string) (*entities.SessionValidationResult, error) {
	// Use manual query - the sessionID is a string (hash) not a UUID
	var result struct {
		SessionID        string    `db:"session_id"`
		SessionUserID    uuid.UUID `db:"session_user_id"`
		SessionExpiresAt time.Time `db:"session_expires_at"`
		SessionCreatedAt time.Time `db:"session_created_at"`
		UserID           uuid.UUID `db:"user_id"`
		UserEmail        string    `db:"user_email"`
		UserName         string    `db:"user_name"`
	}

	err := r.q.db.QueryRow(ctx,
		`SELECT
			sessions.id AS session_id,
			sessions.user_id AS session_user_id,
			sessions.expires_at AS session_expires_at,
			sessions.created_at AS session_created_at,
			users.id AS user_id,
			users.email AS user_email,
			users.name AS user_name
		FROM sessions
		INNER JOIN users ON sessions.user_id = users.id
		WHERE sessions.id = $1`, sessionID).Scan(
		&result.SessionID,
		&result.SessionUserID,
		&result.SessionExpiresAt,
		&result.SessionCreatedAt,
		&result.UserID,
		&result.UserEmail,
		&result.UserName,
	)

	if err != nil {
		// Handle "no rows in result set" error gracefully
		if err.Error() == "no rows in result set" {
			return &entities.SessionValidationResult{
				Session: nil,
				User:    nil,
			}, nil
		}
		return nil, err
	}

	return &entities.SessionValidationResult{
		Session: &entities.Session{
			ID:        result.SessionID,
			UserID:    result.SessionUserID.String(),
			ExpiresAt: result.SessionExpiresAt,
			CreatedAt: result.SessionCreatedAt,
		},
		User: &entities.User{
			ID:    result.UserID.String(),
			Email: result.UserEmail,
			Name:  result.UserName,
		},
	}, nil
}

func (r *SessionRepositoryImpl) UpdateExpiresAt(ctx context.Context, sessionID string, expiresAt time.Time) error {
	// Use manual query - the sessionID is a string (hash) not a UUID
	_, err := r.q.db.Exec(ctx, "UPDATE sessions SET expires_at = $2 WHERE id = $1", sessionID, expiresAt)
	return err
}

func (r *SessionRepositoryImpl) Delete(ctx context.Context, sessionID string) error {
	// Use manual query - the sessionID is a string (hash) not a UUID
	_, err := r.q.db.Exec(ctx, "DELETE FROM sessions WHERE id = $1", sessionID)
	return err
}

func (r *SessionRepositoryImpl) DeleteAllByUserID(ctx context.Context, userID string) error {
	id, err := uuid.Parse(userID)
	if err != nil {
		return err
	}

	return r.q.DeleteAllSessionsByUserID(ctx, id.String())
}
