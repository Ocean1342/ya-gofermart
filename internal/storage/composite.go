package storage

import (
	"context"
	"github.com/jackc/pgx/v5"
)

func (p *PGStorage) UpdateOrderAndBalance(ctx context.Context, orderID, userID, accrual int, status string) error {
	//tx
	tx, err := p.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	//upd order
	updOrderSQL := `UPDATE orders SET status=$1, accrual=$2 WHERE id=$3`
	_, err = tx.Exec(ctx, updOrderSQL, status, accrual, orderID)
	if err != nil {
		return err
	}
	//get balance
	var balance int
	row := tx.QueryRow(ctx, "SELECT balance FROM user_balance WHERE user_id=$1", userID)
	err = row.Scan(&balance)
	if err != nil {
		return err
	}
	newBalance := balance + accrual
	//upd balance
	_, err = tx.Exec(ctx, `UPDATE user_balance SET balance = $1 WHERE user_id=$2`, newBalance, userID)
	if err != nil {
		return err
	}
	return nil
}
