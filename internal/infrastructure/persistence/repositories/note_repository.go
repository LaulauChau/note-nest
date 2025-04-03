package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/LaulauChau/note-nest/internal/domain/entities"
	"github.com/LaulauChau/note-nest/internal/domain/repositories"
	"github.com/jackc/pgx/v5/pgtype"
)

type NoteRepositoryImpl struct {
	q *Queries
}

func NewNoteRepository(q *Queries) repositories.NoteRepository {
	return &NoteRepositoryImpl{q: q}
}

func (r *NoteRepositoryImpl) Create(ctx context.Context, note *entities.Note) error {
	// Parse the user ID (which should be a UUID)
	userID, err := uuid.Parse(note.UserID)
	if err != nil {
		return err
	}

	// Parse the note ID
	noteID, err := uuid.Parse(note.ID)
	if err != nil {
		return err
	}

	params := CreateNoteParams{
		ID:         noteID.String(),
		UserID:     userID.String(),
		Title:      note.Title,
		Content:    note.Content,
		IsArchived: note.IsArchived,
		Label:      pgtype.Text{String: note.Label, Valid: true},
		CreatedAt:  note.CreatedAt,
		UpdatedAt:  note.UpdatedAt,
	}

	_, err = r.q.CreateNote(ctx, params)
	return err
}

func (r *NoteRepositoryImpl) GetByID(ctx context.Context, id string) (*entities.Note, error) {
	noteID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	note, err := r.q.GetNoteByID(ctx, noteID.String())
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, nil
		}
		return nil, err
	}

	return &entities.Note{
		ID:         note.ID,
		UserID:     note.UserID,
		Title:      note.Title,
		Content:    note.Content,
		IsArchived: note.IsArchived,
		Label:      note.Label.String,
		CreatedAt:  note.CreatedAt,
		UpdatedAt:  note.UpdatedAt,
	}, nil
}

func (r *NoteRepositoryImpl) GetByUserID(ctx context.Context, userID string) ([]*entities.Note, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	notes, err := r.q.GetNotesByUserID(ctx, userUUID.String())
	if err != nil {
		return nil, err
	}

	result := make([]*entities.Note, len(notes))
	for i, note := range notes {
		result[i] = &entities.Note{
			ID:         note.ID,
			UserID:     note.UserID,
			Title:      note.Title,
			Content:    note.Content,
			IsArchived: note.IsArchived,
			Label:      note.Label.String,
			CreatedAt:  note.CreatedAt,
			UpdatedAt:  note.UpdatedAt,
		}
	}

	return result, nil
}

func (r *NoteRepositoryImpl) GetArchivedByUserID(ctx context.Context, userID string) ([]*entities.Note, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	notes, err := r.q.GetArchivedNotesByUserID(ctx, userUUID.String())
	if err != nil {
		return nil, err
	}

	result := make([]*entities.Note, len(notes))
	for i, note := range notes {
		result[i] = &entities.Note{
			ID:         note.ID,
			UserID:     note.UserID,
			Title:      note.Title,
			Content:    note.Content,
			IsArchived: note.IsArchived,
			Label:      note.Label.String,
			CreatedAt:  note.CreatedAt,
			UpdatedAt:  note.UpdatedAt,
		}
	}

	return result, nil
}

func (r *NoteRepositoryImpl) Update(ctx context.Context, note *entities.Note) error {
	noteID, err := uuid.Parse(note.ID)
	if err != nil {
		return err
	}

	params := UpdateNoteParams{
		ID:         noteID.String(),
		Title:      note.Title,
		Content:    note.Content,
		IsArchived: note.IsArchived,
		Label:      pgtype.Text{String: note.Label, Valid: true},
		UpdatedAt:  time.Now(),
	}

	return r.q.UpdateNote(ctx, params)
}

func (r *NoteRepositoryImpl) Delete(ctx context.Context, id string) error {
	noteID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	return r.q.DeleteNote(ctx, noteID.String())
}
