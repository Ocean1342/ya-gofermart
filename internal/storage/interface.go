package storage

import (
	"context"
	"github.com/jackc/pgx/v5"
)

type Storage interface {
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
	GetUserByLogin(ctx context.Context, login string) (*User, error)
	CreateUser(ctx context.Context, login string, password string) (*User, error)
	GetOrder(ctx context.Context, orderID int) (*Order, error)
	SaveOrder(ctx context.Context, orderID int, userID int, status string) (*Order, error)
	GetOrders(ctx context.Context, userID int) ([]*Order, error)
	GetUserBalanceWithDraw(ctx context.Context, userID int) (*UserBalanceWithDraw, error)
}
