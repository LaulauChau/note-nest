// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: queries.sql

package repositories

import (
	"context"
	"time"
)

const addLabelToNote = `-- name: AddLabelToNote :exec
INSERT INTO note_labels (note_id, label_id) VALUES ($1, $2) ON CONFLICT DO NOTHING
`

type AddLabelToNoteParams struct {
	NoteID  string `json:"note_id"`
	LabelID string `json:"label_id"`
}

func (q *Queries) AddLabelToNote(ctx context.Context, arg AddLabelToNoteParams) error {
	_, err := q.db.Exec(ctx, addLabelToNote, arg.NoteID, arg.LabelID)
	return err
}

const createLabel = `-- name: CreateLabel :one
INSERT INTO labels (id, user_id, name, color, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, user_id, name, color, created_at, updated_at
`

type CreateLabelParams struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Name      string    `json:"name"`
	Color     string    `json:"color"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (q *Queries) CreateLabel(ctx context.Context, arg CreateLabelParams) (Label, error) {
	row := q.db.QueryRow(ctx, createLabel,
		arg.ID,
		arg.UserID,
		arg.Name,
		arg.Color,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var i Label
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Name,
		&i.Color,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const createNote = `-- name: CreateNote :one
INSERT INTO notes (id, user_id, title, content, is_archived, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, user_id, title, content, is_archived, created_at, updated_at
`

type CreateNoteParams struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	IsArchived bool      `json:"is_archived"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (q *Queries) CreateNote(ctx context.Context, arg CreateNoteParams) (Note, error) {
	row := q.db.QueryRow(ctx, createNote,
		arg.ID,
		arg.UserID,
		arg.Title,
		arg.Content,
		arg.IsArchived,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var i Note
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Title,
		&i.Content,
		&i.IsArchived,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const createSession = `-- name: CreateSession :one
INSERT INTO sessions (user_id, expires_at)
VALUES ($1, $2)
RETURNING id, user_id, expires_at, created_at
`

type CreateSessionParams struct {
	UserID    string    `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
}

func (q *Queries) CreateSession(ctx context.Context, arg CreateSessionParams) (Session, error) {
	row := q.db.QueryRow(ctx, createSession, arg.UserID, arg.ExpiresAt)
	var i Session
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.ExpiresAt,
		&i.CreatedAt,
	)
	return i, err
}

const createUser = `-- name: CreateUser :one
INSERT INTO users (email, name, password)
VALUES ($1, $2, $3)
RETURNING id, email, name, password, created_at, updated_at
`

type CreateUserParams struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, createUser, arg.Email, arg.Name, arg.Password)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Name,
		&i.Password,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteAllSessionsByUserID = `-- name: DeleteAllSessionsByUserID :exec
DELETE FROM sessions WHERE user_id = $1
`

func (q *Queries) DeleteAllSessionsByUserID(ctx context.Context, userID string) error {
	_, err := q.db.Exec(ctx, deleteAllSessionsByUserID, userID)
	return err
}

const deleteLabel = `-- name: DeleteLabel :exec
DELETE FROM labels WHERE id = $1
`

func (q *Queries) DeleteLabel(ctx context.Context, id string) error {
	_, err := q.db.Exec(ctx, deleteLabel, id)
	return err
}

const deleteNote = `-- name: DeleteNote :exec
DELETE FROM notes WHERE id = $1
`

func (q *Queries) DeleteNote(ctx context.Context, id string) error {
	_, err := q.db.Exec(ctx, deleteNote, id)
	return err
}

const deleteSession = `-- name: DeleteSession :exec
DELETE FROM sessions WHERE id = $1
`

func (q *Queries) DeleteSession(ctx context.Context, id string) error {
	_, err := q.db.Exec(ctx, deleteSession, id)
	return err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1
`

func (q *Queries) DeleteUser(ctx context.Context, id string) error {
	_, err := q.db.Exec(ctx, deleteUser, id)
	return err
}

const getArchivedNotesByUserID = `-- name: GetArchivedNotesByUserID :many
SELECT id, user_id, title, content, is_archived, created_at, updated_at FROM notes WHERE user_id = $1 AND is_archived = true ORDER BY updated_at DESC
`

func (q *Queries) GetArchivedNotesByUserID(ctx context.Context, userID string) ([]Note, error) {
	rows, err := q.db.Query(ctx, getArchivedNotesByUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Note
	for rows.Next() {
		var i Note
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Title,
			&i.Content,
			&i.IsArchived,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getLabelByID = `-- name: GetLabelByID :one
SELECT id, user_id, name, color, created_at, updated_at FROM labels WHERE id = $1
`

func (q *Queries) GetLabelByID(ctx context.Context, id string) (Label, error) {
	row := q.db.QueryRow(ctx, getLabelByID, id)
	var i Label
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Name,
		&i.Color,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getLabelByName = `-- name: GetLabelByName :one
SELECT id, user_id, name, color, created_at, updated_at FROM labels WHERE user_id = $1 AND name = $2
`

type GetLabelByNameParams struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
}

func (q *Queries) GetLabelByName(ctx context.Context, arg GetLabelByNameParams) (Label, error) {
	row := q.db.QueryRow(ctx, getLabelByName, arg.UserID, arg.Name)
	var i Label
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Name,
		&i.Color,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getLabelsByUserID = `-- name: GetLabelsByUserID :many
SELECT id, user_id, name, color, created_at, updated_at FROM labels WHERE user_id = $1 ORDER BY name
`

func (q *Queries) GetLabelsByUserID(ctx context.Context, userID string) ([]Label, error) {
	rows, err := q.db.Query(ctx, getLabelsByUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Label
	for rows.Next() {
		var i Label
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Name,
			&i.Color,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getLabelsForNote = `-- name: GetLabelsForNote :many
SELECT l.id, l.user_id, l.name, l.color, l.created_at, l.updated_at FROM labels l
JOIN note_labels nl ON l.id = nl.label_id
WHERE nl.note_id = $1
ORDER BY l.name
`

func (q *Queries) GetLabelsForNote(ctx context.Context, noteID string) ([]Label, error) {
	rows, err := q.db.Query(ctx, getLabelsForNote, noteID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Label
	for rows.Next() {
		var i Label
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Name,
			&i.Color,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getNoteByID = `-- name: GetNoteByID :one
SELECT id, user_id, title, content, is_archived, created_at, updated_at FROM notes WHERE id = $1
`

func (q *Queries) GetNoteByID(ctx context.Context, id string) (Note, error) {
	row := q.db.QueryRow(ctx, getNoteByID, id)
	var i Note
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Title,
		&i.Content,
		&i.IsArchived,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getNotesByUserID = `-- name: GetNotesByUserID :many
SELECT id, user_id, title, content, is_archived, created_at, updated_at FROM notes WHERE user_id = $1 AND is_archived = false ORDER BY updated_at DESC
`

func (q *Queries) GetNotesByUserID(ctx context.Context, userID string) ([]Note, error) {
	rows, err := q.db.Query(ctx, getNotesByUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Note
	for rows.Next() {
		var i Note
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Title,
			&i.Content,
			&i.IsArchived,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getNotesForLabel = `-- name: GetNotesForLabel :many
SELECT note_id FROM note_labels WHERE label_id = $1
`

func (q *Queries) GetNotesForLabel(ctx context.Context, labelID string) ([]string, error) {
	rows, err := q.db.Query(ctx, getNotesForLabel, labelID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var note_id string
		if err := rows.Scan(&note_id); err != nil {
			return nil, err
		}
		items = append(items, note_id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getSessionByID = `-- name: GetSessionByID :one
SELECT id, user_id, expires_at, created_at FROM sessions WHERE id = $1
`

func (q *Queries) GetSessionByID(ctx context.Context, id string) (Session, error) {
	row := q.db.QueryRow(ctx, getSessionByID, id)
	var i Session
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.ExpiresAt,
		&i.CreatedAt,
	)
	return i, err
}

const getSessionWithUser = `-- name: GetSessionWithUser :one
SELECT
    sessions.id AS session_id,
    sessions.user_id AS session_user_id,
    sessions.expires_at AS session_expires_at,
    sessions.created_at AS session_created_at,
    users.id AS user_id,
    users.email AS user_email,
    users.name AS user_name
FROM sessions
INNER JOIN users ON sessions.user_id = users.id
WHERE sessions.id = $1
`

type GetSessionWithUserRow struct {
	SessionID        string    `json:"session_id"`
	SessionUserID    string    `json:"session_user_id"`
	SessionExpiresAt time.Time `json:"session_expires_at"`
	SessionCreatedAt time.Time `json:"session_created_at"`
	UserID           string    `json:"user_id"`
	UserEmail        string    `json:"user_email"`
	UserName         string    `json:"user_name"`
}

func (q *Queries) GetSessionWithUser(ctx context.Context, id string) (GetSessionWithUserRow, error) {
	row := q.db.QueryRow(ctx, getSessionWithUser, id)
	var i GetSessionWithUserRow
	err := row.Scan(
		&i.SessionID,
		&i.SessionUserID,
		&i.SessionExpiresAt,
		&i.SessionCreatedAt,
		&i.UserID,
		&i.UserEmail,
		&i.UserName,
	)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, email, name, password, created_at, updated_at FROM users WHERE email = $1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRow(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Name,
		&i.Password,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUserByID = `-- name: GetUserByID :one
SELECT id, email, name, password, created_at, updated_at FROM users WHERE id = $1
`

func (q *Queries) GetUserByID(ctx context.Context, id string) (User, error) {
	row := q.db.QueryRow(ctx, getUserByID, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Name,
		&i.Password,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const removeLabelFromNote = `-- name: RemoveLabelFromNote :exec
DELETE FROM note_labels WHERE note_id = $1 AND label_id = $2
`

type RemoveLabelFromNoteParams struct {
	NoteID  string `json:"note_id"`
	LabelID string `json:"label_id"`
}

func (q *Queries) RemoveLabelFromNote(ctx context.Context, arg RemoveLabelFromNoteParams) error {
	_, err := q.db.Exec(ctx, removeLabelFromNote, arg.NoteID, arg.LabelID)
	return err
}

const updateLabel = `-- name: UpdateLabel :exec
UPDATE labels SET name = $2, color = $3, updated_at = $4 WHERE id = $1
`

type UpdateLabelParams struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Color     string    `json:"color"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (q *Queries) UpdateLabel(ctx context.Context, arg UpdateLabelParams) error {
	_, err := q.db.Exec(ctx, updateLabel,
		arg.ID,
		arg.Name,
		arg.Color,
		arg.UpdatedAt,
	)
	return err
}

const updateNote = `-- name: UpdateNote :exec
UPDATE notes SET title = $2, content = $3, is_archived = $4, updated_at = $5 WHERE id = $1
`

type UpdateNoteParams struct {
	ID         string    `json:"id"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	IsArchived bool      `json:"is_archived"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (q *Queries) UpdateNote(ctx context.Context, arg UpdateNoteParams) error {
	_, err := q.db.Exec(ctx, updateNote,
		arg.ID,
		arg.Title,
		arg.Content,
		arg.IsArchived,
		arg.UpdatedAt,
	)
	return err
}

const updateSessionExpiresAt = `-- name: UpdateSessionExpiresAt :exec
UPDATE sessions SET expires_at = $2 WHERE id = $1
`

type UpdateSessionExpiresAtParams struct {
	ID        string    `json:"id"`
	ExpiresAt time.Time `json:"expires_at"`
}

func (q *Queries) UpdateSessionExpiresAt(ctx context.Context, arg UpdateSessionExpiresAtParams) error {
	_, err := q.db.Exec(ctx, updateSessionExpiresAt, arg.ID, arg.ExpiresAt)
	return err
}

const updateUser = `-- name: UpdateUser :exec
UPDATE users SET email = $2, name = $3, password = $4 WHERE id = $1
`

type UpdateUserParams struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) error {
	_, err := q.db.Exec(ctx, updateUser,
		arg.ID,
		arg.Email,
		arg.Name,
		arg.Password,
	)
	return err
}
