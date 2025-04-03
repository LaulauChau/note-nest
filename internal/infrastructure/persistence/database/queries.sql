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