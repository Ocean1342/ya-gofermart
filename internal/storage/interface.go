package storage

import (
	"context"
	"github.com/jackc/pgx/v5"
)

type Storage interface {
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
	GetToken()
	GetUserByLogin(ctx context.Context, login string) (*User, error)
	CreateUser(ctx context.Context, login string, password string) (*User, error)
}
