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

// MockLabelRepository mocks the LabelRepository interface
type MockLabelRepository struct {
	mock.Mock
}

// Add mock implementations for all methods in repositories.LabelRepository interface
func (m *MockLabelRepository) Create(ctx context.Context, label *entities.Label) error {
	args := m.Called(ctx, label)
	return args.Error(0)
}

func (m *MockLabelRepository) GetByID(ctx context.Context, id string) (*entities.Label, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Label), args.Error(1)
}

func (m *MockLabelRepository) GetByUserID(ctx context.Context, userID string) ([]*entities.Label, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.Label), args.Error(1)
}

func (m *MockLabelRepository) GetByName(ctx context.Context, userID, name string) (*entities.Label, error) {
	args := m.Called(ctx, userID, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Label), args.Error(1)
}

func (m *MockLabelRepository) Update(ctx context.Context, label *entities.Label) error {
	args := m.Called(ctx, label)
	return args.Error(0)
}

func (m *MockLabelRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockLabelRepository) AddLabelToNote(ctx context.Context, noteID, labelID string) error {
	args := m.Called(ctx, noteID, labelID)
	return args.Error(0)
}

func (m *MockLabelRepository) RemoveLabelFromNote(ctx context.Context, noteID, labelID string) error {
	args := m.Called(ctx, noteID, labelID)
	return args.Error(0)
}

func (m *MockLabelRepository) GetLabelsForNote(ctx context.Context, noteID string) ([]*entities.Label, error) {
	args := m.Called(ctx, noteID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.Label), args.Error(1)
}

func (m *MockLabelRepository) GetNotesForLabel(ctx context.Context, labelID string) ([]string, error) {
	args := m.Called(ctx, labelID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func TestCreateNote(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockNoteRepo := new(MockNoteRepository)
	mockUserRepo := new(MockUserRepository)
	mockLabelRepo := new(MockLabelRepository)

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

	useCase := use_cases.NewNoteUseCase(mockNoteRepo, mockUserRepo, mockLabelRepo)

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
	mockLabelRepo := new(MockLabelRepository)

	userID := uuid.New().String()
	title := "Test Note"
	content := "This is a test note content."
	label := "test-label"

	// Mock user repository to return nil (user not found)
	mockUserRepo.On("GetByID", ctx, userID).Return(nil, nil)

	useCase := use_cases.NewNoteUseCase(mockNoteRepo, mockUserRepo, mockLabelRepo)

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
	mockLabelRepo := new(MockLabelRepository)

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

	useCase := use_cases.NewNoteUseCase(mockNoteRepo, mockUserRepo, mockLabelRepo)

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
	mockLabelRepo := new(MockLabelRepository)

	userID := uuid.New().String()
	noteID := uuid.New().String()

	// Mock note repository to return nil (note not found)
	mockNoteRepo.On("GetByID", ctx, noteID).Return(nil, nil)

	useCase := use_cases.NewNoteUseCase(mockNoteRepo, mockUserRepo, mockLabelRepo)

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
	mockLabelRepo := new(MockLabelRepository)

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

	useCase := use_cases.NewNoteUseCase(mockNoteRepo, mockUserRepo, mockLabelRepo)

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
	mockLabelRepo := new(MockLabelRepository)

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

	useCase := use_cases.NewNoteUseCase(mockNoteRepo, mockUserRepo, mockLabelRepo)

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
	mockLabelRepo := new(MockLabelRepository)

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

	useCase := use_cases.NewNoteUseCase(mockNoteRepo, mockUserRepo, mockLabelRepo)

	// Act
	result, err := useCase.GetArchivedNotes(ctx, userID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, archivedNotes, result)
	assert.Len(t, result, 2)
	mockUserRepo.AssertExpectations(t)
	mockNoteRepo.AssertExpectations(t)
}

func TestUpdateNote(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockNoteRepo := new(MockNoteRepository)
	mockUserRepo := new(MockUserRepository)
	mockLabelRepo := new(MockLabelRepository)

	userID := uuid.New().String()
	noteID := uuid.New().String()

	// Create a timestamp in the past
	pastTime := time.Now().Add(-1 * time.Hour)

	// Existing note
	existingNote := &entities.Note{
		ID:         noteID,
		UserID:     userID,
		Title:      "Original Title",
		Content:    "Original content",
		IsArchived: false,
		Label:      "original-label",
		CreatedAt:  pastTime,
		UpdatedAt:  pastTime,
	}

	// Updated fields
	newTitle := "Updated Title"
	newContent := "Updated content"
	newLabel := "updated-label"
	newIsArchived := true

	// Mock note repository to return the existing note
	mockNoteRepo.On("GetByID", ctx, noteID).Return(existingNote, nil)

	// Use mock.Anything for the note parameter to avoid matching issues
	mockNoteRepo.On("Update", ctx, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		// Verify the updated note properties within the Run function
		updatedNote := args.Get(1).(*entities.Note)
		assert.Equal(t, noteID, updatedNote.ID)
		assert.Equal(t, userID, updatedNote.UserID)
		assert.Equal(t, newTitle, updatedNote.Title)
		assert.Equal(t, newContent, updatedNote.Content)
		assert.Equal(t, newLabel, updatedNote.Label)
		assert.Equal(t, newIsArchived, updatedNote.IsArchived)
		assert.Equal(t, pastTime, updatedNote.CreatedAt)
		assert.True(t, updatedNote.UpdatedAt.After(pastTime))
	})

	useCase := use_cases.NewNoteUseCase(mockNoteRepo, mockUserRepo, mockLabelRepo)

	// Act
	updatedNote, err := useCase.UpdateNote(ctx, noteID, userID, newTitle, newContent, newLabel, newIsArchived)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, updatedNote)
	assert.Equal(t, noteID, updatedNote.ID)
	assert.Equal(t, userID, updatedNote.UserID)
	assert.Equal(t, newTitle, updatedNote.Title)
	assert.Equal(t, newContent, updatedNote.Content)
	assert.Equal(t, newLabel, updatedNote.Label)
	assert.Equal(t, newIsArchived, updatedNote.IsArchived)
	assert.Equal(t, pastTime, updatedNote.CreatedAt)      // Created time should not change
	assert.True(t, updatedNote.UpdatedAt.After(pastTime)) // Updated time should be newer

	mockNoteRepo.AssertExpectations(t)
}

