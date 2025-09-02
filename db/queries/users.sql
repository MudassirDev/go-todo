-- name: CreateUser :one
INSERT INTO users (
  id, name, username, password, created_at, updated_at
) VALUES (
  ?, ?, ?, ?, ?, ?
)
RETURNING id, name, username, created_at, updated_at;

-- name: GetUserWithUsername :one
SELECT * FROM users WHERE username = ?;
