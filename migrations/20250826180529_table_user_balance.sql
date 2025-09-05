-- +goose Up
SELECT 'up SQL query';
CREATE TABLE IF NOT EXISTS user_balance
(
    user_id    int UNIQUE NOT NULL,
    balance    int NOT NULL DEFAULT 0 CHECK (balance >= 0),
    updated_at TIMESTAMP
);

-- +goose Down
SELECT 'down SQL query';
DROP TABLE IF EXISTS user_balance;