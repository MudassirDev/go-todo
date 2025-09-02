-- name: CreateTask :one
INSERT INTO tasks (
  id, userId, task, created_at, updated_at
) VALUES (
  ?, ?, ?, ?, ?
)
RETURNING *;
