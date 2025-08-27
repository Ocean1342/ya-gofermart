-- +goose Up
CREATE table IF NOT EXISTS orders
(
    id          int UNIQUE NOT NULL,
    user_id     int        NOT NULL,
    status      text       NOT NULL,
    accrual     int        NOT NULL DEFAULT 0,
    uploaded_at TIMESTAMP  NOT NULL,
    updated_at  TIMESTAMP
);

-- +goose Down
DROP table IF EXISTS orders;