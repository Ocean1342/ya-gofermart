-- +goose Up
CREATE table IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    login text NOT NULL UNIQUE,
    password text NOT NULL);

-- +goose Down
DROP table IF EXISTS users;