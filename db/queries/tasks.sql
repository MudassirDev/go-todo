-- name: CreateTask :one
INSERT INTO tasks (
  id, user_id, task, created_at, updated_at
) VALUES (
  ?, ?, ?, ?, ?
)
RETURNING *;

-- name: DeleteTaskWithID :exec
DELETE FROM tasks WHERE id = ? AND user_id = ?;
