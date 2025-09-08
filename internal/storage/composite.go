package storage

import (
	"context"
	"github.com/jackc/pgx/v5"
	"strings"
)

func (p *PGStorage) UpdateOrderAndBalance(ctx context.Context, orderID, userID, accrual int, status string) error {
	//tx
	tx, err := p.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	//upd order
	updOrderSQL := `UPDATE orders SET status=$1, accrual=$2 WHERE id=$3`
	_, err = tx.Exec(ctx, updOrderSQL, strings.ToUpper(status), accrual, orderID)
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}
	//get balance
	var balance int
	row := tx.QueryRow(ctx, "SELECT balance FROM user_balance WHERE user_id=$1", userID)
	err = row.Scan(&balance)
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}
	newBalance := balance + accrual
	//upd balance
	_, err = tx.Exec(ctx, `UPDATE user_balance SET balance = $1 WHERE user_id=$2`, newBalance, userID)
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}
	err = tx.Commit(ctx)
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}
	return nil
}

type UnprocessedOrder struct {
	OrderID int `db:"id"`
	UserID  int `db:"user_id"`
}

func (p *PGStorage) GetUnprocessedOrders(ctx context.Context, limit int) ([]*UnprocessedOrder, error) {
	orderSQL := `SELECT id, user_id FROM orders WHERE status !='PROCESSED' and status !='INVALID' LIMIT $1`
	rows, err := p.Connection.Query(ctx, orderSQL, limit)
	if err != nil {
		return nil, err
	}

	var res []*UnprocessedOrder
	for rows.Next() {
		var upOrder UnprocessedOrder
		err = rows.Scan(&upOrder.OrderID, &upOrder.UserID)
		if err != nil {
			return nil, err
		}
		res = append(res, &upOrder)
	}
	return res, nil
}
