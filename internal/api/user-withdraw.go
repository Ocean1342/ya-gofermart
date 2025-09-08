package api

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
	"gofermart/internal/storage"
	"gofermart/pkg/common"
	"io"
	"net/http"
	"strconv"
	"time"
)

type WithdrawRequest struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}

func (h *Handler) UserWithdraw(w http.ResponseWriter, r *http.Request) {
	var withdrawRequest WithdrawRequest
	logger := logrus.WithFields(map[string]interface{}{
		"HANDLER":    "UserWithdraw",
		"REQUEST_ID": r.Context().Value(middleware.RequestIDKey),
	})
	ctxUser, ok := (r.Context().Value(common.CtxUser)).(*storage.User)
	if !ok || ctxUser == nil {
		logger.Errorf("could not extract user from context")
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte("could not define user"))
		if err != nil {
			logger.Errorf("could not write data to response")
		}
		return
	}
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Errorf("user id:`%d` could not read request body: %s", ctxUser.ID, err)
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte("could not read request body"))
		if err != nil {
			logger.Errorf("could not read request body")
		}
		return
	}
	err = json.Unmarshal(bytes, &withdrawRequest)
	if err != nil {
		logger.Errorf("user id:`%d` could not unmarshal body: %s", ctxUser.ID, err)
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte(" could not unmarshal body request"))
		if err != nil {
			logger.Errorf(" could not unmarshal body request")
		}
		return
	}
	intOrder, err := strconv.Atoi(withdrawRequest.Order)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		logger.Errorf("wrong order number: %s", withdrawRequest.Order)
		_, err = w.Write([]byte("wrong order number"))
		return
	}
	if intOrder <= 0 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		logger.Errorf("wrong order number: %s", withdrawRequest.Order)
		_, err = w.Write([]byte("wrong order number"))
		if err != nil {
			logger.Errorf("could not read request body")
		}
		return
	}
	ctxTx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
	tx, err := h.Storage.BeginTx(ctxTx, pgx.TxOptions{})
	defer func() {
		err = tx.Rollback(ctxTx)
		if err != nil {
			logger.Errorf("user id:`%d` could not rollback transaction. err: %v", ctxUser.ID, err)
		}
	}()
	if err != nil {
		logger.Errorf("user id:`%d` could not start transaction: %s", ctxUser.ID, err)
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte("could not withdraw"))
		if err != nil {
			logger.Errorf("could not write data to response")
		}
		return
	}
	var balance int
	err = tx.QueryRow(ctxTx, "SELECT balance FROM user_balance WHERE user_id=$1", ctxUser.ID).Scan(&balance)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		err = tx.Rollback(ctxTx)
		if err != nil {
			logger.Errorf("user id:`%d` could not rollback transaction. err: %v", ctxUser.ID, err)
		}
		logger.Errorf("user id:`%d` could not get balance: %s", ctxUser.ID, err)
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte("could not get balance"))
		if err != nil {
			logger.Errorf("could not write data to response")
		}
		return
	}
	requestedDraw := common.MoneyFloatToInt(withdrawRequest.Sum)
	logger.Infof("user id:`%d` try draw current balance: %d, withdraw:%d", ctxUser.ID, balance, requestedDraw)
	newBalance := balance - requestedDraw
	logger.Infof("user id:`%d` try sasve new balance: %d", ctxUser.ID, newBalance)
	if newBalance < 0 {
		err = tx.Rollback(ctxTx)
		if err != nil {
			logger.Errorf("user id:`%d` could not rollback transaction. err: %v", ctxUser.ID, err)
		}
		logger.Errorf("user id:`%d` could not draw. current balance: %d,withdraw:%d", ctxUser.ID, balance, requestedDraw)
		w.WriteHeader(http.StatusPaymentRequired)
		_, err = w.Write([]byte("could not withdraw"))
		if err != nil {
			logger.Errorf("could not write data to response")
		}
		return
	}
	updateBalanceSQL := `UPDATE user_balance SET balance = $1 WHERE user_id=$2`
	_, err = tx.Exec(r.Context(), updateBalanceSQL, newBalance, ctxUser.ID)
	if err != nil {
		err = tx.Rollback(ctxTx)
		if err != nil {
			logger.Errorf("user id:`%d` could not rollback transaction. err: %v", ctxUser.ID, err)
		}
		logger.Errorf("user id:`%d` could not update balance. err: %v", ctxUser.ID, err)
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte("could not update balance"))
		if err != nil {
			logger.Errorf("could not write data to response")
		}
		return
	}
	addDrawHistorySQL := `INSERT INTO user_draw_history (user_id, transaction_dt, order_id, draw) VALUES($1,$2,$3,$4)`
	logger.Infof("sql:%s userid: %d time%s orderOD:%d money:%d", addDrawHistorySQL, ctxUser.ID, time.Now(), intOrder, common.MoneyFloatToInt(withdrawRequest.Sum))
	_, err = tx.Exec(ctxTx, addDrawHistorySQL, ctxUser.ID, time.Now(), intOrder, common.MoneyFloatToInt(withdrawRequest.Sum))
	if err != nil {
		logger.Errorf("user id:`%d` could not update draw history. err:%v", ctxUser.ID, err)
		err = tx.Rollback(ctxTx)
		if err != nil {
			logger.Errorf("user id:`%d` could not rollback transaction. err: %v", ctxUser.ID, err)
		}
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte("could not update draw history"))
		if err != nil {
			logger.Errorf("could not write data to response")
		}
		return
	}
	err = tx.Commit(ctxTx)
	if err != nil {
		logger.Errorf("user id:`%d` could not commit tx. err: %v", ctxUser.ID, err)
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte("could not save data"))
		if err != nil {
			logger.Errorf("could not write data to response")
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}
