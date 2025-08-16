-- +goose Up
CREATE table IF NOT EXISTS tokens (
     id SERIAL PRIMARY KEY,
     user_id int NOT NULL,
     token text NOT NULL UNIQUE,
     expired_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP table IF EXISTS tokens;