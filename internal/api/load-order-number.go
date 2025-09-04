package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/joeljunstrom/go-luhn"
	"github.com/sirupsen/logrus"
	orderprocessor "gofermart/internal/order-processor"
	"gofermart/internal/storage"
	"gofermart/pkg/common"
	"io"
	"net/http"
	"strconv"
)

func (h *Handler) LoadOrderNumber(w http.ResponseWriter, r *http.Request) {
	logger := logrus.WithFields(map[string]interface{}{
		"HANDLER":    "LoadOrderNumber",
		"REQUEST_ID": r.Context().Value(middleware.RequestIDKey),
	})
	var orderID int
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
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte(fmt.Sprintf("could not read body. err: %v", err)))
		if err != nil {
			logger.Errorf("could not write data to response. err:%s", err)
		}
		return
	}
	err = json.Unmarshal(bodyBytes, &orderID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte(fmt.Sprintf("could not unmarshal body. err: %v", err)))
		if err != nil {
			logger.Errorf("could not write data to response err: %s", err)
		}
		return
	}
	if orderID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte("order id must be greater than 0"))
		if err != nil {
			logger.Errorf("could not write data to response. err: %s", err)
		}
		return
	}
	logger.Infof("user id:%d login:%s trying put order id:%d", ctxUser.ID, ctxUser.Login, orderID)
	if !luhn.Valid(strconv.Itoa(orderID)) {
		logger.Error("luhn err. wrong order number. err")
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	existedOrder, err := h.Storage.GetOrder(r.Context(), orderID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		logger.Errorf("could not get order. err:%s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if existedOrder != nil && ctxUser.ID == existedOrder.UserID {
		w.WriteHeader(http.StatusOK)
		return
	}
	if existedOrder != nil && ctxUser.ID != existedOrder.UserID {
		w.WriteHeader(http.StatusConflict)
		logger.Errorf("order id:%d was already loaded by user id:%d", existedOrder.ID, existedOrder.UserID)
		_, err = w.Write([]byte(fmt.Sprintf("order id:%d was already loaded", existedOrder.ID)))
		if err != nil {
			logger.Errorf("could not write data to response. err:%s", err)
		}
		return
	}

	order, err := h.Storage.SaveOrder(r.Context(), orderID, ctxUser.ID, "NEW")
	if err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) {
			if pgerrcode.IsConnectionException(pgError.Code) {
				w.WriteHeader(http.StatusInternalServerError)
				_, err = w.Write([]byte("connection error"))
				if err != nil {
					logger.Errorf("could not write data to response. err:%s", err)
				}
				return
			}
			if pgerrcode.IsIntegrityConstraintViolation(pgError.Code) {
				w.WriteHeader(http.StatusConflict)
				_, err = w.Write([]byte("order was already loaded"))
				if err != nil {
					logger.Errorf("could not write data to response. err:%s", err)
				}
				return
			}
			if pgerrcode.IsSyntaxErrororAccessRuleViolation(pgError.Code) {
				w.WriteHeader(http.StatusInternalServerError)
				_, err = w.Write([]byte("syntax error"))
				if err != nil {
					logger.Errorf("could not write data to response")
				}
				logger.Errorf("syntax error on save order. err:%s", pgError.Error())
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			_, err = w.Write([]byte("unexpected error"))
			if err != nil {
				logger.Errorf("could not write data to response")
			}
			logger.Errorf("error on save order. err:%s", pgError.Error())
			return
		}

		logger.Infof("user id:%d login:%s saved order id:%d", ctxUser.ID, ctxUser.Login, order.ID)
	}

	h.OrderProcessor.Queue <- orderprocessor.OrderQueueItem{
		OrderID:    orderID,
		UserID:     ctxUser.ID,
		RetryTimes: 0,
	}
	w.WriteHeader(http.StatusAccepted)
}
