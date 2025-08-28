package api

import (
	"encoding/json"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
	"gofermart/internal/storage"
	"gofermart/pkg/common"
	"io"
	"net/http"
	"time"
)

type WithdrawRequest struct {
	Order int     `json:"order"`
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
	if withdrawRequest.Order <= 0 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		logger.Errorf("wrong order number: %d", withdrawRequest.Order)
		_, err = w.Write([]byte("wrong order number"))
		if err != nil {
			logger.Errorf("could not read request body")
		}
		return
	}

	//получить баланс в транзакции, если хватает, то списать, если не хватает, то откатиться
	tx, err := h.Storage.BeginTx(r.Context(), pgx.TxOptions{})
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
	err = tx.QueryRow(r.Context(), "SELECT balance FROM user_balance WHERE user_id=$1", ctxUser.ID).Scan(&balance)
	if err != nil {
		err = tx.Rollback(r.Context())
		if err != nil {
			logger.Errorf("user id:`%d` could not rollback transaction. err:%s", ctxUser.ID, err)
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
	logger.Errorf("user id:`%d` try draw current balance: %d,withdraw:%d", ctxUser.ID, balance, requestedDraw)
	newBalance := balance - requestedDraw
	if newBalance < 0 {
		err = tx.Rollback(r.Context())
		if err != nil {
			logger.Errorf("user id:`%d` could not rollback transaction. err:%s", ctxUser.ID, err)
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
		err = tx.Rollback(r.Context())
		if err != nil {
			logger.Errorf("user id:`%d` could not rollback transaction. err:%s", ctxUser.ID, err)
		}
		logger.Errorf("user id:`%d` could not update balance. err:%s", ctxUser.ID, err)
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte("could not update balance"))
		if err != nil {
			logger.Errorf("could not write data to response")
		}
		return
	}
	addDrawHistorySQL := `INSERT INTO user_draw_history (user_id, transaction_dt, order_id, draw) VALUES($1,$2,$3,$4)`
	_, err = tx.Exec(r.Context(), addDrawHistorySQL, ctxUser.ID, time.Now(), withdrawRequest.Order, common.MoneyFloatToInt(withdrawRequest.Sum))
	if err != nil {
		err = tx.Rollback(r.Context())
		if err != nil {
			logger.Errorf("user id:`%d` could not rollback transaction. err:%s", ctxUser.ID, err)
		}
		logger.Errorf("user id:`%d` could not update draw history. err:%s", ctxUser.ID, err)
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte("could not update draw history"))
		if err != nil {
			logger.Errorf("could not write data to response")
		}
		return
	}
	err = tx.Commit(r.Context())
	if err != nil {
		logger.Errorf("user id:`%d` could not commit tx. err:%s", ctxUser.ID, err)
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte("could not save data"))
		if err != nil {
			logger.Errorf("could not write data to response")
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}
