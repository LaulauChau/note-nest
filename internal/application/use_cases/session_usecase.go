package use_cases

import (
	"context"
	"time"

	"github.com/LaulauChau/note-nest/internal/application/services"
	"github.com/LaulauChau/note-nest/internal/domain/entities"
	"github.com/LaulauChau/note-nest/internal/domain/repositories"
)

type SessionUseCase struct {
	sessionRepo  repositories.SessionRepository
	userRepo     repositories.UserRepository
	tokenService services.TokenService
}

func NewSessionUseCase(
	sessionRepo repositories.SessionRepository,
	userRepo repositories.UserRepository,
	tokenService services.TokenService,
) *SessionUseCase {
	return &SessionUseCase{
		sessionRepo:  sessionRepo,
		userRepo:     userRepo,
		tokenService: tokenService,
	}
}

func (uc *SessionUseCase) GenerateSessionToken(ctx context.Context) (string, error) {
	return uc.tokenService.GenerateToken(ctx)
}

func (uc *SessionUseCase) CreateSession(ctx context.Context, token string, userID string) (*entities.Session, error) {
	sessionID, err := uc.tokenService.HashToken(ctx, token)
	if err != nil {
		return nil, err
	}

	// 30 days from now
	expiresAt := time.Now().Add(30 * 24 * time.Hour)

	session := &entities.Session{
		ID:        sessionID,
		UserID:    userID,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	}

	if err := uc.sessionRepo.Create(ctx, session); err != nil {
		return nil, err
	}

	return session, nil
}

func (uc *SessionUseCase) ValidateSessionToken(ctx context.Context, token string) (*entities.SessionValidationResult, error) {
	sessionID, err := uc.tokenService.HashToken(ctx, token)
	if err != nil {
		return &entities.SessionValidationResult{Session: nil, User: nil}, nil
	}

	result, err := uc.sessionRepo.GetSessionWithUser(ctx, sessionID)
	if err != nil {
		return &entities.SessionValidationResult{Session: nil, User: nil}, nil
	}

	if result.Session == nil || result.User == nil {
		return &entities.SessionValidationResult{Session: nil, User: nil}, nil
	}

	now := time.Now()

	// Check if session is expired
	if now.After(result.Session.ExpiresAt) {
		if err := uc.sessionRepo.Delete(ctx, result.Session.ID); err != nil {
			return nil, err
		}
		return &entities.SessionValidationResult{Session: nil, User: nil}, nil
	}

	// If session is going to expire in 15 days, extend it
	fifteenDaysFromNow := now.Add(15 * 24 * time.Hour)
	if result.Session.ExpiresAt.Before(fifteenDaysFromNow) {
		// Extend session by 30 days
		newExpiresAt := now.Add(30 * 24 * time.Hour)
		result.Session.ExpiresAt = newExpiresAt

		if err := uc.sessionRepo.UpdateExpiresAt(ctx, result.Session.ID, newExpiresAt); err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (uc *SessionUseCase) InvalidateSession(ctx context.Context, sessionID string) error {
	return uc.sessionRepo.Delete(ctx, sessionID)
}

func (uc *SessionUseCase) InvalidateAllSessions(ctx context.Context, userID string) error {
	return uc.sessionRepo.DeleteAllByUserID(ctx, userID)
}
