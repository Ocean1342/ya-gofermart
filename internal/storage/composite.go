package storage

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

func (p *PGStorage) UpdateOrderAndBalance(ctx context.Context, orderID, userID, accrual int, status string) error {
	//tx
	ctxTx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	tx, err := p.BeginTx(ctxTx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	//upd order
	updOrderSQL := `UPDATE orders SET status=$1, accrual=$2 WHERE id=$3`
	_, err = tx.Exec(ctxTx, updOrderSQL, strings.ToUpper(status), accrual, orderID)
	if err != nil {
		errRollback := tx.Rollback(ctxTx)
		if errRollback != nil {
			logrus.Warnf("UPDATE orders could not rollback tx. err: %v", errRollback)
		}
		return err
	}
	//get balance
	var balance int
	row := tx.QueryRow(ctxTx, "SELECT balance FROM user_balance WHERE user_id=$1", userID)
	err = row.Scan(&balance)
	if err != nil {
		errRollback := tx.Rollback(ctxTx)
		if errRollback != nil {
			logrus.Warnf("SELECT balance could not rollback tx. err: %v", errRollback)
		}
		return err
	}
	newBalance := balance + accrual
	//upd balance
	_, err = tx.Exec(ctxTx, `UPDATE user_balance SET balance = $1 WHERE user_id=$2`, newBalance, userID)
	if err != nil {
		errRollback := tx.Rollback(ctxTx)
		if errRollback != nil {
			logrus.Warnf("UPDATE user_balance could not rollback tx. err: %v", errRollback)
		}
		return err
	}
	err = tx.Commit(ctxTx)
	if err != nil {
		errRollback := tx.Rollback(ctxTx)
		if errRollback != nil {
			logrus.Warnf("commit could not rollback tx. err: %v", errRollback)
		}
		return err
	}
	return nil
}

type UnprocessedOrder struct {
	OrderID int `db:"id"`
	UserID  int `db:"user_id"`
}

func (p *PGStorage) GetUnprocessedOrders(ctxTx context.Context, limit int) ([]*UnprocessedOrder, error) {
	orderSQL := `SELECT id, user_id FROM orders WHERE status !='PROCESSED' and status !='INVALID' LIMIT $1`
	rows, err := p.Connection.Query(ctxTx, orderSQL, limit)
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
