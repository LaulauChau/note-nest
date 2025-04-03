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

// MockNoteRepository mocks the NoteRepository interface
type MockNoteRepository struct {
	mock.Mock
}

func (m *MockNoteRepository) Create(ctx context.Context, note *entities.Note) error {
	args := m.Called(ctx, note)
	return args.Error(0)
}

func (m *MockNoteRepository) GetByID(ctx context.Context, id string) (*entities.Note, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Note), args.Error(1)
}

func (m *MockNoteRepository) GetByUserID(ctx context.Context, userID string) ([]*entities.Note, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.Note), args.Error(1)
}

func (m *MockNoteRepository) GetArchivedByUserID(ctx context.Context, userID string) ([]*entities.Note, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.Note), args.Error(1)
}

func (m *MockNoteRepository) Update(ctx context.Context, note *entities.Note) error {
	args := m.Called(ctx, note)
	return args.Error(0)
}

func (m *MockNoteRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestCreateNote(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockNoteRepo := new(MockNoteRepository)
	mockUserRepo := new(MockUserRepository)

	userID := uuid.New().String()
	title := "Test Note"
	content := "This is a test note content."
	label := "test-label"

	// Mock user repository to return a valid user
	user := &entities.User{
		ID:    userID,
		Email: "test@example.com",
		Name:  "Test User",
	}
	mockUserRepo.On("GetByID", ctx, userID).Return(user, nil)

	// Mock note repository to accept the create
	mockNoteRepo.On("Create", ctx, mock.MatchedBy(func(note *entities.Note) bool {
		return note.UserID == userID &&
			note.Title == title &&
			note.Content == content &&
			note.Label == label &&
			note.IsArchived == false
	})).Return(nil)

	useCase := use_cases.NewNoteUseCase(mockNoteRepo, mockUserRepo)

	// Act
	note, err := useCase.CreateNote(ctx, userID, title, content, label)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, note)
	assert.Equal(t, userID, note.UserID)
	assert.Equal(t, title, note.Title)
	assert.Equal(t, content, note.Content)
	assert.Equal(t, label, note.Label)
	assert.False(t, note.IsArchived)
	assert.NotEmpty(t, note.ID)
	assert.NotZero(t, note.CreatedAt)
	assert.NotZero(t, note.UpdatedAt)

	mockUserRepo.AssertExpectations(t)
	mockNoteRepo.AssertExpectations(t)
}

func TestCreateNote_UserNotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockNoteRepo := new(MockNoteRepository)
	mockUserRepo := new(MockUserRepository)

	userID := uuid.New().String()
	title := "Test Note"
	content := "This is a test note content."
	label := "test-label"

	// Mock user repository to return nil (user not found)
	mockUserRepo.On("GetByID", ctx, userID).Return(nil, nil)

	useCase := use_cases.NewNoteUseCase(mockNoteRepo, mockUserRepo)

	// Act
	note, err := useCase.CreateNote(ctx, userID, title, content, label)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, note)
	assert.Contains(t, err.Error(), "user not found")

	mockUserRepo.AssertExpectations(t)
	// Create should not be called if user is not found
	mockNoteRepo.AssertNotCalled(t, "Create")
}

func TestGetNoteByID(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockNoteRepo := new(MockNoteRepository)
	mockUserRepo := new(MockUserRepository)

	userID := uuid.New().String()
	noteID := uuid.New().String()

	note := &entities.Note{
		ID:         noteID,
		UserID:     userID,
		Title:      "Test Note",
		Content:    "This is the content",
		IsArchived: false,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Mock note repository to return a note
	mockNoteRepo.On("GetByID", ctx, noteID).Return(note, nil)

	useCase := use_cases.NewNoteUseCase(mockNoteRepo, mockUserRepo)

	// Act
	result, err := useCase.GetNoteByID(ctx, noteID, userID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, note, result)
	mockNoteRepo.AssertExpectations(t)
}

func TestGetNoteByID_NotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockNoteRepo := new(MockNoteRepository)
	mockUserRepo := new(MockUserRepository)

	userID := uuid.New().String()
	noteID := uuid.New().String()

	// Mock note repository to return nil (note not found)
	mockNoteRepo.On("GetByID", ctx, noteID).Return(nil, nil)

	useCase := use_cases.NewNoteUseCase(mockNoteRepo, mockUserRepo)

	// Act
	result, err := useCase.GetNoteByID(ctx, noteID, userID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "note not found")
	mockNoteRepo.AssertExpectations(t)
}

