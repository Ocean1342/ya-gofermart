-- +goose Up
CREATE table IF NOT EXISTS tokens (
     id int NOT NULL,
     user_id int NOT NULL,
     token text NOT NULL UNIQUE,
     expired_at TIMESTAMP NOT NULL,
     PRIMARY KEY(id)
    );

-- +goose Down
DROP table IF EXISTS tokens;