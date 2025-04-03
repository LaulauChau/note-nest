-- name: CreateUser :one
INSERT INTO users (email, name, password)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: UpdateUser :exec
UPDATE users SET email = $2, name = $3, password = $4 WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;

-- name: CreateSession :one
INSERT INTO sessions (user_id, expires_at)
VALUES ($1, $2)
RETURNING *;

-- name: GetSessionByID :one
SELECT * FROM sessions WHERE id = $1;

-- name: GetSessionWithUser :one
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
WHERE sessions.id = $1;

-- name: UpdateSessionExpiresAt :exec
UPDATE sessions SET expires_at = $2 WHERE id = $1;

-- name: DeleteSession :exec
DELETE FROM sessions WHERE id = $1;

-- name: DeleteAllSessionsByUserID :exec
DELETE FROM sessions WHERE user_id = $1;

-- name: CreateNote :one
INSERT INTO notes (id, user_id, title, content, is_archived, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetNoteByID :one
SELECT * FROM notes WHERE id = $1;

-- name: GetNotesByUserID :many
SELECT * FROM notes WHERE user_id = $1 AND is_archived = false ORDER BY updated_at DESC;

-- name: GetArchivedNotesByUserID :many
SELECT * FROM notes WHERE user_id = $1 AND is_archived = true ORDER BY updated_at DESC;

-- name: UpdateNote :exec
UPDATE notes SET title = $2, content = $3, is_archived = $4, updated_at = $5 WHERE id = $1;

-- name: DeleteNote :exec
DELETE FROM notes WHERE id = $1;

-- name: CreateLabel :one
INSERT INTO labels (id, user_id, name, color, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetLabelByID :one
SELECT * FROM labels WHERE id = $1;

-- name: GetLabelsByUserID :many
SELECT * FROM labels WHERE user_id = $1 ORDER BY name;

-- name: GetLabelByName :one
SELECT * FROM labels WHERE user_id = $1 AND name = $2;

-- name: UpdateLabel :exec
UPDATE labels SET name = $2, color = $3, updated_at = $4 WHERE id = $1;

-- name: DeleteLabel :exec
DELETE FROM labels WHERE id = $1;

-- name: AddLabelToNote :exec
INSERT INTO note_labels (note_id, label_id) VALUES ($1, $2) ON CONFLICT DO NOTHING;

-- name: RemoveLabelFromNote :exec
DELETE FROM note_labels WHERE note_id = $1 AND label_id = $2;

-- name: GetLabelsForNote :many
SELECT l.* FROM labels l
JOIN note_labels nl ON l.id = nl.label_id
WHERE nl.note_id = $1
ORDER BY l.name;

-- name: GetNotesForLabel :many
SELECT note_id FROM note_labels WHERE label_id = $1;