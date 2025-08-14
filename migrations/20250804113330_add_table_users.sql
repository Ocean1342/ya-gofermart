-- +goose Up
CREATE table IF NOT EXISTS users (
    id int NOT NULL,
    login text NOT NULL UNIQUE,
    password text NOT NULL,
    PRIMARY KEY(id)
);

-- +goose Down
DROP table IF EXISTS users;