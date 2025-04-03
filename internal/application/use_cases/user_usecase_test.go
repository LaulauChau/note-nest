package use_cases_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/LaulauChau/note-nest/internal/application/use_cases"
	"github.com/LaulauChau/note-nest/internal/domain/entities"
)

// MockHashService mocks the HashService interface
type MockHashService struct {
	mock.Mock
}

func (m *MockHashService) HashPassword(ctx context.Context, password string) (string, error) {
	args := m.Called(ctx, password)
	return args.String(0), args.Error(1)
}

func (m *MockHashService) VerifyPassword(ctx context.Context, hashedPassword, password string) (bool, error) {
	args := m.Called(ctx, hashedPassword, password)
	return args.Bool(0), args.Error(1)
}

func TestRegisterUser(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockUserRepo := new(MockUserRepository)
	mockHashService := new(MockHashService)

	email := "test@example.com"
	name := "Test User"
	password := "securepassword"
	hashedPassword := "hashed_password_value"

	// Mock GetByEmail to return nil (user doesn't exist)
	mockUserRepo.On("GetByEmail", ctx, email).Return(nil, nil)

	// Mock HashPassword
	mockHashService.On("HashPassword", ctx, password).Return(hashedPassword, nil)

	// Mock Create to simply return nil (success)
	// Note that RegisterUser creates the user object with ID internally
	mockUserRepo.On("Create", ctx, mock.MatchedBy(func(u *entities.User) bool {
		return u.Email == email && u.Name == name && u.Password == hashedPassword && u.ID != ""
	})).Return(nil)

	useCase := use_cases.NewUserUseCase(mockUserRepo, mockHashService)

	// Act
	user, err := useCase.RegisterUser(ctx, email, name, password)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, email, user.Email)
	assert.Equal(t, name, user.Name)
	assert.NotEmpty(t, user.ID)
	assert.Empty(t, user.Password) // Password should not be returned

	mockUserRepo.AssertExpectations(t)
	mockHashService.AssertExpectations(t)
}

func TestRegisterUser_EmailAlreadyExists(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockUserRepo := new(MockUserRepository)
	mockHashService := new(MockHashService)

	email := "existing@example.com"
	name := "Test User"
	password := "securepassword"

	// Mock GetByEmail to return an existing user
	existingUser := &entities.User{
		ID:    uuid.New().String(),
		Email: email,
		Name:  "Existing User",
	}
	mockUserRepo.On("GetByEmail", ctx, email).Return(existingUser, nil)

	useCase := use_cases.NewUserUseCase(mockUserRepo, mockHashService)

	// Act
	user, err := useCase.RegisterUser(ctx, email, name, password)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "email already taken", err.Error())

	mockUserRepo.AssertExpectations(t)
	// HashPassword should not be called if email already exists
	mockHashService.AssertNotCalled(t, "HashPassword")
}

func TestAuthenticateUser_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockUserRepo := new(MockUserRepository)
	mockHashService := new(MockHashService)

	email := "test@example.com"
	password := "securepassword"
	hashedPassword := "hashed_password_value"

	user := &entities.User{
		ID:       uuid.New().String(),
		Email:    email,
		Name:     "Test User",
		Password: hashedPassword,
	}

	// Mock GetByEmail to return a user
	mockUserRepo.On("GetByEmail", ctx, email).Return(user, nil)

	// Mock VerifyPassword to return true
	mockHashService.On("VerifyPassword", ctx, hashedPassword, password).Return(true, nil)

	useCase := use_cases.NewUserUseCase(mockUserRepo, mockHashService)

	// Act
	authenticatedUser, err := useCase.AuthenticateUser(ctx, email, password)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, authenticatedUser)
	assert.Equal(t, user.ID, authenticatedUser.ID)
	assert.Equal(t, email, authenticatedUser.Email)
	assert.Empty(t, authenticatedUser.Password) // Password should not be returned

	mockUserRepo.AssertExpectations(t)
	mockHashService.AssertExpectations(t)
}

func TestAuthenticateUser_InvalidCredentials(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockUserRepo := new(MockUserRepository)
	mockHashService := new(MockHashService)

	email := "test@example.com"
	password := "wrongpassword"
	hashedPassword := "hashed_password_value"

	user := &entities.User{
		ID:       uuid.New().String(),
		Email:    email,
		Name:     "Test User",
		Password: hashedPassword,
	}

	// Mock GetByEmail to return a user
	mockUserRepo.On("GetByEmail", ctx, email).Return(user, nil)

	// Mock VerifyPassword to return false (invalid password)
	mockHashService.On("VerifyPassword", ctx, hashedPassword, password).Return(false, nil)

	useCase := use_cases.NewUserUseCase(mockUserRepo, mockHashService)

	// Act
	authenticatedUser, err := useCase.AuthenticateUser(ctx, email, password)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, authenticatedUser)
	assert.Contains(t, err.Error(), "invalid credentials")

	mockUserRepo.AssertExpectations(t)
	mockHashService.AssertExpectations(t)
}

func TestAuthenticateUser_UserNotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockUserRepo := new(MockUserRepository)
	mockHashService := new(MockHashService)

	email := "nonexistent@example.com"
	password := "anypassword"

	// Mock GetByEmail to return nil (user not found)
	mockUserRepo.On("GetByEmail", ctx, email).Return(nil, nil)

	useCase := use_cases.NewUserUseCase(mockUserRepo, mockHashService)

	// Act
	authenticatedUser, err := useCase.AuthenticateUser(ctx, email, password)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, authenticatedUser)
	assert.Contains(t, err.Error(), "invalid credentials")

	mockUserRepo.AssertExpectations(t)
	// VerifyPassword should not be called if user doesn't exist
	mockHashService.AssertNotCalled(t, "VerifyPassword")
}
