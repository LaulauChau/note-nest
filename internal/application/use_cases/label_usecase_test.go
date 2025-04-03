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

func TestCreateLabel(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockLabelRepo := new(MockLabelRepository)
	mockUserRepo := new(MockUserRepository)
	mockNoteRepo := new(MockNoteRepository)

	userID := uuid.New().String()
	name := "Work"
	color := "#ff5733"

	// Mock user repository to return a valid user
	user := &entities.User{
		ID:    userID,
		Email: "test@example.com",
		Name:  "Test User",
	}
	mockUserRepo.On("GetByID", ctx, userID).Return(user, nil)

	// Mock label repository to check if label name exists
	mockLabelRepo.On("GetByName", ctx, userID, name).Return(nil, nil)

	// Mock label repository to create the label
	mockLabelRepo.On("Create", ctx, mock.MatchedBy(func(label *entities.Label) bool {
		return label.UserID == userID &&
			label.Name == name &&
			label.Color == color
	})).Return(nil)

	useCase := use_cases.NewLabelUseCase(mockLabelRepo, mockUserRepo, mockNoteRepo)

	// Act
	label, err := useCase.CreateLabel(ctx, userID, name, color)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, label)
	assert.Equal(t, userID, label.UserID)
	assert.Equal(t, name, label.Name)
	assert.Equal(t, color, label.Color)
	assert.NotEmpty(t, label.ID)
	assert.NotZero(t, label.CreatedAt)
	assert.NotZero(t, label.UpdatedAt)

	mockUserRepo.AssertExpectations(t)
	mockLabelRepo.AssertExpectations(t)
}

func TestCreateLabel_UserNotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockLabelRepo := new(MockLabelRepository)
	mockUserRepo := new(MockUserRepository)
	mockNoteRepo := new(MockNoteRepository)

	userID := uuid.New().String()
	name := "Work"
	color := "#ff5733"

	// Mock user repository to return nil (user not found)
	mockUserRepo.On("GetByID", ctx, userID).Return(nil, nil)

	useCase := use_cases.NewLabelUseCase(mockLabelRepo, mockUserRepo, mockNoteRepo)

	// Act
	label, err := useCase.CreateLabel(ctx, userID, name, color)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, label)
	assert.Contains(t, err.Error(), "user not found")

	mockUserRepo.AssertExpectations(t)
	// CreateLabel should not be called if user is not found
	mockLabelRepo.AssertNotCalled(t, "Create")
}

func TestCreateLabel_DuplicateName(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockLabelRepo := new(MockLabelRepository)
	mockUserRepo := new(MockUserRepository)
	mockNoteRepo := new(MockNoteRepository)

	userID := uuid.New().String()
	name := "Work"
	color := "#ff5733"

	// Mock user repository to return a valid user
	user := &entities.User{
		ID:    userID,
		Email: "test@example.com",
		Name:  "Test User",
	}
	mockUserRepo.On("GetByID", ctx, userID).Return(user, nil)

	// Mock label repository to return an existing label with the same name
	existingLabel := &entities.Label{
		ID:     uuid.New().String(),
		UserID: userID,
		Name:   name,
		Color:  "#000000",
	}
	mockLabelRepo.On("GetByName", ctx, userID, name).Return(existingLabel, nil)

	useCase := use_cases.NewLabelUseCase(mockLabelRepo, mockUserRepo, mockNoteRepo)

	// Act
	label, err := useCase.CreateLabel(ctx, userID, name, color)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, label)
	assert.Contains(t, err.Error(), "label with this name already exists")

	mockUserRepo.AssertExpectations(t)
	mockLabelRepo.AssertExpectations(t)
	// Create should not be called if label name already exists
	mockLabelRepo.AssertNotCalled(t, "Create")
}

func TestGetLabelByID(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockLabelRepo := new(MockLabelRepository)
	mockUserRepo := new(MockUserRepository)
	mockNoteRepo := new(MockNoteRepository)

	userID := uuid.New().String()
	labelID := uuid.New().String()

	label := &entities.Label{
		ID:     labelID,
		UserID: userID,
		Name:   "Work",
		Color:  "#ff5733",
	}

	// Mock label repository to return a label
	mockLabelRepo.On("GetByID", ctx, labelID).Return(label, nil)

	useCase := use_cases.NewLabelUseCase(mockLabelRepo, mockUserRepo, mockNoteRepo)

	// Act
	result, err := useCase.GetLabelByID(ctx, labelID, userID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, label, result)
	mockLabelRepo.AssertExpectations(t)
}