func TestUpdateNote_NotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockNoteRepo := new(MockNoteRepository)
	mockUserRepo := new(MockUserRepository)
	mockLabelRepo := new(MockLabelRepository)

	userID := uuid.New().String()
	noteID := uuid.New().String()

	// Mock note repository to return nil (note not found)
	mockNoteRepo.On("GetByID", ctx, noteID).Return(nil, nil)

	useCase := use_cases.NewNoteUseCase(mockNoteRepo, mockUserRepo, mockLabelRepo)

	// Act
	updatedNote, err := useCase.UpdateNote(ctx, noteID, userID, "Title", "Content", "label", false)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, updatedNote)
	assert.Contains(t, err.Error(), "note not found")

	mockNoteRepo.AssertExpectations(t)
	// Update should not be called if note is not found
	mockNoteRepo.AssertNotCalled(t, "Update")
}

func TestUpdateNote_WrongUser(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockNoteRepo := new(MockNoteRepository)
	mockUserRepo := new(MockUserRepository)
	mockLabelRepo := new(MockLabelRepository)

	userID := uuid.New().String()
	anotherUserID := uuid.New().String()
	noteID := uuid.New().String()

	// Note belongs to another user
	note := &entities.Note{
		ID:         noteID,
		UserID:     anotherUserID,
		Title:      "Original Title",
		Content:    "Original content",
		IsArchived: false,
		Label:      "original-label",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Mock note repository to return a note that belongs to another user
	mockNoteRepo.On("GetByID", ctx, noteID).Return(note, nil)

	useCase := use_cases.NewNoteUseCase(mockNoteRepo, mockUserRepo, mockLabelRepo)

	// Act
	updatedNote, err := useCase.UpdateNote(ctx, noteID, userID, "Updated Title", "Updated content", "updated-label", true)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, updatedNote)
	assert.Contains(t, err.Error(), "note not found")

	mockNoteRepo.AssertExpectations(t)
	// Update should not be called if note belongs to another user
	mockNoteRepo.AssertNotCalled(t, "Update")
}

func TestDeleteNote(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockNoteRepo := new(MockNoteRepository)
	mockUserRepo := new(MockUserRepository)
	mockLabelRepo := new(MockLabelRepository)

	userID := uuid.New().String()
	noteID := uuid.New().String()

	// Existing note
	note := &entities.Note{
		ID:        noteID,
		UserID:    userID,
		Title:     "Note to Delete",
		Content:   "This note will be deleted",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Mock note repository to return the existing note
	mockNoteRepo.On("GetByID", ctx, noteID).Return(note, nil)

	// Mock note repository to delete the note
	mockNoteRepo.On("Delete", ctx, noteID).Return(nil)

	useCase := use_cases.NewNoteUseCase(mockNoteRepo, mockUserRepo, mockLabelRepo)

	// Act
	err := useCase.DeleteNote(ctx, noteID, userID)

	// Assert
	assert.NoError(t, err)
	mockNoteRepo.AssertExpectations(t)
}

func TestDeleteNote_NotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockNoteRepo := new(MockNoteRepository)
	mockUserRepo := new(MockUserRepository)
	mockLabelRepo := new(MockLabelRepository)

	userID := uuid.New().String()
	noteID := uuid.New().String()

	// Mock note repository to return nil (note not found)
	mockNoteRepo.On("GetByID", ctx, noteID).Return(nil, nil)

	useCase := use_cases.NewNoteUseCase(mockNoteRepo, mockUserRepo, mockLabelRepo)

	// Act
	err := useCase.DeleteNote(ctx, noteID, userID)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "note not found")

	mockNoteRepo.AssertExpectations(t)
	// Delete should not be called if note is not found
	mockNoteRepo.AssertNotCalled(t, "Delete")
}

func TestDeleteNote_WrongUser(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockNoteRepo := new(MockNoteRepository)
	mockUserRepo := new(MockUserRepository)
	mockLabelRepo := new(MockLabelRepository)

	userID := uuid.New().String()
	anotherUserID := uuid.New().String()
	noteID := uuid.New().String()

	// Note belongs to another user
	note := &entities.Note{
		ID:        noteID,
		UserID:    anotherUserID,
		Title:     "Another User's Note",
		Content:   "This note belongs to another user",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Mock note repository to return a note that belongs to another user
	mockNoteRepo.On("GetByID", ctx, noteID).Return(note, nil)

	useCase := use_cases.NewNoteUseCase(mockNoteRepo, mockUserRepo, mockLabelRepo)

	// Act
	err := useCase.DeleteNote(ctx, noteID, userID)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "note not found")

	mockNoteRepo.AssertExpectations(t)
	// Delete should not be called if note belongs to another user
	mockNoteRepo.AssertNotCalled(t, "Delete")
}
