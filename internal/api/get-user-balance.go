package api

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
	"gofermart/internal/storage"
	"gofermart/pkg/common"
	"net/http"
)

type UserBalanceResponse struct {
	Current  float64 `json:"current"`
	Withdraw float64 `json:"withdraw"`
}

func (h *Handler) GetUserBalance(w http.ResponseWriter, r *http.Request) {
	logger := logrus.WithFields(map[string]interface{}{
		"HANDLER":    "GetUserBalance",
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
	userBalance, err := h.Storage.GetUserBalanceWithDraw(r.Context(), ctxUser.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			userBalanceWithDraw := UserBalanceResponse{
				Current:  0,
				Withdraw: 0,
			}
			bytes, err := json.Marshal(userBalanceWithDraw)
			if err != nil {
				logger.Errorf("could not marshal data")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			_, err = w.Write(bytes)
			if err != nil {
				logger.Errorf("could not write data to response")
			}
			return
		}

		logger.Errorf("user id:`%d` error get balanceWithDraw: %s", ctxUser.ID, err)
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte("could not get balanceWithDraw"))
		if err != nil {
			logger.Errorf("could not write data to response")
		}
		return
	}
	var withDraw int
	var withDrawFloat float64
	if userBalance.WithDraw.Valid {
		err = userBalance.WithDraw.Scan(withDraw)
		if err != nil {
			logger.Errorf("could not scan user withdraw")
		}
		withDrawFloat = common.MoneyIntToFloat(withDraw)
	}

	userBalanceWithDraw := UserBalanceResponse{
		Current:  common.MoneyIntToFloat(userBalance.Balance),
		Withdraw: withDrawFloat,
	}
	bytes, err := json.Marshal(userBalanceWithDraw)
	if err != nil {
		logger.Errorf("could not marshal data")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = w.Write(bytes)
	if err != nil {
		logger.Errorf("could not write data to response")
	}
}
