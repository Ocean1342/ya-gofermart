-- +goose Up
SELECT 'up SQL query';
CREATE TABLE IF NOT EXISTS user_balance_history
(
    user_id        int NOT NULL,
    transaction_dt TIMESTAMP,
    order_id       int NOT NULL,
    nominal        int
);


-- +goose Down
SELECT 'down SQL query';
DROP TABLE IF EXISTS user_balance_history;
