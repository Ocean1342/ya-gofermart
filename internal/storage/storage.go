package storage

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/randallmlough/pgxscan"
	log "github.com/sirupsen/logrus"
	"os"
)

/*
интерфейс
инициализация
*/
type PGStorage struct {
	Connection *pgx.Conn
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

func (p *PGStorage) GetToken() {
}

type User struct {
	ID       int    `db:"id"`
	Login    string `db:"login"`
	Password string `db:"password"`
}

func (p *PGStorage) GetUserByLogin(ctx context.Context, login string) (*User, error) {
	var user User
	sql := "SELECT id, login, password FROM users WHERE login = $1"
	row := p.Connection.QueryRow(ctx, sql, login)
	if row == nil {
		return nil, fmt.Errorf("user not found")
	}
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
