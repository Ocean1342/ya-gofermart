package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
	"gofermart/internal/storage"
	"gofermart/pkg/common"
	"net/http"
	"strconv"
	"time"
)

type WithdrawRequest struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}

func (s *Service) UserWithDraw(r *http.Request, user *storage.User, withdrawRequest WithdrawRequest) error {
	logger := logrus.WithFields(map[string]interface{}{
		"HANDLER":    "UserWithdraw",
		"REQUEST_ID": r.Context().Value(middleware.RequestIDKey),
	})

	intOrder, err := strconv.Atoi(withdrawRequest.Order)
	if err != nil {
		return ErrStatusUnprocessableEntity
	}
	if intOrder <= 0 {
		logger.Errorf("wrong order number: %s", withdrawRequest.Order)
		return ErrStatusUnprocessableEntity
	}

	ctxTx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
	tx, err := s.storage.BeginTx(ctxTx, pgx.TxOptions{})
	defer func() {
		err = tx.Rollback(ctxTx)
		if err != nil {
			logger.Errorf("user id:`%d` could not rollback transaction. err: %v", user.ID, err)
		}
	}()
	if err != nil {
		return errors.Join(ErrInternalServer, fmt.Errorf("user id:`%d` could not begin transaction. err: %v", user.ID, err))
	}
	var balance int
	err = tx.QueryRow(ctxTx, "SELECT balance FROM user_balance WHERE user_id=$1", user.ID).Scan(&balance)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		err = tx.Rollback(ctxTx)
		if err != nil {
			logger.Errorf("user id:`%d` could not rollback transaction. err: %v", user.ID, err)
		}
		logger.Errorf("user id:`%d` could not get balance: %s", user.ID, err)
		return errors.Join(ErrInternalServer, err)
	}
	requestedDraw := common.MoneyFloatToInt(withdrawRequest.Sum)
	logger.Infof("user id:`%d` try draw current balance: %d, withdraw:%d", user.ID, balance, requestedDraw)
	newBalance := balance - requestedDraw
	logger.Infof("user id:`%d` try sasve new balance: %d", user.ID, newBalance)
	if newBalance < 0 {
		err = tx.Rollback(ctxTx)
		if err != nil {
			logger.Errorf("user id:`%d` could not rollback transaction. err: %v", user.ID, err)
		}
		logger.Errorf("user id:`%d` could not draw. current balance: %d,withdraw:%d", user.ID, balance, requestedDraw)
		return errors.Join(ErrPaymentRequired, err)
	}
	updateBalanceSQL := `UPDATE user_balance SET balance = $1 WHERE user_id=$2`
	_, err = tx.Exec(r.Context(), updateBalanceSQL, newBalance, user.ID)
	if err != nil {
		err = tx.Rollback(ctxTx)
		if err != nil {
			logger.Errorf("user id:`%d` could not rollback transaction. err: %v", user.ID, err)
		}
		logger.Errorf("user id:`%d` could not update balance. err: %v", user.ID, err)
		return errors.Join(ErrInternalServer, err)
	}
	addDrawHistorySQL := `INSERT INTO user_draw_history (user_id, transaction_dt, order_id, draw) VALUES($1,$2,$3,$4)`
	logger.Infof("sql:%s userid: %d time%s orderOD:%d money:%d", addDrawHistorySQL, user.ID, time.Now(), intOrder, common.MoneyFloatToInt(withdrawRequest.Sum))
	_, err = tx.Exec(ctxTx, addDrawHistorySQL, user.ID, time.Now(), intOrder, common.MoneyFloatToInt(withdrawRequest.Sum))
	if err != nil {
		logger.Errorf("user id:`%d` could not update draw history. err:%v", user.ID, err)
		err = tx.Rollback(ctxTx)
		if err != nil {
			logger.Errorf("user id:`%d` could not rollback transaction. err: %v", user.ID, err)
		}
		return errors.Join(ErrInternalServer, err)
	}
	err = tx.Commit(ctxTx)
	if err != nil {
		logger.Errorf("user id:`%d` could not commit tx. err: %v", user.ID, err)
		return errors.Join(ErrInternalServer, err)
	}
	return nil
}
