package storage

import (
	"context"
	"database/sql"
	"github.com/jackc/pgx/v5"
	"github.com/randallmlough/pgxscan"
	"strings"
	"time"
)

type Order struct {
	ID         int          `db:"id"`
	UserID     int          `db:"user_id"`
	Status     string       `db:"status"`
	Accrual    int          `db:"accrual"`
	UploadedAt time.Time    `db:"uploaded_at"`
	UpdatedAt  sql.NullTime `db:"updated_at"`
}

func (p *PGStorage) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error) {
	return p.Connection.BeginTx(ctx, txOptions)
}

func (p *PGStorage) GetOrder(ctx context.Context, orderID int) (*Order, error) {
	var order Order
	sqlQuery := "SELECT id, user_id, status, uploaded_at, updated_at FROM orders WHERE id = $1"
	row := p.Connection.QueryRow(ctx, sqlQuery, orderID)
	err := pgxscan.NewScanner(row).Scan(&order.ID, &order.UserID, &order.Status, &order.UploadedAt, &order.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (p *PGStorage) GetOrders(ctx context.Context, userID int) ([]*Order, error) {
	var orders []*Order
	sqlQuery := "SELECT id, user_id, status, accrual, uploaded_at, updated_at FROM orders WHERE user_id = $1 ORDER BY uploaded_at DESC"
	row, err := p.Connection.Query(ctx, sqlQuery, userID)
	if err != nil {
		return nil, err
	}
	for row.Next() {
		var order Order
		err = pgxscan.NewScanner(row).Scan(&order.ID, &order.UserID, &order.Status, &order.Accrual, &order.UploadedAt, &order.UpdatedAt)
		if err != nil {
			return nil, err
		}
		orders = append(orders, &order)
	}
	return orders, nil
}

func (p *PGStorage) SaveOrder(ctx context.Context, orderID, userID int, status string) (*Order, error) {
	var order Order
	sqlQuery := "INSERT INTO orders (id,user_id,status,uploaded_at) VALUES ($1,$2,$3,$4) RETURNING id,user_id,status,uploaded_at"
	row := p.Connection.QueryRow(ctx, sqlQuery, orderID, userID, strings.ToUpper(status), time.Now())
	err := pgxscan.NewScanner(row).Scan(&order.ID, &order.UserID, &order.Status, &order.UploadedAt)
	if err != nil {
		return nil, err
	}
	return &order, nil
}
