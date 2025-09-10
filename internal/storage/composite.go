package storage

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

func (p *PGStorage) UpdateOrder(ctx context.Context, externalTx pgx.Tx, orderID, accrual int, status string) error {
	updOrderSQL := `UPDATE orders SET status=$1, accrual=$2 WHERE id=$3`
	if externalTx != nil {
		_, err := externalTx.Exec(ctx, updOrderSQL, strings.ToUpper(status), accrual, orderID)
		if err != nil {
			return err
		}
	} else {
		_, err := p.Connection.Exec(ctx, updOrderSQL, strings.ToUpper(status), accrual, orderID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *PGStorage) UpdateUserBalance(ctx context.Context, externalTx pgx.Tx, userID, accrual int) error {
	var (
		tx            pgx.Tx
		ctxTx         context.Context
		err           error
		useExternalTx bool
	)
	if externalTx != nil {
		tx = externalTx
		ctxTx = ctx
		useExternalTx = true
	} else {
		ctxTx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		tx, err = p.BeginTx(ctxTx, pgx.TxOptions{})
		if err != nil {
			return err
		}
		useExternalTx = false
	}
	//get balance
	var balance int
	row := tx.QueryRow(ctxTx, "SELECT balance FROM user_balance WHERE user_id=$1", userID)
	err = row.Scan(&balance)
	if err != nil {
		if useExternalTx {
			return err
		}
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
		if useExternalTx {
			return err
		}
		errRollback := tx.Rollback(ctxTx)
		if errRollback != nil {
			logrus.Warnf("UPDATE user_balance could not rollback tx. err: %v", errRollback)
		}
		return err
	}
	if useExternalTx {
		return nil
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

func (p *PGStorage) UpdateOrderAndBalance(ctx context.Context, orderID, userID, accrual int, status string) error {
	ctxTx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	tx, err := p.BeginTx(ctxTx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	err = p.UpdateOrder(ctxTx, tx, orderID, accrual, status)
	if err != nil {
		errRollback := tx.Rollback(ctxTx)
		if errRollback != nil {
			logrus.Warnf("UPDATE orders could not rollback tx. err: %v", errRollback)
		}
		return err
	}
	err = p.UpdateUserBalance(ctxTx, tx, userID, accrual)
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
