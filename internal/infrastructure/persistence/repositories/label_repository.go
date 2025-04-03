package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/LaulauChau/note-nest/internal/domain/entities"
	"github.com/LaulauChau/note-nest/internal/domain/repositories"
)

type LabelRepositoryImpl struct {
	q *Queries
}

func NewLabelRepository(q *Queries) repositories.LabelRepository {
	return &LabelRepositoryImpl{q: q}
}

func (r *LabelRepositoryImpl) Create(ctx context.Context, label *entities.Label) error {
	// Parse the user ID
	userID, err := uuid.Parse(label.UserID)
	if err != nil {
		return err
	}

	// Parse the label ID
	labelID, err := uuid.Parse(label.ID)
	if err != nil {
		return err
	}

	params := CreateLabelParams{
		ID:        labelID.String(),
		UserID:    userID.String(),
		Name:      label.Name,
		Color:     label.Color,
		CreatedAt: label.CreatedAt,
		UpdatedAt: label.UpdatedAt,
	}

	_, err = r.q.CreateLabel(ctx, params)
	return err
}

func (r *LabelRepositoryImpl) GetByID(ctx context.Context, id string) (*entities.Label, error) {
	labelID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	label, err := r.q.GetLabelByID(ctx, labelID.String())
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, nil
		}
		return nil, err
	}

	return &entities.Label{
		ID:        label.ID,
		UserID:    label.UserID,
		Name:      label.Name,
		Color:     label.Color,
		CreatedAt: label.CreatedAt,
		UpdatedAt: label.UpdatedAt,
	}, nil
}

func (r *LabelRepositoryImpl) GetByUserID(ctx context.Context, userID string) ([]*entities.Label, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	labels, err := r.q.GetLabelsByUserID(ctx, userUUID.String())
	if err != nil {
		return nil, err
	}

	result := make([]*entities.Label, len(labels))
	for i, label := range labels {
		result[i] = &entities.Label{
			ID:        label.ID,
			UserID:    label.UserID,
			Name:      label.Name,
			Color:     label.Color,
			CreatedAt: label.CreatedAt,
			UpdatedAt: label.UpdatedAt,
		}
	}

	return result, nil
}

func (r *LabelRepositoryImpl) GetByName(ctx context.Context, userID, name string) (*entities.Label, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	params := GetLabelByNameParams{
		UserID: userUUID.String(),
		Name:   name,
	}

	label, err := r.q.GetLabelByName(ctx, params)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, nil
		}
		return nil, err
	}

	return &entities.Label{
		ID:        label.ID,
		UserID:    label.UserID,
		Name:      label.Name,
		Color:     label.Color,
		CreatedAt: label.CreatedAt,
		UpdatedAt: label.UpdatedAt,
	}, nil
}

func (r *LabelRepositoryImpl) Update(ctx context.Context, label *entities.Label) error {
	labelID, err := uuid.Parse(label.ID)
	if err != nil {
		return err
	}

	params := UpdateLabelParams{
		ID:        labelID.String(),
		Name:      label.Name,
		Color:     label.Color,
		UpdatedAt: time.Now(),
	}

	return r.q.UpdateLabel(ctx, params)
}

func (r *LabelRepositoryImpl) Delete(ctx context.Context, id string) error {
	labelID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	return r.q.DeleteLabel(ctx, labelID.String())
}

func (r *LabelRepositoryImpl) AddLabelToNote(ctx context.Context, noteID, labelID string) error {
	noteUUID, err := uuid.Parse(noteID)
	if err != nil {
		return err
	}

	labelUUID, err := uuid.Parse(labelID)
	if err != nil {
		return err
	}

	params := AddLabelToNoteParams{
		NoteID:  noteUUID.String(),
		LabelID: labelUUID.String(),
	}

	return r.q.AddLabelToNote(ctx, params)
}

func (r *LabelRepositoryImpl) RemoveLabelFromNote(ctx context.Context, noteID, labelID string) error {
	noteUUID, err := uuid.Parse(noteID)
	if err != nil {
		return err
	}

	labelUUID, err := uuid.Parse(labelID)
	if err != nil {
		return err
	}

	params := RemoveLabelFromNoteParams{
		NoteID:  noteUUID.String(),
		LabelID: labelUUID.String(),
	}

	return r.q.RemoveLabelFromNote(ctx, params)
}

func (r *LabelRepositoryImpl) GetLabelsForNote(ctx context.Context, noteID string) ([]*entities.Label, error) {
	noteUUID, err := uuid.Parse(noteID)
	if err != nil {
		return nil, err
	}

	labels, err := r.q.GetLabelsForNote(ctx, noteUUID.String())
	if err != nil {
		return nil, err
	}

	result := make([]*entities.Label, len(labels))
	for i, label := range labels {
		result[i] = &entities.Label{
			ID:        label.ID,
			UserID:    label.UserID,
			Name:      label.Name,
			Color:     label.Color,
			CreatedAt: label.CreatedAt,
			UpdatedAt: label.UpdatedAt,
		}
	}

	return result, nil
}

func (r *LabelRepositoryImpl) GetNotesForLabel(ctx context.Context, labelID string) ([]string, error) {
	labelUUID, err := uuid.Parse(labelID)
	if err != nil {
		return nil, err
	}

	noteIDs, err := r.q.GetNotesForLabel(ctx, labelUUID.String())
	if err != nil {
		return nil, err
	}

	result := make([]string, len(noteIDs))
	copy(result, noteIDs)

	return result, nil
}
