-- +goose Up
CREATE TABLE tasks (
  id UUID PRIMARY KEY,
  userId UUID REFERENCES users(id) ON DELETE CASCADE,
  task TEXT NOT NULL,
  is_completed BOOL NOT NULL DEFAULT FALSE,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE tasks;
