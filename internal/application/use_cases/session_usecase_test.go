package use_cases_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/LaulauChau/note-nest/internal/application/use_cases"
	"github.com/LaulauChau/note-nest/internal/domain/entities"
)

// MockSessionRepository mocks the SessionRepository interface
type MockSessionRepository struct {
	mock.Mock
}

func (m *MockSessionRepository) Create(ctx context.Context, session *entities.Session) error {
	args := m.Called(ctx, session)
	return args.Error(0)
}

func (m *MockSessionRepository) GetByID(ctx context.Context, id string) (*entities.Session, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Session), args.Error(1)
}

func (m *MockSessionRepository) GetSessionWithUser(ctx context.Context, sessionID string) (*entities.SessionValidationResult, error) {
	args := m.Called(ctx, sessionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.SessionValidationResult), args.Error(1)
}

func (m *MockSessionRepository) UpdateExpiresAt(ctx context.Context, sessionID string, expiresAt time.Time) error {
	args := m.Called(ctx, sessionID, expiresAt)
	return args.Error(0)
}

func (m *MockSessionRepository) Delete(ctx context.Context, sessionID string) error {
	args := m.Called(ctx, sessionID)
	return args.Error(0)
}

func (m *MockSessionRepository) DeleteAllByUserID(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

// MockUserRepository mocks the UserRepository interface
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *entities.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id string) (*entities.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *entities.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockTokenService mocks the TokenService interface
type MockTokenService struct {
	mock.Mock
}

func (m *MockTokenService) GenerateToken(ctx context.Context) (string, error) {
	args := m.Called(ctx)
	return args.String(0), args.Error(1)
}

func (m *MockTokenService) HashToken(ctx context.Context, token string) (string, error) {
	args := m.Called(ctx, token)
	return args.String(0), args.Error(1)
}

func TestCreateSession(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockSessionRepo := new(MockSessionRepository)
	mockUserRepo := new(MockUserRepository)
	mockTokenService := new(MockTokenService)

	userID := uuid.New().String()
	token := "random-token-string"
	hashedToken := "hashed-token-string"

	mockTokenService.On("HashToken", ctx, token).Return(hashedToken, nil)
	// Use mock.MatchedBy to match any session argument
	mockSessionRepo.On("Create", ctx, mock.AnythingOfType("*entities.Session")).Return(nil)

	useCase := use_cases.NewSessionUseCase(mockSessionRepo, mockUserRepo, mockTokenService)

	// Act
	result, err := useCase.CreateSession(ctx, token, userID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)

	mockTokenService.AssertExpectations(t)
	mockSessionRepo.AssertExpectations(t)
}

func TestValidateSessionToken_ValidSession(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockSessionRepo := new(MockSessionRepository)
	mockUserRepo := new(MockUserRepository)
	mockTokenService := new(MockTokenService)

	token := "random-token-string"
	hashedToken := "hashed-token-string"

	userID := uuid.New().String()
	sessionID := uuid.New().String()

	session := &entities.Session{
		ID:        sessionID,
		UserID:    userID,
		ExpiresAt: time.Now().Add(20 * 24 * time.Hour),
		CreatedAt: time.Now(),
	}

	user := &entities.User{
		ID:    userID,
		Email: "test@example.com",
		Name:  "Test User",
	}

	validationResult := &entities.SessionValidationResult{
		Session: session,
		User:    user,
	}

	mockTokenService.On("HashToken", ctx, token).Return(hashedToken, nil)
	mockSessionRepo.On("GetSessionWithUser", ctx, hashedToken).Return(validationResult, nil)

	useCase := use_cases.NewSessionUseCase(mockSessionRepo, mockUserRepo, mockTokenService)

	// Act
	result, err := useCase.ValidateSessionToken(ctx, token)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, result.Session)
	assert.NotNil(t, result.User)
	assert.Equal(t, sessionID, result.Session.ID)
	assert.Equal(t, userID, result.User.ID)

	mockTokenService.AssertExpectations(t)
	mockSessionRepo.AssertExpectations(t)
}

func TestValidateSessionToken_ExpiredSession(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockSessionRepo := new(MockSessionRepository)
	mockUserRepo := new(MockUserRepository)
	mockTokenService := new(MockTokenService)

	token := "random-token-string"
	hashedToken := "hashed-token-string"

	userID := uuid.New().String()
	sessionID := uuid.New().String()

	// Create an expired session
	session := &entities.Session{
		ID:        sessionID,
		UserID:    userID,
		ExpiresAt: time.Now().Add(-1 * time.Hour), // Expired 1 hour ago
		CreatedAt: time.Now().Add(-30 * 24 * time.Hour),
	}

	user := &entities.User{
		ID:    userID,
		Email: "test@example.com",
		Name:  "Test User",
	}

	validationResult := &entities.SessionValidationResult{
		Session: session,
		User:    user,
	}

	mockTokenService.On("HashToken", ctx, token).Return(hashedToken, nil)
	mockSessionRepo.On("GetSessionWithUser", ctx, hashedToken).Return(validationResult, nil)
	mockSessionRepo.On("Delete", ctx, sessionID).Return(nil)

	useCase := use_cases.NewSessionUseCase(mockSessionRepo, mockUserRepo, mockTokenService)

	// Act
	result, err := useCase.ValidateSessionToken(ctx, token)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Nil(t, result.Session)
	assert.Nil(t, result.User)

	mockTokenService.AssertExpectations(t)
	mockSessionRepo.AssertExpectations(t)
	// Verify that Delete was called for the expired session
	mockSessionRepo.AssertCalled(t, "Delete", ctx, sessionID)
}

func TestInvalidateSession(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockSessionRepo := new(MockSessionRepository)
	mockUserRepo := new(MockUserRepository)
	mockTokenService := new(MockTokenService)

	sessionID := uuid.New().String()

	mockSessionRepo.On("Delete", ctx, sessionID).Return(nil)

	useCase := use_cases.NewSessionUseCase(mockSessionRepo, mockUserRepo, mockTokenService)

	// Act
	err := useCase.InvalidateSession(ctx, sessionID)

	// Assert
	assert.NoError(t, err)
	mockSessionRepo.AssertExpectations(t)
}

func TestInvalidateAllSessions(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockSessionRepo := new(MockSessionRepository)
	mockUserRepo := new(MockUserRepository)
	mockTokenService := new(MockTokenService)

	userID := uuid.New().String()

	mockSessionRepo.On("DeleteAllByUserID", ctx, userID).Return(nil)

	useCase := use_cases.NewSessionUseCase(mockSessionRepo, mockUserRepo, mockTokenService)

	// Act
	err := useCase.InvalidateAllSessions(ctx, userID)

	// Assert
	assert.NoError(t, err)
	mockSessionRepo.AssertExpectations(t)
}