func TestGetLabelByID_NotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockLabelRepo := new(MockLabelRepository)
	mockUserRepo := new(MockUserRepository)
	mockNoteRepo := new(MockNoteRepository)

	userID := uuid.New().String()
	labelID := uuid.New().String()

	// Mock label repository to return nil (label not found)
	mockLabelRepo.On("GetByID", ctx, labelID).Return(nil, nil)

	useCase := use_cases.NewLabelUseCase(mockLabelRepo, mockUserRepo, mockNoteRepo)

	// Act
	result, err := useCase.GetLabelByID(ctx, labelID, userID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "label not found")
	mockLabelRepo.AssertExpectations(t)
}

func TestGetLabelByID_WrongUser(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockLabelRepo := new(MockLabelRepository)
	mockUserRepo := new(MockUserRepository)
	mockNoteRepo := new(MockNoteRepository)

	userID := uuid.New().String()
	anotherUserID := uuid.New().String()
	labelID := uuid.New().String()

	label := &entities.Label{
		ID:     labelID,
		UserID: anotherUserID, // Label belongs to another user
		Name:   "Work",
		Color:  "#ff5733",
	}

	// Mock label repository to return a label that belongs to another user
	mockLabelRepo.On("GetByID", ctx, labelID).Return(label, nil)

	useCase := use_cases.NewLabelUseCase(mockLabelRepo, mockUserRepo, mockNoteRepo)

	// Act
	result, err := useCase.GetLabelByID(ctx, labelID, userID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "label not found")
	mockLabelRepo.AssertExpectations(t)
}

func TestGetLabelsByUser(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockLabelRepo := new(MockLabelRepository)
	mockUserRepo := new(MockUserRepository)
	mockNoteRepo := new(MockNoteRepository)

	userID := uuid.New().String()

	user := &entities.User{
		ID:    userID,
		Email: "test@example.com",
		Name:  "Test User",
	}

	labels := []*entities.Label{
		{
			ID:     uuid.New().String(),
			UserID: userID,
			Name:   "Work",
			Color:  "#ff5733",
		},
		{
			ID:     uuid.New().String(),
			UserID: userID,
			Name:   "Personal",
			Color:  "#33ff57",
		},
	}

	// Mock user repository to return a valid user
	mockUserRepo.On("GetByID", ctx, userID).Return(user, nil)

	// Mock label repository to return labels
	mockLabelRepo.On("GetByUserID", ctx, userID).Return(labels, nil)

	useCase := use_cases.NewLabelUseCase(mockLabelRepo, mockUserRepo, mockNoteRepo)

	// Act
	result, err := useCase.GetLabelsByUser(ctx, userID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, labels, result)
	assert.Len(t, result, 2)
	mockUserRepo.AssertExpectations(t)
	mockLabelRepo.AssertExpectations(t)
}

func TestUpdateLabel(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockLabelRepo := new(MockLabelRepository)
	mockUserRepo := new(MockUserRepository)
	mockNoteRepo := new(MockNoteRepository)

	userID := uuid.New().String()
	labelID := uuid.New().String()

	// Existing label
	existingLabel := &entities.Label{
		ID:        labelID,
		UserID:    userID,
		Name:      "Work",
		Color:     "#ff5733",
		CreatedAt: time.Now().Add(-24 * time.Hour),
		UpdatedAt: time.Now().Add(-24 * time.Hour),
	}

	// Updated fields
	newName := "Job"
	newColor := "#33ff57"

	// Mock label repository to return the existing label
	mockLabelRepo.On("GetByID", ctx, labelID).Return(existingLabel, nil)

	// Mock label repository to check if the new name already exists
	mockLabelRepo.On("GetByName", ctx, userID, newName).Return(nil, nil)

	// Mock label repository to update the label. We verify the result after the call.
	mockLabelRepo.On("Update", ctx, mock.AnythingOfType("*entities.Label")).Run(func(args mock.Arguments) {
		// Get the label passed to Update
		label := args.Get(1).(*entities.Label)
		// Ensure the UpdatedAt time is set to a newer time
		label.UpdatedAt = time.Now()
	}).Return(nil)

	useCase := use_cases.NewLabelUseCase(mockLabelRepo, mockUserRepo, mockNoteRepo)

	// Introduce a small delay to ensure UpdatedAt changes measurably
	time.Sleep(50 * time.Millisecond)

	// Act
	updatedLabel, err := useCase.UpdateLabel(ctx, labelID, userID, newName, newColor)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, updatedLabel)
	assert.Equal(t, labelID, updatedLabel.ID)
	assert.Equal(t, userID, updatedLabel.UserID)
	assert.Equal(t, newName, updatedLabel.Name)
	assert.Equal(t, newColor, updatedLabel.Color)
	assert.Equal(t, existingLabel.CreatedAt, updatedLabel.CreatedAt) // Created time should not change

	mockLabelRepo.AssertExpectations(t)
}

