package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/sirupsen/logrus"
	"gofermart/internal/storage"
	"gofermart/pkg/common"
	"io"
	"net/http"
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
		w.Write([]byte(fmt.Sprintf("could not define user")))
		return
	}
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not read body. err: %v", err)))
		return
	}
	err = json.Unmarshal(bodyBytes, &orderID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("could not unmarshal body. err: %v", err)))
		return
	}
	if orderID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("order id must be greater than 0")))
		return
	}
	logger.Infof("user id:`%d` login:`%s` trying put order id:`%d`", ctxUser.ID, ctxUser.Login, orderID)

	//попытаться воткнуть заказ, помапить ошибку и вернуть корректный статус
	order, err := h.Storage.SaveOrder(r.Context(), ctxUser.ID, orderID, "NEW")
	if err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) {
			if pgerrcode.IsConnectionException(pgError.Code) {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("connection error")))
				return
			}
			if pgerrcode.IsIntegrityConstraintViolation(pgError.Code) {
				w.WriteHeader(http.StatusConflict)
				w.Write([]byte(fmt.Sprintf("order was already loaded")))
				return
			}
			if pgerrcode.IsSyntaxErrororAccessRuleViolation(pgError.Code) {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("syntax error")))
				logger.Errorf("syntax error on save order. err:`%s`", pgError.Error())
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("unexpected error")))
			logger.Errorf("error on save order. err:`%s`", pgError.Error())
			return
		}

		logger.Infof("user id:`%d` login:`%s` saved order id:`%d`", ctxUser.ID, ctxUser.Login, order.ID)
	}
}