func TestGetNoteByID_WrongUser(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockNoteRepo := new(MockNoteRepository)
	mockUserRepo := new(MockUserRepository)

	userID := uuid.New().String()
	anotherUserID := uuid.New().String()
	noteID := uuid.New().String()

	note := &entities.Note{
		ID:         noteID,
		UserID:     anotherUserID, // Note belongs to another user
		Title:      "Test Note",
		Content:    "This is the content",
		IsArchived: false,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Mock note repository to return a note that belongs to another user
	mockNoteRepo.On("GetByID", ctx, noteID).Return(note, nil)

	useCase := use_cases.NewNoteUseCase(mockNoteRepo, mockUserRepo)

	// Act
	result, err := useCase.GetNoteByID(ctx, noteID, userID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "note not found")
	mockNoteRepo.AssertExpectations(t)
}

func TestGetActiveNotes(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockNoteRepo := new(MockNoteRepository)
	mockUserRepo := new(MockUserRepository)

	userID := uuid.New().String()

	user := &entities.User{
		ID:    userID,
		Email: "test@example.com",
		Name:  "Test User",
	}

	notes := []*entities.Note{
		{
			ID:         uuid.New().String(),
			UserID:     userID,
			Title:      "Note 1",
			Content:    "Content 1",
			IsArchived: false,
		},
		{
			ID:         uuid.New().String(),
			UserID:     userID,
			Title:      "Note 2",
			Content:    "Content 2",
			IsArchived: false,
		},
	}

	// Mock user repository to return a valid user
	mockUserRepo.On("GetByID", ctx, userID).Return(user, nil)

	// Mock note repository to return notes
	mockNoteRepo.On("GetByUserID", ctx, userID).Return(notes, nil)

	useCase := use_cases.NewNoteUseCase(mockNoteRepo, mockUserRepo)

	// Act
	result, err := useCase.GetActiveNotes(ctx, userID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, notes, result)
	assert.Len(t, result, 2)
	mockUserRepo.AssertExpectations(t)
	mockNoteRepo.AssertExpectations(t)
}

func TestGetArchivedNotes(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockNoteRepo := new(MockNoteRepository)
	mockUserRepo := new(MockUserRepository)

	userID := uuid.New().String()

	user := &entities.User{
		ID:    userID,
		Email: "test@example.com",
		Name:  "Test User",
	}

	archivedNotes := []*entities.Note{
		{
			ID:         uuid.New().String(),
			UserID:     userID,
			Title:      "Archived Note 1",
			Content:    "Content 1",
			IsArchived: true,
		},
		{
			ID:         uuid.New().String(),
			UserID:     userID,
			Title:      "Archived Note 2",
			Content:    "Content 2",
			IsArchived: true,
		},
	}

	// Mock user repository to return a valid user
	mockUserRepo.On("GetByID", ctx, userID).Return(user, nil)

	// Mock note repository to return archived notes
	mockNoteRepo.On("GetArchivedByUserID", ctx, userID).Return(archivedNotes, nil)

	useCase := use_cases.NewNoteUseCase(mockNoteRepo, mockUserRepo)

	// Act
	result, err := useCase.GetArchivedNotes(ctx, userID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, archivedNotes, result)
	assert.Len(t, result, 2)
	mockUserRepo.AssertExpectations(t)
	mockNoteRepo.AssertExpectations(t)
}