func TestUpdateLabel_NotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockLabelRepo := new(MockLabelRepository)
	mockUserRepo := new(MockUserRepository)
	mockNoteRepo := new(MockNoteRepository)

	userID := uuid.New().String()
	labelID := uuid.New().String()
	newName := "Job"
	newColor := "#33ff57"

	// Mock label repository to return nil (label not found)
	mockLabelRepo.On("GetByID", ctx, labelID).Return(nil, nil)

	useCase := use_cases.NewLabelUseCase(mockLabelRepo, mockUserRepo, mockNoteRepo)

	// Act
	updatedLabel, err := useCase.UpdateLabel(ctx, labelID, userID, newName, newColor)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, updatedLabel)
	assert.Contains(t, err.Error(), "label not found")

	mockLabelRepo.AssertExpectations(t)
	// Update should not be called if label is not found
	mockLabelRepo.AssertNotCalled(t, "Update")
}

func TestUpdateLabel_WrongUser(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockLabelRepo := new(MockLabelRepository)
	mockUserRepo := new(MockUserRepository)
	mockNoteRepo := new(MockNoteRepository)

	userID := uuid.New().String()
	anotherUserID := uuid.New().String()
	labelID := uuid.New().String()
	newName := "Job"
	newColor := "#33ff57"

	// Label belongs to another user
	label := &entities.Label{
		ID:     labelID,
		UserID: anotherUserID,
		Name:   "Work",
		Color:  "#ff5733",
	}

	// Mock label repository to return a label that belongs to another user
	mockLabelRepo.On("GetByID", ctx, labelID).Return(label, nil)

	useCase := use_cases.NewLabelUseCase(mockLabelRepo, mockUserRepo, mockNoteRepo)

	// Act
	updatedLabel, err := useCase.UpdateLabel(ctx, labelID, userID, newName, newColor)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, updatedLabel)
	assert.Contains(t, err.Error(), "label not found")

	mockLabelRepo.AssertExpectations(t)
	// Update should not be called if label belongs to another user
	mockLabelRepo.AssertNotCalled(t, "Update")
}

func TestUpdateLabel_DuplicateName(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockLabelRepo := new(MockLabelRepository)
	mockUserRepo := new(MockUserRepository)
	mockNoteRepo := new(MockNoteRepository)

	userID := uuid.New().String()
	labelID := uuid.New().String()
	anotherLabelID := uuid.New().String()

	// Existing label
	existingLabel := &entities.Label{
		ID:     labelID,
		UserID: userID,
		Name:   "Work",
		Color:  "#ff5733",
	}

	// Another label with the name we want to update to
	anotherLabel := &entities.Label{
		ID:     anotherLabelID,
		UserID: userID,
		Name:   "Personal", // This is the name we're trying to update to
		Color:  "#33ff57",
	}

	newName := "Personal" // This name is already taken by another label
	newColor := "#33ff57"

	// Mock label repository to return the existing label
	mockLabelRepo.On("GetByID", ctx, labelID).Return(existingLabel, nil)

	// Mock label repository to check if the new name already exists
	mockLabelRepo.On("GetByName", ctx, userID, newName).Return(anotherLabel, nil)

	useCase := use_cases.NewLabelUseCase(mockLabelRepo, mockUserRepo, mockNoteRepo)

	// Act
	updatedLabel, err := useCase.UpdateLabel(ctx, labelID, userID, newName, newColor)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, updatedLabel)
	assert.Contains(t, err.Error(), "label with this name already exists")

	mockLabelRepo.AssertExpectations(t)
	// Update should not be called if name is duplicate
	mockLabelRepo.AssertNotCalled(t, "Update")
}

