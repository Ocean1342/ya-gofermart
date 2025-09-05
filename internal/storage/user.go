package storage

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/randallmlough/pgxscan"
	"time"
)

type User struct {
	ID       int    `db:"id"`
	Login    string `db:"login"`
	Password string `db:"password"`
}

type UserBalanceWithDraw struct {
	UserID   int           `db:"user_id"`
	Balance  int           `db:"balance"`
	WithDraw sql.NullInt64 `db:"draw"`
}

type UserWithDraw struct {
	UserID        int           `db:"user_id"`
	WithDraw      sql.NullInt64 `db:"draw"`
	OrderID       int           `db:"order_id"`
	TransactionDT time.Time     `db:"transaction_dt"`
}

func (p *PGStorage) GetUserByLogin(ctx context.Context, login string) (*User, error) {
	var user User
	sqlQuery := "SELECT id, login, password FROM users WHERE login = $1"
	row := p.Connection.QueryRow(ctx, sqlQuery, login)
	err := pgxscan.NewScanner(row).Scan(&user.ID, &user.Login, &user.Password)
	if err != nil {
		return nil, fmt.Errorf("scan rows error. err: %v", err)
	}
	return &user, nil
}

func (p *PGStorage) CreateUser(ctx context.Context, login string, password string) (*User, error) {
	var user User
	sqlQuery := `INSERT INTO users (login,password) VALUES($1,$2) RETURNING id, login, password`
	row := p.Connection.QueryRow(ctx, sqlQuery, login, password)
	err := pgxscan.NewScanner(row).Scan(&user.ID, &user.Login, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (p *PGStorage) CreateUserBalance(ctx context.Context, userID int) error {
	sqlQuery := `INSERT INTO user_balance (user_id, balance, updated_at) VALUES($1,$2,$3)`
	_, err := p.Connection.Exec(ctx, sqlQuery, userID, 0, time.Now())
	if err != nil {
		return err
	}
	return nil
}

func (p *PGStorage) GetUserBalanceWithDraw(ctx context.Context, userID int) (*UserBalanceWithDraw, error) {
	var user UserBalanceWithDraw
	sqlQuery := `SELECT  ub.balance, SUM(udh.draw) as draw
			FROM user_balance AS ub
         	LEFT JOIN user_draw_history AS udh 
         	ON ub.user_id = udh.user_id
			WHERE ub.user_id = $1
			GROUP BY ub.balance`
	row := p.Connection.QueryRow(ctx, sqlQuery, userID)
	err := pgxscan.NewScanner(row).Scan(&user.Balance, &user.WithDraw)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (p *PGStorage) GetUserDrawHistory(ctx context.Context, userID int) ([]*UserWithDraw, error) {
	var res []*UserWithDraw
	sqlQuery := `SELECT user_id, transaction_dt, order_id, draw 
				 FROM user_draw_history
				 WHERE user_id=$1`
	rows, err := p.Connection.Query(ctx, sqlQuery, userID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var el UserWithDraw
		err = rows.Scan(&el.UserID, &el.TransactionDT, &el.OrderID, &el.WithDraw)
		if err != nil {
			return nil, err
		}
		res = append(res, &el)
	}
	return res, nil
}
