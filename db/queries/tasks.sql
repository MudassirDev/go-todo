-- name: CreateTask :one
INSERT INTO tasks (
  id, user_id, task, created_at, updated_at
) VALUES (
  ?, ?, ?, ?, ?
)
RETURNING *;

-- name: DeleteTaskWithID :exec
DELETE FROM tasks WHERE id = ? AND user_id = ?;

-- name: UpdateTaskWithID :one
UPDATE tasks SET is_completed = ? WHERE id = ? AND user_id = ? RETURNING *;

-- name: GetTasksWithUserID :many
SELECT * FROM tasks WHERE user_id = ?;