func TestDeleteLabel(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockLabelRepo := new(MockLabelRepository)
	mockUserRepo := new(MockUserRepository)
	mockNoteRepo := new(MockNoteRepository)

	userID := uuid.New().String()
	labelID := uuid.New().String()

	// Existing label
	label := &entities.Label{
		ID:     labelID,
		UserID: userID,
		Name:   "Work",
		Color:  "#ff5733",
	}

	// Mock label repository to return the existing label
	mockLabelRepo.On("GetByID", ctx, labelID).Return(label, nil)

	// Mock label repository to delete the label
	mockLabelRepo.On("Delete", ctx, labelID).Return(nil)

	useCase := use_cases.NewLabelUseCase(mockLabelRepo, mockUserRepo, mockNoteRepo)

	// Act
	err := useCase.DeleteLabel(ctx, labelID, userID)

	// Assert
	assert.NoError(t, err)
	mockLabelRepo.AssertExpectations(t)
}

func TestDeleteLabel_NotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockLabelRepo := new(MockLabelRepository)
	mockUserRepo := new(MockUserRepository)
	mockNoteRepo := new(MockNoteRepository)

	userID := uuid.New().String()
	labelID := uuid.New().String()

	// Mock label repository to return nil (label not found)
	mockLabelRepo.On("GetByID", ctx, labelID).Return(nil, nil)

	useCase := use_cases.NewLabelUseCase(mockLabelRepo, mockUserRepo, mockNoteRepo)

	// Act
	err := useCase.DeleteLabel(ctx, labelID, userID)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "label not found")

	mockLabelRepo.AssertExpectations(t)
	// Delete should not be called if label is not found
	mockLabelRepo.AssertNotCalled(t, "Delete")
}

func TestDeleteLabel_WrongUser(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockLabelRepo := new(MockLabelRepository)
	mockUserRepo := new(MockUserRepository)
	mockNoteRepo := new(MockNoteRepository)

	userID := uuid.New().String()
	anotherUserID := uuid.New().String()
	labelID := uuid.New().String()

	// Label belongs to another user
	label := &entities.Label{
		ID:     labelID,
		UserID: anotherUserID,
		Name:   "Work",
		Color:  "#ff5733",
	}

	// Mock label repository to return a label that belongs to another user
	mockLabelRepo.On("GetByID", ctx, labelID).Return(label, nil)

	useCase := use_cases.NewLabelUseCase(mockLabelRepo, mockUserRepo, mockNoteRepo)

	// Act
	err := useCase.DeleteLabel(ctx, labelID, userID)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "label not found")

	mockLabelRepo.AssertExpectations(t)
	// Delete should not be called if label belongs to another user
	mockLabelRepo.AssertNotCalled(t, "Delete")
}

func TestAddLabelToNote(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockLabelRepo := new(MockLabelRepository)
	mockUserRepo := new(MockUserRepository)
	mockNoteRepo := new(MockNoteRepository)

	userID := uuid.New().String()
	noteID := uuid.New().String()
	labelID := uuid.New().String()

	// Mock note repository to return a valid note owned by the user
	note := &entities.Note{
		ID:     noteID,
		UserID: userID,
		Title:  "Test Note",
	}
	mockNoteRepo.On("GetByID", ctx, noteID).Return(note, nil)

	// Mock label repository to return a valid label owned by the user
	label := &entities.Label{
		ID:     labelID,
		UserID: userID,
		Name:   "Work",
		Color:  "#ff5733",
	}
	mockLabelRepo.On("GetByID", ctx, labelID).Return(label, nil)

	// Mock label repository to add the label to the note
	mockLabelRepo.On("AddLabelToNote", ctx, noteID, labelID).Return(nil)

	useCase := use_cases.NewLabelUseCase(mockLabelRepo, mockUserRepo, mockNoteRepo)

	// Act
	err := useCase.AddLabelToNote(ctx, noteID, labelID, userID)

	// Assert
	assert.NoError(t, err)
	mockNoteRepo.AssertExpectations(t)
	mockLabelRepo.AssertExpectations(t)
}

