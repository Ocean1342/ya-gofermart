-- +goose Up
SELECT 'up SQL query';
CREATE TABLE IF NOT EXISTS user_draw_history
(
    user_id        int NOT NULL,
    transaction_dt TIMESTAMP,
    order_id       bigint NOT NULL,
    draw           int
);


-- +goose Down
SELECT 'down SQL query';
DROP TABLE IF EXISTS user_draw_history;
