package storage

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/randallmlough/pgxscan"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

/*
интерфейс
инициализация
*/
type PGStorage struct {
	Connection *pgx.Conn
}

type User struct {
	ID       int    `db:"id"`
	Login    string `db:"login"`
	Password string `db:"password"`
}

type Order struct {
	ID         int          `db:"id"`
	UserID     int          `db:"user_id"`
	Status     string       `db:"status"`
	Accrual    int          `db:"accrual"`
	UploadedAt time.Time    `db:"uploaded_at"`
	UpdatedAt  sql.NullTime `db:"updated_at"`
}

func New(ctx context.Context, databaseURL string) Storage {
	conn, err := pgx.Connect(ctx, databaseURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	return &PGStorage{Connection: conn}
}

func (p *PGStorage) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error) {
	return p.Connection.BeginTx(ctx, txOptions)
}

func (p *PGStorage) GetUserByLogin(ctx context.Context, login string) (*User, error) {
	var user User
	sql := "SELECT id, login, password FROM users WHERE login = $1"
	row := p.Connection.QueryRow(ctx, sql, login)
	err := pgxscan.NewScanner(row).Scan(&user.ID, &user.Login, &user.Password)
	if err != nil {
		return nil, fmt.Errorf("scan rows error. err: %s", err)
	}
	return &user, nil
}

func (p *PGStorage) CreateUser(ctx context.Context, login string, password string) (*User, error) {
	var user User
	sql := "INSERT INTO users (login,password) VALUES($1,$2) RETURNING id, login, password"
	row := p.Connection.QueryRow(ctx, sql, login, password)
	err := pgxscan.NewScanner(row).Scan(&user.ID, &user.Login, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (p *PGStorage) GetOrder(ctx context.Context, orderID int) (*Order, error) {
	var order Order
	sql := "SELECT id, user_id, status, uploaded_at, updated_at FROM orders WHERE id = $1"
	row := p.Connection.QueryRow(ctx, sql, orderID)
	err := pgxscan.NewScanner(row).Scan(&order.ID, &order.UserID, &order.Status, &order.UploadedAt, &order.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("scan rows error. err: %s", err)
	}
	return &order, nil
}

func (p *PGStorage) GetOrders(ctx context.Context, userID int) ([]*Order, error) {
	var orders []*Order
	sql := "SELECT id, user_id, status, accrual, uploaded_at, updated_at FROM orders WHERE user_id = $1 ORDER BY uploaded_at DESC"
	row, err := p.Connection.Query(ctx, sql, userID)
	if err != nil {
		return nil, err
	}
	for row.Next() {
		var order Order
		err = pgxscan.NewScanner(row).Scan(&order.ID, &order.UserID, &order.Status, &order.Accrual, &order.UploadedAt, &order.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("scan rows error. err: %s", err)
		}
		orders = append(orders, &order)
	}
	return orders, nil
}

func (p *PGStorage) SaveOrder(ctx context.Context, orderID, userID int, status string) (*Order, error) {
	var order Order
	sql := "INSERT INTO orders (id,user_id,status,uploaded_at) VALUES ($1,$2,$3,$4) RETURNING id,user_id,status,uploaded_at"
	row := p.Connection.QueryRow(ctx, sql, orderID, userID, status, time.Now())
	err := pgxscan.NewScanner(row).Scan(&order.ID, &order.UserID, &order.Status, &order.UploadedAt)
	if err != nil {
		return nil, err
	}
	return &order, nil
}