func TestAddLabelToNote_NoteNotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockLabelRepo := new(MockLabelRepository)
	mockUserRepo := new(MockUserRepository)
	mockNoteRepo := new(MockNoteRepository)

	userID := uuid.New().String()
	noteID := uuid.New().String()
	labelID := uuid.New().String()

	// Mock note repository to return nil (note not found)
	mockNoteRepo.On("GetByID", ctx, noteID).Return(nil, nil)

	useCase := use_cases.NewLabelUseCase(mockLabelRepo, mockUserRepo, mockNoteRepo)

	// Act
	err := useCase.AddLabelToNote(ctx, noteID, labelID, userID)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "note not found")
	mockNoteRepo.AssertExpectations(t)
	// Label repository methods should not be called if note is not found
	mockLabelRepo.AssertNotCalled(t, "GetByID")
	mockLabelRepo.AssertNotCalled(t, "AddLabelToNote")
}

func TestAddLabelToNote_NoteWrongUser(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockLabelRepo := new(MockLabelRepository)
	mockUserRepo := new(MockUserRepository)
	mockNoteRepo := new(MockNoteRepository)

	userID := uuid.New().String()
	anotherUserID := uuid.New().String()
	noteID := uuid.New().String()
	labelID := uuid.New().String()

	// Mock note repository to return a note that belongs to another user
	note := &entities.Note{
		ID:     noteID,
		UserID: anotherUserID,
		Title:  "Test Note",
	}
	mockNoteRepo.On("GetByID", ctx, noteID).Return(note, nil)

	useCase := use_cases.NewLabelUseCase(mockLabelRepo, mockUserRepo, mockNoteRepo)

	// Act
	err := useCase.AddLabelToNote(ctx, noteID, labelID, userID)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "note not found")
	mockNoteRepo.AssertExpectations(t)
	// Label repository methods should not be called if note belongs to another user
	mockLabelRepo.AssertNotCalled(t, "GetByID")
	mockLabelRepo.AssertNotCalled(t, "AddLabelToNote")
}

func TestAddLabelToNote_LabelNotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockLabelRepo := new(MockLabelRepository)
	mockUserRepo := new(MockUserRepository)
	mockNoteRepo := new(MockNoteRepository)

	userID := uuid.New().String()
	noteID := uuid.New().String()
	labelID := uuid.New().String()

	// Mock note repository to return a valid note owned by the user
	note := &entities.Note{
		ID:     noteID,
		UserID: userID,
		Title:  "Test Note",
	}
	mockNoteRepo.On("GetByID", ctx, noteID).Return(note, nil)

	// Mock label repository to return nil (label not found)
	mockLabelRepo.On("GetByID", ctx, labelID).Return(nil, nil)

	useCase := use_cases.NewLabelUseCase(mockLabelRepo, mockUserRepo, mockNoteRepo)

	// Act
	err := useCase.AddLabelToNote(ctx, noteID, labelID, userID)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "label not found")
	mockNoteRepo.AssertExpectations(t)
	mockLabelRepo.AssertExpectations(t)
	// AddLabelToNote should not be called if label is not found
	mockLabelRepo.AssertNotCalled(t, "AddLabelToNote")
}

func TestAddLabelToNote_LabelWrongUser(t *testing.T) {
	// Arrange
	ctx := context.Background()
	mockLabelRepo := new(MockLabelRepository)
	mockUserRepo := new(MockUserRepository)
	mockNoteRepo := new(MockNoteRepository)

	userID := uuid.New().String()
	anotherUserID := uuid.New().String()
	noteID := uuid.New().String()
	labelID := uuid.New().String()

	// Mock note repository to return a valid note owned by the user
	note := &entities.Note{
		ID:     noteID,
		UserID: userID,
		Title:  "Test Note",
	}
	mockNoteRepo.On("GetByID", ctx, noteID).Return(note, nil)

	// Mock label repository to return a label that belongs to another user
	label := &entities.Label{
		ID:     labelID,
		UserID: anotherUserID,
		Name:   "Work",
		Color:  "#ff5733",
	}
	mockLabelRepo.On("GetByID", ctx, labelID).Return(label, nil)

	useCase := use_cases.NewLabelUseCase(mockLabelRepo, mockUserRepo, mockNoteRepo)

	// Act
	err := useCase.AddLabelToNote(ctx, noteID, labelID, userID)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "label not found")
	mockNoteRepo.AssertExpectations(t)
	mockLabelRepo.AssertExpectations(t)
	// AddLabelToNote should not be called if label belongs to another user
	mockLabelRepo.AssertNotCalled(t, "AddLabelToNote")
}
